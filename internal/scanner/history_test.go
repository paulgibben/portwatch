package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeDiff(added, removed []int) map[string][]int {
	return map[string][]int{
		"added":   added,
		"removed": removed,
	}
}

func TestHistory_Record_AppendsEntry(t *testing.T) {
	h := NewHistory(0)
	h.Record(makeDiff([]int{8080}, nil))
	if len(h.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h.Entries))
	}
	if len(h.Entries[0].Added) != 1 || h.Entries[0].Added[0] != 8080 {
		t.Errorf("unexpected added ports: %v", h.Entries[0].Added)
	}
}

func TestHistory_Record_RespectsMaxSize(t *testing.T) {
	h := NewHistory(3)
	for i := 0; i < 5; i++ {
		h.Record(makeDiff([]int{i}, nil))
	}
	if len(h.Entries) != 3 {
		t.Fatalf("expected 3 entries after trimming, got %d", len(h.Entries))
	}
	// Should keep the last 3
	if h.Entries[0].Added[0] != 2 {
		t.Errorf("expected first kept entry added=2, got %d", h.Entries[0].Added[0])
	}
}

func TestHistory_Record_SetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	h := NewHistory(0)
	h.Record(makeDiff(nil, []int{443}))
	after := time.Now().UTC()
	ts := h.Entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ts, before, after)
	}
}

func TestSaveLoadHistory_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h := NewHistory(10)
	h.Record(makeDiff([]int{80, 443}, nil))
	h.Record(makeDiff(nil, []int{80}))

	if err := SaveHistory(path, h); err != nil {
		t.Fatalf("SaveHistory: %v", err)
	}

	loaded, err := LoadHistory(path, 10)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded.Entries))
	}
	if loaded.Entries[1].Removed[0] != 80 {
		t.Errorf("unexpected removed port: %v", loaded.Entries[1].Removed)
	}
}

func TestLoadHistory_Missing(t *testing.T) {
	h, err := LoadHistory("/nonexistent/history.json", 0)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(h.Entries))
	}
}

func TestLoadHistory_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0644)
	_, err := LoadHistory(path, 0)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
