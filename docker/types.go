package docker

import (
	"errors"
	"time"
)

// sentinel errors for docker operations
var (
	ErrContainerNotFound   = errors.New("container not found")
	ErrContainerNotRunning = errors.New("container not running")
)

// FileEntry represents a file or directory inside a container
type FileEntry struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Type       string    `json:"type"` // file | dir | symlink
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modifiedAt"`
}
