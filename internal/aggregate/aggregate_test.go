package aggregate

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntries() []*parser.Entry {
	now := time.Now()
	return []*parser.Entry{
		{Timestamp: now, Fields: map[string]interface{}{"level": "info", "svc": "api"}},
		{Timestamp: now, Fields: map[string]interface{}{"level": "error", "svc": "api"}},
		{Timestamp: now, Fields: map[string]interface{}{"level": "info", "svc": "worker"}},
		{Timestamp: now, Fields: map[string]interface{}{"level": "warn"}},
		{Timestamp: now, Fields: map[string]interface{}{"level": "info"}},
	}
}

func TestByField_Counts(t *testing.T) {
	entries := makeEntries()
	res, err := ByField(entries, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Counts["info"] != 3 {
		t.Errorf("expected 3 info, got %d", res.Counts["info"])
	}
	if res.Counts["error"] != 1 {
		t.Errorf("expected 1 error, got %d", res.Counts["error"])
	}
	if res.Total != 5 {
		t.Errorf("expected total 5, got %d", res.Total)
	}
}

func TestByField_MissingField(t *testing.T) {
	entries := makeEntries()
	res, err := ByField(entries, "svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Counts["<missing>"] != 2 {
		t.Errorf("expected 2 missing, got %d", res.Counts["<missing>"])
	}
}

func TestByField_EmptyField(t *testing.T) {
	_, err := ByField(makeEntries(), "")
	if err == nil {
		t.Error("expected error for empty field name")
	}
}

func TestByField_EmptyEntries(t *testing.T) {
	res, err := ByField([]*parser.Entry{}, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 0 {
		t.Errorf("expected total 0, got %d", res.Total)
	}
}

func TestSortedKeys(t *testing.T) {
	res, _ := ByField(makeEntries(), "level")
	keys := res.SortedKeys()
	if len(keys) == 0 {
		t.Fatal("expected non-empty keys")
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
		}
	}
}
