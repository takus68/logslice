package flatten

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) *parser.Entry {
	return &parser.Entry{
		Timestamp: time.Now(),
		Raw:       "{}",
		Fields:    fields,
	}
}

func TestRun_FlatEntry(t *testing.T) {
	entries := []*parser.Entry{makeEntry(map[string]interface{}{"level": "info", "msg": "ok"})}
	out := Run(entries, ".")
	if out[0].Fields["level"] != "info" {
		t.Errorf("expected info, got %v", out[0].Fields["level"])
	}
}

func TestRun_NestedEntry(t *testing.T) {
	entries := []*parser.Entry{makeEntry(map[string]interface{}{
		"http": map[string]interface{}{"method": "GET", "status": "200"},
	})}
	out := Run(entries, ".")
	if out[0].Fields["http.method"] != "GET" {
		t.Errorf("expected GET, got %v", out[0].Fields["http.method"])
	}
	if out[0].Fields["http.status"] != "200" {
		t.Errorf("expected 200, got %v", out[0].Fields["http.status"])
	}
}

func TestRun_DeeplyNested(t *testing.T) {
	entries := []*parser.Entry{makeEntry(map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "deep",
			},
		},
	})}
	out := Run(entries, ".")
	if out[0].Fields["a.b.c"] != "deep" {
		t.Errorf("expected deep, got %v", out[0].Fields["a.b.c"])
	}
}

func TestRun_CustomSeparator(t *testing.T) {
	entries := []*parser.Entry{makeEntry(map[string]interface{}{
		"req": map[string]interface{}{"id": "abc"},
	})}
	out := Run(entries, "_")
	if out[0].Fields["req_id"] != "abc" {
		t.Errorf("expected abc, got %v", out[0].Fields["req_id"])
	}
}

func TestRun_DefaultSeparator(t *testing.T) {
	entries := []*parser.Entry{makeEntry(map[string]interface{}{
		"x": map[string]interface{}{"y": "val"},
	})}
	out := Run(entries, "")
	if out[0].Fields["x.y"] != "val" {
		t.Errorf("expected val, got %v", out[0].Fields["x.y"])
	}
}

func TestRun_Empty(t *testing.T) {
	out := Run([]*parser.Entry{}, ".")
	if len(out) != 0 {
		t.Errorf("expected empty, got %d", len(out))
	}
}
