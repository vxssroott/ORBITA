package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	keys := []string{
		"ORBITA_HTTP_ADDR",
		"ORBITA_READ_HEADER_TIMEOUT",
		"ORBITA_READ_TIMEOUT",
		"ORBITA_WRITE_TIMEOUT",
		"ORBITA_IDLE_TIMEOUT",
		"ORBITA_SHUTDOWN_TIMEOUT",
		"ORBITA_MAX_REQUEST_BODY_BYTES",
		"ORBITA_LOG_LEVEL",
		"ORBITA_ENVIRONMENT",
	}

	for _, key := range keys {
		t.Setenv(key, "")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("unexpected HTTP address: %s", cfg.HTTPAddr)
	}

	if cfg.ReadTimeout != 15*time.Second {
		t.Fatalf("unexpected read timeout: %s", cfg.ReadTimeout)
	}

	if cfg.MaxRequestBodyBytes != 1<<20 {
		t.Fatalf("unexpected request body limit: %d", cfg.MaxRequestBodyBytes)
	}
}

func TestLoadRejectsInvalidBodyLimit(t *testing.T) {
	t.Setenv("ORBITA_MAX_REQUEST_BODY_BYTES", "0")

	if _, err := Load(); err == nil {
		t.Fatal("expected invalid body limit to fail")
	}
}
