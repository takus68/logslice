package sort_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	logsort "github.com/yourorg/logslice/internal/sort"
)

func makeEntries(fields []map[string]any) []parser.Entry {
	entries := make([]parser.Entry, len(fields))
	for i, f := range fields {
		entries[i] = parser.Entry{Fields: f, Timestamp: time.Time{}}
	}
	return entries
}

func TestRun_SortAscending(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"level": "warn"},
		{"level": "error"},
		{"level": "info"},
	})
	cfg := logsort.Config{Field: "level", Descending: false}
	out := logsort.Run(entries, cfg)
	expect := []string{"error", "info", "warn"}
	for i, e := range out {
		if e.Fields["level"] != expect[i] {
			t.Errorf("index %d: got %v, want %v", i, e.Fields["level"], expect[i])
		}
	}
}

func TestRun_SortDescending(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"level": "info"},
		{"level": "error"},
		{"level": "warn"},
	})
	cfg := logsort.Config{Field: "level", Descending: true}
	out := logsort.Run(entries, cfg)
	expect := []string{"warn", "info", "error"}
	for i, e := range out {
		if e.Fields["level"] != expect[i] {
			t.Errorf("index %d: got %v, want %v", i, e.Fields["level"], expect[i])
		}
	}
}

func TestRun_MissingFieldSortsToEnd(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"level": "info"},
		{"msg": "no level here"},
		{"level": "error"},
	})
	cfg := logsort.Config{Field: "level", Descending: false}
	out := logsort.Run(entries, cfg)
	if out[2].Fields["msg"] != "no level here" {
		t.Errorf("expected missing-field entry at end, got %v", out[2].Fields)
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"level": "warn"},
		{"level": "info"},
	})
	origFirst := entries[0].Fields["level"]
	cfg := logsort.Config{Field: "level", Descending: false}
	logsort.Run(entries, cfg)
	if entries[0].Fields["level"] != origFirst {
		t.Errorf("original slice was mutated")
	}
}

func TestParseConfig_Defaults(t *testing.T) {
	cfg, err := logsort.ParseConfig(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Field != "timestamp" {
		t.Errorf("expected default field 'timestamp', got %q", cfg.Field)
	}
	if cfg.Descending {
		t.Error("expected ascending by default")
	}
}

func TestParseConfig_CustomField(t *testing.T) {
	cfg, err := logsort.ParseConfig([]string{"field=level", "order=desc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Field != "level" {
		t.Errorf("expected field 'level', got %q", cfg.Field)
	}
	if !cfg.Descending {
		t.Error("expected descending")
	}
}

func TestParseConfig_InvalidOrder(t *testing.T) {
	_, err := logsort.ParseConfig([]string{"order=sideways"})
	if err == nil {
		t.Error("expected error for invalid order")
	}
}

func TestParseConfig_UnknownOption(t *testing.T) {
	_, err := logsort.ParseConfig([]string{"foo=bar"})
	if err == nil {
		t.Error("expected error for unknown option")
	}
}
