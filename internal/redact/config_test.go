package redact

import (
	"testing"
)

func TestParseConfig_Fields(t *testing.T) {
	cfg, err := ParseConfig(Options{RedactFields: "password,token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(cfg.Fields))
	}
}

func TestParseConfig_FieldsTrimmed(t *testing.T) {
	cfg, err := ParseConfig(Options{RedactFields: " password , token "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Fields[0] != "password" || cfg.Fields[1] != "token" {
		t.Errorf("fields not trimmed: %v", cfg.Fields)
	}
}

func TestParseConfig_MaskPattern(t *testing.T) {
	cfg, err := ParseConfig(Options{MaskPatterns: []string{"email=[^@]+"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Patterns["email"]; !ok {
		t.Error("expected email pattern")
	}
}

func TestParseConfig_InvalidPattern(t *testing.T) {
	_, err := ParseConfig(Options{MaskPatterns: []string{"email=[invalid"}})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestParseConfig_MissingEquals(t *testing.T) {
	_, err := ParseConfig(Options{MaskPatterns: []string{"emailnoequals"}})
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseConfig_Empty(t *testing.T) {
	cfg, err := ParseConfig(Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Fields) != 0 || len(cfg.Patterns) != 0 {
		t.Error("expected empty config")
	}
}
