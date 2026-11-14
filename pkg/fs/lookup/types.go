// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"fmt"
	iofs "io/fs"
	"time"
)

const (
	// FileGitIgnore is the standard Git ignore rules filename.
	FileGitIgnore = ".gitignore"
	// FileIgnore is the generic ignore rules filename used by search/lint tools.
	FileIgnore = ".ignore"
	// FileEditorConfig is the EditorConfig formatting configuration filename.
	FileEditorConfig = ".editorconfig"

	// VCSMarkerGit is the Git repository/worktree marker name.
	VCSMarkerGit = ".git"
	// VCSMarkerHg is the Mercurial repository marker name.
	VCSMarkerHg = ".hg"

	// FilterForceExclude is the filter name returned for force-exclude matches.
	FilterForceExclude = "--force-exclude"
	// FilterExclude is the filter name returned for exclude matches.
	FilterExclude = "--exclude"
	// FilterExtendExclude is the filter name returned for extend-exclude matches.
	FilterExtendExclude = "--extend-exclude"
	// FilterInclude is the filter name used when reporting include filter skips.
	FilterInclude = "--include"

	// DefaultMaxDepth is the default maximum upward directory traversal depth.
	DefaultMaxDepth = 64
	// DefaultMaxFileSize is the default maximum allowed utility file size (1 MiB).
	DefaultMaxFileSize int64 = 1 << 20
)

// BoundaryReason identifies why upward directory traversal stopped.
type BoundaryReason int

const (
	// BoundaryNone indicates no boundary was determined.
	BoundaryNone BoundaryReason = iota
	// BoundaryFilesystemRoot indicates traversal reached the filesystem root.
	BoundaryFilesystemRoot
	// BoundaryVCSRoot indicates traversal stopped at a VCS repository root.
	BoundaryVCSRoot
	// BoundaryMountPoint indicates the parent directory resides on a different mount/device.
	BoundaryMountPoint
	// BoundaryOwnership indicates the parent directory has an untrusted owner UID.
	BoundaryOwnership
	// BoundaryCeiling indicates traversal reached a configured ceiling directory.
	BoundaryCeiling
	// BoundaryMaxDepth indicates traversal reached the configured max depth limit.
	BoundaryMaxDepth
)

// String returns a human-readable description of the boundary reason.
func (r BoundaryReason) String() string {
	switch r {
	case BoundaryNone:
		return "none"
	case BoundaryFilesystemRoot:
		return "filesystem-root"
	case BoundaryVCSRoot:
		return "vcs-root"
	case BoundaryMountPoint:
		return "mount-point"
	case BoundaryOwnership:
		return "ownership"
	case BoundaryCeiling:
		return "ceiling"
	case BoundaryMaxDepth:
		return "max-depth"
	default:
		return fmt.Sprintf("unknown(%d)", int(r))
	}
}

// BoundaryInfo describes the outermost directory permitted for a lookup.
type BoundaryInfo struct {
	RootDir string
	Reason  BoundaryReason
	VCSPath string
}

// SkipReason classifies why a candidate directory or utility file was skipped.
type SkipReason int

const (
	// SkipNotRegular indicates a candidate path is not a regular file (e.g. FIFO, socket, directory).
	SkipNotRegular SkipReason = iota + 1
	// SkipFileTooLarge indicates a candidate file exceeds the configured MaxFileSize.
	SkipFileTooLarge
	// SkipSymlinkDisabled indicates symlinks are disallowed by configuration.
	SkipSymlinkDisabled
	// SkipBrokenSymlink indicates a symlink target does not exist or cannot be resolved.
	SkipBrokenSymlink
	// SkipSymlinkEscape indicates a symlink resolves outside the active boundary root.
	SkipSymlinkEscape
	// SkipSymlinkCrossDevice indicates a symlink resolves onto a different mount/device.
	SkipSymlinkCrossDevice
	// SkipPermissionDenied indicates a filesystem permission or stat error occurred.
	SkipPermissionDenied
	// SkipExcluded indicates a path matched exclude, extend-exclude, or force-exclude filters.
	SkipExcluded
)

// String returns a human-readable label for the skip reason.
func (r SkipReason) String() string {
	switch r {
	case SkipNotRegular:
		return "not-regular-file"
	case SkipFileTooLarge:
		return "file-too-large"
	case SkipSymlinkDisabled:
		return "symlink-disabled"
	case SkipBrokenSymlink:
		return "broken-symlink"
	case SkipSymlinkEscape:
		return "symlink-escape"
	case SkipSymlinkCrossDevice:
		return "symlink-cross-device"
	case SkipPermissionDenied:
		return "permission-denied"
	case SkipExcluded:
		return "excluded"
	default:
		return fmt.Sprintf("unknown(%d)", int(r))
	}
}

// SkipCallback is invoked whenever a candidate file or directory is rejected
// by a security guard or path exclusion filter during upward discovery.
type SkipCallback func(path string, reason SkipReason, detail error)

// IgnoreCallback is invoked during recursive collection when a file or
// directory is ignored due to inclusion or exclusion filters.
type IgnoreCallback func(path string, isDir bool, filterName string)

// SecurityError is returned by discovery methods when StrictSecurity is
// enabled and a candidate file violates a security check.
type SecurityError struct {
	Path   string
	Reason SkipReason
	Err    error
}

// Error implements the error interface.
func (e *SecurityError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf(
			"lookup security check failed (%s) for %s: %v",
			e.Reason,
			e.Path,
			e.Err,
		)
	}
	return fmt.Sprintf(
		"lookup security check failed (%s) for %s",
		e.Reason,
		e.Path,
	)
}

// Unwrap returns the underlying error associated with this security failure.
func (e *SecurityError) Unwrap() error {
	return e.Err
}

// DiscoveredFile holds the resolved path and validated metadata for a file.
type DiscoveredFile struct {
	Path         string
	ResolvedPath string
	Dir          string
	Name         string
	Size         int64
	ModTime      time.Time
	Mode         iofs.FileMode
	IsSymlink    bool
	Meta         FileMeta
}

// CacheStats reports cache hit, miss, and invalidation counters.
type CacheStats struct {
	DirHits       uint64
	DirMisses     uint64
	FileHits      uint64
	FileMisses    uint64
	Invalidations uint64
}
