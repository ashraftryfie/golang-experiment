package config

import (
	"errors"
	"testing"
	"time"
)

func TestServerConfig(t *testing.T) {
	cfg, err := NewServerConfig("127.0.0.1", 8080, 30)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: NewServerConfig is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Host() != "127.0.0.1" {
		t.Errorf("Host() = %q, want '127.0.0.1'", cfg.Host())
	}
	if cfg.Port() != 8080 {
		t.Errorf("Port() = %d, want 8080", cfg.Port())
	}
	if cfg.Timeout() != 30*time.Second {
		t.Errorf("Timeout() = %v, want 30s", cfg.Timeout())
	}

	// Validation tests
	if _, err := NewServerConfig("", 8080, 10); !errors.Is(err, ErrEmptyHost) {
		t.Errorf("expected ErrEmptyHost, got %v", err)
	}
	if _, err := NewServerConfig("localhost", 70000, 10); !errors.Is(err, ErrInvalidPort) {
		t.Errorf("expected ErrInvalidPort, got %v", err)
	}
	if _, err := NewServerConfig("localhost", 80, 0); !errors.Is(err, ErrInvalidTimeout) {
		t.Errorf("expected ErrInvalidTimeout, got %v", err)
	}
}
