package compute

import (
	"fmt"
	"strings"
)

var validOps = map[string]bool{"+": true, "-": true, "*": true, "/": true}

// ParseRules parses specs of the form "dest=left+right" or "dest=left/right", etc.
// Each spec must contain exactly one operator character after the '=' separator.
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		eq := strings.Index(spec, "=")
		if eq < 1 {
			return nil, fmt.Errorf("compute: missing '=' in spec %q", spec)
		}
		dest := strings.TrimSpace(spec[:eq])
		expr := strings.TrimSpace(spec[eq+1:])
		if dest == "" {
			return nil, fmt.Errorf("compute: empty destination in spec %q", spec)
		}
		op, left, right, err := parseExpr(expr, spec)
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{Dest: dest, Left: left, Op: op, Right: right})
	}
	return rules, nil
}

func parseExpr(expr, spec string) (op, left, right string, err error) {
	// Scan for an operator character, respecting that field names won't contain them.
	for i := 1; i < len(expr); i++ {
		ch := string(expr[i])
		if validOps[ch] {
			left = strings.TrimSpace(expr[:i])
			right = strings.TrimSpace(expr[i+1:])
			if left == "" || right == "" {
				return "", "", "", fmt.Errorf("compute: empty operand in spec %q", spec)
			}
			return ch, left, right, nil
		}
	}
	return "", "", "", fmt.Errorf("compute: no valid operator (+,-,*,/) found in spec %q", spec)
}
