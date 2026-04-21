package mask

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseRules parses masking rule specs of the form:
//
//	"field:full[:placeholder]"
//	"field:partial:keepPrefix:keepSuffix[:placeholder]"
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
	parts := strings.Split(spec, ":")
	if len(parts) < 2 {
		return Rule{}, fmt.Errorf("mask: invalid rule %q: expected field:strategy", spec)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return Rule{}, fmt.Errorf("mask: empty field name in rule %q", spec)
	}
	strategy := Strategy(strings.TrimSpace(parts[1]))

	switch strategy {
	case StrategyFull:
		r := Rule{Field: field, Strategy: StrategyFull}
		if len(parts) >= 3 {
			r.Placeholder = strings.TrimSpace(parts[2])
		}
		return r, nil

	case StrategyPartial:
		if len(parts) < 4 {
			return Rule{}, fmt.Errorf("mask: partial rule %q requires keepPrefix and keepSuffix", spec)
		}
		pre, err := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err != nil || pre < 0 {
			return Rule{}, fmt.Errorf("mask: invalid keepPrefix in rule %q", spec)
		}
		suf, err := strconv.Atoi(strings.TrimSpace(parts[3]))
		if err != nil || suf < 0 {
			return Rule{}, fmt.Errorf("mask: invalid keepSuffix in rule %q", spec)
		}
		r := Rule{Field: field, Strategy: StrategyPartial, KeepPrefix: pre, KeepSuffix: suf}
		if len(parts) >= 5 {
			r.Placeholder = strings.TrimSpace(parts[4])
		}
		return r, nil

	default:
		return Rule{}, fmt.Errorf("mask: unknown strategy %q in rule %q", strategy, spec)
	}
}
