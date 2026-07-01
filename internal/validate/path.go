// Package validate provides path-safety checks for untrusted CLI inputs.
// CLI arguments frequently originate from AI agents, so any file path must be
// validated before it reaches file I/O.
package validate

import (
	"path/filepath"
	"strings"

	"github.com/gmvstudio/adex-cli/errs"
)

// SafeInputPath validates a path the CLI will read from. It rejects empty
// paths and paths containing a parent-directory (`..`) traversal segment, and
// returns the cleaned, absolute path.
func SafeInputPath(path string) (string, error) {
	return safePath(path, "input")
}

// SafeOutputPath validates a path the CLI will write to, applying the same
// traversal checks as SafeInputPath.
func SafeOutputPath(path string) (string, error) {
	return safePath(path, "output")
}

func safePath(path, kind string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"%s path must not be empty", kind).WithParam("path")
	}
	cleaned := filepath.Clean(path)
	if hasTraversal(cleaned) {
		return "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"%s path %q must not contain a '..' traversal segment", kind, path).
			WithParam("path")
	}
	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"cannot resolve %s path %q: %v", kind, path, err).
			WithParam("path").WithCause(err)
	}
	return abs, nil
}

// hasTraversal reports whether a cleaned path contains a ".." segment.
func hasTraversal(cleaned string) bool {
	for _, seg := range strings.Split(cleaned, string(filepath.Separator)) {
		if seg == ".." {
			return true
		}
	}
	return false
}
