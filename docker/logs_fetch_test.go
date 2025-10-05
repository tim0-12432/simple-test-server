package docker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tim0-12432/simple-test-server/db/dtos"
)

func TestFetchContainerLogs_Success(t *testing.T) {
	oldGet := getContainerFn
	oldExec := ExecCommandContextFn
	defer func() {
		getContainerFn = oldGet
		ExecCommandContextFn = oldExec
	}()

	getContainerFn = func(name string) (*dtos.Container, error) {
		return &dtos.Container{
			Name:   name,
			Status: dtos.Running,
			Type:   "web",
		}, nil
	}

	ExecCommandContextFn = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		// simulate docker logs output
		return []byte("2025-10-05T11:59:58Z first-line\nsecond-line-without-ts\n"), nil
	}

	ctx := context.Background()
	lines, truncated, err := FetchContainerLogs(ctx, "my-container", 500, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if truncated {
		t.Fatalf("expected truncated=false")
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].Line != "2025-10-05T11:59:58Z first-line" {
		t.Fatalf("unexpected first line: %q", lines[0].Line)
	}
	if lines[1].Line != "second-line-without-ts" {
		t.Fatalf("unexpected second line: %q", lines[1].Line)
	}
}

func TestFetchContainerLogs_NotFoundByLookup(t *testing.T) {
	oldGet := getContainerFn
	oldExec := ExecCommandContextFn
	defer func() {
		getContainerFn = oldGet
		ExecCommandContextFn = oldExec
	}()

	getContainerFn = func(name string) (*dtos.Container, error) {
		return nil, errors.New("not found")
	}

	// Exec shouldn't be called, but stub to be safe
	ExecCommandContextFn = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return nil, nil
	}

	ctx := context.Background()
	_, _, err := FetchContainerLogs(ctx, "missing", 10, nil)
	if !errors.Is(err, ErrContainerNotFound) {
		t.Fatalf("expected ErrContainerNotFound, got %v", err)
	}
}

func TestFetchContainerLogs_NotRunning(t *testing.T) {
	oldGet := getContainerFn
	oldExec := ExecCommandContextFn
	defer func() {
		getContainerFn = oldGet
		ExecCommandContextFn = oldExec
	}()

	getContainerFn = func(name string) (*dtos.Container, error) {
		return &dtos.Container{
			Name:   name,
			Status: dtos.Discarded,
			Type:   "web",
		}, nil
	}

	ExecCommandContextFn = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return []byte("line-a\nline-b\n"), nil
	}

	ctx := context.Background()
	lines, truncated, err := FetchContainerLogs(ctx, "stopped", 100, nil)
	if !errors.Is(err, ErrContainerNotRunning) {
		t.Fatalf("expected ErrContainerNotRunning, got %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines returned even when not running, got %d", len(lines))
	}
	_ = truncated
}

func TestFetchContainerLogs_DockerNoSuchContainer(t *testing.T) {
	oldGet := getContainerFn
	oldExec := ExecCommandContextFn
	defer func() {
		getContainerFn = oldGet
		ExecCommandContextFn = oldExec
	}()

	getContainerFn = func(name string) (*dtos.Container, error) {
		return &dtos.Container{
			Name:   name,
			Status: dtos.Running,
			Type:   "web",
		}, nil
	}

	ExecCommandContextFn = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		return []byte("Error: No such container: foo"), errors.New("exit status 1")
	}

	ctx := context.Background()
	_, _, err := FetchContainerLogs(ctx, "my-container", 10, nil)
	if !errors.Is(err, ErrContainerNotFound) {
		t.Fatalf("expected ErrContainerNotFound from docker output, got %v", err)
	}
}

func TestFetchContainerLogs_WithSinceOption(t *testing.T) {
	oldGet := getContainerFn
	oldExec := ExecCommandContextFn
	defer func() {
		getContainerFn = oldGet
		ExecCommandContextFn = oldExec
	}()

	capturedArgs := []string{}
	getContainerFn = func(name string) (*dtos.Container, error) {
		return &dtos.Container{Name: name, Status: dtos.Running, Type: "web"}, nil
	}

	ExecCommandContextFn = func(ctx context.Context, name string, args ...string) ([]byte, error) {
		capturedArgs = append(capturedArgs, args...)
		return []byte("ok\n"), nil
	}

	since := time.Date(2025, 10, 5, 10, 0, 0, 0, time.UTC)
	_, _, err := FetchContainerLogs(context.Background(), "c", 10, &since)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// ensure --since was passed
	found := false
	for i := 0; i < len(capturedArgs)-1; i++ {
		if capturedArgs[i] == "--since" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected --since in docker args, got %v", capturedArgs)
	}
}
