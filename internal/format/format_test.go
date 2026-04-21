package format

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(fields []map[string]interface{}) []parser.Entry {
	entries := make([]parser.Entry, len(fields))
	for i, f := range fields {
		entries[i] = parser.Entry{Timestamp: time.Now(), Fields: f}
	}
	return entries
}

func TestRun_UpperFormat(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"},
	})
	rules := []Rule{{Field: "level", Format: "upper"}}
	out := Run(entries, rules)
	if got := out[0].Fields["level"]; got != "INFO" {
		t.Errorf("expected INFO, got %v", got)
	}
}

func TestRun_LowerFormat(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "WARNING"},
	})
	rules := []Rule{{Field: "level", Format: "lower"}}
	out := Run(entries, rules)
	if got := out[0].Fields["level"]; got != "warning" {
		t.Errorf("expected warning, got %v", got)
	}
}

func TestRun_TrimFormat(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"msg": "  hello world  "},
	})
	rules := []Rule{{Field: "msg", Format: "trim"}}
	out := Run(entries, rules)
	if got := out[0].Fields["msg"]; got != "hello world" {
		t.Errorf("expected 'hello world', got %v", got)
	}
}

func TestRun_MissingFieldSkipped(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"},
	})
	rules := []Rule{{Field: "nonexistent", Format: "upper"}}
	out := Run(entries, rules)
	if _, ok := out[0].Fields["nonexistent"]; ok {
		t.Error("expected missing field to remain absent")
	}
}

func TestRun_NoRulesReturnsOriginal(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"},
	})
	out := Run(entries, nil)
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"level=upper", "msg=trim"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Field != "level" || rules[0].Format != "upper" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseRules_UnknownFormat(t *testing.T) {
	_, err := ParseRules([]string{"level=bold"})
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestParseRules_MissingEquals(t *testing.T) {
	_, err := ParseRules([]string{"levelupper"})
	if err == nil {
		t.Error("expected error for missing equals")
	}
}

func TestParseRules_SkipsEmpty(t *testing.T) {
	rules, err := ParseRules([]string{"", "  ", "level=lower"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}
