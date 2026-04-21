// Package rename implements field-key renaming for structured log entries.
//
// Rules are expressed as "old=new" strings and can be supplied via CLI flags.
// Each matching key in an entry is renamed while preserving its value.
// Keys that are absent in a given entry are silently skipped.
//
// Example usage:
//
//	rules, err := rename.ParseRules([]string{"msg=message", "lvl=level"})
//	if err != nil { ... }
//	result := rename.Run(entries, rules)
package rename
