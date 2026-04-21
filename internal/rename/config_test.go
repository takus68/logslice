package rename

import "testing"

func TestParseRules_Valid(t *testing.T) {
	rules, err := ParseRules([]string{"msg=message", "lvl=level"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "msg" || rules[0].To != "message" {
		t.Errorf("unexpected rule[0]: %+v", rules[0])
	}
}

func TestParseRules_Trimmed(t *testing.T) {
	rules, err := ParseRules([]string{" msg = message "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules[0].From != "msg" || rules[0].To != "message" {
		t.Errorf("expected trimmed keys, got %+v", rules[0])
	}
}

func TestParseRules_MissingEquals(t *testing.T) {
	_, err := ParseRules([]string{"msgmessage"})
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseRules_EmptyFrom(t *testing.T) {
	_, err := ParseRules([]string{"=message"})
	if err == nil {
		t.Error("expected error for empty source key")
	}
}

func TestParseRules_EmptyTo(t *testing.T) {
	_, err := ParseRules([]string{"msg="})
	if err == nil {
		t.Error("expected error for empty destination key")
	}
}

func TestParseRules_SkipsEmptySpecs(t *testing.T) {
	rules, err := ParseRules([]string{"", "  ", "msg=message"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(rules))
	}
}
