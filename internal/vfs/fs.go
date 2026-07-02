// Package vfs abstracts filesystem operations so code paths that read or write
// files can be exercised in tests without touching the real disk. Production
// code uses vfs.Default; tests may swap in a fake implementation.
package vfs

import "io/fs"

// FS is the filesystem surface used across the CLI. Implementations must behave
// identically to the corresponding os package functions.
type FS interface {
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm fs.FileMode) error
	MkdirAll(path string, perm fs.FileMode) error
	Stat(name string) (fs.FileInfo, error)
	Remove(name string) error
	Rename(oldpath, newpath string) error
	UserHomeDir() (string, error)
	Executable() (string, error)
	EvalSymlinks(path string) (string, error)
}

// Default is the process-wide filesystem. Reassign in tests if needed.
var Default FS = OS{}
