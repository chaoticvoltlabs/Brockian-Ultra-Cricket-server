package ui

import (
	"buc/internal/config"
)

func BuildWindStrip(
	_ *config.AllConfig,
	_ string,
	componentName string,
	componentType string,
	sourceName string,
	data map[string]interface{},
	options map[string]interface{},
) ComponentEnvelope {
	env := ComponentEnvelope{
		Component: componentName,
		Type:      componentType,
		Source:    sourceName,
		Status: Status{
			OK:       true,
			Warnings: []string{},
		},
		Data:     data,
		Options:  options,
		Resolved: map[string]interface{}{},
	}

	mode := stringOption(options, "mode", "gusts")
	unit := stringOption(options, "unit", "bft")
	hours := intOption(options, "hours", 12)

	env.Resolved["display_mode"] = mode
	env.Resolved["display_unit"] = unit

	rawForecast, ok := data["forecast"].([]interface{})
	if !ok {
		env.Status.OK = false
		env.Status.Warnings = append(env.Status.Warnings, "forecast missing or not an array")
		env.Resolved["items"] = []interface{}{}
		return env
	}

	out := make([]interface{}, 0, hours)

	for i, item := range rawForecast {
		if i >= hours {
			break
		}

		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		row := map[string]interface{}{}

		if v, exists := obj["hour"]; exists {
			row["hour"] = v
		}
		if v, exists := obj["wd_label"]; exists {
			row["direction_label"] = v
		}

		var windValue interface{}
		var gustValue interface{}

		if unit == "bft" {
			windValue = obj["ws_bft"]
			if v, exists := obj["wg_bft"]; exists && v != nil {
				gustValue = v
			} else {
				gustValue = obj["ws_bft"]
			}
		} else {
			windValue = obj["ws"]
			if v, exists := obj["wg"]; exists && v != nil {
				gustValue = v
			} else {
				gustValue = obj["ws"]
			}
		}

		row["wind_value"] = windValue
		row["gust_value"] = gustValue

		out = append(out, row)
	}

	env.Resolved["items"] = out
	return env
}
