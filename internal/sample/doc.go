// Package sample implements log entry sampling for logslice.
//
// Two strategies are supported:
//
//   - nth: keeps every Nth log entry (deterministic, order-preserving)
//   - random: keeps each entry with a given probability using a seeded RNG
//
// Use Run with an Options struct to apply sampling to a slice of parser.Entry.
package sample
