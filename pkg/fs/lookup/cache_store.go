// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	iofs "io/fs"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

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

type dirSnapshot struct {
	entry     *dirCacheEntry
	meta      FileMeta
	vcsProbed bool
	isVCSRoot bool
	vcsPath   string
	markers   []vcsMarkerState
}

type dirStore struct {
	cacheTTL      time.Duration
	mu            sync.RWMutex
	dirs          map[string]*dirCacheEntry
	dirHits       atomic.Uint64
	dirMisses     atomic.Uint64
	fileHits      atomic.Uint64
	fileMisses    atomic.Uint64
	invalidations atomic.Uint64
}

func newDirStore(cacheTTL time.Duration) *dirStore {
	return &dirStore{
		cacheTTL: cacheTTL,
		dirs:     make(map[string]*dirCacheEntry),
	}
}

func (s *dirStore) stats() CacheStats {
	return CacheStats{
		DirHits:       s.dirHits.Load(),
		DirMisses:     s.dirMisses.Load(),
		FileHits:      s.fileHits.Load(),
		FileMisses:    s.fileMisses.Load(),
		Invalidations: s.invalidations.Load(),
	}
}

func (s *dirStore) reset() {
	s.mu.Lock()
	clear(s.dirs)
	s.mu.Unlock()
	s.invalidations.Add(1)
}

func (s *dirStore) invalidateTargets(targets []string) {
	var removed uint64
	s.mu.Lock()
	for dir := range s.dirs {
		if matchesAnyInvalidationTarget(dir, targets) {
			delete(s.dirs, dir)
			removed++
		}
	}
	s.mu.Unlock()

	if removed > 0 {
		s.invalidations.Add(removed)
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

func (s *dirStore) lookupDir(
	dir string,
	now time.Time,
) (dirSnapshot, bool, bool) {
	s.mu.RLock()
	entry, ok := s.dirs[dir]
	if !ok {
		s.mu.RUnlock()
		return dirSnapshot{}, false, false
	}
	snap := dirSnapshot{
		entry:     entry,
		meta:      entry.meta,
		vcsProbed: entry.vcsProbed,
		isVCSRoot: entry.isVCSRoot,
		vcsPath:   entry.vcsPath,
		markers:   slices.Clone(entry.markers),
	}
	if entry.vcsProbed && s.cacheTTL > 0 &&
		now.Sub(entry.validatedAt) < s.cacheTTL {
		s.mu.RUnlock()
		s.dirHits.Add(1)
		return snap, true, true
	}
	s.mu.RUnlock()
	return snap, false, true
}

func (s *dirStore) touchDir(
	dir string,
	expected *dirCacheEntry,
	now time.Time,
) {
	s.mu.Lock()
	if cur, exists := s.dirs[dir]; exists && cur == expected {
		cur.validatedAt = now
	}
	s.mu.Unlock()
	s.dirHits.Add(1)
}

func (s *dirStore) updateDirVCS(
	dir string,
	expected *dirCacheEntry,
	isVCS bool,
	vcsPath string,
	markers []vcsMarkerState,
	now time.Time,
) {
	s.mu.Lock()
	if cur, exists := s.dirs[dir]; exists && cur == expected {
		cur.vcsProbed = true
		cur.isVCSRoot = isVCS
		cur.vcsPath = vcsPath
		cur.markers = markers
		cur.validatedAt = now
	}
	s.mu.Unlock()
	s.dirHits.Add(1)
}

func (s *dirStore) evictDir(dir string, expected *dirCacheEntry) {
	s.mu.Lock()
	if cur, exists := s.dirs[dir]; exists && cur == expected {
		delete(s.dirs, dir)
		s.invalidations.Add(1)
	}
	s.mu.Unlock()
}

func (s *dirStore) saveNewDir(
	dir string,
	meta FileMeta,
	isVCS bool,
	vcsPath string,
	markers []vcsMarkerState,
	now time.Time,
	hadPrevious bool,
) {
	if hadPrevious {
		s.invalidations.Add(1)
	}
	s.dirMisses.Add(1)
	newEntry := &dirCacheEntry{
		meta:        meta,
		vcsProbed:   true,
		isVCSRoot:   isVCS,
		vcsPath:     vcsPath,
		markers:     markers,
		validatedAt: now,
		files:       make(map[string]fileCacheEntry),
	}
	s.mu.Lock()
	s.dirs[dir] = newEntry
	s.mu.Unlock()
}

func (s *dirStore) ensureDirEntry(
	dir string,
	curMeta FileMeta,
	now time.Time,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cur, exists := s.dirs[dir]; exists {
		if cur.meta.sameFingerprint(curMeta) {
			cur.validatedAt = now
			return
		}
		s.invalidations.Add(1)
	}
	s.dirMisses.Add(1)
	s.dirs[dir] = &dirCacheEntry{
		meta:        curMeta,
		validatedAt: now,
		files:       make(map[string]fileCacheEntry),
	}
}

func (s *dirStore) lookupDirEntries(
	cacheKey string,
	now time.Time,
) ([]iofs.DirEntry, bool, bool) {
	s.mu.RLock()
	entry, ok := s.dirs[cacheKey]
	if ok && entry.dirEntriesCached && s.cacheTTL > 0 &&
		now.Sub(entry.dirEntriesAt) < s.cacheTTL {
		entries := slices.Clone(entry.dirEntries)
		s.mu.RUnlock()
		s.dirHits.Add(1)
		return entries, true, true
	}
	s.mu.RUnlock()
	return nil, false, ok
}

func (s *dirStore) saveDirEntries(
	cacheKey string,
	entries []iofs.DirEntry,
	now time.Time,
) {
	s.mu.Lock()
	if cur, exists := s.dirs[cacheKey]; exists {
		cur.dirEntries = slices.Clone(entries)
		cur.dirEntriesCached = true
		cur.dirEntriesAt = now
	}
	s.mu.Unlock()
}

func (s *dirStore) lookupFile(
	dir, name string,
) (*dirCacheEntry, fileCacheEntry, bool) {
	s.mu.RLock()
	dirEntry, dirOK := s.dirs[dir]
	var fentry fileCacheEntry
	var fOK bool
	if dirOK {
		fentry, fOK = dirEntry.files[name]
	}
	s.mu.RUnlock()
	return dirEntry, fentry, fOK
}

func (s *dirStore) storeFile(
	dir string,
	expectedDir *dirCacheEntry,
	name string,
	entry fileCacheEntry,
) {
	s.mu.Lock()
	cur, ok := s.dirs[dir]
	if ok && (expectedDir == nil || cur == expectedDir) {
		if cur.files == nil {
			cur.files = make(map[string]fileCacheEntry)
		}
		cur.files[name] = entry
	}
	s.mu.Unlock()
}

func (s *dirStore) recordFileHit() {
	s.fileHits.Add(1)
}

func (s *dirStore) recordFileMiss(invalidated bool) {
	if invalidated {
		s.invalidations.Add(1)
	}
	s.fileMisses.Add(1)
}
