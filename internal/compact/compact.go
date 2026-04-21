// Package compact removes nil/empty fields from log entries.
package compact

import (
	"fmt"
	"strings"
)

// Config controls which fields are considered empty.
type Config struct {
	// RemoveEmpty removes fields whose string representation is "".
	RemoveEmpty bool
	// RemoveNull removes fields with nil values.
	RemoveNull bool
	// Fields restricts compaction to a specific set of keys; empty means all keys.
	Fields []string
}

// ParseConfig parses key=value option strings into a Config.
// Supported options: remove_empty=true|false, remove_null=true|false, fields=a,b,c
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{RemoveNull: true}
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("compact: invalid option %q, expected key=value", opt)
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "remove_empty":
			cfg.RemoveEmpty = val == "true"
		case "remove_null":
			cfg.RemoveNull = val == "true"
		case "fields":
			for _, f := range strings.Split(val, ",") {
				f = strings.TrimSpace(f)
				if f != "" {
					cfg.Fields = append(cfg.Fields, f)
				}
			}
		default:
			return Config{}, fmt.Errorf("compact: unknown option %q", key)
		}
	}
	return cfg, nil
}

// Run removes empty or null fields from each entry according to cfg.
func Run(entries []map[string]any, cfg Config) []map[string]any {
	result := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		result = append(result, compact(entry, cfg))
	}
	return result
}

func compact(entry map[string]any, cfg Config) map[string]any {
	out := make(map[string]any, len(entry))
	for k, v := range entry {
		if len(cfg.Fields) > 0 && !contains(cfg.Fields, k) {
			out[k] = v
			continue
		}
		if cfg.RemoveNull && v == nil {
			continue
		}
		if cfg.RemoveEmpty {
			if s, ok := v.(string); ok && s == "" {
				continue
			}
		}
		out[k] = v
	}
	return out
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
