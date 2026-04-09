package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadAll(configDir string) (*AllConfig, error) {
	cfg := &AllConfig{}

	if err := loadJSON(filepath.Join(configDir, "sources.json"), &cfg.Sources); err != nil {
		return nil, fmt.Errorf("load sources.json: %w", err)
	}
	if err := loadJSON(filepath.Join(configDir, "components.json"), &cfg.Components); err != nil {
		return nil, fmt.Errorf("load components.json: %w", err)
	}
	if err := loadJSON(filepath.Join(configDir, "screens.json"), &cfg.Screens); err != nil {
		return nil, fmt.Errorf("load screens.json: %w", err)
	}
	if err := loadJSON(filepath.Join(configDir, "devices.json"), &cfg.Devices); err != nil {
		return nil, fmt.Errorf("load devices.json: %w", err)
	}
	if err := loadJSON(filepath.Join(configDir, "device_types.json"), &cfg.DeviceTypes); err != nil {
		return nil, fmt.Errorf("load device_types.json: %w", err)
	}
	if err := loadJSON(filepath.Join(configDir, "panel_devices.json"), &cfg.PanelDevices); err != nil {
		return nil, fmt.Errorf("load panel_devices.json: %w", err)
	}
	if err := loadJSON(filepath.Join(configDir, "panel_profiles.json"), &cfg.PanelProfiles); err != nil {
		return nil, fmt.Errorf("load panel_profiles.json: %w", err)
	}
	if err := loadJSON(filepath.Join(configDir, "panel_commands.json"), &cfg.PanelCommands); err != nil {
		return nil, fmt.Errorf("load panel_commands.json: %w", err)
	}
	if err := loadJSON(filepath.Join(configDir, "themes.json"), &cfg.Themes); err != nil {
		return nil, fmt.Errorf("load themes.json: %w", err)
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadJSON(path string, dst interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, dst); err != nil {
		return err
	}
	return nil
}
