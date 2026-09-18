package contextpropagation

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestExecutePipeline_Starter(t *testing.T) {
	var count int32
	step1 := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	err := ExecutePipeline(context.Background(), []StepFunc{step1})
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: ExecutePipeline")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}
}
