package enrich

import (
	"testing"
)

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"env=prod", "region=us-east"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Key != "env" || rules[0].Value != "prod" {
		t.Errorf("rule 0 mismatch: %+v", rules[0])
	}
	if rules[1].Key != "region" || rules[1].Value != "us-east" {
		t.Errorf("rule 1 mismatch: %+v", rules[1])
	}
}

func TestParseRules_Template(t *testing.T) {
	rules, err := ParseRules([]string{"source={host}:{port}"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].Value != "{host}:{port}" {
		t.Errorf("expected template value, got %q", rules[0].Value)
	}
}

func TestParseRules_MissingEquals(t *testing.T) {
	_, err := ParseRules([]string{"badspec"})
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseRules_EmptyKey(t *testing.T) {
	_, err := ParseRules([]string{"=value"})
	if err == nil {
		t.Error("expected error for empty key")
	}
}

func TestParseRules_SkipsEmptySpecs(t *testing.T) {
	rules, err := ParseRules([]string{"", "  ", "env=prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}

func TestParseRules_ValueWithEquals(t *testing.T) {
	rules, err := ParseRules([]string{"expr=a=b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].Value != "a=b" {
		t.Errorf("expected value 'a=b', got %q", rules[0].Value)
	}
}
