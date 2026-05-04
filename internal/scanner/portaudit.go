package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// AuditEventType describes the kind of audit event recorded.
type AuditEventType string

const (
	AuditEventOpened  AuditEventType = "opened"
	AuditEventClosed  AuditEventType = "closed"
	AuditEventScanned AuditEventType = "scanned"
	AuditEventLabeled AuditEventType = "labeled"
	AuditEventTagged  AuditEventType = "tagged"
)

// AuditEntry represents a single audit log entry for a port event.
type AuditEntry struct {
	Timestamp time.Time      `json:"timestamp"`
	EventType AuditEventType `json:"event_type"`
	Port      int            `json:"port"`
	Protocol  string         `json:"protocol"`
	Detail    string         `json:"detail,omitempty"`
}

// PortAuditLog maintains an ordered, bounded audit trail of port events.
type PortAuditLog struct {
	mu      sync.RWMutex
	entries []AuditEntry
	maxSize int
}

// NewPortAuditLog creates a new PortAuditLog with the given maximum number of
// retained entries. If maxSize is <= 0, a default of 500 is used.
func NewPortAuditLog(maxSize int) *PortAuditLog {
	if maxSize <= 0 {
		maxSize = 500
	}
	return &PortAuditLog{
		entries: make([]AuditEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Record appends an audit entry. If the log is at capacity the oldest entry
// is evicted to make room.
func (a *PortAuditLog) Record(event AuditEventType, port int, protocol, detail string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		EventType: event,
		Port:      port,
		Protocol:  protocol,
		Detail:    detail,
	}

	if len(a.entries) >= a.maxSize {
		a.entries = a.entries[1:]
	}
	a.entries = append(a.entries, entry)
}

// Entries returns a copy of all audit entries in chronological order.
func (a *PortAuditLog) Entries() []AuditEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	copy := make([]AuditEntry, len(a.entries))
	copy_ := copy // avoid shadowing builtin
	_ = copy_
	result := make([]AuditEntry, len(a.entries))
	for i, e := range a.entries {
		result[i] = e
	}
	return result
}

// FilterByPort returns audit entries matching the given port number.
func (a *PortAuditLog) FilterByPort(port int) []AuditEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var result []AuditEntry
	for _, e := range a.entries {
		if e.Port == port {
			result = append(result, e)
		}
	}
	return result
}

// FilterByEvent returns audit entries matching the given event type.
func (a *PortAuditLog) FilterByEvent(event AuditEventType) []AuditEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var result []AuditEntry
	for _, e := range a.entries {
		if e.EventType == event {
			result = append(result, e)
		}
	}
	return result
}

// SavePortAuditLog serialises the audit log to a JSON file at the given path.
func SavePortAuditLog(path string, log *PortAuditLog) error {
	entries := log.Entries()
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("portaudit: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("portaudit: write %s: %w", path, err)
	}
	return nil
}

// LoadPortAuditLog deserialises a PortAuditLog from a JSON file. If the file
// does not exist an empty log is returned without error.
func LoadPortAuditLog(path string, maxSize int) (*PortAuditLog, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewPortAuditLog(maxSize), nil
	}
	if err != nil {
		return nil, fmt.Errorf("portaudit: read %s: %w", path, err)
	}

	var entries []AuditEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("portaudit: unmarshal: %w", err)
	}

	log := NewPortAuditLog(maxSize)
	log.entries = entries
	return log, nil
}
