package domainentities

import (
	"errors"
	"testing"
)

func TestStarter_NotImplemented(t *testing.T) {
	_, err := NewMoney(100, "USD")
	if !errors.Is(err, ErrNotImplemented) {
		t.Logf("Exercise has been worked on or implemented. Current error: %v", err)
	}
}
