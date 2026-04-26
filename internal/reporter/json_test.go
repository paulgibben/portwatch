package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"portwatch/internal/scanner"
)

// TestWriteJSON_BasicReport verifies that a report is serialized to valid JSON.
func TestWriteJSON_BasicReport(t *testing.T) {
	r := makeReport()

	var buf bytes.Buffer
	if err := writeJSON(&buf, r); err != nil {
		t.Fatalf("writeJSON returned error: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("expected non-empty JSON output")
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
}

// TestWriteJSON_ContainsTimestamp checks that the JSON output includes a timestamp field.
func TestWriteJSON_ContainsTimestamp(t *testing.T) {
	r := makeReport()

	var buf bytes.Buffer
	if err := writeJSON(&buf, r); err != nil {
		t.Fatalf("writeJSON returned error: %v", err)
	}

	if !strings.Contains(buf.String(), "timestamp") {
		t.Errorf("expected JSON to contain 'timestamp', got: %s", buf.String())
	}
}

// TestWriteJSON_ContainsPorts checks that open ports appear in the JSON output.
func TestWriteJSON_ContainsPorts(t *testing.T) {
	r := makeReport()

	var buf bytes.Buffer
	if err := writeJSON(&buf, r); err != nil {
		t.Fatalf("writeJSON returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "open_ports") {
		t.Errorf("expected JSON to contain 'open_ports', got: %s", output)
	}
}

// TestParseReport_RoundTrip verifies that a report written as JSON can be parsed back.
func TestParseReport_RoundTrip(t *testing.T) {
	original := makeReport()

	var buf bytes.Buffer
	if err := writeJSON(&buf, original); err != nil {
		t.Fatalf("writeJSON returned error: %v", err)
	}

	parsed, err := ParseReport(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ParseReport returned error: %v", err)
	}

	if len(parsed.OpenPorts) != len(original.OpenPorts) {
		t.Errorf("expected %d open ports, got %d", len(original.OpenPorts), len(parsed.OpenPorts))
	}
}

// TestParseReport_InvalidJSON ensures ParseReport returns an error on bad input.
func TestParseReport_InvalidJSON(t *testing.T) {
	_, err := ParseReport(strings.NewReader("not json at all"))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// TestParseReport_WithChanges verifies that diff changes survive a JSON round-trip.
func TestParseReport_WithChanges(t *testing.T) {
	r := makeReport()
	r.Changes = &scanner.Diff{
		Added: []scanner.Port{{Number: 9090, Protocol: "tcp"}},
		Removed: []scanner.Port{{Number: 22, Protocol: "tcp"}},
		Timestamp: time.Now().UTC(),
	}

	var buf bytes.Buffer
	if err := writeJSON(&buf, r); err != nil {
		t.Fatalf("writeJSON returned error: %v", err)
	}

	parsed, err := ParseReport(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ParseReport returned error: %v", err)
	}

	if parsed.Changes == nil {
		t.Fatal("expected non-nil Changes after round-trip")
	}
	if len(parsed.Changes.Added) != 1 {
		t.Errorf("expected 1 added port, got %d", len(parsed.Changes.Added))
	}
	if len(parsed.Changes.Removed) != 1 {
		t.Errorf("expected 1 removed port, got %d", len(parsed.Changes.Removed))
	}
}
