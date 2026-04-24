package notifier

import (
	"encoding/json"
	"fmt"
	"os"
)

// Type represents the kind of notifier backend.
type Type string

const (
	TypeStdout  Type = "stdout"
	TypeCommand Type = "command"
	TypeLog     Type = "log"
	TypeWebhook Type = "webhook"
	TypeEmail   Type = "email"
)

// Config holds the top-level notifier configuration.
type Config struct {
	Type        Type        `json:"type"`
	CommandPath string      `json:"command,omitempty"`
	LogFile     string      `json:"log_file,omitempty"`
	WebhookURL  string      `json:"webhook_url,omitempty"`
	Email       EmailConfig `json:"email,omitempty"`
}

// Defaults fills in zero-value fields with sensible defaults.
func (c *Config) Defaults() {
	if c.Type == "" {
		c.Type = TypeStdout
	}
}

// Validate returns an error if the configuration is incomplete.
func (c *Config) Validate() error {
	switch c.Type {
	case TypeStdout:
		// no extra fields required
	case TypeCommand:
		if c.CommandPath == "" {
			return fmt.Errorf("notifier: command path is required for type %q", c.Type)
		}
	case TypeLog:
		if c.LogFile == "" {
			return fmt.Errorf("notifier: log_file is required for type %q", c.Type)
		}
	case TypeWebhook:
		if c.WebhookURL == "" {
			return fmt.Errorf("notifier: webhook_url is required for type %q", c.Type)
		}
	case TypeEmail:
		if c.Email.SMTPHost == "" || c.Email.From == "" || len(c.Email.To) == 0 {
			return fmt.Errorf("notifier: email config is incomplete for type %q", c.Type)
		}
	default:
		return fmt.Errorf("notifier: unknown type %q", c.Type)
	}
	return nil
}

// LoadConfig reads a Config from a JSON file at path.
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
	return &cfg, nil
}

// SaveConfig writes cfg to a JSON file at path.
func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("notifier: marshal config: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}
