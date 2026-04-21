package compact

import (
	"testing"
)

func makeEntries(maps ...map[string]any) []map[string]any {
	return maps
}

func TestRun_RemovesNullByDefault(t *testing.T) {
	entries := makeEntries(map[string]any{"a": "hello", "b": nil, "c": 42})
	cfg := Config{RemoveNull: true}
	out := Run(entries, cfg)
	if _, ok := out[0]["b"]; ok {
		t.Error("expected nil field 'b' to be removed")
	}
	if out[0]["a"] != "hello" {
		t.Errorf("expected 'a' to be 'hello', got %v", out[0]["a"])
	}
}

func TestRun_RemovesEmptyStrings(t *testing.T) {
	entries := makeEntries(map[string]any{"msg": "", "level": "info"})
	cfg := Config{RemoveEmpty: true, RemoveNull: false}
	out := Run(entries, cfg)
	if _, ok := out[0]["msg"]; ok {
		t.Error("expected empty string field 'msg' to be removed")
	}
	if out[0]["level"] != "info" {
		t.Errorf("expected 'level' to be 'info', got %v", out[0]["level"])
	}
}

func TestRun_RestrictsToFields(t *testing.T) {
	entries := makeEntries(map[string]any{"a": nil, "b": nil, "c": "keep"})
	cfg := Config{RemoveNull: true, Fields: []string{"a"}}
	out := Run(entries, cfg)
	if _, ok := out[0]["a"]; ok {
		t.Error("expected field 'a' to be removed")
	}
	// 'b' is nil but not in Fields, so it should remain
	if _, ok := out[0]["b"]; !ok {
		t.Error("expected field 'b' to remain (not in restricted fields)")
	}
}

func TestRun_EmptyEntries(t *testing.T) {
	out := Run([]map[string]any{}, Config{RemoveNull: true})
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d entries", len(out))
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	orig := map[string]any{"x": nil, "y": "val"}
	entries := makeEntries(orig)
	Run(entries, Config{RemoveNull: true})
	if _, ok := orig["x"]; !ok {
		t.Error("original entry was mutated")
	}
}

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := ParseConfig([]string{"remove_empty=true", "remove_null=false", "fields=a,b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.RemoveEmpty {
		t.Error("expected RemoveEmpty=true")
	}
	if cfg.RemoveNull {
		t.Error("expected RemoveNull=false")
	}
	if len(cfg.Fields) != 2 || cfg.Fields[0] != "a" || cfg.Fields[1] != "b" {
		t.Errorf("unexpected fields: %v", cfg.Fields)
	}
}

func TestParseConfig_InvalidOption(t *testing.T) {
	_, err := ParseConfig([]string{"badoption"})
	if err == nil {
		t.Error("expected error for missing '=' in option")
	}
}

func TestParseConfig_UnknownKey(t *testing.T) {
	_, err := ParseConfig([]string{"unknown=value"})
	if err == nil {
		t.Error("expected error for unknown key")
	}
}
