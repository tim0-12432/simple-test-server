package dtos

import "time"

// LogLine represents a single log line with a timestamp and the raw text.
type LogLine struct {
	TS   time.Time `json:"ts"`
	Line string    `json:"line"`
}
