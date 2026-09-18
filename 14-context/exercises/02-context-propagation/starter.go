package contextpropagation

import (
	"context"
	"errors"
	"time"
)

var ErrNotImplemented = errors.New("exercise not implemented yet")

type StepFunc func(ctx context.Context) error

// ExecutePipeline executes a list of steps sequentially, respecting context cancellation.
// If the context is canceled before or during any step, it halts and returns ctx.Err().
func ExecutePipeline(ctx context.Context, steps []StepFunc) error {
	return ErrNotImplemented
}

// RunWithTimeout wraps ExecutePipeline in a child context with the given timeout.
func RunWithTimeout(ctx context.Context, timeout time.Duration, steps []StepFunc) error {
	return ErrNotImplemented
}
