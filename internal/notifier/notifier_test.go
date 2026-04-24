package notifier_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/notifier"
)

func makeAlert(kind, proto string, port int) *alert.Alert {
	return &alert.Alert{
		Events: []alert.Event{
			{Kind: kind, Port: port, Proto: proto},
		},
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	n := notifier.New(notifier.Config{})
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNew_WithType(t *testing.T) {
	n := notifier.New(notifier.Config{Type: notifier.TypeStdout})
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSend_Stdout(t *testing.T) {
	n := notifier.New(notifier.Config{Type: notifier.TypeStdout})
	a := makeAlert("added", "tcp", 8080)
	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSend_Command(t *testing.T) {
	n := notifier.New(notifier.Config{
		Type:   notifier.TypeCommand,
		Target: "cat",
	})
	a := makeAlert("removed", "tcp", 443)
	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error sending via command: %v", err)
	}
}

func TestSend_Command_NoTarget(t *testing.T) {
	n := notifier.New(notifier.Config{Type: notifier.TypeCommand})
	a := makeAlert("added", "tcp", 22)
	if err := n.Send(a); err == nil {
		t.Fatal("expected error when command target is empty")
	}
}

func TestSend_UnknownType(t *testing.T) {
	n := notifier.New(notifier.Config{Type: "unknown"})
	a := makeAlert("added", "tcp", 80)
	if err := n.Send(a); err == nil {
		t.Fatal("expected error for unsupported notifier type")
	}
}
