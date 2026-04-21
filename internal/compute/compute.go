// Package compute provides field-level arithmetic operations on log entries.
// It supports adding, subtracting, multiplying, and dividing numeric fields,
// writing the result into a destination field.
package compute

import (
	"fmt"
	"math"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule describes a single compute operation.
type Rule struct {
	Dest string
	Left string
	Op   string // +, -, *, /
	Right string
}

// Run applies each rule to every entry, writing the result into Dest.
// Entries where either operand is missing or non-numeric are left unchanged.
func Run(entries []*parser.Entry, rules []Rule) []*parser.Entry {
	out := make([]*parser.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, applyRules(e, rules))
	}
	return out
}

func applyRules(e *parser.Entry, rules []Rule) *parser.Entry {
	copy := copyEntry(e)
	for _, r := range rules {
		lv, ok1 := toFloat(copy.Fields[r.Left])
		rv, ok2 := toFloat(copy.Fields[r.Right])
		if !ok1 || !ok2 {
			continue
		}
		result, err := compute(lv, rv, r.Op)
		if err != nil {
			continue
		}
		// Store as float64; callers can cast downstream if needed.
		copy.Fields[r.Dest] = result
	}
	return copy
}

func compute(l, r float64, op string) (float64, error) {
	switch op {
	case "+":
		return l + r, nil
	case "-":
		return l - r, nil
	case "*":
		return l * r, nil
	case "/":
		if r == 0 {
			return math.NaN(), fmt.Errorf("division by zero")
		}
		return l / r, nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", op)
	}
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	}
	return 0, false
}

func copyEntry(e *parser.Entry) *parser.Entry {
	fields := make(map[string]interface{}, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	return &parser.Entry{Timestamp: e.Timestamp, Fields: fields}
}
