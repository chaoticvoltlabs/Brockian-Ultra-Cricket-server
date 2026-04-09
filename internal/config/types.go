package config

type SourcesFile struct {
	Sources map[string]SourceConfig `json:"sources"`
}

type SourceConfig struct {
	Type     string `json:"type"`
	EntityID string `json:"entity_id,omitempty"`
}

type ComponentsFile struct {
	Components map[string]ComponentConfig `json:"components"`
}

type ComponentConfig struct {
	Type    string                 `json:"type"`
	Source  string                 `json:"source,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type ScreensFile struct {
	Screens map[string]ScreenConfig `json:"screens"`
}

type ScreenConfig struct {
	Layout  string              `json:"layout"`
	Title   string              `json:"title,omitempty"`
	Regions map[string][]string `json:"regions"`
}

type DevicesFile struct {
	Devices map[string]DeviceConfig `json:"devices"`
}

type DeviceConfig struct {
	Mode           string           `json:"mode"`
	Screen         string           `json:"screen"`
	Theme          string           `json:"theme"`
	Orientation    string           `json:"orientation"`
	Resolution     ResolutionConfig `json:"resolution"`
	RefreshSeconds int              `json:"refresh_seconds"`
}

type ResolutionConfig struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DeviceTypesFile struct {
	DeviceTypes map[string]DeviceTypeConfig `json:"device_types"`
}

type DeviceTypeConfig struct {
	Family string `json:"family"`
}

type PanelDevicesFile struct {
	Default string                       `json:"default"`
	Devices map[string]PanelDeviceConfig `json:"devices"`
}

type PanelDeviceConfig struct {
	Identity   PanelIdentityConfig `json:"identity"`
	DeviceType string              `json:"device_type"`
	Profile    string              `json:"profile"`
}

type PanelIdentityConfig struct {
	MAC string `json:"mac"`
}

type PanelProfilesFile struct {
	Profiles map[string]PanelProfileConfig `json:"profiles"`
}

type PanelProfileConfig struct {
	Page3 PanelProfilePage3Config `json:"page3"`
}

type PanelProfilePage3Config struct {
	Scenes    []PanelSlotConfig `json:"scenes"`
	Targets   []PanelSlotConfig `json:"targets"`
	LongPress PanelSlotConfig   `json:"long_press"`
}

type PanelSlotConfig struct {
	Label  string `json:"label"`
	Target string `json:"target"`
	Action string `json:"action"`
}

type PanelCommandsFile struct {
	Commands map[string]PanelCommandConfig `json:"commands"`
}

type PanelCommandConfig struct {
	Domain   string `json:"domain"`
	Service  string `json:"service"`
	EntityID string `json:"entity_id"`
}

type ThemesFile struct {
	Themes map[string]ThemeConfig `json:"themes"`
}

type ThemeConfig struct {
	Tokens            map[string]string              `json:"tokens"`
	TemperatureScales map[string]TemperatureScaleDef `json:"temperature_scales"`
}

type TemperatureScaleDef struct {
	Bands []TemperatureBand `json:"bands"`
}

type TemperatureBand struct {
	Max   float64 `json:"max"`
	Token string  `json:"token"`
}

type AllConfig struct {
	Sources       SourcesFile
	Components    ComponentsFile
	Screens       ScreensFile
	Devices       DevicesFile
	DeviceTypes   DeviceTypesFile
	PanelDevices  PanelDevicesFile
	PanelProfiles PanelProfilesFile
	PanelCommands PanelCommandsFile
	Themes        ThemesFile
}
