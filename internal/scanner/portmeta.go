package scanner

import (
	"encoding/json"
	"os"
	"time"
)

// PortMeta holds metadata about a known port, including when it was first
// and last seen, and an optional human-readable label.
type PortMeta struct {
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Label     string    `json:"label,omitempty"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
	SeenCount int       `json:"seen_count"`
}

// PortMetaStore maps port numbers to their metadata.
type PortMetaStore struct {
	Entries map[int]*PortMeta `json:"entries"`
}

// NewPortMetaStore creates an empty PortMetaStore.
func NewPortMetaStore() *PortMetaStore {
	return &PortMetaStore{
		Entries: make(map[int]*PortMeta),
	}
}

// Observe records a port as seen at the given time, creating or updating its entry.
func (s *PortMetaStore) Observe(port int, protocol string, at time.Time) {
	if m, ok := s.Entries[port]; ok {
		m.LastSeen = at
		m.SeenCount++
		return
	}
	s.Entries[port] = &PortMeta{
		Port:      port,
		Protocol:  protocol,
		FirstSeen: at,
		LastSeen:  at,
		SeenCount: 1,
	}
}

// SetLabel assigns a human-readable label to a port entry if it exists.
func (s *PortMetaStore) SetLabel(port int, label string) bool {
	if m, ok := s.Entries[port]; ok {
		m.Label = label
		return true
	}
	return false
}

// SavePortMeta writes the store to a JSON file at the given path.
func SavePortMeta(path string, store *PortMetaStore) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(store)
}

// LoadPortMeta reads a PortMetaStore from a JSON file. Returns an empty store
// if the file does not exist.
func LoadPortMeta(path string) (*PortMetaStore, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewPortMetaStore(), nil
		}
		return nil, err
	}
	defer f.Close()
	var store PortMetaStore
	if err := json.NewDecoder(f).Decode(&store); err != nil {
		return nil, err
	}
	if store.Entries == nil {
		store.Entries = make(map[int]*PortMeta)
	}
	return &store, nil
}
