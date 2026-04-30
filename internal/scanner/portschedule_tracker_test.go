package scanner

import (
	"testing"
	"time"
)

func makeScanScheduleTracker() *ScanScheduleTracker {
	cfg := ScanScheduleConfig{
		Interval:  10 * time.Second,
		MaxMissed: 2,
		Enabled:   true,
	}
	return NewScanScheduleTracker(cfg)
}

func TestScanScheduleTracker_RecordsScans(t *testing.T) {
	tr := makeScanScheduleTracker()
	now := time.Now()
	tr.Record(now)
	tr.Record(now.Add(10 * time.Second))
	if tr.ScansRun() != 2 {
		t.Errorf("expected 2 scans, got %d", tr.ScansRun())
	}
}

func TestScanScheduleTracker_LastScan(t *testing.T) {
	tr := makeScanScheduleTracker()
	now := time.Now()
	tr.Record(now)
	if !tr.LastScan().Equal(now) {
		t.Errorf("expected LastScan=%v, got %v", now, tr.LastScan())
	}
}

func TestScanScheduleTracker_DetectsMissedScan(t *testing.T) {
	tr := makeScanScheduleTracker()
	now := time.Now()
	tr.Record(now)
	// Simulate a scan that arrived 25s late (interval=10s, so >15s gap triggers missed)
	tr.Record(now.Add(25 * time.Second))
	if tr.MissedCount() != 1 {
		t.Errorf("expected 1 missed scan, got %d", tr.MissedCount())
	}
}

func TestScanScheduleTracker_NoMissedOnTime(t *testing.T) {
	tr := makeScanScheduleTracker()
	now := time.Now()
	tr.Record(now)
	tr.Record(now.Add(10 * time.Second))
	tr.Record(now.Add(20 * time.Second))
	if tr.MissedCount() != 0 {
		t.Errorf("expected 0 missed scans, got %d", tr.MissedCount())
	}
}

func TestScanScheduleTracker_ExceededThreshold(t *testing.T) {
	tr := makeScanScheduleTracker()
	now := time.Now()
	tr.Record(now)
	tr.Record(now.Add(30 * time.Second)) // missed 1
	tr.Record(now.Add(70 * time.Second)) // missed 2
	if !tr.ExceededMissedThreshold() {
		t.Error("expected threshold to be exceeded")
	}
}

func TestScanScheduleTracker_EmptyLastScan(t *testing.T) {
	tr := makeScanScheduleTracker()
	if !tr.LastScan().IsZero() {
		t.Error("expected zero LastScan before any record")
	}
}
