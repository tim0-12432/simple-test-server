package common

import (
	"errors"
	"github.com/tim0-12432/simple-test-server/config"
)

// defaultUploadSize is the compile-time default (10 MB) used when no env config is set.
const defaultUploadSize int64 = 10 << 20 // 10 MB

// GetMaxUploadSize returns the configured maximum upload size in bytes.
// It prefers the value from `config.EnvConfig.UploadMaxBytes` when available,
// otherwise falling back to the compile-time default.
func GetMaxUploadSize() int64 {
	if config.EnvConfig != nil && config.EnvConfig.UploadMaxBytes > 0 {
		return config.EnvConfig.UploadMaxBytes
	}
	return defaultUploadSize
}

// AllowedMIMEs lists permitted MIME types for uploads.
var AllowedMIMEs = map[string]struct{}{
	"image/jpeg":      {},
	"image/png":       {},
	"image/gif":       {},
	"application/pdf": {},
	"text/plain":      {},
	"text/html":       {},
}

// Package-level errors returned by the upload service.
var (
	ErrMissingFile   = errors.New("missing file")
	ErrTooLarge      = errors.New("file too large")
	ErrInvalidType   = errors.New("file type not allowed")
	ErrSaveFailed    = errors.New("failed to save uploaded file")
	ErrReopenFailure = errors.New("failed to reopen uploaded file")
)
