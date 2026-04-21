// Package pivot provides functionality to pivot log entries by a field,
// collecting values of another field into a grouped structure.
package pivot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Config holds the configuration for a pivot operation.
type Config struct {
	// KeyField is the field whose values become the pivot keys (columns).
	KeyField string
	// ValueField is the field whose values are collected under each key.
	ValueField string
	// GroupField is the field used to group rows (the row identifier).
	GroupField string
}

// ParseConfig parses pivot options from a slice of "key=value" strings.
// Required options: key=<field>, value=<field>, group=<field>
func ParseConfig(opts []string) (Config, error) {
	var cfg Config
	for _, opt := range opts {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 {
			return cfg, fmt.Errorf("pivot: invalid option %q, expected key=value", opt)
		}
		k, v := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch k {
		case "key":
			cfg.KeyField = v
		case "value":
			cfg.ValueField = v
		case "group":
			cfg.GroupField = v
		default:
			return cfg, fmt.Errorf("pivot: unknown option %q", k)
		}
	}
	if cfg.KeyField == "" {
		return cfg, fmt.Errorf("pivot: missing required option 'key'")
	}
	if cfg.ValueField == "" {
		return cfg, fmt.Errorf("pivot: missing required option 'value'")
	}
	if cfg.GroupField == "" {
		return cfg, fmt.Errorf("pivot: missing required option 'group'")
	}
	return cfg, nil
}

// Run pivots the entries according to cfg. It returns one entry per unique
// GroupField value. Each returned entry contains the group value plus one
// field per unique KeyField value, set to the corresponding ValueField value
// from the last matching source entry.
func Run(entries []parser.Entry, cfg Config) []parser.Entry {
	// Maintain insertion order for group keys.
	groupOrder := []string{}
	groupSeen := map[string]bool{}
	// rows maps groupValue -> (keyValue -> valueValue)
	rows := map[string]map[string]interface{}{}

	for _, e := range entries {
		groupVal, ok := e.Fields[cfg.GroupField]
		if !ok {
			continue
		}
		groupStr := fmt.Sprintf("%v", groupVal)

		keyVal, ok := e.Fields[cfg.KeyField]
		if !ok {
			continue
		}
		keyStr := fmt.Sprintf("%v", keyVal)

		valVal := e.Fields[cfg.ValueField]

		if !groupSeen[groupStr] {
			groupOrder = append(groupOrder, groupStr)
			groupSeen[groupStr] = true
			rows[groupStr] = map[string]interface{}{}
		}
		rows[groupStr][keyStr] = valVal
	}

	// Collect all key column names in sorted order for determinism.
	keySet := map[string]bool{}
	for _, row := range rows {
		for k := range row {
			keySet[k] = true
		}
	}
	allKeys := make([]string, 0, len(keySet))
	for k := range keySet {
		allKeys = append(allKeys, k)
	}
	sort.Strings(allKeys)

	result := make([]parser.Entry, 0, len(groupOrder))
	for _, g := range groupOrder {
		fields := map[string]interface{}{
			cfg.GroupField: g,
		}
		for _, k := range allKeys {
			if v, ok := rows[g][k]; ok {
				fields[k] = v
			}
		}
		result = append(result, parser.Entry{Fields: fields})
	}
	return result
}
