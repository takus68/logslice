package coalesce

import (
	"github.com/logslice/logslice/internal/parser"
)

// Run applies coalesce rules to each entry, returning new entries with the
// destination fields populated from the first non-empty source field.
func Run(entries []*parser.Entry, rules []Rule) []*parser.Entry {
	out := make([]*parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = copyEntry(e)
		for _, rule := range rules {
			val := firstNonEmpty(e.Fields, rule.Sources)
			if val != nil {
				out[i].Fields[rule.Dest] = val
			}
		}
	}
	return out
}

// firstNonEmpty returns the first field value from fields that is non-nil and
// non-empty string, or nil if none qualify.
func firstNonEmpty(fields map[string]any, sources []string) any {
	for _, src := range sources {
		v, ok := fields[src]
		if !ok || v == nil {
			continue
		}
		if s, isStr := v.(string); isStr && s == "" {
			continue
		}
		return v
	}
	return nil
}

// copyEntry creates a shallow copy of an entry's Fields map.
func copyEntry(e *parser.Entry) *parser.Entry {
	newFields := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		newFields[k] = v
	}
	return &parser.Entry{
		Timestamp: e.Timestamp,
		Raw:       e.Raw,
		Fields:    newFields,
	}
}
