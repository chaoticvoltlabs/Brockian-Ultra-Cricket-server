package ha

import (
	"fmt"

	"buc/internal/support"
)

func BuildWeatherCurrent(name string, entity *EntityState) SourceResult {
	res := SourceResult{
		Name:   name,
		Type:   "ha_entity",
		Status: Status{OK: true, Warnings: []string{}},
		Data:   map[string]interface{}{},
	}
	if entity == nil {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "missing entity")
		return res
	}

	res.EntityID = entity.EntityID
	norm := NormalizeEntity(entity)

	raw, ok := norm.Attributes["current"]
	if !ok {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "missing attribute current")
		return res
	}

	obj, ok := raw.(map[string]interface{})
	if !ok {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "attribute current is not an object")
		return res
	}

	copyField := func(key string) {
		if v, exists := obj[key]; exists {
			res.Data[key] = v
		}
	}

	for _, k := range []string{"t", "feels_like", "rh", "pressure", "ws", "ws_bft", "wg", "wg_bft", "wd", "precip"} {
		copyField(k)
	}

	if wd, ok := support.ToFloat(obj["wd"]); ok {
		res.Data["wd"] = wd
		res.Data["wd_label"] = support.WindLabel16(wd)
	}

	if ws, ok := support.ToFloat(obj["ws"]); ok {
		res.Data["ws"] = ws
		if _, exists := res.Data["ws_bft"]; !exists {
			res.Data["ws_bft"] = support.Beaufort(ws)
		}
	}

	if wg, ok := support.ToFloat(obj["wg"]); ok {
		res.Data["wg"] = wg
		if _, exists := res.Data["wg_bft"]; !exists {
			res.Data["wg_bft"] = support.Beaufort(wg)
		}
	}

	res.Data["updated_at"] = norm.LastUpdated
	return res
}

func BuildSource(name string, entityID string, entity *EntityState) (SourceResult, error) {
	switch name {
	case "weather_current":
		return BuildWeatherCurrent(name, entity), nil
	case "weather_hourly":
		return BuildWeatherHourly(name, entity), nil
	case "weather_daily":
		return BuildWeatherDaily(name, entity), nil
	case "indoor_payload":
		return BuildIndoorPayload(name, entity), nil
	case "overview_payload":
		return BuildOverviewPayload(name, entity), nil
	default:
		return SourceResult{}, fmt.Errorf("no source adapter for %s", name)
	}
}
