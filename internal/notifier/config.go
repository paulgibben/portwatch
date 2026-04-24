package notifier

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Type identifies the notification backend.
type Type string

const (
	TypeStdout  Type = "stdout"
	TypeCommand Type = "command"
	TypeLog     Type = "log"
	TypeWebhook Type = "webhook"
)

// Config holds notifier configuration.
type Config struct {
	Type           Type          `json:"type"`
	Command        string        `json:"command,omitempty"`
	LogFile        string        `json:"log_file,omitempty"`
	WebhookURL     string        `json:"webhook_url,omitempty"`
	WebhookTimeout time.Duration `json:"webhook_timeout_ms,omitempty"`
}

// Defaults fills in zero-value fields with sensible defaults.
func (c *Config) Defaults() {
	if c.Type == "" {
		c.Type = TypeStdout
	}
	if c.Type == TypeWebhook && c.WebhookTimeout == 0 {
		c.WebhookTimeout = 10 * time.Second
	}
}

// Validate returns an error if the config is invalid for its type.
func (c *Config) Validate() error {
	switch c.Type {
	case TypeStdout:
		return nil
	case TypeCommand:
		if c.Command == "" {
			return fmt.Errorf("notifier: command type requires 'command' field")
		}
		return nil
	case TypeLog:
		if c.LogFile == "" {
			return fmt.Errorf("notifier: log type requires 'log_file' field")
		}
		return nil
	case TypeWebhook:
		if c.WebhookURL == "" {
			return fmt.Errorf("notifier: webhook type requires 'webhook_url' field")
		}
		return nil
	default:
		return fmt.Errorf("notifier: unknown type %q", c.Type)
	}
}

// LoadConfig reads a Config from a JSON file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("notifier: load config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("notifier: parse config: %w", err)
	}
	cfg.Defaults()
	return &cfg, nil
}

// SaveConfig writes a Config to a JSON file.
func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("notifier: save config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}
