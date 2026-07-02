package vfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

// OS is the production FS backed by the standard os package.
type OS struct{}

func (OS) ReadFile(name string) ([]byte, error) { return os.ReadFile(name) }

func (OS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (OS) MkdirAll(path string, perm fs.FileMode) error { return os.MkdirAll(path, perm) }

func (OS) Stat(name string) (fs.FileInfo, error) { return os.Stat(name) }

func (OS) Remove(name string) error { return os.Remove(name) }

func (OS) Rename(oldpath, newpath string) error { return os.Rename(oldpath, newpath) }

func (OS) UserHomeDir() (string, error) { return os.UserHomeDir() }

func (OS) Executable() (string, error) { return os.Executable() }

func (OS) EvalSymlinks(path string) (string, error) { return filepath.EvalSymlinks(path) }
