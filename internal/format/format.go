// Package format provides field value formatting transformations for log entries.
package format

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule describes how to format a single field.
type Rule struct {
	Field  string
	Format string // e.g. "upper", "lower", "title", "trim", "quote"
}

// Run applies all formatting rules to each log entry.
// Entries that do not contain a referenced field are left unchanged.
func Run(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 {
		return entries
	}

	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		result = append(result, applyRules(e, rules))
	}
	return result
}

func applyRules(e parser.Entry, rules []Rule) parser.Entry {
	copy := make(map[string]interface{}, len(e.Fields))
	for k, v := range e.Fields {
		copy[k] = v
	}

	for _, r := range rules {
		val, ok := copy[r.Field]
		if !ok {
			continue
		}
		s := fmt.Sprintf("%v", val)
		copy[r.Field] = applyFormat(s, r.Format)
	}

	return parser.Entry{Timestamp: e.Timestamp, Fields: copy, Raw: e.Raw}
}

func applyFormat(s, format string) string {
	switch strings.ToLower(format) {
	case "upper":
		return strings.ToUpper(s)
	case "lower":
		return strings.ToLower(s)
	case "title":
		return strings.Title(strings.ToLower(s)) //nolint:staticcheck
	case "trim":
		return strings.TrimSpace(s)
	case "quote":
		return fmt.Sprintf("%q", s)
	default:
		return s
	}
}
