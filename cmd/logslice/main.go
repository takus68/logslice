package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
)

func main() {
	var (
		filePath  = flag.String("file", "", "Path to log file (required)")
		startStr  = flag.String("start", "", "Start time (RFC3339), e.g. 2024-01-01T00:00:00Z")
		endStr    = flag.String("end", "", "End time (RFC3339), e.g. 2024-12-31T23:59:59Z")
		fieldKey  = flag.String("field", "", "Field key to filter on")
		fieldVal  = flag.String("value", "", "Field value to match (substring, case-insensitive)")
		fmt_      = flag.String("format", "json", "Output format: json, pretty, compact")
	)
	flag.Parse()

	if *filePath == "" {
		fmt.Fprintln(os.Stderr, "error: -file is required")
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	entries, err := parser.Parse(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing file: %v\n", err)
		os.Exit(1)
	}

	var filters []filter.FilterFunc

	if *startStr != "" || *endStr != "" {
		var start, end time.Time
		if *startStr != "" {
			start, err = time.Parse(time.RFC3339, *startStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid -start: %v\n", err)
				os.Exit(1)
			}
		}
		if *endStr != "" {
			end, err = time.Parse(time.RFC3339, *endStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "invalid -end: %v\n", err)
				os.Exit(1)
			}
		}
		filters = append(filters, filter.ByTimeRange(start, end))
	}

	if *fieldKey != "" {
		filters = append(filters, filter.ByField(*fieldKey, *fieldVal))
	}

	result := filter.Apply(entries, filters...)

	if err := output.Write(os.Stdout, result, *fmt_); err != nil {
		fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
		os.Exit(1)
	}
}
