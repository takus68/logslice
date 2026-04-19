package sample

import "testing"

func TestParseOptions_Nth(t *testing.T) {
	opts, err := ParseOptions("nth:3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Strategy != StrategyNth || opts.N != 3 {
		t.Fatalf("unexpected opts: %+v", opts)
	}
}

func TestParseOptions_NthMissingN(t *testing.T) {
	_, err := ParseOptions("nth")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseOptions_NthInvalidN(t *testing.T) {
	_, err := ParseOptions("nth:abc")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseOptions_Random(t *testing.T) {
	opts, err := ParseOptions("random:0.25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Strategy != StrategyRandom || opts.Rate != 0.25 {
		t.Fatalf("unexpected opts: %+v", opts)
	}
}

func TestParseOptions_RandomWithSeed(t *testing.T) {
	opts, err := ParseOptions("random:0.5:seed=42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Seed != 42 {
		t.Fatalf("expected seed 42, got %d", opts.Seed)
	}
}

func TestParseOptions_RandomInvalidRate(t *testing.T) {
	_, err := ParseOptions("random:1.5")
	if err == nil {
		t.Fatal("expected error for rate > 1")
	}
}

func TestParseOptions_UnknownStrategy(t *testing.T) {
	_, err := ParseOptions("reservoir:10")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseOptions_Empty(t *testing.T) {
	_, err := ParseOptions("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}
