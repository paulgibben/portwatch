package scanner

import (
	"testing"
	"time"
)

func TestNewThrottle_Defaults(t *testing.T) {
	th := NewThrottle(ThrottleConfig{})
	if th.cfg.MinInterval != 5*time.Second {
		t.Errorf("expected default MinInterval 5s, got %v", th.cfg.MinInterval)
	}
	if th.cfg.BurstAllowance != 1 {
		t.Errorf("expected default BurstAllowance 1, got %d", th.cfg.BurstAllowance)
	}
}

func TestNewThrottle_CustomConfig(t *testing.T) {
	cfg := ThrottleConfig{
		MinInterval:    100 * time.Millisecond,
		BurstAllowance: 3,
	}
	th := NewThrottle(cfg)
	if th.cfg.MinInterval != 100*time.Millisecond {
		t.Errorf("expected 100ms, got %v", th.cfg.MinInterval)
	}
	if th.cfg.BurstAllowance != 3 {
		t.Errorf("expected burst 3, got %d", th.cfg.BurstAllowance)
	}
}

func TestThrottle_FirstWait_ImmediatelyReturns(t *testing.T) {
	th := NewThrottle(ThrottleConfig{MinInterval: 10 * time.Second, BurstAllowance: 1})
	start := time.Now()
	th.Wait()
	if time.Since(start) > 50*time.Millisecond {
		t.Error("first Wait should return immediately")
	}
	if th.LastScan().IsZero() {
		t.Error("LastScan should be set after Wait")
	}
}

func TestThrottle_EnforcesInterval(t *testing.T) {
	cfg := ThrottleConfig{
		MinInterval:    80 * time.Millisecond,
		BurstAllowance: 1,
	}
	th := NewThrottle(cfg)
	th.Wait() // first call — immediate

	start := time.Now()
	th.Wait() // second call — should block ~80ms
	elapsed := time.Since(start)

	if elapsed < 60*time.Millisecond {
		t.Errorf("expected throttle delay, got %v", elapsed)
	}
}

func TestThrottle_BurstAllowance(t *testing.T) {
	cfg := ThrottleConfig{
		MinInterval:    500 * time.Millisecond,
		BurstAllowance: 3,
	}
	th := NewThrottle(cfg)

	// Three calls should all complete quickly due to burst allowance.
	start := time.Now()
	th.Wait()
	th.Wait()
	th.Wait()
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("burst calls should be fast, took %v", elapsed)
	}
}

func TestDefaultThrottleConfig(t *testing.T) {
	cfg := DefaultThrottleConfig()
	if cfg.MinInterval <= 0 {
		t.Error("default MinInterval must be positive")
	}
	if cfg.BurstAllowance <= 0 {
		t.Error("default BurstAllowance must be positive")
	}
}
