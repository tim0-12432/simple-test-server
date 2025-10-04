package docker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// CopyFileToContainer copies a file from the host into the container using `docker cp` with a timeout.
func CopyFileToContainer(ctx context.Context, containerName, srcPath, destPath string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// docker cp <src> <container>:<dest>
	destination := fmt.Sprintf("%s:%s", containerName, destPath)
	cmd := exec.CommandContext(ctx, "docker", "cp", srcPath, destination)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker cp failed: %w - %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
