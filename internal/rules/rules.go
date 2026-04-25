package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// PortRange defines an inclusive range of port numbers to monitor.
type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Config holds the rules configuration for portwatch.
type Config struct {
	PortRanges  []PortRange `json:"port_ranges"`
	IgnorePorts []int       `json:"ignore_ports"`
	AlertOnNew  bool        `json:"alert_on_new"`
	AlertOnGone bool        `json:"alert_on_gone"`
}

// Evaluate checks a port against the config and returns whether it should
// trigger an alert, and whether it should be ignored entirely.
func (c *Config) Evaluate(port int) (alert bool, ignore bool) {
	for _, ignored := range c.IgnorePorts {
		if ignored == port {
			return false, true
		}
	}
	if len(c.PortRanges) == 0 {
		return true, false
	}
	for _, r := range c.PortRanges {
		if port >= r.Start && port <= r.End {
			return true, false
		}
	}
	return false, false
}

// LoadConfig reads and parses a JSON rules config from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config file not found: %s", path)
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// SaveConfig serializes the config to JSON and writes it to the given path.
func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}
