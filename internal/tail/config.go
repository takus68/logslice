package tail

import (
	"fmt"
	"strconv"
	"strings"
)

// Mode indicates whether to take from the head or tail.
type Mode int

const (
	ModeTail Mode = iota
	ModeHead
)

// Config holds parsed CLI options for the tail/head feature.
type Config struct {
	Mode Mode
	N    int
}

// ParseConfig parses a raw option string of the form "tail=N" or "head=N".
func ParseConfig(raw string) (Config, error) {
	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 {
		return Config{}, fmt.Errorf("tail: invalid option %q, expected head=N or tail=N", raw)
	}
	key := strings.TrimSpace(strings.ToLower(parts[0]))
	val := strings.TrimSpace(parts[1])

	n, err := strconv.Atoi(val)
	if err != nil || n < 0 {
		return Config{}, fmt.Errorf("tail: N must be a non-negative integer, got %q", val)
	}

	switch key {
	case "tail":
		return Config{Mode: ModeTail, N: n}, nil
	case "head":
		return Config{Mode: ModeHead, N: n}, nil
	default:
		return Config{}, fmt.Errorf("tail: unknown mode %q, expected head or tail", key)
	}
}
