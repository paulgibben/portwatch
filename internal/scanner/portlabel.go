package scanner

import (
	"encoding/json"
	"fmt"
	"os"
)

// PortLabel maps a port number to a human-readable service label.
type PortLabel struct {
	Port    int    `json:"port"`
	Label   string `json:"label"`
	Comment string `json:"comment,omitempty"`
}

// PortLabelStore holds labels for known ports.
type PortLabelStore struct {
	Labels map[int]PortLabel `json:"labels"`
}

// NewPortLabelStore creates an empty PortLabelStore.
func NewPortLabelStore() *PortLabelStore {
	return &PortLabelStore{
		Labels: make(map[int]PortLabel),
	}
}

// Set adds or updates the label for a port.
func (s *PortLabelStore) Set(port int, label, comment string) {
	s.Labels[port] = PortLabel{Port: port, Label: label, Comment: comment}
}

// Get returns the label for a port, or an empty string if not found.
func (s *PortLabelStore) Get(port int) (PortLabel, bool) {
	l, ok := s.Labels[port]
	return l, ok
}

// Delete removes the label for a port.
func (s *PortLabelStore) Delete(port int) {
	delete(s.Labels, port)
}

// SavePortLabels writes the store to a JSON file.
func SavePortLabels(path string, store *PortLabelStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal port labels: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPortLabels reads a PortLabelStore from a JSON file.
func LoadPortLabels(path string) (*PortLabelStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewPortLabelStore(), nil
		}
		return nil, fmt.Errorf("read port labels: %w", err)
	}
	var store PortLabelStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("unmarshal port labels: %w", err)
	}
	if store.Labels == nil {
		store.Labels = make(map[int]PortLabel)
	}
	return &store, nil
}
