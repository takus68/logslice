// Package coalesce provides log entry field coalescing.
//
// For each configured rule, Run scans a list of source fields in order and
// writes the first non-nil, non-empty value into the destination field of a
// new entry copy. This is useful for normalising logs that use different field
// names for the same concept across services (e.g. "msg", "message",
// "log_message").
//
// Rules are expressed as "dest=src1,src2,..." strings and parsed with
// ParseRules.
package coalesce
