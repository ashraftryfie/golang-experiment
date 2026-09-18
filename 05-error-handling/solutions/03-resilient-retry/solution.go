package retry

import (
	"errors"
	"time"
)

var ErrMaxAttemptsReached = errors.New("maximum retry attempts exceeded")

type TransientError interface {
	IsTransient() bool
}

type NetworkError struct {
	Msg string
}

func (e NetworkError) Error() string {
	return e.Msg
}

func (e NetworkError) IsTransient() bool {
	return true
}

type AuthError struct {
	Msg string
}

func (e AuthError) Error() string {
	return e.Msg
}

func (e AuthError) IsTransient() bool {
	return false
}

// Retry executes fn up to maxAttempts, aborting immediately on permanent errors.
func Retry(maxAttempts int, initialDelay time.Duration, fn func() error) error {
	delay := initialDelay
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err

		// Check if error is explicitly permanent
		var transientErr TransientError
		if errors.As(err, &transientErr) && !transientErr.IsTransient() {
			return err // Abort immediately!
		}

		if attempt < maxAttempts {
			time.Sleep(delay)
			delay *= 2 // Exponential backoff
		}
	}
	return lastErr
}
