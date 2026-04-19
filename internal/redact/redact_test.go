package redact

import (
	"regexp"
	"testing"
	"time"

	"github.com/yourusername/logslice/internal/parser"
)

func makeEntries() []parser.Entry {
	now := time.Now()
	return []parser.Entry{
		{Timestamp: now, Raw: map[string]interface{}{"level": "info", "password": "secret123", "email": "user@example.com", "msg": "login"}},
		{Timestamp: now, Raw: map[string]interface{}{"level": "warn", "token": "abc-token", "msg": "request"}},
	}
}

func TestRun_RedactField(t *testing.T) {
	entries := makeEntries()
	cfg := Config{Fields: []string{"password"}}
	out := Run(entries, cfg)
	if out[0].Raw["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", out[0].Raw["password"])
	}
	if out[0].Raw["email"] == "[REDACTED]" {
		t.Error("email should not be redacted")
	}
}

func TestRun_RedactMultipleFields(t *testing.T) {
	entries := makeEntries()
	cfg := Config{Fields: []string{"password", "token"}}
	out := Run(entries, cfg)
	if out[0].Raw["password"] != "[REDACTED]" {
		t.Error("password should be redacted")
	}
	if out[1].Raw["token"] != "[REDACTED]" {
		t.Error("token should be redacted")
	}
}

func TestRun_MissingField(t *testing.T) {
	entries := makeEntries()
	cfg := Config{Fields: []string{"nonexistent"}}
	out := Run(entries, cfg)
	if _, ok := out[0].Raw["nonexistent"]; ok {
		t.Error("nonexistent field should not appear")
	}
}

func TestRun_PatternMask(t *testing.T) {
	entries := makeEntries()
	cfg := Config{
		Patterns: map[string]*regexp.Regexp{
			"email": regexp.MustCompile(`[^@]+`),
		},
	}
	out := Run(entries, cfg)
	val, _ := out[0].Raw["email"].(string)
	if val == "user@example.com" {
		t.Error("email should be masked")
	}
}

func TestRun_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries()
	cfg := Config{Fields: []string{"password"}}
	Run(entries, cfg)
	if entries[0].Raw["password"] == "[REDACTED]" {
		t.Error("original entry should not be mutated")
	}
}
