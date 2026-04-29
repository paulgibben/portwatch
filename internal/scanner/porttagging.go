package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

// PortTagStore manages tags associated with port/protocol pairs.
type PortTagStore struct {
	mu   sync.RWMutex
	tags map[string][]string // key: "port/proto"
}

// NewPortTagStore creates an empty PortTagStore.
func NewPortTagStore() *PortTagStore {
	return &PortTagStore{
		tags: make(map[string][]string),
	}
}

func tagKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// AddTag adds a tag to a port/proto pair (no duplicates).
func (s *PortTagStore) AddTag(port int, proto, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := tagKey(port, proto)
	for _, t := range s.tags[key] {
		if t == tag {
			return
		}
	}
	s.tags[key] = append(s.tags[key], tag)
	sort.Strings(s.tags[key])
}

// RemoveTag removes a tag from a port/proto pair.
func (s *PortTagStore) RemoveTag(port int, proto, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := tagKey(port, proto)
	existing := s.tags[key]
	updated := existing[:0]
	for _, t := range existing {
		if t != tag {
			updated = append(updated, t)
		}
	}
	if len(updated) == 0 {
		delete(s.tags, key)
	} else {
		s.tags[key] = updated
	}
}

// GetTags returns all tags for a port/proto pair.
func (s *PortTagStore) GetTags(port int, proto string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := tagKey(port, proto)
	copy := append([]string(nil), s.tags[key]...)
	return copy
}

// All returns a snapshot of all tags keyed by "port/proto".
func (s *PortTagStore) All() map[string][]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string][]string, len(s.tags))
	for k, v := range s.tags {
		out[k] = append([]string(nil), v...)
	}
	return out
}

// SavePortTags persists the tag store to a JSON file.
func SavePortTags(path string, s *PortTagStore) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := json.MarshalIndent(s.tags, "", "  ")
	if err != nil {
		return fmt.Errorf("porttagging: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadPortTags loads a tag store from a JSON file.
func LoadPortTags(path string) (*PortTagStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewPortTagStore(), nil
		}
		return nil, fmt.Errorf("porttagging: read: %w", err)
	}
	var tags map[string][]string
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("porttagging: unmarshal: %w", err)
	}
	return &PortTagStore{tags: tags}, nil
}
