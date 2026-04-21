// Package unique provides functionality to extract unique values
// for a given field across a set of log entries.
package unique

import (
	"fmt"
	"sort"
	"strings"
)

// Entry represents a single parsed log entry.
type Entry = map[string]interface{}

// Config holds options for the unique operation.
type Config struct {
	Field  string
	Sorted bool
}

// ParseConfig parses key=value option strings into a Config.
// Supported options:
//
//	field=<name>   — the field to extract unique values from (required)
//	sorted=true    — whether to sort the output (default: false)
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{}
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("unique: invalid option %q, expected key=value", opt)
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "field":
			if val == "" {
				return Config{}, fmt.Errorf("unique: field value must not be empty")
			}
			cfg.Field = val
		case "sorted":
			switch val {
			case "true":
				cfg.Sorted = true
			case "false":
				cfg.Sorted = false
			default:
				return Config{}, fmt.Errorf("unique: invalid value for sorted: %q", val)
			}
		default:
			return Config{}, fmt.Errorf("unique: unknown option %q", key)
		}
	}
	if cfg.Field == "" {
		return Config{}, fmt.Errorf("unique: field option is required")
	}
	return cfg, nil
}

// Run returns the distinct values found for cfg.Field across entries.
// Values are returned as strings. Entries missing the field are skipped.
func Run(entries []Entry, cfg Config) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, e := range entries {
		v, ok := e[cfg.Field]
		if !ok {
			continue
		}
		s := fmt.Sprintf("%v", v)
		if _, exists := seen[s]; !exists {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	if cfg.Sorted {
		sort.Strings(result)
	}
	return result
}
