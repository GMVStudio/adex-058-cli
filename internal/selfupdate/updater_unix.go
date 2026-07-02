//go:build !windows

package selfupdate

func (u *Updater) PrepareSelfReplace() (restore func(), err error) {
	return func() {}, nil
}

func (u *Updater) CleanupStaleFiles() {}

func (u *Updater) CanRestorePreviousVersion() bool {
	if u.RestoreAvailableOverride != nil {
		return u.RestoreAvailableOverride()
	}
	return u.backupCreated
}
