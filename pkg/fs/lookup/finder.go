// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"fmt"
	iofs "io/fs"
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
		rf, err := newFilterFromConfig(&cfg)
		if err != nil {
			return nil, err
		}
		f = rf
	}

	c := cfg.cache
	if c == nil {
		c = newCacheFromConfig(&cfg)
	}

	im := cfg.ignoreMatcher
	if im == nil && len(cfg.ignoreFiles) > 0 {
		im = newIgnoreMatcherFromConfig(&cfg)
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

type collectContext struct {
	rootSrc          string
	boundaryRoot     string
	evalBoundaryRoot string
	startMeta        FileMeta
	seen             map[string]bool
	files            *[]string
}

func (f *fileFinder) ignoreEnabled() bool {
	return f.ignoreMatcher != nil && len(f.ignoreFiles) > 0
}

func (f *fileFinder) Collect(sources []string) ([]string, error) {
	seen := make(map[string]bool)
	var files []string

	for _, src := range sources {
		if f.filter.IsForceExcluded(src) {
			isDir := false
			if info, err := f.cache.Stat(src); err == nil &&
				info.Mode.IsDir() {
				isDir = true
			}
			f.notifyIgnore(filepath.Clean(src), isDir, FilterForceExclude)
			continue
		}

		info, err := f.cache.Stat(src)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", src, err)
		}

		if info.Mode.IsDir() {
			cleanSrc := filepath.Clean(src)
			if flt := f.filter.MatchExclude(src, src, true); flt != "" {
				f.notifyIgnore(cleanSrc, true, flt)
				continue
			}

			ctx := &collectContext{
				rootSrc: src,
				seen:    seen,
				files:   &files,
			}
			var stack []IgnoreRuleSet
			absSrc := cleanSrc
			if abs, absErr := f.cache.Abs(src); absErr == nil {
				absSrc = filepath.Clean(abs)
			}

			if f.ignoreEnabled() {
				initStack, bRoot, evalBRoot, sMeta, startDir, initErr := f.initCollectIgnoreStack(
					src,
				)
				if initErr != nil {
					return nil, initErr
				}
				stack = initStack
				ctx.boundaryRoot = bRoot
				ctx.evalBoundaryRoot = evalBRoot
				ctx.startMeta = sMeta
				absSrc = startDir
				if len(stack) > 0 {
					if igName := f.ignoreMatcher.Match(
						stack,
						absSrc,
						true,
					); igName != "" {
						f.notifyIgnore(cleanSrc, true, igName)
						continue
					}
				}
			}

			if err := f.walkCollectDir(
				ctx,
				src,
				absSrc,
				stack,
				false,
			); err != nil {
				return nil, err
			}
		} else {
			cleanPath := filepath.Clean(src)
			if !seen[cleanPath] {
				seen[cleanPath] = true
				files = append(files, cleanPath)
			}
		}
	}

	slices.Sort(files)
	return files, nil
}

func (f *fileFinder) initCollectIgnoreStack(
	src string,
) ([]IgnoreRuleSet, string, string, FileMeta, string, error) {
	startDir, _, _, startMeta, err := f.cache.ResolveStart(src)
	if err != nil {
		return nil, "", "", FileMeta{}, "", err
	}

	dirs, boundary, resolvedMeta, err := f.cache.ResolveDirChain(
		startDir,
		startMeta,
	)
	if err != nil {
		return nil, "", "", FileMeta{}, "", err
	}

	evalBoundaryRoot := f.cache.EvalBoundaryRoot(boundary.RootDir)
	var stack []IgnoreRuleSet
	for _, dir := range slices.Backward(dirs) {
		loaded, loadErr := f.loadDirIgnoreRules(
			dir,
			boundary.RootDir,
			evalBoundaryRoot,
			resolvedMeta,
		)
		if loadErr != nil {
			return nil, "", "", FileMeta{}, "", loadErr
		}
		stack = append(stack, loaded...)
	}

	return stack, boundary.RootDir, evalBoundaryRoot, resolvedMeta, startDir, nil
}

func (f *fileFinder) loadDirIgnoreRules(
	dir string,
	boundaryRoot string,
	evalBoundaryRoot string,
	startMeta FileMeta,
) ([]IgnoreRuleSet, error) {
	found, skipErr := f.probeDir(
		dir,
		boundaryRoot,
		evalBoundaryRoot,
		startMeta,
		f.ignoreFiles,
		false,
	)
	if skipErr != nil {
		return nil, skipErr
	}
	if len(found) == 0 {
		return nil, nil
	}

	rules := make([]IgnoreRuleSet, 0, len(found))
	for _, df := range found {
		rs, loadErr := f.ignoreMatcher.LoadFile(df)
		if loadErr != nil {
			return nil, fmt.Errorf(
				"reading ignore file %s: %w",
				df.Path,
				loadErr,
			)
		}
		rules = append(rules, rs)
	}
	return rules, nil
}

func dirHasAnyCandidate(entries []iofs.DirEntry, names []string) bool {
	for _, de := range entries {
		if slices.Contains(names, de.Name()) {
			return true
		}
	}
	return false
}

func (f *fileFinder) walkCollectDir(
	ctx *collectContext,
	currDir string,
	absCurrDir string,
	stack []IgnoreRuleSet,
	probeSelf bool,
) error {
	entries, err := f.cache.DirEntries(currDir)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", currDir, err)
	}

	activeStack := stack
	if probeSelf && f.ignoreEnabled() &&
		dirHasAnyCandidate(entries, f.ignoreFiles) {
		loaded, loadErr := f.loadDirIgnoreRules(
			absCurrDir,
			ctx.boundaryRoot,
			ctx.evalBoundaryRoot,
			ctx.startMeta,
		)
		if loadErr != nil {
			return loadErr
		}
		if len(loaded) > 0 {
			activeStack = append(slices.Clone(stack), loaded...)
		}
	}

	for _, d := range entries {
		path := filepath.Join(currDir, d.Name())
		absPath := filepath.Join(absCurrDir, d.Name())
		rel, relErr := filepath.Rel(ctx.rootSrc, path)
		if relErr != nil {
			rel = path
		}
		cleanPath := filepath.Clean(path)

		if d.IsDir() {
			if flt := f.filter.MatchExclude(rel, path, true); flt != "" {
				f.notifyIgnore(cleanPath, true, flt)
				continue
			}
			if f.ignoreMatcher != nil && len(activeStack) > 0 {
				if igName := f.ignoreMatcher.Match(
					activeStack,
					absPath,
					true,
				); igName != "" {
					f.notifyIgnore(cleanPath, true, igName)
					continue
				}
			}
			if err := f.walkCollectDir(
				ctx,
				path,
				absPath,
				activeStack,
				true,
			); err != nil {
				return err
			}
			continue
		}

		if d.Type()&iofs.ModeSymlink != 0 {
			targetInfo, statErr := f.cache.Stat(path)
			if statErr != nil || targetInfo.Mode.IsDir() {
				continue
			}
		}

		if flt := f.filter.MatchExclude(rel, path, false); flt != "" {
			f.notifyIgnore(cleanPath, false, flt)
			continue
		}

		if f.ignoreMatcher != nil && len(activeStack) > 0 {
			if igName := f.ignoreMatcher.Match(
				activeStack,
				absPath,
				false,
			); igName != "" {
				f.notifyIgnore(cleanPath, false, igName)
				continue
			}
		}

		if !f.filter.IsIncluded(rel, path) {
			f.notifyIgnore(cleanPath, false, FilterInclude)
			continue
		}

		if !ctx.seen[cleanPath] {
			ctx.seen[cleanPath] = true
			*ctx.files = append(*ctx.files, cleanPath)
		}
	}

	return nil
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
