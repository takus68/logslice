package filter

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func makeTimeEntries() []parser.Entry {
	parse := func(s string) time.Time {
		t, _ := time.Parse(time.RFC3339, s)
		return t
	}
	return []parser.Entry{
		{Timestamp: parse("2024-01-01T08:00:00Z"), Fields: map[string]interface{}{"msg": "early"}},
		{Timestamp: parse("2024-01-01T10:00:00Z"), Fields: map[string]interface{}{"msg": "in-range-start"}},
		{Timestamp: parse("2024-01-01T12:00:00Z"), Fields: map[string]interface{}{"msg": "in-range-mid"}},
		{Timestamp: parse("2024-01-01T14:00:00Z"), Fields: map[string]interface{}{"msg": "in-range-end"}},
		{Timestamp: parse("2024-01-01T16:00:00Z"), Fields: map[string]interface{}{"msg": "late"}},
	}
}

func TestByTimeRange_IncludesBoundaries(t *testing.T) {
	start, _ := time.Parse(time.RFC3339, "2024-01-01T10:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2024-01-01T14:00:00Z")

	result := ByTimeRange(makeTimeEntries(), start, end)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].Fields["msg"] != "in-range-start" {
		t.Errorf("unexpected first entry: %v", result[0].Fields["msg"])
	}
	if result[2].Fields["msg"] != "in-range-end" {
		t.Errorf("unexpected last entry: %v", result[2].Fields["msg"])
	}
}

func TestByTimeRange_EmptyResult(t *testing.T) {
	start, _ := time.Parse(time.RFC3339, "2024-01-02T00:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2024-01-02T23:59:59Z")

	result := ByTimeRange(makeTimeEntries(), start, end)
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestByTimeRange_AllEntries(t *testing.T) {
	start, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2024-01-01T23:59:59Z")

	result := ByTimeRange(makeTimeEntries(), start, end)
	if len(result) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(result))
	}
}

func TestByTimeRange_ZeroEnd(t *testing.T) {
	start, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")

	result := ByTimeRange(makeTimeEntries(), start, time.Time{})
	if len(result) != 3 {
		t.Fatalf("expected 3 entries (from mid onwards), got %d", len(result))
	}
}
