package docker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os/exec"
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

// StreamCommandRunner is the interface for running streaming commands. Tests can inject a mock.
type StreamCommandRunner interface {
	Start() error
	StdoutPipe() (StreamReader, error)
	StderrPipe() (StreamReader, error)
	Wait() error
}

// StreamReader abstracts reading from command output streams.
type StreamReader interface {
	Read(p []byte) (n int, error error)
}

// ExecCommandStreamFn is injected to start a streaming command. Tests replace this to inject mock behavior.
var ExecCommandStreamFn = func(ctx context.Context, name string, args ...string) StreamCommandRunner {
	return &realCommandRunner{
		cmd: exec.CommandContext(ctx, name, args...),
	}
}

// realCommandRunner wraps exec.Cmd to implement StreamCommandRunner.
type realCommandRunner struct {
	cmd *exec.Cmd
}

func (r *realCommandRunner) Start() error {
	return r.cmd.Start()
}

func (r *realCommandRunner) StdoutPipe() (StreamReader, error) {
	return r.cmd.StdoutPipe()
}

func (r *realCommandRunner) StderrPipe() (StreamReader, error) {
	return r.cmd.StderrPipe()
}

func (r *realCommandRunner) Wait() error {
	return r.cmd.Wait()
}

// StreamContainerLogs streams docker logs for a container in real-time using `docker logs -f`.
// It sends each log line to the onLine callback. Blocks until context is cancelled or command finishes.
// Returns error if container not found or command fails to start.
func StreamContainerLogs(ctx context.Context, containerID string, onLine func(line string)) error {
	// Validate container exists in our registry
	container, err := getContainerFn(containerID)
	if err != nil {
		return ErrContainerNotFound
	}

	// build command: docker logs -f --tail=50 <containerName>
	args := []string{"logs", "-f", "--tail", "50", container.Name}

	// start streaming command
	cmd := ExecCommandStreamFn(ctx, "docker", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start docker logs: %w", err)
	}

	// goroutine to read stdout
	stdoutDone := make(chan struct{})
	go func() {
		defer close(stdoutDone)
		scanner := bufio.NewScanner(stdout)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, MaxLineLen)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				line := scanner.Text()
				if len(line) > MaxLineLen {
					line = line[:MaxLineLen]
				}
				onLine(line)
			}
		}
	}()

	// goroutine to read stderr
	stderrDone := make(chan struct{})
	go func() {
		defer close(stderrDone)
		scanner := bufio.NewScanner(stderr)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, MaxLineLen)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				line := scanner.Text()
				if len(line) > MaxLineLen {
					line = line[:MaxLineLen]
				}
				onLine(line)
			}
		}
	}()

	// wait for context cancellation or command completion
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// context cancelled, wait for goroutines to finish
		<-stdoutDone
		<-stderrDone
		return ctx.Err()
	case err := <-done:
		// command finished
		<-stdoutDone
		<-stderrDone
		if err != nil {
			// check if container not found
			if strings.Contains(err.Error(), "No such container") {
				return ErrContainerNotFound
			}
			return fmt.Errorf("docker logs command failed: %w", err)
		}
		return nil
	}
}
