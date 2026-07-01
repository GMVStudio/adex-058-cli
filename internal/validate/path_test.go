package validate

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/gmvstudio/adex-cli/errs"
)

func TestSafeInputPathRejectsEmpty(t *testing.T) {
	_, err := SafeInputPath("   ")
	var ve *errs.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("error type = %T, want *errs.ValidationError", err)
	}
	if ve.Param != "path" {
		t.Errorf("param = %q, want path", ve.Param)
	}
}

func TestSafeInputPathRejectsTraversal(t *testing.T) {
	for _, p := range []string{"../etc/passwd", "a/../../b", "foo/../../bar"} {
		if _, err := SafeInputPath(p); err == nil {
			t.Errorf("SafeInputPath(%q) should reject traversal", p)
		}
	}
}

func TestSafeInputPathAcceptsCleanPath(t *testing.T) {
	got, err := SafeInputPath("data/report.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !filepath.IsAbs(got) {
		t.Errorf("result %q should be absolute", got)
	}
}

func TestSafeOutputPathSameRules(t *testing.T) {
	if _, err := SafeOutputPath("../x"); err == nil {
		t.Error("SafeOutputPath should reject traversal")
	}
	if _, err := SafeOutputPath("out.csv"); err != nil {
		t.Errorf("SafeOutputPath clean path: %v", err)
	}
}
