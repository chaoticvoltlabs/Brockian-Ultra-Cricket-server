package ha

func BuildOverviewPayload(name string, entity *EntityState) SourceResult {
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

	for _, k := range []string{"map_asset", "indoor", "weather", "weather_meta"} {
		if v, exists := norm.Attributes[k]; exists {
			res.Data[k] = v
		}
	}

	if _, exists := res.Data["map_asset"]; !exists {
		res.Status.Warnings = append(res.Status.Warnings, "missing attribute map_asset")
	}

	res.Data["updated_at"] = norm.LastUpdated
	return res
}
