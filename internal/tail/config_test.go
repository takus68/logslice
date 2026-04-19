package tail

import (
	"testing"
)

func TestParseConfig_DefaultTail(t *testing.T) {
	cfg, err := ParseConfig(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.N != 10 {
		t.Errorf("expected default N=10, got %d", cfg.N)
	}
	if cfg.Head {
		t.Error("expected Head=false by default")
	}
}

func TestParseConfig_CustomN(t *testing.T) {
	cfg, err := ParseConfig(map[string]string{"n": "25"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.N != 25 {
		t.Errorf("expected N=25, got %d", cfg.N)
	}
}

func TestParseConfig_HeadMode(t *testing.T) {
	cfg, err := ParseConfig(map[string]string{"head": "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Head {
		t.Error("expected Head=true")
	}
}

func TestParseConfig_InvalidN(t *testing.T) {
	_, err := ParseConfig(map[string]string{"n": "abc"})
	if err == nil {
		t.Error("expected error for invalid n")
	}
}

func TestParseConfig_NegativeN(t *testing.T) {
	_, err := ParseConfig(map[string]string{"n": "-5"})
	if err == nil {
		t.Error("expected error for negative n")
	}
}

func TestParseConfig_ZeroN(t *testing.T) {
	_, err := ParseConfig(map[string]string{"n": "0"})
	if err == nil {
		t.Error("expected error for n=0")
	}
}
