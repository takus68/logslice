package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries() []parser.LogEntry {
	return []parser.LogEntry{
		{
			Timestamp: time.Now(),
			Fields: map[string]interface{}{
				"level":   "error",
				"message": "disk full",
			},
		},
		{
			Timestamp: time.Now(),
			Fields: map[string]interface{}{
				"level":   "info",
				"message": "service started",
			},
		},
		{
			Timestamp: time.Now(),
			Fields: map[string]interface{}{
				"level":   "error",
				"message": "connection refused",
			},
		},
	}
}

func TestByField_ExactMatch(t *testing.T) {
	entries := makeEntries()
	result := ByField(entries, FieldMatcher{Key: "level", Value: "error", Exact: true})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestByField_SubstringMatch(t *testing.T) {
	entries := makeEntries()
	result := ByField(entries, FieldMatcher{Key: "message", Value: "disk", Exact: false})
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
}

func TestByField_CaseInsensitive(t *testing.T) {
	entries := makeEntries()
	result := ByField(entries, FieldMatcher{Key: "level", Value: "ERROR", Exact: false})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestByField_MissingKey(t *testing.T) {
	entries := makeEntries()
	result := ByField(entries, FieldMatcher{Key: "nonexistent", Value: "x", Exact: true})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestByField_NoMatches(t *testing.T) {
	entries := makeEntries()
	result := ByField(entries, FieldMatcher{Key: "level", Value: "warn", Exact: true})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}
