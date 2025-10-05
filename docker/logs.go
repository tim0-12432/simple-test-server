package docker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// GetContainerLogs returns the container logs (stdout/stderr) with timestamps using `docker logs`.
func GetContainerLogs(ctx context.Context, containerId string, tail int) (string, error) {
	if containerId == "" {
		return "", fmt.Errorf("container id empty")
	}

	if tail <= 0 {
		tail = 100
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	args := []string{"logs", "--timestamps", "--tail", fmt.Sprintf("%d", tail), containerId}
	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker logs failed: %v - %s", err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
