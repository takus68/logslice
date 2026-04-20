package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseRules parses validation rule specs of the form:
//
//	"field:required"
//	"field:type=string"
//	"field:pattern=^[a-z]+$"
//	"field:required,type=number"
//
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		parts := strings.SplitN(spec, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rule spec %q: expected field:constraints", spec)
		}
		field := strings.TrimSpace(parts[0])
		if field == "" {
			return nil, fmt.Errorf("invalid rule spec %q: empty field name", spec)
		}
		r := Rule{Field: field}
		constraints := strings.Split(parts[1], ",")
		for _, c := range constraints {
			c = strings.TrimSpace(c)
			switch {
			case c == "required":
				r.Required = true
			case strings.HasPrefix(c, "type="):
				t := strings.TrimPrefix(c, "type=")
				if t != "string" && t != "number" && t != "bool" {
					return nil, fmt.Errorf("unknown type %q in spec %q", t, spec)
				}
				r.TypeName = t
			case strings.HasPrefix(c, "pattern="):
				patStr := strings.TrimPrefix(c, "pattern=")
				pat, err := regexp.Compile(patStr)
				if err != nil {
					return nil, fmt.Errorf("invalid pattern %q in spec %q: %w", patStr, spec, err)
				}
				r.Pattern = pat
			default:
				return nil, fmt.Errorf("unknown constraint %q in spec %q", c, spec)
			}
		}
		rules = append(rules, r)
	}
	return rules, nil
}
