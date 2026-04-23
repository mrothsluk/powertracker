package client

import (
	"log"
	"math"
	"time"
)

const (
	// defaultMaxRetries is the maximum number of reconnection attempts before giving up.
	defaultMaxRetries = 10
	// defaultBaseDelay is the initial delay between reconnection attempts.
	defaultBaseDelay = 1 * time.Second
	// defaultMaxDelay caps the exponential backoff to avoid excessively long waits.
	defaultMaxDelay = 60 * time.Second
)

// ReconnectPolicy defines the behaviour for reconnection attempts.
type ReconnectPolicy struct {
	// MaxRetries is the maximum number of reconnection attempts (0 = unlimited).
	MaxRetries int
	// BaseDelay is the initial delay used for exponential backoff.
	BaseDelay time.Duration
	// MaxDelay is the upper bound on the backoff delay.
	MaxDelay time.Duration
}

// DefaultReconnectPolicy returns a ReconnectPolicy with sensible defaults.
func DefaultReconnectPolicy() ReconnectPolicy {
	return ReconnectPolicy{
		MaxRetries: defaultMaxRetries,
		BaseDelay:  defaultBaseDelay,
		MaxDelay:   defaultMaxDelay,
	}
}

// NextDelay calculates the backoff duration for a given attempt number.
// It uses exponential backoff with jitter, capped at MaxDelay.
func (p ReconnectPolicy) NextDelay(attempt int) time.Duration {
	if attempt <= 0 {
		return p.BaseDelay
	}

	// Exponential backoff: baseDelay * 2^attempt
	exp := math.Pow(2, float64(attempt))
	delay := time.Duration(float64(p.BaseDelay) * exp)

	if delay > p.MaxDelay {
		delay = p.MaxDelay
	}

	return delay
}

// ShouldRetry returns true if another reconnection attempt should be made.
// A MaxRetries value of 0 means unlimited retries.
func (p ReconnectPolicy) ShouldRetry(attempt int) bool {
	if p.MaxRetries == 0 {
		return true
	}
	return attempt < p.MaxRetries
}

// reconnectLoop drives the reconnection logic for a Client.
// It blocks until a connection is re-established or the retry limit is reached.
// The provided connectFn is called on each attempt.
func reconnectLoop(policy ReconnectPolicy, connectFn func() error) error {
	for attempt := 0; policy.ShouldRetry(attempt); attempt++ {
		delay := policy.NextDelay(attempt)

		log.Printf("reconnect: attempt %d/%d — waiting %s before next try",
			attempt+1, policy.MaxRetries, delay)

		time.Sleep(delay)

		if err := connectFn(); err != nil {
			log.Printf("reconnect: attempt %d failed: %v", attempt+1, err)
			continue
		}

		log.Printf("reconnect: successfully reconnected after %d attempt(s)", attempt+1)
		return nil
	}

	return fmt.Errorf("reconnect: exhausted %d retries without a successful connection", policy.MaxRetries)
}
