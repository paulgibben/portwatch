package scanner

import (
	"sync"
	"time"
)

// ScanScheduleTracker tracks scan execution times and detects missed scans.
type ScanScheduleTracker struct {
	mu        sync.Mutex
	cfg       ScanScheduleConfig
	lastScan  time.Time
	missed    int
	scansRun  int
}

// NewScanScheduleTracker creates a new tracker with the given config.
func NewScanScheduleTracker(cfg ScanScheduleConfig) *ScanScheduleTracker {
	return &ScanScheduleTracker{cfg: cfg}
}

// Record marks a scan as having occurred at the given time.
func (t *ScanScheduleTracker) Record(at time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.lastScan.IsZero() {
		expected := t.lastScan.Add(t.cfg.Interval)
		if at.After(expected.Add(t.cfg.Interval / 2)) {
			t.missed++
		}
	}
	t.lastScan = at
	t.scansRun++
}

// MissedCount returns the number of detected missed scans.
func (t *ScanScheduleTracker) MissedCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.missed
}

// ScansRun returns the total number of scans recorded.
func (t *ScanScheduleTracker) ScansRun() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.scansRun
}

// LastScan returns the time of the most recent recorded scan.
func (t *ScanScheduleTracker) LastScan() time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastScan
}

// ExceededMissedThreshold returns true if missed scans exceed MaxMissed.
func (t *ScanScheduleTracker) ExceededMissedThreshold() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.missed >= t.cfg.MaxMissed
}
