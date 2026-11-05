// Command score serves the risk scoring API.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"workshop/scoring/internal/event"
	"workshop/scoring/internal/scoring"
	"workshop/scoring/internal/telemetry"
)

func main() {
	ctx := context.Background()

	shutdown, err := telemetry.Setup(ctx, "scoring")
	if err != nil {
		slog.Error("telemetry setup failed", "err", err)
		os.Exit(1)
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(stopCtx); err != nil {
			slog.Error("telemetry shutdown failed", "err", err)
		}
	}()

	rules := scoring.DefaultRules()
	tracer := otel.Tracer("scoring")
	// A no-op unless DEBUG_MCP_ADDR is set; see debug.go.
	record := startDebug(tracer, rules)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /score", func(w http.ResponseWriter, r *http.Request) {
		var e event.Event
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			slog.WarnContext(r.Context(), "malformed event", "err", err)
			http.Error(w, "malformed event", http.StatusBadRequest)
			return
		}

		res, t := process(r.Context(), tracer, rules, e)
		record(e, res, t)

		slog.InfoContext(r.Context(), "scored event",
			"event.id", e.ID, "score", res.Score, "items", len(e.Items))

		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: otelhttp.NewHandler(mux, "score"),
	}
	slog.Info("scoring listening", "addr", srv.Addr, "rules", len(rules))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}

// took is the time spent in each step of one event.
type took struct {
	Score  time.Duration
	Enrich time.Duration
}

// process scores one event. The handler and the debug MCP's replay both call
// it, so a replay runs the same code as a request.
func process(ctx context.Context, tracer trace.Tracer, rules []scoring.Rule, e event.Event) (event.Response, took) {
	var t took
	start := time.Now()
	res := scoring.Score(rules, e)
	t.Score = time.Since(start)

	if e.Device != nil {
		// Enrichment outlives the response, so a client that hangs up mid
		// request cannot abort the device bookkeeping.
		start = time.Now()
		enrichCtx, span := tracer.Start(context.WithoutCancel(ctx), "enrich")
		res.Score += scoring.EnrichDevice(enrichCtx, e.Device)
		span.End()
		t.Enrich = time.Since(start)
	}
	return res, t
}
