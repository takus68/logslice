// Package highlight provides field-based terminal highlighting for log entries.
package highlight

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Color ANSI codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

// Rule defines a highlight rule: if field equals (or contains) value, apply color.
type Rule struct {
	Field     string
	Value     string
	Color     string
	Substring bool
}

// Apply applies highlight rules to a log entry, returning a colorized string.
func Apply(entry parser.Entry, rules []Rule) string {
	color := ""
	for _, r := range rules {
		v, ok := entry.Fields[r.Field]
		if !ok {
			continue
		}
		s := fmt.Sprintf("%v", v)
		matched := false
		if r.Substring {
			matched = strings.Contains(strings.ToLower(s), strings.ToLower(r.Value))
		} else {
			matched = strings.EqualFold(s, r.Value)
		}
		if matched {
			color = r.Color
			break
		}
	}
	line := entry.Raw
	if color != "" {
		line = color + line + Reset
	}
	return line
}

// ApplyAll applies rules to all entries, returning colorized lines.
func ApplyAll(entries []parser.Entry, rules []Rule) []string {
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		out = append(out, Apply(e, rules))
	}
	return out
}
