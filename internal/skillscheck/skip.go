package skillscheck

import (
	"os"
	"regexp"
)

var semverPattern = regexp.MustCompile(`^v?\d+\.\d+\.\d+`)

func shouldSkip(version string) bool {
	if os.Getenv("ADEX_NO_SKILLS_NOTIFIER") != "" {
		return true
	}
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		return true
	}
	if version == "DEV" || version == "dev" || version == "" {
		return true
	}
	return !semverPattern.MatchString(version)
}
