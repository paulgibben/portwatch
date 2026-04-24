package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Baseline represents a saved reference state of open ports used
// to detect deviations on subsequent scans.
type Baseline struct {
	CreatedAt time.Time `json:"created_at"`
	Host      string    `json:"host"`
	Ports     []int     `json:"ports"`
}

// NewBaseline creates a Baseline from the current open ports on host.
func NewBaseline(host string, ports []int) *Baseline {
	copy := make([]int, len(ports))
	for i, p := range ports {
		copy[i] = p
	}
	return &Baseline{
		CreatedAt: time.Now().UTC(),
		Host:      host,
		Ports:     copy,
	}
}

// SaveBaseline writes the baseline to a JSON file at path.
func SaveBaseline(b *Baseline, path string) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("baseline: write %s: %w", path, err)
	}
	return nil
}

// LoadBaseline reads a baseline from a JSON file at path.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("baseline: file not found: %s", path)
		}
		return nil, fmt.Errorf("baseline: read %s: %w", path, err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// PortSet returns the baseline ports as a map for O(1) lookup.
func (b *Baseline) PortSet() map[int]struct{} {
	set := make(map[int]struct{}, len(b.Ports))
	for _, p := range b.Ports {
		set[p] = struct{}{}
	}
	return set
}

// Unexpected returns ports present in current that are absent from the baseline.
func (b *Baseline) Unexpected(current []int) []int {
	known := b.PortSet()
	var unexpected []int
	for _, p := range current {
		if _, ok := known[p]; !ok {
			unexpected = append(unexpected, p)
		}
	}
	return unexpected
}
