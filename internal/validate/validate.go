// Package validate provides field validation for log entries.
// It checks that required fields are present and that field values
// match expected types or patterns.
package validate

import (
	"fmt"
	"regexp"

	"github.com/yourorg/logslice/internal/parser"
)

// Result holds the outcome of validating a single log entry.
type Result struct {
	Entry  parser.Entry
	Errors []string
}

// Rule describes a single validation constraint.
type Rule struct {
	Field    string
	Required bool
	Pattern  *regexp.Regexp // optional: value must match
	TypeName string         // optional: "string", "number", "bool"
}

// Run validates each entry against the provided rules.
// It returns one Result per entry; entries with no errors are still included.
func Run(entries []parser.Entry, rules []Rule) []Result {
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		var errs []string
		for _, r := range rules {
			val, ok := e.Fields[r.Field]
			if !ok || val == nil {
				if r.Required {
					errs = append(errs, fmt.Sprintf("missing required field %q", r.Field))
				}
				continue
			}
			str := fmt.Sprintf("%v", val)
			if r.Pattern != nil && !r.Pattern.MatchString(str) {
				errs = append(errs, fmt.Sprintf("field %q value %q does not match pattern %s", r.Field, str, r.Pattern))
			}
			if r.TypeName != "" {
				if typeErr := checkType(r.Field, val, r.TypeName); typeErr != "" {
					errs = append(errs, typeErr)
				}
			}
		}
		results = append(results, Result{Entry: e, Errors: errs})
	}
	return results
}

func checkType(field string, val interface{}, typeName string) string {
	switch typeName {
	case "string":
		if _, ok := val.(string); !ok {
			return fmt.Sprintf("field %q expected string, got %T", field, val)
		}
	case "number":
		switch val.(type) {
		case float64, int, int64:
		default:
			return fmt.Sprintf("field %q expected number, got %T", field, val)
		}
	case "bool":
		if _, ok := val.(bool); !ok {
			return fmt.Sprintf("field %q expected bool, got %T", field, val)
		}
	}
	return ""
}
