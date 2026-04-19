package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// LogEntry represents a single parsed log line.
type LogEntry struct {
	Timestamp time.Time
	Fields    map[string]interface{}
	Raw       string
}

// TimeField is the JSON key used to extract the timestamp.
const TimeField = "time"

// Parse reads newline-delimited JSON log entries from r and returns them.
func Parse(r io.Reader) ([]LogEntry, error) {
	var entries []LogEntry
	scanner := bufio.NewScanner(r)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}
		entry, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func parseLine(line string) (LogEntry, error) {
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return LogEntry{}, fmt.Errorf("invalid JSON: %w", err)
	}

	entry := LogEntry{Fields: fields, Raw: line}

	if v, ok := fields[TimeField]; ok {
		switch val := v.(type) {
		case string:
			t, err := time.Parse(time.RFC3339Nano, val)
			if err != nil {
				t, err = time.Parse(time.RFC3339, val)
				if err != nil {
					return LogEntry{}, fmt.Errorf("cannot parse timestamp %q", val)
				}
			}
			entry.Timestamp = t
		default:
			return LogEntry{}, fmt.Errorf("timestamp field %q is not a string", TimeField)
		}
	}

	return entry, nil
}
