package rules

import (
	"encoding/json"
	"fmt"
	"os"
)

// Action defines what to do when a rule matches.
type Action string

const (
	ActionAlert  Action = "alert"
	ActionIgnore Action = "ignore"
)

// Rule defines a configurable rule for port monitoring.
type Rule struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"` // "tcp" or "udp"
	Action   Action `json:"action"`
	Comment  string `json:"comment,omitempty"`
}

// Config holds the full rules configuration.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON rules file from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading rules file: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing rules file: %w", err)
	}
	return &cfg, nil
}

// Evaluate returns the Action for a given port/protocol pair.
// If no rule matches, ActionAlert is returned by default.
func (c *Config) Evaluate(port int, protocol string) Action {
	for _, r := range c.Rules {
		if r.Port == port && r.Protocol == protocol {
			return r.Action
		}
	}
	return ActionAlert
}

// SaveConfig writes the config as JSON to the given path.
func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling rules: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing rules file: %w", err)
	}
	return nil
}
