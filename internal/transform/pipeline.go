package transform

import "github.com/logslice/logslice/internal/parser"

// TransformFunc is a function that transforms a slice of log entries.
type TransformFunc func([]parser.Entry) []parser.Entry

// Pipeline applies a sequence of TransformFuncs to a slice of entries.
type Pipeline struct {
	steps []TransformFunc
}

// NewPipeline creates a new empty Pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Add appends a TransformFunc to the pipeline.
func (p *Pipeline) Add(fn TransformFunc) *Pipeline {
	p.steps = append(p.steps, fn)
	return p
}

// Run executes all pipeline steps in order and returns the final entries.
func (p *Pipeline) Run(entries []parser.Entry) []parser.Entry {
	result := entries
	for _, step := range p.steps {
		result = step(result)
	}
	return result
}
