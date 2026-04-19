// Package aggregate provides grouping and counting of log entries by field value.
package aggregate

import (
	"fmt"
	"sort"

	"github.com/user/logslice/internal/parser"
)

// Result holds aggregation results for a single field.
type Result struct {
	Field  string
	Counts map[string]int
	Total  int
}

// ByField groups log entries by the value of the given field and counts occurrences.
func ByField(entries []*parser.Entry, field string) (*Result, error) {
	if field == "" {
		return nil, fmt.Errorf("aggregate: field name must not be empty")
	}

	counts := make(map[string]int)
	for _, e := range entries {
		val, ok := e.Fields[field]
		if !ok {
			counts["<missing>"]++
			continue
		}
		counts[fmt.Sprintf("%v", val)]++
	}

	total := 0
	for _, c := range counts {
		total += c
	}

	return &Result{
		Field:  field,
		Counts: counts,
		Total:  total,
	}, nil
}

// SortedKeys returns the keys of the result sorted alphabetically.
func (r *Result) SortedKeys() []string {
	keys := make([]string, 0, len(r.Counts))
	for k := range r.Counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
