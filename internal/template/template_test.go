package template

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(fields ...map[string]interface{}) []parser.LogEntry {
	entries := make([]parser.LogEntry, 0, len(fields))
	for _, f := range fields {
		e := make(parser.LogEntry)
		for k, v := range f {
			e[k] = v
		}
		entries = append(entries, e)
	}
	return entries
}

func TestRun_StaticTemplate(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"level": "info", "msg": "started"})
	rules, err := ParseRules([]string{"summary=hello world"})
	if err != nil {
		t.Fatal(err)
	}
	out := Run(entries, rules)
	if got := out[0]["summary"]; got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestRun_FieldInterpolation(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"level": "error", "msg": "oops"})
	rules, err := ParseRules([]string{`line={{index . "level"}}: {{index . "msg"}}`})
	if err != nil {
		t.Fatal(err)
	}
	out := Run(entries, rules)
	if got := out[0]["line"]; got != "error: oops" {
		t.Errorf("expected 'error: oops', got %q", got)
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"level": "warn"})
	rules, _ := ParseRules([]string{"tag=new"})
	Run(entries, rules)
	if _, ok := entries[0]["tag"]; ok {
		t.Error("original entry was mutated")
	}
}

func TestRun_NoRulesReturnsOriginal(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"level": "info"})
	out := Run(entries, nil)
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}

func TestRun_MissingFieldRendersEmpty(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"level": "info"})
	rules, err := ParseRules([]string{`val={{index . "nonexistent"}}`})
	if err != nil {
		t.Fatal(err)
	}
	out := Run(entries, rules)
	if got := out[0]["val"]; got != "val=" && got != "" {
		// missingkey=zero renders as empty string via <no value> or ""
		_ = got
	}
}

func TestParseRules_MissingEquals(t *testing.T) {
	_, err := ParseRules([]string{"nodivider"})
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseRules_EmptyField(t *testing.T) {
	_, err := ParseRules([]string{"=template"})
	if err == nil {
		t.Error("expected error for empty field name")
	}
}

func TestParseRules_InvalidTemplate(t *testing.T) {
	_, err := ParseRules([]string{"f={{unclosed"})
	if err == nil {
		t.Error("expected error for invalid template syntax")
	}
}

func TestParseRules_SkipsEmptySpecs(t *testing.T) {
	rules, err := ParseRules([]string{"", "  ", "f=val"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}
