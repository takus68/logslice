// Package filter provides utilities for filtering parsed log entries
// by time range or field value.
package filter

import (
	"time"

	"github.com/user/logslice/internal/parser"
)

// Options holds all filtering criteria that can be applied together.
type Options struct {
	From    time.Time
	To      time.Time
	Matcher *FieldMatcher // optional field filter
}

// Apply runs all active filters in Options against the provided entries
// and returns the surviving subset. Time bounds are applied first,
// followed by the optional field matcher.
func Apply(entries []parser.LogEntry, opts Options) []parser.LogEntry {
	result := entries

	zero := time.Time{}
	if opts.From != zero || opts.To != zero {
		result = ByTimeRange(result, opts.From, opts.To)
	}

	if opts.Matcher != nil {
		result = ByField(result, *opts.Matcher)
	}

	return result
}
