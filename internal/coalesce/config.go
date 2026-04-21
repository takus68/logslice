package coalesce

import (
	"fmt"
	"strings"
)

// Rule describes a coalesce operation: write the first non-empty value from
// Sources into Dest.
type Rule struct {
	Dest    string
	Sources []string
}

// ParseRules parses specs of the form "dest=src1,src2,...".
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("coalesce: invalid rule %q: missing '='" , spec)
		}
		dest := strings.TrimSpace(parts[0])
		if dest == "" {
			return nil, fmt.Errorf("coalesce: invalid rule %q: empty destination key", spec)
		}
		rawSources := strings.Split(parts[1], ",")
		var sources []string
		for _, s := range rawSources {
			s = strings.TrimSpace(s)
			if s != "" {
				sources = append(sources, s)
			}
		}
		if len(sources) == 0 {
			return nil, fmt.Errorf("coalesce: invalid rule %q: no source fields", spec)
		}
		rules = append(rules, Rule{Dest: dest, Sources: sources})
	}
	return rules, nil
}
