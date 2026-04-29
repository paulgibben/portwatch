package scanner

import (
	"fmt"
	"net"
	"time"
)

// ProbeConfig holds configuration for port probing behavior.
type ProbeConfig struct {
	Timeout    time.Duration `json:"timeout"`
	Retries    int           `json:"retries"`
	BannerGrab bool          `json:"banner_grab"`
}

// DefaultProbeConfig returns sensible defaults for probing.
func DefaultProbeConfig() ProbeConfig {
	return ProbeConfig{
		Timeout:    2 * time.Second,
		Retries:    1,
		BannerGrab: false,
	}
}

// ProbeResult holds the result of probing a single port.
type ProbeResult struct {
	Port   int
	Open   bool
	Banner string
	Error  error
}

// Prober probes individual ports for liveness and optional banner.
type Prober struct {
	cfg ProbeConfig
}

// NewProber creates a Prober with the given config.
func NewProber(cfg ProbeConfig) *Prober {
	return &Prober{cfg: cfg}
}

// Probe checks whether the given TCP port is open on localhost.
// If BannerGrab is enabled, it attempts to read an initial banner.
func (p *Prober) Probe(port int) ProbeResult {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	var lastErr error

	for attempt := 0; attempt <= p.cfg.Retries; attempt++ {
		conn, err := net.DialTimeout("tcp", addr, p.cfg.Timeout)
		if err != nil {
			lastErr = err
			continue
		}
		defer conn.Close()

		result := ProbeResult{Port: port, Open: true}

		if p.cfg.BannerGrab {
			_ = conn.SetReadDeadline(time.Now().Add(p.cfg.Timeout))
			buf := make([]byte, 256)
			n, _ := conn.Read(buf)
			if n > 0 {
				result.Banner = string(buf[:n])
			}
		}

		return result
	}

	return ProbeResult{Port: port, Open: false, Error: lastErr}
}
