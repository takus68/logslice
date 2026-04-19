package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries() []parser.Entry {
	t1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 1, 15, 10, 1, 0, 0, time.UTC)
	return []parser.Entry{
		{
			Timestamp: t1,
			Raw:       map[string]interface{}{"time": t1.Format(time.RFC3339), "level": "info", "msg": "started"},
		},
		{
			Timestamp: t2,
			Raw:       map[string]interface{}{"time": t2.Format(time.RFC3339), "level": "error", "msg": "failed"},
		},
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, makeEntries(), FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "started") || !strings.Contains(out, "failed") {
		t.Errorf("expected both entries in output, got: %s", out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 JSON lines, got %d", len(lines))
	}
}

func TestWrite_Compact(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, makeEntries(), FormatCompact); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "10:00:00") {
		t.Errorf("expected timestamp prefix in compact output, got: %s", out)
	}
}

func TestWrite_Pretty(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, makeEntries(), FormatPretty); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2024-01-15") {
		t.Errorf("expected date in pretty output, got: %s", out)
	}
}

func TestWrite_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, makeEntries(), Format("xml"))
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestWrite_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, []parser.Entry{}, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for no entries, got %d bytes", buf.Len())
	}
}
