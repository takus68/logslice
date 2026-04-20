package convert

import (
	"fmt"
	"strings"
)

// ParseRules parses conversion rule specs of the form "field=type".
// Example: ["duration=int", "ratio=float", "active=bool"]
func ParseRules(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rule %q: expected field=type", spec)
		}
		field := strings.TrimSpace(parts[0])
		typeName := TargetType(strings.TrimSpace(parts[1]))
		if field == "" {
			return nil, fmt.Errorf("invalid rule %q: field name is empty", spec)
		}
		if !isValidType(typeName) {
			return nil, fmt.Errorf("invalid rule %q: unknown type %q", spec, typeName)
		}
		rules = append(rules, Rule{Field: field, Type: typeName})
	}
	return rules, nil
}

func isValidType(t TargetType) bool {
	switch t {
	case TypeInt, TypeFloat, TypeBool, TypeString:
		return true
	}
	return false
}
