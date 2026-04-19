package aggregate

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_NonEmpty(t *testing.T) {
	entries := makeEntries()
	res, _ := ByField(entries, "level")

	var buf bytes.Buffer
	Print(&buf, res)
	out := buf.String()

	if !strings.Contains(out, "level") {
		t.Error("expected field name in output")
	}
	if !strings.Contains(out, "info") {
		t.Error("expected 'info' in output")
	}
	if !strings.Contains(out, "Total") {
		t.Error("expected Total in output")
	}
}

func TestPrint_Nil(t *testing.T) {
	var buf bytes.Buffer
	Print(&buf, nil)
	if !strings.Contains(buf.String(), "no aggregation") {
		t.Error("expected fallback message for nil result")
	}
}

func TestPrint_SortedOutput(t *testing.T) {
	entries := makeEntries()
	res, _ := ByField(entries, "level")

	var buf bytes.Buffer
	Print(&buf, res)
	lines := strings.Split(buf.String(), "\n")

	var valueLines []string
	for _, l := range lines {
		if strings.HasPrefix(l, "error") || strings.HasPrefix(l, "info") || strings.HasPrefix(l, "warn") {
			valueLines = append(valueLines, l)
		}
	}

	if len(valueLines) < 3 {
		t.Errorf("expected at least 3 value lines, got %d", len(valueLines))
	}

	for i := 1; i < len(valueLines); i++ {
		if valueLines[i] < valueLines[i-1] {
			t.Errorf("output not sorted: %v", valueLines)
		}
	}
}
