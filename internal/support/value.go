package support

import "strconv"

func ToFloat(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func Beaufort(kmh float64) int {
	switch {
	case kmh < 1:
		return 0
	case kmh < 6:
		return 1
	case kmh < 12:
		return 2
	case kmh < 20:
		return 3
	case kmh < 29:
		return 4
	case kmh < 39:
		return 5
	case kmh < 50:
		return 6
	case kmh < 62:
		return 7
	case kmh < 75:
		return 8
	case kmh < 89:
		return 9
	case kmh < 103:
		return 10
	case kmh < 118:
		return 11
	default:
		return 12
	}
}

func WindLabel16(bearing float64) string {
	dirs := []string{"N", "NNO", "NO", "ONO", "O", "OZO", "ZO", "ZZO", "Z", "ZZW", "ZW", "WZW", "W", "WNW", "NW", "NNW"}
	idx := int((bearing + 11.25) / 22.5)
	idx = idx % 16
	return dirs[idx]
}
