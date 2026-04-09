package ui

import (
	"fmt"

	"buc/internal/config"
	"buc/internal/support"
	"buc/internal/theme"
)

func BuildOutsideSummary(
	cfg *config.AllConfig,
	themeName string,
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

	tempMode := stringOption(options, "temperature_color_mode", "actual")
	tempScale := stringOption(options, "temperature_color_scale", "comfort_default")
	windUnit := stringOption(options, "wind_unit", "bft")

	var tempValue float64
	var ok bool

	switch tempMode {
	case "actual":
		tempValue, ok = support.ToFloat(data["t"])
	case "feels_like":
		tempValue, ok = support.ToFloat(data["feels_like"])
	default:
		ok = false
	}

	if ok {
		token, color, err := theme.ResolveTemperatureToken(cfg, themeName, tempScale, tempValue)
		if err != nil {
			env.Status.Warnings = append(env.Status.Warnings, err.Error())
		} else {
			env.Resolved["temperature_color_token"] = token
			env.Resolved["temperature_color"] = color
		}
	} else {
		env.Status.Warnings = append(env.Status.Warnings, "missing temperature value for color resolution")
	}

	if windUnit == "bft" {
		if v, exists := data["ws_bft"]; exists && v != nil {
			env.Resolved["wind_display_value"] = v
			env.Resolved["wind_display_unit"] = "bft"
		} else if v, exists := data["ws"]; exists && v != nil {
			env.Resolved["wind_display_value"] = v
			env.Resolved["wind_display_unit"] = "km/h"
		}

		if v, exists := data["wg_bft"]; exists && v != nil {
			env.Resolved["gust_display_value"] = v
			env.Resolved["gust_display_unit"] = "bft"
		} else if v, exists := data["wg"]; exists && v != nil {
			env.Resolved["gust_display_value"] = v
			env.Resolved["gust_display_unit"] = "km/h"
		}
	} else {
		if v, exists := data["ws"]; exists && v != nil {
			env.Resolved["wind_display_value"] = v
			env.Resolved["wind_display_unit"] = "km/h"
		}
		if v, exists := data["wg"]; exists && v != nil {
			env.Resolved["gust_display_value"] = v
			env.Resolved["gust_display_unit"] = "km/h"
		}
	}

	if label, lok := data["wd_label"].(string); lok {
		if wd, wok := support.ToFloat(data["wd"]); wok {
			env.Resolved["direction_human"] = fmt.Sprintf("%s (%.0f°)", label, wd)
		} else {
			env.Resolved["direction_human"] = label
		}
	}

	return env
}
