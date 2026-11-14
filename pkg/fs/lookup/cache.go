// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"errors"
	"fmt"
	iofs "io/fs"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
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
	Stat(path string) (FileMeta, error)
	Lstat(path string) (FileMeta, error)
	Abs(path string) (string, error)
	HandleSkip(path string, reason SkipReason, detail error) error
	Invalidate(path string)
	Reset()
	Stats() CacheStats
}

type vcsMarkerState struct {
	marker string
	exists bool
	meta   FileMeta
}

type fileCacheEntry struct {
	validatedAt  time.Time
	exists       bool
	lstatMeta    FileMeta
	isSymlink    bool
	resolvedPath string
	targetMeta   FileMeta
	boundaryRoot string
	discovered   DiscoveredFile
	skipReason   SkipReason
	hasSkip      bool
	skipErr      error
}

type dirCacheEntry struct {
	meta             FileMeta
	vcsProbed        bool
	isVCSRoot        bool
	vcsPath          string
	markers          []vcsMarkerState
	validatedAt      time.Time
	dirEntries       []iofs.DirEntry
	dirEntriesCached bool
	dirEntriesAt     time.Time
	files            map[string]fileCacheEntry
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
	cacheTTL        time.Duration
	strictSecurity  bool
	onSkip          SkipCallback

	mu            sync.RWMutex
	dirs          map[string]*dirCacheEntry
	dirHits       atomic.Uint64
	dirMisses     atomic.Uint64
	fileHits      atomic.Uint64
	fileMisses    atomic.Uint64
	invalidations atomic.Uint64
}

// NewCache creates a standalone [Cache] instance configured with the given options.
func NewCache(opts ...Option) Cache {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return newCacheFromConfig(&cfg)
}

func newCacheFromConfig(cfg *config) *securityCache {
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
		cacheTTL:        cfg.cacheTTL,
		strictSecurity:  cfg.strictSecurity,
		onSkip:          cfg.onSkip,
		dirs:            make(map[string]*dirCacheEntry),
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
	return CacheStats{
		DirHits:       c.dirHits.Load(),
		DirMisses:     c.dirMisses.Load(),
		FileHits:      c.fileHits.Load(),
		FileMisses:    c.fileMisses.Load(),
		Invalidations: c.invalidations.Load(),
	}
}

func (c *securityCache) Reset() {
	c.mu.Lock()
	clear(c.dirs)
	c.mu.Unlock()
	c.invalidations.Add(1)
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

	var removed uint64
	c.mu.Lock()
	for dir := range c.dirs {
		if matchesAnyInvalidationTarget(dir, targets) {
			delete(c.dirs, dir)
			removed++
		}
	}
	c.mu.Unlock()

	if removed > 0 {
		c.invalidations.Add(removed)
	}
}

func matchesAnyInvalidationTarget(dir string, targets []string) bool {
	sep := string(filepath.Separator)
	for _, target := range targets {
		parentDir := filepath.Dir(target)
		prefix := target + sep
		if dir == target || dir == parentDir ||
			strings.HasPrefix(dir, prefix) {
			return true
		}
	}
	return false
}

func (c *securityCache) ResolveBoundary(
	startPath string,
) (BoundaryInfo, error) {
	startDir, _, _, startMeta, err := c.ResolveStart(startPath)
	if err != nil {
		return BoundaryInfo{}, err
	}
	_, boundary, _, err := c.ResolveDirChain(startDir, startMeta)
	if err != nil {
		return BoundaryInfo{}, err
	}
	return boundary, nil
}

func (c *securityCache) ResolveStart(
	startPath string,
) (string, string, bool, FileMeta, error) {
	if startPath == "" {
		startPath = "."
	}
	absPath, err := c.fsys.Abs(startPath)
	if err != nil {
		return "", "", false, FileMeta{}, err
	}
	absPath = filepath.Clean(absPath)

	meta, err := c.fsys.Stat(absPath)
	if err == nil {
		if meta.Mode.IsDir() {
			return absPath, absPath, false, meta, nil
		}
		startDir := filepath.Dir(absPath)
		dirMeta, dirErr := c.fsys.Stat(startDir)
		if dirErr != nil {
			return "", "", false, FileMeta{}, dirErr
		}
		return startDir, absPath, true, dirMeta, nil
	}

	lmeta, lerr := c.fsys.Lstat(absPath)
	if lerr != nil {
		return "", "", false, FileMeta{}, err
	}
	if lmeta.Mode.IsDir() {
		return absPath, absPath, false, lmeta, nil
	}
	startDir := filepath.Dir(absPath)
	dirMeta, dirErr := c.fsys.Stat(startDir)
	if dirErr != nil {
		return "", "", false, FileMeta{}, dirErr
	}
	return startDir, absPath, true, dirMeta, nil
}

func (c *securityCache) ResolveDirChain(
	startDir string,
	startMeta FileMeta,
) ([]string, BoundaryInfo, FileMeta, error) {
	dirs := make([]string, 0, 8)
	currentDir := startDir
	depth := 0
	euid := c.fsys.Geteuid()

	curMeta, isVCS, vcsPath, err := c.getOrRefreshDir(currentDir)
	if err != nil {
		return nil, BoundaryInfo{}, FileMeta{}, err
	}
	startMeta = curMeta

	for {
		dirs = append(dirs, currentDir)

		if c.isCeiling(currentDir) {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryCeiling,
			}, startMeta, nil
		}
		if c.stopAtVCSRoot && isVCS {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryVCSRoot,
				VCSPath: vcsPath,
			}, startMeta, nil
		}
		if depth >= c.maxDepth {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryMaxDepth,
			}, startMeta, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryFilesystemRoot,
			}, startMeta, nil
		}

		parentMeta, parentIsVCS, parentVCSPath, parentErr := c.getOrRefreshDir(
			parentDir,
		)
		if parentErr != nil {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryOwnership,
			}, startMeta, nil
		}

		if c.stopAtMount && startMeta.HasDevice && parentMeta.HasDevice &&
			parentMeta.DeviceID != startMeta.DeviceID {
			return dirs, BoundaryInfo{
				RootDir: currentDir,
				Reason:  BoundaryMountPoint,
			}, startMeta, nil
		}

		if c.verifyOwnership && startMeta.HasUID && parentMeta.HasUID {
			if !isTrustedOwner(parentMeta.UID, startMeta.UID, euid) {
				return dirs, BoundaryInfo{
					RootDir: currentDir,
					Reason:  BoundaryOwnership,
				}, startMeta, nil
			}
		}

		currentDir = parentDir
		isVCS = parentIsVCS
		vcsPath = parentVCSPath
		depth++
	}
}

func isTrustedOwner(parentUID, startUID uint32, euid int) bool {
	if parentUID == startUID || parentUID == 0 {
		return true
	}
	if euid >= 0 && int64(euid) <= math.MaxUint32 {
		// #nosec G115 -- guarded against negative and > MaxUint32
		return parentUID == uint32(euid)
	}
	return false
}

func (c *securityCache) isCeiling(dir string) bool {
	if len(c.ceilingDirs) == 0 {
		return false
	}
	if _, ok := c.ceilingDirs[dir]; ok {
		return true
	}
	if evalDir, err := c.fsys.EvalSymlinks(dir); err == nil {
		_, ok := c.ceilingDirs[filepath.Clean(evalDir)]
		return ok
	}
	return false
}

func (c *securityCache) getOrRefreshDir(
	dir string,
) (FileMeta, bool, string, error) {
	now := c.clock.Now()

	c.mu.RLock()
	entry, ok := c.dirs[dir]
	if ok && entry.vcsProbed && c.cacheTTL > 0 &&
		now.Sub(entry.validatedAt) < c.cacheTTL {
		meta := entry.meta
		isVCS := entry.isVCSRoot
		vcsPath := entry.vcsPath
		c.mu.RUnlock()
		c.dirHits.Add(1)
		return meta, isVCS, vcsPath, nil
	}
	var cachedMeta FileMeta
	var cachedVCSProbed bool
	var cachedMarkers []vcsMarkerState
	var cachedVCS bool
	var cachedVCSPath string
	if ok {
		cachedMeta = entry.meta
		cachedVCSProbed = entry.vcsProbed
		cachedMarkers = slices.Clone(entry.markers)
		cachedVCS = entry.isVCSRoot
		cachedVCSPath = entry.vcsPath
	}
	c.mu.RUnlock()

	curMeta, err := c.fsys.Stat(dir)
	if err != nil {
		if ok {
			c.mu.Lock()
			if cur, exists := c.dirs[dir]; exists && cur == entry {
				delete(c.dirs, dir)
				c.invalidations.Add(1)
			}
			c.mu.Unlock()
		}
		return FileMeta{}, false, "", err
	}

	if ok && cachedMeta.sameFingerprint(curMeta) {
		if cachedVCSProbed && c.markersUnchanged(dir, cachedMarkers) {
			c.mu.Lock()
			if cur, exists := c.dirs[dir]; exists && cur == entry {
				cur.validatedAt = now
			}
			c.mu.Unlock()
			c.dirHits.Add(1)
			return cachedMeta, cachedVCS, cachedVCSPath, nil
		}
		if !cachedVCSProbed {
			isVCS, vcsPath, markers := c.probeVCSMarkers(dir)
			c.mu.Lock()
			if cur, exists := c.dirs[dir]; exists && cur == entry {
				cur.vcsProbed = true
				cur.isVCSRoot = isVCS
				cur.vcsPath = vcsPath
				cur.markers = markers
				cur.validatedAt = now
			}
			c.mu.Unlock()
			c.dirHits.Add(1)
			return cachedMeta, isVCS, vcsPath, nil
		}
	}

	if ok {
		c.invalidations.Add(1)
	}
	c.dirMisses.Add(1)

	isVCS, vcsPath, markers := c.probeVCSMarkers(dir)
	newEntry := &dirCacheEntry{
		meta:        curMeta,
		vcsProbed:   true,
		isVCSRoot:   isVCS,
		vcsPath:     vcsPath,
		markers:     markers,
		validatedAt: now,
		files:       make(map[string]fileCacheEntry),
	}

	c.mu.Lock()
	c.dirs[dir] = newEntry
	c.mu.Unlock()

	return curMeta, isVCS, vcsPath, nil
}

func (c *securityCache) ensureDirEntry(dir string) {
	now := c.clock.Now()
	curMeta, err := c.fsys.Stat(dir)
	if err != nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if cur, exists := c.dirs[dir]; exists {
		if cur.meta.sameFingerprint(curMeta) {
			cur.validatedAt = now
			return
		}
		c.invalidations.Add(1)
	}
	c.dirMisses.Add(1)
	c.dirs[dir] = &dirCacheEntry{
		meta:        curMeta,
		validatedAt: now,
		files:       make(map[string]fileCacheEntry),
	}
}

func (c *securityCache) probeVCSMarkers(
	dir string,
) (bool, string, []vcsMarkerState) {
	if !c.stopAtVCSRoot || len(c.vcsMarkers) == 0 {
		return false, "", nil
	}
	states := make([]vcsMarkerState, 0, len(c.vcsMarkers))
	isVCS := false
	vcsPath := ""

	for _, marker := range c.vcsMarkers {
		p := filepath.Join(dir, marker)
		lmeta, err := c.fsys.Lstat(p)
		if err != nil {
			states = append(states, vcsMarkerState{
				marker: marker,
				exists: false,
			})
			continue
		}
		states = append(states, vcsMarkerState{
			marker: marker,
			exists: true,
			meta:   lmeta,
		})
		if isVCS {
			continue
		}
		if lmeta.Mode.IsDir() || lmeta.Mode.IsRegular() {
			isVCS = true
			vcsPath = p
			continue
		}
		if lmeta.Mode&iofs.ModeSymlink != 0 {
			if smeta, serr := c.fsys.Stat(p); serr == nil &&
				(smeta.Mode.IsDir() || smeta.Mode.IsRegular()) {
				isVCS = true
				vcsPath = p
			}
		}
	}

	return isVCS, vcsPath, states
}

func (c *securityCache) markersUnchanged(
	dir string,
	states []vcsMarkerState,
) bool {
	if !c.stopAtVCSRoot {
		return len(states) == 0
	}
	if len(states) != len(c.vcsMarkers) {
		return false
	}
	for _, st := range states {
		p := filepath.Join(dir, st.marker)
		lmeta, err := c.fsys.Lstat(p)
		if err != nil {
			if st.exists {
				return false
			}
			continue
		}
		if !st.exists || !st.meta.sameFingerprint(lmeta) {
			return false
		}
	}
	return true
}

func (c *securityCache) DirEntries(dir string) ([]iofs.DirEntry, error) {
	now := c.clock.Now()
	cacheKey := filepath.Clean(dir)
	if absDir, err := c.fsys.Abs(dir); err == nil {
		cacheKey = filepath.Clean(absDir)
	}

	c.mu.RLock()
	entry, ok := c.dirs[cacheKey]
	if ok && entry.dirEntriesCached && c.cacheTTL > 0 &&
		now.Sub(entry.dirEntriesAt) < c.cacheTTL {
		entries := slices.Clone(entry.dirEntries)
		c.mu.RUnlock()
		c.dirHits.Add(1)
		return entries, nil
	}
	c.mu.RUnlock()

	if !ok {
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

	c.mu.Lock()
	if cur, exists := c.dirs[cacheKey]; exists {
		cur.dirEntries = slices.Clone(entries)
		cur.dirEntriesCached = true
		cur.dirEntriesAt = now
	}
	c.mu.Unlock()

	return entries, nil
}

func (c *securityCache) ProbeCandidate(
	dir string,
	name string,
	candidatePath string,
	evalBoundaryRoot string,
	startMeta FileMeta,
) (DiscoveredFile, bool, error) {
	now := c.clock.Now()

	c.mu.RLock()
	dirEntry, dirOK := c.dirs[dir]
	var fentry fileCacheEntry
	var fOK bool
	if dirOK {
		fentry, fOK = dirEntry.files[name]
	}
	c.mu.RUnlock()

	if fOK && c.cacheTTL > 0 && now.Sub(fentry.validatedAt) < c.cacheTTL &&
		fentry.boundaryRoot == evalBoundaryRoot {
		c.fileHits.Add(1)
		return c.resultFromCacheEntry(candidatePath, fentry)
	}

	lmeta, err := c.fsys.Lstat(candidatePath)
	if err != nil {
		if os.IsNotExist(err) || errors.Is(err, iofs.ErrNotExist) {
			if fOK && !fentry.exists {
				c.fileHits.Add(1)
				c.storeFileCache(dir, dirEntry, name, fileCacheEntry{
					validatedAt:  now,
					exists:       false,
					boundaryRoot: evalBoundaryRoot,
				})
				return DiscoveredFile{}, false, nil
			}
			if fOK && fentry.exists {
				c.invalidations.Add(1)
			}
			c.fileMisses.Add(1)
			c.storeFileCache(dir, dirEntry, name, fileCacheEntry{
				validatedAt:  now,
				exists:       false,
				boundaryRoot: evalBoundaryRoot,
			})
			return DiscoveredFile{}, false, nil
		}

		if fOK && fentry.exists {
			c.invalidations.Add(1)
		}
		c.fileMisses.Add(1)
		return DiscoveredFile{}, false, c.HandleSkip(
			candidatePath,
			SkipPermissionDenied,
			err,
		)
	}

	if fOK && fentry.exists && fentry.boundaryRoot == evalBoundaryRoot &&
		fentry.lstatMeta.sameFingerprint(lmeta) {
		if !fentry.isSymlink {
			c.fileHits.Add(1)
			fentry.validatedAt = now
			c.storeFileCache(dir, dirEntry, name, fentry)
			return c.resultFromCacheEntry(candidatePath, fentry)
		}
		if resolved, rerr := c.fsys.EvalSymlinks(candidatePath); rerr == nil {
			resolved = filepath.Clean(resolved)
			if resolved == fentry.resolvedPath {
				if smeta, serr := c.fsys.Stat(resolved); serr == nil &&
					fentry.targetMeta.sameFingerprint(smeta) {
					c.fileHits.Add(1)
					fentry.validatedAt = now
					c.storeFileCache(dir, dirEntry, name, fentry)
					return c.resultFromCacheEntry(candidatePath, fentry)
				}
			}
		}
	}

	if fOK {
		c.invalidations.Add(1)
	}
	c.fileMisses.Add(1)

	newEntry := c.evaluateCandidate(
		dir,
		name,
		candidatePath,
		lmeta,
		evalBoundaryRoot,
		startMeta,
		now,
	)
	c.storeFileCache(dir, dirEntry, name, newEntry)
	return c.resultFromCacheEntry(candidatePath, newEntry)
}

func (c *securityCache) storeFileCache(
	dir string,
	expectedDir *dirCacheEntry,
	name string,
	entry fileCacheEntry,
) {
	c.mu.Lock()
	cur, ok := c.dirs[dir]
	if ok && (expectedDir == nil || cur == expectedDir) {
		if cur.files == nil {
			cur.files = make(map[string]fileCacheEntry)
		}
		cur.files[name] = entry
	}
	c.mu.Unlock()
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

func (c *securityCache) EvalBoundaryRoot(boundaryRoot string) string {
	if evalRoot, err := c.fsys.EvalSymlinks(boundaryRoot); err == nil {
		return filepath.Clean(evalRoot)
	}
	return filepath.Clean(boundaryRoot)
}

func isWithinBoundary(boundaryRoot, target string) bool {
	cleanRoot := filepath.Clean(boundaryRoot)
	cleanTarget := filepath.Clean(target)
	rel, err := filepath.Rel(cleanRoot, cleanTarget)
	if err != nil {
		return false
	}
	if filepath.IsAbs(rel) || rel == ".." {
		return false
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator))
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
