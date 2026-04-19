package merge_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/merge"
	"github.com/user/logslice/internal/parser"
)

func makeEntries(timestamps []string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(timestamps))
	for _, ts := range timestamps {
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			panic(err)
		}
		entries = append(entries, parser.Entry{
			Timestamp: t,
			Raw:       map[string]interface{}{"ts": ts},
		})
	}
	return entries
}

func TestRun_EmptyStreams(t *testing.T) {
	result := merge.Run(nil, merge.Options{})
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestRun_SingleStream(t *testing.T) {
	s := makeEntries([]string{"2024-01-01T00:00:01Z", "2024-01-01T00:00:02Z"})
	result := merge.Run([][]parser.Entry{s}, merge.Options{})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if !result[0].Timestamp.Before(result[1].Timestamp) {
		t.Error("expected ascending order")
	}
}

func TestRun_MergesTwoStreams(t *testing.T) {
	a := makeEntries([]string{"2024-01-01T00:00:01Z", "2024-01-01T00:00:03Z"})
	b := makeEntries([]string{"2024-01-01T00:00:02Z", "2024-01-01T00:00:04Z"})
	result := merge.Run([][]parser.Entry{a, b}, merge.Options{})
	if len(result) != 4 {
		t.Fatalf("expected 4, got %d", len(result))
	}
	for i := 1; i < len(result); i++ {
		if result[i].Timestamp.Before(result[i-1].Timestamp) {
			t.Errorf("out of order at index %d", i)
		}
	}
}

func TestRun_StableOrdering(t *testing.T) {
	ts := "2024-01-01T00:00:01Z"
	a := makeEntries([]string{ts})
	a[0].Raw["src"] = "a"
	b := makeEntries([]string{ts})
	b[0].Raw["src"] = "b"

	result := merge.Run([][]parser.Entry{a, b}, merge.Options{Stable: true})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0].Raw["src"] != "a" {
		t.Errorf("expected stream 0 first, got src=%v", result[0].Raw["src"])
	}
}

func TestRun_AllEmptyStreams(t *testing.T) {
	result := merge.Run([][]parser.Entry{{}, {}}, merge.Options{})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}
