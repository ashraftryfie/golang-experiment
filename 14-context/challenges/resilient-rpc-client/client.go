package resilientrpc

import (
	"context"
	"errors"
	"time"
)

var (
	ErrZeroAttempts = errors.New("max attempts must be at least 1")
)

type RPCCallFunc func(ctx context.Context) (any, error)

type RetryConfig struct {
	MaxAttempts       int
	PerAttemptTimeout time.Duration
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffFactor     float64
}

// DefaultRetryConfig provides sensible production defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:       3,
		PerAttemptTimeout: 100 * time.Millisecond,
		InitialBackoff:    10 * time.Millisecond,
		MaxBackoff:        100 * time.Millisecond,
		BackoffFactor:     2.0,
	}
}

type ResilientClient struct {
	config RetryConfig
}

func NewResilientClient(cfg RetryConfig) (*ResilientClient, error) {
	if cfg.MaxAttempts < 1 {
		return nil, ErrZeroAttempts
	}
	if cfg.BackoffFactor < 1.0 {
		cfg.BackoffFactor = 1.5
	}
	if cfg.InitialBackoff <= 0 {
		cfg.InitialBackoff = 10 * time.Millisecond
	}
	if cfg.MaxBackoff < cfg.InitialBackoff {
		cfg.MaxBackoff = cfg.InitialBackoff
	}
	return &ResilientClient{config: cfg}, nil
}

// Invoke attempts to execute fn up to MaxAttempts with per-attempt timeouts and exponential backoff.
// It stops immediately if ctx is cancelled or times out.
func (c *ResilientClient) Invoke(ctx context.Context, fn RPCCallFunc) (any, error) {
	var lastErr error
	currentBackoff := c.config.InitialBackoff

	for attempt := 1; attempt <= c.config.MaxAttempts; attempt++ {
		// Proactive check before launching attempt
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		// Perform attempt with child timeout context
		res, err := c.executeAttempt(ctx, fn)
		if err == nil {
			return res, nil
		}

		lastErr = err

		// If parent context is dead, don't retry
		if parentErr := ctx.Err(); parentErr != nil {
			return nil, parentErr
		}

		// Don't sleep after the final attempt
		if attempt == c.config.MaxAttempts {
			break
		}

		// Context-aware sleep for backoff
		if err := c.sleep(ctx, currentBackoff); err != nil {
			return nil, err
		}

		// Compute next exponential backoff
		nextBackoff := time.Duration(float64(currentBackoff) * c.config.BackoffFactor)
		if nextBackoff > c.config.MaxBackoff {
			nextBackoff = c.config.MaxBackoff
		}
		currentBackoff = nextBackoff
	}

	return nil, lastErr
}

func (c *ResilientClient) executeAttempt(ctx context.Context, fn RPCCallFunc) (any, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, c.config.PerAttemptTimeout)
	defer cancel()

	return fn(attemptCtx)
}

func (c *ResilientClient) sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
