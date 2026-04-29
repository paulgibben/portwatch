package scanner

import (
	"encoding/json"
	"os"
)

// SeverityLevel represents the importance of a port change event.
type SeverityLevel string

const (
	SeverityLow      SeverityLevel = "low"
	SeverityMedium   SeverityLevel = "medium"
	SeverityHigh     SeverityLevel = "high"
	SeverityCritical SeverityLevel = "critical"
)

// PortSeverityRule maps a port (and optional protocol) to a severity level.
type PortSeverityRule struct {
	Port     int           `json:"port"`
	Protocol string        `json:"protocol,omitempty"`
	Severity SeverityLevel `json:"severity"`
}

// PortSeverityConfig holds all severity rules and a default level.
type PortSeverityConfig struct {
	DefaultSeverity SeverityLevel      `json:"default_severity"`
	Rules           []PortSeverityRule  `json:"rules"`
}

// DefaultPortSeverityConfig returns a sensible default configuration.
func DefaultPortSeverityConfig() PortSeverityConfig {
	return PortSeverityConfig{
		DefaultSeverity: SeverityLow,
		Rules: []PortSeverityRule{
			{Port: 22, Protocol: "tcp", Severity: SeverityHigh},
			{Port: 443, Protocol: "tcp", Severity: SeverityMedium},
			{Port: 3306, Protocol: "tcp", Severity: SeverityCritical},
			{Port: 5432, Protocol: "tcp", Severity: SeverityCritical},
		},
	}
}

// SavePortSeverityConfig writes the config to a JSON file.
func SavePortSeverityConfig(path string, cfg PortSeverityConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPortSeverityConfig reads the config from a JSON file.
func LoadPortSeverityConfig(path string) (PortSeverityConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultPortSeverityConfig(), nil
		}
		return PortSeverityConfig{}, err
	}
	var cfg PortSeverityConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return PortSeverityConfig{}, err
	}
	return cfg, nil
}
