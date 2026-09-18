package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open; failing fast")
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

type CircuitBreaker struct {
	config Config

	state        CircuitState
	failures     int
	successes    int
	lastOpenedAt time.Time

	mu sync.Mutex
}

func NewCircuitBreaker(cfg Config) *CircuitBreaker {
	if cfg.MaxFailures <= 0 {
		cfg.MaxFailures = 3
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 1 * time.Second
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = 1
	}

	return &CircuitBreaker{
		config: cfg,
		state:  StateClosed,
	}
}

func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.checkStateTransition()
	return cb.state
}

func (cb *CircuitBreaker) checkStateTransition() {
	if cb.state == StateOpen {
		if time.Since(cb.lastOpenedAt) >= cb.config.Timeout {
			cb.state = StateHalfOpen
			cb.successes = 0
		}
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	cb.checkStateTransition()

	if cb.state == StateOpen {
		cb.mu.Unlock()
		return ErrCircuitOpen
	}
	cb.mu.Unlock()

	// Execute operation outside mutex
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
		return err
	}

	cb.onSuccess()
	return nil
}

func (cb *CircuitBreaker) onFailure() {
	switch cb.state {
	case StateClosed:
		cb.failures++
		if cb.failures >= cb.config.MaxFailures {
			cb.state = StateOpen
			cb.lastOpenedAt = time.Now()
		}
	case StateHalfOpen:
		// Any failure in half-open trips back to open
		cb.state = StateOpen
		cb.lastOpenedAt = time.Now()
	}
}

func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateClosed:
		cb.failures = 0
	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.SuccessThreshold {
			cb.state = StateClosed
			cb.failures = 0
			cb.successes = 0
		}
	}
}
