package extract

import (
	"testing"
)

func makeEntries() []Entry {
	return []Entry{
		{"level": "info", "msg": "started", "host": "web-1"},
		{"level": "error", "msg": "failed", "host": "web-2", "code": 500},
		{"level": "debug", "msg": "ping"},
	}
}

func TestRun_ExtractsFields(t *testing.T) {
	entries := makeEntries()
	cfg := Config{Fields: []string{"level", "msg"}}
	out := Run(entries, cfg)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	for _, e := range out {
		if _, ok := e["host"]; ok {
			t.Error("expected 'host' to be excluded")
		}
		if _, ok := e["level"]; !ok {
			t.Error("expected 'level' to be present")
		}
	}
}

func TestRun_MissingFieldOmitted(t *testing.T) {
	entries := makeEntries()
	cfg := Config{Fields: []string{"code"}, KeepMissing: false}
	out := Run(entries, cfg)
	// Only the second entry has 'code'
	if v, ok := out[0]["code"]; ok {
		t.Errorf("expected 'code' absent in entry 0, got %v", v)
	}
	if _, ok := out[1]["code"]; !ok {
		t.Error("expected 'code' present in entry 1")
	}
}

func TestRun_KeepMissing(t *testing.T) {
	entries := makeEntries()
	cfg := Config{Fields: []string{"code"}, KeepMissing: true}
	out := Run(entries, cfg)
	for i, e := range out {
		if _, ok := e["code"]; !ok {
			t.Errorf("entry %d: expected 'code' key to be present with keep_missing=true", i)
		}
	}
	if out[0]["code"] != nil {
		t.Errorf("entry 0: expected nil for missing 'code', got %v", out[0]["code"])
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries()
	cfg := Config{Fields: []string{"level"}}
	Run(entries, cfg)
	if _, ok := entries[0]["msg"]; !ok {
		t.Error("original entry was mutated")
	}
}

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := ParseConfig([]string{"fields=level,msg", "keep_missing=true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(cfg.Fields))
	}
	if !cfg.KeepMissing {
		t.Error("expected KeepMissing=true")
	}
}

func TestParseConfig_MissingFields(t *testing.T) {
	_, err := ParseConfig([]string{"keep_missing=false"})
	if err == nil {
		t.Error("expected error when no fields specified")
	}
}

func TestParseConfig_InvalidOption(t *testing.T) {
	_, err := ParseConfig([]string{"fields=x", "unknown=yes"})
	if err == nil {
		t.Error("expected error for unknown option")
	}
}

func TestParseConfig_InvalidKeepMissing(t *testing.T) {
	_, err := ParseConfig([]string{"fields=x", "keep_missing=maybe"})
	if err == nil {
		t.Error("expected error for invalid keep_missing value")
	}
}

func TestParseConfig_MissingEquals(t *testing.T) {
	_, err := ParseConfig([]string{"fields"})
	if err == nil {
		t.Error("expected error for option missing '='")
	}
}
