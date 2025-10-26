package db

import (
	"reflect"
	"testing"
)

func TestToIntMap(t *testing.T) {
	mapStringInt := map[string]int{"8025": 8025, "1025": 1025}

	cases := []struct {
		name string
		in   any
		want map[int]int
	}{
		{"empty map", make(map[string]int), map[int]int{}},
		{"map[string]int", mapStringInt, map[int]int{8025: 8025, 1025: 1025}},
		{"invalid string", "not-a-map", map[int]int{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ToIntMap(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("case %q: got %#v, want %#v", tc.name, got, tc.want)
			}
		})
	}
}

func TestToStringMap(t *testing.T) {
	mapStringString := map[string]string{"key1": "value1", "key2": "value2"}

	cases := []struct {
		name string
		in   any
		want map[string]string
	}{
		{"empty map", make(map[string]string), map[string]string{}},
		{"map[string]string", mapStringString, map[string]string{"key1": "value1", "key2": "value2"}},
		{"invalid string", "not-a-map", map[string]string{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ToStringMap(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("case %q: got %#v, want %#v", tc.name, got, tc.want)
			}
		})
	}
}
