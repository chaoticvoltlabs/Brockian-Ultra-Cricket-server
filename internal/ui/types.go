package ui

type ComponentEnvelope struct {
	Component string                 `json:"component"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source,omitempty"`
	Status    Status                 `json:"status"`
	Data      map[string]interface{} `json:"data"`
	Options   map[string]interface{} `json:"options"`
	Resolved  map[string]interface{} `json:"resolved"`
}

type Status struct {
	OK       bool     `json:"ok"`
	Warnings []string `json:"warnings"`
}

type ScreenModel struct {
	Screen ScreenMeta                     `json:"screen"`
	Layout LayoutMeta                     `json:"layout"`
	Theme  ThemeMeta                      `json:"theme"`
	Regions map[string][]ComponentEnvelope `json:"regions"`
}

type ScreenMeta struct {
	Name       string `json:"name"`
	Title      string `json:"title,omitempty"`
	Layout     string `json:"layout"`
	Theme      string `json:"theme"`
	GeneratedAt string `json:"generated_at"`
	DeviceMode string `json:"device_mode"`
}

type LayoutMeta struct {
	Regions []string `json:"regions"`
}

type ThemeMeta struct {
	Name   string            `json:"name"`
	Tokens map[string]string `json:"tokens"`
}
