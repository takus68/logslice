package split

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries(fields []map[string]interface{}) []*parser.LogEntry {
	entries := make([]*parser.LogEntry, len(fields))
	for i, f := range fields {
		entries[i] = &parser.LogEntry{Timestamp: time.Now(), Raw: f}
	}
	return entries
}

func TestRun_SingleChunkDefault(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"msg": "a"}, {"msg": "b"}, {"msg": "c"},
	})
	chunks := Run(entries, Config{})
	if len(chunks) != 1 || len(chunks[0]) != 3 {
		t.Fatalf("expected 1 chunk of 3, got %d chunks", len(chunks))
	}
}

func TestRun_EmptyEntries(t *testing.T) {
	chunks := Run(nil, Config{ChunkSize: 2})
	if chunks != nil {
		t.Fatalf("expected nil, got %v", chunks)
	}
}

func TestRun_BySize_Even(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"msg": "a"}, {"msg": "b"}, {"msg": "c"}, {"msg": "d"},
	})
	chunks := Run(entries, Config{ChunkSize: 2})
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if len(chunks[0]) != 2 || len(chunks[1]) != 2 {
		t.Fatalf("unexpected chunk sizes")
	}
}

func TestRun_BySize_Remainder(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"msg": "a"}, {"msg": "b"}, {"msg": "c"},
	})
	chunks := Run(entries, Config{ChunkSize: 2})
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if len(chunks[1]) != 1 {
		t.Fatalf("expected last chunk of 1, got %d", len(chunks[1]))
	}
}

func TestRun_ByField(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"level": "info"}, {"level": "info"}, {"level": "error"}, {"level": "error"}, {"level": "info"},
	})
	chunks := Run(entries, Config{FieldBoundary: "level"})
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	if len(chunks[0]) != 2 || len(chunks[1]) != 2 || len(chunks[2]) != 1 {
		t.Fatalf("unexpected chunk sizes: %d %d %d", len(chunks[0]), len(chunks[1]), len(chunks[2]))
	}
}

func TestRun_ByField_MissingField(t *testing.T) {
	entries := makeEntries([]map[string]interface{}{
		{"msg": "a"}, {"msg": "b"},
	})
	chunks := Run(entries, Config{FieldBoundary: "level"})
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk when field missing, got %d", len(chunks))
	}
}

func TestParseConfig_Size(t *testing.T) {
	cfg, err := ParseConfig([]string{"size=10"})
	if err != nil || cfg.ChunkSize != 10 {
		t.Fatalf("expected size 10, got %d err %v", cfg.ChunkSize, err)
	}
}

func TestParseConfig_Boundary(t *testing.T) {
	cfg, err := ParseConfig([]string{"boundary=level"})
	if err != nil || cfg.FieldBoundary != "level" {
		t.Fatalf("expected boundary=level, got %q err %v", cfg.FieldBoundary, err)
	}
}

func TestParseConfig_InvalidSize(t *testing.T) {
	_, err := ParseConfig([]string{"size=0"})
	if err == nil {
		t.Fatal("expected error for size=0")
	}
}

func TestParseConfig_UnknownOption(t *testing.T) {
	_, err := ParseConfig([]string{"foo=bar"})
	if err == nil {
		t.Fatal("expected error for unknown option")
	}
}
