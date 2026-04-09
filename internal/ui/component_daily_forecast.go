package ui

import (
	"time"

	"buc/internal/config"
)

func BuildDailyForecast(
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

	daysLimit := intOption(options, "days", 4)
	windMode := stringOption(options, "wind_mode", "gust_max")
	windUnit := stringOption(options, "wind_unit", "bft")

	env.Resolved["display_wind_label"] = "Wind"
	env.Resolved["display_wind_unit"] = windUnit
	env.Resolved["display_wind_field"] = "wg_bft_max"
	env.Resolved["display_gust_label"] = "Piek"

	if windMode == "wind_max" && windUnit == "bft" {
		env.Resolved["display_wind_label"] = "Wind"
		env.Resolved["display_wind_unit"] = "bft"
		env.Resolved["display_wind_field"] = "ws_bft_max"
	} else if windMode == "wind_max" && windUnit == "kmh" {
		env.Resolved["display_wind_label"] = "Wind"
		env.Resolved["display_wind_unit"] = "km/h"
		env.Resolved["display_wind_field"] = "ws_max"
	} else if windMode == "gust_max" && windUnit == "kmh" {
		env.Resolved["display_wind_label"] = "Stoten"
		env.Resolved["display_wind_unit"] = "km/h"
		env.Resolved["display_wind_field"] = "wg_max"
	}

	rawDays, ok := data["days"].([]interface{})
	if !ok {
		env.Status.OK = false
		env.Status.Warnings = append(env.Status.Warnings, "days missing or not an array")
		env.Resolved["items"] = []interface{}{}
		return env
	}

	out := make([]interface{}, 0, daysLimit)
	for i, item := range rawDays {
		if i >= daysLimit {
			break
		}

		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		row := map[string]interface{}{}

		if dateStr, ok := obj["date"].(string); ok {
			row["date"] = dateStr
			if _, exists := obj["day_label"]; !exists {
				if t, err := time.Parse("2006-01-02", dateStr); err == nil {
					row["day_label"] = weekdayLabelNL(t.Weekday())
				}
			}
		}

		for _, k := range []string{"day_label", "weather_code", "t_min", "t_max", "pop"} {
			if v, exists := obj[k]; exists {
				row[k] = v
			}
		}

		row["icon_key"] = iconKeyFromWeatherCode(obj["weather_code"])
		if windUnit == "bft" {
			row["wind_value"] = obj["ws_bft_max"]
			if v, exists := obj["wg_bft_max"]; exists && v != nil {
				row["gust_value"] = v
			} else {
				row["gust_value"] = obj["ws_bft_max"]
			}
		} else {
			row["wind_value"] = obj["ws_max"]
			if v, exists := obj["wg_max"]; exists && v != nil {
				row["gust_value"] = v
			} else {
				row["gust_value"] = obj["ws_max"]
			}
		}

		out = append(out, row)
	}

	env.Resolved["items"] = out
	return env
}

func iconKeyFromWeatherCode(v interface{}) string {
	code, ok := intFromAnyUI(v)
	if !ok {
		return "cloudy"
	}

	switch code {
	case 0:
		return "sunny"
	case 1, 2:
		return "partly_cloudy"
	case 3:
		return "cloudy"
	case 45, 48:
		return "fog"
	case 51, 53, 55, 61, 63, 80:
		return "rain"
	case 65, 81, 82:
		return "pouring"
	case 56, 57, 66, 67:
		return "sleet"
	case 71, 73, 75, 77, 85, 86:
		return "snow"
	case 95, 96, 99:
		return "thunder"
	default:
		return "cloudy"
	}
}
