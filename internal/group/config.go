package group

import (
	"fmt"
	"strings"
)

// ParseConfig parses a slice of option strings into a Config.
// Supported options:
//
//	field=<name>   – field to group by (required)
//	sorted=true    – sort group keys alphabetically (default: false)
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{}

	for _, opt := range opts {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			continue
		}

		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("group: invalid option %q: expected key=value", opt)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "field":
			if val == "" {
				return Config{}, fmt.Errorf("group: field name must not be empty")
			}
			cfg.Field = val
		case "sorted":
			switch val {
			case "true":
				cfg.Sorted = true
			case "false":
				cfg.Sorted = false
			default:
				return Config{}, fmt.Errorf("group: invalid value for sorted: %q", val)
			}
		default:
			return Config{}, fmt.Errorf("group: unknown option %q", key)
		}
	}

	if cfg.Field == "" {
		return Config{}, fmt.Errorf("group: field option is required")
	}

	return cfg, nil
}
