package scanner

import (
	"errors"
	"time"
)

// RetryConfig holds configuration for retry behaviour during port scanning.
type RetryConfig struct {
	MaxAttempts int           `json:"max_attempts"`
	Delay       time.Duration `json:"delay"`
	Multiplier  float64       `json:"multiplier"`
}

// DefaultRetryConfig returns a RetryConfig with sensible defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		Delay:       100 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Retryer executes an operation with exponential back-off retry logic.
type Retryer struct {
	cfg RetryConfig
}

// NewRetryer creates a Retryer with the given config.
// Zero-value fields fall back to DefaultRetryConfig values.
func NewRetryer(cfg RetryConfig) *Retryer {
	def := DefaultRetryConfig()
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = def.MaxAttempts
	}
	if cfg.Delay <= 0 {
		cfg.Delay = def.Delay
	}
	if cfg.Multiplier <= 0 {
		cfg.Multiplier = def.Multiplier
	}
	return &Retryer{cfg: cfg}
}

// Do runs fn up to MaxAttempts times, backing off between attempts.
// It returns the last error if all attempts fail.
func (r *Retryer) Do(fn func() error) error {
	var err error
	delay := r.cfg.Delay
	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		if attempt < r.cfg.MaxAttempts {
			time.Sleep(delay)
			delay = time.Duration(float64(delay) * r.cfg.Multiplier)
		}
	}
	return errors.New("all attempts failed: " + err.Error())
}
