package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// PortState represents the observed state of a single port over time.
type PortState struct {
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Open      bool      `json:"open"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	SeenCount int       `json:"seen_count"`
}

// PortStateStore tracks state transitions for monitored ports.
type PortStateStore struct {
	Entries map[string]*PortState `json:"entries"`
}

// NewPortStateStore creates an empty PortStateStore.
func NewPortStateStore() *PortStateStore {
	return &PortStateStore{
		Entries: make(map[string]*PortState),
	}
}

func stateKey(port int, protocol string) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}

// Observe records that a port was seen open or closed at the given time.
func (s *PortStateStore) Observe(port int, protocol string, open bool, at time.Time) {
	key := stateKey(port, protocol)
	entry, exists := s.Entries[key]
	if !exists {
		entry = &PortState{
			Port:      port,
			Protocol:  protocol,
			FirstSeen: at,
		}
		s.Entries[key] = entry
	}
	entry.Open = open
	entry.LastSeen = at
	entry.SeenCount++
}

// Get returns the state for a port, or nil if not tracked.
func (s *PortStateStore) Get(port int, protocol string) *PortState {
	return s.Entries[stateKey(port, protocol)]
}

// SavePortState persists the store to a JSON file.
func SavePortState(path string, store *PortStateStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal port state: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPortState loads a PortStateStore from a JSON file.
func LoadPortState(path string) (*PortStateStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewPortStateStore(), nil
		}
		return nil, fmt.Errorf("read port state: %w", err)
	}
	var store PortStateStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("unmarshal port state: %w", err)
	}
	if store.Entries == nil {
		store.Entries = make(map[string]*PortState)
	}
	return &store, nil
}
