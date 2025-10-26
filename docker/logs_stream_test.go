package docker

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/tim0-12432/simple-test-server/db/dtos"
)

// mockStreamReader implements StreamReader for testing.
type mockStreamReader struct {
	data []byte
	pos  int
}

func (m *mockStreamReader) Read(p []byte) (int, error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	n := copy(p, m.data[m.pos:])
	m.pos += n
	return n, nil
}

// mockStreamCommand implements StreamCommandRunner for testing.
type mockStreamCommand struct {
	stdout *mockStreamReader
	stderr *mockStreamReader
	err    error
}

func (m *mockStreamCommand) Start() error {
	return nil
}

func (m *mockStreamCommand) StdoutPipe() (StreamReader, error) {
	return m.stdout, nil
}

func (m *mockStreamCommand) StderrPipe() (StreamReader, error) {
	return m.stderr, nil
}

func (m *mockStreamCommand) Wait() error {
	return m.err
}

func TestStreamContainerLogs_Success(t *testing.T) {
	// Save original functions
	origExecStreamFn := ExecCommandStreamFn
	origGetContainerFn := getContainerFn
	defer func() {
		ExecCommandStreamFn = origExecStreamFn
		getContainerFn = origGetContainerFn
	}()

	// Mock container lookup
	getContainerFn = func(id string) (*dtos.Container, error) {
		return &dtos.Container{
			ID:     "test-id",
			Name:   "test-otel-container",
			Type:   "OTEL",
			Status: dtos.Running,
		}, nil
	}

	// Mock streaming command that outputs a few lines
	stdoutData := "2024-01-01T10:00:00Z log line 1\n2024-01-01T10:00:01Z log line 2\n"
	stderrData := "2024-01-01T10:00:02Z error line\n"

	ExecCommandStreamFn = func(ctx context.Context, name string, args ...string) StreamCommandRunner {
		// verify command
		if name != "docker" {
			t.Errorf("expected command 'docker', got %q", name)
		}
		if len(args) < 3 || args[0] != "logs" || args[1] != "-f" {
			t.Errorf("expected args starting with 'logs -f', got %v", args)
		}

		return &mockStreamCommand{
			stdout: &mockStreamReader{data: []byte(stdoutData)},
			stderr: &mockStreamReader{data: []byte(stderrData)},
			err:    nil,
		}
	}

	// Collect lines
	var lines []string
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := StreamContainerLogs(ctx, "test-id", func(line string) {
		lines = append(lines, line)
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify we got lines
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d: %v", len(lines), lines)
	}

	// Check content (order is not guaranteed due to concurrent stdout/stderr reading)
	allLines := strings.Join(lines, "\n")
	expectedContents := []string{"log line 1", "log line 2", "error line"}
	for _, expected := range expectedContents {
		if !strings.Contains(allLines, expected) {
			t.Errorf("expected output to contain %q, got lines: %v", expected, lines)
		}
	}
}

func TestStreamContainerLogs_ContainerNotFound(t *testing.T) {
	// Save original functions
	origGetContainerFn := getContainerFn
	defer func() {
		getContainerFn = origGetContainerFn
	}()

	// Mock container lookup to return error
	getContainerFn = func(id string) (*dtos.Container, error) {
		return nil, errors.New("container not found")
	}

	ctx := context.Background()
	err := StreamContainerLogs(ctx, "nonexistent", func(line string) {})

	if err != ErrContainerNotFound {
		t.Errorf("expected ErrContainerNotFound, got %v", err)
	}
}

func TestStreamContainerLogs_ContextCancelled(t *testing.T) {
	// Save original functions
	origExecStreamFn := ExecCommandStreamFn
	origGetContainerFn := getContainerFn
	defer func() {
		ExecCommandStreamFn = origExecStreamFn
		getContainerFn = origGetContainerFn
	}()

	// Mock container lookup
	getContainerFn = func(id string) (*dtos.Container, error) {
		return &dtos.Container{
			ID:     "test-id",
			Name:   "test-container",
			Type:   "OTEL",
			Status: dtos.Running,
		}, nil
	}

	// Mock streaming command that never ends
	ExecCommandStreamFn = func(ctx context.Context, name string, args ...string) StreamCommandRunner {
		// infinite stream - will be stopped by context cancellation
		return &mockStreamCommand{
			stdout: &mockStreamReader{data: []byte{}}, // empty stream
			stderr: &mockStreamReader{data: []byte{}},
			err:    nil,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	// cancel immediately
	cancel()

	err := StreamContainerLogs(ctx, "test-id", func(line string) {})

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
