// Package sample provides log entry sampling strategies.
package sample

import (
	"math/rand"

	"github.com/yourorg/logslice/internal/parser"
)

// Strategy defines how entries are sampled.
type Strategy string

const (
	StrategyNth    Strategy = "nth"
	StrategyRandom Strategy = "random"
)

// Options configures the sampler.
type Options struct {
	Strategy Strategy
	// N is used by StrategyNth: keep every Nth entry.
	N int
	// Rate is used by StrategyRandom: probability [0,1] of keeping an entry.
	Rate float64
	// Seed for random sampling; 0 means use default source.
	Seed int64
}

// Run applies the sampling strategy to entries and returns the sampled subset.
func Run(entries []parser.Entry, opts Options) []parser.Entry {
	switch opts.Strategy {
	case StrategyNth:
		return nthSample(entries, opts.N)
	case StrategyRandom:
		return randomSample(entries, opts.Rate, opts.Seed)
	default:
		return entries
	}
}

func nthSample(entries []parser.Entry, n int) []parser.Entry {
	if n <= 0 {
		return entries
	}
	out := make([]parser.Entry, 0, len(entries)/n+1)
	for i, e := range entries {
		if i%n == 0 {
			out = append(out, e)
		}
	}
	return out
}

func randomSample(entries []parser.Entry, rate float64, seed int64) []parser.Entry {
	if rate <= 0 {
		return nil
	}
	if rate >= 1 {
		return entries
	}
	r := rand.New(rand.NewSource(seed))
	out := make([]parser.Entry, 0, int(float64(len(entries))*rate)+1)
	for _, e := range entries {
		if r.Float64() < rate {
			out = append(out, e)
		}
	}
	return out
}
