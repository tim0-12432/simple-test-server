package docker

import "time"

// FileEntry represents a file or directory inside a container
type FileEntry struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Type       string    `json:"type"` // file | dir | symlink
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modifiedAt"`
}
