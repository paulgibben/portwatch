package reporter

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeReport() Report {
	return Report{
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Ports: []scanner.Port{
			{Number: 80, Protocol: "tcp"},
			{Number: 443, Protocol: "tcp"},
		},
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	r := New(nil, "")
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
	if r.format != FormatText {
		t.Errorf("expected default format text, got %q", r.format)
	}
}

func TestWrite_Text(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	if err := r.Write(makeReport()); err != nil {
		t.Fatalf("Write: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Port Report") {
		t.Errorf("expected header in output, got: %s", out)
	}
	if !strings.Contains(out, "tcp/80") {
		t.Errorf("expected port 80 in output, got: %s", out)
	}
}

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatJSON)
	if err := r.Write(makeReport()); err != nil {
		t.Fatalf("Write: %v", err)
	}
	parsed, err := ParseReport(&buf)
	if err != nil {
		t.Fatalf("ParseReport: %v", err)
	}
	if len(parsed.Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(parsed.Ports))
	}
}

func TestWrite_TextWithChanges(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf, FormatText)
	rep := makeReport()
	rep.Changes = &scanner.Diff{
		Added:   []scanner.Port{{Number: 8080, Protocol: "tcp"}},
		Removed: []scanner.Port{},
	}
	if err := r.Write(rep); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if !strings.Contains(buf.String(), "Added: 1") {
		t.Errorf("expected change summary, got: %s", buf.String())
	}
}
