// Package limit provides functionality to cap the number of log entries
// returned from a pipeline, useful for previewing large log files.
package limit

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Config holds the configuration for the limit operation.
type Config struct {
	// Max is the maximum number of entries to return. Zero means no limit.
	Max int
	// Offset is the number of entries to skip before collecting.
	Offset int
}

// ParseConfig parses options of the form "max=N" and "offset=N".
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{}
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("limit: invalid option %q, expected key=value", opt)
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		n, err := strconv.Atoi(val)
		if err != nil {
			return Config{}, fmt.Errorf("limit: option %q has non-integer value %q", key, val)
		}
		if n < 0 {
			return Config{}, fmt.Errorf("limit: option %q must be non-negative, got %d", key, n)
		}
		switch key {
		case "max":
			cfg.Max = n
		case "offset":
			cfg.Offset = n
		default:
			return Config{}, fmt.Errorf("limit: unknown option %q", key)
		}
	}
	return cfg, nil
}

// Run applies the limit and offset to the given entries.
// If cfg.Max is zero, all entries after the offset are returned.
func Run(entries []*parser.Entry, cfg Config) []*parser.Entry {
	if cfg.Offset >= len(entries) {
		return []*parser.Entry{}
	}
	sliced := entries[cfg.Offset:]
	if cfg.Max == 0 || cfg.Max >= len(sliced) {
		return sliced
	}
	return sliced[:cfg.Max]
}
