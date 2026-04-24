package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single alert event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      int
	Proto     string
}

// Notifier sends alerts to a destination.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier writing to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes an alert to the output writer.
func (n *Notifier) Notify(a Alert) error {
	_, err := fmt.Fprintf(
		n.out,
		"%s [%s] port=%d proto=%s msg=%s\n",
		a.Timestamp.UTC().Format(time.RFC3339),
		a.Level,
		a.Port,
		a.Proto,
		a.Message,
	)
	return err
}

// FromDiff converts scanner diff results into a slice of Alerts.
func FromDiff(added, removed []scanner.PortInfo) []Alert {
	var alerts []Alert
	for _, p := range added {
		alerts = append(alerts, Alert{
			Timestamp: time.Now(),
			Level:     LevelAlert,
			Message:   "new port opened",
			Port:      p.Port,
			Proto:     p.Proto,
		})
	}
	for _, p := range removed {
		alerts = append(alerts, Alert{
			Timestamp: time.Now(),
			Level:     LevelWarn,
			Message:   "port closed",
			Port:      p.Port,
			Proto:     p.Proto,
		})
	}
	return alerts
}
