package httpapi

import (
	"log"
	"net/http"
	"strings"

	"buc/internal/app"
	"buc/internal/config"
	"buc/internal/support"
)

type panelConfigSlot struct {
	Label  string `json:"label"`
	Target string `json:"target"`
	Action string `json:"action"`
}

type panelConfigResponse struct {
	Profile string           `json:"profile"`
	Page3   panelConfigPage3 `json:"page3"`
}

type panelConfigPage3 struct {
	Scenes    []panelConfigSlot `json:"scenes"`
	Targets   []panelConfigSlot `json:"targets"`
	LongPress panelConfigSlot   `json:"long_press"`
}

func PanelConfigHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := recordPanelRequest(r)
		profile, resolvedDevice, resolvedType := panelProfileForRequest(a, info.PanelMAC)
		log.Printf("panel config resolved panel_mac=%s panel_ip=%s profile=%s",
			info.PanelMAC, info.PanelIP, profile.Profile)
		log.Printf("panel device resolved panel_mac=%s device=%s device_type=%s profile=%s",
			info.PanelMAC, resolvedDevice, resolvedType, profile.Profile)
		support.JSON(w, http.StatusOK, profile)
	}
}

func panelProfileForRequest(a *app.App, mac string) (panelConfigResponse, string, string) {
	if a == nil || a.Config == nil {
		return panelFallbackProfile()
	}

	deviceName, deviceCfg, ok := panelDeviceForMAC(a.Config, mac)
	if !ok {
		return panelDefaultProfile(a.Config)
	}

	profileCfg, ok := a.Config.PanelProfiles.Profiles[deviceCfg.Profile]
	if !ok {
		return panelDefaultProfile(a.Config)
	}

	return panelResponseFromConfig(deviceCfg.Profile, profileCfg), deviceName, deviceCfg.DeviceType
}

func panelDeviceForMAC(cfg *config.AllConfig, mac string) (string, config.PanelDeviceConfig, bool) {
	normalized := strings.ToLower(strings.TrimSpace(mac))
	if normalized == "" {
		return "", config.PanelDeviceConfig{}, false
	}

	for name, dev := range cfg.PanelDevices.Devices {
		if strings.EqualFold(strings.TrimSpace(dev.Identity.MAC), normalized) {
			return name, dev, true
		}
	}

	return "", config.PanelDeviceConfig{}, false
}

func panelDefaultProfile(cfg *config.AllConfig) (panelConfigResponse, string, string) {
	if cfg == nil {
		return panelFallbackProfile()
	}

	defaultName := strings.TrimSpace(cfg.PanelDevices.Default)
	if defaultName != "" {
		if dev, ok := cfg.PanelDevices.Devices[defaultName]; ok {
			if profileCfg, ok := cfg.PanelProfiles.Profiles[dev.Profile]; ok {
				return panelResponseFromConfig(dev.Profile, profileCfg), defaultName, dev.DeviceType
			}
		}
	}

	return panelFallbackProfile()
}

func panelResponseFromConfig(profileName string, profileCfg config.PanelProfileConfig) panelConfigResponse {
	return panelConfigResponse{
		Profile: profileName,
		Page3: panelConfigPage3{
			Scenes:    panelSlotsFromConfig(profileCfg.Page3.Scenes),
			Targets:   panelSlotsFromConfig(profileCfg.Page3.Targets),
			LongPress: panelSlotFromConfig(profileCfg.Page3.LongPress),
		},
	}
}

func panelSlotsFromConfig(slots []config.PanelSlotConfig) []panelConfigSlot {
	out := make([]panelConfigSlot, 0, len(slots))
	for _, slot := range slots {
		out = append(out, panelSlotFromConfig(slot))
	}
	return out
}

func panelSlotFromConfig(slot config.PanelSlotConfig) panelConfigSlot {
	return panelConfigSlot{
		Label:  slot.Label,
		Target: slot.Target,
		Action: slot.Action,
	}
}

func panelFallbackProfile() (panelConfigResponse, string, string) {
	return panelConfigResponse{
		Profile: "fallback",
		Page3:   panelConfigPage3{},
	}, "", ""
}
