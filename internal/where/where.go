// Package where provides conditional filtering of log entries
// based on simple field comparison expressions.
package where

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Rule represents a single where-clause condition.
type Rule struct {
	Field string
	Op    string
	Value string
}

// Run filters entries that satisfy all provided rules.
func Run(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 {
		return entries
	}
	var out []parser.Entry
	for _, e := range entries {
		if matchesAll(e, rules) {
			out = append(out, e)
		}
	}
	return out
}

func matchesAll(e parser.Entry, rules []Rule) bool {
	for _, r := range rules {
		if !matchesRule(e, r) {
			return false
		}
	}
	return true
}

func matchesRule(e parser.Entry, r Rule) bool {
	raw, ok := e.Fields[r.Field]
	if !ok {
		return false
	}
	actual := fmt.Sprintf("%v", raw)
	switch r.Op {
	case "==", "=":
		return actual == r.Value
	case "!=":
		return actual != r.Value
	case "contains":
		return strings.Contains(actual, r.Value)
	case ">", ">=", "<", "<=":
		return compareNumeric(actual, r.Op, r.Value)
	}
	return false
}

func compareNumeric(actual, op, expected string) bool {
	a, err1 := strconv.ParseFloat(actual, 64)
	b, err2 := strconv.ParseFloat(expected, 64)
	if err1 != nil || err2 != nil {
		return false
	}
	switch op {
	case ">":
		return a > b
	case ">=":
		return a >= b
	case "<":
		return a < b
	case "<=":
		return a <= b
	}
	return false
}
