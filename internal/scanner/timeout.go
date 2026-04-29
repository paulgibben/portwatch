package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TimeoutConfig controls per-port and global scan timeouts.
type TimeoutConfig struct {
	// PortTimeout is the maximum duration to wait when probing a single port.
	PortTimeout time.Duration `json:"port_timeout_ms"`
	// GlobalTimeout is the maximum total duration for a full scan.
	GlobalTimeout time.Duration `json:"global_timeout_ms"`
}

// DefaultTimeoutConfig returns sensible defaults.
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		PortTimeout:   500 * time.Millisecond,
		GlobalTimeout: 30 * time.Second,
	}
}

type timeoutConfigJSON struct {
	PortTimeoutMS   int64 `json:"port_timeout_ms"`
	GlobalTimeoutMS int64 `json:"global_timeout_ms"`
}

// SaveTimeoutConfig writes a TimeoutConfig to a JSON file.
func SaveTimeoutConfig(path string, cfg TimeoutConfig) error {
	raw := timeoutConfigJSON{
		PortTimeoutMS:   cfg.PortTimeout.Milliseconds(),
		GlobalTimeoutMS: cfg.GlobalTimeout.Milliseconds(),
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal timeout config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write timeout config: %w", err)
	}
	return nil
}

// LoadTimeoutConfig reads a TimeoutConfig from a JSON file.
// If the file does not exist, DefaultTimeoutConfig is returned.
func LoadTimeoutConfig(path string) (TimeoutConfig, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DefaultTimeoutConfig(), nil
	}
	if err != nil {
		return TimeoutConfig{}, fmt.Errorf("read timeout config: %w", err)
	}
	var raw timeoutConfigJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return TimeoutConfig{}, fmt.Errorf("parse timeout config: %w", err)
	}
	return TimeoutConfig{
		PortTimeout:   time.Duration(raw.PortTimeoutMS) * time.Millisecond,
		GlobalTimeout: time.Duration(raw.GlobalTimeoutMS) * time.Millisecond,
	}, nil
}
