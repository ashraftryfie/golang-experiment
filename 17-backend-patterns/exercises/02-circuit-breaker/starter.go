package circuitbreaker

import (
	"errors"
	"time"
)

var (
	ErrNotImplemented = errors.New("exercise not implemented yet")
	ErrCircuitOpen     = errors.New("circuit breaker is open; failing fast")
)

type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

type Config struct {
	MaxFailures      int
	Timeout          time.Duration
	SuccessThreshold int
}

type CircuitBreaker struct{}

func NewCircuitBreaker(cfg Config) *CircuitBreaker {
	return nil
}

func (cb *CircuitBreaker) State() CircuitState {
	return StateClosed
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	return ErrNotImplemented
}
