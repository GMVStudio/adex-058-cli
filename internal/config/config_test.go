package config

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/vfs"
)

// failFS makes WriteFile fail so we can assert Save returns a typed error.
type failFS struct{ vfs.OS }

func (failFS) WriteFile(string, []byte, fs.FileMode) error { return errors.New("disk full") }

func TestSaveReturnsTypedErrorOnWriteFailure(t *testing.T) {
	t.Setenv("ADEX_CONFIG_DIR", t.TempDir())
	orig := vfs.Default
	vfs.Default = failFS{}
	defer func() { vfs.Default = orig }()

	err := Save(&Config{BaseURL: "http://x"})
	var ie *errs.InternalError
	if !errors.As(err, &ie) {
		t.Fatalf("error type = %T, want *errs.InternalError", err)
	}
	if ie.Subtype != errs.SubtypeFileIO {
		t.Errorf("subtype = %q, want file_io", ie.Subtype)
	}
}

func TestNormalizeToken(t *testing.T) {
	cases := map[string]string{
		"Bearer adex_abc": "adex_abc",
		"bearer adex_abc": "adex_abc",
		"  adex_abc  ":    "adex_abc",
		"adex_abc":        "adex_abc",
		"Bearer   x":      "x",
		"":                "",
	}
	for in, want := range cases {
		if got := NormalizeToken(in); got != want {
			t.Errorf("NormalizeToken(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ADEX_CONFIG_DIR", dir)

	cfg := &Config{BaseURL: "http://example.test", Authorization: "adex_secret"}
	if err := Save(cfg); err != nil {
		t.Fatalf("save: %v", err)
	}

	loaded := Load()
	if loaded.BaseURL != "http://example.test" {
		t.Errorf("BaseURL = %q", loaded.BaseURL)
	}
	if loaded.Authorization != "adex_secret" {
		t.Errorf("Authorization = %q", loaded.Authorization)
	}
}

func TestEnvOverridesFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ADEX_CONFIG_DIR", dir)
	if err := Save(&Config{BaseURL: "http://file", Authorization: "file_key"}); err != nil {
		t.Fatalf("save: %v", err)
	}

	t.Setenv("ADEX_API_BASE_URL", "http://env")
	t.Setenv("ADEX_AUTHORIZATION", "Bearer env_key")

	loaded := Load()
	if loaded.BaseURL != "http://env" {
		t.Errorf("BaseURL = %q, want http://env", loaded.BaseURL)
	}
	if loaded.Authorization != "env_key" {
		t.Errorf("Authorization = %q, want env_key (normalized)", loaded.Authorization)
	}
}

func TestLoadDefaultsWhenNoFile(t *testing.T) {
	t.Setenv("ADEX_CONFIG_DIR", t.TempDir())
	loaded := Load()
	if loaded.BaseURL != DefaultBaseURL {
		t.Errorf("BaseURL = %q, want default %q", loaded.BaseURL, DefaultBaseURL)
	}
	if loaded.Authorization != "" {
		t.Errorf("Authorization = %q, want empty", loaded.Authorization)
	}
}
