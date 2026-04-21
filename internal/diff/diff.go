// Package diff compares two streams of log entries and reports added, removed,
// or changed entries based on a key field.
package diff

import (
	"fmt"

	"github.com/yourorg/logslice/internal/parser"
)

// Mode controls what the diff output includes.
type Mode string

const (
	ModeAll     Mode = "all"
	ModeAdded   Mode = "added"
	ModeRemoved Mode = "removed"
	ModeChanged Mode = "changed"
)

// Result holds a single diff entry with a tag indicating its status.
type Result struct {
	Tag   string
	Entry parser.Entry
}

// Run compares entries from left and right slices keyed by keyField.
// It returns a slice of Results tagged with "+", "-", or "~".
func Run(left, right []parser.Entry, keyField string, mode Mode) ([]Result, error) {
	if keyField == "" {
		return nil, fmt.Errorf("diff: keyField must not be empty")
	}

	leftMap := indexByField(left, keyField)
	rightMap := indexByField(right, keyField)

	var results []Result

	for key, le := range leftMap {
		re, exists := rightMap[key]
		if !exists {
			if mode == ModeAll || mode == ModeRemoved {
				results = append(results, Result{Tag: "-", Entry: le})
			}
			continue
		}
		if !equalRaw(le.Raw, re.Raw) {
			if mode == ModeAll || mode == ModeChanged {
				results = append(results, Result{Tag: "~", Entry: re})
			}
		}
	}

	for key, re := range rightMap {
		if _, exists := leftMap[key]; !exists {
			if mode == ModeAll || mode == ModeAdded {
				results = append(results, Result{Tag: "+", Entry: re})
			}
		}
	}

	return results, nil
}

func indexByField(entries []parser.Entry, field string) map[string]parser.Entry {
	m := make(map[string]parser.Entry, len(entries))
	for _, e := range entries {
		if v, ok := e.Fields[field]; ok {
			m[fmt.Sprintf("%v", v)] = e
		}
	}
	return m
}

func equalRaw(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, av := range a {
		bv, ok := b[k]
		if !ok {
			return false
		}
		if fmt.Sprintf("%v", av) != fmt.Sprintf("%v", bv) {
			return false
		}
	}
	return true
}
