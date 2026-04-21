package window

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(timestamps []string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(timestamps))
	for _, ts := range timestamps {
		entries = append(entries, parser.Entry{
			Raw: map[string]interface{}{"time": ts, "msg": "hello"},
		})
	}
	return entries
}

func TestRun_SingleWindow(t *testing.T) {
	entries := makeEntries([]string{
		"2024-01-01T00:00:10Z",
		"2024-01-01T00:00:30Z",
		"2024-01-01T00:00:50Z",
	})
	cfg := Config{Size: time.Minute, Field: "time"}
	windows, err := Run(entries, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(windows))
	}
	if len(windows[0].Entries) != 3 {
		t.Errorf("expected 3 entries in window, got %d", len(windows[0].Entries))
	}
}

func TestRun_MultipleWindows(t *testing.T) {
	entries := makeEntries([]string{
		"2024-01-01T00:00:10Z",
		"2024-01-01T00:01:05Z",
		"2024-01-01T00:02:20Z",
	})
	cfg := Config{Size: time.Minute, Field: "time"}
	windows, err := Run(entries, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(windows))
	}
}

func TestRun_SkipsMissingField(t *testing.T) {
	entries := []parser.Entry{
		{Raw: map[string]interface{}{"msg": "no timestamp"}},
		{Raw: map[string]interface{}{"time": "2024-01-01T00:00:05Z"}},
	}
	cfg := Config{Size: time.Minute, Field: "time"}
	windows, err := Run(entries, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(windows) != 1 || len(windows[0].Entries) != 1 {
		t.Errorf("expected 1 window with 1 entry, got %v", windows)
	}
}

func TestRun_InvalidSize(t *testing.T) {
	cfg := Config{Size: 0}
	_, err := Run(nil, cfg)
	if err == nil {
		t.Fatal("expected error for zero size")
	}
}

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := ParseConfig([]string{"size=1m", "field=timestamp", "tumbling=false"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Size != time.Minute {
		t.Errorf("expected 1m, got %v", cfg.Size)
	}
	if cfg.Field != "timestamp" {
		t.Errorf("expected 'timestamp', got %q", cfg.Field)
	}
	if cfg.Tumbling {
		t.Errorf("expected tumbling=false")
	}
}

func TestParseConfig_MissingSize(t *testing.T) {
	_, err := ParseConfig([]string{"field=time"})
	if err == nil {
		t.Fatal("expected error when size is missing")
	}
}

func TestParseConfig_InvalidDuration(t *testing.T) {
	_, err := ParseConfig([]string{"size=notaduration"})
	if err == nil {
		t.Fatal("expected error for invalid duration")
	}
}

func TestParseConfig_UnknownOption(t *testing.T) {
	_, err := ParseConfig([]string{"size=1m", "unknown=foo"})
	if err == nil {
		t.Fatal("expected error for unknown option")
	}
}

func TestParseConfig_InvalidTumblingValue(t *testing.T) {
	_, err := ParseConfig([]string{"size=30s", "tumbling=yes"})
	if err == nil {
		t.Fatal("expected error for invalid tumbling value")
	}
}
