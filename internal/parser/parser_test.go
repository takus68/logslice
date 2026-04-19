package parser

import (
	"strings"
	"testing"
	"time"
)

func TestParse_ValidEntries(t *testing.T) {
	input := `{"time":"2024-01-15T10:00:00Z","level":"info","msg":"started"}
{"time":"2024-01-15T10:01:00Z","level":"error","msg":"failed"}
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	expected := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !entries[0].Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, entries[0].Timestamp)
	}
	if entries[1].Fields["level"] != "error" {
		t.Errorf("expected level=error, got %v", entries[1].Fields["level"])
	}
}

func TestParse_SkipsEmptyLines(t *testing.T) {
	input := `{"time":"2024-01-15T10:00:00Z","msg":"a"}

{"time":"2024-01-15T10:01:00Z","msg":"b"}
`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	input := `not json`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParse_MissingTimestamp(t *testing.T) {
	input := `{"level":"info","msg":"no time"}`
	entries, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !entries[0].Timestamp.IsZero() {
		t.Errorf("expected zero timestamp for entry without time field")
	}
}

func TestParse_InvalidTimestamp(t *testing.T) {
	input := `{"time":"not-a-date","msg":"bad"}`
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for invalid timestamp")
	}
}
