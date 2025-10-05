package docker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tim0-12432/simple-test-server/db/dtos"
	"github.com/tim0-12432/simple-test-server/db/services"
)

// MaxLineLen limits the size of a single log line returned.
const MaxLineLen = 8192

// ExecCommandContextFn is the signature used to run external commands. It returns combined stdout/stderr
// and an error if the command failed. Tests may replace this variable to inject output.
var ExecCommandContextFn = func(ctx context.Context, name string, args ...string) ([]byte, error) {
	// default implementation uses CombinedOutput
	cmd := execCommand(ctx, name, args...)
	return cmd.CombinedOutput()
}

// execCommand is a tiny wrapper so the default implementation can be replaced in tests if needed.
var execCommand = func(ctx context.Context, name string, args ...string) *execCmdWrapper {
	return &execCmdWrapper{ctx: ctx, name: name, args: args}
}

// execCmdWrapper abstracts exec.CommandContext to avoid importing os/exec in tests that don't need it.
// It is implemented using the real os/exec package in production via the default execCommand above.
// In tests we can replace ExecCommandContextFn to a stub and don't need to touch this.

// FetchContainerLogs runs `docker logs` for the given container name and returns the last `tail` lines.
// If `since` is non-nil, it will add the --since option. Returns a slice of LogLine ordered oldest->newest.
func FetchContainerLogs(ctx context.Context, containerName string, tail int, since *time.Time) ([]dtos.LogLine, bool, error) {
	if tail <= 0 {
		tail = 500
	}
	if tail > 5000 {
		tail = 5000
	}

	// Validate container exists in our registry
	container, err := getContainerFn(containerName)
	if err != nil {
		return nil, false, ErrContainerNotFound
	}

	// build command
	args := []string{"logs", "--tail", fmt.Sprint(tail)}
	if since != nil {
		args = append(args, "--since", since.UTC().Format(time.RFC3339))
	}
	args = append(args, container.Name)

	// run with timeout
	runCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	out, err := ExecCommandContextFn(runCtx, "docker", args...)
	if err != nil {
		low := strings.ToLower(string(out))
		if strings.Contains(low, "no such container") {
			return nil, false, ErrContainerNotFound
		}
		return nil, false, fmt.Errorf("docker logs failed: %w", err)
	}

	lines, truncated, err := parseLogOutput(out, tail, nil)
	if err != nil {
		return nil, false, err
	}

	// Determine running state
	if strings.ToUpper(container.Type) != "WEB" {
		return nil, false, errors.New("container is not a web server")
	}

	if container.Status != dtos.Running {
		return lines, truncated, ErrContainerNotRunning
	}

	return lines, truncated, nil
}

// getContainerFn is a variable so tests can override lookup behavior.
var getContainerFn = services.GetContainer

// parseLogOutput parses combined output bytes from docker logs and returns parsed LogLine slice and truncated flag.
// If nowFunc is nil, time.Now will be used for fallback timestamps. nowFunc is provided for tests.
func parseLogOutput(output []byte, tail int, nowFunc func() time.Time) ([]dtos.LogLine, bool, error) {
	if nowFunc == nil {
		nowFunc = time.Now
	}
	reader := strings.NewReader(string(output))
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, MaxLineLen)

	lines := make([]dtos.LogLine, 0)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) > MaxLineLen {
			text = text[:MaxLineLen]
		}
		// try to parse common timestamp formats at line start
		ts := nowFunc().UTC()
		first := ""
		fields := strings.Fields(text)
		if len(fields) > 0 {
			first = fields[0]
		}
		if first != "" {
			// try a few layouts
			layouts := []string{time.RFC3339, time.RFC3339Nano, "2006-01-02 15:04:05", "2006/01/02 15:04:05"}
			for _, layout := range layouts {
				if t, err := time.Parse(layout, first); err == nil {
					ts = t.UTC()
					break
				}
			}
		}
		lines = append(lines, dtos.LogLine{TS: ts, Line: text})
	}
	if err := scanner.Err(); err != nil {
		return nil, false, fmt.Errorf("failed scanning logs: %w", err)
	}

	truncated := false
	if len(lines) >= tail {
		truncated = true
	}
	return lines, truncated, nil
}
