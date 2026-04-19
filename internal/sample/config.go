package sample

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseOptions parses a strategy string like "nth:5" or "random:0.25:seed=42"
// into an Options struct.
func ParseOptions(s string) (Options, error) {
	parts := strings.Split(s, ":")
	if len(parts) == 0 || parts[0] == "" {
		return Options{}, fmt.Errorf("empty strategy")
	}

	opts := Options{Strategy: Strategy(parts[0])}

	switch opts.Strategy {
	case StrategyNth:
		if len(parts) < 2 {
			return Options{}, fmt.Errorf("nth strategy requires N, e.g. nth:5")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n <= 0 {
			return Options{}, fmt.Errorf("nth: invalid N %q", parts[1])
		}
		opts.N = n

	case StrategyRandom:
		if len(parts) < 2 {
			return Options{}, fmt.Errorf("random strategy requires rate, e.g. random:0.5")
		}
		rate, err := strconv.ParseFloat(parts[1], 64)
		if err != nil || rate < 0 || rate > 1 {
			return Options{}, fmt.Errorf("random: invalid rate %q", parts[1])
		}
		opts.Rate = rate
		if len(parts) >= 3 {
			seedStr := strings.TrimPrefix(parts[2], "seed=")
			seed, err := strconv.ParseInt(seedStr, 10, 64)
			if err != nil {
				return Options{}, fmt.Errorf("random: invalid seed %q", parts[2])
			}
			opts.Seed = seed
		}

	default:
		return Options{}, fmt.Errorf("unknown strategy %q", opts.Strategy)
	}

	return opts, nil
}
