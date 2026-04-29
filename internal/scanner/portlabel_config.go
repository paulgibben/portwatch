package scanner

import (
	"encoding/json"
	"fmt"
	"os"
)

// PortLabelConfig holds default well-known port labels for bootstrapping.
type PortLabelConfig struct {
	Defaults []PortLabel `json:"defaults"`
}

// DefaultPortLabelConfig returns a config populated with common service labels.
func DefaultPortLabelConfig() PortLabelConfig {
	return PortLabelConfig{
		Defaults: []PortLabel{
			{Port: 21, Label: "FTP", Comment: "File Transfer Protocol"},
			{Port: 22, Label: "SSH", Comment: "Secure Shell"},
			{Port: 25, Label: "SMTP", Comment: "Simple Mail Transfer Protocol"},
			{Port: 53, Label: "DNS", Comment: "Domain Name System"},
			{Port: 80, Label: "HTTP", Comment: "Hypertext Transfer Protocol"},
			{Port: 443, Label: "HTTPS", Comment: "HTTP Secure"},
			{Port: 3306, Label: "MySQL", Comment: "MySQL Database"},
			{Port: 5432, Label: "PostgreSQL", Comment: "PostgreSQL Database"},
			{Port: 6379, Label: "Redis", Comment: "Redis In-Memory Store"},
			{Port: 8080, Label: "HTTP-Alt", Comment: "Alternate HTTP port"},
		},
	}
}

// ApplyDefaults populates a PortLabelStore from a PortLabelConfig without
// overwriting entries that already exist in the store.
func ApplyDefaults(store *PortLabelStore, cfg PortLabelConfig) {
	for _, entry := range cfg.Defaults {
		if _, exists := store.Get(entry.Port); !exists {
			store.Set(entry.Port, entry.Label, entry.Comment)
		}
	}
}

// SavePortLabelConfig writes a PortLabelConfig to a JSON file.
func SavePortLabelConfig(path string, cfg PortLabelConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal port label config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPortLabelConfig reads a PortLabelConfig from a JSON file.
// Returns the default config if the file does not exist.
func LoadPortLabelConfig(path string) (PortLabelConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultPortLabelConfig(), nil
		}
		return PortLabelConfig{}, fmt.Errorf("read port label config: %w", err)
	}
	var cfg PortLabelConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return PortLabelConfig{}, fmt.Errorf("unmarshal port label config: %w", err)
	}
	return cfg, nil
}
