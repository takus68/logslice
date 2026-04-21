// Package coalesce provides functionality to merge or fall back across
// multiple fields, returning the first non-empty value found.
//
// This is useful when log entries may use different field names for the
// same semantic value (e.g. "msg", "message", "text") and you want to
// normalize them into a single canonical field.
package coalesce

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Rule describes a coalesce operation: read the first non-empty value
// from Sources and write it to Target.
type Rule struct {
	Target  string
	Sources []string
}

// ParseRules parses coalesce rule specs of the form:
//
//	target=src1,src2,src3
//
// Multiple rules may be provided. Whitespace around tokens is trimmed.
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("coalesce: invalid rule %q: expected target=src1,src2,...", spec)
		}
		target := strings.TrimSpace(parts[0])
		if target == "" {
			return nil, fmt.Errorf("coalesce: invalid rule %q: target field must not be empty", spec)
		}
		rawSources := strings.Split(parts[1], ",")
		var sources []string
		for _, s := range rawSources {
			s = strings.TrimSpace(s)
			if s != "" {
				sources = append(sources, s)
			}
		}
		if len(sources) == 0 {
			return nil, fmt.Errorf("coalesce: invalid rule %q: at least one source field is required", spec)
		}
		rules = append(rules, Rule{Target: target, Sources: sources})
	}
	return rules, nil
}

// Run applies each Rule to every entry. For each rule, the first source
// field that exists and has a non-empty string value is written to the
// target field. If no source yields a value, the target field is left
// unchanged (or absent). The original entries are not mutated.
func Run(entries []parser.LogEntry, rules []Rule) []parser.LogEntry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]parser.LogEntry, len(entries))
	for i, entry := range entries {
		fields := make(map[string]interface{}, len(entry.Fields))
		for k, v := range entry.Fields {
			fields[k] = v
		}
		for _, rule := range rules {
			for _, src := range rule.Sources {
				v, ok := fields[src]
				if !ok {
					continue
				}
				s, isStr := v.(string)
				if isStr && strings.TrimSpace(s) == "" {
					continue
				}
				// Non-string non-nil values are also accepted.
				fields[rule.Target] = v
				break
			}
		}
		out[i] = parser.LogEntry{
			Timestamp: entry.Timestamp,
			Raw:       entry.Raw,
			Fields:    fields,
		}
	}
	return out
}
