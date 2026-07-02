//go:build windows

package selfupdate

import (
	"fmt"

	"github.com/gmvstudio/adex-cli/internal/vfs"
)

func (u *Updater) resolveExe() (string, error) {
	exe, err := vfs.Default.Executable()
	if err != nil {
		return "", err
	}
	return vfs.Default.EvalSymlinks(exe)
}

func (u *Updater) PrepareSelfReplace() (restore func(), err error) {
	noop := func() {}

	exe, err := u.resolveExe()
	if err != nil {
		return noop, nil
	}

	oldPath := exe + ".old"

	vfs.Default.Remove(oldPath)

	if err := vfs.Default.Rename(exe, oldPath); err != nil {
		return noop, fmt.Errorf("cannot rename binary for update: %w", err)
	}
	u.backupCreated = true

	restore = func() {
		if _, err := vfs.Default.Stat(oldPath); err != nil {
			u.backupCreated = false
			return
		}
		vfs.Default.Remove(exe)
		if err := vfs.Default.Rename(oldPath, exe); err != nil {
			u.backupCreated = false
		}
	}

	return restore, nil
}

func (u *Updater) CleanupStaleFiles() {
	exe, err := u.resolveExe()
	if err != nil {
		return
	}
	oldPath := exe + ".old"

	if _, err := vfs.Default.Stat(oldPath); err != nil {
		return
	}

	if _, err := vfs.Default.Stat(exe); err != nil {
		vfs.Default.Rename(oldPath, exe)
		return
	}

	vfs.Default.Remove(oldPath)
}

func (u *Updater) CanRestorePreviousVersion() bool {
	if u.RestoreAvailableOverride != nil {
		return u.RestoreAvailableOverride()
	}
	return u.backupCreated
}
