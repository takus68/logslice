package group

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries(fields []map[string]any) []parser.Entry {
	entries := make([]parser.Entry, len(fields))
	for i, f := range fields {
		entries[i] = parser.Entry{Timestamp: time.Now(), Fields: f}
	}
	return entries
}

func TestRun_GroupsByField(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"level": "info", "msg": "a"},
		{"level": "error", "msg": "b"},
		{"level": "info", "msg": "c"},
	})

	result := Run(entries, Config{Field: "level"})

	if len(result.Groups["info"]) != 2 {
		t.Errorf("expected 2 info entries, got %d", len(result.Groups["info"]))
	}
	if len(result.Groups["error"]) != 1 {
		t.Errorf("expected 1 error entry, got %d", len(result.Groups["error"]))
	}
}

func TestRun_MissingFieldKey(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"msg": "no level here"},
	})

	result := Run(entries, Config{Field: "level"})

	if len(result.Groups["<missing>"]) != 1 {
		t.Errorf("expected 1 missing entry, got %d", len(result.Groups["<missing>"]))
	}
}

func TestRun_SortedKeys(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"level": "warn"},
		{"level": "debug"},
		{"level": "error"},
		{"level": "info"},
	})

	result := Run(entries, Config{Field: "level", Sorted: true})

	want := []string{"debug", "error", "info", "warn"}
	for i, k := range result.Keys {
		if k != want[i] {
			t.Errorf("key[%d]: got %q, want %q", i, k, want[i])
		}
	}
}

func TestRun_EmptyEntries(t *testing.T) {
	result := Run(nil, Config{Field: "level"})
	if len(result.Keys) != 0 {
		t.Errorf("expected no keys, got %d", len(result.Keys))
	}
}
