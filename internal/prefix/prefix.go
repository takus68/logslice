// Package prefix provides functionality to add a static or dynamic
// string prefix to a specified field value in log entries.
package prefix

import (
	"fmt"
	"strings"
)

// Rule defines a single prefix operation: add Prefix to the value of Field.
type Rule struct {
	Field  string
	Prefix string
}

// ParseRules parses a slice of spec strings of the form "field=prefix".
// Example: "message=INFO: "
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		idx := strings.Index(spec, "=")
		if idx < 0 {
			return nil, fmt.Errorf("prefix: missing '=' in spec %q", spec)
		}
		field := strings.TrimSpace(spec[:idx])
		pfx := spec[idx+1:]
		if field == "" {
			return nil, fmt.Errorf("prefix: empty field name in spec %q", spec)
		}
		rules = append(rules, Rule{Field: field, Prefix: pfx})
	}
	return rules, nil
}

// Run applies the given prefix rules to each entry.
// For each entry a shallow copy is made so originals are not mutated.
// If a field is missing or its value is not a string, the entry is left unchanged.
func Run(entries []map[string]interface{}, rules []Rule) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		copy := copyEntry(e)
		for _, r := range rules {
			v, ok := copy[r.Field]
			if !ok {
				continue
			}
			s, ok := v.(string)
			if !ok {
				continue
			}
			copy[r.Field] = r.Prefix + s
		}
		result = append(result, copy)
	}
	return result
}

func copyEntry(e map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(e))
	for k, v := range e {
		out[k] = v
	}
	return out
}
