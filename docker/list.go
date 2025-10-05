package docker

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ListContainerDir lists the immediate children of the given relative path
// inside the container's webroot. relPath must be relative (no leading '/')
// and must not contain path traversal ("..").
//
// Returns the entries, a boolean indicating whether the result was truncated
// due to the maxEntries limit, and an error if something went wrong.
func ListContainerDir(ctx context.Context, containerId string, relPath string, maxEntries int) ([]FileEntry, bool, error) {
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

	webroot := "/usr/share/nginx/html"
	var target string
	if relPath == "" {
		target = webroot
	} else {
		target = filepath.Join(webroot, relPath)
	}

	// set a conservative timeout for the exec
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Use find to list immediate children with metadata (type|size|mtime|path)
	// Note: This requires `find` in the container to support -printf. Most
	// standard GNU find implementations do; some minimal images may not.
	cmd := exec.CommandContext(ctx, "docker", "exec", containerId, "find", target, "-maxdepth", "1", "-mindepth", "1", "-printf", "%y|%s|%T@|%p\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// return output to help debugging, but wrap it
		return nil, false, fmt.Errorf("docker exec find failed: %v - %s", err, strings.TrimSpace(string(out)))
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return []FileEntry{}, false, nil
	}

	entries, truncated := parseFindOutput(raw, webroot, maxEntries)
	return entries, truncated, nil
}

// parseFindOutput parses the raw output from find -printf (one entry per line)
// format: "%y|%s|%T@|%p\n" and returns FileEntry slice and truncated flag.
func parseFindOutput(raw string, webroot string, maxEntries int) ([]FileEntry, bool) {
	lines := strings.Split(raw, "\n")
	truncated := false
	entries := make([]FileEntry, 0, len(lines))

	limit := maxEntries
	if limit <= 0 {
		limit = 500
	}

	count := 0
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			// skip malformed line
			continue
		}
		if count >= limit {
			truncated = true
			break
		}
		typeChar := parts[0]
		sizeStr := parts[1]
		mtimeStr := parts[2]
		pathStr := parts[3]

		sz, _ := strconv.ParseInt(sizeStr, 10, 64)
		mtf, _ := strconv.ParseFloat(mtimeStr, 64)
		sec := int64(mtf)
		nsec := int64((mtf - float64(sec)) * 1e9)
		modTime := time.Unix(sec, nsec).UTC()

		var typ string
		switch typeChar {
		case "f":
			typ = "file"
		case "d":
			typ = "dir"
		case "l":
			typ = "symlink"
		default:
			typ = "unknown"
		}

		// compute relative path inside webroot
		rel := strings.TrimPrefix(pathStr, webroot)
		rel = strings.TrimPrefix(rel, "/")
		name := filepath.Base(pathStr)

		entries = append(entries, FileEntry{
			Name:       name,
			Path:       rel,
			Type:       typ,
			Size:       sz,
			ModifiedAt: modTime,
		})
		count++
	}

	return entries, truncated
}
