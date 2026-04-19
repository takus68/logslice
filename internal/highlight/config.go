package highlight

import (
	"fmt"
	"strings"
)

var namedColors = map[string]string{
	"red":    Red,
	"green":  Green,
	"yellow": Yellow,
	"blue":   Blue,
	"cyan":   Cyan,
}

// ParseRules parses highlight rules from strings of the form:
//   field=value:color
//   field~value:color   (substring match)
func ParseRules(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		r, err := parseRule(spec)
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func parseRule(spec string) (Rule, error) {
	// format: field=value:color or field~value:color
	substring := false
	sep := "="
	if strings.Contains(spec, "~") && (!strings.Contains(spec, "=") || strings.Index(spec, "~") < strings.Index(spec, "=")) {
		sep = "~"
		substring = true
	}
	parts := strings.SplitN(spec, sep, 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("highlight: invalid rule %q, expected field=value:color", spec)
	}
	field := strings.TrimSpace(parts[0])
	rest := parts[1]
	ci := strings.LastIndex(rest, ":")
	if ci < 0 {
		return Rule{}, fmt.Errorf("highlight: missing color in rule %q", spec)
	}
	value := rest[:ci]
	colorName := strings.ToLower(strings.TrimSpace(rest[ci+1:]))
	color, ok := namedColors[colorName]
	if !ok {
		return Rule{}, fmt.Errorf("highlight: unknown color %q in rule %q", colorName, spec)
	}
	return Rule{Field: field, Value: value, Color: color, Substring: substring}, nil
}
