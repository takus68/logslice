package count_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/count"
)

func makeEntries() []map[string]any {
	return []map[string]any{
		{"level": "error", "msg": "disk full"},
		{"level": "info", "msg": "started"},
		{"level": "error", "msg": "connection refused"},
		{"level": "warn", "msg": "low memory"},
		{"level": "error", "msg": "timeout"},
	}
}

func TestRun_CountsMatches(t *testing.T) {
	var buf bytes.Buffer
	n, err := count.Run(makeEntries(), count.Config{Field: "level", Value: "error"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
	if strings.TrimSpace(buf.String()) != "3" {
		t.Errorf("expected output \"3\", got %q", buf.String())
	}
}

func TestRun_NoMatches(t *testing.T) {
	var buf bytes.Buffer
	n, err := count.Run(makeEntries(), count.Config{Field: "level", Value: "debug"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestRun_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	n, err := count.Run([]map[string]any{}, count.Config{Field: "level", Value: "error"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := count.ParseConfig([]string{"field=level", "value=error", "exact=true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Field != "level" || cfg.Value != "error" || !cfg.Exact {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestParseConfig_MissingField(t *testing.T) {
	_, err := count.ParseConfig([]string{"value=error"})
	if err == nil {
		t.Fatal("expected error for missing field option")
	}
}

func TestParseConfig_UnknownOption(t *testing.T) {
	_, err := count.ParseConfig([]string{"field=level", "bogus=x"})
	if err == nil {
		t.Fatal("expected error for unknown option")
	}
}

func TestParseConfig_InvalidFormat(t *testing.T) {
	_, err := count.ParseConfig([]string{"fieldlevel"})
	if err == nil {
		t.Fatal("expected error for missing '=' in option")
	}
}
