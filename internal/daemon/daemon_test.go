package daemon

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"portwatch/internal/rules"
)

func writeTempConfig(t *testing.T, cfg *rules.Config) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	return f.Name()
}

func TestNew_LoadsConfig(t *testing.T) {
	cfg := &rules.Config{
		ScanIntervalSeconds: 5,
		Ports:               []int{8080},
	}
	path := writeTempConfig(t, cfg)

	d, err := New(path)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Daemon")
	}
	if d.cfg.ScanIntervalSeconds != 5 {
		t.Errorf("expected interval 5, got %d", d.cfg.ScanIntervalSeconds)
	}
}

func TestNew_MissingConfig(t *testing.T) {
	_, err := New("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestRun_CancelContext(t *testing.T) {
	cfg := &rules.Config{
		ScanIntervalSeconds: 60,
		Ports:               []int{},
	}
	path := writeTempConfig(t, cfg)

	d, err := New(path)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- d.Run(ctx) }()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Run() unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Error("Run() did not return after context cancellation")
	}
}
