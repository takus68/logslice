// Package count provides functionality to count log entries
// matching specified field conditions.
package count

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// Config holds the configuration for a count operation.
type Config struct {
	// Field is the log entry field to match against.
	Field string
	// Value is the value to match (substring, case-insensitive).
	Value string
	// Exact requires an exact match instead of substring.
	Exact bool
}

// ParseConfig parses a count config from key=value option strings.
// Supported options: field=<name>, value=<val>, exact=true|false
func ParseConfig(opts []string) (Config, error) {
	var cfg Config
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("count: invalid option %q, expected key=value", opt)
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "field":
			cfg.Field = val
		case "value":
			cfg.Value = val
		case "exact":
			cfg.Exact = strings.EqualFold(val, "true")
		default:
			return Config{}, fmt.Errorf("count: unknown option %q", key)
		}
	}
	if cfg.Field == "" {
		return Config{}, fmt.Errorf("count: field option is required")
	}
	return cfg, nil
}

// Run counts entries matching the given Config and writes the result to w.
func Run(entries []map[string]any, cfg Config, w io.Writer) (int, error) {
	matched := filter.Apply(entries, filter.ByField(cfg.Field, cfg.Value, cfg.Exact))
	n := len(matched)
	_, err := fmt.Fprintf(w, "%d\n", n)
	return n, err
}
