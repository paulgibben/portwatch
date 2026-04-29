package scanner

import (
	"encoding/json"
	"fmt"
	"os"
)

// PortTaggingConfig holds default tags applied to well-known ports on startup.
type PortTaggingConfig struct {
	// DefaultTags maps "port/proto" to a list of tags to apply automatically.
	DefaultTags map[string][]string `json:"default_tags,omitempty"`
	// PersistPath is the file path used to save/load the live tag store.
	PersistPath string `json:"persist_path,omitempty"`
}

// DefaultPortTaggingConfig returns a sensible default configuration.
func DefaultPortTaggingConfig() PortTaggingConfig {
	return PortTaggingConfig{
		DefaultTags: map[string][]string{
			"22/tcp":  {"ssh"},
			"80/tcp":  {"http", "web"},
			"443/tcp": {"https", "web", "tls"},
		},
		PersistPath: "portwatch_tags.json",
	}
}

// SavePortTaggingConfig writes the config to a JSON file.
func SavePortTaggingConfig(path string, cfg PortTaggingConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("porttagging_config: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadPortTaggingConfig reads a PortTaggingConfig from a JSON file.
// If the file does not exist, the default config is returned.
func LoadPortTaggingConfig(path string) (PortTaggingConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultPortTaggingConfig(), nil
		}
		return PortTaggingConfig{}, fmt.Errorf("porttagging_config: read: %w", err)
	}
	var cfg PortTaggingConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return PortTaggingConfig{}, fmt.Errorf("porttagging_config: unmarshal: %w", err)
	}
	return cfg, nil
}

// ApplyDefaultTags seeds a PortTagStore with the tags defined in the config.
func ApplyDefaultTags(store *PortTagStore, cfg PortTaggingConfig) {
	for key, tags := range cfg.DefaultTags {
		var port int
		var proto string
		fmt.Sscanf(key, "%d/%s", &port, &proto)
		if port == 0 || proto == "" {
			continue
		}
		for _, tag := range tags {
			store.AddTag(port, proto, tag)
		}
	}
}
