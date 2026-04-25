package reporter

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestScheduled_RunsAndStops(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	calls := 0
	produce := func() (Report, error) {
		calls++
		return Report{
			Timestamp: time.Now(),
			Ports:     []scanner.Port{{Number: 22, Protocol: "tcp"}},
		}, nil
	}
	sr := NewScheduled(r, 20*time.Millisecond, produce)
	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)
	defer cancel()
	err := sr.Run(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if calls < 2 {
		t.Errorf("expected at least 2 produce calls, got %d", calls)
	}
}

func TestScheduled_SkipsOnProduceError(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	produce := func() (Report, error) {
		return Report{}, errors.New("scan failed")
	}
	sr := NewScheduled(r, 15*time.Millisecond, produce)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_ = sr.Run(ctx)
	// Output buffer should be empty since produce always errors.
	if buf.Len() != 0 {
		t.Errorf("expected empty output on produce error, got: %s", buf.String())
	}
}
