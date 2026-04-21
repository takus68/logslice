package coalesce_test

import (
	"testing"

	"github.com/logslice/logslice/internal/coalesce"
)

func TestParseRules_Valid(t *testing.T) {
	rules, err := coalesce.ParseRules([]string{"result=a,b,c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Dest != "result" {
		t.Errorf("expected dest 'result', got %q", rules[0].Dest)
	}
	if len(rules[0].Sources) != 3 {
		t.Errorf("expected 3 sources, got %d", len(rules[0].Sources))
	}
}

func TestParseRules_Trimmed(t *testing.T) {
	rules, err := coalesce.ParseRules([]string{" out = a , b "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].Dest != "out" {
		t.Errorf("expected 'out', got %q", rules[0].Dest)
	}
	if rules[0].Sources[0] != "a" || rules[0].Sources[1] != "b" {
		t.Errorf("unexpected sources: %v", rules[0].Sources)
	}
}

func TestParseRules_MissingEquals(t *testing.T) {
	_, err := coalesce.ParseRules([]string{"nodest"})
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseRules_EmptyDest(t *testing.T) {
	_, err := coalesce.ParseRules([]string{"=a,b"})
	if err == nil {
		t.Error("expected error for empty destination")
	}
}

func TestParseRules_NoSources(t *testing.T) {
	_, err := coalesce.ParseRules([]string{"dest="})
	if err == nil {
		t.Error("expected error for empty sources")
	}
}

func TestParseRules_SkipsEmpty(t *testing.T) {
	rules, err := coalesce.ParseRules([]string{"", "  ", "out=a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}
