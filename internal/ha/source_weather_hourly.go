package ha

import (
	"buc/internal/support"
)

func BuildWeatherHourly(name string, entity *EntityState) SourceResult {
	res := SourceResult{
		Name:   name,
		Type:   "ha_entity",
		Status: Status{OK: true, Warnings: []string{}},
		Data: map[string]interface{}{
			"forecast": []interface{}{},
		},
	}
	if entity == nil {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "missing entity")
		return res
	}

	res.EntityID = entity.EntityID
	norm := NormalizeEntity(entity)

	raw, ok := norm.Attributes["forecast"]
	if !ok {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "missing attribute forecast")
		res.Data["updated_at"] = norm.LastUpdated
		return res
	}

	items, ok := raw.([]interface{})
	if !ok {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "attribute forecast is not an array")
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

		for _, k := range []string{"hour", "t", "ws", "ws_bft", "wg", "wg_bft", "wd", "wd_label"} {
			if v, exists := obj[k]; exists {
				row[k] = v
			}
		}

		if wd, ok := support.ToFloat(obj["wd"]); ok {
			row["wd"] = wd
			if _, exists := row["wd_label"]; !exists {
				row["wd_label"] = support.WindLabel16(wd)
			}
		}

		if ws, ok := support.ToFloat(obj["ws"]); ok {
			row["ws"] = ws
			if _, exists := row["ws_bft"]; !exists {
				row["ws_bft"] = support.Beaufort(ws)
			}
		}

		if wg, ok := support.ToFloat(obj["wg"]); ok {
			row["wg"] = wg
			if _, exists := row["wg_bft"]; !exists {
				row["wg_bft"] = support.Beaufort(wg)
			}
		}

		out = append(out, row)
	}

	res.Data["forecast"] = out
	res.Data["updated_at"] = norm.LastUpdated
	return res
}
