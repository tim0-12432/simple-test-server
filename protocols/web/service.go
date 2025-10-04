package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// UploadResult contains metadata about the saved uploaded file.
type UploadResult struct {
	LocalPath   string
	SafeName    string
	Size        int64
	ContentType string
}

// SaveUploadedFileToTmp validates the multipart file header, sniffs the content type, enforces size limits, saves to a temp file and returns metadata.
func SaveUploadedFileToTmp(ctx context.Context, fh *multipart.FileHeader) (*UploadResult, error) {
	if fh == nil {
		return nil, ErrMissingFile
	}

	upFile, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer upFile.Close()

	// read first 512 bytes to sniff content type
	head := make([]byte, 512)
	n, _ := upFile.Read(head)
	contentType := detectContentType(head[:n])
	if _, ok := AllowedMIMEs[contentType]; !ok {
		return nil, ErrInvalidType
	}

	// build safe name
	safeName := buildSafeName(fh.Filename, contentType)

	// rebuild reader that yields the header we consumed + remainder
	fullReader := io.MultiReader(bytes.NewReader(head[:n]), upFile)
	// enforce max size
	limited := io.LimitReader(fullReader, MaxUploadSize+1)

	tmpDir := os.TempDir()
	localPath := filepath.Join(tmpDir, safeName)
	outF, err := os.Create(localPath)
	if err != nil {
		return nil, ErrSaveFailed
	}
	defer func() {
		_ = outF.Close()
	}()

	written, err := io.Copy(outF, limited)
	if err != nil {
		_ = os.Remove(localPath)
		return nil, ErrSaveFailed
	}

	if written > MaxUploadSize {
		_ = os.Remove(localPath)
		return nil, ErrTooLarge
	}

	return &UploadResult{
		LocalPath:   localPath,
		SafeName:    safeName,
		Size:        written,
		ContentType: contentType,
	}, nil
}
