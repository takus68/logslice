// Package redact provides field redaction for log entries.
package redact

import (
	"regexp"
	"strings"

	"github.com/yourusername/logslice/internal/parser"
)

// Config holds redaction configuration.
type Config struct {
	// Fields to fully redact (value replaced with "[REDACTED]").
	Fields []string
	// Patterns maps field names to regex patterns; matching values are masked.
	Patterns map[string]*regexp.Regexp
}

// Run applies redaction rules to a slice of log entries, returning new entries.
func Run(entries []parser.Entry, cfg Config) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		result = append(result, redactEntry(e, cfg))
	}
	return result
}

func redactEntry(e parser.Entry, cfg Config) parser.Entry {
	newRaw := make(map[string]interface{}, len(e.Raw))
	for k, v := range e.Raw {
		newRaw[k] = v
	}

	for _, field := range cfg.Fields {
		if _, ok := newRaw[field]; ok {
			newRaw[field] = "[REDACTED]"
		}
	}

	for field, re := range cfg.Patterns {
		if val, ok := newRaw[field]; ok {
			if s, ok := val.(string); ok {
				newRaw[field] = maskString(s, re)
			}
		}
	}

	return parser.Entry{Timestamp: e.Timestamp, Raw: newRaw}
}

// maskString replaces matched groups or full matches with "***".
func maskString(s string, re *regexp.Regexp) string {
	return re.ReplaceAllStringFunc(s, func(match string) string {
		return strings.Repeat("*", len(match))
	})
}
