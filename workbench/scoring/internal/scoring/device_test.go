package scoring

import (
	"context"
	"testing"

	"workshop/scoring/internal/event"
)

func TestEnrichDeviceWeighsEmulator(t *testing.T) {
	d := &event.Device{Fingerprint: "generic-x86-0001", History: []string{"a", "b"}}

	if got := EnrichDevice(context.Background(), d); got != 30 {
		t.Fatalf("EnrichDevice() = %d, want 30", got)
	}
}

func TestEnrichDeviceIgnoresUnknownFingerprint(t *testing.T) {
	d := &event.Device{Fingerprint: "pixel-8-0002"}

	if got := EnrichDevice(context.Background(), d); got != 0 {
		t.Fatalf("EnrichDevice() = %d, want 0", got)
	}
}
