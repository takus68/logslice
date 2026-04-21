package compute

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(fields ...map[string]interface{}) []*parser.Entry {
	var out []*parser.Entry
	for _, f := range fields {
		out = append(out, &parser.Entry{Timestamp: time.Now(), Fields: f})
	}
	return out
}

func TestRun_Addition(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"a": float64(3), "b": float64(4)})
	rules := []Rule{{Dest: "sum", Left: "a", Op: "+", Right: "b"}}
	out := Run(entries, rules)
	if v, ok := out[0].Fields["sum"]; !ok || v.(float64) != 7 {
		t.Fatalf("expected sum=7, got %v", v)
	}
}

func TestRun_Subtraction(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"x": float64(10), "y": float64(3)})
	rules := []Rule{{Dest: "diff", Left: "x", Op: "-", Right: "y"}}
	out := Run(entries, rules)
	if v := out[0].Fields["diff"].(float64); v != 7 {
		t.Fatalf("expected 7, got %v", v)
	}
}

func TestRun_Division_ByZero(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"a": float64(5), "b": float64(0)})
	rules := []Rule{{Dest: "res", Left: "a", Op: "/", Right: "b"}}
	out := Run(entries, rules)
	if _, ok := out[0].Fields["res"]; ok {
		t.Fatal("expected res field to be absent on divide-by-zero")
	}
}

func TestRun_MissingOperand(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"a": float64(5)})
	rules := []Rule{{Dest: "res", Left: "a", Op: "+", Right: "missing"}}
	out := Run(entries, rules)
	if _, ok := out[0].Fields["res"]; ok {
		t.Fatal("expected res to be absent when operand missing")
	}
}

func TestRun_NonNumericField(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"a": "hello", "b": float64(2)})
	rules := []Rule{{Dest: "res", Left: "a", Op: "*", Right: "b"}}
	out := Run(entries, rules)
	if _, ok := out[0].Fields["res"]; ok {
		t.Fatal("expected res to be absent for non-numeric field")
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	original := makeEntries(map[string]interface{}{"a": float64(1), "b": float64(2)})
	rules := []Rule{{Dest: "c", Left: "a", Op: "+", Right: "b"}}
	Run(original, rules)
	if _, ok := original[0].Fields["c"]; ok {
		t.Fatal("original entry should not be mutated")
	}
}

func TestRun_MultipleRules(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"a": float64(6), "b": float64(2)})
	rules := []Rule{
		{Dest: "sum", Left: "a", Op: "+", Right: "b"},
		{Dest: "product", Left: "a", Op: "*", Right: "b"},
	}
	out := Run(entries, rules)
	if out[0].Fields["sum"].(float64) != 8 {
		t.Fatalf("expected sum=8")
	}
	if out[0].Fields["product"].(float64) != 12 {
		t.Fatalf("expected product=12")
	}
}
