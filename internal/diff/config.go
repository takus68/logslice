package diff

import (
	"fmt"
	"strings"
)

// Config holds parsed options for the diff operation.
type Config struct {
	KeyField string
	Mode     Mode
}

// ParseConfig parses key=value option strings into a Config.
// Supported options:
//   - key=<fieldName>  (required)
//   - mode=all|added|removed|changed  (default: all)
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{
		Mode: ModeAll,
	}

	for _, opt := range opts {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			continue
		}
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("diff: invalid option %q, expected key=value", opt)
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		switch k {
		case "key":
			if v == "" {
				return Config{}, fmt.Errorf("diff: key field must not be empty")
			}
			cfg.KeyField = v
		case "mode":
			switch Mode(v) {
			case ModeAll, ModeAdded, ModeRemoved, ModeChanged:
				cfg.Mode = Mode(v)
			default:
				return Config{}, fmt.Errorf("diff: unknown mode %q, expected all|added|removed|changed", v)
			}
		default:
			return Config{}, fmt.Errorf("diff: unknown option %q", k)
		}
	}

	if cfg.KeyField == "" {
		return Config{}, fmt.Errorf("diff: option 'key' is required")
	}

	return cfg, nil
}
