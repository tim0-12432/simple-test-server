package ftp

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/tim0-12432/simple-test-server/docker"
)

// ListDirectory lists files inside FTP container relative path
func ListDirectory(ctx context.Context, containerId string, relPath string, maxEntries int) ([]docker.FileEntry, bool, error) {
	entries, truncated, err := docker.ListFtpDir(ctx, containerId, relPath, maxEntries)
	if err != nil {
		return nil, false, fmt.Errorf("list directory %q: %w", relPath, err)
	}
	return entries, truncated, nil
}

// GetLogs retrieves container logs
func GetLogs(ctx context.Context, containerId string, tail int) (string, error) {
	out, err := docker.GetContainerLogs(ctx, containerId, tail)
	if err != nil {
		return "", fmt.Errorf("get logs: %w", err)
	}
	return out, nil
}

// UploadFileToContainer copies a local file into the container's FTP home (/home/user)
func UploadFileToContainer(ctx context.Context, containerId string, destRelPath string, localPath string) error {
	if containerId == "" {
		return fmt.Errorf("container id empty")
	}

	// sanitize destination: do not allow absolute path
	if filepath.IsAbs(destRelPath) {
		return fmt.Errorf("destination path must be relative")
	}
	if destRelPath == "" {
		return fmt.Errorf("destination path empty")
	}

	destination := filepath.Join("/home/user", destRelPath)
	// give a short timeout for copy
	if err := docker.CopyFileToContainer(ctx, containerId, localPath, destination, 30*time.Second); err != nil {
		return fmt.Errorf("copy file to container: %w", err)
	}
	return nil
}
