package scanner

import (
	"testing"
)

func makeSeverityDiff() Diff {
	return Diff{
		Added: []Port{
			{Port: 22, Protocol: "tcp"},
			{Port: 8080, Protocol: "tcp"},
		},
		Removed: []Port{
			{Port: 3306, Protocol: "tcp"},
		},
	}
}

func TestSeverityTracker_Track_RecordsEvents(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	tracker := NewPortSeverityTracker(store)
	tracker.Track(makeSeverityDiff())

	events := tracker.Events()
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
}

func TestSeverityTracker_Track_CorrectChangeType(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	tracker := NewPortSeverityTracker(store)
	tracker.Track(makeSeverityDiff())

	events := tracker.Events()
	addedCount, removedCount := 0, 0
	for _, e := range events {
		switch e.ChangeType {
		case "added":
			addedCount++
		case "removed":
			removedCount++
		}
	}
	if addedCount != 2 {
		t.Errorf("expected 2 added events, got %d", addedCount)
	}
	if removedCount != 1 {
		t.Errorf("expected 1 removed event, got %d", removedCount)
	}
}

func TestSeverityTracker_Track_AssignsSeverity(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	tracker := NewPortSeverityTracker(store)
	tracker.Track(Diff{
		Added: []Port{{Port: 22, Protocol: "tcp"}},
	})
	events := tracker.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event")
	}
	if events[0].Severity != SeverityHigh {
		t.Errorf("expected high severity for port 22, got %s", events[0].Severity)
	}
}

func TestSeverityTracker_EventsAbove_FiltersCorrectly(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	tracker := NewPortSeverityTracker(store)
	tracker.Track(makeSeverityDiff())

	high := tracker.EventsAbove(SeverityHigh)
	// port 22 = high, port 3306 = critical; port 8080 = low (default)
	if len(high) != 2 {
		t.Errorf("expected 2 events at or above high, got %d", len(high))
	}
}

func TestSeverityTracker_Events_ReturnsCopy(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	tracker := NewPortSeverityTracker(store)
	tracker.Track(makeSeverityDiff())

	e1 := tracker.Events()
	e1[0].Severity = SeverityLow
	e2 := tracker.Events()
	if e2[0].Severity == SeverityLow {
		t.Error("Events() should return a copy, not a reference")
	}
}
