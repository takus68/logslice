// Package template renders log entries into user-defined string templates.
// Each rule maps a destination field to a Go text/template expression that
// may reference any field from the source entry.
package template

import (
	"bytes"
	"fmt"
	"strings"
	text "text/template"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule describes a single field-to-template mapping.
type Rule struct {
	Field    string
	Template string
	tmpl     *text.Template
}

// ParseRules parses specs of the form "field=template expression".
// Multiple specs may be passed; each becomes one Rule.
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		idx := strings.Index(spec, "=")
		if idx < 0 {
			return nil, fmt.Errorf("template: missing '=' in spec %q", spec)
		}
		field := strings.TrimSpace(spec[:idx])
		expr := strings.TrimSpace(spec[idx+1:])
		if field == "" {
			return nil, fmt.Errorf("template: empty field name in spec %q", spec)
		}
		if expr == "" {
			return nil, fmt.Errorf("template: empty template expression in spec %q", spec)
		}
		t, err := text.New(field).Option("missingkey=zero").Parse(expr)
		if err != nil {
			return nil, fmt.Errorf("template: invalid template for field %q: %w", field, err)
		}
		rules = append(rules, Rule{Field: field, Template: expr, tmpl: t})
	}
	return rules, nil
}

// Run applies all rules to every entry, writing the rendered string into the
// destination field.  The original entry is not mutated.
func Run(entries []parser.LogEntry, rules []Rule) []parser.LogEntry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]parser.LogEntry, 0, len(entries))
	for _, e := range entries {
		copy := make(parser.LogEntry, len(e))
		for k, v := range e {
			copy[k] = v
		}
		for _, r := range rules {
			var buf bytes.Buffer
			if err := r.tmpl.Execute(&buf, map[string]interface{}(copy)); err == nil {
				copy[r.Field] = buf.String()
			}
		}
		out = append(out, copy)
	}
	return out
}
