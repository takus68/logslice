package merge_test

import (
	"testing"

	"github.com/user/logslice/internal/merge"
)

func TestParseConfig_Defaults(t *testing.T) {
	cfg, err := merge.ParseConfig([]string{"a.log", "b.log"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Stable {
		t.Error("expected stable=false by default")
	}
	if len(cfg.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(cfg.Files))
	}
}

func TestParseConfig_StableTrue(t *testing.T) {
	cfg, err := merge.ParseConfig(nil, "stable=true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Stable {
		t.Error("expected stable=true")
	}
}

func TestParseConfig_StableFalse(t *testing.T) {
	cfg, err := merge.ParseConfig(nil, "stable=false")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Stable {
		t.Error("expected stable=false")
	}
}

func TestParseConfig_InvalidStableValue(t *testing.T) {
	_, err := merge.ParseConfig(nil, "stable=yes")
	if err == nil {
		t.Fatal("expected error for invalid stable value")
	}
}

func TestParseConfig_UnknownOption(t *testing.T) {
	_, err := merge.ParseConfig(nil, "foo=bar")
	if err == nil {
		t.Fatal("expected error for unknown option")
	}
}

func TestParseConfig_MissingEquals(t *testing.T) {
	_, err := merge.ParseConfig(nil, "stable")
	if err == nil {
		t.Fatal("expected error for missing equals")
	}
}
