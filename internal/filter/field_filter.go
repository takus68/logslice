package filter

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// FieldMatcher defines how a field value should be matched.
type FieldMatcher struct {
	Key      string
	Value    string
	Exact    bool // if false, substring match is used
}

// ByField filters log entries where the given field matches the expected value.
// If matcher.Exact is true, an exact string match is performed;
// otherwise a case-insensitive substring match is used.
func ByField(entries []parser.LogEntry, matcher FieldMatcher) []parser.LogEntry {
	var result []parser.LogEntry
	for _, entry := range entries {
		raw, ok := entry.Fields[matcher.Key]
		if !ok {
			continue
		}
		val, ok := raw.(string)
		if !ok {
			continue
		}
		if matcher.Exact {
			if val == matcher.Value {
				result = append(result, entry)
			}
		} else {
			if strings.Contains(strings.ToLower(val), strings.ToLower(matcher.Value)) {
				result = append(result, entry)
			}
		}
	}
	return result
}
