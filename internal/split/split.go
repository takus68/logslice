// Package split provides functionality to split log entries into chunks
// based on a maximum count or a field value boundary.
package split

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Config holds options for splitting log entries.
type Config struct {
	// ChunkSize splits entries into groups of at most ChunkSize entries.
	ChunkSize int
	// FieldBoundary splits whenever the value of this field changes.
	FieldBoundary string
}

// ParseConfig parses key=value options into a Config.
// Supported keys: size, boundary.
func ParseConfig(opts []string) (Config, error) {
	cfg := Config{ChunkSize: 0}
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return cfg, fmt.Errorf("split: invalid option %q, expected key=value", opt)
		}
		key, val := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "size":
			n, err := strconv.Atoi(val)
			if err != nil || n <= 0 {
				return cfg, fmt.Errorf("split: invalid size %q, must be positive integer", val)
			}
			cfg.ChunkSize = n
		case "boundary":
			cfg.FieldBoundary = val
		default:
			return cfg, fmt.Errorf("split: unknown option %q", key)
		}
	}
	return cfg, nil
}

// Run splits entries into chunks according to cfg.
// If ChunkSize > 0, entries are grouped into slices of that size.
// If FieldBoundary is set, a new chunk begins whenever the field value changes.
// If neither is set, all entries are returned as a single chunk.
func Run(entries []*parser.LogEntry, cfg Config) [][]*parser.LogEntry {
	if len(entries) == 0 {
		return nil
	}
	if cfg.FieldBoundary != "" {
		return splitByField(entries, cfg.FieldBoundary)
	}
	if cfg.ChunkSize > 0 {
		return splitBySize(entries, cfg.ChunkSize)
	}
	return [][]*parser.LogEntry{entries}
}

func splitBySize(entries []*parser.LogEntry, size int) [][]*parser.LogEntry {
	var chunks [][]*parser.LogEntry
	for i := 0; i < len(entries); i += size {
		end := i + size
		if end > len(entries) {
			end = len(entries)
		}
		chunks = append(chunks, entries[i:end])
	}
	return chunks
}

func splitByField(entries []*parser.LogEntry, field string) [][]*parser.LogEntry {
	var chunks [][]*parser.LogEntry
	var current []*parser.LogEntry
	var lastVal interface{}

	for _, e := range entries {
		val := e.Raw[field]
		if current == nil {
			lastVal = val
		} else if val != lastVal {
			chunks = append(chunks, current)
			current = nil
			lastVal = val
		}
		current = append(current, e)
	}
	if len(current) > 0 {
		chunks = append(chunks, current)
	}
	return chunks
}
