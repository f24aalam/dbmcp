package projectscan

import (
	"encoding/json"
	"strconv"
)

// flattenJSONStringMap unmarshals JSON and collects string leaf values keyed by a
// normalized path (e.g. spring.datasource.url -> SPRING_DATASOURCE_URL).
func flattenJSONStringMap(data []byte) map[string]string {
	var root any
	if err := json.Unmarshal(data, &root); err != nil {
		return nil
	}

	out := make(map[string]string)
	var walk func(prefix string, v any)
	walk = func(prefix string, v any) {
		switch t := v.(type) {
		case map[string]any:
			for k, v2 := range t {
				next := k
				if prefix != "" {
					next = prefix + "_" + k
				}
				walk(next, v2)
			}
		case []any:
			for i, v2 := range t {
				next := prefix + "_" + strconv.Itoa(i)
				walk(next, v2)
			}
		case string:
			if prefix != "" {
				out[normalizeKey(prefix)] = t
			}
		case float64, bool, nil:
			// ignore non-string leaves
		default:
			// json.Number etc.
		}
	}
	walk("", root)
	if len(out) == 0 {
		return nil
	}
	return out
}
