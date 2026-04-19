// Package merge provides functionality to merge multiple sorted log streams
// into a single time-ordered sequence.
package merge

import (
	"sort"

	"github.com/user/logslice/internal/parser"
)

// Options controls merge behaviour.
type Options struct {
	// Stable preserves original order for entries with identical timestamps.
	Stable bool
}

// Run merges multiple slices of log entries into one slice sorted by timestamp.
// Entries with equal timestamps are ordered by their source index when Stable
// is true, otherwise their relative order is undefined.
func Run(streams [][]parser.Entry, opts Options) []parser.Entry {
	var total int
	for _, s := range streams {
		total += len(s)
	}
	if total == 0 {
		return nil
	}

	type indexed struct {
		entry  parser.Entry
		stream int
		pos    int
	}

	all := make([]indexed, 0, total)
	for si, s := range streams {
		for pi, e := range s {
			all = append(all, indexed{entry: e, stream: si, pos: pi})
		}
	}

	sort.SliceStable(all, func(i, j int) bool {
		if all[i].entry.Timestamp.Equal(all[j].entry.Timestamp) {
			if opts.Stable {
				if all[i].stream != all[j].stream {
					return all[i].stream < all[j].stream
				}
				return all[i].pos < all[j].pos
			}
			return false
		}
		return all[i].entry.Timestamp.Before(all[j].entry.Timestamp)
	})

	result := make([]parser.Entry, len(all))
	for i, a := range all {
		result[i] = a.entry
	}
	return result
}
