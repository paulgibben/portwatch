package scanner

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of open ports.
type Snapshot struct {
	Timestamp time.Time   `json:"timestamp"`
	Host      string      `json:"host"`
	Ports     []PortState `json:"ports"`
}

// NewSnapshot creates a snapshot from a list of port states.
func NewSnapshot(host string, states []PortState) *Snapshot {
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Ports:     OpenPorts(states),
	}
}

// SaveToFile serializes the snapshot to a JSON file.
func (s *Snapshot) SaveToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadFromFile deserializes a snapshot from a JSON file.
func LoadFromFile(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

// Diff compares two snapshots and returns newly opened and closed ports.
func Diff(prev, curr *Snapshot) (opened, closed []PortState) {
	prevMap := make(map[int]struct{}, len(prev.Ports))
	currMap := make(map[int]struct{}, len(curr.Ports))

	for _, p := range prev.Ports {
		prevMap[p.Port] = struct{}{}
	}
	for _, p := range curr.Ports {
		currMap[p.Port] = struct{}{}
	}
	for _, p := range curr.Ports {
		if _, found := prevMap[p.Port]; !found {
			opened = append(opened, p)
		}
	}
	for _, p := range prev.Ports {
		if _, found := currMap[p.Port]; !found {
			closed = append(closed, p)
		}
	}
	return
}
