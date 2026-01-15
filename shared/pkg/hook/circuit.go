// Package hook provides circuit breaker for hook resilience.
// circuit.go: Circuit breaker types and creation.
// DACE: Prevents cascade failures when hooks fail repeatedly.
package hook

import (
	"sync"
	"time"
)

// CircuitState represents the circuit breaker state.
type CircuitState int

const (
	StateClosed   CircuitState = iota // Normal operation
	StateOpen                         // Failing, bypass hooks
	StateHalfOpen                     // Testing if recovered
)

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	mu              sync.RWMutex
	state           CircuitState
	failures        int
	lastFailure     time.Time
	threshold       int           // failures before opening
	timeout         time.Duration // time before half-open
	successRequired int           // successes to close
	successes       int
}

// DefaultCircuit is the global circuit breaker for hooks.
var DefaultCircuit = NewCircuitBreaker(3, 30*time.Second, 1)

// NewCircuitBreaker creates a circuit breaker.
func NewCircuitBreaker(threshold int, timeout time.Duration, successRequired int) *CircuitBreaker {
	return &CircuitBreaker{
		state:           StateClosed,
		threshold:       threshold,
		timeout:         timeout,
		successRequired: successRequired,
	}
}

// State returns the current circuit state.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// IsOpen returns true if circuit is open (failing).
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.State() == StateOpen
}
