package stats

import (
	"time"

	"github.com/user/logslice/internal/parser"
)

// Summary holds aggregate statistics for a set of log entries.
type Summary struct {
	Total      int
	Earliest   time.Time
	Latest     time.Time
	LevelCounts map[string]int
	FieldKeys  []string
}

// Compute calculates statistics over the provided log entries.
func Compute(entries []parser.Entry) Summary {
	if len(entries) == 0 {
		return Summary{LevelCounts: map[string]int{}}
	}

	s := Summary{
		Total:       len(entries),
		Earliest:    entries[0].Timestamp,
		Latest:      entries[0].Timestamp,
		LevelCounts: map[string]int{},
	}

	keySet := map[string]struct{}{}

	for _, e := range entries {
		if e.Timestamp.Before(s.Earliest) {
			s.Earliest = e.Timestamp
		}
		if e.Timestamp.After(s.Latest) {
			s.Latest = e.Timestamp
		}
		if lvl, ok := e.Fields["level"]; ok {
			s.LevelCounts[lvl]++
		}
		for k := range e.Fields {
			keySet[k] = struct{}{}
		}
	}

	for k := range keySet {
		s.FieldKeys = append(s.FieldKeys, k)
	}

	return s
}
