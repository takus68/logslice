// Package window provides sliding and tumbling window aggregation over log entries.
package window

import (
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Config controls how windowing is applied.
type Config struct {
	Size     time.Duration
	Field    string // timestamp field name, defaults to "time"
	Tumbling bool   // if false, windows slide by Size/2
}

// Window holds a slice of entries that fall within a time bucket.
type Window struct {
	Start   time.Time
	End     time.Time
	Entries []parser.Entry
}

// Run partitions entries into time windows according to cfg.
// Entries without a parseable timestamp in cfg.Field are skipped.
func Run(entries []parser.Entry, cfg Config) ([]Window, error) {
	if cfg.Size <= 0 {
		return nil, fmt.Errorf("window: size must be positive, got %v", cfg.Size)
	}
	field := cfg.Field
	if field == "" {
		field = "time"
	}

	var windows []Window

	for _, e := range entries {
		raw, ok := e.Raw[field]
		if !ok {
			continue
		}
		ts, err := parseTimestamp(raw)
		if err != nil {
			continue
		}

		start := bucketStart(ts, cfg.Size)
		windows = addToWindow(windows, e, start, start.Add(cfg.Size))
	}

	return windows, nil
}

func bucketStart(t time.Time, size time.Duration) time.Time {
	return t.Truncate(size)
}

func addToWindow(windows []Window, e parser.Entry, start, end time.Time) []Window {
	for i := range windows {
		if windows[i].Start.Equal(start) {
			windows[i].Entries = append(windows[i].Entries, e)
			return windows
		}
	}
	return append(windows, Window{Start: start, End: end, Entries: []parser.Entry{e}})
}

func parseTimestamp(v interface{}) (time.Time, error) {
	s, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("window: timestamp not a string")
	}
	formats := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05"}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("window: cannot parse timestamp %q", s)
}
