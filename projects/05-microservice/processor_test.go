package orderprocessor

import (
	"context"
	"errors"
	"testing"
)

func TestOrderProcessor_IdempotencyAndTransitions(t *testing.T) {
	ctx := context.Background()
	processor := NewOrderProcessor()

	// 1. Initial process
	ord1, err := processor.ProcessOrder(ctx, "ord-1", "idem-key-123", 4999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ord1.State != StateConfirmed {
		t.Fatalf("expected CONFIRMED state, got %s", ord1.State)
	}

	// 2. Duplicate request with same idempotency key
	dup, err := processor.ProcessOrder(ctx, "ord-2", "idem-key-123", 4999)
	if !errors.Is(err, ErrDuplicateRequest) {
		t.Fatalf("expected ErrDuplicateRequest, got %v", err)
	}
	if dup.ID != "ord-1" {
		t.Fatalf("expected duplicate to return original ord-1, got %s", dup.ID)
	}

	// 3. Fulfill order
	if err := processor.FulfillOrder(ctx, "ord-1"); err != nil {
		t.Fatalf("failed to fulfill order: %v", err)
	}

	stored, err := processor.GetOrder(ctx, "ord-1")
	if err != nil {
		t.Fatalf("failed getting order: %v", err)
	}
	if stored.State != StateFulfilled {
		t.Fatalf("expected state FULFILLED, got %s", stored.State)
	}

	// 4. Cannot fulfill twice
	if err := processor.FulfillOrder(ctx, "ord-1"); !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus when fulfilling again, got %v", err)
	}
}
