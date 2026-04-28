package scanner

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts 3, got %d", cfg.MaxAttempts)
	}
	if cfg.Delay != 100*time.Millisecond {
		t.Errorf("expected Delay 100ms, got %v", cfg.Delay)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier 2.0, got %v", cfg.Multiplier)
	}
}

func TestNewRetryer_Defaults(t *testing.T) {
	r := NewRetryer(RetryConfig{})
	if r.cfg.MaxAttempts != 3 {
		t.Errorf("expected default MaxAttempts 3, got %d", r.cfg.MaxAttempts)
	}
	if r.cfg.Delay != 100*time.Millisecond {
		t.Errorf("expected default Delay 100ms, got %v", r.cfg.Delay)
	}
}

func TestNewRetryer_CustomConfig(t *testing.T) {
	r := NewRetryer(RetryConfig{MaxAttempts: 5, Delay: 10 * time.Millisecond, Multiplier: 1.5})
	if r.cfg.MaxAttempts != 5 {
		t.Errorf("expected MaxAttempts 5, got %d", r.cfg.MaxAttempts)
	}
}

func TestRetryer_Do_SucceedsFirstAttempt(t *testing.T) {
	r := NewRetryer(RetryConfig{MaxAttempts: 3, Delay: 1 * time.Millisecond, Multiplier: 1.0})
	calls := 0
	err := r.Do(func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryer_Do_RetriesOnFailure(t *testing.T) {
	r := NewRetryer(RetryConfig{MaxAttempts: 3, Delay: 1 * time.Millisecond, Multiplier: 1.0})
	calls := 0
	err := r.Do(func() error {
		calls++
		if calls < 3 {
			return errors.New("transient error")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success on 3rd attempt, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetryer_Do_ExhaustsAttempts(t *testing.T) {
	r := NewRetryer(RetryConfig{MaxAttempts: 3, Delay: 1 * time.Millisecond, Multiplier: 1.0})
	calls := 0
	err := r.Do(func() error {
		calls++
		return errors.New("permanent error")
	})
	if err == nil {
		t.Fatal("expected error after exhausting attempts")
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}
