// Package enrich provides functionality to add derived or static fields
// to log entries based on configurable rules.
package enrich

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule defines a single enrichment operation.
type Rule struct {
	// Key is the field name to add or overwrite.
	Key string
	// Value is the static value or a template like "{field1}+{field2}".
	Value string
}

// Run applies all enrichment rules to each entry, returning a new slice.
func Run(entries []*parser.Entry, rules []Rule) []*parser.Entry {
	result := make([]*parser.Entry, 0, len(entries))
	for _, e := range entries {
		ne := &parser.Entry{
			Timestamp: e.Timestamp,
			Raw:       e.Raw,
			Fields:    copyFields(e.Fields),
		}
		for _, r := range rules {
			ne.Fields[r.Key] = resolveValue(r.Value, ne.Fields)
		}
		result = append(result, ne)
	}
	return result
}

// resolveValue replaces {fieldName} placeholders with actual field values.
func resolveValue(template string, fields map[string]interface{}) string {
	result := template
	for k, v := range fields {
		placeholder := fmt.Sprintf("{%s}", k)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", v))
	}
	return result
}

func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
