package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func makeTagStore() *PortTagStore {
	s := NewPortTagStore()
	s.AddTag(80, "tcp", "web")
	s.AddTag(80, "tcp", "http")
	s.AddTag(443, "tcp", "web")
	s.AddTag(443, "tcp", "tls")
	return s
}

func TestNewPortTagStore_Empty(t *testing.T) {
	s := NewPortTagStore()
	if s == nil {
		t.Fatal("expected non-nil store")
	}
	if len(s.All()) != 0 {
		t.Errorf("expected empty store, got %d entries", len(s.All()))
	}
}

func TestPortTagStore_AddAndGet(t *testing.T) {
	s := NewPortTagStore()
	s.AddTag(22, "tcp", "ssh")
	s.AddTag(22, "tcp", "admin")
	tags := s.GetTags(22, "tcp")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "admin" || tags[1] != "ssh" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestPortTagStore_AddDuplicate(t *testing.T) {
	s := NewPortTagStore()
	s.AddTag(22, "tcp", "ssh")
	s.AddTag(22, "tcp", "ssh")
	if len(s.GetTags(22, "tcp")) != 1 {
		t.Error("expected duplicate tag to be ignored")
	}
}

func TestPortTagStore_RemoveTag(t *testing.T) {
	s := makeTagStore()
	s.RemoveTag(80, "tcp", "http")
	tags := s.GetTags(80, "tcp")
	if len(tags) != 1 || tags[0] != "web" {
		t.Errorf("expected [web], got %v", tags)
	}
}

func TestPortTagStore_RemoveLastTag_DeletesKey(t *testing.T) {
	s := NewPortTagStore()
	s.AddTag(9000, "tcp", "custom")
	s.RemoveTag(9000, "tcp", "custom")
	all := s.All()
	if _, ok := all["9000/tcp"]; ok {
		t.Error("expected key to be removed when last tag deleted")
	}
}

func TestPortTagStore_GetMissing(t *testing.T) {
	s := NewPortTagStore()
	tags := s.GetTags(9999, "tcp")
	if tags != nil && len(tags) != 0 {
		t.Errorf("expected empty slice for missing key, got %v", tags)
	}
}

func TestSaveLoadPortTags_RoundTrip(t *testing.T) {
	s := makeTagStore()
	dir := t.TempDir()
	path := filepath.Join(dir, "tags.json")
	if err := SavePortTags(path, s); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadPortTags(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	for key, expected := range s.All() {
		got := loaded.All()[key]
		if len(got) != len(expected) {
			t.Errorf("key %s: expected %v, got %v", key, expected, got)
		}
	}
}

func TestLoadPortTags_Missing(t *testing.T) {
	s, err := LoadPortTags("/nonexistent/tags.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.All()) != 0 {
		t.Error("expected empty store for missing file")
	}
}

func TestLoadPortTags_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadPortTags(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
