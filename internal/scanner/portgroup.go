package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// PortGroup represents a named collection of ports with an optional description.
type PortGroup struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Ports       []int  `json:"ports"`
}

// PortGroupStore holds multiple named port groups.
type PortGroupStore struct {
	Groups []PortGroup `json:"groups"`
}

// NewPortGroupStore returns an empty PortGroupStore.
func NewPortGroupStore() *PortGroupStore {
	return &PortGroupStore{}
}

// Add adds or replaces a group by name.
func (s *PortGroupStore) Add(g PortGroup) {
	for i, existing := range s.Groups {
		if existing.Name == g.Name {
			s.Groups[i] = g
			return
		}
	}
	s.Groups = append(s.Groups, g)
}

// Get returns the group with the given name, or false if not found.
func (s *PortGroupStore) Get(name string) (PortGroup, bool) {
	for _, g := range s.Groups {
		if g.Name == name {
			return g, true
		}
	}
	return PortGroup{}, false
}

// AllPorts returns a deduplicated, sorted slice of all ports across all groups.
func (s *PortGroupStore) AllPorts() []int {
	seen := make(map[int]struct{})
	for _, g := range s.Groups {
		for _, p := range g.Ports {
			seen[p] = struct{}{}
		}
	}
	result := make([]int, 0, len(seen))
	for p := range seen {
		result = append(result, p)
	}
	sort.Ints(result)
	return result
}

// SavePortGroups writes the store to a JSON file.
func SavePortGroups(path string, store *PortGroupStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("portgroup: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPortGroups reads a PortGroupStore from a JSON file.
func LoadPortGroups(path string) (*PortGroupStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewPortGroupStore(), nil
		}
		return nil, fmt.Errorf("portgroup: read: %w", err)
	}
	var store PortGroupStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("portgroup: unmarshal: %w", err)
	}
	return &store, nil
}
