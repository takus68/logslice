package convert

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(fields []map[string]interface{}) []*parser.LogEntry {
	entries := make([]*parser.LogEntry, len(fields))
	for i, f := range fields {
		entries[i] = &parser.LogEntry{
			Timestamp: time.Now(),
			Fields:    f,
		}
	}
	return entries
}

func TestRun_ConvertToInt(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"count": "42"},
	})
	rules := []Rule{{Field: "count", Type: TypeInt}}
	out := Run(entries, rules)
	if v, ok := out[0].Fields["count"].(int64); !ok || v != 42 {
		t.Errorf("expected int64(42), got %v (%T)", out[0].Fields["count"], out[0].Fields["count"])
	}
}

func TestRun_ConvertToFloat(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"ratio": "3.14"},
	})
	rules := []Rule{{Field: "ratio", Type: TypeFloat}}
	out := Run(entries, rules)
	if v, ok := out[0].Fields["ratio"].(float64); !ok || v != 3.14 {
		t.Errorf("expected float64(3.14), got %v", out[0].Fields["ratio"])
	}
}

func TestRun_ConvertToBool(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"active": "true"},
	})
	rules := []Rule{{Field: "active", Type: TypeBool}}
	out := Run(entries, rules)
	if v, ok := out[0].Fields["active"].(bool); !ok || !v {
		t.Errorf("expected bool(true), got %v", out[0].Fields["active"])
	}
}

func TestRun_ConvertToString(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"code": 404},
	})
	rules := []Rule{{Field: "code", Type: TypeString}}
	out := Run(entries, rules)
	if v, ok := out[0].Fields["code"].(string); !ok || v != "404" {
		t.Errorf("expected string \"404\", got %v", out[0].Fields["code"])
	}
}

func TestRun_MissingFieldSkipped(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"},
	})
	rules := []Rule{{Field: "missing", Type: TypeInt}}
	out := Run(entries, rules)
	if _, ok := out[0].Fields["missing"]; ok {
		t.Error("expected missing field to remain absent")
	}
}

func TestRun_InvalidConversionSkipped(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"count": "not-a-number"},
	})
	rules := []Rule{{Field: "count", Type: TypeInt}}
	out := Run(entries, rules)
	if v := out[0].Fields["count"]; v != "not-a-number" {
		t.Errorf("expected original value preserved, got %v", v)
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"n": "7"},
	})
	rules := []Rule{{Field: "n", Type: TypeInt}}
	Run(entries, rules)
	if v := entries[0].Fields["n"]; v != "7" {
		t.Errorf("original entry mutated, got %v", v)
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"count=int", "ratio=float", "active=bool", "name=string"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 4 {
		t.Errorf("expected 4 rules, got %d", len(rules))
	}
}

func TestParseRules_UnknownType(t *testing.T) {
	_, err := ParseRules([]string{"field=datetime"})
	if err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestParseRules_MissingEquals(t *testing.T) {
	_, err := ParseRules([]string{"fieldonly"})
	if err == nil {
		t.Error("expected error for missing equals")
	}
}

func TestParseRules_SkipsEmpty(t *testing.T) {
	rules, err := ParseRules([]string{"", "  ", "n=int"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}
