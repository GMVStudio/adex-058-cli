package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/gmvstudio/adex-cli/errs"
)

// TestHandleErrorTypedEnvelope asserts a typed error is rendered as a JSON
// envelope on stderr with the right exit code.
func TestHandleErrorTypedEnvelope(t *testing.T) {
	var out, errOut bytes.Buffer
	f := newTestFactory(&out, &errOut)

	err := errs.NewAuthError(errs.SubtypeAuthRequired, "no token").WithHint("run adex init")
	code := handleError(f, err)

	if code != errs.ExitCodeForCategory(errs.CategoryUnauthorized) {
		t.Errorf("exit code = %d, want %d", code, errs.ExitCodeForCategory(errs.CategoryUnauthorized))
	}
	var env struct {
		OK    bool `json:"ok"`
		Error struct {
			Type    string `json:"type"`
			Subtype string `json:"subtype"`
			Hint    string `json:"hint"`
		} `json:"error"`
	}
	if e := json.Unmarshal(errOut.Bytes(), &env); e != nil {
		t.Fatalf("stderr is not a JSON envelope: %v (%q)", e, errOut.String())
	}
	if env.OK {
		t.Error("ok should be false")
	}
	if env.Error.Type != string(errs.CategoryUnauthorized) {
		t.Errorf("type = %q", env.Error.Type)
	}
	if env.Error.Hint != "run adex init" {
		t.Errorf("hint = %q", env.Error.Hint)
	}
}

// TestHandleErrorUntypedWrapped asserts a plain error is wrapped as validation.
func TestHandleErrorUntypedWrapped(t *testing.T) {
	var out, errOut bytes.Buffer
	f := newTestFactory(&out, &errOut)

	code := handleError(f, errors.New("boom"))
	if code != errs.ExitCodeForCategory(errs.CategoryValidation) {
		t.Errorf("exit code = %d, want validation", code)
	}
	if !strings.Contains(errOut.String(), string(errs.CategoryValidation)) {
		t.Errorf("stderr should carry a validation envelope: %s", errOut.String())
	}
}

// TestSkillsList asserts the embedded skills list command returns an envelope.
// It is skipped when skills are not embedded (unit-test build without embed).
func TestSkillsList(t *testing.T) {
	if embeddedSkillContent == nil {
		t.Skip("skills not embedded in this build")
	}
	res := runCmd(t, "skills", "list")
	if res.ExecErr != nil {
		t.Fatalf("unexpected error: %v", res.ExecErr)
	}
	if !strings.Contains(res.Out, `"ok":true`) {
		t.Errorf("skills list output = %q", res.Out)
	}
}

// TestUnknownCommand asserts an unknown command produces an error.
func TestUnknownCommand(t *testing.T) {
	res := runCmd(t, "definitely-not-a-command")
	if res.ExecErr == nil {
		t.Fatal("expected an error for an unknown command")
	}
}
