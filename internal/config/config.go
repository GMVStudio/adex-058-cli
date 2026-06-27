package config

import (
	"os"
	"strconv"
)

const (
	DefaultBaseURL   = "http://47.99.131.55:8000"
	DefaultPage      = 1
	DefaultPageSize  = 20
	DefaultOrderBy   = "charge"
	DefaultOrderDesc = true
	DefaultStatHour  = -1
)

// Config holds the CLI configuration resolved from env + flags.
type Config struct {
	BaseURL string
}

// Load resolves configuration from environment variables with sensible defaults.
func Load() *Config {
	cfg := &Config{
		BaseURL: DefaultBaseURL,
	}
	if v := os.Getenv("ADEX_API_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	return cfg
}

// EnvInt reads an integer from an env var, falling back to def.
func EnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
