package transform

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntries() []parser.Entry {
	t := time.Now()
	return []parser.Entry{
		{Timestamp: t, Fields: map[string]interface{}{"Level": "info", "msg": "hello", "svc": "api"}},
		{Timestamp: t, Fields: map[string]interface{}{"Level": "error", "msg": "oops", "svc": "db"}},
	}
}

func TestFieldRename(t *testing.T) {
	entries := makeEntries()
	out := FieldRename(entries, "msg", "message")
	for _, e := range out {
		if _, ok := e.Fields["msg"]; ok {
			t.Error("old key 'msg' should not exist")
		}
		if _, ok := e.Fields["message"]; !ok {
			t.Error("new key 'message' should exist")
		}
	}
}

func TestFieldRename_MissingKey(t *testing.T) {
	entries := makeEntries()
	out := FieldRename(entries, "nonexistent", "new")
	for i, e := range out {
		if len(e.Fields) != len(entries[i].Fields) {
			t.Error("field count should be unchanged when key not found")
		}
	}
}

func TestFieldDrop(t *testing.T) {
	entries := makeEntries()
	out := FieldDrop(entries, "svc", "Level")
	for _, e := range out {
		if _, ok := e.Fields["svc"]; ok {
			t.Error("'svc' should be dropped")
		}
		if _, ok := e.Fields["Level"]; ok {
			t.Error("'Level' should be dropped")
		}
		if _, ok := e.Fields["msg"]; !ok {
			t.Error("'msg' should remain")
		}
	}
}

func TestFieldDrop_NoKeys(t *testing.T) {
	entries := makeEntries()
	out := FieldDrop(entries)
	for i, e := range out {
		if len(e.Fields) != len(entries[i].Fields) {
			t.Error("field count should be unchanged with no keys to drop")
		}
	}
}

func TestFieldNormalize(t *testing.T) {
	entries := makeEntries()
	out := FieldNormalize(entries)
	for _, e := range out {
		if _, ok := e.Fields["Level"]; ok {
			t.Error("'Level' should be normalized to 'level'")
		}
		if _, ok := e.Fields["level"]; !ok {
			t.Error("'level' should exist after normalization")
		}
	}
}
