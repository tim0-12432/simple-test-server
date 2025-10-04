package web

import "errors"

// MaxUploadSize is the maximum allowed upload size in bytes (10 MB).
const MaxUploadSize int64 = 10 << 20 // 10 MB

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
