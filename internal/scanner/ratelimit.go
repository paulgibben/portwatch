package scanner

import (
	"sync"
	"time"
)

// RateLimiter controls the rate of port scan operations to avoid
// overwhelming the system or triggering network defenses.
type RateLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	last     time.Time
	tokens   int
	max      int
}

// RateLimitConfig holds configuration for the rate limiter.
type RateLimitConfig struct {
	// MaxConcurrent is the maximum number of concurrent scan tokens.
	MaxConcurrent int `json:"max_concurrent"`
	// IntervalMs is the minimum interval in milliseconds between scans.
	IntervalMs int `json:"interval_ms"`
}

// NewRateLimiter creates a RateLimiter from the given config.
// If cfg is nil or values are zero, sensible defaults are applied.
func NewRateLimiter(cfg *RateLimitConfig) *RateLimiter {
	max := 10
	interval := 50 * time.Millisecond

	if cfg != nil {
		if cfg.MaxConcurrent > 0 {
			max = cfg.MaxConcurrent
		}
		if cfg.IntervalMs > 0 {
			interval = time.Duration(cfg.IntervalMs) * time.Millisecond
		}
	}

	return &RateLimiter{
		interval: interval,
		max:      max,
		tokens:   max,
	}
}

// Acquire blocks until a scan token is available and the interval has elapsed.
func (r *RateLimiter) Acquire() {
	for {
		r.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(r.last)
		if elapsed >= r.interval && r.tokens > 0 {
			r.tokens--
			r.last = now
			r.mu.Unlock()
			return
		}
		r.mu.Unlock()
		time.Sleep(r.interval / 4)
	}
}

// Release returns a token to the pool.
func (r *RateLimiter) Release() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.tokens < r.max {
		r.tokens++
	}
}

// Available returns the current number of available tokens.
func (r *RateLimiter) Available() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tokens
}
