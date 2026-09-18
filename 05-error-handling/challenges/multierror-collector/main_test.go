package main

import (
	"errors"
	"strings"
	"testing"
)

var (
	ErrTimeout = errors.New("timeout connecting")
	ErrAuth    = errors.New("invalid credentials")
)

func TestMultiError(t *testing.T) {
	var m MultiError

	// Empty ErrorOrNil should return untyped nil
	if err := m.ErrorOrNil(); err != nil {
		t.Errorf("expected empty MultiError to return nil, got %v", err)
	}

	// Add errors
	m.Add(ErrTimeout)
	m.Add(nil) // Should be ignored
	m.Add(ErrAuth)

	if len(m.Errors()) != 2 {
		t.Errorf("expected 2 errors, got %d", len(m.Errors()))
	}

	err := m.ErrorOrNil()
	if err == nil {
		t.Fatalf("expected non-nil error")
	}

	// Test Go 1.20+ multi-error inspection with errors.Is
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("errors.Is should find ErrTimeout in MultiError")
	}
	if !errors.Is(err, ErrAuth) {
		t.Errorf("errors.Is should find ErrAuth in MultiError")
	}

	// Verify formatted output
	errMsg := err.Error()
	if !strings.Contains(errMsg, "encountered 2 errors") {
		t.Errorf("error output should contain summary count, got %s", errMsg)
	}
}
