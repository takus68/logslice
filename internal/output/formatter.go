package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yourorg/logslice/internal/parser"
)

// Format controls how log entries are rendered.
type Format string

const (
	FormatJSON    Format = "json"
	FormatPretty  Format = "pretty"
	FormatCompact Format = "compact"
)

// Write writes log entries to w in the requested format.
func Write(w io.Writer, entries []parser.Entry, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, entries)
	case FormatPretty:
		return writePretty(w, entries)
	case FormatCompact:
		return writeCompact(w, entries)
	default:
		return fmt.Errorf("unknown format: %q", format)
	}
}

func writeJSON(w io.Writer, entries []parser.Entry) error {
	enc := json.NewEncoder(w)
	for _, e := range entries {
		if err := enc.Encode(e.Raw); err != nil {
			return err
		}
	}
	return nil
}

func writePretty(w io.Writer, entries []parser.Entry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, e := range entries {
		line, err := json.MarshalIndent(e.Raw, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintf(tw, "%s\t%s\n", e.Timestamp.Format("2006-01-02T15:04:05Z07:00"), line)
	}
	return tw.Flush()
}

func writeCompact(w io.Writer, entries []parser.Entry) error {
	for _, e := range entries {
		line, err := json.Marshal(e.Raw)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%s %s\n", e.Timestamp.Format("15:04:05"), line)
	}
	return nil
}
