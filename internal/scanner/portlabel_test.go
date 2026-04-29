package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func makePortLabelStore() *PortLabelStore {
	s := NewPortLabelStore()
	s.Set(80, "HTTP", "web traffic")
	s.Set(443, "HTTPS", "secure web traffic")
	s.Set(22, "SSH", "")
	return s
}

func TestNewPortLabelStore_Empty(t *testing.T) {
	s := NewPortLabelStore()
	if len(s.Labels) != 0 {
		t.Errorf("expected empty store, got %d entries", len(s.Labels))
	}
}

func TestPortLabelStore_SetAndGet(t *testing.T) {
	s := NewPortLabelStore()
	s.Set(8080, "HTTP-Alt", "alternate HTTP")
	l, ok := s.Get(8080)
	if !ok {
		t.Fatal("expected label to exist")
	}
	if l.Label != "HTTP-Alt" {
		t.Errorf("expected HTTP-Alt, got %s", l.Label)
	}
	if l.Comment != "alternate HTTP" {
		t.Errorf("expected 'alternate HTTP', got %s", l.Comment)
	}
}

func TestPortLabelStore_Get_Missing(t *testing.T) {
	s := NewPortLabelStore()
	_, ok := s.Get(9999)
	if ok {
		t.Error("expected missing label, got found")
	}
}

func TestPortLabelStore_Delete(t *testing.T) {
	s := makePortLabelStore()
	s.Delete(80)
	_, ok := s.Get(80)
	if ok {
		t.Error("expected label to be deleted")
	}
}

func TestSaveLoadPortLabels_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")

	orig := makePortLabelStore()
	if err := SavePortLabels(path, orig); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded, err := LoadPortLabels(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Labels) != len(orig.Labels) {
		t.Errorf("expected %d labels, got %d", len(orig.Labels), len(loaded.Labels))
	}
	l, ok := loaded.Get(443)
	if !ok || l.Label != "HTTPS" {
		t.Errorf("expected HTTPS label for port 443, got %+v", l)
	}
}

func TestLoadPortLabels_Missing(t *testing.T) {
	store, err := LoadPortLabels("/nonexistent/path/labels.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(store.Labels) != 0 {
		t.Errorf("expected empty store for missing file")
	}
}

func TestLoadPortLabels_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")
	if err := os.WriteFile(path, []byte("not-json"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadPortLabels(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
