package notifier

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_JSONRoundtrip(t *testing.T) {
	cfg := &Config{
		Type:    TypeCommand,
		Command: "/usr/bin/notify-send",
		Format:  "text",
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "notifier.json")

	if err := SaveConfig(path, cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if loaded.Type != cfg.Type {
		t.Errorf("Type: got %q, want %q", loaded.Type, cfg.Type)
	}
	if loaded.Command != cfg.Command {
		t.Errorf("Command: got %q, want %q", loaded.Command, cfg.Command)
	}
}

func TestConfig_Defaults(t *testing.T) {
	cfg := &Config{}
	cfg.Defaults()
	if cfg.Type != TypeStdout {
		t.Errorf("expected default type stdout, got %q", cfg.Type)
	}
	if cfg.Format != "text" {
		t.Errorf("expected default format text, got %q", cfg.Format)
	}
}

func TestType_Constants(t *testing.T) {
	if TypeStdout != "stdout" {
		t.Errorf("TypeStdout = %q, want \"stdout\"", TypeStdout)
	}
	if TypeCommand != "command" {
		t.Errorf("TypeCommand = %q, want \"command\"", TypeCommand)
	}
	if TypeLog != "log" {
		t.Errorf("TypeLog = %q, want \"log\"", TypeLog)
	}
}

func TestConfig_Validate_MissingCommand(t *testing.T) {
	cfg := &Config{Type: TypeCommand}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for command type with empty command")
	}
}

func TestConfig_Validate_MissingLogFile(t *testing.T) {
	cfg := &Config{Type: TypeLog}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for log type with empty log_file")
	}
}

func TestLoadConfig_Missing(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Error("expected error loading missing config file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0644)
	_, err := LoadConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
