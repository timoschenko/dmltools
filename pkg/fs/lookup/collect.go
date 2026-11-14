// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"fmt"
	iofs "io/fs"
	"path/filepath"
	"slices"
)

type collectContext struct {
	rootSrc          string
	boundaryRoot     string
	evalBoundaryRoot string
	startMeta        FileMeta
	seen             map[string]bool
	files            *[]string
}

type ignoreStackInfo struct {
	stack            []IgnoreRuleSet
	boundaryRoot     string
	evalBoundaryRoot string
	startMeta        FileMeta
	startDir         string
	absStart         string
	isFile           bool
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
				igInfo, initErr := f.initCollectIgnoreStack(src)
				if initErr != nil {
					return nil, initErr
				}
				stack = igInfo.stack
				ctx.boundaryRoot = igInfo.boundaryRoot
				ctx.evalBoundaryRoot = igInfo.evalBoundaryRoot
				ctx.startMeta = igInfo.startMeta
				absSrc = igInfo.startDir
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

func (f *fileFinder) MatchIgnore(path string) (string, error) {
	if !f.ignoreEnabled() {
		return "", nil
	}
	info, err := f.initCollectIgnoreStack(path)
	if err != nil {
		return "", err
	}
	if len(info.stack) == 0 {
		return "", nil
	}
	return f.ignoreMatcher.Match(
		info.stack,
		info.absStart,
		!info.isFile,
	), nil
}

func (f *fileFinder) initCollectIgnoreStack(
	src string,
) (ignoreStackInfo, error) {
	startDir, absStart, isFile, startMeta, err := f.cache.ResolveStart(src)
	if err != nil {
		return ignoreStackInfo{}, err
	}

	dirs, boundary, resolvedMeta, err := f.cache.ResolveDirChain(
		startDir,
		startMeta,
	)
	if err != nil {
		return ignoreStackInfo{}, err
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
			return ignoreStackInfo{}, loadErr
		}
		stack = append(stack, loaded...)
	}

	return ignoreStackInfo{
		stack:            stack,
		boundaryRoot:     boundary.RootDir,
		evalBoundaryRoot: evalBoundaryRoot,
		startMeta:        resolvedMeta,
		startDir:         startDir,
		absStart:         absStart,
		isFile:           isFile,
	}, nil
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
