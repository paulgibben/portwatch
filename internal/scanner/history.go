package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// HistoryEntry records a snapshot diff event with a timestamp.
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Added     []int     `json:"added"`
	Removed   []int     `json:"removed"`
}

// History holds an ordered list of diff entries.
type History struct {
	Entries []HistoryEntry `json:"entries"`
	MaxSize int            `json:"-"`
}

// NewHistory creates a History with an optional max size (0 = unlimited).
func NewHistory(maxSize int) *History {
	return &History{MaxSize: maxSize}
}

// Record appends a new entry derived from a Diff result.
func (h *History) Record(diff map[string][]int) {
	entry := HistoryEntry{
		Timestamp: time.Now().UTC(),
		Added:     diff["added"],
		Removed:   diff["removed"],
	}
	h.Entries = append(h.Entries, entry)
	if h.MaxSize > 0 && len(h.Entries) > h.MaxSize {
		h.Entries = h.Entries[len(h.Entries)-h.MaxSize:]
	}
}

// SaveHistory writes history to a JSON file.
func SaveHistory(path string, h *History) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadHistory reads history from a JSON file.
// Returns an empty History if the file does not exist.
func LoadHistory(path string, maxSize int) (*History, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewHistory(maxSize), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read history: %w", err)
	}
	h := NewHistory(maxSize)
	if err := json.Unmarshal(data, h); err != nil {
		return nil, fmt.Errorf("unmarshal history: %w", err)
	}
	h.MaxSize = maxSize
	return h, nil
}
