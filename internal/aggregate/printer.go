package aggregate

import (
	"fmt"
	"io"
)

// Print writes a human-readable aggregation table to w.
func Print(w io.Writer, r *Result) {
	if r == nil {
		fmt.Fprintln(w, "no aggregation result")
		return
	}

	fmt.Fprintf(w, "Aggregation by field: %q\n", r.Field)
	fmt.Fprintf(w, "%-30s %s\n", "Value", "Count")
	fmt.Fprintln(w, "-------------------------------+-------")

	for _, k := range r.SortedKeys() {
		fmt.Fprintf(w, "%-30s %d\n", k, r.Counts[k])
	}

	fmt.Fprintf(w, "\nTotal: %d\n", r.Total)
}
