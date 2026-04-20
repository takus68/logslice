package validate

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(fields ...map[string]interface{}) []parser.Entry {
	entries := make([]parser.Entry, len(fields))
	for i, f := range fields {
		entries[i] = parser.Entry{Timestamp: time.Now(), Fields: f}
	}
	return entries
}

func TestRun_RequiredFieldPresent(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"level": "info"})
	rules, _ := ParseRules([]string{"level:required"})
	results := Run(entries, rules)
	if len(results[0].Errors) != 0 {
		t.Fatalf("expected no errors, got %v", results[0].Errors)
	}
}

func TestRun_RequiredFieldMissing(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"msg": "hello"})
	rules, _ := ParseRules([]string{"level:required"})
	results := Run(entries, rules)
	if len(results[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %v", results[0].Errors)
	}
}

func TestRun_PatternMatch(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"level": "info"})
	rules, _ := ParseRules([]string{"level:pattern=^(info|warn|error)$"})
	results := Run(entries, rules)
	if len(results[0].Errors) != 0 {
		t.Fatalf("expected no errors, got %v", results[0].Errors)
	}
}

func TestRun_PatternNoMatch(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"level": "debug"})
	rules, _ := ParseRules([]string{"level:pattern=^(info|warn|error)$"})
	results := Run(entries, rules)
	if len(results[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %v", results[0].Errors)
	}
}

func TestRun_TypeCheckNumber(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"latency": float64(42)})
	rules, _ := ParseRules([]string{"latency:type=number"})
	results := Run(entries, rules)
	if len(results[0].Errors) != 0 {
		t.Fatalf("expected no errors, got %v", results[0].Errors)
	}
}

func TestRun_TypeCheckFails(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"latency": "fast"})
	rules, _ := ParseRules([]string{"latency:type=number"})
	results := Run(entries, rules)
	if len(results[0].Errors) != 1 {
		t.Fatalf("expected 1 error, got %v", results[0].Errors)
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"level:required,type=string", "code:pattern=^[0-9]+$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_InvalidType(t *testing.T) {
	_, err := ParseRules([]string{"level:type=object"})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestParseRules_MissingSep(t *testing.T) {
	_, err := ParseRules([]string{"levelrequired"})
	if err == nil {
		t.Fatal("expected error for missing colon separator")
	}
}

func TestParseRules_InvalidPattern(t *testing.T) {
	_, err := ParseRules([]string{"field:pattern=[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}
