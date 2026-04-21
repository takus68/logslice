// Package join provides functionality to join two streams of log entries
// on a common field, similar to a SQL JOIN operation.
package join

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Config holds configuration for a join operation.
type Config struct {
	// JoinField is the field name used to match entries across streams.
	JoinField string

	// Mode is the join type: "inner", "left", or "outer".
	Mode string

	// Prefix is prepended to field names from the right stream to avoid collisions.
	Prefix string
}

// ParseConfig parses join options from a slice of key=value strings.
// Supported options:
//
//	on=<field>       — field to join on (required)
//	mode=<type>      — join mode: inner (default), left, outer
//	prefix=<string>  — prefix for right-stream fields (default: "r_")
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{
		Mode:   "inner",
		Prefix: "r_",
	}

	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return Config{}, fmt.Errorf("invalid option %q: expected key=value", opt)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "on":
			if val == "" {
				return Config{}, fmt.Errorf("join field 'on' must not be empty")
			}
			cfg.JoinField = val
		case "mode":
			switch val {
			case "inner", "left", "outer":
				cfg.Mode = val
			default:
				return Config{}, fmt.Errorf("unknown join mode %q: must be inner, left, or outer", val)
			}
		case "prefix":
			cfg.Prefix = val
		default:
			return Config{}, fmt.Errorf("unknown option %q", key)
		}
	}

	if cfg.JoinField == "" {
		return Config{}, fmt.Errorf("option 'on' is required")
	}

	return cfg, nil
}

// Run joins two slices of log entries on a common field.
//
// For each entry in left, it looks for matching entries in right where
// the join field value is equal. Matched right-side fields are merged
// into the result entry with the configured prefix applied.
//
// Modes:
//   - inner: only entries with a match in both streams are returned.
//   - left:  all left entries are returned; unmatched entries have no right fields.
//   - outer: all left entries plus unmatched right entries are returned.
func Run(left, right []parser.Entry, cfg Config) []parser.Entry {
	// Index right entries by join field value.
	rightIndex := make(map[string][]parser.Entry)
	for _, e := range right {
		val, ok := e.Fields[cfg.JoinField]
		if !ok {
			continue
		}
		key := fmt.Sprintf("%v", val)
		rightIndex[key] = append(rightIndex[key], e)
	}

	matched := make(map[string]bool)
	var results []parser.Entry

	for _, le := range left {
		joinVal, ok := le.Fields[cfg.JoinField]
		if !ok {
			if cfg.Mode == "left" || cfg.Mode == "outer" {
				results = append(results, le)
			}
			continue
		}
		key := fmt.Sprintf("%v", joinVal)
		matches, found := rightIndex[key]

		if !found {
			if cfg.Mode == "left" || cfg.Mode == "outer" {
				results = append(results, le)
			}
			continue
		}

		matched[key] = true
		for _, re := range matches {
			merged := mergeEntries(le, re, cfg.Prefix, cfg.JoinField)
			results = append(results, merged)
		}
	}

	// For outer join, append unmatched right entries.
	if cfg.Mode == "outer" {
		for _, re := range right {
			val, ok := re.Fields[cfg.JoinField]
			if !ok {
				results = append(results, re)
				continue
			}
			key := fmt.Sprintf("%v", val)
			if !matched[key] {
				results = append(results, re)
			}
		}
	}

	return results
}

// mergeEntries combines a left and right entry into a single entry.
// Fields from the right entry are prefixed to avoid collisions.
// The join field is not duplicated from the right side.
func mergeEntries(left, right parser.Entry, prefix, joinField string) parser.Entry {
	fields := make(map[string]interface{}, len(left.Fields)+len(right.Fields))

	for k, v := range left.Fields {
		fields[k] = v
	}
	for k, v := range right.Fields {
		if k == joinField {
			continue
		}
		fields[prefix+k] = v
	}

	return parser.Entry{
		Timestamp: left.Timestamp,
		Raw:       left.Raw,
		Fields:    fields,
	}
}
