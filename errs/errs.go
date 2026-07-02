package errs

import (
	"encoding/json"
	"errors"
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
	SubtypeInvalidResponse  Subtype = "invalid_response"
)

// formatMessage applies fmt.Sprintf only when args are present, so a caller
// passing a literal message containing a stray "%" (e.g. "disk 100% full") is
// not rendered as "%!(NOVERB)".
func formatMessage(format string, args []interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

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

// Unwrap exposes the wrapped cause so errors.Is / errors.As traverse the chain.
func (p *Problem) Unwrap() error {
	return p.Cause
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
			Message:  formatMessage(format, args),
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
			Message:  formatMessage(format, args),
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
			Message:  formatMessage(format, args),
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
			Message:  formatMessage(format, args),
		},
	}
}

func (e *APIError) WithHint(hint string) *APIError {
	e.Hint = hint
	return e
}

func (e *APIError) WithCause(cause error) *APIError {
	e.Cause = cause
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
			Message:  formatMessage(format, args),
		},
	}
}

func (e *AuthError) WithHint(hint string) *AuthError {
	e.Hint = hint
	return e
}

func (e *AuthError) WithCause(cause error) *AuthError {
	e.Cause = cause
	return e
}

// IsTyped returns true if err is (or wraps) one of the typed *errs.* errors.
func IsTyped(err error) bool {
	_, ok := ProblemOf(err)
	return ok
}

// ProblemOf extracts the Problem from a typed error anywhere in the chain,
// using errors.As so it works through wrapped errors.
func ProblemOf(err error) (*Problem, bool) {
	var ve *ValidationError
	if errors.As(err, &ve) {
		return &ve.Problem, true
	}
	var ne *NetworkError
	if errors.As(err, &ne) {
		return &ne.Problem, true
	}
	var ie *InternalError
	if errors.As(err, &ie) {
		return &ie.Problem, true
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return &apiErr.Problem, true
	}
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return &authErr.Problem, true
	}
	return nil, false
}

// Envelope is the JSON error envelope written to stderr.
type Envelope struct {
	OK    bool        `json:"ok"`
	Error interface{} `json:"error"`
}

// RenderEnvelope serializes a typed error into the JSON stderr envelope.
// It preserves type-specific fields (e.g. ValidationError.Param) by marshaling
// the concrete typed error when possible, instead of just the embedded Problem.
func RenderEnvelope(err error) []byte {
	// ValidationError has a Param field not present on Problem; marshal the
	// full struct so param reaches the JSON output.
	var ve *ValidationError
	if errors.As(err, &ve) {
		b, _ := json.Marshal(Envelope{OK: false, Error: ve})
		return b
	}

	// Other typed errors (NetworkError, InternalError, APIError, AuthError) have
	// no fields beyond Problem, so ProblemOf is sufficient.
	prob, ok := ProblemOf(err)
	if !ok {
		prob = &Problem{
			Category: CategoryInternal,
			Subtype:  SubtypeUnknown,
			Message:  err.Error(),
		}
	}
	b, _ := json.Marshal(Envelope{OK: false, Error: prob})
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
