package pivot

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(rows []map[string]interface{}) []parser.Entry {
	entries := make([]parser.Entry, len(rows))
	for i, r := range rows {
		entries[i] = parser.Entry{Fields: r}
	}
	return entries
}

func TestRun_BasicPivot(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"host": "web1", "metric": "cpu", "value": "80"},
		{"host": "web1", "metric": "mem", "value": "60"},
		{"host": "web2", "metric": "cpu", "value": "40"},
		{"host": "web2", "metric": "mem", "value": "55"},
	})
	cfg := Config{GroupField: "host", KeyField: "metric", ValueField: "value"}
	out := Run(entries, cfg)
	if len(out) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(out))
	}
	if out[0].Fields["host"] != "web1" {
		t.Errorf("expected first row host=web1, got %v", out[0].Fields["host"])
	}
	if out[0].Fields["cpu"] != "80" {
		t.Errorf("expected cpu=80 for web1, got %v", out[0].Fields["cpu"])
	}
	if out[1].Fields["mem"] != "55" {
		t.Errorf("expected mem=55 for web2, got %v", out[1].Fields["mem"])
	}
}

func TestRun_MissingGroupField(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"metric": "cpu", "value": "80"},
	})
	cfg := Config{GroupField: "host", KeyField: "metric", ValueField: "value"}
	out := Run(entries, cfg)
	if len(out) != 0 {
		t.Errorf("expected 0 rows when group field missing, got %d", len(out))
	}
}

func TestRun_MissingKeyField(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"host": "web1", "value": "80"},
	})
	cfg := Config{GroupField: "host", KeyField: "metric", ValueField: "value"}
	out := Run(entries, cfg)
	if len(out) != 0 {
		t.Errorf("expected 0 rows when key field missing, got %d", len(out))
	}
}

func TestRun_EmptyEntries(t *testing.T) {
	out := Run([]parser.Entry{}, Config{GroupField: "host", KeyField: "metric", ValueField: "value"})
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d", len(out))
	}
}

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := ParseConfig([]string{"key=metric", "value=count", "group=host"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.KeyField != "metric" || cfg.ValueField != "count" || cfg.GroupField != "host" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestParseConfig_MissingKey(t *testing.T) {
	_, err := ParseConfig([]string{"value=count", "group=host"})
	if err == nil {
		t.Error("expected error for missing key option")
	}
}

func TestParseConfig_UnknownOption(t *testing.T) {
	_, err := ParseConfig([]string{"key=metric", "value=count", "group=host", "bogus=x"})
	if err == nil {
		t.Error("expected error for unknown option")
	}
}

func TestParseConfig_InvalidFormat(t *testing.T) {
	_, err := ParseConfig([]string{"keymetric"})
	if err == nil {
		t.Error("expected error for malformed option")
	}
}
