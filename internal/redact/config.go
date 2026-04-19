package redact

import (
	"fmt"
	"regexp"
	"strings"
)

// Options holds raw CLI/config input for redaction.
type Options struct {
	// RedactFields is a comma-separated list of field names to fully redact.
	RedactFields string
	// MaskPatterns is a slice of "field=pattern" strings.
	MaskPatterns []string
}

// ParseConfig converts Options into a Config ready for use.
func ParseConfig(opts Options) (Config, error) {
	cfg := Config{
		Patterns: make(map[string]*regexp.Regexp),
	}

	if opts.RedactFields != "" {
		for _, f := range strings.Split(opts.RedactFields, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				cfg.Fields = append(cfg.Fields, f)
			}
		}
	}

	for _, mp := range opts.MaskPatterns {
		parts := strings.SplitN(mp, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("invalid mask pattern %q: expected field=pattern", mp)
		}
		field := strings.TrimSpace(parts[0])
		pattern := strings.TrimSpace(parts[1])
		re, err := regexp.Compile(pattern)
		if err != nil {
			return Config{}, fmt.Errorf("invalid regex for field %q: %w", field, err)
		}
		cfg.Patterns[field] = re
	}

	return cfg, nil
}
