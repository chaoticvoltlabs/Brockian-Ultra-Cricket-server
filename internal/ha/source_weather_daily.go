package ha

import (
	"time"

	"buc/internal/support"
)

func BuildWeatherDaily(name string, entity *EntityState) SourceResult {
	res := SourceResult{
		Name:   name,
		Type:   "ha_entity",
		Status: Status{OK: true, Warnings: []string{}},
		Data: map[string]interface{}{
			"days": []interface{}{},
		},
	}
	if entity == nil {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "missing entity")
		return res
	}

	res.EntityID = entity.EntityID
	norm := NormalizeEntity(entity)

	raw, ok := norm.Attributes["days"]
	if !ok {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "missing attribute days")
		res.Data["updated_at"] = norm.LastUpdated
		return res
	}

	items, ok := raw.([]interface{})
	if !ok {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "attribute days is not an array")
		res.Data["updated_at"] = norm.LastUpdated
		return res
	}

	out := make([]interface{}, 0, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		row := map[string]interface{}{}

		for _, k := range []string{
			"date", "day_label", "weather_code", "t_min", "t_max",
			"pop", "ws_max", "ws_bft_max", "wg_max", "wg_bft_max", "wd", "wd_label",
		} {
			if v, exists := obj[k]; exists {
				row[k] = v
			}
		}

		if dateStr, ok := obj["date"].(string); ok {
			if _, exists := row["day_label"]; !exists {
				if t, err := time.Parse("2006-01-02", dateStr); err == nil {
					row["day_label"] = weekdayLabelNL(t.Weekday())
				}
			}
		}

		if wd, ok := support.ToFloat(obj["wd"]); ok {
			row["wd"] = wd
			if _, exists := row["wd_label"]; !exists {
				row["wd_label"] = support.WindLabel16(wd)
			}
		}

		if ws, ok := support.ToFloat(obj["ws_max"]); ok {
			row["ws_max"] = ws
			if _, exists := row["ws_bft_max"]; !exists {
				row["ws_bft_max"] = support.Beaufort(ws)
			}
		}

		if wg, ok := support.ToFloat(obj["wg_max"]); ok {
			row["wg_max"] = wg
			if _, exists := row["wg_bft_max"]; !exists {
				row["wg_bft_max"] = support.Beaufort(wg)
			}
		}

		out = append(out, row)
	}

	res.Data["days"] = out
	res.Data["updated_at"] = norm.LastUpdated
	return res
}

func weekdayLabelNL(wd time.Weekday) string {
	switch wd {
	case time.Sunday:
		return "ZO"
	case time.Monday:
		return "MA"
	case time.Tuesday:
		return "DI"
	case time.Wednesday:
		return "WO"
	case time.Thursday:
		return "DO"
	case time.Friday:
		return "VR"
	case time.Saturday:
		return "ZA"
	default:
		return ""
	}
}
