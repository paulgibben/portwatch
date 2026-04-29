package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTestSeverityConfig() PortSeverityConfig {
	return PortSeverityConfig{
		DefaultSeverity: SeverityLow,
		Rules: []PortSeverityRule{
			{Port: 22, Protocol: "tcp", Severity: SeverityHigh},
			{Port: 3306, Protocol: "tcp", Severity: SeverityCritical},
			{Port: 80, Severity: SeverityMedium},
		},
	}
}

func TestEvaluate_MatchesRule(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	if got := store.Evaluate(22, "tcp"); got != SeverityHigh {
		t.Errorf("expected high, got %s", got)
	}
}

func TestEvaluate_CriticalPort(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	if got := store.Evaluate(3306, "tcp"); got != SeverityCritical {
		t.Errorf("expected critical, got %s", got)
	}
}

func TestEvaluate_DefaultFallback(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	if got := store.Evaluate(9999, "tcp"); got != SeverityLow {
		t.Errorf("expected low, got %s", got)
	}
}

func TestEvaluate_ProtocolWildcard(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	// Rule for port 80 has no protocol set — should match any
	if got := store.Evaluate(80, "udp"); got != SeverityMedium {
		t.Errorf("expected medium, got %s", got)
	}
}

func TestEvaluate_CachesResult(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	_ = store.Evaluate(22, "tcp")
	if got := store.Evaluate(22, "tcp"); got != SeverityHigh {
		t.Errorf("expected cached high, got %s", got)
	}
}

func TestUpdateConfig_ClearsCache(t *testing.T) {
	store := NewPortSeverityStore(makeTestSeverityConfig())
	_ = store.Evaluate(22, "tcp") // populate cache

	newCfg := PortSeverityConfig{
		DefaultSeverity: SeverityMedium,
		Rules:           []PortSeverityRule{},
	}
	store.UpdateConfig(newCfg)

	if got := store.Evaluate(22, "tcp"); got != SeverityMedium {
		t.Errorf("expected medium after config update, got %s", got)
	}
}

func TestSaveLoadPortSeverityConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "severity.json")
	cfg := makeTestSeverityConfig()

	if err := SavePortSeverityConfig(path, cfg); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := LoadPortSeverityConfig(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded.DefaultSeverity != cfg.DefaultSeverity {
		t.Errorf("default severity mismatch")
	}
	if len(loaded.Rules) != len(cfg.Rules) {
		t.Errorf("rules count mismatch: got %d want %d", len(loaded.Rules), len(cfg.Rules))
	}
}

func TestLoadPortSeverityConfig_Missing(t *testing.T) {
	cfg, err := LoadPortSeverityConfig("/nonexistent/path/severity.json")
	if err != nil {
		t.Fatalf("expected default on missing file, got error: %v", err)
	}
	if cfg.DefaultSeverity == "" {
		t.Error("expected non-empty default severity")
	}
}

func TestLoadPortSeverityConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("{invalid"), 0644)
	_, err := LoadPortSeverityConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
