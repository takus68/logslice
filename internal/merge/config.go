package merge

import (
	"fmt"
	"strings"
)

// Config holds parsed CLI options for the merge command.
type Config struct {
	Files  []string
	Stable bool
}

// ParseConfig parses a flat key=value option string and a list of file paths
// into a Config.
//
// Supported options:
//
//	stable=true   – preserve source-stream order for equal timestamps
func ParseConfig(files []string, opts string) (Config, error) {
	cfg := Config{Files: files}
	if opts == "" {
		return cfg, nil
	}
	for _, part := range strings.Split(opts, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return Config{}, fmt.Errorf("merge: invalid option %q (expected key=value)", part)
		}
		key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		switch key {
		case "stable":
			switch val {
			case "true":
				cfg.Stable = true
			case "false":
				cfg.Stable = false
			default:
				return Config{}, fmt.Errorf("merge: stable must be true or false, got %q", val)
			}
		default:
			return Config{}, fmt.Errorf("merge: unknown option %q", key)
		}
	}
	return cfg, nil
}
