package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeDiffForChange() Diff {
	return Diff{
		Added:   []Port{{Port: 8080, Proto: "tcp"}, {Port: 443, Proto: "tcp"}},
		Removed: []Port{{Port: 22, Proto: "tcp"}},
	}
}

func TestPortChangeLog_Record(t *testing.T) {
	log := NewPortChangeLog()
	before := time.Now().UTC()
	log.Record(80, "tcp", ChangeAdded)
	after := time.Now().UTC()

	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Port != 80 || e.Proto != "tcp" || e.Change != ChangeAdded {
		t.Errorf("unexpected entry: %+v", e)
	}
	if e.DetectedAt.Before(before) || e.DetectedAt.After(after) {
		t.Errorf("timestamp out of range: %v", e.DetectedAt)
	}
}

func TestSaveLoadPortChangeLog_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "changes.json")

	log := NewPortChangeLog()
	log.Record(9000, "tcp", ChangeAdded)
	log.Record(22, "tcp", ChangeRemoved)

	if err := SavePortChangeLog(path, log); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadPortChangeLog(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded.Entries))
	}
}

func TestLoadPortChangeLog_Missing(t *testing.T) {
	log, err := LoadPortChangeLog(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(log.Entries) != 0 {
		t.Errorf("expected empty log")
	}
}

func TestLoadPortChangeLog_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not-json"), 0644)
	_, err := LoadPortChangeLog(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestPortChangeTracker_Apply(t *testing.T) {
	tracker := NewPortChangeTracker(nil)
	tracker.Apply(makeDiffForChange())

	if tracker.CountByType(ChangeAdded) != 2 {
		t.Errorf("expected 2 added, got %d", tracker.CountByType(ChangeAdded))
	}
	if tracker.CountByType(ChangeRemoved) != 1 {
		t.Errorf("expected 1 removed, got %d", tracker.CountByType(ChangeRemoved))
	}
}

func TestPortChangeTracker_Since(t *testing.T) {
	tracker := NewPortChangeTracker(nil)
	tracker.Apply(makeDiffForChange())

	past := time.Now().Add(-time.Hour).Unix()
	future := time.Now().Add(time.Hour).Unix()

	if len(tracker.Since(past)) != 3 {
		t.Errorf("expected 3 entries since past")
	}
	if len(tracker.Since(future)) != 0 {
		t.Errorf("expected 0 entries since future")
	}
}
