package otel

import (
	"context"

	"github.com/tim0-12432/simple-test-server/docker"
)

// StreamOtelLogs streams real-time OTEL collector logs for the given container.
// Calls onLine for each log line received. Blocks until ctx is cancelled or stream ends.
func StreamOtelLogs(ctx context.Context, containerID string, onLine func(line string)) error {
	return docker.StreamContainerLogs(ctx, containerID, onLine)
}
