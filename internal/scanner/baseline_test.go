package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeBaseline(host string, ports []int) *Baseline {
	return &Baseline{
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Host:      host,
		Ports:     ports,
	}
}

func TestNewBaseline_CopiesPorts(t *testing.T) {
	original := []int{80, 443, 8080}
	b := NewBaseline("localhost", original)
	original[0] = 9999 // mutate original slice
	if b.Ports[0] == 9999 {
		t.Error("NewBaseline should copy ports slice, not reference it")
	}
	if b.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", b.Host)
	}
}

func TestSaveLoadBaseline_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	orig := makeBaseline("192.168.1.1", []int{22, 80, 443})
	if err := SaveBaseline(orig, path); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	loaded, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	if loaded.Host != orig.Host {
		t.Errorf("host mismatch: got %s, want %s", loaded.Host, orig.Host)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("ports length mismatch: got %d, want %d", len(loaded.Ports), len(orig.Ports))
	}
}

func TestLoadBaseline_Missing(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/path/baseline.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadBaseline(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestBaseline_Unexpected(t *testing.T) {
	b := makeBaseline("localhost", []int{80, 443})

	current := []int{80, 443, 8080, 9090}
	unexpected := b.Unexpected(current)

	if len(unexpected) != 2 {
		t.Fatalf("expected 2 unexpected ports, got %d: %v", len(unexpected), unexpected)
	}
	for _, p := range unexpected {
		if p != 8080 && p != 9090 {
			t.Errorf("unexpected port %d not in expected set", p)
		}
	}
}

func TestBaseline_Unexpected_NoneNew(t *testing.T) {
	b := makeBaseline("localhost", []int{80, 443})
	unexpected := b.Unexpected([]int{80, 443})
	if len(unexpected) != 0 {
		t.Errorf("expected no unexpected ports, got %v", unexpected)
	}
}
