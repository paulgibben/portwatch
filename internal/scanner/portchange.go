package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ChangeType describes the kind of port change observed.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
)

// PortChange records a single port change event.
type PortChange struct {
	Port      int        `json:"port"`
	Proto     string     `json:"proto"`
	Change    ChangeType `json:"change"`
	DetectedAt time.Time `json:"detected_at"`
}

// PortChangeLog holds an ordered list of port change events.
type PortChangeLog struct {
	Entries []PortChange `json:"entries"`
}

// NewPortChangeLog returns an empty PortChangeLog.
func NewPortChangeLog() *PortChangeLog {
	return &PortChangeLog{}
}

// Record appends a new change event to the log.
func (l *PortChangeLog) Record(port int, proto string, ct ChangeType) {
	l.Entries = append(l.Entries, PortChange{
		Port:       port,
		Proto:      proto,
		Change:     ct,
		DetectedAt: time.Now().UTC(),
	})
}

// SavePortChangeLog writes the log to a JSON file.
func SavePortChangeLog(path string, log *PortChangeLog) error {
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("portchange: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPortChangeLog reads a PortChangeLog from a JSON file.
// Returns an empty log if the file does not exist.
func LoadPortChangeLog(path string) (*PortChangeLog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewPortChangeLog(), nil
		}
		return nil, fmt.Errorf("portchange: read: %w", err)
	}
	var log PortChangeLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("portchange: unmarshal: %w", err)
	}
	return &log, nil
}
