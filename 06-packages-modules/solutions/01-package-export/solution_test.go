package config

import (
	"testing"
	"time"
)

func TestSolutionConfig(t *testing.T) {
	cfg, err := NewServerConfig("0.0.0.0", 443, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Host() != "0.0.0.0" {
		t.Errorf("got %s, want '0.0.0.0'", cfg.Host())
	}
	if cfg.Port() != 443 {
		t.Errorf("got %d, want 443", cfg.Port())
	}
	if cfg.Timeout() != 60*time.Second {
		t.Errorf("got %v, want 60s", cfg.Timeout())
	}
}
