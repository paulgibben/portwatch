package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultPortTaggingConfig(t *testing.T) {
	cfg := DefaultPortTaggingConfig()
	if cfg.PersistPath == "" {
		t.Error("expected non-empty PersistPath")
	}
	if len(cfg.DefaultTags) == 0 {
		t.Error("expected non-empty DefaultTags")
	}
	if _, ok := cfg.DefaultTags["443/tcp"]; !ok {
		t.Error("expected 443/tcp in DefaultTags")
	}
}

func TestSaveLoadPortTaggingConfig_RoundTrip(t *testing.T) {
	cfg := PortTaggingConfig{
		PersistPath: "/var/portwatch/tags.json",
		DefaultTags: map[string][]string{
			"8080/tcp": {"proxy", "web"},
		},
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "tagging.json")
	if err := SavePortTaggingConfig(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadPortTaggingConfig(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.PersistPath != cfg.PersistPath {
		t.Errorf("PersistPath: got %q, want %q", loaded.PersistPath, cfg.PersistPath)
	}
	tags := loaded.DefaultTags["8080/tcp"]
	if len(tags) != 2 {
		t.Errorf("expected 2 tags for 8080/tcp, got %v", tags)
	}
}

func TestLoadPortTaggingConfig_Missing(t *testing.T) {
	cfg, err := LoadPortTaggingConfig("/nonexistent/tagging.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg.PersistPath == "" {
		t.Error("expected default PersistPath")
	}
}

func TestLoadPortTaggingConfig_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("!!!"), 0o644)
	_, err := LoadPortTaggingConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApplyDefaultTags(t *testing.T) {
	cfg := DefaultPortTaggingConfig()
	store := NewPortTagStore()
	ApplyDefaultTags(store, cfg)

	tags := store.GetTags(80, "tcp")
	if len(tags) == 0 {
		t.Error("expected tags for 80/tcp after applying defaults")
	}
	found := false
	for _, tag := range tags {
		if tag == "http" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'http' tag for 80/tcp, got %v", tags)
	}
}

func TestApplyDefaultTags_SkipsBadKeys(t *testing.T) {
	cfg := PortTaggingConfig{
		DefaultTags: map[string][]string{
			"badkey": {"sometag"},
			"22/tcp": {"ssh"},
		},
	}
	store := NewPortTagStore()
	ApplyDefaultTags(store, cfg)
	// Should not panic; only valid key should be applied
	if len(store.GetTags(22, "tcp")) != 1 {
		t.Error("expected ssh tag on 22/tcp")
	}
}
