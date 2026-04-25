package reporter

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format defines the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report holds a snapshot of port state at a point in time.
type Report struct {
	Timestamp time.Time         `json:"timestamp"`
	Ports     []scanner.Port    `json:"ports"`
	Changes   *scanner.Diff     `json:"changes,omitempty"`
}

// Reporter writes port reports to an output destination.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter writing to w in the given format.
// If w is nil, os.Stdout is used.
func New(w io.Writer, format Format) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Reporter{out: w, format: format}
}

// Write renders r to the reporter's output.
func (r *Reporter) Write(rep Report) error {
	switch r.format {
	case FormatJSON:
		return writeJSON(r.out, rep)
	default:
		return writeText(r.out, rep)
	}
}

func writeText(w io.Writer, rep Report) error {
	_, err := fmt.Fprintf(w, "=== Port Report [%s] ===\n", rep.Timestamp.Format(time.RFC3339))
	if err != nil {
		return err
	}
	for _, p := range rep.Ports {
		if _, err := fmt.Fprintf(w, "  %s/%d\n", p.Protocol, p.Number); err != nil {
			return err
		}
	}
	if rep.Changes != nil {
		if _, err := fmt.Fprintf(w, "Added: %d  Removed: %d\n",
			len(rep.Changes.Added), len(rep.Changes.Removed)); err != nil {
			return err
		}
	}
	return nil
}
