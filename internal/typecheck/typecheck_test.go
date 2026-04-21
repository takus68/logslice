package typecheck

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries(fields []map[string]interface{}) []parser.Entry {
	entries := make([]parser.Entry, len(fields))
	for i, f := range fields {
		entries[i] = parser.Entry{
			Timestamp: time.Now(),
			Fields:    f,
		}
	}
	return entries
}

func TestRun_InfersStringType(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"},
		{"level": "warn"},
	})
	reports := Run(entries, nil)
	r, ok := reports["level"]
	if !ok {
		t.Fatal("expected report for 'level'")
	}
	if r.Types[TypeString] != 2 {
		t.Errorf("expected 2 strings, got %d", r.Types[TypeString])
	}
	if r.Total != 2 {
		t.Errorf("expected total 2, got %d", r.Total)
	}
}

func TestRun_InfersIntType(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"status": "200"},
		{"status": "404"},
	})
	reports := Run(entries, nil)
	r := reports["status"]
	if r.Types[TypeInt] != 2 {
		t.Errorf("expected 2 ints, got %d", r.Types[TypeInt])
	}
}

func TestRun_InfersFloatType(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"latency": "1.23"},
	})
	reports := Run(entries, nil)
	r := reports["latency"]
	if r.Types[TypeFloat] != 1 {
		t.Errorf("expected float, got %v", r.Types)
	}
}

func TestRun_InfersBoolType(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"ok": "true"},
		{"ok": "false"},
	})
	reports := Run(entries, nil)
	r := reports["ok"]
	if r.Types[TypeBool] != 2 {
		t.Errorf("expected 2 bools, got %d", r.Types[TypeBool])
	}
}

func TestRun_MixedTypes(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"val": "hello"},
		{"val": "42"},
		{"val": "true"},
	})
	reports := Run(entries, nil)
	r := reports["val"]
	if r.Total != 3 {
		t.Errorf("expected total 3, got %d", r.Total)
	}
	if r.Types[TypeString] != 1 || r.Types[TypeInt] != 1 || r.Types[TypeBool] != 1 {
		t.Errorf("unexpected type distribution: %v", r.Types)
	}
}

func TestRun_FiltersByFields(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info", "status": "200"},
	})
	reports := Run(entries, []string{"level"})
	if _, ok := reports["status"]; ok {
		t.Error("expected 'status' to be excluded")
	}
	if _, ok := reports["level"]; !ok {
		t.Error("expected 'level' to be included")
	}
}

func TestRun_EmptyEntries(t *testing.T) {
	reports := Run([]parser.Entry{}, nil)
	if len(reports) != 0 {
		t.Errorf("expected empty reports, got %d entries", len(reports))
	}
}
