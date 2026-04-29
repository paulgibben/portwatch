package scanner

import (
	"encoding/json"
	"fmt"
	"os"
)

// PortGroupConfig holds configuration for port group management,
// including the file path used to persist groups.
type PortGroupConfig struct {
	StorePath    string `json:"store_path"`
	AutoMerge    bool   `json:"auto_merge"`
	MaxGroupSize int    `json:"max_group_size"`
}

// DefaultPortGroupConfig returns a PortGroupConfig with sensible defaults.
func DefaultPortGroupConfig() PortGroupConfig {
	return PortGroupConfig{
		StorePath:    "portgroups.json",
		AutoMerge:    false,
		MaxGroupSize: 256,
	}
}

// Validate checks that the config is valid.
func (c PortGroupConfig) Validate() error {
	if c.StorePath == "" {
		return fmt.Errorf("portgroup config: store_path must not be empty")
	}
	if c.MaxGroupSize <= 0 {
		return fmt.Errorf("portgroup config: max_group_size must be positive")
	}
	return nil
}

// SavePortGroupConfig writes the config to a JSON file.
func SavePortGroupConfig(path string, cfg PortGroupConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("portgroup config: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPortGroupConfig reads a PortGroupConfig from a JSON file.
// Returns defaults if the file does not exist.
func LoadPortGroupConfig(path string) (PortGroupConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultPortGroupConfig(), nil
		}
		return PortGroupConfig{}, fmt.Errorf("portgroup config: read: %w", err)
	}
	var cfg PortGroupConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return PortGroupConfig{}, fmt.Errorf("portgroup config: unmarshal: %w", err)
	}
	return cfg, nil
}
