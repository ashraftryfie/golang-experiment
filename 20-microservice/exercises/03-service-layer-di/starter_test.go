package servicelayer

import (
	"context"
	"errors"
	"testing"
)

func TestStarter_NotImplemented(t *testing.T) {
	svc := NewTransferService(nil, nil)
	err := svc.ExecuteTransfer(context.Background(), "a1", "a2", 100)
	if !errors.Is(err, ErrNotImplemented) {
		t.Logf("Exercise has been worked on or implemented. Current error: %v", err)
	}
}
