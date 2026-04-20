package limit

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries(n int) []*parser.Entry {
	entries := make([]*parser.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = &parser.Entry{
			Timestamp: time.Unix(int64(i), 0),
			Fields:    map[string]interface{}{"index": i},
		}
	}
	return entries
}

func TestRun_NoLimit(t *testing.T) {
	entries := makeEntries(5)
	result := Run(entries, Config{})
	if len(result) != 5 {
		t.Errorf("expected 5 entries, got %d", len(result))
	}
}

func TestRun_MaxCapsEntries(t *testing.T) {
	entries := makeEntries(10)
	result := Run(entries, Config{Max: 3})
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
}

func TestRun_MaxLargerThanEntries(t *testing.T) {
	entries := makeEntries(4)
	result := Run(entries, Config{Max: 100})
	if len(result) != 4 {
		t.Errorf("expected 4 entries, got %d", len(result))
	}
}

func TestRun_OffsetSkipsEntries(t *testing.T) {
	entries := makeEntries(5)
	result := Run(entries, Config{Offset: 2})
	if len(result) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result))
	}
	if result[0].Fields["index"] != 2 {
		t.Errorf("expected first entry index=2, got %v", result[0].Fields["index"])
	}
}

func TestRun_OffsetAndMax(t *testing.T) {
	entries := makeEntries(10)
	result := Run(entries, Config{Offset: 3, Max: 4})
	if len(result) != 4 {
		t.Errorf("expected 4 entries, got %d", len(result))
	}
	if result[0].Fields["index"] != 3 {
		t.Errorf("expected first entry index=3, got %v", result[0].Fields["index"])
	}
}

func TestRun_OffsetBeyondLength(t *testing.T) {
	entries := makeEntries(3)
	result := Run(entries, Config{Offset: 10})
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := ParseConfig([]string{"max=5", "offset=2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Max != 5 || cfg.Offset != 2 {
		t.Errorf("expected max=5 offset=2, got max=%d offset=%d", cfg.Max, cfg.Offset)
	}
}

func TestParseConfig_NegativeValue(t *testing.T) {
	_, err := ParseConfig([]string{"max=-1"})
	if err == nil {
		t.Error("expected error for negative value")
	}
}

func TestParseConfig_UnknownOption(t *testing.T) {
	_, err := ParseConfig([]string{"limit=5"})
	if err == nil {
		t.Error("expected error for unknown option")
	}
}

func TestParseConfig_MissingEquals(t *testing.T) {
	_, err := ParseConfig([]string{"max5"})
	if err == nil {
		t.Error("expected error for missing equals")
	}
}
