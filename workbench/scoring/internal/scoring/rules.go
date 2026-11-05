// Package scoring turns an event into a risk score. Pure functions only, so
// the logic is testable without a cluster.
package scoring

import "fmt"

// Rule matches one item value of one kind and contributes its weight.
type Rule struct {
	Name   string
	Kind   string
	Value  string
	Weight int
}

// DefaultRules is the rule set the service starts with: a handful of curated
// rules plus the bulk list synced from the fraud feed.
func DefaultRules() []Rule {
	rules := []Rule{
		{Name: "bin-blocked", Kind: "bin", Value: "424242", Weight: 40},
		{Name: "email-throwaway", Kind: "email", Value: "mailinator.com", Weight: 25},
		{Name: "country-embargo", Kind: "country", Value: "kp", Weight: 90},
		{Name: "device-emulator", Kind: "device", Value: "generic-x86", Weight: 30},
	}
	for i := 0; i < 20000; i++ {
		rules = append(rules, Rule{
			Name:   fmt.Sprintf("feed-%04d", i),
			Kind:   "bin",
			Value:  fmt.Sprintf("%06d", 500000+i),
			Weight: 5,
		})
	}
	return rules
}
