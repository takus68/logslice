// Package group provides functionality to group log entries by a field value.
package group

import (
	"fmt"
	"sort"

	"github.com/user/logslice/internal/parser"
)

// Result holds grouped log entries keyed by field value.
type Result struct {
	Keys   []string
	Groups map[string][]parser.Entry
}

// Config controls grouping behaviour.
type Config struct {
	// Field is the entry field to group by.
	Field string
	// Sorted determines whether the result keys are sorted alphabetically.
	Sorted bool
}

// Run groups entries by the value of the configured field.
// Entries that are missing the field are placed under the key "<missing>".
func Run(entries []parser.Entry, cfg Config) Result {
	groups := make(map[string][]parser.Entry)

	for _, e := range entries {
		key := "<missing>"
		if v, ok := e.Fields[cfg.Field]; ok {
			key = fmt.Sprintf("%v", v)
		}
		groups[key] = append(groups[key], e)
	}

	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}

	if cfg.Sorted {
		sort.Strings(keys)
	}

	return Result{
		Keys:   keys,
		Groups: groups,
	}
}
