package scanner

import (
	"net"
	"testing"
	"time"
)

func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScanPort_Open(t *testing.T) {
	port, stop := startTestServer(t)
	defer stop()

	s := New("127.0.0.1", time.Second)
	state := s.ScanPort(port)

	if !state.Open {
		t.Errorf("expected port %d to be open", port)
	}
	if state.Port != port {
		t.Errorf("expected port %d, got %d", port, state.Port)
	}
	if state.Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", state.Protocol)
	}
}

func TestScanPort_Closed(t *testing.T) {
	s := New("127.0.0.1", 100*time.Millisecond)
	// Port 1 is almost certainly closed in test environments
	state := s.ScanPort(1)
	if state.Open {
		t.Skip("port 1 unexpectedly open, skipping test")
	}
	if state.Port != 1 {
		t.Errorf("expected port 1, got %d", state.Port)
	}
}

func TestOpenPorts_Filter(t *testing.T) {
	states := []PortState{
		{Port: 80, Open: true},
		{Port: 81, Open: false},
		{Port: 443, Open: true},
		{Port: 8080, Open: false},
	}
	open := OpenPorts(states)
	if len(open) != 2 {
		t.Errorf("expected 2 open ports, got %d", len(open))
	}
	for _, s := range open {
		if !s.Open {
			t.Errorf("expected open port, got closed port %d", s.Port)
		}
	}
}
