// Package sort provides log entry sorting by a specified field.
package sort

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Config holds sorting configuration.
type Config struct {
	Field     string
	Descending bool
}

// ParseConfig parses options like "field=timestamp", "order=desc".
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{Field: "timestamp"}
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("sort: invalid option %q, expected key=value", opt)
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "field":
			if val == "" {
				return Config{}, fmt.Errorf("sort: field name must not be empty")
			}
			cfg.Field = val
		case "order":
			switch strings.ToLower(val) {
			case "asc":
				cfg.Descending = false
			case "desc":
				cfg.Descending = true
			default:
				return Config{}, fmt.Errorf("sort: unknown order %q, expected asc or desc", val)
			}
		default:
			return Config{}, fmt.Errorf("sort: unknown option %q", key)
		}
	}
	return cfg, nil
}

// Run sorts entries by the configured field using lexicographic comparison.
// Entries missing the field are placed at the end.
func Run(entries []parser.Entry, cfg Config) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)

	sort.SliceStable(out, func(i, j int) bool {
		vi, oki := out[i].Fields[cfg.Field]
		vj, okj := out[j].Fields[cfg.Field]

		// Missing field sorts to the end.
		if !oki && !okj {
			return false
		}
		if !oki {
			return false
		}
		if !okj {
			return true
		}

		si := fmt.Sprintf("%v", vi)
		sj := fmt.Sprintf("%v", vj)

		if cfg.Descending {
			return si > sj
		}
		return si < sj
	})
	return out
}
