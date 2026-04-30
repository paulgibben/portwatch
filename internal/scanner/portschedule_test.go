package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultScanScheduleConfig(t *testing.T) {
	cfg := DefaultScanScheduleConfig()
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected 60s interval, got %v", cfg.Interval)
	}
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if cfg.MaxMissed != 3 {
		t.Errorf("expected MaxMissed=3, got %d", cfg.MaxMissed)
	}
}

func TestSaveLoadScanScheduleConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schedule.json")

	orig := ScanScheduleConfig{
		Interval:   30 * time.Second,
		StartDelay: 5 * time.Second,
		MaxMissed:  5,
		Enabled:    true,
	}
	if err := SaveScanScheduleConfig(path, orig); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadScanScheduleConfig(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Interval != orig.Interval || loaded.MaxMissed != orig.MaxMissed {
		t.Errorf("round-trip mismatch: %+v vs %+v", loaded, orig)
	}
}

func TestLoadScanScheduleConfig_Missing(t *testing.T) {
	cfg, err := LoadScanScheduleConfig("/nonexistent/path/schedule.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected default interval, got %v", cfg.Interval)
	}
}

func TestLoadScanScheduleConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json{"), 0644)
	_, err := LoadScanScheduleConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestScanScheduleConfig_JSONFields(t *testing.T) {
	cfg := DefaultScanScheduleConfig()
	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	for _, key := range []string{"interval", "start_delay", "max_missed", "enabled"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing JSON key: %s", key)
		}
	}
}
