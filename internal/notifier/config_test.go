package notifier_test

import (
	"encoding/json"
	"testing"

	"github.com/user/portwatch/internal/notifier"
)

func TestConfig_JSONRoundtrip(t *testing.T) {
	orig := notifier.Config{
		Type:   notifier.TypeCommand,
		Target: "/usr/local/bin/notify",
		Format: "json",
	}

	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got notifier.Config
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if got.Type != orig.Type {
		t.Errorf("Type: got %q, want %q", got.Type, orig.Type)
	}
	if got.Target != orig.Target {
		t.Errorf("Target: got %q, want %q", got.Target, orig.Target)
	}
	if got.Format != orig.Format {
		t.Errorf("Format: got %q, want %q", got.Format, orig.Format)
	}
}

func TestConfig_Defaults(t *testing.T) {
	n := notifier.New(notifier.Config{})
	if n == nil {
		t.Fatal("expected non-nil notifier with empty config")
	}
}

func TestType_Constants(t *testing.T) {
	types := []notifier.Type{
		notifier.TypeStdout,
		notifier.TypeWebhook,
		notifier.TypeCommand,
	}
	for _, ty := range types {
		if ty == "" {
			t.Errorf("notifier type constant should not be empty")
		}
	}
}
