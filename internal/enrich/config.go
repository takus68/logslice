package enrich

import (
	"fmt"
	"strings"
)

// ParseRules parses enrichment rule strings of the form "key=value".
// The value may contain {field} placeholders.
//
// Example:
//
//	"env=production"
//	"source={host}:{service}"
func ParseRules(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		idx := strings.IndexByte(spec, '=')
		if idx < 0 {
			return nil, fmt.Errorf("enrich: invalid rule %q: missing '='" , spec)
		}
		key := strings.TrimSpace(spec[:idx])
		val := strings.TrimSpace(spec[idx+1:])
		if key == "" {
			return nil, fmt.Errorf("enrich: invalid rule %q: empty key", spec)
		}
		rules = append(rules, Rule{Key: key, Value: val})
	}
	return rules, nil
}
