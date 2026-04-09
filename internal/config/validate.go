package config

import "fmt"

func Validate(cfg *AllConfig) error {
	for name, comp := range cfg.Components.Components {
		if comp.Source != "" {
			if _, ok := cfg.Sources.Sources[comp.Source]; !ok {
				return fmt.Errorf("component %q references unknown source %q", name, comp.Source)
			}
		}
	}

	for name, screen := range cfg.Screens.Screens {
		for _, comps := range screen.Regions {
			for _, compName := range comps {
				if _, ok := cfg.Components.Components[compName]; !ok {
					return fmt.Errorf("screen %q references unknown component %q", name, compName)
				}
			}
		}
	}

	for name, dev := range cfg.Devices.Devices {
		if _, ok := cfg.Screens.Screens[dev.Screen]; !ok {
			return fmt.Errorf("device %q references unknown screen %q", name, dev.Screen)
		}
		if _, ok := cfg.Themes.Themes[dev.Theme]; !ok {
			return fmt.Errorf("device %q references unknown theme %q", name, dev.Theme)
		}
		if dev.Resolution.Width <= 0 || dev.Resolution.Height <= 0 {
			return fmt.Errorf("device %q has invalid resolution", name)
		}
	}

	if cfg.PanelDevices.Default != "" {
		if _, ok := cfg.PanelDevices.Devices[cfg.PanelDevices.Default]; !ok {
			return fmt.Errorf("panel device default %q not found", cfg.PanelDevices.Default)
		}
	}

	for name, dev := range cfg.PanelDevices.Devices {
		if dev.DeviceType == "" {
			return fmt.Errorf("panel device %q missing device_type", name)
		}
		if _, ok := cfg.DeviceTypes.DeviceTypes[dev.DeviceType]; !ok {
			return fmt.Errorf("panel device %q references unknown device_type %q", name, dev.DeviceType)
		}
		if dev.Profile == "" {
			return fmt.Errorf("panel device %q missing profile", name)
		}
		if _, ok := cfg.PanelProfiles.Profiles[dev.Profile]; !ok {
			return fmt.Errorf("panel device %q references unknown profile %q", name, dev.Profile)
		}
	}

	for name, profile := range cfg.PanelProfiles.Profiles {
		if len(profile.Page3.Scenes) == 0 && len(profile.Page3.Targets) == 0 {
			return fmt.Errorf("panel profile %q has no page3 content", name)
		}
		for _, slot := range profile.Page3.Scenes {
			if err := validatePanelSlot(cfg, name, "scene", slot); err != nil {
				return err
			}
		}
		for _, slot := range profile.Page3.Targets {
			if err := validatePanelSlot(cfg, name, "target", slot); err != nil {
				return err
			}
		}
		if profile.Page3.LongPress.Target != "" || profile.Page3.LongPress.Action != "" {
			if err := validatePanelSlot(cfg, name, "long_press", profile.Page3.LongPress); err != nil {
				return err
			}
		}
	}

	for themeName, theme := range cfg.Themes.Themes {
		for scaleName, scale := range theme.TemperatureScales {
			for _, band := range scale.Bands {
				if _, ok := theme.Tokens[band.Token]; !ok {
					return fmt.Errorf("theme %q scale %q references unknown token %q", themeName, scaleName, band.Token)
				}
			}
		}
	}

	return nil
}

func validatePanelSlot(cfg *AllConfig, profileName string, slotKind string, slot PanelSlotConfig) error {
	key := slot.Target + ":" + slot.Action
	if slot.Target == "" || slot.Action == "" {
		return fmt.Errorf("panel profile %q %s slot missing target or action", profileName, slotKind)
	}
	if _, ok := cfg.PanelCommands.Commands[key]; !ok {
		return fmt.Errorf("panel profile %q %s slot references unknown command %q", profileName, slotKind, key)
	}
	return nil
}
