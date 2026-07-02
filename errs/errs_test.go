package errs

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestUnwrapTraversesCause(t *testing.T) {
	sentinel := errors.New("root cause")
	err := NewInternalError(SubtypeUnknown, "wrap: %v", sentinel).WithCause(sentinel)

	if !errors.Is(err, sentinel) {
		t.Fatalf("errors.Is should find the wrapped sentinel through Unwrap()")
	}
	if got := errors.Unwrap(err); !errors.Is(got, sentinel) {
		t.Fatalf("Unwrap() = %v, want %v", got, sentinel)
	}
}

func TestUnwrapAllTypesCarryCause(t *testing.T) {
	sentinel := errors.New("cause")
	cases := map[string]error{
		"validation": NewValidationError(SubtypeInvalidArgument, "x").WithCause(sentinel),
		"network":    NewNetworkError(SubtypeNetworkTransport, "x").WithCause(sentinel),
		"internal":   NewInternalError(SubtypeUnknown, "x").WithCause(sentinel),
		"api":        NewAPIError(500, "x").WithCause(sentinel),
		"auth":       NewAuthError(SubtypeAuthRequired, "x").WithCause(sentinel),
	}
	for name, err := range cases {
		if !errors.Is(err, sentinel) {
			t.Errorf("%s: errors.Is could not reach cause", name)
		}
	}
}

func TestFormatMessageNoArgsPreservesPercent(t *testing.T) {
	err := NewValidationError(SubtypeInvalidArgument, "disk 100% full")
	if strings.Contains(err.Error(), "NOVERB") {
		t.Fatalf("literal message with %% was mangled: %q", err.Error())
	}
	if err.Error() != "disk 100% full" {
		t.Fatalf("message = %q, want %q", err.Error(), "disk 100% full")
	}
}

func TestFormatMessageWithArgs(t *testing.T) {
	err := NewValidationError(SubtypeInvalidArgument, "bad %s=%d", "n", 3)
	if err.Error() != "bad n=3" {
		t.Fatalf("message = %q, want %q", err.Error(), "bad n=3")
	}
}

func TestRenderEnvelopeTyped(t *testing.T) {
	err := NewValidationError(SubtypeInvalidArgument, "bad").WithParam("--tenant").WithHint("try again")
	var env struct {
		OK    bool `json:"ok"`
		Error struct {
			Type    string `json:"type"`
			Subtype string `json:"subtype"`
			Message string `json:"message"`
			Hint    string `json:"hint"`
			Param   string `json:"param"`
		} `json:"error"`
	}
	if e := json.Unmarshal(RenderEnvelope(err), &env); e != nil {
		t.Fatalf("unmarshal: %v", e)
	}
	if env.OK {
		t.Error("ok should be false")
	}
	if env.Error.Type != string(CategoryValidation) {
		t.Errorf("type = %q", env.Error.Type)
	}
	if env.Error.Subtype != string(SubtypeInvalidArgument) {
		t.Errorf("subtype = %q", env.Error.Subtype)
	}
	if env.Error.Hint != "try again" {
		t.Errorf("hint = %q", env.Error.Hint)
	}
	if env.Error.Param != "--tenant" {
		t.Errorf("param = %q, want --tenant", env.Error.Param)
	}
}

func TestRenderEnvelopeUntyped(t *testing.T) {
	env := RenderEnvelope(errors.New("plain"))
	if !strings.Contains(string(env), string(SubtypeUnknown)) {
		t.Errorf("untyped error should map to unknown subtype: %s", env)
	}
	if !strings.Contains(string(env), string(CategoryInternal)) {
		t.Errorf("untyped error should map to internal category: %s", env)
	}
}

func TestExitCodeForCategory(t *testing.T) {
	cases := map[Category]int{
		CategoryValidation:   2,
		CategoryUnauthorized: 3,
		CategoryNetwork:      4,
		CategoryAPI:          5,
		CategoryInternal:     1,
	}
	for c, want := range cases {
		if got := ExitCodeForCategory(c); got != want {
			t.Errorf("ExitCodeForCategory(%q) = %d, want %d", c, got, want)
		}
	}
}
