package web

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// sanitizeFilename produces a safe filename (without extension) derived from the original name.
func sanitizeFilename(name string) string {
	base := filepath.Base(name)
	ext := filepath.Ext(base)
	nameOnly := strings.TrimSuffix(base, ext)
	// remove disallowed characters
	re := regexp.MustCompile(`[^A-Za-z0-9._-]`)
	s := re.ReplaceAllString(nameOnly, "_")
	if len(s) > 64 {
		s = s[:64]
	}
	return s
}

// extensionForContentType returns a sensible extension for a detected content type.
func extensionForContentType(contentType string) string {
	switch contentType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "application/pdf":
		return ".pdf"
	case "text/plain":
		return ".txt"
	case "text/html":
		return ".html"
	default:
		return ""
	}
}

// buildSafeName creates a filename using sanitized name + extension.
func buildSafeName(rawFilename string, contentType string) string {
	ext := filepath.Ext(rawFilename)
	name := strings.TrimSuffix(filepath.Base(rawFilename), ext)
	if ext == "" {
		ext = extensionForContentType(contentType)
	}
	safe := sanitizeFilename(name)
	return fmt.Sprintf("%s%s", safe, ext)
}
