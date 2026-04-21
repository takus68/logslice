// Package typecheck provides functionality to inspect and report
// the inferred types of field values across log entries.
package typecheck

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// FieldType represents the inferred type of a field value.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeFloat   FieldType = "float"
	TypeBool    FieldType = "bool"
	TypeNull    FieldType = "null"
	TypeUnknown FieldType = "unknown"
)

// FieldReport holds the observed types and count for a single field.
type FieldReport struct {
	Field  string
	Types  map[FieldType]int
	Total  int
}

// Run inspects each entry and returns a report of inferred types per field.
// If fields is non-empty, only those fields are inspected.
func Run(entries []parser.Entry, fields []string) map[string]*FieldReport {
	reports := make(map[string]*FieldReport)

	for _, e := range entries {
		for k, v := range e.Fields {
			if len(fields) > 0 && !contains(fields, k) {
				continue
			}
			if _, ok := reports[k]; !ok {
				reports[k] = &FieldReport{
					Field: k,
					Types: make(map[FieldType]int),
				}
			}
			t := inferType(fmt.Sprintf("%v", v))
			reports[k].Types[t]++
			reports[k].Total++
		}
	}
	return reports
}

// inferType returns the best-guess FieldType for a raw string value.
func inferType(raw string) FieldType {
	if raw == "" || raw == "null" || raw == "<nil>" {
		return TypeNull
	}
	if raw == "true" || raw == "false" {
		return TypeBool
	}
	if _, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return TypeInt
	}
	if _, err := strconv.ParseFloat(raw, 64); err == nil {
		return TypeFloat
	}
	return TypeString
}

func contains(fields []string, key string) bool {
	for _, f := range fields {
		if strings.EqualFold(f, key) {
			return true
		}
	}
	return false
}
