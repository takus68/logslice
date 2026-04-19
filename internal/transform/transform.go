package transform

import (
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// FieldRename renames a field key in each log entry.
func FieldRename(entries []parser.Entry, from, to string) []parser.Entry {
	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		newFields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			if k == from {
				newFields[to] = v
			} else {
				newFields[k] = v
			}
		}
		result[i] = parser.Entry{Timestamp: e.Timestamp, Fields: newFields}
	}
	return result
}

// FieldDrop removes the specified field keys from each log entry.
func FieldDrop(entries []parser.Entry, keys ...string) []parser.Entry {
	dropSet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		dropSet[k] = struct{}{}
	}
	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		newFields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			if _, skip := dropSet[k]; !skip {
				newFields[k] = v
			}
		}
		result[i] = parser.Entry{Timestamp: e.Timestamp, Fields: newFields}
	}
	return result
}

// FieldNormalize lowercases all field keys in each log entry.
func FieldNormalize(entries []parser.Entry) []parser.Entry {
	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		newFields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			newFields[strings.ToLower(k)] = v
		}
		result[i] = parser.Entry{Timestamp: e.Timestamp, Fields: newFields}
	}
	return result
}
