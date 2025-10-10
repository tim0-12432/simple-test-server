package docker

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ListSmbDir lists the immediate children of the given relative path
// inside the container's SMB share root (/shares). relPath must be relative
// (no leading '/') and must not contain path traversal ("..").
//
// Returns the entries, a boolean indicating whether the result was truncated
// due to the maxEntries limit, and an error if something went wrong.
func ListSmbDir(ctx context.Context, containerId string, relPath string, maxEntries int) ([]FileEntry, bool, error) {
	if containerId == "" {
		return nil, false, fmt.Errorf("container id empty")
	}

	// basic validation: disallow absolute paths and traversal
	if strings.HasPrefix(relPath, "/") {
		return nil, false, fmt.Errorf("relPath must be relative")
	}
	if strings.Contains(relPath, "..") {
		return nil, false, fmt.Errorf("relPath must not contain '..'")
	}

	// verify container exists by running inspect
	if err := func() error {
		ci := exec.Command("docker", "inspect", containerId)
		if err := ci.Run(); err != nil {
			return fmt.Errorf("container not found: %s", containerId)
		}
		return nil
	}(); err != nil {
		return nil, false, err
	}

	smbroot := "/shares"
	var target string
	if relPath == "" {
		target = smbroot
	} else {
		target = filepath.Join(smbroot, relPath)
	}

	// set a conservative timeout for the exec
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Use find to list immediate children with metadata (type|size|mtime|path)
	cmd := exec.CommandContext(ctx, "docker", "exec", containerId, "find", target, "-maxdepth", "1", "-mindepth", "1", "-printf", "%y|%s|%T@|%p\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		outStr := strings.TrimSpace(string(out))
		// fallback for minimal find implementations that don't support -printf
		if strings.Contains(strings.ToLower(outStr), "unrecognized -printf") || strings.Contains(strings.ToLower(outStr), "unknown primary") || strings.Contains(strings.ToLower(outStr), "-printf") {
			fallback := `for f in ` + target + `/* ; do [ -e "$f" ] || continue ; t='?'; if [ -d "$f" ]; then t='d'; elif [ -L "$f" ]; then t='l'; elif [ -f "$f" ]; then t='f'; fi; s=0; if [ -f "$f" ]; then s=$(wc -c <"$f" 2>/dev/null || echo 0); fi; printf "%s|%s|%s|%s\n" "$t" "$s" "0.0" "$f"; done`
			fbCmd := exec.CommandContext(ctx, "docker", "exec", containerId, "sh", "-c", fallback)
			out2, err2 := fbCmd.CombinedOutput()
			if err2 != nil {
				return nil, false, fmt.Errorf("docker exec fallback listing failed: %v - %s", err2, strings.TrimSpace(string(out2)))
			}
			raw := strings.TrimSpace(string(out2))
			if raw == "" {
				return []FileEntry{}, false, nil
			}
			entries, truncated := parseFindOutput(raw, smbroot, maxEntries)
			return entries, truncated, nil
		}
		return nil, false, fmt.Errorf("docker exec find failed: %v - %s", err, outStr)
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return []FileEntry{}, false, nil
	}

	entries, truncated := parseFindOutput(raw, smbroot, maxEntries)
	return entries, truncated, nil
}
