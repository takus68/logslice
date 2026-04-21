package rename

import (
	"fmt"
	"strings"
)

// ParseRules parses a slice of rule specs in the form "old=new".
// Each spec must contain exactly one "=" separator with non-empty
// keys on both sides. Whitespace around keys is trimmed.
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("rename: invalid rule %q: must be in form old=new", spec)
		}
		from := strings.TrimSpace(parts[0])
		to := strings.TrimSpace(parts[1])
		if from == "" {
			return nil, fmt.Errorf("rename: invalid rule %q: source key must not be empty", spec)
		}
		if to == "" {
			return nil, fmt.Errorf("rename: invalid rule %q: destination key must not be empty", spec)
		}
		rules = append(rules, Rule{From: from, To: to})
	}
	return rules, nil
}
