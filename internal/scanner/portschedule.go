package scanner

import (
	"encoding/json"
	"os"
	"time"
)

// ScanScheduleConfig defines when and how often port scanning should occur.
type ScanScheduleConfig struct {
	Interval    time.Duration `json:"interval"`
	StartDelay  time.Duration `json:"start_delay"`
	MaxMissed   int           `json:"max_missed"`
	Enabled     bool          `json:"enabled"`
}

// DefaultScanScheduleConfig returns a sensible default schedule.
func DefaultScanScheduleConfig() ScanScheduleConfig {
	return ScanScheduleConfig{
		Interval:   60 * time.Second,
		StartDelay: 0,
		MaxMissed:  3,
		Enabled:    true,
	}
}

// SaveScanScheduleConfig writes the schedule config to a JSON file.
func SaveScanScheduleConfig(path string, cfg ScanScheduleConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadScanScheduleConfig reads a schedule config from a JSON file.
// If the file does not exist, returns the default config.
func LoadScanScheduleConfig(path string) (ScanScheduleConfig, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DefaultScanScheduleConfig(), nil
	}
	if err != nil {
		return ScanScheduleConfig{}, err
	}
	var cfg ScanScheduleConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return ScanScheduleConfig{}, err
	}
	return cfg, nil
}
