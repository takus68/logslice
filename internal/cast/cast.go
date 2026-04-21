// Package cast provides functionality for casting log entry field values
// between Go types without failing the entire pipeline on partial errors.
// Unlike convert, cast is lenient: fields that cannot be cast are left unchanged.
package cast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule describes a single cast operation: a field name and the target type.
type Rule struct {
	Field  string
	Target string // "string", "int", "float", "bool"
}

// ParseRules parses cast rule specs of the form "field=type".
// Example: ["status=int", "active=bool"]
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("cast: invalid rule %q: expected field=type", spec)
		}
		field := strings.TrimSpace(parts[0])
		target := strings.TrimSpace(parts[1])
		if field == "" {
			return nil, fmt.Errorf("cast: empty field name in rule %q", spec)
		}
		if !isValidType(target) {
			return nil, fmt.Errorf("cast: unknown type %q in rule %q (want string|int|float|bool)", target, spec)
		}
		rules = append(rules, Rule{Field: field, Target: target})
	}
	return rules, nil
}

func isValidType(t string) bool {
	switch t {
	case "string", "int", "float", "bool":
		return true
	}
	return false
}

// Run applies the cast rules to each entry. Fields that cannot be cast are
// left at their original value; no error is returned for individual failures.
func Run(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, castEntry(e, rules))
	}
	return out
}

func castEntry(e parser.Entry, rules []Rule) parser.Entry {
	fields := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	for _, r := range rules {
		val, ok := fields[r.Field]
		if !ok {
			continue
		}
		if casted, err := castValue(val, r.Target); err == nil {
			fields[r.Field] = casted
		}
		// on error: leave original value unchanged
	}
	return parser.Entry{Timestamp: e.Timestamp, Fields: fields}
}

func castValue(val any, target string) (any, error) {
	str := fmt.Sprintf("%v", val)
	switch target {
	case "string":
		return str, nil
	case "int":
		// handle float strings like "3.0"
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			return i, nil
		}
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			return int64(f), nil
		}
		return nil, fmt.Errorf("cannot cast %q to int", str)
	case "float":
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot cast %q to float", str)
		}
		return f, nil
	case "bool":
		b, err := strconv.ParseBool(str)
		if err != nil {
			return nil, fmt.Errorf("cannot cast %q to bool", str)
		}
		return b, nil
	}
	return nil, fmt.Errorf("unknown type %q", target)
}
