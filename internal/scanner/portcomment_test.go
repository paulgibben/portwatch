package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func makeCommentStore() *PortCommentStore {
	s := NewPortCommentStore()
	s.Set(80, "tcp", "HTTP traffic")
	s.Set(443, "tcp", "HTTPS traffic")
	return s
}

func TestNewPortCommentStore_Empty(t *testing.T) {
	s := NewPortCommentStore()
	if len(s.All()) != 0 {
		t.Fatalf("expected empty store, got %d entries", len(s.All()))
	}
}

func TestPortCommentStore_SetAndGet(t *testing.T) {
	s := NewPortCommentStore()
	s.Set(8080, "tcp", "dev server")
	c := s.Get(8080, "tcp")
	if c == nil {
		t.Fatal("expected comment, got nil")
	}
	if c.Comment != "dev server" {
		t.Errorf("expected 'dev server', got %q", c.Comment)
	}
	if c.Port != 8080 || c.Proto != "tcp" {
		t.Errorf("unexpected port/proto: %d/%s", c.Port, c.Proto)
	}
}

func TestPortCommentStore_Get_Missing(t *testing.T) {
	s := NewPortCommentStore()
	if got := s.Get(9999, "tcp"); got != nil {
		t.Errorf("expected nil for missing key, got %+v", got)
	}
}

func TestPortCommentStore_Delete(t *testing.T) {
	s := makeCommentStore()
	s.Delete(80, "tcp")
	if s.Get(80, "tcp") != nil {
		t.Error("expected nil after delete")
	}
	if s.Get(443, "tcp") == nil {
		t.Error("expected 443 to still exist")
	}
}

func TestPortCommentStore_All(t *testing.T) {
	s := makeCommentStore()
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(all))
	}
}

func TestSaveLoadPortComments_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "comments.json")

	orig := makeCommentStore()
	if err := SavePortComments(path, orig); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := LoadPortComments(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(loaded.All()) != len(orig.All()) {
		t.Errorf("expected %d comments, got %d", len(orig.All()), len(loaded.All()))
	}
	c := loaded.Get(80, "tcp")
	if c == nil || c.Comment != "HTTP traffic" {
		t.Errorf("unexpected comment for 80/tcp: %+v", c)
	}
}

func TestLoadPortComments_Missing(t *testing.T) {
	s, err := LoadPortComments("/nonexistent/path/comments.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.All()) != 0 {
		t.Error("expected empty store for missing file")
	}
}

func TestLoadPortComments_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0644)
	_, err := LoadPortComments(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
