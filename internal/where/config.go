package where

import (
	"fmt"
	"strings"
)

var validOps = []string{">=", "<=", "!=", "==", "=", ">", "<", "contains"}

// ParseRules parses a slice of expression strings into Rules.
// Each spec must be of the form: field<op>value
// Supported operators: ==, =, !=, >, >=, <, <=, contains
//
// Example specs:
//
//	"level==error"
//	"status>=400"
//	"message contains timeout"
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		r, err := parseRule(spec)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func parseRule(spec string) (Rule, error) {
	for _, op := range validOps {
		idx := strings.Index(spec, op)
		if idx < 0 {
			continue
		}
		field := strings.TrimSpace(spec[:idx])
		value := strings.TrimSpace(spec[idx+len(op):])
		if field == "" {
			return Rule{}, fmt.Errorf("where: empty field in expression %q", spec)
		}
		return Rule{Field: field, Op: op, Value: value}, nil
	}
	return Rule{}, fmt.Errorf("where: no valid operator found in expression %q", spec)
}
