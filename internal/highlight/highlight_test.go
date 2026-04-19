package highlight

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
		Raw:       `{"level":"error","msg":"boom"}`,
	}
}

func TestApply_NoMatch(t *testing.T) {
	e := makeEntry(map[string]any{"level": "info"})
	rules := []Rule{{Field: "level", Value: "error", Color: Red}}
	out := Apply(e, rules)
	if out != e.Raw {
		t.Errorf("expected no color, got %q", out)
	}
}

func TestApply_ExactMatch(t *testing.T) {
	e := makeEntry(map[string]any{"level": "error"})
	rules := []Rule{{Field: "level", Value: "error", Color: Red}}
	out := Apply(e, rules)
	expected := Red + e.Raw + Reset
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	e := makeEntry(map[string]any{"level": "ERROR"})
	rules := []Rule{{Field: "level", Value: "error", Color: Yellow}}
	out := Apply(e, rules)
	if out != Yellow+e.Raw+Reset {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestApply_SubstringMatch(t *testing.T) {
	e := makeEntry(map[string]any{"msg": "connection refused"})
	rules := []Rule{{Field: "msg", Value: "refused", Color: Red, Substring: true}}
	out := Apply(e, rules)
	if out != Red+e.Raw+Reset {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestApply_MissingField(t *testing.T) {
	e := makeEntry(map[string]any{"level": "info"})
	rules := []Rule{{Field: "service", Value: "api", Color: Blue}}
	out := Apply(e, rules)
	if out != e.Raw {
		t.Errorf("expected raw, got %q", out)
	}
}

func TestApplyAll(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(map[string]any{"level": "error"}),
		makeEntry(map[string]any{"level": "info"}),
	}
	rules := []Rule{{Field: "level", Value: "error", Color: Red}}
	out := ApplyAll(entries, rules)
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	if out[1] != entries[1].Raw {
		t.Errorf("second entry should be uncolored")
	}
}
