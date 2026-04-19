// Package tail provides functionality to read the last N log entries.
package tail

import "github.com/yourorg/logslice/internal/parser"

// Options configures tail behaviour.
type Options struct {
	N int // number of entries to return; 0 means all
}

// Run returns the last N entries from the provided slice.
// If N <= 0 or N >= len(entries), all entries are returned.
func Run(entries []*parser.Entry, opts Options) []*parser.Entry {
	if len(entries) == 0 {
		return entries
	}
	if opts.N <= 0 || opts.N >= len(entries) {
		return entries
	}
	return entries[len(entries)-opts.N:]
}

// Head returns the first N entries from the provided slice.
// If N <= 0 or N >= len(entries), all entries are returned.
func Head(entries []*parser.Entry, opts Options) []*parser.Entry {
	if len(entries) == 0 {
		return entries
	}
	if opts.N <= 0 || opts.N >= len(entries) {
		return entries
	}
	return entries[:opts.N]
}
