// Package rename provides functionality to rename keys in log entries
// based on a set of configurable mapping rules.
package rename

import "github.com/yourorg/logslice/internal/parser"

// Rule maps an old field name to a new field name.
type Rule struct {
	From string
	To   string
}

// Run applies the given rename rules to each entry in the slice.
// If the From field exists, it is moved to the To field and the
// original key is removed. Entries that do not contain a matched
// key are returned unchanged.
func Run(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 || len(entries) == 0 {
		return entries
	}

	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		result = append(result, applyRules(e, rules))
	}
	return result
}

func applyRules(e parser.Entry, rules []Rule) parser.Entry {
	fields := make(map[string]interface{}, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}

	for _, r := range rules {
		if val, ok := fields[r.From]; ok {
			fields[r.To] = val
			delete(fields, r.From)
		}
	}

	return parser.Entry{
		Timestamp: e.Timestamp,
		Raw:       e.Raw,
		Fields:    fields,
	}
}
