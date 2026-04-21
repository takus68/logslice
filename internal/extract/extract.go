// Package extract provides functionality to extract a subset of fields
// from log entries, producing new entries containing only the specified keys.
package extract

import (
	"fmt"
	"strings"
)

// Entry represents a single parsed log entry.
type Entry = map[string]interface{}

// Config holds the configuration for field extraction.
type Config struct {
	// Fields is the ordered list of field names to keep.
	Fields []string
	// KeepMissing, when true, includes keys with a null value if absent.
	KeepMissing bool
}

// ParseConfig parses option strings of the form:
//
//	"fields=a,b,c"
//	"keep_missing=true"
func ParseConfig(opts []string) (Config, error) {
	var cfg Config
	for _, opt := range opts {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			continue
		}
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("extract: invalid option %q: missing '='" , opt)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "fields":
			for _, f := range strings.Split(val, ",") {
				f = strings.TrimSpace(f)
				if f != "" {
					cfg.Fields = append(cfg.Fields, f)
				}
			}
		case "keep_missing":
			switch val {
			case "true":
				cfg.KeepMissing = true
			case "false":
				cfg.KeepMissing = false
			default:
				return Config{}, fmt.Errorf("extract: invalid value for keep_missing: %q", val)
			}
		default:
			return Config{}, fmt.Errorf("extract: unknown option %q", key)
		}
	}
	if len(cfg.Fields) == 0 {
		return Config{}, fmt.Errorf("extract: at least one field must be specified via 'fields='")
	}
	return cfg, nil
}

// Run returns a new slice of entries containing only the configured fields.
// If KeepMissing is true, absent fields are included with a nil value.
func Run(entries []Entry, cfg Config) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		out := make(Entry, len(cfg.Fields))
		for _, f := range cfg.Fields {
			v, ok := e[f]
			if ok {
				out[f] = v
			} else if cfg.KeepMissing {
				out[f] = nil
			}
		}
		result = append(result, out)
	}
	return result
}
