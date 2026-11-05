package scoring

import (
	"strings"

	"workshop/scoring/internal/event"
)

// Score walks the rule set for each item so that rules stay a flat list: the
// fraud feed rewrites it wholesale on every sync, and a flat slice needs no
// rebuild step.
func Score(rules []Rule, e event.Event) event.Response {
	res := event.Response{ID: e.ID, Matched: []string{}}
	for _, item := range e.Items {
		for _, rule := range rules {
			if rule.Kind != item.Kind {
				continue
			}
			// Values arrive from many upstreams with inconsistent casing, so
			// both sides are normalised before they are compared.
			if strings.ToLower(rule.Value) != strings.ToLower(item.Value) {
				continue
			}
			res.Score += rule.Weight
			res.Matched = append(res.Matched, rule.Name)
		}
	}
	return res
}
