package tail

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(n int) []*parser.Entry {
	entries := make([]*parser.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = &parser.Entry{
			Timestamp: time.Unix(int64(i), 0),
			Raw:       map[string]interface{}{"i": i},
		}
	}
	return entries
}

func TestTail_LastN(t *testing.T) {
	entries := makeEntries(10)
	got := Run(entries, Options{N: 3})
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
	if got[0].Raw["i"] != 7 {
		t.Errorf("expected first element index 7, got %v", got[0].Raw["i"])
	}
}

func TestTail_ZeroN(t *testing.T) {
	entries := makeEntries(5)
	got := Run(entries, Options{N: 0})
	if len(got) != 5 {
		t.Fatalf("expected all 5, got %d", len(got))
	}
}

func TestTail_NGreaterThanLen(t *testing.T) {
	entries := makeEntries(3)
	got := Run(entries, Options{N: 100})
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func TestHead_FirstN(t *testing.T) {
	entries := makeEntries(10)
	got := Head(entries, Options{N: 4})
	if len(got) != 4 {
		t.Fatalf("expected 4, got %d", len(got))
	}
	if got[0].Raw["i"] != 0 {
		t.Errorf("expected first element index 0, got %v", got[0].Raw["i"])
	}
}

func TestParseConfig_Tail(t *testing.T) {
	cfg, err := ParseConfig("tail=5")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Mode != ModeTail || cfg.N != 5 {
		t.Errorf("unexpected config %+v", cfg)
	}
}

func TestParseConfig_Head(t *testing.T) {
	cfg, err := ParseConfig("head=10")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Mode != ModeHead || cfg.N != 10 {
		t.Errorf("unexpected config %+v", cfg)
	}
}

func TestParseConfig_Invalid(t *testing.T) {
	_, err := ParseConfig("notvalid")
	if err == nil {
		t.Error("expected error for missing equals")
	}
}

func TestParseConfig_BadN(t *testing.T) {
	_, err := ParseConfig("tail=abc")
	if err == nil {
		t.Error("expected error for non-integer N")
	}
}
