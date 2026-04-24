package notifier

import (
	"encoding/json"
	"fmt"
	"os"
)

// Type represents the notifier backend type.
type Type string

const (
	TypeStdout  Type = "stdout"
	TypeCommand Type = "command"
	TypeLog     Type = "log"
)

// Config holds configuration for a notifier.
type Config struct {
	Type    Type   `json:"type"`
	Command string `json:"command,omitempty"`
	LogFile string `json:"log_file,omitempty"`
	Format  string `json:"format,omitempty"`
}

// Defaults applies default values to a Config.
func (c *Config) Defaults() {
	if c.Type == "" {
		c.Type = TypeStdout
	}
	if c.Format == "" {
		c.Format = "text"
	}
}

// Validate checks that the Config is valid for its type.
func (c *Config) Validate() error {
	switch c.Type {
	case TypeStdout:
		// no extra fields required
	case TypeCommand:
		if c.Command == "" {
			return fmt.Errorf("notifier: command type requires a non-empty command")
		}
	case TypeLog:
		if c.LogFile == "" {
			return fmt.Errorf("notifier: log type requires a non-empty log_file")
		}
	default:
		return fmt.Errorf("notifier: unknown type %q", c.Type)
	}
	return nil
}

// LoadConfig reads a notifier Config from a JSON file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("notifier: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("notifier: parse config: %w", err)
	}
	cfg.Defaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig writes a Config to a JSON file.
func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("notifier: marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("notifier: write config: %w", err)
	}
	return nil
}
