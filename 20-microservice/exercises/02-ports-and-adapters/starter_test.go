package portsadapters

import (
	"errors"
	"testing"
)

func TestStarter_NotImplemented(t *testing.T) {
	_, err := NewProduct("prod-1", "SKU-ABC", 10)
	if !errors.Is(err, ErrNotImplemented) {
		t.Logf("Exercise has been worked on or implemented. Current error: %v", err)
	}
}
