package skillscheck

import "strings"

func Init(currentVersion string) {
	SetPending(nil)
	if shouldSkip(currentVersion) {
		return
	}
	version, ok := ReadSyncedVersion()
	if !ok {
		return
	}
	if strings.TrimPrefix(strings.TrimPrefix(version, "v"), "V") == strings.TrimPrefix(strings.TrimPrefix(currentVersion, "v"), "V") {
		return
	}
	SetPending(&StaleNotice{
		Current: version,
		Target:  currentVersion,
	})
}
