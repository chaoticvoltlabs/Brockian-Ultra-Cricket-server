package ui

import (
	"time"

	"buc/internal/support"
)

func stringOption(m map[string]interface{}, key, fallback string) string {
	if v, ok := m[key].(string); ok && v != "" {
		return v
	}
	return fallback
}

func intOption(m map[string]interface{}, key string, fallback int) int {
	if v, ok := support.ToFloat(m[key]); ok {
		return int(v)
	}
	return fallback
}

func intFromAnyUI(v interface{}) (int, bool) {
	if f, ok := support.ToFloat(v); ok {
		return int(f), true
	}
	return 0, false
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
