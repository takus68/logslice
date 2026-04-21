package rename

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

func makeEntries(fields ...map[string]interface{}) []parser.Entry {
	var entries []parser.Entry
	for _, f := range fields {
		entries = append(entries, parser.Entry{
			Timestamp: time.Now(),
			Fields:    f,
		})
	}
	return entries
}

func TestRun_RenamesField(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"msg": "hello", "level": "info"})
	rules := []Rule{{From: "msg", To: "message"}}
	out := Run(entries, rules)
	if _, ok := out[0].Fields["message"]; !ok {
		t.Error("expected 'message' key to exist")
	}
	if _, ok := out[0].Fields["msg"]; ok {
		t.Error("expected 'msg' key to be removed")
	}
}

func TestRun_MissingSourceKey(t *testing.T) {
	entries := make{}{"level": "info"})
	rules := []Rule{{From: "msg", To: "message"}}
	out := Run(entries, rules)
	if _, ok := out[0].Fields["message"]; ok {
		t.Error("expected 'message' key to not be added")
	}
}

func TestRun_MultipleRules(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"msg": "hi", "lvl": "warn"})
	rules := []Rule{{From: "msg", To: "message"}, {From: "lvl", To: "level"}}
	out := Run(entries, rules)
	if out[0].Fields["message"] != "hi" {
		t.Errorf("expected message=hi, got %v", out[0].Fields["message"])
	}
	if out[0].Fields["level"] != "warn" {
		t.Errorf("expected level=warn, got %v", out[0].Fields["level"])
	}
}

func TestRun_EmptyRules(t *testing.T) {
	entries := makeEntries(map[string]interface{}{"msg": "hello"})
	out := Run(entries, nil)
	if out[0].Fields["msg"] != "hello" {
		t.Error("expected entry to be unchanged with no rules")
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	origFields := map[string]interface{}{"msg": "hello"}
	entries := makeEntries(origFields)
	rules := []Rule{{From: "msg", To: "message"}}
	Run(entries, rules)
	if _, ok := entries[0].Fields["msg"]; !ok {
		t.Error("original entry should not be mutated")
	}
}
