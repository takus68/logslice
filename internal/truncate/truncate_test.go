package truncate

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries(fields []map[string]interface{}) []*parser.Entry {
	entries := make([]*parser.Entry, len(fields))
	for i, f := range fields {
		entries[i] = &parser.Entry{Timestamp: time.Now(), Fields: f}
	}
	return entries
}

func TestRun_TruncatesLongValue(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"msg": "this is a very long message that exceeds the limit"},
	})
	cfg := Config{Fields: []string{"msg"}, MaxLength: 10, Suffix: "..."}
	out := Run(entries, cfg)
	got := out[0].Fields["msg"].(string)
	if got != "this is a ..." {
		t.Errorf("expected truncated value, got %q", got)
	}
}

func TestRun_ShortValueUnchanged(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"msg": "short"},
	})
	cfg := Config{Fields: []string{"msg"}, MaxLength: 20, Suffix: "..."}
	out := Run(entries, cfg)
	if out[0].Fields["msg"] != "short" {
		t.Errorf("expected unchanged value")
	}
}

func TestRun_MissingFieldSkipped(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"},
	})
	cfg := Config{Fields: []string{"msg"}, MaxLength: 5, Suffix: "..."}
	out := Run(entries, cfg)
	if _, ok := out[0].Fields["msg"]; ok {
		t.Error("expected missing field to remain absent")
	}
}

func TestRun_NoFieldsReturnsOriginal(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"msg": "hello"},
	})
	cfg := Config{Fields: []string{}, MaxLength: 3}
	out := Run(entries, cfg)
	if out[0].Fields["msg"] != "hello" {
		t.Error("expected original entries returned")
	}
}

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := ParseConfig([]string{"fields=msg,body", "max=50", "suffix=~~"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Fields) != 2 || cfg.Fields[0] != "msg" || cfg.Fields[1] != "body" {
		t.Errorf("unexpected fields: %v", cfg.Fields)
	}
	if cfg.MaxLength != 50 {
		t.Errorf("expected max 50, got %d", cfg.MaxLength)
	}
	if cfg.Suffix != "~~" {
		t.Errorf("expected suffix ~~, got %q", cfg.Suffix)
	}
}

func TestParseConfig_InvalidMax(t *testing.T) {
	_, err := ParseConfig([]string{"max=abc"})
	if err == nil {
		t.Error("expected error for invalid max")
	}
}

func TestParseConfig_UnknownOption(t *testing.T) {
	_, err := ParseConfig([]string{"unknown=val"})
	if err == nil {
		t.Error("expected error for unknown option")
	}
}

func TestParseConfig_MissingEquals(t *testing.T) {
	_, err := ParseConfig([]string{"fieldsonly"})
	if err == nil {
		t.Error("expected error for missing equals")
	}
}
