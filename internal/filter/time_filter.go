package filter

import (
	"time"

	"github.com/user/logslice/internal/parser"
)

// TimeRange holds optional start and end bounds for filtering log entries.
type TimeRange struct {
	From *time.Time
	To   *time.Time
}

// ByTimeRange returns only the log entries whose timestamps fall within the
// given range. A nil bound means unbounded in that direction.
// Entries with a zero timestamp are excluded when any bound is set.
func ByTimeRange(entries []parser.LogEntry, tr TimeRange) []parser.LogEntry {
	if tr.From == nil && tr.To == nil {
		return entries
	}

	var result []parser.LogEntry
	for _, e := range entries {
		if e.Timestamp.IsZero() {
			continue
		}
		if tr.From != nil && e.Timestamp.Before(*tr.From) {
			continue
		}
		if tr.To != nil && e.Timestamp.After(*tr.To) {
			continue
		}
		result = append(result, e)
	}
	return result
}
