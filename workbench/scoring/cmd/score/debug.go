// The debug MCP: a view into the running process for an agent. Off unless
// DEBUG_MCP_ADDR is set, and then only on a loopback address: inside the pod
// nothing reaches it but `kubectl port-forward`, so the cluster's RBAC decides
// who gets in. The agent's read-only kubeconfig cannot port-forward; you open
// the tunnel, the agent uses it.

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"runtime"
	"slices"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel/trace"

	"workshop/scoring/internal/event"
	"workshop/scoring/internal/scoring"
)

const (
	keep     = 500  // requests kept for recent_requests
	maxItems = 2000 // largest event replay scores
)

// request is what recent_requests keeps of one event: sizes and timings,
// never the event itself, so no card or email sits in this process's memory
// longer than the request.
type request struct {
	ID       string    `json:"id"`
	At       time.Time `json:"at"`
	Items    int       `json:"items"`
	Device   bool      `json:"device"`
	Score    int       `json:"score"`
	ScoreMs  float64   `json:"score_ms"`
	EnrichMs float64   `json:"enrich_ms"`
	TotalMs  float64   `json:"total_ms"`
}

// recent holds the last requests in a ring.
type recent struct {
	mu   sync.Mutex
	buf  []request
	next int
}

func (r *recent) add(q request) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.buf) < keep {
		r.buf = append(r.buf, q)
		return
	}
	r.buf[r.next] = q
	r.next = (r.next + 1) % keep
}

func (r *recent) slowest(n int) []request {
	r.mu.Lock()
	all := slices.Clone(r.buf)
	r.mu.Unlock()
	slices.SortFunc(all, func(a, b request) int {
		switch {
		case a.TotalMs > b.TotalMs:
			return -1
		case a.TotalMs < b.TotalMs:
			return 1
		}
		return 0
	})
	if len(all) == 0 {
		return []request{}
	}
	return all[:min(n, len(all))]
}

func ms(d time.Duration) float64 { return float64(d.Microseconds()) / 1000 }

func toRequest(e event.Event, res event.Response, t took) request {
	return request{
		ID: e.ID, At: time.Now(), Items: len(e.Items), Device: e.Device != nil,
		Score: res.Score, ScoreMs: ms(t.Score), EnrichMs: ms(t.Enrich), TotalMs: ms(t.Score + t.Enrich),
	}
}

type recentIn struct {
	Slowest int `json:"slowest,omitempty" jsonschema:"how many requests to return, slowest first; default 10, at most 50"`
}

type recentOut struct {
	Kept     int       `json:"kept"`
	Requests []request `json:"requests"`
}

type replayIn struct {
	Event event.Event `json:"event" jsonschema:"the event, the same JSON that POST /score takes"`
	Items int         `json:"items,omitempty" jsonschema:"repeat the event's items until it has this many; 0 keeps them as they are"`
}

type replayOut struct {
	Response event.Response `json:"response"`
	ScoreMs  float64        `json:"score_ms"`
	Items    int            `json:"items"`
}

type runtimeOut struct {
	GoVersion  string  `json:"go_version"`
	Uptime     string  `json:"uptime"`
	Goroutines int     `json:"goroutines"`
	HeapMB     float64 `json:"heap_mb"`
	GCRuns     uint32  `json:"gc_runs"`
	CPUs       int     `json:"cpus"`
}

type configOut struct {
	Rules  int            `json:"rules"`
	ByKind map[string]int `json:"rules_by_kind"`
	First  []scoring.Rule `json:"first_rules"`
}

func startDebug(tracer trace.Tracer, rules []scoring.Rule) func(event.Event, event.Response, took) {
	addr := os.Getenv("DEBUG_MCP_ADDR")
	if addr == "" {
		return func(event.Event, event.Response, took) {}
	}
	// The MCP has no login: it never listens beyond loopback.
	if host, _, err := net.SplitHostPort(addr); err != nil || !(host == "localhost" || net.ParseIP(host).IsLoopback()) {
		slog.Error("debug mcp: DEBUG_MCP_ADDR is not a loopback address", "addr", addr)
		os.Exit(1)
	}

	started := time.Now()
	reqs := &recent{}
	// One replay at a time: it runs next to real traffic, on the same CPU.
	var replaying sync.Mutex
	readOnly := &mcp.ToolAnnotations{ReadOnlyHint: true}
	server := mcp.NewServer(&mcp.Implementation{Name: "scoring-debug", Version: "1"}, nil)

	// The audit log shows who opened the tunnel; this shows what went through it.
	server.AddReceivingMiddleware(func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			if call, ok := req.(*mcp.CallToolRequest); ok {
				slog.InfoContext(ctx, "debug mcp call", "tool", call.Params.Name)
			}
			return next(ctx, method, req)
		}
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "recent_requests",
		Description: "Read only. The last requests this process served, slowest first: item count, device, score, and milliseconds spent scoring and enriching.",
		Annotations: readOnly,
	}, func(_ context.Context, _ *mcp.CallToolRequest, in recentIn) (*mcp.CallToolResult, recentOut, error) {
		n := in.Slowest
		if n <= 0 {
			n = 10
		}
		return nil, recentOut{Kept: keep, Requests: reqs.slowest(min(n, 50))}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "replay",
		Description: fmt.Sprintf("Runs code: scores one event with the rules this process loaded and returns the response "+
			"and the milliseconds it took. Scoring only: the device enrichment is skipped. Stores nothing; uses CPU next to "+
			"real traffic, so one replay at a time, at most %d items.", maxItems),
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in replayIn) (*mcp.CallToolResult, replayOut, error) {
		e := in.Event
		if in.Items > 0 && len(e.Items) > 0 {
			grown := make([]event.Item, 0, in.Items)
			for len(grown) < in.Items {
				grown = append(grown, e.Items[len(grown)%len(e.Items)])
			}
			e.Items = grown
		}
		if len(e.Items) > maxItems {
			return nil, replayOut{}, fmt.Errorf("event has %d items, replay takes at most %d", len(e.Items), maxItems)
		}
		if !replaying.TryLock() {
			return nil, replayOut{}, fmt.Errorf("another replay is running, try again")
		}
		defer replaying.Unlock()

		_, span := tracer.Start(ctx, "debug replay")
		start := time.Now()
		res := scoring.Score(rules, e)
		took := time.Since(start)
		span.End()
		return nil, replayOut{Response: res, ScoreMs: ms(took), Items: len(e.Items)}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "runtime",
		Description: "Read only. Go runtime state of this process: goroutines, heap, GC runs, uptime.",
		Annotations: readOnly,
	}, func(context.Context, *mcp.CallToolRequest, struct{}) (*mcp.CallToolResult, runtimeOut, error) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		return nil, runtimeOut{
			GoVersion: runtime.Version(), Uptime: time.Since(started).Round(time.Second).String(),
			Goroutines: runtime.NumGoroutine(), HeapMB: float64(m.HeapAlloc) / (1 << 20),
			GCRuns: m.NumGC, CPUs: runtime.GOMAXPROCS(0),
		}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "config",
		Description: "Read only. The rule set this process loaded: how many rules, how many per item kind, and the first ten.",
		Annotations: readOnly,
	}, func(context.Context, *mcp.CallToolRequest, struct{}) (*mcp.CallToolResult, configOut, error) {
		out := configOut{Rules: len(rules), ByKind: map[string]int{}, First: rules[:min(10, len(rules))]}
		for _, r := range rules {
			out.ByKind[r.Kind]++
		}
		return nil, out, nil
	})

	// Stateless: after a pod restart, the agent's next call works without a
	// new session.
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server },
		&mcp.StreamableHTTPOptions{Stateless: true})
	go func() {
		slog.Info("debug mcp listening", "addr", addr)
		mux := http.NewServeMux()
		mux.Handle("/mcp", handler)
		if err := http.ListenAndServe(addr, mux); err != nil {
			slog.Error("debug mcp failed", "err", err)
			os.Exit(1)
		}
	}()

	return func(e event.Event, res event.Response, t took) { reqs.add(toRequest(e, res, t)) }
}
