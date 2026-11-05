package scoring

import (
	"testing"

	"workshop/scoring/internal/event"
)

func TestScoreSumsMatchedWeights(t *testing.T) {
	rules := []Rule{
		{Name: "bin-blocked", Kind: "bin", Value: "424242", Weight: 40},
		{Name: "email-throwaway", Kind: "email", Value: "mailinator.com", Weight: 25},
	}
	e := event.Event{
		ID: "e1",
		Items: []event.Item{
			{Kind: "bin", Value: "424242"},
			{Kind: "email", Value: "MAILINATOR.COM"},
			{Kind: "bin", Value: "555555"},
		},
	}

	got := Score(rules, e)

	if got.Score != 65 {
		t.Fatalf("score = %d, want 65", got.Score)
	}
	if len(got.Matched) != 2 {
		t.Fatalf("matched = %v, want two entries", got.Matched)
	}
}

func TestScoreIgnoresRulesOfOtherKind(t *testing.T) {
	rules := []Rule{{Name: "bin-blocked", Kind: "bin", Value: "424242", Weight: 40}}
	e := event.Event{ID: "e2", Items: []event.Item{{Kind: "email", Value: "424242"}}}

	if got := Score(rules, e); got.Score != 0 {
		t.Fatalf("score = %d, want 0", got.Score)
	}
}

func TestDefaultRulesAreNonEmpty(t *testing.T) {
	if n := len(DefaultRules()); n < 1000 {
		t.Fatalf("DefaultRules() = %d rules, want at least 1000", n)
	}
}
