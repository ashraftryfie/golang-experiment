package contextpropagation

import (
	"context"
	"time"
)

type StepFunc func(ctx context.Context) error

// ExecutePipeline executes a list of steps sequentially, respecting context cancellation.
// If the context is canceled before or during any step, it halts and returns ctx.Err().
func ExecutePipeline(ctx context.Context, steps []StepFunc) error {
	for _, step := range steps {
		// Proactive check before starting the step
		if err := ctx.Err(); err != nil {
			return err
		}

		// Execute step
		if err := step(ctx); err != nil {
			return err
		}
	}
	return nil
}

// RunWithTimeout wraps ExecutePipeline in a child context with the given timeout.
func RunWithTimeout(ctx context.Context, timeout time.Duration, steps []StepFunc) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return ExecutePipeline(ctxTimeout, steps)
}
