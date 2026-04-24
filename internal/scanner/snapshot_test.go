package scanner

import (
	"os"
	"testing"
	"time"
)

func makeSnapshot(host string, ports []int) *Snapshot {
	states := make([]PortState, len(ports))
	for i, p := range ports {
		states[i] = PortState{Port: p, Protocol: "tcp", Open: true, Address: host}
	}
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Ports:     states,
	}
}

func TestSnapshotSaveLoad(t *testing.T) {
	snap := makeSnapshot("127.0.0.1", []int{80, 443, 8080})

	tmp, err := os.CreateTemp("", "portwatch-snap-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := snap.SaveToFile(tmp.Name()); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	loaded, err := LoadFromFile(tmp.Name())
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}
	if loaded.Host != snap.Host {
		t.Errorf("host mismatch: got %s, want %s", loaded.Host, snap.Host)
	}
	if len(loaded.Ports) != len(snap.Ports) {
		t.Errorf("port count mismatch: got %d, want %d", len(loaded.Ports), len(snap.Ports))
	}
}

func TestDiff(t *testing.T) {
	prev := makeSnapshot("127.0.0.1", []int{80, 443})
	curr := makeSnapshot("127.0.0.1", []int{443, 8080})

	opened, closed := Diff(prev, curr)

	if len(opened) != 1 || opened[0].Port != 8080 {
		t.Errorf("expected opened port 8080, got %+v", opened)
	}
	if len(closed) != 1 || closed[0].Port != 80 {
		t.Errorf("expected closed port 80, got %+v", closed)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	prev := makeSnapshot("127.0.0.1", []int{80, 443})
	curr := makeSnapshot("127.0.0.1", []int{80, 443})
	opened, closed := Diff(prev, curr)
	if len(opened) != 0 || len(closed) != 0 {
		t.Errorf("expected no diff, got opened=%v closed=%v", opened, closed)
	}
}
