package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// probeConfigFile is the serialisable form of ProbeConfig.
type probeConfigFile struct {
	TimeoutMS  int64 `json:"timeout_ms"`
	Retries    int   `json:"retries"`
	BannerGrab bool  `json:"banner_grab"`
}

// SaveProbeConfig writes a ProbeConfig to a JSON file at path.
func SaveProbeConfig(path string, cfg ProbeConfig) error {
	pf := probeConfigFile{
		TimeoutMS:  cfg.Timeout.Milliseconds(),
		Retries:    cfg.Retries,
		BannerGrab: cfg.BannerGrab,
	}
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal probe config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write probe config: %w", err)
	}
	return nil
}

// LoadProbeConfig reads a ProbeConfig from a JSON file at path.
// If the file does not exist, DefaultProbeConfig is returned.
func LoadProbeConfig(path string) (ProbeConfig, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DefaultProbeConfig(), nil
	}
	if err != nil {
		return ProbeConfig{}, fmt.Errorf("read probe config: %w", err)
	}

	var pf probeConfigFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return ProbeConfig{}, fmt.Errorf("parse probe config: %w", err)
	}

	return ProbeConfig{
		Timeout:    time.Duration(pf.TimeoutMS) * time.Millisecond,
		Retries:    pf.Retries,
		BannerGrab: pf.BannerGrab,
	}, nil
}
