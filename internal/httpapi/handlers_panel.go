package httpapi

import (
	"encoding/json"
	"math"
	"net/http"
	"reflect"
	"strings"

	"buc/internal/app"
	"buc/internal/config"
	"buc/internal/ha"
	"buc/internal/support"
)

type panelWeatherResponse struct {
	OutsideTempC     float64           `json:"outside_temp_c"`
	FeelsLikeC       float64           `json:"feels_like_c"`
	WindBft          int               `json:"wind_bft"`
	WindKmh          int               `json:"wind_kmh"`
	GustBft          int               `json:"gust_bft"`
	GustKmh          int               `json:"gust_kmh"`
	WindDirDeg       int               `json:"wind_dir_deg"`
	HumidityPct      int               `json:"humidity_pct"`
	PressureHPA      int               `json:"pressure_hpa"`
	PressureTrend24h []float64         `json:"pressure_trend_24h"`
	IndoorZones      []panelIndoorZone `json:"indoor_zones"`
}

type panelIndoorZone struct {
	TempC *float64 `json:"temp_c"`
	RHPct *int     `json:"rh_pct"`
}

const (
	panelWeatherCurrentSource  = "weather_current"
	panelOverviewPayloadSource = "overview_payload"
	panelIndoorPayloadSource   = "indoor_payload"
	panelIndoorPage            = 2
	panelIndoorRows            = 3
	panelIndoorCols            = 4
)

var panelIndoorZoneOrder = []string{
	"office",
	"bathroom",
	"bedroom",
	"wardrobe",
	"kitchen",
	"living",
	"library",
	"sunroom",
	"servers",
	"laundry",
	"utility",
	"studio",
}

func PanelWeatherHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recordPanelRequest(r)

		if a == nil || a.HAClient == nil {
			support.JSON(w, http.StatusServiceUnavailable, map[string]interface{}{
				"error": "ha client not configured",
			})
			return
		}

		currentEntityID := sourceEntityID(a.Config, panelWeatherCurrentSource, "sensor.panel_weather_current")
		currentEntity, err := a.HAClient.GetState(currentEntityID)
		if err != nil {
			support.JSON(w, http.StatusBadGateway, map[string]interface{}{
				"error":     "failed to load weather_current from HA",
				"entity_id": currentEntityID,
				"details":   err.Error(),
			})
			return
		}

		currentSource, err := ha.BuildSource(panelWeatherCurrentSource, currentEntityID, currentEntity)
		if err != nil {
			support.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		resp := panelWeatherResponse{
			OutsideTempC:     extractFloat(currentSource.Data, "t"),
			FeelsLikeC:       extractFloat(currentSource.Data, "feels_like"),
			WindBft:          extractRoundedInt(currentSource.Data, "ws_bft"),
			WindKmh:          extractRoundedInt(currentSource.Data, "ws"),
			GustBft:          extractRoundedInt(currentSource.Data, "wg_bft"),
			GustKmh:          extractRoundedInt(currentSource.Data, "wg"),
			WindDirDeg:       extractRoundedInt(currentSource.Data, "wd"),
			HumidityPct:      extractRoundedInt(currentSource.Data, "rh"),
			PressureHPA:      extractRoundedInt(currentSource.Data, "pressure"),
			PressureTrend24h: extractPressureTrend24hFromEntity(currentEntity),
			IndoorZones:      emptyIndoorZones(panelIndoorRows * panelIndoorCols),
		}

		overviewEntityID := sourceEntityID(a.Config, panelOverviewPayloadSource, "sensor.panel_overview_payload")
		if overviewEntity, err := a.HAClient.GetState(overviewEntityID); err == nil {
			if overviewSource, buildErr := ha.BuildSource(panelOverviewPayloadSource, overviewEntityID, overviewEntity); buildErr == nil {
				if values := extractPressureTrend24h(overviewSource.Data); len(values) > 0 {
					resp.PressureTrend24h = values
				}
			}
			if len(resp.PressureTrend24h) == 0 {
				if values := extractPressureTrend24hFromEntity(overviewEntity); len(values) > 0 {
					resp.PressureTrend24h = values
				}
			}
		}

		indoorEntityID := sourceEntityID(a.Config, panelIndoorPayloadSource, "sensor.panel_indoor_payload")
		if indoorEntity, err := a.HAClient.GetState(indoorEntityID); err == nil {
			resp.IndoorZones = extractIndoorZonesFromEntity(indoorEntity, panelIndoorPage, panelIndoorRows, panelIndoorCols)
		}

		support.JSON(w, http.StatusOK, resp)
	}
}

func sourceEntityID(cfg *config.AllConfig, sourceName string, fallback string) string {
	if cfg == nil {
		return fallback
	}
	srcCfg, ok := cfg.Sources.Sources[sourceName]
	if !ok || srcCfg.EntityID == "" {
		return fallback
	}
	return srcCfg.EntityID
}

func extractFloat(data map[string]interface{}, key string) float64 {
	if v, ok := support.ToFloat(data[key]); ok {
		return v
	}
	return 0
}

func extractRoundedInt(data map[string]interface{}, key string) int {
	if v, ok := support.ToFloat(data[key]); ok {
		return int(math.Round(v))
	}
	return 0
}

func extractPressureTrend24h(data map[string]interface{}) []float64 {
	if len(data) == 0 {
		return []float64{}
	}

	for _, containerKey := range []string{"weather_meta", "weather", ""} {
		container := data
		if containerKey != "" {
			obj, ok := data[containerKey].(map[string]interface{})
			if !ok {
				continue
			}
			container = obj
		}

		for _, key := range []string{"pressure_trend_24h", "pressure_trend", "pressure_history_24h", "pressure_history"} {
			values := parseNumericSeries(container[key])
			if len(values) > 0 {
				return values
			}
		}
	}

	return []float64{}
}

func extractPressureTrend24hFromEntity(entity *ha.EntityState) []float64 {
	if entity == nil {
		return []float64{}
	}

	norm := ha.NormalizeEntity(entity)
	if norm == nil {
		return []float64{}
	}

	for _, candidate := range []string{"weather_meta", "weather", "current", "payload"} {
		if obj, ok := norm.Attributes[candidate].(map[string]interface{}); ok {
			if values := extractPressureTrend24h(obj); len(values) > 0 {
				return values
			}
		}
	}

	return extractPressureTrend24h(norm.Attributes)
}

func parseNumericSeries(v interface{}) []float64 {
	switch x := v.(type) {
	case nil:
		return []float64{}
	case string:
		return parseNumericSeriesString(x)
	}

	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return []float64{}
	}
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return []float64{}
	}

	out := make([]float64, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		if num, ok := support.ToFloat(rv.Index(i).Interface()); ok {
			out = append(out, num)
		}
	}
	return out
}

func parseNumericSeriesString(raw string) []float64 {
	s := strings.TrimSpace(raw)
	if s == "" {
		return []float64{}
	}

	if len(s) > 1 && s[0] == '[' && s[len(s)-1] == ']' {
		var parsed interface{}
		if err := json.Unmarshal([]byte(s), &parsed); err == nil {
			return parseNumericSeries(parsed)
		}
	}

	parts := strings.Split(s, ",")
	out := make([]float64, 0, len(parts))
	for _, part := range parts {
		if num, ok := support.ToFloat(strings.TrimSpace(part)); ok {
			out = append(out, num)
		}
	}
	return out
}

func emptyIndoorZones(count int) []panelIndoorZone {
	return make([]panelIndoorZone, count)
}

func extractIndoorZonesFromEntity(entity *ha.EntityState, page int, rows int, cols int) []panelIndoorZone {
	zones := emptyIndoorZones(rows * cols)
	if entity == nil {
		return zones
	}

	source, err := ha.BuildSource(panelIndoorPayloadSource, entity.EntityID, entity)
	if err != nil {
		return zones
	}

	tiles, ok := source.Data["tiles"].([]interface{})
	if !ok {
		return zones
	}

	filledByGrid := false
	for _, item := range tiles {
		tile, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if intFromData(tile["page"]) != page {
			continue
		}

		row := intFromData(tile["row"])
		col := intFromData(tile["col"])
		if row < 0 || row >= rows || col < 0 || col >= cols {
			continue
		}

		if key, ok := tile["key"].(string); ok && strings.EqualFold(strings.TrimSpace(key), "empty") {
			continue
		}

		idx := row*cols + col
		zones[idx] = panelIndoorZoneFromTile(tile)
		filledByGrid = true
	}

	if filledByGrid {
		return zones
	}

	indexByKey := map[string]int{}
	for i, key := range panelIndoorZoneOrder {
		indexByKey[key] = i
	}

	for _, item := range tiles {
		tile, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := tile["key"].(string)
		key = strings.ToLower(strings.TrimSpace(key))
		if key == "" || key == "empty" {
			continue
		}
		idx, ok := indexByKey[key]
		if !ok || idx < 0 || idx >= len(zones) {
			continue
		}
		zones[idx] = panelIndoorZoneFromTile(tile)
	}

	return zones
}

func intFromData(v interface{}) int {
	if f, ok := support.ToFloat(v); ok {
		return int(f)
	}
	return 0
}

func panelIndoorZoneFromTile(tile map[string]interface{}) panelIndoorZone {
	zone := panelIndoorZone{}
	if v, ok := support.ToFloat(tile["temp"]); ok {
		temp := v
		zone.TempC = &temp
	}
	if v, ok := support.ToFloat(tile["hum"]); ok {
		rh := int(math.Round(v))
		zone.RHPct = &rh
	}
	return zone
}
