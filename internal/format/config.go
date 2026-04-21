package format

import (
	"fmt"
	"strings"
)

var validFormats = map[string]bool{
	"upper": true,
	"lower": true,
	"title": true,
	"trim":  true,
	"quote": true,
}

// ParseRules parses formatting rule specs of the form "field=format".
// Example: ["level=upper", "message=trim"]
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("format: invalid rule %q: expected field=format", spec)
		}
		field := strings.TrimSpace(parts[0])
		fmt_ := strings.TrimSpace(strings.ToLower(parts[1]))

		if field == "" {
			return nil, fmt.Errorf("format: empty field name in rule %q", spec)
		}
		if !validFormats[fmt_] {
			return nil, fmt.Errorf("format: unknown format %q in rule %q (valid: upper, lower, title, trim, quote)", fmt_, spec)
		}
		rules = append(rules, Rule{Field: field, Format: fmt_})
	}
	return rules, nil
}
