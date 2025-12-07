package resilience

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrCircuitOpen is returned when the circuit breaker is open
	ErrCircuitOpen = errors.New("circuit breaker is open")
	// ErrTimeout is returned when operation times out
	ErrTimeout = errors.New("operation timed out")
)

// CircuitState represents the state of the circuit breaker
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateHalfOpen
	StateOpen
)

func (s CircuitState) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	name            string
	maxFailures     int
	timeout         time.Duration
	halfOpenMaxCalls int

	mu              sync.RWMutex
	state           CircuitState
	failureCount    int
	successCount    int
	lastFailureTime time.Time
	halfOpenCalls   int
}

// CircuitBreakerConfig configures the circuit breaker
type CircuitBreakerConfig struct {
	Name             string
	MaxFailures      int           // Failures before opening
	Timeout          time.Duration // Time before trying half-open
	HalfOpenMaxCalls int           // Max calls in half-open state
}

// DefaultCircuitConfig returns sensible defaults
func DefaultCircuitConfig(name string) CircuitBreakerConfig {
	return CircuitBreakerConfig{
		Name:             name,
		MaxFailures:      5,
		Timeout:          30 * time.Second,
		HalfOpenMaxCalls: 3,
	}
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		name:             cfg.Name,
		maxFailures:      cfg.MaxFailures,
		timeout:          cfg.Timeout,
		halfOpenMaxCalls: cfg.HalfOpenMaxCalls,
		state:            StateClosed,
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.allowRequest() {
		return ErrCircuitOpen
	}

	err := fn()

	cb.recordResult(err)
	return err
}

// allowRequest checks if a request should be allowed
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true

	case StateOpen:
		// Check if timeout has elapsed
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.halfOpenCalls = 0
			cb.successCount = 0
			cb.failureCount = 0
			return true
		}
		return false

	case StateHalfOpen:
		if cb.halfOpenCalls < cb.halfOpenMaxCalls {
			cb.halfOpenCalls++
			return true
		}
		return false

	default:
		return true
	}
}

// recordResult records the result of an operation
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		switch cb.state {
		case StateClosed:
			if cb.failureCount >= cb.maxFailures {
				cb.state = StateOpen
			}
		case StateHalfOpen:
			cb.state = StateOpen
		}
	} else {
		cb.successCount++

		switch cb.state {
		case StateHalfOpen:
			if cb.successCount >= cb.halfOpenMaxCalls {
				cb.state = StateClosed
				cb.failureCount = 0
				cb.successCount = 0
			}
		case StateClosed:
			// Reset failure count on success in closed state
			if cb.failureCount > 0 {
				cb.failureCount = 0
			}
		}
	}
}

// State returns the current state
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Stats returns circuit breaker statistics
func (cb *CircuitBreaker) Stats() CircuitBreakerStats {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return CircuitBreakerStats{
		Name:         cb.name,
		State:        cb.state.String(),
		FailureCount: cb.failureCount,
		SuccessCount: cb.successCount,
	}
}

// CircuitBreakerStats holds circuit breaker metrics
type CircuitBreakerStats struct {
	Name         string
	State        string
	FailureCount int
	SuccessCount int
}

// Reset forces the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.state = StateClosed
	cb.failureCount = 0
	cb.successCount = 0
}

// Retry executes a function with automatic retries
type Retry struct {
	MaxAttempts int
	Delay       time.Duration
	MaxDelay    time.Duration
	Multiplier  float64
}

// DefaultRetry returns sensible retry defaults
func DefaultRetry() *Retry {
	return &Retry{
		MaxAttempts: 3,
		Delay:       100 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  2.0,
	}
}

// Execute runs the function with retries
func (r *Retry) Execute(fn func() error) error {
	var lastErr error
	delay := r.Delay

	for attempt := 1; attempt <= r.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt < r.MaxAttempts {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * r.Multiplier)
			if delay > r.MaxDelay {
				delay = r.MaxDelay
			}
		}
	}

	return lastErr
}
