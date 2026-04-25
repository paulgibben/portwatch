package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTestConfig() *Config {
	return &Config{
		PortRanges:  []PortRange{{Start: 1, End: 1024}, {Start: 8000, End: 9000}},
		IgnorePorts: []int{22, 80},
		AlertOnNew:  true,
		AlertOnGone: false,
	}
}

func TestSaveLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.json")
	cfg := makeTestConfig()

	if err := SaveConfig(path, cfg); err != nil {
		t.Fatalf("SaveConfig: %v", err)
	}
	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if loaded.AlertOnNew != cfg.AlertOnNew {
		t.Errorf("AlertOnNew mismatch: got %v, want %v", loaded.AlertOnNew, cfg.AlertOnNew)
	}
	if len(loaded.IgnorePorts) != len(cfg.IgnorePorts) {
		t.Errorf("IgnorePorts length mismatch")
	}
}

func TestLoadConfig_Missing(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/rules.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestEvaluate_MatchIgnore(t *testing.T) {
	cfg := makeTestConfig()
	alert, ignore := cfg.Evaluate(22)
	if !ignore {
		t.Error("expected port 22 to be ignored")
	}
	if alert {
		t.Error("expected no alert for ignored port")
	}
}

func TestEvaluate_MatchRange(t *testing.T) {
	cfg := makeTestConfig()
	alert, ignore := cfg.Evaluate(443)
	if ignore {
		t.Error("expected port 443 not to be ignored")
	}
	if !alert {
		t.Error("expected alert for port 443 in range")
	}
}

func TestEvaluate_OutOfRange(t *testing.T) {
	cfg := makeTestConfig()
	alert, ignore := cfg.Evaluate(5000)
	if ignore {
		t.Error("expected port 5000 not to be ignored")
	}
	if alert {
		t.Error("expected no alert for port 5000 out of range")
	}
}

func TestEvaluate_NoRanges(t *testing.T) {
	cfg := &Config{AlertOnNew: true}
	alert, ignore := cfg.Evaluate(9999)
	if ignore {
		t.Error("expected port not to be ignored")
	}
	if !alert {
		t.Error("expected alert when no ranges configured")
	}
}
