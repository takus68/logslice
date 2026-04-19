// Package redact implements field-level redaction and pattern-based masking
// for structured log entries.
//
// It supports two modes:
//
//   - Full redaction: replaces a field's value with the literal string "[REDACTED]".
//   - Pattern masking: applies a regular expression to a field's string value,
//     replacing each match with an equal-length sequence of '*' characters.
//
// Redaction is non-destructive — original entries are never modified.
package redact
