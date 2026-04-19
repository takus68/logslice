package stats

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Print writes a human-readable summary to w.
func Print(w io.Writer, s Summary) {
	fmt.Fprintf(w, "Total entries : %d\n", s.Total)
	if s.Total == 0 {
		return
	}
	fmt.Fprintf(w, "Earliest      : %s\n", s.Earliest.Format("2006-01-02T15:04:05Z07:00"))
	fmt.Fprintf(w, "Latest        : %s\n", s.Latest.Format("2006-01-02T15:04:05Z07:00"))

	if len(s.LevelCounts) > 0 {
		levels := make([]string, 0, len(s.LevelCounts))
		for l := range s.LevelCounts {
			levels = append(levels, l)
		}
		sort.Strings(levels)
		parts := make([]string, 0, len(levels))
		for _, l := range levels {
			parts = append(parts, fmt.Sprintf("%s=%d", l, s.LevelCounts[l]))
		}
		fmt.Fprintf(w, "Levels        : %s\n", strings.Join(parts, ", "))
	}

	if len(s.FieldKeys) > 0 {
		sort.Strings(s.FieldKeys)
		fmt.Fprintf(w, "Fields        : %s\n", strings.Join(s.FieldKeys, ", "))
	}
}
