package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultTimeoutConfig(t *testing.T) {
	cfg := DefaultTimeoutConfig()
	if cfg.PortTimeout <= 0 {
		t.Errorf("expected positive PortTimeout, got %v", cfg.PortTimeout)
	}
	if cfg.GlobalTimeout <= 0 {
		t.Errorf("expected positive GlobalTimeout, got %v", cfg.GlobalTimeout)
	}
}

func TestSaveLoadTimeoutConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "timeout.json")

	orig := TimeoutConfig{
		PortTimeout:   250 * time.Millisecond,
		GlobalTimeout: 10 * time.Second,
	}

	if err := SaveTimeoutConfig(path, orig); err != nil {
		t.Fatalf("SaveTimeoutConfig: %v", err)
	}

	loaded, err := LoadTimeoutConfig(path)
	if err != nil {
		t.Fatalf("LoadTimeoutConfig: %v", err)
	}

	if loaded.PortTimeout != orig.PortTimeout {
		t.Errorf("PortTimeout: want %v, got %v", orig.PortTimeout, loaded.PortTimeout)
	}
	if loaded.GlobalTimeout != orig.GlobalTimeout {
		t.Errorf("GlobalTimeout: want %v, got %v", orig.GlobalTimeout, loaded.GlobalTimeout)
	}
}

func TestLoadTimeoutConfig_Missing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.json")

	cfg, err := LoadTimeoutConfig(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	def := DefaultTimeoutConfig()
	if cfg.PortTimeout != def.PortTimeout {
		t.Errorf("PortTimeout: want %v, got %v", def.PortTimeout, cfg.PortTimeout)
	}
}

func TestLoadTimeoutConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadTimeoutConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestSaveTimeoutConfig_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "timeout.json")

	cfg := DefaultTimeoutConfig()
	if err := SaveTimeoutConfig(path, cfg); err != nil {
		t.Fatalf("SaveTimeoutConfig: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to be created")
	}
}
