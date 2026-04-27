package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// throttleConfigFile is the on-disk representation for JSON marshalling.
type throttleConfigFile struct {
	MinIntervalSeconds int `json:"min_interval_seconds"`
	BurstAllowance     int `json:"burst_allowance"`
}

// SaveThrottleConfig persists a ThrottleConfig to the given file path as JSON.
func SaveThrottleConfig(path string, cfg ThrottleConfig) error {
	fc := throttleConfigFile{
		MinIntervalSeconds: int(cfg.MinInterval.Seconds()),
		BurstAllowance:     cfg.BurstAllowance,
	}
	data, err := json.MarshalIndent(fc, "", "  ")
	if err != nil {
		return fmt.Errorf("throttle: marshal config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("throttle: write config: %w", err)
	}
	return nil
}

// LoadThrottleConfig reads a ThrottleConfig from the given file path.
// If the file does not exist, DefaultThrottleConfig is returned.
func LoadThrottleConfig(path string) (ThrottleConfig, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DefaultThrottleConfig(), nil
	}
	if err != nil {
		return ThrottleConfig{}, fmt.Errorf("throttle: read config: %w", err)
	}

	var fc throttleConfigFile
	if err := json.Unmarshal(data, &fc); err != nil {
		return ThrottleConfig{}, fmt.Errorf("throttle: parse config: %w", err)
	}

	cfg := ThrottleConfig{
		MinInterval:    time.Duration(fc.MinIntervalSeconds) * time.Second,
		BurstAllowance: fc.BurstAllowance,
	}
	// Apply defaults for zero values.
	def := DefaultThrottleConfig()
	if cfg.MinInterval <= 0 {
		cfg.MinInterval = def.MinInterval
	}
	if cfg.BurstAllowance <= 0 {
		cfg.BurstAllowance = def.BurstAllowance
	}
	return cfg, nil
}
