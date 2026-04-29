package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makePortStateStore() *PortStateStore {
	s := NewPortStateStore()
	now := time.Now()
	s.Observe(80, "tcp", true, now)
	s.Observe(443, "tcp", true, now)
	return s
}

func TestNewPortStateStore_Empty(t *testing.T) {
	s := NewPortStateStore()
	if s.Entries == nil {
		t.Fatal("expected non-nil Entries map")
	}
	if len(s.Entries) != 0 {
		t.Fatalf("expected empty store, got %d entries", len(s.Entries))
	}
}

func TestObserve_CreatesEntry(t *testing.T) {
	s := NewPortStateStore()
	now := time.Now()
	s.Observe(8080, "tcp", true, now)

	entry := s.Get(8080, "tcp")
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if !entry.Open {
		t.Error("expected port to be open")
	}
	if entry.SeenCount != 1 {
		t.Errorf("expected SeenCount=1, got %d", entry.SeenCount)
	}
	if entry.FirstSeen.IsZero() {
		t.Error("expected FirstSeen to be set")
	}
}

func TestObserve_UpdatesExisting(t *testing.T) {
	s := NewPortStateStore()
	t1 := time.Now()
	t2 := t1.Add(time.Minute)

	s.Observe(22, "tcp", true, t1)
	s.Observe(22, "tcp", false, t2)

	entry := s.Get(22, "tcp")
	if entry == nil {
		t.Fatal("expected entry")
	}
	if entry.Open {
		t.Error("expected port to be closed after second observe")
	}
	if entry.SeenCount != 2 {
		t.Errorf("expected SeenCount=2, got %d", entry.SeenCount)
	}
	if !entry.FirstSeen.Equal(t1) {
		t.Error("FirstSeen should not change on update")
	}
	if !entry.LastSeen.Equal(t2) {
		t.Errorf("expected LastSeen=%v, got %v", t2, entry.LastSeen)
	}
}

func TestGet_Missing(t *testing.T) {
	s := NewPortStateStore()
	if s.Get(9999, "tcp") != nil {
		t.Error("expected nil for missing port")
	}
}

func TestSaveLoadPortState_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portstate.json")

	orig := makePortStateStore()
	if err := SavePortState(path, orig); err != nil {
		t.Fatalf("SavePortState: %v", err)
	}

	loaded, err := LoadPortState(path)
	if err != nil {
		t.Fatalf("LoadPortState: %v", err)
	}
	if len(loaded.Entries) != len(orig.Entries) {
		t.Errorf("expected %d entries, got %d", len(orig.Entries), len(loaded.Entries))
	}
}

func TestLoadPortState_Missing(t *testing.T) {
	store, err := LoadPortState("/nonexistent/portstate.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(store.Entries) != 0 {
		t.Error("expected empty store for missing file")
	}
}

func TestLoadPortState_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not-json"), 0644)

	_, err := LoadPortState(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
