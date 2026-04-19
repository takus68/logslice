package dedupe

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/user/logslice/internal/parser"
)

// Strategy defines how deduplication is performed.
type Strategy int

const (
	// ByFullEntry deduplicates based on the entire log entry.
	ByFullEntry Strategy = iota
	// ByFields deduplicates based on a subset of fields.
	ByFields
)

// Options configures deduplication behaviour.
type Options struct {
	Strategy Strategy
	// Fields used when Strategy is ByFields.
	Fields []string
}

// Run removes duplicate log entries according to opts.
// The first occurrence of each entry is kept; subsequent duplicates are dropped.
func Run(entries []parser.LogEntry, opts Options) []parser.LogEntry {
	seen := make(map[string]struct{})
	result := make([]parser.LogEntry, 0, len(entries))

	for _, e := range entries {
		key := computeKey(e, opts)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, e)
	}

	return result
}

func computeKey(e parser.LogEntry, opts Options) string {
	var payload any

	if opts.Strategy == ByFields && len(opts.Fields) > 0 {
		subset := make(map[string]any, len(opts.Fields))
		for _, f := range opts.Fields {
			if v, ok := e.Fields[f]; ok {
				subset[f] = v
			}
		}
		payload = subset
	} else {
		payload = e.Fields
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
