package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr            string
	ReadHeaderTimeout   time.Duration
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
	ShutdownTimeout     time.Duration
	MaxRequestBodyBytes int64
	LogLevel            string
	Environment         string
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:            envString("ORBITA_HTTP_ADDR", ":8080"),
		ReadHeaderTimeout:   envDuration("ORBITA_READ_HEADER_TIMEOUT", 5*time.Second),
		ReadTimeout:         envDuration("ORBITA_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:        envDuration("ORBITA_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:         envDuration("ORBITA_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:     envDuration("ORBITA_SHUTDOWN_TIMEOUT", 10*time.Second),
		MaxRequestBodyBytes: envInt64("ORBITA_MAX_REQUEST_BODY_BYTES", 1<<20),
		LogLevel:            envString("ORBITA_LOG_LEVEL", "info"),
		Environment:         envString("ORBITA_ENVIRONMENT", "development"),
	}

	if cfg.MaxRequestBodyBytes <= 0 {
		return Config{}, fmt.Errorf("ORBITA_MAX_REQUEST_BODY_BYTES must be positive")
	}

	if cfg.ReadHeaderTimeout <= 0 ||
		cfg.ReadTimeout <= 0 ||
		cfg.WriteTimeout <= 0 ||
		cfg.IdleTimeout <= 0 ||
		cfg.ShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP timeout configuration must be positive")
	}

	return cfg, nil
}

func envString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func envInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}
