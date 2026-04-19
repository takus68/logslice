package filter

import (
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// ByTimeRange returns entries whose timestamps fall within [start, end].
// If end is zero, only a lower bound is applied.
func ByTimeRange(entries []parser.Entry, start, end time.Time) []parser.Entry {
	var result []parser.Entry
	for _, e := range entries {
		if e.Timestamp.Before(start) {
			continue
		}
		if !end.IsZero() && e.Timestamp.After(end) {
			continue
		}
		result = append(result, e)
	}
	if result == nil {
		return []parser.Entry{}
	}
	return result
}
