package mask

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(fields []map[string]interface{}) []parser.Entry {
	entries := make([]parser.Entry, len(fields))
	for i, f := range fields {
		entries[i] = parser.Entry{Timestamp: time.Now(), Raw: f}
	}
	return entries
}

func TestRun_FullMask(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"email": "user@example.com", "level": "info"},
	})
	rules := []Rule{{Field: "email", Strategy: StrategyFull}}
	out := Run(entries, rules)
	if out[0].Raw["email"] != "***" {
		t.Errorf("expected *** got %v", out[0].Raw["email"])
	}
	if out[0].Raw["level"] != "info" {
		t.Error("level should be unchanged")
	}
}

func TestRun_FullMask_CustomPlaceholder(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"token": "abc123"},
	})
	rules := []Rule{{Field: "token", Strategy: StrategyFull, Placeholder: "[REDACTED]"}}
	out := Run(entries, rules)
	if out[0].Raw["token"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED] got %v", out[0].Raw["token"])
	}
}

func TestRun_PartialMask(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"card": "1234567890123456"},
	})
	rules := []Rule{{Field: "card", Strategy: StrategyPartial, KeepPrefix: 4, KeepSuffix: 4}}
	out := Run(entries, rules)
	if out[0].Raw["card"] != "1234***3456" {
		t.Errorf("unexpected partial mask: %v", out[0].Raw["card"])
	}
}

func TestRun_MissingField(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "warn"},
	})
	rules := []Rule{{Field: "email", Strategy: StrategyFull}}
	out := Run(entries, rules)
	if _, ok := out[0].Raw["email"]; ok {
		t.Error("email should not be added")
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"secret": "mysecret"},
	})
	rules := []Rule{{Field: "secret", Strategy: StrategyFull}}
	Run(entries, rules)
	if entries[0].Raw["secret"] != "mysecret" {
		t.Error("original entry should not be mutated")
	}
}

func TestParseRules_Full(t *testing.T) {
	rules, err := ParseRules([]string{"email:full"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 || rules[0].Field != "email" || rules[0].Strategy != StrategyFull {
		t.Errorf("unexpected rule: %+v", rules)
	}
}

func TestParseRules_Partial(t *testing.T) {
	rules, err := ParseRules([]string{"card:partial:4:4"})
	if err != nil {
		t.Fatal(err)
	}
	if rules[0].KeepPrefix != 4 || rules[0].KeepSuffix != 4 {
		t.Errorf("unexpected rule: %+v", rules)
	}
}

func TestParseRules_UnknownStrategy(t *testing.T) {
	_, err := ParseRules([]string{"field:hash"})
	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestParseRules_MissingPartialArgs(t *testing.T) {
	_, err := ParseRules([]string{"field:partial:2"})
	if err == nil {
		t.Error("expected error when keepSuffix is missing")
	}
}
