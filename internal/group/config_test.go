package group

import (
	"testing"
)

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := ParseConfig([]string{"field=level"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Field != "level" {
		t.Errorf("expected field=level, got %q", cfg.Field)
	}
	if cfg.Sorted {
		t.Errorf("expected sorted=false by default")
	}
}

func TestParseConfig_WithSorted(t *testing.T) {
	cfg, err := ParseConfig([]string{"field=level", "sorted=true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Sorted {
		t.Errorf("expected sorted=true")
	}
}

func TestParseConfig_MissingField(t *testing.T) {
	_, err := ParseConfig([]string{"sorted=true"})
	if err == nil {
		t.Fatal("expected error for missing field option")
	}
}

func TestParseConfig_EmptyFieldValue(t *testing.T) {
	_, err := ParseConfig([]string{"field="})
	if err == nil {
		t.Fatal("expected error for empty field value")
	}
}

func TestParseConfig_InvalidSortedValue(t *testing.T) {
	_, err := ParseConfig([]string{"field=level", "sorted=yes"})
	if err == nil {
		t.Fatal("expected error for invalid sorted value")
	}
}

func TestParseConfig_UnknownOption(t *testing.T) {
	_, err := ParseConfig([]string{"field=level", "foo=bar"})
	if err == nil {
		t.Fatal("expected error for unknown option")
	}
}

func TestParseConfig_MissingEquals(t *testing.T) {
	_, err := ParseConfig([]string{"fieldlevel"})
	if err == nil {
		t.Fatal("expected error for missing equals separator")
	}
}

func TestParseConfig_SkipsEmptyOpts(t *testing.T) {
	cfg, err := ParseConfig([]string{"  ", "field=ts", ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Field != "ts" {
		t.Errorf("expected field=ts, got %q", cfg.Field)
	}
}
