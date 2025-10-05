package docker

import (
	"context"
	"os/exec"
)

// execCmdWrapper uses the real os/exec.CommandContext to run commands.
type execCmdWrapper struct {
	ctx  context.Context
	name string
	args []string
}

// CombinedOutput runs the command and returns combined stdout/stderr.
func (e *execCmdWrapper) CombinedOutput() ([]byte, error) {
	cmd := exec.CommandContext(e.ctx, e.name, e.args...)
	return cmd.CombinedOutput()
}

// Note: tests can replace ExecCommandContextFn with a stub that returns the desired output.
