// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	iofs "io/fs"
	"os"
	"path/filepath"
	"time"
)

// Clock supplies the current time for cache TTL evaluation.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

// FileMeta captures portable and platform-specific filesystem metadata used for
// security boundary checks and cache fingerprint validation.
type FileMeta struct {
	Name      string
	Size      int64
	Mode      iofs.FileMode
	ModTime   time.Time
	DeviceID  uint64
	Inode     uint64
	UID       uint32
	HasDevice bool
	HasInode  bool
	HasUID    bool
}

func (m FileMeta) sameFingerprint(other FileMeta) bool {
	if m.Mode != other.Mode || m.Size != other.Size ||
		!m.ModTime.Equal(other.ModTime) {
		return false
	}
	if m.HasDevice != other.HasDevice ||
		(m.HasDevice && m.DeviceID != other.DeviceID) {
		return false
	}
	if m.HasInode != other.HasInode ||
		(m.HasInode && m.Inode != other.Inode) {
		return false
	}
	if m.HasUID != other.HasUID || (m.HasUID && m.UID != other.UID) {
		return false
	}
	return true
}

// FileSystem abstracts OS filesystem operations so boundary, mount-point,
// ownership, and cache behavior can be tested deterministically across OSes.
type FileSystem interface {
	Stat(path string) (FileMeta, error)
	Lstat(path string) (FileMeta, error)
	ReadDir(path string) ([]iofs.DirEntry, error)
	ReadFile(path string) ([]byte, error)
	EvalSymlinks(path string) (string, error)
	Abs(path string) (string, error)
	Geteuid() int
}

type osFileSystem struct{}

// NewOSFileSystem returns the default [FileSystem] backed by the OS.
func NewOSFileSystem() FileSystem {
	return osFileSystem{}
}

func (osFileSystem) Stat(path string) (FileMeta, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileMeta{}, err
	}
	return fileInfoToMeta(path, info), nil
}

func (osFileSystem) Lstat(path string) (FileMeta, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return FileMeta{}, err
	}
	return fileInfoToMeta(path, info), nil
}

func (osFileSystem) ReadDir(path string) ([]iofs.DirEntry, error) {
	return os.ReadDir(path)
}

func (osFileSystem) ReadFile(path string) ([]byte, error) {
	//nolint:gosec // callers validate candidate paths via Cache security guards
	return os.ReadFile(filepath.Clean(path))
}

func (osFileSystem) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

func (osFileSystem) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (osFileSystem) Geteuid() int {
	return osGeteuid()
}

func fileInfoToMeta(path string, info os.FileInfo) FileMeta {
	meta := FileMeta{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
	}
	fillSysMeta(path, info, &meta)
	return meta
}
