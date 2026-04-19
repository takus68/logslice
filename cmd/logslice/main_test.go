package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp log: %v", err)
	}
	return p
}

const sampleLog = `{"timestamp":"2024-03-01T10:00:00Z","level":"info","msg":"started"}
{"timestamp":"2024-03-01T11:00:00Z","level":"error","msg":"failed"}
{"timestamp":"2024-03-01T12:00:00Z","level":"info","msg":"done"}
`

func buildBinary(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	bin := filepath.Join(dir, "logslice")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return bin
}

func TestMain_NoFile(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, _ := cmd.CombinedOutput()
	if !strings.Contains(string(out), "-file is required") {
		t.Errorf("expected usage error, got: %s", out)
	}
}

func TestMain_FilterByField(t *testing.T) {
	bin := buildBinary(t)
	log := writeTempLog(t, sampleLog)
	cmd := exec.Command(bin, "-file", log, "-field", "level", "-value", "error", "-format", "compact")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "failed") {
		t.Errorf("expected 'failed' in output, got: %s", out)
	}
	if strings.Contains(string(out), "started") {
		t.Errorf("expected 'started' to be filtered out")
	}
}

func TestMain_FilterByTimeRange(t *testing.T) {
	bin := buildBinary(t)
	log := writeTempLog(t, sampleLog)
	cmd := exec.Command(bin, "-file", log,
		"-start", "2024-03-01T10:30:00Z",
		"-end", "2024-03-01T11:30:00Z",
		"-format", "compact")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "failed") {
		t.Errorf("expected 'failed' entry in output")
	}
	if strings.Contains(string(out), "started") || strings.Contains(string(out), "done") {
		t.Errorf("expected only middle entry, got: %s", out)
	}
}
