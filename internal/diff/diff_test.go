package diff

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(fields ...map[string]interface{}) []parser.Entry {
	var entries []parser.Entry
	for _, f := range fields {
		entries = append(entries, parser.Entry{
			Timestamp: time.Now(),
			Fields:    f,
			Raw:       f,
		})
	}
	return entries
}

func TestRun_NoChanges(t *testing.T) {
	left := makeEntries(map[string]interface{}{"id": "1", "level": "info"})
	right := makeEntries(map[string]interface{}{"id": "1", "level": "info"})
	res, err := Run(left, right, "id", ModeAll)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("expected 0 results, got %d", len(res))
	}
}

func TestRun_DetectsAdded(t *testing.T) {
	left := makeEntries(map[string]interface{}{"id": "1", "msg": "a"})
	right := makeEntries(
		map[string]interface{}{"id": "1", "msg": "a"},
		map[string]interface{}{"id": "2", "msg": "b"},
	)
	res, err := Run(left, right, "id", ModeAdded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Tag != "+" {
		t.Fatalf("expected 1 added result, got %v", res)
	}
}

func TestRun_DetectsRemoved(t *testing.T) {
	left := makeEntries(
		map[string]interface{}{"id": "1", "msg": "a"},
		map[string]interface{}{"id": "2", "msg": "b"},
	)
	right := makeEntries(map[string]interface{}{"id": "1", "msg": "a"})
	res, err := Run(left, right, "id", ModeRemoved)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Tag != "-" {
		t.Fatalf("expected 1 removed result, got %v", res)
	}
}

func TestRun_DetectsChanged(t *testing.T) {
	left := makeEntries(map[string]interface{}{"id": "1", "level": "info"})
	right := makeEntries(map[string]interface{}{"id": "1", "level": "error"})
	res, err := Run(left, right, "id", ModeChanged)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 1 || res[0].Tag != "~" {
		t.Fatalf("expected 1 changed result, got %v", res)
	}
}

func TestRun_EmptyKeyField(t *testing.T) {
	_, err := Run(nil, nil, "", ModeAll)
	if err == nil {
		t.Fatal("expected error for empty keyField")
	}
}

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := ParseConfig([]string{"key=id", "mode=added"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.KeyField != "id" {
		t.Errorf("expected key=id, got %q", cfg.KeyField)
	}
	if cfg.Mode != ModeAdded {
		t.Errorf("expected mode=added, got %q", cfg.Mode)
	}
}

func TestParseConfig_MissingKey(t *testing.T) {
	_, err := ParseConfig([]string{"mode=all"})
	if err == nil {
		t.Fatal("expected error when key option is missing")
	}
}

func TestParseConfig_InvalidMode(t *testing.T) {
	_, err := ParseConfig([]string{"key=id", "mode=unknown"})
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestParseConfig_DefaultMode(t *testing.T) {
	cfg, err := ParseConfig([]string{"key=trace_id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Mode != ModeAll {
		t.Errorf("expected default mode=all, got %q", cfg.Mode)
	}
}
