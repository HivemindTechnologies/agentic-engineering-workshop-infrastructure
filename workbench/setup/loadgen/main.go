// Command loadgen is the outside world: customers sending events to the
// scoring API. It lives outside scoring/ so that deploying a scoring fix never
// restarts it.
//
// The traffic profile is process state set over HTTP and absent from git, so
// Argo CD never reverts a changed profile.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// The wire schema of POST /score, copied from scoring/ so that loadgen builds
// on its own. It holds only the fields loadgen sends.
type item struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type event struct {
	ID    string `json:"id"`
	Items []item `json:"items"`
}

// traffic is the current profile and rate. Part 1 sends only profile a; the
// Part 2 material adds profiles for its incidents. Profile names are opaque
// letters, so the name alone tells a diagnosing agent nothing.
type traffic struct {
	Profile string `json:"profile"`
	Rate    int    `json:"rate"`
}

var profiles = map[string]bool{"a": true}

func build(profile string, n int) event {
	return event{ID: fmt.Sprintf("evt-%d", n), Items: items(3)}
}

func items(n int) []item {
	out := make([]item, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, item{Kind: "bin", Value: fmt.Sprintf("%06d", 400000+rand.Intn(200000))})
	}
	return out
}

func main() {
	target := os.Getenv("TARGET_URL")
	if target == "" {
		target = "http://scoring.demo.svc.cluster.local:8080/score"
	}

	var mu sync.Mutex
	cur := traffic{Profile: "a", Rate: 5}
	changed := make(chan struct{}, 1)

	// GET shows the traffic, PUT ?profile=&rate= changes it. Validation is
	// strict, so a typo fails loudly instead of silently running the baseline.
	http.HandleFunc("/traffic", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			p := r.URL.Query().Get("profile")
			rate, err := strconv.Atoi(r.URL.Query().Get("rate"))
			if !profiles[p] || err != nil || rate < 1 || rate > 100 {
				http.Error(w, "want ?profile=a&rate=1..100", http.StatusBadRequest)
				return
			}
			mu.Lock()
			cur = traffic{Profile: p, Rate: rate}
			mu.Unlock()
			select {
			case changed <- struct{}{}:
			default:
			}
		}
		mu.Lock()
		defer mu.Unlock()
		_ = json.NewEncoder(w).Encode(cur)
	})
	go func() { log.Fatal(http.ListenAndServe(":8081", nil)) }()

	client := &http.Client{Timeout: 30 * time.Second}
	mu.Lock()
	t := cur
	mu.Unlock()
	// A ticker drops ticks while a request is slow, so a degraded service
	// receives less load, as with real clients that wait for an answer.
	ticker := time.NewTicker(time.Second / time.Duration(t.Rate))
	for n := 0; ; n++ {
		select {
		case <-changed:
			mu.Lock()
			t = cur
			mu.Unlock()
			ticker.Reset(time.Second / time.Duration(t.Rate))
		case <-ticker.C:
		}
		body, _ := json.Marshal(build(t.Profile, n))
		resp, err := client.Post(target, "application/json", bytes.NewReader(body))
		if err != nil {
			log.Printf("post failed: %v", err)
			continue
		}
		// Draining and closing keeps the connection reusable, so the generator
		// does not exhaust ephemeral ports.
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
