package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"buc/internal/app"
	"buc/internal/config"
	"buc/internal/ha"
)

func TestPanelWeatherHandler(t *testing.T) {
	haServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/states/sensor.panel_weather_current":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "sensor.panel_weather_current",
				"attributes": map[string]interface{}{
					"current": map[string]interface{}{
						"t":          11.7,
						"feels_like": 0.9,
						"rh":         88,
						"pressure":   996,
						"ws":         4,
						"ws_bft":     1,
						"wg":         7,
						"wg_bft":     1,
						"wd":         304,
					},
				},
			})
		case "/api/states/sensor.panel_overview_payload":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "sensor.panel_overview_payload",
				"attributes": map[string]interface{}{
					"night_mode": true,
					"weather_meta": map[string]interface{}{
						"pressure_trend_24h": "1002.0, 1001.4, 1001.4, 1000.8, invalid, 999.1, 998.5",
					},
				},
			})
		case "/api/states/sensor.panel_indoor_payload":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "sensor.panel_indoor_payload",
				"attributes": map[string]interface{}{
					"tiles": []map[string]interface{}{
						{"page": 2, "row": 0, "col": 0, "key": "office", "temp": 17.1, "hum": 37},
						{"page": 2, "row": 0, "col": 1, "key": "bathroom", "temp": 20.3, "hum": 29},
						{"page": 2, "row": 0, "col": 2, "key": "bedroom", "temp": 17.0, "hum": 53},
						{"page": 2, "row": 0, "col": 3, "key": "wardrobe", "temp": 19.3, "hum": 36},
						{"page": 2, "row": 1, "col": 0, "key": "kitchen", "temp": 20.1, "hum": 33},
						{"page": 2, "row": 1, "col": 1, "key": "living", "temp": 19.9, "hum": 37},
						{"page": 2, "row": 1, "col": 2, "key": "library", "temp": 17.1, "hum": 42},
						{"page": 2, "row": 1, "col": 3, "key": "sunroom", "temp": 18.3, "hum": 34},
						{"page": 2, "row": 2, "col": 0, "key": "empty"},
						{"page": 2, "row": 2, "col": 1, "key": "laundry", "temp": 17.3, "hum": 48},
						{"page": 2, "row": 2, "col": 2, "key": "empty"},
						{"page": 2, "row": 2, "col": 3, "key": "studio", "temp": 22.2, "hum": 42},
						{"page": 1, "row": 0, "col": 0, "key": "ignored", "temp": 99.9, "hum": 99},
					},
				},
			})
		case "/api/states/light.room_light_a":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "light.room_light_a",
				"state":     "on",
				"attributes": map[string]interface{}{},
			})
		case "/api/states/switch.room_light_b":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "switch.room_light_b",
				"state":     "off",
				"attributes": map[string]interface{}{},
			})
		case "/api/states/light.room_light_c":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "light.room_light_c",
				"state":     "unavailable",
				"attributes": map[string]interface{}{},
			})
		case "/api/states/switch.media_power":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "switch.media_power",
				"state":     "off",
				"attributes": map[string]interface{}{},
			})
		case "/api/states/light.room_ambient":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "light.room_ambient",
				"state":     "on",
				"attributes": map[string]interface{}{
					"brightness": 191,
					"rgb_color":  []interface{}{255, 136, 32},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer haServer.Close()

	cfg := &config.AllConfig{
		Sources: config.SourcesFile{
			Sources: map[string]config.SourceConfig{
				"weather_current": {
					Type:     "ha_entity",
					EntityID: "sensor.panel_weather_current",
				},
				"overview_payload": {
					Type:     "ha_entity",
					EntityID: "sensor.panel_overview_payload",
				},
				"indoor_payload": {
					Type:     "ha_entity",
					EntityID: "sensor.panel_indoor_payload",
				},
			},
		},
		PanelDevices: config.PanelDevicesFile{
			Default: "room_alpha",
			Devices: map[string]config.PanelDeviceConfig{
				"room_alpha": {
					DeviceType: "panel_4_3",
					Profile:    "room_alpha",
				},
			},
		},
		PanelProfiles: config.PanelProfilesFile{
			Profiles: map[string]config.PanelProfileConfig{
				"room_alpha": {
					Page3: config.PanelProfilePage3Config{
						Targets: []config.PanelSlotConfig{
							{Label: "Light A", Target: "light_a", Action: "toggle"},
							{Label: "Light B", Target: "light_b", Action: "toggle"},
							{Label: "Light C", Target: "light_c", Action: "toggle"},
							{Label: "Media", Target: "media_power", Action: "toggle"},
						},
					},
				},
			},
		},
		PanelCommands: config.PanelCommandsFile{
			Commands: map[string]config.PanelCommandConfig{
				"light_a:toggle": {
					Domain:   "light",
					Service:  "toggle",
					EntityID: "light.room_light_a",
				},
				"light_b:toggle": {
					Domain:   "switch",
					Service:  "toggle",
					EntityID: "switch.room_light_b",
				},
				"light_c:toggle": {
					Domain:   "light",
					Service:  "toggle",
					EntityID: "light.room_light_c",
				},
				"media_power:toggle": {
					Domain:   "switch",
					Service:  "toggle",
					EntityID: "switch.media_power",
				},
				"ambient:toggle": {
					Domain:   "light",
					Service:  "toggle",
					EntityID: "light.room_ambient",
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/panel/weather", nil)
	rec := httptest.NewRecorder()

	Router(app.New(cfg, ha.NewClient(haServer.URL, "test-token"))).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("expected json content type, got %q", got)
	}

	var resp panelWeatherResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.OutsideTempC != 11.7 {
		t.Fatalf("expected outside_temp_c 11.7, got %v", resp.OutsideTempC)
	}

	if len(resp.PressureTrend24h) != 6 {
		t.Fatalf("expected 6 pressure points, got %d", len(resp.PressureTrend24h))
	}

	if resp.PressureTrend24h[0] != 1002.0 || resp.PressureTrend24h[5] != 998.5 {
		t.Fatalf("unexpected pressure trend: %#v", resp.PressureTrend24h)
	}

	if len(resp.IndoorZones) != 12 {
		t.Fatalf("expected 12 indoor zones, got %d", len(resp.IndoorZones))
	}
	if !resp.NightMode {
		t.Fatalf("expected night_mode true, got %#v", resp.NightMode)
	}
	if resp.AmbientBrightnessPct == nil || *resp.AmbientBrightnessPct != 75 {
		t.Fatalf("expected ambient brightness 75, got %#v", resp.AmbientBrightnessPct)
	}
	if len(resp.AmbientRGB) != 3 || resp.AmbientRGB[0] != 255 || resp.AmbientRGB[1] != 136 || resp.AmbientRGB[2] != 32 {
		t.Fatalf("unexpected ambient rgb: %#v", resp.AmbientRGB)
	}
	if got := resp.Page3TargetStates["light_a"]; got != "on" {
		t.Fatalf("expected light_a state on, got %#v", got)
	}
	if got := resp.Page3TargetStates["light_b"]; got != "off" {
		t.Fatalf("expected light_b state off, got %#v", got)
	}
	if got := resp.Page3TargetStates["light_c"]; got != "unavailable" {
		t.Fatalf("expected light_c state unavailable, got %#v", got)
	}
	if got := resp.Page3TargetStates["media_power"]; got != "off" {
		t.Fatalf("expected media_power state off, got %#v", got)
	}

	assertIndoorZone(t, resp.IndoorZones[0], 17.1, 37)
	assertIndoorZone(t, resp.IndoorZones[7], 18.3, 34)
	assertIndoorZoneNull(t, resp.IndoorZones[8])
	assertIndoorZone(t, resp.IndoorZones[9], 17.3, 48)
	assertIndoorZoneNull(t, resp.IndoorZones[10])
	assertIndoorZone(t, resp.IndoorZones[11], 22.2, 42)
}

func TestPanelWeatherHandlerFallsBackToCurrentEntityPressureTrend(t *testing.T) {
	haServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/states/sensor.panel_weather_current":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "sensor.panel_weather_current",
				"attributes": map[string]interface{}{
					"current": map[string]interface{}{
						"t":                  9.96,
						"feels_like":         0.5,
						"rh":                 45,
						"pressure":           993,
						"ws":                 33,
						"ws_bft":             5,
						"wg":                 72,
						"wg_bft":             8,
						"wd":                 251,
						"pressure_trend_24h": []float64{995.4, 994.8, 994.2, 993.7},
					},
				},
			})
		case "/api/states/sensor.panel_overview_payload":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"entity_id": "sensor.panel_overview_payload",
				"attributes": map[string]interface{}{},
			})
		case "/api/states/sensor.panel_indoor_payload":
			http.NotFound(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer haServer.Close()

	cfg := &config.AllConfig{
		Sources: config.SourcesFile{
			Sources: map[string]config.SourceConfig{
				"weather_current": {
					Type:     "ha_entity",
					EntityID: "sensor.panel_weather_current",
				},
				"overview_payload": {
					Type:     "ha_entity",
					EntityID: "sensor.panel_overview_payload",
				},
				"indoor_payload": {
					Type:     "ha_entity",
					EntityID: "sensor.panel_indoor_payload",
				},
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/panel/weather", nil)
	rec := httptest.NewRecorder()

	Router(app.New(cfg, ha.NewClient(haServer.URL, "test-token"))).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp panelWeatherResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(resp.PressureTrend24h) != 4 {
		t.Fatalf("expected 4 pressure points, got %d (%#v)", len(resp.PressureTrend24h), resp.PressureTrend24h)
	}

	if resp.PressureTrend24h[0] != 995.4 || resp.PressureTrend24h[3] != 993.7 {
		t.Fatalf("unexpected pressure trend: %#v", resp.PressureTrend24h)
	}

	if len(resp.IndoorZones) != 12 {
		t.Fatalf("expected 12 indoor zones, got %d", len(resp.IndoorZones))
	}
	for i, zone := range resp.IndoorZones {
		if zone.TempC != nil || zone.RHPct != nil {
			t.Fatalf("expected indoor zone %d to be null-filled, got %#v", i, zone)
		}
	}
}

func TestExtractPressureTrend24h(t *testing.T) {
	tests := []struct {
		name string
		data map[string]interface{}
		want []float64
	}{
		{
			name: "typed slice",
			data: map[string]interface{}{
				"weather_meta": map[string]interface{}{
					"pressure_trend_24h": []float64{991.9, 991.8, 992.1},
				},
			},
			want: []float64{991.9, 991.8, 992.1},
		},
		{
			name: "json string array",
			data: map[string]interface{}{
				"weather_meta": map[string]interface{}{
					"pressure_trend_24h": `[991.9, 991.8, 992.1]`,
				},
			},
			want: []float64{991.9, 991.8, 992.1},
		},
		{
			name: "comma separated string",
			data: map[string]interface{}{
				"weather_meta": map[string]interface{}{
					"pressure_trend_24h": `991.9, 991.8, bad, 992.1`,
				},
			},
			want: []float64{991.9, 991.8, 992.1},
		},
		{
			name: "missing value",
			data: map[string]interface{}{
				"weather_meta": map[string]interface{}{},
			},
			want: []float64{},
		},
		{
			name: "malformed string",
			data: map[string]interface{}{
				"weather_meta": map[string]interface{}{
					"pressure_trend_24h": `nope`,
				},
			},
			want: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPressureTrend24h(tt.data)
			if len(got) != len(tt.want) {
				t.Fatalf("expected %d pressure points, got %d (%#v)", len(tt.want), len(got), got)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("expected point %d to be %v, got %v (%#v)", i, tt.want[i], got[i], got)
				}
			}
		})
	}
}

func TestExtractIndoorZonesFromEntityFallsBackToZoneKeyOrder(t *testing.T) {
	entity := &ha.EntityState{
		EntityID: "sensor.panel_indoor_payload",
		Attributes: map[string]interface{}{
			"tiles": []interface{}{
				map[string]interface{}{"key": "office", "temp": 17.1, "hum": 37},
				map[string]interface{}{"key": "bathroom", "temp": 20.3, "hum": 29},
				map[string]interface{}{"key": "bedroom", "temp": 17.0, "hum": 53},
				map[string]interface{}{"key": "wardrobe", "temp": 19.3, "hum": 36},
				map[string]interface{}{"key": "kitchen", "temp": 20.1, "hum": 33},
				map[string]interface{}{"key": "living", "temp": 19.9, "hum": 37},
				map[string]interface{}{"key": "library", "temp": 17.1, "hum": 42},
				map[string]interface{}{"key": "sunroom", "temp": 18.3, "hum": 34},
				map[string]interface{}{"key": "laundry", "temp": 17.3, "hum": 48},
				map[string]interface{}{"key": "studio", "temp": 22.2, "hum": 42},
			},
		},
	}

	got := extractIndoorZonesFromEntity(entity, panelIndoorPage, panelIndoorRows, panelIndoorCols)
	if len(got) != 12 {
		t.Fatalf("expected 12 indoor zones, got %d", len(got))
	}

	assertIndoorZone(t, got[0], 17.1, 37)
	assertIndoorZone(t, got[7], 18.3, 34)
	assertIndoorZoneNull(t, got[8])
	assertIndoorZone(t, got[9], 17.3, 48)
	assertIndoorZoneNull(t, got[10])
	assertIndoorZone(t, got[11], 22.2, 42)
}

func assertIndoorZone(t *testing.T, zone panelIndoorZone, wantTemp float64, wantRH int) {
	t.Helper()
	if zone.TempC == nil || zone.RHPct == nil {
		t.Fatalf("expected populated indoor zone, got %#v", zone)
	}
	if *zone.TempC != wantTemp || *zone.RHPct != wantRH {
		t.Fatalf("expected indoor zone temp=%v rh=%d, got %#v", wantTemp, wantRH, zone)
	}
}

func assertIndoorZoneNull(t *testing.T, zone panelIndoorZone) {
	t.Helper()
	if zone.TempC != nil || zone.RHPct != nil {
		t.Fatalf("expected null indoor zone, got %#v", zone)
	}
}
