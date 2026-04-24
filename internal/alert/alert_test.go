package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeAlert(level Level, port int, proto, msg string) Alert {
	return Alert{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     level,
		Message:   msg,
		Port:      port,
		Proto:     proto,
	}
}

func TestNotify_Output(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)
	a := makeAlert(LevelAlert, 8080, "tcp", "new port opened")

	if err := n.Notify(a); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[ALERT]") {
		t.Errorf("expected [ALERT] in output, got: %s", got)
	}
	if !strings.Contains(got, "port=8080") {
		t.Errorf("expected port=8080 in output, got: %s", got)
	}
	if !strings.Contains(got, "proto=tcp") {
		t.Errorf("expected proto=tcp in output, got: %s", got)
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	n := New(nil)
	if n.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestFromDiff_Added(t *testing.T) {
	added := []scanner.PortInfo{
		{Port: 443, Proto: "tcp"},
		{Port: 8443, Proto: "tcp"},
	}
	alerts := FromDiff(added, nil)

	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
	for _, a := range alerts {
		if a.Level != LevelAlert {
			t.Errorf("expected ALERT level for added port, got %s", a.Level)
		}
	}
}

func TestFromDiff_Removed(t *testing.T) {
	removed := []scanner.PortInfo{
		{Port: 22, Proto: "tcp"},
	}
	alerts := FromDiff(nil, removed)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != LevelWarn {
		t.Errorf("expected WARN level for removed port, got %s", alerts[0].Level)
	}
}

func TestFromDiff_Empty(t *testing.T) {
	alerts := FromDiff(nil, nil)
	if len(alerts) != 0 {
		t.Errorf("expected no alerts, got %d", len(alerts))
	}
}
