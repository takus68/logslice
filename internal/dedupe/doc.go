// Package dedupe provides deduplication of parsed log entries.
//
// Entries can be deduplicated by their full content (all fields) or by a
// specific subset of fields. The first occurrence of each unique entry is
// retained; all subsequent duplicates are dropped.
//
// Example usage:
//
//	result := dedupe.Run(entries, dedupe.Options{
//		Strategy: dedupe.ByFields,
//		Fields:   []string{"level", "msg"},
//	})
package dedupe
