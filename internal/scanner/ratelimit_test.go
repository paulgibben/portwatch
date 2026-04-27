package scanner

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRateLimiter_Defaults(t *testing.T) {
	rl := NewRateLimiter(nil)
	if rl.max != 10 {
		t.Errorf("expected max=10, got %d", rl.max)
	}
	if rl.interval != 50*time.Millisecond {
		t.Errorf("expected interval=50ms, got %v", rl.interval)
	}
}

func TestNewRateLimiter_CustomConfig(t *testing.T) {
	cfg := &RateLimitConfig{MaxConcurrent: 5, IntervalMs: 100}
	rl := NewRateLimiter(cfg)
	if rl.max != 5 {
		t.Errorf("expected max=5, got %d", rl.max)
	}
	if rl.interval != 100*time.Millisecond {
		t.Errorf("expected interval=100ms, got %v", rl.interval)
	}
}

func TestRateLimiter_AcquireRelease(t *testing.T) {
	rl := NewRateLimiter(&RateLimitConfig{MaxConcurrent: 3, IntervalMs: 1})

	if got := rl.Available(); got != 3 {
		t.Fatalf("expected 3 tokens, got %d", got)
	}

	rl.Acquire()
	if got := rl.Available(); got != 2 {
		t.Errorf("expected 2 tokens after acquire, got %d", got)
	}

	rl.Release()
	if got := rl.Available(); got != 3 {
		t.Errorf("expected 3 tokens after release, got %d", got)
	}
}

func TestRateLimiter_ReleaseDoesNotExceedMax(t *testing.T) {
	rl := NewRateLimiter(&RateLimitConfig{MaxConcurrent: 2, IntervalMs: 1})
	rl.Release()
	rl.Release()
	if got := rl.Available(); got != 2 {
		t.Errorf("expected tokens capped at max=2, got %d", got)
	}
}

func TestRateLimiter_ConcurrentAcquire(t *testing.T) {
	rl := NewRateLimiter(&RateLimitConfig{MaxConcurrent: 5, IntervalMs: 1})

	var count int64
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rl.Acquire()
			atomic.AddInt64(&count, 1)
			time.Sleep(2 * time.Millisecond)
			rl.Release()
		}()
	}

	wg.Wait()
	if count != 5 {
		t.Errorf("expected 5 acquisitions, got %d", count)
	}
}
