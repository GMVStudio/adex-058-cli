package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/vfs"
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
	// TenantID is the default tenant used when --tenant is not passed.
	TenantID int `json:"tenant_id,omitempty"`
}

// Dir returns the config directory, honoring ADEX_CONFIG_DIR for tests,
// falling back to ~/.adex.
func Dir() string {
	if v := os.Getenv("ADEX_CONFIG_DIR"); v != "" {
		return v
	}
	home, err := vfs.Default.UserHomeDir()
	if err != nil || home == "" {
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

	if data, err := vfs.Default.ReadFile(Path()); err == nil {
		var fromFile Config
		if json.Unmarshal(data, &fromFile) == nil {
			if fromFile.BaseURL != "" {
				cfg.BaseURL = fromFile.BaseURL
			}
			if fromFile.Authorization != "" {
				cfg.Authorization = fromFile.Authorization
			}
			if fromFile.TenantID > 0 {
				cfg.TenantID = fromFile.TenantID
			}
		}
	}

	if v := os.Getenv("ADEX_API_BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("ADEX_AUTHORIZATION"); v != "" {
		cfg.Authorization = NormalizeToken(v)
	}
	if v := os.Getenv("ADEX_TENANT_ID"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.TenantID = n
		}
	}
	return cfg
}

// Save persists the config to the config file (0600), creating the
// directory if needed. Failures are returned as typed *errs.InternalError.
func Save(cfg *Config) error {
	dir := Dir()
	if err := vfs.Default.MkdirAll(dir, 0o700); err != nil {
		return errs.NewInternalError(errs.SubtypeFileIO,
			"failed to create config dir %q: %v", dir, err).
			WithCause(err).
			WithHint("set ADEX_CONFIG_DIR to a writable directory (e.g. /tmp/adex), " +
				"or skip adex init and use env vars: ADEX_AUTHORIZATION and ADEX_API_BASE_URL")
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return errs.NewInternalError(errs.SubtypeUnknown,
			"failed to marshal config: %v", err).WithCause(err)
	}
	if err := vfs.Default.WriteFile(Path(), data, 0o600); err != nil {
		return errs.NewInternalError(errs.SubtypeFileIO,
			"failed to write config %q: %v", Path(), err).
			WithCause(err).
			WithHint("set ADEX_CONFIG_DIR to a writable directory (e.g. /tmp/adex), " +
				"or skip adex init and use env vars: ADEX_AUTHORIZATION and ADEX_API_BASE_URL")
	}
	return nil
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
