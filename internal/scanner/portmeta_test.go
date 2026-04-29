package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeMetaStore() *PortMetaStore {
	s := NewPortMetaStore()
	now := time.Now()
	s.Observe(80, "tcp", now)
	s.Observe(443, "tcp", now)
	return s
}

func TestNewPortMetaStore_Empty(t *testing.T) {
	s := NewPortMetaStore()
	if s.Entries == nil {
		t.Fatal("expected non-nil Entries map")
	}
	if len(s.Entries) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(s.Entries))
	}
}

func TestObserve_CreatesEntry(t *testing.T) {
	s := NewPortMetaStore()
	now := time.Now()
	s.Observe(8080, "tcp", now)

	m, ok := s.Entries[8080]
	if !ok {
		t.Fatal("expected entry for port 8080")
	}
	if m.SeenCount != 1 {
		t.Errorf("expected SeenCount=1, got %d", m.SeenCount)
	}
	if !m.FirstSeen.Equal(now) {
		t.Errorf("unexpected FirstSeen")
	}
}

func TestObserve_UpdatesExisting(t *testing.T) {
	s := NewPortMetaStore()
	t1 := time.Now()
	t2 := t1.Add(time.Minute)
	s.Observe(22, "tcp", t1)
	s.Observe(22, "tcp", t2)

	m := s.Entries[22]
	if m.SeenCount != 2 {
		t.Errorf("expected SeenCount=2, got %d", m.SeenCount)
	}
	if !m.FirstSeen.Equal(t1) {
		t.Errorf("FirstSeen should not change on update")
	}
	if !m.LastSeen.Equal(t2) {
		t.Errorf("LastSeen should be updated")
	}
}

func TestSetLabel(t *testing.T) {
	s := makeMetaStore()
	ok := s.SetLabel(80, "HTTP")
	if !ok {
		t.Fatal("expected SetLabel to return true for existing port")
	}
	if s.Entries[80].Label != "HTTP" {
		t.Errorf("expected label HTTP, got %s", s.Entries[80].Label)
	}

	ok = s.SetLabel(9999, "unknown")
	if ok {
		t.Error("expected SetLabel to return false for missing port")
	}
}

func TestSaveLoadPortMeta_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portmeta.json")

	orig := makeMetaStore()
	if err := SavePortMeta(path, orig); err != nil {
		t.Fatalf("SavePortMeta: %v", err)
	}

	loaded, err := LoadPortMeta(path)
	if err != nil {
		t.Fatalf("LoadPortMeta: %v", err)
	}
	if len(loaded.Entries) != len(orig.Entries) {
		t.Errorf("expected %d entries, got %d", len(orig.Entries), len(loaded.Entries))
	}
}

func TestLoadPortMeta_Missing(t *testing.T) {
	store, err := LoadPortMeta("/nonexistent/path/portmeta.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(store.Entries) != 0 {
		t.Error("expected empty store for missing file")
	}
}

func TestLoadPortMeta_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not-json{"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadPortMeta(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
