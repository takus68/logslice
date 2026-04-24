package where

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries(fields []map[string]interface{}) []parser.Entry {
	var entries []parser.Entry
	for _, f := range fields {
		entries = append(entries, parser.Entry{Timestamp: time.Now(), Fields: f})
	}
	return entries
}

func TestRun_EqualMatch(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "error", "status": "500"},
		{"level": "info", "status": "200"},
	})
	rules := []Rule{{Field: "level", Op: "==", Value: "error"}}
	out := Run(entries, rules)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestRun_NotEqual(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "error"},
		{"level": "info"},
		{"level": "debug"},
	})
	rules := []Rule{{Field: "level", Op: "!=", Value: "error"}}
	out := Run(entries, rules)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestRun_NumericGreaterThan(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"status": "200"},
		{"status": "404"},
		{"status": "500"},
	})
	rules := []Rule{{Field: "status", Op: ">=", Value: "400"}}
	out := Run(entries, rules)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestRun_Contains(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"message": "connection timeout occurred"},
		{"message": "all good"},
	})
	rules := []Rule{{Field: "message", Op: "contains", Value: "timeout"}}
	out := Run(entries, rules)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestRun_MissingField(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "error"},
	})
	rules := []Rule{{Field: "status", Op: "==", Value: "500"}}
	out := Run(entries, rules)
	if len(out) != 0 {
		t.Fatalf("expected 0, got %d", len(out))
	}
}

func TestRun_NoRules(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"},
	})
	out := Run(entries, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"level==error", "status>=400"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Field != "level" || rules[0].Op != "==" || rules[0].Value != "error" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseRules_InvalidExpression(t *testing.T) {
	_, err := ParseRules([]string{"level"})
	if err == nil {
		t.Fatal("expected error for expression without operator")
	}
}

func TestParseRules_SkipsEmpty(t *testing.T) {
	rules, err := ParseRules([]string{"", "  ", "level==warn"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
}
