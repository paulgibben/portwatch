package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// PortComment holds a user-defined comment/annotation for a specific port.
type PortComment struct {
	Port      int       `json:"port"`
	Proto     string    `json:"proto"`
	Comment   string    `json:"comment"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PortCommentStore manages comments keyed by port+proto.
type PortCommentStore struct {
	mu       sync.RWMutex
	comments map[string]*PortComment
}

// NewPortCommentStore returns an empty PortCommentStore.
func NewPortCommentStore() *PortCommentStore {
	return &PortCommentStore{
		comments: make(map[string]*PortComment),
	}
}

func commentKey(port int, proto string) string {
	return fmt.Sprintf("%d/%s", port, proto)
}

// Set adds or updates a comment for the given port and protocol.
func (s *PortCommentStore) Set(port int, proto, comment string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.comments[commentKey(port, proto)] = &PortComment{
		Port:      port,
		Proto:     proto,
		Comment:   comment,
		UpdatedAt: time.Now().UTC(),
	}
}

// Get retrieves the comment for a port/proto pair. Returns nil if not found.
func (s *PortCommentStore) Get(port int, proto string) *PortComment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.comments[commentKey(port, proto)]
}

// Delete removes the comment for a port/proto pair.
func (s *PortCommentStore) Delete(port int, proto string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.comments, commentKey(port, proto))
}

// All returns a slice of all stored comments.
func (s *PortCommentStore) All() []*PortComment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*PortComment, 0, len(s.comments))
	for _, c := range s.comments {
		out = append(out, c)
	}
	return out
}

// SavePortComments writes the store to a JSON file.
func SavePortComments(path string, s *PortCommentStore) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := json.MarshalIndent(s.comments, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPortComments reads a PortCommentStore from a JSON file.
func LoadPortComments(path string) (*PortCommentStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewPortCommentStore(), nil
		}
		return nil, err
	}
	var comments map[string]*PortComment
	if err := json.Unmarshal(data, &comments); err != nil {
		return nil, err
	}
	return &PortCommentStore{comments: comments}, nil
}
