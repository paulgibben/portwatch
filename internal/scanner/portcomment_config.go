package scanner

import (
	"encoding/json"
	"os"
)

// PortCommentConfig holds the file path configuration for persisting port comments.
type PortCommentConfig struct {
	FilePath string `json:"file_path"`
}

// DefaultPortCommentConfig returns a PortCommentConfig with sensible defaults.
func DefaultPortCommentConfig() PortCommentConfig {
	return PortCommentConfig{
		FilePath: "portcomments.json",
	}
}

// SavePortCommentConfig writes the config to a JSON file.
func SavePortCommentConfig(path string, cfg PortCommentConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadPortCommentConfig reads a PortCommentConfig from a JSON file.
// Returns the default config if the file does not exist.
func LoadPortCommentConfig(path string) (PortCommentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultPortCommentConfig(), nil
		}
		return PortCommentConfig{}, err
	}
	var cfg PortCommentConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return PortCommentConfig{}, err
	}
	if cfg.FilePath == "" {
		cfg.FilePath = DefaultPortCommentConfig().FilePath
	}
	return cfg, nil
}
