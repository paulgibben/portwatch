package scanner

import (
	"sync"
	"time"
)

// TagEvent records a tag being added or removed from a port.
type TagEvent struct {
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Tag       string    `json:"tag"`
	Action    string    `json:"action"` // "added" or "removed"
	Timestamp time.Time `json:"timestamp"`
}

// PortTagTracker wraps a PortTagStore and records tag change events.
type PortTagTracker struct {
	mu     sync.Mutex
	store  *PortTagStore
	events []TagEvent
}

// NewPortTagTracker creates a tracker backed by the given store.
func NewPortTagTracker(store *PortTagStore) *PortTagTracker {
	return &PortTagTracker{store: store}
}

// AddTag adds a tag and records the event.
func (t *PortTagTracker) AddTag(port int, proto, tag string) {
	before := t.store.GetTags(port, proto)
	t.store.AddTag(port, proto, tag)
	after := t.store.GetTags(port, proto)
	if len(after) > len(before) {
		t.mu.Lock()
		t.events = append(t.events, TagEvent{
			Port:      port,
			Proto:     proto,
			Tag:       tag,
			Action:    "added",
			Timestamp: time.Now(),
		})
		t.mu.Unlock()
	}
}

// RemoveTag removes a tag and records the event.
func (t *PortTagTracker) RemoveTag(port int, proto, tag string) {
	before := t.store.GetTags(port, proto)
	t.store.RemoveTag(port, proto, tag)
	after := t.store.GetTags(port, proto)
	if len(after) < len(before) {
		t.mu.Lock()
		t.events = append(t.events, TagEvent{
			Port:      port,
			Proto:     proto,
			Tag:       tag,
			Action:    "removed",
			Timestamp: time.Now(),
		})
		t.mu.Unlock()
	}
}

// Events returns a copy of all recorded tag events.
func (t *PortTagTracker) Events() []TagEvent {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]TagEvent, len(t.events))
	copy(out, t.events)
	return out
}

// Store returns the underlying PortTagStore.
func (t *PortTagTracker) Store() *PortTagStore {
	return t.store
}
