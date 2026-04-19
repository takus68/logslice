// Package tail provides utilities for selecting the last or first N log entries
// from a slice of parsed log entries.
//
// Run returns the last N entries (tail behaviour), while Head returns the first N
// entries. Both functions are safe to call with N larger than the length of the
// input slice — they simply return all available entries.
//
// ParseConfig parses a string map (typically derived from CLI flags) into a
// Config struct that controls which mode and how many lines are selected.
package tail
