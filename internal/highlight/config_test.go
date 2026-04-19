package highlight

import "testing"

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"level=error:red", "level=warn:yellow"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Field != "level" || rules[0].Value != "error" || rules[0].Color != Red {
		t.Errorf("rule[0] mismatch: %+v", rules[0])
	}
}

func TestParseRules_Substring(t *testing.T) {
	rules, err := ParseRules([]string{"msg~timeout:cyan"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rules[0].Substring {
		t.Error("expected substring=true")
	}
	if rules[0].Value != "timeout" {
		t.Errorf("expected value=timeout, got %q", rules[0].Value)
	}
}

func TestParseRules_UnknownColor(t *testing.T) {
	_, err := ParseRules([]string{"level=error:purple"})
	if err == nil {
		t.Error("expected error for unknown color")
	}
}

func TestParseRules_MissingColor(t *testing.T) {
	_, err := ParseRules([]string{"level=error"})
	if err == nil {
		t.Error("expected error for missing color")
	}
}

func TestParseRules_MissingSep(t *testing.T) {
	_, err := ParseRules([]string{"levelred"})
	if err == nil {
		t.Error("expected error for missing separator")
	}
}

func TestParseRules_Empty(t *testing.T) {
	rules, err := ParseRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules")
	}
}
