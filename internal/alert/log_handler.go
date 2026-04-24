package alert

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogHandler writes alerts to a rotating log file.
type LogHandler struct {
	mu      sync.Mutex
	path    string
	notifier *Notifier
}

// NewLogHandler creates a LogHandler that appends alerts to the given file path.
// The directory is created if it does not exist.
func NewLogHandler(path string) (*LogHandler, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("alert: create log dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("alert: open log file: %w", err)
	}
	return &LogHandler{
		path:     path,
		notifier: New(f),
	}, nil
}

// Handle writes a batch of alerts to the log file, thread-safely.
func (h *LogHandler) Handle(alerts []Alert) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, a := range alerts {
		if a.Timestamp.IsZero() {
			a.Timestamp = time.Now()
		}
		if err := h.notifier.Notify(a); err != nil {
			return fmt.Errorf("alert: write log: %w", err)
		}
	}
	return nil
}

// Path returns the log file path used by this handler.
func (h *LogHandler) Path() string {
	return h.path
}
