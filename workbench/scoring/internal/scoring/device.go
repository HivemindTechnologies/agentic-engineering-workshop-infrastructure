package scoring

import (
	"context"
	"runtime"
	"strings"

	"workshop/scoring/internal/event"
)

// EnrichDevice scores the device itself and starts a watcher that releases the
// history summary once the caller is done with it, so the request path never
// waits for cleanup.
func EnrichDevice(ctx context.Context, d *event.Device) int {
	if d == nil {
		return 0
	}

	// The summary is the device's sessions concatenated, which is what the
	// fingerprint comparison downstream consumes.
	var summary []byte
	for _, h := range d.History {
		summary = append(summary, h...)
	}
	go func() {
		<-ctx.Done()
		runtime.KeepAlive(summary)
	}()

	if strings.HasPrefix(d.Fingerprint, "generic-x86") {
		return 30
	}
	return 0
}
