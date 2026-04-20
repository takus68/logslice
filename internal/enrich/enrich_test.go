package enrich

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries() []*parser.Entry {
	return []*parser.Entry{
		{
			Timestamp: time.Now(),
			Fields: map[string]interface{}{
				"host":    "web-01",
				"service": "api",
				"level":   "info",
			},
		},
		{
			Timestamp: time.Now(),
			Fields: map[string]interface{}{
				"host":    "db-02",
				"service": "store",
				"level":   "error",
			},
		},
	}
}

func TestRun_StaticField(t *testing.T) {
	entries := makeEntries()
	rules := []Rule{{Key: "env", Value: "production"}}
	out := Run(entries, rules)
	for _, e := range out {
		if got, ok := e.Fields["env"]; !ok || got != "production" {
			t.Errorf("expected env=production, got %v", got)
		}
	}
}

func TestRun_TemplateField(t *testing.T) {
	entries := makeEntries()
	rules := []Rule{{Key: "source", Value: "{host}/{service}"}}
	out := Run(entries, rules)

	expected := []string{"web-01/api", "db-02/store"}
	for i, e := range out {
		if got := e.Fields["source"]; got != expected[i] {
			t.Errorf("entry %d: expected source=%q, got %q", i, expected[i], got)
		}
	}
}

func TestRun_OverwritesExistingField(t *testing.T) {
	entries := makeEntries()
	rules := []Rule{{Key: "level", Value: "overridden"}}
	out := Run(entries, rules)
	for _, e := range out {
		if got := e.Fields["level"]; got != "overridden" {
			t.Errorf("expected level=overridden, got %v", got)
		}
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries()
	rules := []Rule{{Key: "env", Value: "test"}}
	Run(entries, rules)
	for _, e := range entries {
		if _, ok := e.Fields["env"]; ok {
			t.Error("original entry was mutated")
		}
	}
}

func TestRun_EmptyRules(t *testing.T) {
	entries := makeEntries()
	out := Run(entries, nil)
	if len(out) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(out))
	}
}
