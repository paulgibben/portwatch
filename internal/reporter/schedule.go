package reporter

import (
	"context"
	"time"
)

// ScheduledReporter periodically invokes a callback to produce and write reports.
type ScheduledReporter struct {
	reporter *Reporter
	interval time.Duration
	produce  func() (Report, error)
}

// NewScheduled creates a ScheduledReporter that calls produce every interval
// and writes the result via r.
func NewScheduled(r *Reporter, interval time.Duration, produce func() (Report, error)) *ScheduledReporter {
	return &ScheduledReporter{
		reporter: r,
		interval: interval,
		produce:  produce,
	}
}

// Run starts the scheduled reporting loop. It blocks until ctx is cancelled.
func (s *ScheduledReporter) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			rep, err := s.produce()
			if err != nil {
				// non-fatal: continue running
				continue
			}
			_ = s.reporter.Write(rep)
		}
	}
}
