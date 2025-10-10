package smb

import (
	"context"
	"testing"
	"time"
)

func TestListDirectory_EmptyContainerID(t *testing.T) {
	ctx := context.Background()
	_, _, err := ListDirectory(ctx, "", "", 10)
	if err == nil {
		t.Fatalf("expected error for empty container id")
	}
}

func TestUploadFileToContainer_InvalidDest(t *testing.T) {
	ctx := context.Background()
	err := UploadFileToContainer(ctx, "container-1", "/abs/path.txt", "/tmp/somefile")
	if err == nil {
		t.Fatalf("expected error for absolute destination path")
	}

	err = UploadFileToContainer(ctx, "container-1", "", "/tmp/somefile")
	if err == nil {
		t.Fatalf("expected error for empty destination path")
	}
}

func TestUploadFileToContainer_EmptyContainerID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := UploadFileToContainer(ctx, "", "somefile.txt", "/tmp/somefile")
	if err == nil {
		t.Fatalf("expected error for empty container id")
	}
}
