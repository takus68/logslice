package compute

import (
	"testing"
)

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"total=requests+errors"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	r := rules[0]
	if r.Dest != "total" || r.Left != "requests" || r.Op != "+" || r.Right != "errors" {
		t.Fatalf("unexpected rule: %+v", r)
	}
}

func TestParseRules_Division(t *testing.T) {
	rules, err := ParseRules([]string{"rate=count/total"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].Op != "/" {
		t.Fatalf("expected op '/', got %q", rules[0].Op)
	}
}

func TestParseRules_MissingEquals(t *testing.T) {
	_, err := ParseRules([]string{"totala+b"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseRules_EmptyDest(t *testing.T) {
	_, err := ParseRules([]string{"=a+b"})
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestParseRules_NoOperator(t *testing.T) {
	_, err := ParseRules([]string{"dest=leftonlynooperator"})
	if err == nil {
		t.Fatal("expected error when no operator found")
	}
}

func TestParseRules_SkipsEmptySpecs(t *testing.T) {
	rules, err := ParseRules([]string{"", "  ", "res=a+b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
}

func TestParseRules_MultipleRules(t *testing.T) {
	rules, err := ParseRules([]string{"s=a+b", "p=a*b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}
