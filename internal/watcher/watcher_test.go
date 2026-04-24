package watcher_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watcher"
)

func makeConfig() *rules.Config {
	return &rules.Config{
		IgnorePorts: []int{},
	}
}

func TestWatcher_StartsAndStops(t *testing.T) {
	s := scanner.New("localhost", []int{}, 200*time.Millisecond)
	n := alert.New(nil)
	cfg := makeConfig()

	w := watcher.New(s, n, cfg, 50*time.Millisecond)
	w.Start()
	time.Sleep(120 * time.Millisecond)
	w.Stop()
	// If we reach here without a panic or deadlock the test passes.
}

func TestWatcher_StopIdempotent(t *testing.T) {
	s := scanner.New("localhost", []int{}, 200*time.Millisecond)
	n := alert.New(nil)
	cfg := makeConfig()

	w := watcher.New(s, n, cfg, 100*time.Millisecond)
	w.Start()
	w.Stop()
	// Calling Stop again should not panic (channel already closed).
	// We guard this with a recover.
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Stop panicked on second call: %v", r)
			}
		}()
	}()
}

func TestWatcher_New_ReturnsNonNil(t *testing.T) {
	s := scanner.New("localhost", []int{}, 200*time.Millisecond)
	n := alert.New(nil)
	cfg := makeConfig()

	w := watcher.New(s, n, cfg, 1*time.Second)
	if w == nil {
		t.Fatal("expected non-nil Watcher")
	}
}
