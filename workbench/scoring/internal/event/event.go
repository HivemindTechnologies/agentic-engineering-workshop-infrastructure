// Package event holds the wire schema shared by the scoring API and the load
// generator, so both agree on field names without a code generator.
package event

type Item struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Device struct {
	Fingerprint string   `json:"fingerprint"`
	History     []string `json:"history,omitempty"`
}

type Event struct {
	ID     string  `json:"id"`
	Items  []Item  `json:"items"`
	Device *Device `json:"device,omitempty"`
}

type Response struct {
	ID      string   `json:"id"`
	Score   int      `json:"score"`
	Matched []string `json:"matched"`
}
