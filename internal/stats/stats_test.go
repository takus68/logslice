package stats

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries() []parser.Entry {
	t1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	return []parser.Entry{
		{Timestamp: t1, Fields: map[string]string{"level": "info", "msg": "start"}},
		{Timestamp: t2, Fields: map[string]string{"level": "warn", "msg": "slow"}},
		{Timestamp: t3, Fields: map[string]string{"level": "info", "msg": "done"}},
	}
}

func TestCompute_Total(t *testing.T) {
	entries := makeEntries()
	s := Compute(entries)
	if s.Total != 3 {
		t.Errorf("expected 3, got %d", s.Total)
	}
}

func TestCompute_TimeRange(t *testing.T) {
	entries := makeEntries()
	s := Compute(entries)
	expectedEarliest := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	expectedLatest := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	if !s.Earliest.Equal(expectedEarliest) {
		t.Errorf("earliest mismatch: %v", s.Earliest)
	}
	if !s.Latest.Equal(expectedLatest) {
		t.Errorf("latest mismatch: %v", s.Latest)
	}
}

func TestCompute_LevelCounts(t *testing.T) {
	entries := makeEntries()
	s := Compute(entries)
	if s.LevelCounts["info"] != 2 {
		t.Errorf("expected 2 info, got %d", s.LevelCounts["info"])
	}
	if s.LevelCounts["warn"] != 1 {
		t.Errorf("expected 1 warn, got %d", s.LevelCounts["warn"])
	}
}

func TestCompute_Empty(t *testing.T) {
	s := Compute(nil)
	if s.Total != 0 {
		t.Errorf("expected 0, got %d", s.Total)
	}
	if len(s.LevelCounts) != 0 {
		t.Error("expected empty level counts")
	}
}

func TestCompute_FieldKeys(t *testing.T) {
	entries := makeEntries()
	s := Compute(entries)
	if len(s.FieldKeys) == 0 {
		t.Error("expected non-empty field keys")
	}
}
