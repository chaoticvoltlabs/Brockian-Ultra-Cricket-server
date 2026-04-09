package ha

import (
	"encoding/json"
	"strings"
)

func NormalizeEntity(e *EntityState) *EntityState {
	if e == nil {
		return nil
	}
	e.Attributes = normalizeMap(e.Attributes)
	return e
}

func normalizeMap(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = normalizeValue(v)
	}
	return out
}

func normalizeSlice(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = normalizeValue(v)
	}
	return out
}

func normalizeValue(v interface{}) interface{} {
	switch x := v.(type) {
	case map[string]interface{}:
		return normalizeMap(x)
	case []interface{}:
		return normalizeSlice(x)
	case string:
		s := strings.TrimSpace(x)
		if len(s) > 0 && ((s[0] == '{' && s[len(s)-1] == '}') || (s[0] == '[' && s[len(s)-1] == ']')) {
			var parsed interface{}
			if err := json.Unmarshal([]byte(s), &parsed); err == nil {
				return normalizeValue(parsed)
			}
		}
		return x
	default:
		return v
	}
}
