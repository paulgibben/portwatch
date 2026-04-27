package scanner

import (
	"sync"
	"time"
)

// ThrottleConfig holds configuration for scan throttling.
type ThrottleConfig struct {
	// MinInterval is the minimum time between successive scans.
	MinInterval time.Duration `json:"min_interval"`
	// BurstAllowance is how many scans can run back-to-back before throttling kicks in.
	BurstAllowance int `json:"burst_allowance"`
}

// DefaultThrottleConfig returns a ThrottleConfig with sensible defaults.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		MinInterval:    5 * time.Second,
		BurstAllowance: 1,
	}
}

// Throttle enforces a minimum interval between scan cycles.
type Throttle struct {
	mu       sync.Mutex
	cfg      ThrottleConfig
	lastScan time.Time
	burst    int
}

// NewThrottle creates a new Throttle with the given config.
// Zero-value fields are replaced with defaults.
func NewThrottle(cfg ThrottleConfig) *Throttle {
	def := DefaultThrottleConfig()
	if cfg.MinInterval <= 0 {
		cfg.MinInterval = def.MinInterval
	}
	if cfg.BurstAllowance <= 0 {
		cfg.BurstAllowance = def.BurstAllowance
	}
	return &Throttle{cfg: cfg}
}

// Wait blocks until a scan is allowed to proceed, then records the scan time.
func (t *Throttle) Wait() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.lastScan.IsZero() {
		t.lastScan = time.Now()
		t.burst = 1
		return
	}

	elapsed := time.Since(t.lastScan)
	if elapsed >= t.cfg.MinInterval || t.burst < t.cfg.BurstAllowance {
		if elapsed < t.cfg.MinInterval {
			t.burst++
		} else {
			t.burst = 1
		}
		t.lastScan = time.Now()
		return
	}

	sleepFor := t.cfg.MinInterval - elapsed
	t.mu.Unlock()
	time.Sleep(sleepFor)
	t.mu.Lock()
	t.lastScan = time.Now()
	t.burst = 1
}

// LastScan returns the timestamp of the most recent scan.
func (t *Throttle) LastScan() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastScan
}
