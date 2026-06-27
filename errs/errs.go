package errs

import (
	"encoding/json"
	"fmt"
)

// Category is the top-level error classification.
type Category string

const (
	CategoryValidation   Category = "validation"
	CategoryNetwork      Category = "network"
	CategoryInternal     Category = "internal"
	CategoryAPI          Category = "api"
	CategoryUnauthorized Category = "unauthorized"
)

// Subtype is a stable, machine-readable error identifier.
type Subtype string

const (
	SubtypeInvalidArgument  Subtype = "invalid_argument"
	SubtypeMissingConfig    Subtype = "missing_config"
	SubtypeNetworkTransport Subtype = "network_transport"
	SubtypeFileIO           Subtype = "file_io"
	SubtypeUnknown          Subtype = "unknown"
	SubtypeAPIError         Subtype = "api_error"
	SubtypeAuthRequired     Subtype = "auth_required"
)

// Problem is the typed error embedded in every CLI error.
type Problem struct {
	Category Category `json:"type"`
	Subtype  Subtype  `json:"subtype"`
	Code     int      `json:"code,omitempty"`
	Message  string   `json:"message"`
	Hint     string   `json:"hint,omitempty"`
	Cause    error    `json:"-"`
}

func (p *Problem) Error() string {
	return p.Message
}

// ValidationError is for user-facing input failures.
type ValidationError struct {
	Problem
	Param string `json:"param,omitempty"`
}

func NewValidationError(subtype Subtype, format string, args ...interface{}) *ValidationError {
	return &ValidationError{
		Problem: Problem{
			Category: CategoryValidation,
			Subtype:  subtype,
			Message:  fmt.Sprintf(format, args...),
		},
	}
}

func (e *ValidationError) WithParam(param string) *ValidationError {
	e.Param = param
	return e
}

func (e *ValidationError) WithHint(hint string) *ValidationError {
	e.Hint = hint
	return e
}

func (e *ValidationError) WithCause(cause error) *ValidationError {
	e.Cause = cause
	return e
}

// NetworkError is for transport-level failures.
type NetworkError struct {
	Problem
}

func NewNetworkError(subtype Subtype, format string, args ...interface{}) *NetworkError {
	return &NetworkError{
		Problem: Problem{
			Category: CategoryNetwork,
			Subtype:  subtype,
			Message:  fmt.Sprintf(format, args...),
		},
	}
}

func (e *NetworkError) WithHint(hint string) *NetworkError {
	e.Hint = hint
	return e
}

func (e *NetworkError) WithCause(cause error) *NetworkError {
	e.Cause = cause
	return e
}

// InternalError is for unexpected internal failures.
type InternalError struct {
	Problem
}

func NewInternalError(subtype Subtype, format string, args ...interface{}) *InternalError {
	return &InternalError{
		Problem: Problem{
			Category: CategoryInternal,
			Subtype:  subtype,
			Message:  fmt.Sprintf(format, args...),
		},
	}
}

func (e *InternalError) WithHint(hint string) *InternalError {
	e.Hint = hint
	return e
}

func (e *InternalError) WithCause(cause error) *InternalError {
	e.Cause = cause
	return e
}

// APIError is for non-2xx API responses.
type APIError struct {
	Problem
}

func NewAPIError(code int, format string, args ...interface{}) *APIError {
	return &APIError{
		Problem: Problem{
			Category: CategoryAPI,
			Subtype:  SubtypeAPIError,
			Code:     code,
			Message:  fmt.Sprintf(format, args...),
		},
	}
}

func (e *APIError) WithHint(hint string) *APIError {
	e.Hint = hint
	return e
}

// AuthError is for authentication/authorization failures.
type AuthError struct {
	Problem
}

func NewAuthError(subtype Subtype, format string, args ...interface{}) *AuthError {
	return &AuthError{
		Problem: Problem{
			Category: CategoryUnauthorized,
			Subtype:  subtype,
			Message:  fmt.Sprintf(format, args...),
		},
	}
}

func (e *AuthError) WithHint(hint string) *AuthError {
	e.Hint = hint
	return e
}

// IsTyped returns true if err is one of the typed *errs.* errors.
func IsTyped(err error) bool {
	switch err.(type) {
	case *ValidationError, *NetworkError, *InternalError, *APIError, *AuthError:
		return true
	}
	return false
}

// ProblemOf extracts the Problem from a typed error, if any.
func ProblemOf(err error) (*Problem, bool) {
	switch e := err.(type) {
	case *ValidationError:
		return &e.Problem, true
	case *NetworkError:
		return &e.Problem, true
	case *InternalError:
		return &e.Problem, true
	case *APIError:
		return &e.Problem, true
	case *AuthError:
		return &e.Problem, true
	}
	return nil, false
}

// Envelope is the JSON error envelope written to stderr.
type Envelope struct {
	OK    bool     `json:"ok"`
	Error *Problem `json:"error"`
}

// RenderEnvelope serializes a typed error into the JSON stderr envelope.
func RenderEnvelope(err error) []byte {
	prob, ok := ProblemOf(err)
	if !ok {
		prob = &Problem{
			Category: CategoryInternal,
			Subtype:  SubtypeUnknown,
			Message:  err.Error(),
		}
	}
	env := Envelope{OK: false, Error: prob}
	b, _ := json.Marshal(env)
	return b
}

// ExitCodeForCategory maps a Category to a process exit code.
func ExitCodeForCategory(c Category) int {
	switch c {
	case CategoryValidation:
		return 2
	case CategoryUnauthorized:
		return 3
	case CategoryNetwork:
		return 4
	case CategoryAPI:
		return 5
	default:
		return 1
	}
}
