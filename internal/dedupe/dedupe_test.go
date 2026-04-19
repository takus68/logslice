package dedupe

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries(fields []map[string]any) []parser.LogEntry {
	entries := make([]parser.LogEntry, len(fields))
	for i, f := range fields {
		entries[i] = parser.LogEntry{
			Timestamp: time.Now(),
			Fields:    f,
		}
	}
	return entries
}

func TestRun_NoDuplicates(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"msg": "a"},
		{"msg": "b"},
	})
	result := Run(entries, Options{Strategy: ByFullEntry})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestRun_RemovesDuplicates(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"msg": "hello"},
		{"msg": "hello"},
		{"msg": "world"},
	})
	result := Run(entries, Options{Strategy: ByFullEntry})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestRun_ByFields_Deduplicates(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"level": "info", "msg": "start"},
		{"level": "info", "msg": "different"},
		{"level": "error", "msg": "boom"},
	})
	opts := Options{Strategy: ByFields, Fields: []string{"level"}}
	result := Run(entries, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestRun_ByFields_MissingField(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"msg": "a"},
		{"msg": "b"},
	})
	opts := Options{Strategy: ByFields, Fields: []string{"level"}}
	result := Run(entries, opts)
	// both entries have no "level" field, so they hash the same — only 1 kept
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestRun_Empty(t *testing.T) {
	result := Run(nil, Options{})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}
