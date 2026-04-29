package scanner

import (
	"sync"
	"time"
)

// SeverityEvent records a severity evaluation for a port change.
type SeverityEvent struct {
	Port      int           `json:"port"`
	Protocol  string        `json:"protocol"`
	Severity  SeverityLevel `json:"severity"`
	ChangeType string       `json:"change_type"` // "added" or "removed"
	Timestamp time.Time     `json:"timestamp"`
}

// PortSeverityTracker records severity events derived from diffs.
type PortSeverityTracker struct {
	mu     sync.RWMutex
	store  *PortSeverityStore
	events []SeverityEvent
}

// NewPortSeverityTracker creates a tracker backed by the given store.
func NewPortSeverityTracker(store *PortSeverityStore) *PortSeverityTracker {
	return &PortSeverityTracker{store: store}
}

// Track evaluates and records severity for each port in a Diff.
func (t *PortSeverityTracker) Track(d Diff) {
	now := time.Now().UTC()

	t.mu.Lock()
	defer t.mu.Unlock()

	for _, p := range d.Added {
		t.events = append(t.events, SeverityEvent{
			Port:       p.Port,
			Protocol:   p.Protocol,
			Severity:   t.store.Evaluate(p.Port, p.Protocol),
			ChangeType: "added",
			Timestamp:  now,
		})
	}
	for _, p := range d.Removed {
		t.events = append(t.events, SeverityEvent{
			Port:       p.Port,
			Protocol:   p.Protocol,
			Severity:   t.store.Evaluate(p.Port, p.Protocol),
			ChangeType: "removed",
			Timestamp:  now,
		})
	}
}

// Events returns a copy of all recorded severity events.
func (t *PortSeverityTracker) Events() []SeverityEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]SeverityEvent, len(t.events))
	copy(out, t.events)
	return out
}

// EventsAbove returns events at or above the given severity level.
func (t *PortSeverityTracker) EventsAbove(min SeverityLevel) []SeverityEvent {
	order := map[SeverityLevel]int{
		SeverityLow: 0, SeverityMedium: 1, SeverityHigh: 2, SeverityCritical: 3,
	}
	t.mu.RLock()
	defer t.mu.RUnlock()
	var out []SeverityEvent
	for _, e := range t.events {
		if order[e.Severity] >= order[min] {
			out = append(out, e)
		}
	}
	return out
}
