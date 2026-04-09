package ha

import (
	"sort"

	"buc/internal/support"
)

func BuildIndoorPayload(name string, entity *EntityState) SourceResult {
	res := SourceResult{
		Name:   name,
		Type:   "ha_entity",
		Status: Status{OK: true, Warnings: []string{}},
		Data: map[string]interface{}{
			"tiles": []interface{}{},
		},
	}
	if entity == nil {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "missing entity")
		return res
	}

	res.EntityID = entity.EntityID
	norm := NormalizeEntity(entity)

	raw, ok := norm.Attributes["tiles"]
	if !ok {
		raw, ok = norm.Attributes["indoor"]
	}
	if !ok {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "missing attribute tiles/indoor")
		res.Data["updated_at"] = norm.LastUpdated
		return res
	}

	items, ok := raw.([]interface{})
	if !ok {
		res.Status.OK = false
		res.Status.Warnings = append(res.Status.Warnings, "attribute tiles/indoor is not an array")
		res.Data["updated_at"] = norm.LastUpdated
		return res
	}

	type tileWrap struct {
		Page int
		Row  int
		Col  int
		Tile map[string]interface{}
	}
	tmp := make([]tileWrap, 0, len(items))

	for _, item := range items {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		row := map[string]interface{}{}
		for _, k := range []string{"page", "row", "col", "key", "label", "temp", "temp_unit", "hum", "hum_unit"} {
			if v, exists := obj[k]; exists {
				row[k] = v
			}
		}

		page := intFromAny(obj["page"])
		r := intFromAny(obj["row"])
		c := intFromAny(obj["col"])

		if v, ok := support.ToFloat(obj["temp"]); ok {
			row["temp"] = v
		}
		if v, ok := support.ToFloat(obj["hum"]); ok {
			row["hum"] = v
		}

		tmp = append(tmp, tileWrap{
			Page: page,
			Row:  r,
			Col:  c,
			Tile: row,
		})
	}

	sort.Slice(tmp, func(i, j int) bool {
		if tmp[i].Page != tmp[j].Page {
			return tmp[i].Page < tmp[j].Page
		}
		if tmp[i].Row != tmp[j].Row {
			return tmp[i].Row < tmp[j].Row
		}
		return tmp[i].Col < tmp[j].Col
	})

	out := make([]interface{}, 0, len(tmp))
	for _, t := range tmp {
		out = append(out, t.Tile)
	}

	res.Data["tiles"] = out
	res.Data["updated_at"] = norm.LastUpdated
	return res
}

func intFromAny(v interface{}) int {
	if f, ok := support.ToFloat(v); ok {
		return int(f)
	}
	return 0
}
