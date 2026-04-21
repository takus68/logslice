package window

import (
	"fmt"
	"strings"
	"time"
)

// ParseConfig parses key=value option strings into a Config.
//
// Supported options:
//
//	size=<duration>   required, e.g. size=1m, size=30s
//	field=<name>      optional, timestamp field (default: "time")
//	tumbling=true     optional, use tumbling windows (default: true)
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{
		Field:    "time",
		Tumbling: true,
	}

	for _, opt := range opts {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			continue
		}
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("window: invalid option %q, expected key=value", opt)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "size":
			d, err := time.ParseDuration(val)
			if err != nil {
				return Config{}, fmt.Errorf("window: invalid size %q: %w", val, err)
			}
			if d <= 0 {
				return Config{}, fmt.Errorf("window: size must be positive")
			}
			cfg.Size = d
		case "field":
			if val == "" {
				return Config{}, fmt.Errorf("window: field value must not be empty")
			}
			cfg.Field = val
		case "tumbling":
			switch val {
			case "true":
				cfg.Tumbling = true
			case "false":
				cfg.Tumbling = false
			default:
				return Config{}, fmt.Errorf("window: tumbling must be true or false, got %q", val)
			}
		default:
			return Config{}, fmt.Errorf("window: unknown option %q", key)
		}
	}

	if cfg.Size == 0 {
		return Config{}, fmt.Errorf("window: size is required")
	}
	return cfg, nil
}
