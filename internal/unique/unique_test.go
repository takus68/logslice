package unique

import (
	"testing"
)

func makeEntries(fields []map[string]interface{}) []Entry {
	entries := make([]Entry, len(fields))
	for i, f := range fields {
		entries[i] = f
	}
	return entries
}

func TestRun_UniqueValues(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"},
		{"level": "warn"},
		{"level": "info"},
		{"level": "error"},
	})
	result := Run(entries, Config{Field: "level"})
	if len(result) != 3 {
		t.Fatalf("expected 3 unique values, got %d", len(result))
	}
}

func TestRun_Sorted(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "warn"},
		{"level": "info"},
		{"level": "error"},
	})
	result := Run(entries, Config{Field: "level", Sorted: true})
	expected := []string{"error", "info", "warn"}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("index %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestRun_MissingFieldSkipped(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"},
		{"msg": "no level here"},
		{"level": "info"},
	})
	result := Run(entries, Config{Field: "level"})
	if len(result) != 1 {
		t.Fatalf("expected 1 unique value, got %d", len(result))
	}
}

func TestRun_EmptyEntries(t *testing.T) {
	result := Run([]Entry{}, Config{Field: "level"})
	if len(result) != 0 {
		t.Fatalf("expected 0 results, got %d", len(result))
	}
}

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := ParseConfig([]string{"field=level", "sorted=true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Field != "level" {
		t.Errorf("expected field=level, got %q", cfg.Field)
	}
	if !cfg.Sorted {
		t.Error("expected sorted=true")
	}
}

func TestParseConfig_MissingField(t *testing.T) {
	_, err := ParseConfig([]string{"sorted=true"})
	if err == nil {
		t.Fatal("expected error for missing field option")
	}
}

func TestParseConfig_InvalidSorted(t *testing.T) {
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
		t.Fatal("expected error for missing equals")
	}
}
