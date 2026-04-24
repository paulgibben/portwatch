package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a scanned port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
	Address  string
}

// Scanner scans local open ports.
type Scanner struct {
	Host    string
	Timeout time.Duration
}

// New creates a new Scanner with the given host and timeout.
func New(host string, timeout time.Duration) *Scanner {
	return &Scanner{
		Host:    host,
		Timeout: timeout,
	}
}

// ScanPort checks whether a single TCP port is open.
func (s *Scanner) ScanPort(port int) PortState {
	address := fmt.Sprintf("%s:%d", s.Host, port)
	conn, err := net.DialTimeout("tcp", address, s.Timeout)
	state := PortState{
		Port:     port,
		Protocol: "tcp",
		Address:  address,
	}
	if err != nil {
		state.Open = false
		return state
	}
	conn.Close()
	state.Open = true
	return state
}

// ScanRange scans a range of ports [start, end] inclusive.
func (s *Scanner) ScanRange(start, end int) []PortState {
	results := make([]PortState, 0, end-start+1)
	for port := start; port <= end; port++ {
		results = append(results, s.ScanPort(port))
	}
	return results
}

// OpenPorts filters and returns only open ports from a scan result.
func OpenPorts(states []PortState) []PortState {
	open := make([]PortState, 0)
	for _, s := range states {
		if s.Open {
			open = append(open, s)
		}
	}
	return open
}
