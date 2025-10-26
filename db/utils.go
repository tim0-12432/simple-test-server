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

func ToInt64(v any) int64 {
	if v == nil {
		return 0
	}
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case string:
		i, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0
		}
		return i
	default:
		bs, _ := json.Marshal(v)
		i, _ := strconv.ParseInt(string(bs), 10, 64)
		return i
	}
}

func ToStringMap(v any) map[string]string {
	out := map[string]string{}
	if v == nil {
		return out
	}
	if m, ok := v.(map[string]string); ok {
		return m
	}
	return out
}

func ToIntMap(v any) map[int]int {
	out := map[int]int{}
	if v == nil {
		return out
	}
	if msi, ok := v.(map[string]int); ok {
		for ks, vi := range msi {
			ki, err := strconv.Atoi(ks)
			if err != nil {
				// try float-ish key
				kf, _ := strconv.ParseFloat(ks, 64)
				ki = int(kf)
			}
			out[ki] = vi
		}
		return out
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
