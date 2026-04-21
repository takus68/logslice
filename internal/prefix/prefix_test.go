package prefix

import (
	"testing"
)

func makeEntries(data []map[string]interface{}) []map[string]interface{} {
	return data
}

func TestRun_AddsPrefix(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"msg": "hello", "level": "info"},
	})
	rules := []Rule{{Field: "msg", Prefix: "[APP] "}}
	out := Run(entries, rules)
	if got := out[0]["msg"]; got != "[APP] hello" {
		t.Errorf("expected '[APP] hello', got %q", got)
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	original := map[string]interface{}{"msg": "world"}
	entries := []map[string]interface{}{original}
	rules := []Rule{{Field: "msg", Prefix: "PRE:"}}
	Run(entries, rules)
	if original["msg"] != "world" {
		t.Errorf("original entry was mutated")
	}
}

func TestRun_MissingFieldSkipped(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "warn"},
	})
	rules := []Rule{{Field: "msg", Prefix: "X"}}
	out := Run(entries, rules)
	if _, ok := out[0]["msg"]; ok {
		t.Errorf("expected missing field to remain absent")
	}
}

func TestRun_NonStringFieldSkipped(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"count": 42},
	})
	rules := []Rule{{Field: "count", Prefix: "n="}}
	out := Run(entries, rules)
	if got := out[0]["count"]; got != 42 {
		t.Errorf("expected non-string field to be unchanged, got %v", got)
	}
}

func TestRun_MultipleRules(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"a": "foo", "b": "bar"},
	})
	rules := []Rule{
		{Field: "a", Prefix: "A:"},
		{Field: "b", Prefix: "B:"},
	}
	out := Run(entries, rules)
	if got := out[0]["a"]; got != "A:foo" {
		t.Errorf("expected 'A:foo', got %q", got)
	}
	if got := out[0]["b"]; got != "B:bar" {
		t.Errorf("expected 'B:bar', got %q", got)
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"msg=INFO: ", "level=LVL-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Field != "msg" || rules[0].Prefix != "INFO: " {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseRules_MissingEquals(t *testing.T) {
	_, err := ParseRules([]string{"msgINFO"})
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseRules_EmptyField(t *testing.T) {
	_, err := ParseRules([]string{"=prefix"})
	if err == nil {
		t.Error("expected error for empty field name")
	}
}

func TestParseRules_SkipsEmptySpecs(t *testing.T) {
	rules, err := ParseRules([]string{"", "  ", "msg=X"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}
