package sample

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(n int) []parser.Entry {
	entries := make([]parser.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = parser.Entry{
			Timestamp: time.Now(),
			Raw:       map[string]interface{}{"index": i},
		}
	}
	return entries
}

func TestNth_EverySecond(t *testing.T) {
	entries := makeEntries(10)
	result := Run(entries, Options{Strategy: StrategyNth, N: 2})
	if len(result) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(result))
	}
}

func TestNth_ZeroN(t *testing.T) {
	entries := makeEntries(6)
	result := Run(entries, Options{Strategy: StrategyNth, N: 0})
	if len(result) != 6 {
		t.Fatalf("expected all entries, got %d", len(result))
	}
}

func TestNth_One(t *testing.T) {
	entries := makeEntries(5)
	result := Run(entries, Options{Strategy: StrategyNth, N: 1})
	if len(result) != 5 {
		t.Fatalf("expected 5, got %d", len(result))
	}
}

func TestRandom_RateZero(t *testing.T) {
	entries := makeEntries(20)
	result := Run(entries, Options{Strategy: StrategyRandom, Rate: 0, Seed: 42})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestRandom_RateOne(t *testing.T) {
	entries := makeEntries(10)
	result := Run(entries, Options{Strategy: StrategyRandom, Rate: 1.0, Seed: 1})
	if len(result) != 10 {
		t.Fatalf("expected 10 entries, got %d", len(result))
	}
}

func TestRandom_Deterministic(t *testing.T) {
	entries := makeEntries(100)
	r1 := Run(entries, Options{Strategy: StrategyRandom, Rate: 0.5, Seed: 99})
	r2 := Run(entries, Options{Strategy: StrategyRandom, Rate: 0.5, Seed: 99})
	if len(r1) != len(r2) {
		t.Fatalf("expected deterministic results, got %d vs %d", len(r1), len(r2))
	}
}

func TestUnknownStrategy(t *testing.T) {
	entries := makeEntries(4)
	result := Run(entries, Options{Strategy: "unknown"})
	if len(result) != 4 {
		t.Fatalf("expected passthrough, got %d", len(result))
	}
}
