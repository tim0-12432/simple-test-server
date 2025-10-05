package common

import (
	"os"
	"testing"

	"github.com/tim0-12432/simple-test-server/config"
)

func TestGetMaxUploadSizeDefault(t *testing.T) {
	// ensure config.EnvConfig is nil to simulate default behavior
	prev := config.EnvConfig
	config.EnvConfig = nil
	defer func() { config.EnvConfig = prev }()

	got := GetMaxUploadSize()
	if got != 10<<20 {
		t.Fatalf("expected default 10MB, got %d", got)
	}
}

func TestGetMaxUploadSizeFromEnv(t *testing.T) {
	// set via env var and reinitialize config
	os.Setenv("UPLOAD_MAX_BYTES", "2097152") // 2MB
	config.InitializeEnvConfig()
	defer os.Unsetenv("UPLOAD_MAX_BYTES")

	if config.EnvConfig == nil {
		t.Fatal("expected EnvConfig to be initialized")
	}

	if config.EnvConfig.UploadMaxBytes != 2097152 {
		t.Fatalf("expected config.UploadMaxBytes to be 2097152, got %d", config.EnvConfig.UploadMaxBytes)
	}

	got := GetMaxUploadSize()
	if got != 2097152 {
		t.Fatalf("expected GetMaxUploadSize to return 2097152, got %d", got)
	}
}
