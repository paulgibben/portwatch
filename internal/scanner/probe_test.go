package scanner

import (
	"net"
	"testing"
	"time"
)

func startProbeServer(t *testing.T, banner string) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			if banner != "" {
				_, _ = conn.Write([]byte(banner))
			}
			conn.Close()
		}
	}()

	return ln.Addr().(*net.TCPAddr).Port
}

func TestDefaultProbeConfig(t *testing.T) {
	cfg := DefaultProbeConfig()
	if cfg.Timeout != 2*time.Second {
		t.Errorf("expected 2s timeout, got %v", cfg.Timeout)
	}
	if cfg.Retries != 1 {
		t.Errorf("expected 1 retry, got %d", cfg.Retries)
	}
	if cfg.BannerGrab {
		t.Error("expected BannerGrab to be false by default")
	}
}

func TestProbe_OpenPort(t *testing.T) {
	port := startProbeServer(t, "")
	p := NewProber(DefaultProbeConfig())
	result := p.Probe(port)

	if !result.Open {
		t.Errorf("expected port %d to be open", port)
	}
	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}
}

func TestProbe_ClosedPort(t *testing.T) {
	// Port 1 is almost certainly closed and unbound in test envs.
	cfg := DefaultProbeConfig()
	cfg.Timeout = 200 * time.Millisecond
	cfg.Retries = 0
	p := NewProber(cfg)
	result := p.Probe(1)

	if result.Open {
		t.Error("expected port 1 to be closed")
	}
}

func TestProbe_BannerGrab(t *testing.T) {
	expected := "HELLO portwatch"
	port := startProbeServer(t, expected)

	cfg := DefaultProbeConfig()
	cfg.BannerGrab = true
	p := NewProber(cfg)
	result := p.Probe(port)

	if !result.Open {
		t.Fatalf("expected port %d to be open", port)
	}
	if result.Banner != expected {
		t.Errorf("expected banner %q, got %q", expected, result.Banner)
	}
}

func TestProbe_NoBannerWhenDisabled(t *testing.T) {
	port := startProbeServer(t, "SHOULD NOT READ")

	cfg := DefaultProbeConfig()
	cfg.BannerGrab = false
	p := NewProber(cfg)
	result := p.Probe(port)

	if result.Banner != "" {
		t.Errorf("expected empty banner, got %q", result.Banner)
	}
}
