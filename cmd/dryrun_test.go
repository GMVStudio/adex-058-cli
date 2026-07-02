package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/gmvstudio/adex-cli/errs"
	"github.com/gmvstudio/adex-cli/internal/config"
)

// TestDryRunPaths asserts that each command builds the expected request path
// and core params, without touching the network.
func TestDryRunPaths(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantPath string
	}{
		{"ks accounts", []string{"ks", "accounts", "--tenant", "6", "--dry-run"}, "/v1/ks/ad-accounts"},
		{"ks campaigns", []string{"ks", "campaigns", "--tenant", "6", "--dry-run"}, "/v1/ks/campaigns"},
		{"ks campaign get", []string{"ks", "campaigns", "get", "999", "--tenant", "6", "--dry-run"}, "/v1/ks/campaigns/999"},
		{"ks campaigns top", []string{"ks", "campaigns", "top", "--tenant", "6", "--range", "7d", "--dry-run"}, "/v1/ks/campaigns/top"},
		{"ks campaign-reports daily", []string{"ks", "campaign-reports", "daily", "--tenant", "6", "--range", "7d", "--dry-run"}, "/v1/ks/campaign-reports/daily"},
		{"ks campaign-reports summary", []string{"ks", "campaign-reports", "summary", "--tenant", "6", "--range", "7d", "--dry-run"}, "/v1/ks/campaign-reports/summary"},
		{"ks dashboard", []string{"ks", "dashboard", "--tenant", "6", "--range", "7d", "--dry-run"}, "/v1/ks/dashboard"},
		{"ks metric-meta", []string{"ks", "report-metric-meta", "--level", "account", "--dry-run"}, "/v1/ks/report-metric-meta"},
		{"oe projects", []string{"oe", "projects", "--tenant", "6", "--dry-run"}, "/v1/oe/projects"},
		{"oe budget-vs-actual", []string{"oe", "account-budget-vs-actual", "--tenant", "6", "--range", "7d", "--dry-run"}, "/v1/oe/account-budget-vs-actual"},
		{"tenant", []string{"tenant", "--dry-run"}, "/v1/tenants"},
		{"user", []string{"user", "--dry-run"}, "/v1/users/me"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := runCmd(t, tc.args...)
			if res.ExecErr != nil {
				t.Fatalf("unexpected error: %v (stderr: %s)", res.ExecErr, res.Err)
			}
			if res.DryRunPath != tc.wantPath {
				t.Errorf("path = %q, want %q", res.DryRunPath, tc.wantPath)
			}
		})
	}
}

// TestDryRunParams asserts filter flags are translated into the right params.
func TestDryRunParams(t *testing.T) {
	res := runCmd(t, "ks", "campaigns",
		"--tenant", "6",
		"--campaign", "abc",
		"--campaign-name", "promo",
		"--put-status", "1",
		"--dry-run")
	if res.ExecErr != nil {
		t.Fatalf("unexpected error: %v", res.ExecErr)
	}
	if v, _ := res.paramString("tenant_id"); v != "6" {
		t.Errorf("tenant_id = %q, want 6", v)
	}
	if v, _ := res.paramString("campaign_id"); v != "abc" {
		t.Errorf("campaign_id = %q, want abc", v)
	}
	if v, _ := res.paramString("campaign_name"); v != "promo" {
		t.Errorf("campaign_name = %q, want promo", v)
	}
	if v, _ := res.paramString("put_status"); v != "1" {
		t.Errorf("put_status = %q, want 1", v)
	}
}

// TestDryRunDateRangeResolved asserts --range expands to stat_date bounds.
func TestDryRunDateRangeResolved(t *testing.T) {
	res := runCmd(t, "ks", "campaign-reports", "summary",
		"--tenant", "6", "--range", "7d", "--dry-run")
	if res.ExecErr != nil {
		t.Fatalf("unexpected error: %v", res.ExecErr)
	}
	if _, ok := res.paramString("stat_date_begin"); !ok {
		t.Error("stat_date_begin missing after --range 7d")
	}
	if _, ok := res.paramString("stat_date_end"); !ok {
		t.Error("stat_date_end missing after --range 7d")
	}
}

// TestNonPositiveTenantIsValidationError asserts the requireTenant guard
// rejects a non-positive --tenant with a typed error naming the param.
func TestNonPositiveTenantIsValidationError(t *testing.T) {
	res := runCmd(t, "ks", "accounts", "--tenant", "0", "--dry-run")
	if res.ExecErr == nil {
		t.Fatal("expected an error when --tenant is 0")
	}
	var ve *errs.ValidationError
	if !errors.As(res.ExecErr, &ve) {
		t.Fatalf("error type = %T, want *errs.ValidationError", res.ExecErr)
	}
	if ve.Param != "--tenant" {
		t.Errorf("param = %q, want --tenant", ve.Param)
	}
}

// TestMissingRequiredFlagWrapped asserts cobra's required-flag error is
// surfaced (and, via handleError, rendered as a validation envelope).
func TestMissingRequiredFlagWrapped(t *testing.T) {
	res := runCmd(t, "ks", "accounts", "--dry-run")
	if res.ExecErr == nil {
		t.Fatal("expected a required-flag error when --tenant is missing")
	}
	// The raw cobra error is untyped; handleError promotes it to validation.
	var out, errOut bytes.Buffer
	code := handleError(newTestFactory(&out, &errOut), res.ExecErr)
	if code != errs.ExitCodeForCategory(errs.CategoryValidation) {
		t.Errorf("exit code = %d, want validation exit code", code)
	}
}

// TestRequiredDateRangeError asserts summary/top require a date range.
func TestRequiredDateRangeError(t *testing.T) {
	res := runCmd(t, "ks", "campaign-reports", "summary", "--tenant", "6", "--dry-run")
	if res.ExecErr == nil {
		t.Fatal("expected an error when no date range is provided")
	}
	var ve *errs.ValidationError
	if !errors.As(res.ExecErr, &ve) {
		t.Fatalf("error type = %T, want *errs.ValidationError", res.ExecErr)
	}
}

// TestInvalidDateRejected asserts a malformed --begin is a validation error.
func TestInvalidDateRejected(t *testing.T) {
	res := runCmd(t, "ks", "campaign-reports", "daily",
		"--tenant", "6", "--begin", "2026/07/01", "--dry-run")
	if res.ExecErr == nil {
		t.Fatal("expected an error for malformed --begin")
	}
	var ve *errs.ValidationError
	if !errors.As(res.ExecErr, &ve) {
		t.Fatalf("error type = %T, want *errs.ValidationError", res.ExecErr)
	}
	if ve.Param != "--begin" {
		t.Errorf("param = %q, want --begin", ve.Param)
	}
}

// TestInvalidJqRejected asserts a bad --jq expression fails before the network.
func TestInvalidJqRejected(t *testing.T) {
	res := runCmd(t, "tenant", "--jq", "this is not valid", "--dry-run")
	if res.ExecErr == nil {
		t.Fatal("expected an error for invalid --jq")
	}
	var ve *errs.ValidationError
	if !errors.As(res.ExecErr, &ve) {
		t.Fatalf("error type = %T, want *errs.ValidationError", res.ExecErr)
	}
}

// TestMissingConfigError asserts that running a command without --dry-run and
// no API key returns a missing_config validation error, not a network error.
// It also verifies the hint mentions the env var alternative for sandbox agents.
func TestMissingConfigError(t *testing.T) {
	var out, errOut bytes.Buffer
	f := &Factory{
		Config: &config.Config{BaseURL: "http://test.local", Authorization: ""},
		Out:    &out,
		ErrOut: &errOut,
	}
	root := NewRootCmd(f)
	root.SetOut(&out)
	root.SetErr(&errOut)
	root.SetArgs([]string{"tenant"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected missing_config error when Authorization is empty")
	}
	var ve *errs.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("error type = %T, want *errs.ValidationError", err)
	}
	if ve.Subtype != errs.SubtypeMissingConfig {
		t.Errorf("subtype = %q, want %q", ve.Subtype, errs.SubtypeMissingConfig)
	}
	if !strings.Contains(ve.Hint, "ADEX_AUTHORIZATION") {
		t.Errorf("hint = %q, expected to mention ADEX_AUTHORIZATION for sandbox recovery", ve.Hint)
	}
}
