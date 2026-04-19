package stats

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrint_NonEmpty(t *testing.T) {
	entries := makeEntries()
	s := Compute(entries)
	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "Total entries : 3") {
		t.Errorf("missing total: %s", out)
	}
	if !strings.Contains(out, "Earliest") {
		t.Errorf("missing earliest: %s", out)
	}
	if !strings.Contains(out, "Latest") {
		t.Errorf("missing latest: %s", out)
	}
	if !strings.Contains(out, "info=2") {
		t.Errorf("missing level counts: %s", out)
	}
	if !strings.Contains(out, "warn=1") {
		t.Errorf("missing warn count: %s", out)
	}
}

func TestPrint_Empty(t *testing.T) {
	s := Compute(nil)
	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "Total entries : 0") {
		t.Errorf("expected zero total: %s", out)
	}
	if strings.Contains(out, "Earliest") {
		t.Errorf("should not print earliest for empty: %s", out)
	}
}

func TestPrint_FieldKeys(t *testing.T) {
	entries := makeEntries()
	s := Compute(entries)
	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "Fields") {
		t.Errorf("expected Fields line: %s", out)
	}
	if !strings.Contains(out, "level") {
		t.Errorf("expected 'level' in fields: %s", out)
	}
}
