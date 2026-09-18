package contextpropagation

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestExecutePipeline_AllSuccess(t *testing.T) {
	var count int32
	step := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	err := ExecutePipeline(context.Background(), []StepFunc{step, step, step})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if atomic.LoadInt32(&count) != 3 {
		t.Fatalf("expected 3 steps executed, got %d", count)
	}
}

func TestExecutePipeline_StepFailureAborts(t *testing.T) {
	errBoom := errors.New("boom")
	var count int32

	s1 := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}
	s2 := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return errBoom
	}
	s3 := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	err := ExecutePipeline(context.Background(), []StepFunc{s1, s2, s3})
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected errBoom, got %v", err)
	}
	if atomic.LoadInt32(&count) != 2 {
		t.Fatalf("expected 2 steps executed before failure, got %d", count)
	}
}

func TestExecutePipeline_CancelledBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var executed bool
	step := func(ctx context.Context) error {
		executed = true
		return nil
	}

	err := ExecutePipeline(ctx, []StepFunc{step})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if executed {
		t.Errorf("step should not have been executed")
	}
}

func TestRunWithTimeout_TimeoutAbortsPipeline(t *testing.T) {
	var executedCount int32

	stepSlow := func(ctx context.Context) error {
		atomic.AddInt32(&executedCount, 1)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return nil
		}
	}

	stepShouldNotRun := func(ctx context.Context) error {
		atomic.AddInt32(&executedCount, 1)
		return nil
	}

	err := RunWithTimeout(context.Background(), 20*time.Millisecond, []StepFunc{stepSlow, stepShouldNotRun})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
	if atomic.LoadInt32(&executedCount) != 1 {
		t.Errorf("expected only 1 step attempted, got %d", executedCount)
	}
}
