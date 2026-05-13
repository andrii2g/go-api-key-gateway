package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment         string
	DBDSN               string
	PepperBase64        string
	PepperFile          string
	SecretBytes         int
	AdminToken          string
	HTTPAddr            string
	UsageQueueSize      int
	UsageFlushTimeout   time.Duration
	MarkUsedTimeout     time.Duration
	TrustedProxyHeaders bool
}

func Load() (Config, error) {
	cfg := Config{
		Environment:       os.Getenv("APIKEY_ENV"),
		DBDSN:             os.Getenv("APIKEY_DB_DSN"),
		PepperBase64:      os.Getenv("APIKEY_PEPPER_BASE64"),
		PepperFile:        os.Getenv("APIKEY_PEPPER_FILE"),
		AdminToken:        os.Getenv("APIKEY_ADMIN_TOKEN"),
		HTTPAddr:          getenvDefault("APIKEY_HTTP_ADDR", ":8080"),
		SecretBytes:       getenvIntDefault("APIKEY_SECRET_BYTES", 32),
		UsageQueueSize:    getenvIntDefault("APIKEY_USAGE_QUEUE_SIZE", 1024),
		UsageFlushTimeout: time.Duration(getenvIntDefault("APIKEY_USAGE_FLUSH_TIMEOUT_MS", 250)) * time.Millisecond,
		MarkUsedTimeout:   500 * time.Millisecond,
	}

	if cfg.Environment == "" {
		return Config{}, fmt.Errorf("APIKEY_ENV is required")
	}
	if cfg.DBDSN == "" {
		return Config{}, fmt.Errorf("APIKEY_DB_DSN is required")
	}
	if cfg.PepperBase64 == "" && cfg.PepperFile == "" {
		return Config{}, fmt.Errorf("one of APIKEY_PEPPER_BASE64 or APIKEY_PEPPER_FILE is required")
	}
	if cfg.SecretBytes < 32 {
		return Config{}, fmt.Errorf("APIKEY_SECRET_BYTES must be at least 32")
	}
	if cfg.AdminToken == "" {
		return Config{}, fmt.Errorf("APIKEY_ADMIN_TOKEN is required")
	}
	return cfg, nil
}

func getenvDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getenvIntDefault(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
