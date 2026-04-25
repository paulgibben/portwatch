package scanner

import (
	"testing"

	"github.com/user/portwatch/internal/rules"
)

func makeFilterConfig(ignorePorts []int, ranges []rules.PortRange) *rules.Config {
	return &rules.Config{
		IgnorePorts: ignorePorts,
		PortRanges:  ranges,
	}
}

func TestFilter_Apply_NoConfig(t *testing.T) {
	f := NewFilter(nil)
	ports := []int{22, 80, 443, 8080}
	result := f.Apply(ports)
	if len(result) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(result))
	}
}

func TestFilter_Apply_IgnorePorts(t *testing.T) {
	cfg := makeFilterConfig([]int{22, 8080}, nil)
	f := NewFilter(cfg)
	ports := []int{22, 80, 443, 8080}
	result := f.Apply(ports)
	expected := []int{80, 443}
	if len(result) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
	for i, p := range result {
		if p != expected[i] {
			t.Errorf("expected port %d at index %d, got %d", expected[i], i, p)
		}
	}
}

func TestFilter_Apply_PortRange(t *testing.T) {
	cfg := makeFilterConfig(nil, []rules.PortRange{{Start: 80, End: 443}})
	f := NewFilter(cfg)
	ports := []int{22, 80, 443, 8080}
	result := f.Apply(ports)
	expected := []int{80, 443}
	if len(result) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestFilter_Apply_IgnoreAndRange(t *testing.T) {
	cfg := makeFilterConfig([]int{80}, []rules.PortRange{{Start: 1, End: 1024}})
	f := NewFilter(cfg)
	ports := []int{22, 80, 443, 8080}
	result := f.Apply(ports)
	expected := []int{22, 443}
	if len(result) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
	for i, p := range result {
		if p != expected[i] {
			t.Errorf("expected port %d at index %d, got %d", expected[i], i, p)
		}
	}
}

func TestFilter_Apply_EmptyPorts(t *testing.T) {
	cfg := makeFilterConfig([]int{22}, []rules.PortRange{{Start: 1, End: 1024}})
	f := NewFilter(cfg)
	result := f.Apply([]int{})
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}
