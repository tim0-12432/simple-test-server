package db

import (
	"encoding/json"
	"strconv"
	"strings"
)

func ToString(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case *string:
		if t == nil {
			return ""
		}
		return *t
	default:
		bs, _ := json.Marshal(v)
		return string(bs)
	}
}

func ToStringMap(v any) map[string]string {
	out := map[string]string{}
	if v == nil {
		return out
	}

	// try direct map[string]string
	if m, ok := v.(map[string]string); ok {
		return m
	}

	// try map[string]interface{}
	if m2, ok := v.(map[string]interface{}); ok {
		for k, vv := range m2 {
			out[k] = ToString(vv)
		}
		return out
	}

	// try JSON string
	if s, ok := v.(string); ok {
		var tmp map[string]string
		if err := json.Unmarshal([]byte(s), &tmp); err == nil {
			return tmp
		}
	}

	return out
}

func ToStringSlice(v any) []string {
	out := []string{}
	if v == nil {
		return out
	}
	if s, ok := v.([]string); ok {
		return s
	}
	if sa, ok := v.([]interface{}); ok {
		for _, e := range sa {
			out = append(out, ToString(e))
		}
		return out
	}
	if s, ok := v.(string); ok {
		// try decode json array
		var tmp []string
		if err := json.Unmarshal([]byte(s), &tmp); err == nil {
			return tmp
		}
		// comma separated fallback
		for _, p := range strings.Split(s, ",") {
			out = append(out, strings.TrimSpace(p))
		}
		return out
	}
	return out
}

func ToIntMap(v any) map[int]int {
	out := map[int]int{}
	if v == nil {
		return out
	}

	// try map[string]interface{} -> convert keys to int
	if m2, ok := v.(map[string]interface{}); ok {
		for ks, vv := range m2 {
			ki, err := strconv.Atoi(ks)
			if err != nil {
				// try float key
				kf, _ := strconv.ParseFloat(ks, 64)
				ki = int(kf)
			}
			vi := 0
			switch t := vv.(type) {
			case float64:
				vi = int(t)
			case int:
				vi = t
			case string:
				vi, _ = strconv.Atoi(t)
			}
			out[ki] = vi
		}
		return out
	}

	// try map[string]string JSON
	if m3, ok := v.(map[string]string); ok {
		for ks, vs := range m3 {
			ki, err := strconv.Atoi(ks)
			if err != nil {
				continue
			}
			vi, _ := strconv.Atoi(vs)
			out[ki] = vi
		}
		return out
	}

	// try JSON string
	if s, ok := v.(string); ok {
		var tmp map[string]int
		if err := json.Unmarshal([]byte(s), &tmp); err == nil {
			for ks, vi := range tmp {
				ki, err := strconv.Atoi(ks)
				if err != nil {
					continue
				}
				out[ki] = vi
			}
			return out
		}
	}

	return out
}

func EscapeSQL(s string) string {
	// quick path
	if s == "" {
		return s
	}

	// remove ASCII control characters (including null bytes)
	b := make([]rune, 0, len(s))
	for _, r := range s {
		if r >= 0x20 || r == '\n' || r == '\r' || r == '\t' {
			b = append(b, r)
		}
	}
	clean := string(b)

	// double single quotes for SQL string literal escaping
	clean = strings.ReplaceAll(clean, "'", "''")

	// trim to a reasonable length to avoid ridiculously long inputs
	if len(clean) > 2048 {
		return clean[:2048]
	}
	return clean
}
