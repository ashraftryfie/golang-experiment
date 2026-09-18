package retry

import (
	"errors"
	"time"
)

var ErrMaxAttemptsReached = errors.New("maximum retry attempts exceeded")
var ErrNotImplemented = errors.New("TODO: implement")

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
	// TODO: Implement retry loop with exponential backoff and TransientError check
	return ErrNotImplemented
}
