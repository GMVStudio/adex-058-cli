// Package daterange resolves report stat-date ranges from either a relative
// range expression (like "7d", "4w", "1m") or explicit begin/end dates.
package daterange

import (
	"strconv"
	"strings"
	"time"

	"github.com/gmvstudio/adex-cli/errs"
)

const dateLayout = "2006-01-02"

// now is overridable in tests.
var now = time.Now

// Resolve returns (begin, end) as YYYY-MM-DD strings.
//
// Precedence:
//   - If rangeStr is set, it is expanded relative to today and takes priority
//     over begin/end.
//   - Otherwise explicit begin/end are validated and returned as-is.
//
// All inputs are optional; when nothing is provided, empty strings are returned
// so the caller can decide whether the range is required.
func Resolve(rangeStr, begin, end string) (string, string, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	begin = strings.TrimSpace(begin)
	end = strings.TrimSpace(end)

	if rangeStr != "" {
		return expandRange(rangeStr)
	}

	if begin != "" {
		if err := validateDate(begin, "--begin"); err != nil {
			return "", "", err
		}
	}
	if end != "" {
		if err := validateDate(end, "--end"); err != nil {
			return "", "", err
		}
	}
	if begin != "" && end != "" && begin > end {
		return "", "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"--begin %q must not be after --end %q", begin, end).WithParam("--begin")
	}
	return begin, end, nil
}

// expandRange converts a relative expression like "7d" / "4w" / "1m" into an
// inclusive [begin, end] window ending today.
func expandRange(r string) (string, string, error) {
	unit := r[len(r)-1]
	numStr := r[:len(r)-1]
	num, err := strconv.Atoi(numStr)
	if err != nil || num <= 0 {
		return "", "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"invalid range %q: expected format like 7d, 4w, 1m", r).WithParam("--range")
	}

	today := now()
	var begin time.Time
	switch unit {
	case 'd':
		begin = today.AddDate(0, 0, -(num - 1))
	case 'w':
		begin = today.AddDate(0, 0, -(num*7 - 1))
	case 'm':
		begin = today.AddDate(0, -num, 0).AddDate(0, 0, 1)
	default:
		return "", "", errs.NewValidationError(errs.SubtypeInvalidArgument,
			"unsupported range unit %q in %q: use d (days), w (weeks), or m (months)", string(unit), r).
			WithParam("--range")
	}

	return begin.Format(dateLayout), today.Format(dateLayout), nil
}

func validateDate(s, param string) error {
	if _, err := time.Parse(dateLayout, s); err != nil {
		return errs.NewValidationError(errs.SubtypeInvalidArgument,
			"invalid date %q: expected YYYY-MM-DD", s).WithParam(param).WithCause(err)
	}
	return nil
}
