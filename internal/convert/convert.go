// Package convert provides field type conversion for log entries.
// It supports converting string fields to int, float, bool, or string types.
package convert

import (
	"fmt"
	"strconv"

	"github.com/yourorg/logslice/internal/parser"
)

// TargetType represents the type to convert a field value to.
type TargetType string

const (
	TypeInt    TargetType = "int"
	TypeFloat  TargetType = "float"
	TypeBool   TargetType = "bool"
	TypeString TargetType = "string"
)

// Rule defines a single field conversion rule.
type Rule struct {
	Field string
	Type  TargetType
}

// Run applies all conversion rules to each log entry.
// Fields that cannot be converted are left unchanged.
func Run(entries []*parser.LogEntry, rules []Rule) []*parser.LogEntry {
	result := make([]*parser.LogEntry, len(entries))
	for i, e := range entries {
		result[i] = convertEntry(e, rules)
	}
	return result
}

func convertEntry(e *parser.LogEntry, rules []Rule) *parser.LogEntry {
	newFields := make(map[string]interface{}, len(e.Fields))
	for k, v := range e.Fields {
		newFields[k] = v
	}
	for _, r := range rules {
		val, ok := newFields[r.Field]
		if !ok {
			continue
		}
		converted, err := convertValue(val, r.Type)
		if err != nil {
			continue
		}
		newFields[r.Field] = converted
	}
	return &parser.LogEntry{
		Timestamp: e.Timestamp,
		Fields:    newFields,
		Raw:       e.Raw,
	}
}

func convertValue(val interface{}, t TargetType) (interface{}, error) {
	s := fmt.Sprintf("%v", val)
	switch t {
	case TypeInt:
		return strconv.ParseInt(s, 10, 64)
	case TypeFloat:
		return strconv.ParseFloat(s, 64)
	case TypeBool:
		return strconv.ParseBool(s)
	case TypeString:
		return s, nil
	}
	return nil, fmt.Errorf("unknown type: %s", t)
}
