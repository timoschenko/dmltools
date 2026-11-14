// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"errors"
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// Cache manages stat-fingerprint-validated filesystem metadata, directory
// boundary resolution, and candidate file security checks.
type Cache interface {
	ResolveStart(startPath string) (
		startDir string,
		absStart string,
		isFile bool,
		startMeta FileMeta,
		err error,
	)
	ResolveDirChain(
		startDir string,
		startMeta FileMeta,
	) ([]string, BoundaryInfo, FileMeta, error)
	ResolveBoundary(startPath string) (BoundaryInfo, error)
	DirEntries(dir string) ([]iofs.DirEntry, error)
	ProbeCandidate(
		dir string,
		name string,
		candidatePath string,
		evalBoundaryRoot string,
		startMeta FileMeta,
	) (DiscoveredFile, bool, error)
	EvalBoundaryRoot(boundaryRoot string) string
	IsWithinBoundary(boundaryRoot, targetPath string) bool
	Stat(path string) (FileMeta, error)
	Lstat(path string) (FileMeta, error)
	Abs(path string) (string, error)
	HandleSkip(path string, reason SkipReason, detail error) error
	Invalidate(path string)
	Reset()
	Stats() CacheStats
}

type securityCache struct {
	fsys            FileSystem
	clock           Clock
	vcsMarkers      []string
	stopAtVCSRoot   bool
	stopAtMount     bool
	verifyOwnership bool
	allowSymlinks   bool
	allowSymlinkEsc bool
	ceilingDirs     map[string]struct{}
	maxDepth        int
	maxFileSize     int64
	strictSecurity  bool
	onSkip          SkipCallback

	store *dirStore
}

// NewCache creates a standalone [Cache] instance configured with the given options.
func NewCache(opts ...Option) Cache {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return newCacheFromConfig(&cfg.cacheCfg)
}

func newCacheFromConfig(cfg *cacheConfig) *securityCache {
	ceilings := make(map[string]struct{}, len(cfg.ceilingDirs)*2)
	for _, dir := range cfg.ceilingDirs {
		if dir == "" {
			continue
		}
		absDir, absErr := cfg.fsys.Abs(dir)
		if absErr != nil {
			absDir = filepath.Clean(dir)
		} else {
			absDir = filepath.Clean(absDir)
		}
		ceilings[absDir] = struct{}{}
		if evalDir, evalErr := cfg.fsys.EvalSymlinks(absDir); evalErr == nil {
			ceilings[filepath.Clean(evalDir)] = struct{}{}
		}
	}

	return &securityCache{
		fsys:            cfg.fsys,
		clock:           cfg.clock,
		vcsMarkers:      cfg.vcsMarkers,
		stopAtVCSRoot:   cfg.stopAtVCSRoot,
		stopAtMount:     cfg.stopAtMountPoint,
		verifyOwnership: cfg.verifyOwnership,
		allowSymlinks:   cfg.allowSymlinks,
		allowSymlinkEsc: cfg.allowSymlinkEsc,
		ceilingDirs:     ceilings,
		maxDepth:        cfg.maxDepth,
		maxFileSize:     cfg.maxFileSize,
		strictSecurity:  cfg.strictSecurity,
		onSkip:          cfg.onSkip,
		store:           newDirStore(cfg.cacheTTL),
	}
}

func (c *securityCache) Stat(path string) (FileMeta, error) {
	return c.fsys.Stat(path)
}

func (c *securityCache) Lstat(path string) (FileMeta, error) {
	return c.fsys.Lstat(path)
}

func (c *securityCache) Abs(path string) (string, error) {
	return c.fsys.Abs(path)
}

func (c *securityCache) Stats() CacheStats {
	return c.store.stats()
}

func (c *securityCache) Reset() {
	c.store.reset()
}

func (c *securityCache) Invalidate(path string) {
	if path == "" {
		return
	}
	absPath, err := c.fsys.Abs(path)
	if err != nil {
		absPath = filepath.Clean(path)
	} else {
		absPath = filepath.Clean(absPath)
	}
	targets := []string{absPath}
	if evalPath, evalErr := c.fsys.EvalSymlinks(absPath); evalErr == nil {
		cleanEval := filepath.Clean(evalPath)
		if cleanEval != absPath {
			targets = append(targets, cleanEval)
		}
	}
	c.store.invalidateTargets(targets)
}

func (c *securityCache) DirEntries(dir string) ([]iofs.DirEntry, error) {
	now := c.clock.Now()
	cacheKey := filepath.Clean(dir)
	if absDir, err := c.fsys.Abs(dir); err == nil {
		cacheKey = filepath.Clean(absDir)
	}

	if entries, ttlHit, ok := c.store.lookupDirEntries(
		cacheKey,
		now,
	); ok && ttlHit {
		return entries, nil
	} else if !ok {
		c.ensureDirEntry(cacheKey)
	}

	dirEntries, err := c.fsys.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	entries := slices.Clone(dirEntries)
	slices.SortFunc(entries, func(a, b iofs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	c.store.saveDirEntries(cacheKey, entries, now)
	return entries, nil
}

func (c *securityCache) ensureDirEntry(dir string) {
	now := c.clock.Now()
	curMeta, err := c.fsys.Stat(dir)
	if err != nil {
		return
	}
	c.store.ensureDirEntry(dir, curMeta, now)
}

func (c *securityCache) ProbeCandidate(
	dir string,
	name string,
	candidatePath string,
	evalBoundaryRoot string,
	startMeta FileMeta,
) (DiscoveredFile, bool, error) {
	now := c.clock.Now()

	dirEntry, fentry, fOK := c.store.lookupFile(dir, name)

	if fOK && c.store.cacheTTL > 0 &&
		now.Sub(fentry.validatedAt) < c.store.cacheTTL &&
		fentry.boundaryRoot == evalBoundaryRoot {
		c.store.recordFileHit()
		return c.resultFromCacheEntry(candidatePath, fentry)
	}

	lmeta, err := c.fsys.Lstat(candidatePath)
	if err != nil {
		if os.IsNotExist(err) || errors.Is(err, iofs.ErrNotExist) {
			if fOK && !fentry.exists {
				c.store.recordFileHit()
				c.store.storeFile(dir, dirEntry, name, fileCacheEntry{
					validatedAt:  now,
					exists:       false,
					boundaryRoot: evalBoundaryRoot,
				})
				return DiscoveredFile{}, false, nil
			}
			c.store.recordFileMiss(fOK && fentry.exists)
			c.store.storeFile(dir, dirEntry, name, fileCacheEntry{
				validatedAt:  now,
				exists:       false,
				boundaryRoot: evalBoundaryRoot,
			})
			return DiscoveredFile{}, false, nil
		}

		c.store.recordFileMiss(fOK && fentry.exists)
		return DiscoveredFile{}, false, c.HandleSkip(
			candidatePath,
			SkipPermissionDenied,
			err,
		)
	}

	if fOK && fentry.exists && fentry.boundaryRoot == evalBoundaryRoot &&
		fentry.lstatMeta.sameFingerprint(lmeta) {
		if !fentry.isSymlink {
			c.store.recordFileHit()
			fentry.validatedAt = now
			c.store.storeFile(dir, dirEntry, name, fentry)
			return c.resultFromCacheEntry(candidatePath, fentry)
		}
		if resolved, rerr := c.fsys.EvalSymlinks(candidatePath); rerr == nil {
			resolved = filepath.Clean(resolved)
			if resolved == fentry.resolvedPath {
				if smeta, serr := c.fsys.Stat(resolved); serr == nil &&
					fentry.targetMeta.sameFingerprint(smeta) {
					c.store.recordFileHit()
					fentry.validatedAt = now
					c.store.storeFile(dir, dirEntry, name, fentry)
					return c.resultFromCacheEntry(candidatePath, fentry)
				}
			}
		}
	}

	c.store.recordFileMiss(fOK)

	newEntry := c.evaluateCandidate(
		dir,
		name,
		candidatePath,
		lmeta,
		evalBoundaryRoot,
		startMeta,
		now,
	)
	c.store.storeFile(dir, dirEntry, name, newEntry)
	return c.resultFromCacheEntry(candidatePath, newEntry)
}

func (c *securityCache) resultFromCacheEntry(
	candidatePath string,
	entry fileCacheEntry,
) (DiscoveredFile, bool, error) {
	if !entry.exists {
		return DiscoveredFile{}, false, nil
	}
	if entry.hasSkip {
		return DiscoveredFile{}, false, c.HandleSkip(
			candidatePath,
			entry.skipReason,
			entry.skipErr,
		)
	}
	return entry.discovered, true, nil
}

func (c *securityCache) evaluateCandidate(
	dir string,
	name string,
	candidatePath string,
	lmeta FileMeta,
	evalBoundaryRoot string,
	startMeta FileMeta,
	now time.Time,
) fileCacheEntry {
	base := fileCacheEntry{
		validatedAt:  now,
		exists:       true,
		lstatMeta:    lmeta,
		boundaryRoot: evalBoundaryRoot,
	}

	isSymlink := lmeta.Mode&iofs.ModeSymlink != 0
	base.isSymlink = isSymlink

	targetMeta := lmeta
	resolvedPath := candidatePath

	if isSymlink {
		if !c.allowSymlinks {
			base.hasSkip = true
			base.skipReason = SkipSymlinkDisabled
			return base
		}
		resolved, err := c.fsys.EvalSymlinks(candidatePath)
		if err != nil {
			base.hasSkip = true
			base.skipReason = SkipBrokenSymlink
			base.skipErr = err
			return base
		}
		resolvedPath = filepath.Clean(resolved)
		base.resolvedPath = resolvedPath

		smeta, err := c.fsys.Stat(resolvedPath)
		if err != nil {
			base.hasSkip = true
			base.skipReason = SkipBrokenSymlink
			base.skipErr = err
			return base
		}
		targetMeta = smeta
		base.targetMeta = smeta

		if c.stopAtMount && startMeta.HasDevice && smeta.HasDevice &&
			smeta.DeviceID != startMeta.DeviceID {
			base.hasSkip = true
			base.skipReason = SkipSymlinkCrossDevice
			base.skipErr = fmt.Errorf(
				"symlink target %s crosses device boundary (%d != %d)",
				resolvedPath,
				smeta.DeviceID,
				startMeta.DeviceID,
			)
			return base
		}

		if !c.allowSymlinkEsc &&
			!isWithinBoundary(evalBoundaryRoot, resolvedPath) {
			base.hasSkip = true
			base.skipReason = SkipSymlinkEscape
			base.skipErr = fmt.Errorf(
				"symlink target %s escapes boundary %s",
				resolvedPath,
				evalBoundaryRoot,
			)
			return base
		}
	}

	if !targetMeta.Mode.IsRegular() {
		base.hasSkip = true
		base.skipReason = SkipNotRegular
		base.skipErr = fmt.Errorf("mode %v is not regular", targetMeta.Mode)
		return base
	}

	if c.maxFileSize > 0 && targetMeta.Size > c.maxFileSize {
		base.hasSkip = true
		base.skipReason = SkipFileTooLarge
		base.skipErr = fmt.Errorf(
			"size %d exceeds limit %d",
			targetMeta.Size,
			c.maxFileSize,
		)
		return base
	}

	base.discovered = DiscoveredFile{
		Path:         candidatePath,
		ResolvedPath: resolvedPath,
		Dir:          dir,
		Name:         name,
		Size:         targetMeta.Size,
		ModTime:      targetMeta.ModTime,
		Mode:         targetMeta.Mode,
		IsSymlink:    isSymlink,
		Meta:         targetMeta,
	}
	return base
}

func (c *securityCache) HandleSkip(
	path string,
	reason SkipReason,
	detail error,
) error {
	if c.onSkip != nil {
		c.onSkip(path, reason, detail)
	}
	if c.strictSecurity && reason != SkipExcluded {
		return &SecurityError{
			Path:   path,
			Reason: reason,
			Err:    detail,
		}
	}
	return nil
}
