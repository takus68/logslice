// Package enrich adds derived or static fields to structured log entries.
//
// Rules are expressed as "key=value" strings where value may include
// {fieldName} placeholders that are substituted from each entry's fields.
//
// Example usage:
//
//	rules, err := enrich.ParseRules([]string{"env=prod", "id={host}-{pid}"})
//	if err != nil { ... }
//	enriched := enrich.Run(entries, rules)
package enrich
