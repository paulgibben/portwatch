package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func makePortGroupStore() *PortGroupStore {
	s := NewPortGroupStore()
	s.Add(PortGroup{Name: "web", Description: "Web ports", Ports: []int{80, 443}})
	s.Add(PortGroup{Name: "db", Description: "Database ports", Ports: []int{5432, 3306}})
	return s
}

func TestNewPortGroupStore_Empty(t *testing.T) {
	s := NewPortGroupStore()
	if len(s.Groups) != 0 {
		t.Fatalf("expected empty store, got %d groups", len(s.Groups))
	}
}

func TestPortGroupStore_Add_And_Get(t *testing.T) {
	s := makePortGroupStore()
	g, ok := s.Get("web")
	if !ok {
		t.Fatal("expected to find group 'web'")
	}
	if len(g.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(g.Ports))
	}
}

func TestPortGroupStore_Add_Replaces(t *testing.T) {
	s := makePortGroupStore()
	s.Add(PortGroup{Name: "web", Ports: []int{8080}})
	g, _ := s.Get("web")
	if len(g.Ports) != 1 || g.Ports[0] != 8080 {
		t.Fatalf("expected replaced group with port 8080, got %v", g.Ports)
	}
	if len(s.Groups) != 2 {
		t.Fatalf("expected 2 groups after replace, got %d", len(s.Groups))
	}
}

func TestPortGroupStore_Get_Missing(t *testing.T) {
	s := makePortGroupStore()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestPortGroupStore_AllPorts_Deduplicated(t *testing.T) {
	s := NewPortGroupStore()
	s.Add(PortGroup{Name: "a", Ports: []int{80, 443, 8080}})
	s.Add(PortGroup{Name: "b", Ports: []int{443, 5432}})
	ports := s.AllPorts()
	expected := []int{80, 443, 5432, 8080}
	if len(ports) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, ports)
	}
	for i, p := range ports {
		if p != expected[i] {
			t.Fatalf("expected %v, got %v", expected, ports)
		}
	}
}

func TestSaveLoadPortGroups_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "groups.json")
	s := makePortGroupStore()
	if err := SavePortGroups(path, s); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadPortGroups(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Groups) != len(s.Groups) {
		t.Fatalf("expected %d groups, got %d", len(s.Groups), len(loaded.Groups))
	}
}

func TestLoadPortGroups_Missing(t *testing.T) {
	store, err := LoadPortGroups("/nonexistent/path/groups.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(store.Groups) != 0 {
		t.Fatal("expected empty store")
	}
}

func TestLoadPortGroups_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := LoadPortGroups(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
