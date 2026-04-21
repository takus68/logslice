package coalesce_test

import (
	"testing"

	"github.com/logslice/logslice/internal/coalesce"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntries(fields []map[string]any) []*parser.Entry {
	entries := make([]*parser.Entry, len(fields))
	for i, f := range fields {
		entries[i] = &parser.Entry{Fields: f}
	}
	return entries
}

func TestRun_FirstNonEmpty(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"a": "", "b": nil, "c": "found"},
	})
	rules, _ := coalesce.ParseRules([]string{"result=a,b,c"})
	out := coalesce.Run(entries, rules)
	if got := out[0].Fields["result"]; got != "found" {
		t.Errorf("expected 'found', got %v", got)
	}
}

func TestRun_AllEmpty(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"a": "", "b": nil},
	})
	rules, _ := coalesce.ParseRules([]string{"result=a,b"})
	out := coalesce.Run(entries, rules)
	if got, ok := out[0].Fields["result"]; ok && got != nil && got != "" {
		t.Errorf("expected nil or empty, got %v", got)
	}
}

func TestRun_FirstWins(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"a": "first", "b": "second"},
	})
	rules, _ := coalesce.ParseRules([]string{"result=a,b"})
	out := coalesce.Run(entries, rules)
	if got := out[0].Fields["result"]; got != "first" {
		t.Errorf("expected 'first', got %v", got)
	}
}

func TestRun_MultipleRules(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"x": nil, "y": "yval", "p": "pval", "q": ""},
	})
	rules, _ := coalesce.ParseRules([]string{"out1=x,y", "out2=q,p"})
	out := coalesce.Run(entries, rules)
	if got := out[0].Fields["out1"]; got != "yval" {
		t.Errorf("expected 'yval', got %v", got)
	}
	if got := out[0].Fields["out2"]; got != "pval" {
		t.Errorf("expected 'pval', got %v", got)
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries([]map[string]any{
		{"a": "val"},
	})
	rules, _ := coalesce.ParseRules([]string{"result=a"})
	_ = coalesce.Run(entries, rules)
	if _, ok := entries[0].Fields["result"]; ok {
		t.Error("original entry should not be mutated")
	}
}

func TestRun_EmptyEntries(t *testing.T) {
	rules, _ := coalesce.ParseRules([]string{"result=a,b"})
	out := coalesce.Run([]*parser.Entry{}, rules)
	if len(out) != 0 {
		t.Errorf("expected 0 entries, got %d", len(out))
	}
}
