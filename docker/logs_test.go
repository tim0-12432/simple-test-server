package docker

import (
	"strings"
	"testing"
	"time"
)

func TestParseLogOutput_Basic(t *testing.T) {
	now := time.Date(2025, 10, 5, 12, 0, 0, 0, time.UTC)
	out := []byte("2025-10-05T11:59:58Z GET / 200\nno-ts-line here\n")
	lines, truncated, err := parseLogOutput(out, 500, func() time.Time { return now })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if truncated {
		t.Fatalf("expected truncated=false")
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !lines[0].TS.Equal(time.Date(2025, 10, 5, 11, 59, 58, 0, time.UTC)) {
		t.Fatalf("expected parsed timestamp, got %v", lines[0].TS)
	}
	if strings.TrimSpace(lines[1].Line) != "no-ts-line here" {
		t.Fatalf("unexpected content: %q", lines[1].Line)
	}
}

func TestParseLogOutput_Truncation(t *testing.T) {
	// create many lines to hit tail boundary
	var b strings.Builder
	for i := 0; i < 100; i++ {
		b.WriteString("line-\n")
	}
	out := []byte(b.String())
	lines, truncated, err := parseLogOutput(out, 50, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !truncated {
		t.Fatalf("expected truncated=true when lines >= tail")
	}
	if len(lines) != 100 {
		t.Fatalf("expected parsed all lines, got %d", len(lines))
	}
}
