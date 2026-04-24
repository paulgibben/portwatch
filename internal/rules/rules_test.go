package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTestConfig() *Config {
	return &Config{
		Rules: []Rule{
			{Port: 22, Protocol: "tcp", Action: ActionIgnore, Comment: "SSH"},
			{Port: 80, Protocol: "tcp", Action: ActionIgnore, Comment: "HTTP"},
			{Port: 9999, Protocol: "tcp", Action: ActionAlert, Comment: "suspicious"},
		},
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

	if len(loaded.Rules) != len(cfg.Rules) {
		t.Fatalf("expected %d rules, got %d", len(cfg.Rules), len(loaded.Rules))
	}
	for i, r := range loaded.Rules {
		if r.Port != cfg.Rules[i].Port || r.Protocol != cfg.Rules[i].Protocol || r.Action != cfg.Rules[i].Action {
			t.Errorf("rule %d mismatch: got %+v, want %+v", i, r, cfg.Rules[i])
		}
	}
}

func TestLoadConfig_Missing(t *testing.T) {
	_, err := LoadConfig("/nonexistent/rules.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestEvaluate_MatchIgnore(t *testing.T) {
	cfg := makeTestConfig()
	if got := cfg.Evaluate(22, "tcp"); got != ActionIgnore {
		t.Errorf("expected ignore for port 22/tcp, got %s", got)
	}
}

func TestEvaluate_MatchAlert(t *testing.T) {
	cfg := makeTestConfig()
	if got := cfg.Evaluate(9999, "tcp"); got != ActionAlert {
		t.Errorf("expected alert for port 9999/tcp, got %s", got)
	}
}

func TestEvaluate_NoMatch_DefaultsToAlert(t *testing.T) {
	cfg := makeTestConfig()
	if got := cfg.Evaluate(12345, "tcp"); got != ActionAlert {
		t.Errorf("expected default alert for unknown port, got %s", got)
	}
}
