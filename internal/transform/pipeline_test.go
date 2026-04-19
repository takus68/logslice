package transform

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func makePipelineEntries() []parser.Entry {
	t := time.Now()
	return []parser.Entry{
		{Timestamp: t, Fields: map[string]interface{}{"Level": "info", "Msg": "start", "svc": "web"}},
	}
}

func TestPipeline_Empty(t *testing.T) {
	entries := makePipelineEntries()
	p := NewPipeline()
	out := p.Run(entries)
	if len(out) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(out))
	}
}

func TestPipeline_SingleStep(t *testing.T) {
	entries := makePipelineEntries()
	p := NewPipeline().Add(FieldNormalize)
	out := p.Run(entries)
	if _, ok := out[0].Fields["level"]; !ok {
		t.Error("expected 'level' after normalize step")
	}
}

func TestPipeline_MultiStep(t *testing.T) {
	entries := makePipelineEntries()
	p := NewPipeline().
		Add(FieldNormalize).
		Add(func(e []parser.Entry) []parser.Entry { return FieldDrop(e, "svc") }).
		Add(func(e []parser.Entry) []parser.Entry { return FieldRename(e, "msg", "message") })
	out := p.Run(entries)
	e := out[0]
	if _, ok := e.Fields["svc"]; ok {
		t.Error("'svc' should have been dropped")
	}
	if _, ok := e.Fields["message"]; !ok {
		t.Error("'message' should exist after rename")
	}
	if _, ok := e.Fields["level"]; !ok {
		t.Error("'level' should exist after normalize")
	}
}

func TestPipeline_ChainReturnsPointer(t *testing.T) {
	p := NewPipeline()
	if p.Add(FieldNormalize) != p {
		t.Error("Add should return the same pipeline pointer")
	}
}
