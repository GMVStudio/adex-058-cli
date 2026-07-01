package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	DefaultBaseURL   = "http://47.99.131.55:8000"
	DefaultPage      = 1
	DefaultPageSize  = 20
	DefaultOrderBy   = "charge"
	DefaultOrderDesc = true
	DefaultStatHour  = -1
)

// Config holds the CLI configuration resolved from file + env + flags.
type Config struct {
	BaseURL string `json:"base_url,omitempty"`
	// Authorization is the API key (without the "Bearer " prefix).
	Authorization string `json:"authorization,omitempty"`
}

// Dir returns the config directory, honoring ADEX_CONFIG_DIR for tests,
// falling back to ~/.adex.
func Dir() string {
	if v := os.Getenv("ADEX_CONFIG_DIR"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".adex"
	}
	return filepath.Join(home, ".adex")
}

// Path returns the absolute path to the config file.
func Path() string {
	return filepath.Join(Dir(), "config.json")
}

// Load resolves configuration from the config file, then overlays
// environment variables. Env vars take precedence over the file.
func Load() *Config {
	cfg := &Config{
		BaseURL: DefaultBaseURL,
	}

	if data, err := os.ReadFile(Path()); err == nil {
		var fromFile Config
		if json.Unmarshal(data, &fromFile) == nil {
			if fromFile.BaseURL != "" {
				cfg.BaseURL = fromFile.BaseURL
			}
			if fromFile.Authorization != "" {
				cfg.Authorization = fromFile.Authorization
			}
		}
	}

	if v := os.Getenv("ADEX_API_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("ADEX_AUTHORIZATION"); v != "" {
		cfg.Authorization = NormalizeToken(v)
	}
	return cfg
}

// Save persists the config to the config file (0600), creating the
// directory if needed.
func Save(cfg *Config) error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0o600)
}

// NormalizeToken strips an optional "Bearer " prefix and surrounding
// whitespace, returning the bare API key.
func NormalizeToken(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 7 && strings.EqualFold(s[:7], "bearer ") {
		s = strings.TrimSpace(s[7:])
	}
	return s
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
