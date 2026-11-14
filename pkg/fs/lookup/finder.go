// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"fmt"
	"path/filepath"
	"slices"
	"sync"
)

// Finder combines [Cache] and [Filter] to perform upward configuration
// discovery and downward recursive source collection.
type Finder interface {
	// FindChain walks upward from startPath to the effective boundary and
	// returns all matching files ordered Root-to-Leaf (outermost first,
	// nearest to startPath last). Within the same directory, matches follow
	// the order of names (or lexical order when matching via Include regex).
	FindChain(startPath string, names ...string) ([]DiscoveredFile, error)

	// FindNearest walks upward from startPath to the effective boundary and
	// returns the closest matching file to startPath.
	FindNearest(
		startPath string,
		names ...string,
	) (DiscoveredFile, bool, error)

	// Collect discovers, filters, deduplicates, and lexically sorts file paths
	// from the provided source files and directories.
	Collect(sources []string) ([]string, error)

	// MatchIgnore evaluates configured ignore files (.gitignore, .ignore) for
	// path within its resolved boundary, returning the matching ignore file
	// name or an empty string if not ignored.
	MatchIgnore(path string) (string, error)

	// SetOnIgnore registers the callback invoked when Collect skips a path
	// due to inclusion or exclusion filters.
	SetOnIgnore(cb IgnoreCallback)

	// Cache returns the underlying [Cache] instance.
	Cache() Cache

	// Filter returns the underlying [Filter] instance.
	Filter() Filter
}

type fileFinder struct {
	cache         Cache
	filter        Filter
	ignoreMatcher IgnoreMatcher
	ignoreFiles   []string
	mu            sync.RWMutex
	onIgnore      IgnoreCallback
}

// NewFinder constructs a [Finder] using injected [WithCache], [WithFilter], or
// [WithIgnoreMatcher] instances, or creating defaults from the provided options.
func NewFinder(opts ...Option) (Finder, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	f := cfg.filter
	if f == nil {
		rf, err := newFilterFromConfig(&cfg.filterCfg)
		if err != nil {
			return nil, err
		}
		f = rf
	}

	c := cfg.cache
	if c == nil {
		c = newCacheFromConfig(&cfg.cacheCfg)
	}

	im := cfg.ignoreMatcher
	if im == nil && len(cfg.ignoreFiles) > 0 {
		im = newIgnoreMatcher(cfg.cacheCfg.fsys)
	}

	return &fileFinder{
		cache:         c,
		filter:        f,
		ignoreMatcher: im,
		ignoreFiles:   slices.Clone(cfg.ignoreFiles),
		onIgnore:      cfg.onIgnore,
	}, nil
}

func (f *fileFinder) Cache() Cache {
	return f.cache
}

func (f *fileFinder) Filter() Filter {
	return f.filter
}

func (f *fileFinder) SetOnIgnore(cb IgnoreCallback) {
	f.mu.Lock()
	f.onIgnore = cb
	f.mu.Unlock()
}

func (f *fileFinder) notifyIgnore(path string, isDir bool, filterName string) {
	f.mu.RLock()
	cb := f.onIgnore
	f.mu.RUnlock()
	if cb != nil {
		cb(path, isDir, filterName)
	}
}

func (f *fileFinder) FindChain(
	startPath string,
	names ...string,
) ([]DiscoveredFile, error) {
	startDir, absStart, isFile, startMeta, err := f.cache.ResolveStart(
		startPath,
	)
	if err != nil {
		return nil, err
	}

	if isFile {
		if flt := f.filter.MatchExclude(
			startPath,
			absStart,
			false,
		); flt != "" {
			_ = f.cache.HandleSkip(
				absStart,
				SkipExcluded,
				fmt.Errorf("matched %s", flt),
			)
			return nil, nil
		}
	}

	dirs, boundary, resolvedMeta, err := f.cache.ResolveDirChain(
		startDir,
		startMeta,
	)
	if err != nil {
		return nil, err
	}

	evalBoundaryRoot := f.cache.EvalBoundaryRoot(boundary.RootDir)
	var results []DiscoveredFile

	for _, dir := range slices.Backward(dirs) {
		found, skipErr := f.probeDir(
			dir,
			boundary.RootDir,
			evalBoundaryRoot,
			resolvedMeta,
			names,
			false,
		)
		if skipErr != nil {
			return nil, skipErr
		}
		results = append(results, found...)
	}

	return results, nil
}

func (f *fileFinder) FindNearest(
	startPath string,
	names ...string,
) (DiscoveredFile, bool, error) {
	startDir, absStart, isFile, startMeta, err := f.cache.ResolveStart(
		startPath,
	)
	if err != nil {
		return DiscoveredFile{}, false, err
	}

	if isFile {
		if flt := f.filter.MatchExclude(
			startPath,
			absStart,
			false,
		); flt != "" {
			_ = f.cache.HandleSkip(
				absStart,
				SkipExcluded,
				fmt.Errorf("matched %s", flt),
			)
			return DiscoveredFile{}, false, nil
		}
	}

	dirs, boundary, resolvedMeta, err := f.cache.ResolveDirChain(
		startDir,
		startMeta,
	)
	if err != nil {
		return DiscoveredFile{}, false, err
	}

	evalBoundaryRoot := f.cache.EvalBoundaryRoot(boundary.RootDir)

	for _, dir := range dirs {
		found, skipErr := f.probeDir(
			dir,
			boundary.RootDir,
			evalBoundaryRoot,
			resolvedMeta,
			names,
			true,
		)
		if skipErr != nil {
			return DiscoveredFile{}, false, skipErr
		}
		if len(found) > 0 {
			return found[0], true, nil
		}
	}

	return DiscoveredFile{}, false, nil
}

func (f *fileFinder) probeDir(
	dir string,
	boundaryRoot string,
	evalBoundaryRoot string,
	startMeta FileMeta,
	names []string,
	stopOnFirst bool,
) ([]DiscoveredFile, error) {
	relDir := relPathForFilter(boundaryRoot, dir)
	if flt := f.filter.MatchExclude(relDir, dir, true); flt != "" {
		_ = f.cache.HandleSkip(
			dir,
			SkipExcluded,
			fmt.Errorf("matched %s", flt),
		)
		return nil, nil
	}

	candidates, err := f.candidatesForDir(dir, boundaryRoot, names)
	if err != nil {
		return nil, err
	}

	var results []DiscoveredFile
	for _, name := range candidates {
		if name == "" {
			continue
		}
		candidatePath := filepath.Join(dir, name)
		relFile := relPathForFilter(boundaryRoot, candidatePath)
		if flt := f.filter.MatchExclude(
			relFile,
			candidatePath,
			false,
		); flt != "" {
			_ = f.cache.HandleSkip(
				candidatePath,
				SkipExcluded,
				fmt.Errorf("matched %s", flt),
			)
			continue
		}

		df, ok, skipErr := f.cache.ProbeCandidate(
			dir,
			name,
			candidatePath,
			evalBoundaryRoot,
			startMeta,
		)
		if skipErr != nil {
			return nil, skipErr
		}
		if ok {
			results = append(results, df)
			if stopOnFirst {
				return results, nil
			}
		}
	}
	return results, nil
}

func (f *fileFinder) candidatesForDir(
	dir string,
	boundaryRoot string,
	names []string,
) ([]string, error) {
	if len(names) > 0 {
		out := make([]string, 0, len(names))
		seen := make(map[string]struct{}, len(names))
		for _, n := range names {
			if n == "" {
				continue
			}
			if _, ok := seen[n]; ok {
				continue
			}
			seen[n] = struct{}{}
			out = append(out, n)
		}
		return out, nil
	}
	if !f.filter.HasInclude() {
		return nil, nil
	}

	entries, err := f.cache.DirEntries(dir)
	if err != nil {
		return nil, f.cache.HandleSkip(dir, SkipPermissionDenied, err)
	}

	var matched []string
	for _, de := range entries {
		if de.IsDir() {
			continue
		}
		full := filepath.Join(dir, de.Name())
		rel := relPathForFilter(boundaryRoot, full)
		if f.filter.IsIncluded(rel, full) {
			matched = append(matched, de.Name())
		}
	}
	return matched, nil
}

func relPathForFilter(boundaryRoot, fullPath string) string {
	if rel, err := filepath.Rel(boundaryRoot, fullPath); err == nil {
		if rel != "." {
			return rel
		}
	}
	return fullPath
}
