// Package mask provides field-level value masking for log entries,
// supporting full replacement and partial masking strategies.
package mask

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Strategy defines how a field value is masked.
type Strategy string

const (
	StrategyFull    Strategy = "full"    // replace entire value with placeholder
	StrategyPartial Strategy = "partial" // keep first/last N chars, mask middle
)

// Rule describes a masking rule for a single field.
type Rule struct {
	Field       string
	Strategy    Strategy
	Placeholder string // used for full masking, default "***"
	KeepPrefix  int    // chars to keep at start for partial masking
	KeepSuffix  int    // chars to keep at end for partial masking
}

// Run applies all masking rules to every entry and returns new entries.
func Run(entries []parser.Entry, rules []Rule) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		result = append(result, applyRules(e, rules))
	}
	return result
}

func applyRules(e parser.Entry, rules []Rule) parser.Entry {
	copy := parser.Entry{
		Timestamp: e.Timestamp,
		Raw:       make(map[string]interface{}, len(e.Raw)),
	}
	for k, v := range e.Raw {
		copy.Raw[k] = v
	}
	for _, r := range rules {
		v, ok := copy.Raw[r.Field]
		if !ok {
			continue
		}
		copy.Raw[r.Field] = maskValue(fmt.Sprintf("%v", v), r)
	}
	return copy
}

func maskValue(s string, r Rule) string {
	switch r.Strategy {
	case StrategyPartial:
		ph := r.Placeholder
		if ph == "" {
			ph = "***"
		}
		pre := r.KeepPrefix
		suf := r.KeepSuffix
		if pre+suf >= len(s) {
			return s
		}
		return s[:pre] + ph + s[len(s)-suf:]
	default: // StrategyFull
		ph := r.Placeholder
		if ph == "" {
			ph = "***"
		}
		return strings.Repeat(string(ph[0]), len(ph))
	}
}
