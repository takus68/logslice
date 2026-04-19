// Package truncate provides field-level value truncation for log entries.
package truncate

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Config holds truncation settings.
type Config struct {
	Fields    []string
	MaxLength int
	Suffix    string
}

// Run truncates specified fields in each entry to MaxLength characters.
// If a field value exceeds MaxLength, it is trimmed and Suffix is appended.
func Run(entries []*parser.Entry, cfg Config) []*parser.Entry {
	if len(cfg.Fields) == 0 || cfg.MaxLength <= 0 {
		return entries
	}
	suffix := cfg.Suffix
	if suffix == "" {
		suffix = "..."
	}
	result := make([]*parser.Entry, len(entries))
	for i, e := range entries {
		fields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		for _, f := range cfg.Fields {
			val, ok := fields[f]
			if !ok {
				continue
			}
			s := fmt.Sprintf("%v", val)
			if len(s) > cfg.MaxLength {
				s = s[:cfg.MaxLength] + suffix
			}
			fields[f] = s
		}
		result[i] = &parser.Entry{Timestamp: e.Timestamp, Fields: fields}
	}
	return result
}

// ParseConfig parses key=value options for truncation.
// Supported keys: fields (comma-separated), max (int), suffix (string).
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{MaxLength: 80, Suffix: "..."}
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return cfg, fmt.Errorf("truncate: invalid option %q, expected key=value", opt)
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "fields":
			for _, f := range strings.Split(val, ",") {
				f = strings.TrimSpace(f)
				if f != "" {
					cfg.Fields = append(cfg.Fields, f)
				}
			}
		case "max":
			n := 0
			if _, err := fmt.Sscanf(val, "%d", &n); err != nil || n <= 0 {
				return cfg, fmt.Errorf("truncate: invalid max value %q", val)
			}
			cfg.MaxLength = n
		case "suffix":
			cfg.Suffix = val
		default:
			return cfg, fmt.Errorf("truncate: unknown option %q", key)
		}
	}
	return cfg, nil
}
