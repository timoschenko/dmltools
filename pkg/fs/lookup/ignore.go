// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

// IgnoreRuleSet holds the compiled ignore patterns parsed from a single ignore
// file (such as .gitignore or .ignore).
type IgnoreRuleSet struct {
	BaseDir    string
	SourceName string
	rules      []ignoreRule
}

// IgnoreMatcher parses, caches, and evaluates hierarchical ignore rule files
// (such as .gitignore and .ignore).
type IgnoreMatcher interface {
	LoadFile(df DiscoveredFile) (IgnoreRuleSet, error)
	Match(stack []IgnoreRuleSet, fullPath string, isDir bool) string
}

type cachedRuleSet struct {
	meta         FileMeta
	resolvedPath string
	baseDir      string
	sourceName   string
	ruleSet      IgnoreRuleSet
}

type gitignoreMatcher struct {
	fsys  FileSystem
	mu    sync.RWMutex
	cache map[string]cachedRuleSet
}

// NewIgnoreMatcher creates a standalone [IgnoreMatcher] with stat-validated
// rule caching.
func NewIgnoreMatcher(opts ...Option) IgnoreMatcher {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return newIgnoreMatcher(cfg.cacheCfg.fsys)
}

func newIgnoreMatcher(fsys FileSystem) *gitignoreMatcher {
	return &gitignoreMatcher{
		fsys:  fsys,
		cache: make(map[string]cachedRuleSet),
	}
}

func (m *gitignoreMatcher) LoadFile(
	df DiscoveredFile,
) (IgnoreRuleSet, error) {
	readPath := df.Path
	if readPath == "" {
		readPath = df.ResolvedPath
	}
	cacheKey := filepath.Clean(readPath)

	sourceName := df.Name
	if sourceName == "" {
		sourceName = filepath.Base(readPath)
	}
	baseDir := df.Dir
	if baseDir == "" {
		baseDir = filepath.Dir(readPath)
	}
	baseDir = filepath.Clean(baseDir)

	meta, err := m.resolveFileMeta(readPath, df)
	if err != nil {
		return IgnoreRuleSet{}, err
	}

	m.mu.RLock()
	entry, ok := m.cache[cacheKey]
	m.mu.RUnlock()
	if ok && entry.resolvedPath == df.ResolvedPath &&
		entry.baseDir == baseDir &&
		entry.sourceName == sourceName &&
		entry.meta.sameFingerprint(meta) {
		return entry.ruleSet, nil
	}

	data, err := m.fsys.ReadFile(readPath)
	if err != nil {
		return IgnoreRuleSet{}, err
	}

	rs := IgnoreRuleSet{
		BaseDir:    baseDir,
		SourceName: sourceName,
		rules:      parseIgnoreRules(data),
	}

	m.mu.Lock()
	m.cache[cacheKey] = cachedRuleSet{
		meta:         meta,
		resolvedPath: df.ResolvedPath,
		baseDir:      baseDir,
		sourceName:   sourceName,
		ruleSet:      rs,
	}
	m.mu.Unlock()

	return rs, nil
}

func (m *gitignoreMatcher) resolveFileMeta(
	readPath string,
	df DiscoveredFile,
) (FileMeta, error) {
	if df.Meta.Mode != 0 {
		return df.Meta, nil
	}
	if df.Mode != 0 || !df.ModTime.IsZero() || df.Size != 0 {
		return FileMeta{
			Name:    df.Name,
			Size:    df.Size,
			Mode:    df.Mode,
			ModTime: df.ModTime,
		}, nil
	}
	return m.fsys.Stat(readPath)
}

func (m *gitignoreMatcher) Match(
	stack []IgnoreRuleSet,
	fullPath string,
	isDir bool,
) string {
	if len(stack) == 0 {
		return ""
	}

	cleanFull := filepath.Clean(fullPath)
	for _, ancDir := range ancestorDirsUnderStack(stack, cleanFull) {
		if src := matchSinglePath(stack, ancDir, true); src != "" {
			return src
		}
	}

	return matchSinglePath(stack, cleanFull, isDir)
}

func ancestorDirsUnderStack(
	stack []IgnoreRuleSet,
	cleanFull string,
) []string {
	var ancestors []string
	for dir := filepath.Dir(cleanFull); dir != "" &&
		dir != "." &&
		dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		underAny := false
		for _, rs := range stack {
			if len(rs.rules) == 0 {
				continue
			}
			if _, ok := relUnderBaseDir(rs.BaseDir, dir); ok {
				underAny = true
				break
			}
		}
		if !underAny {
			break
		}
		ancestors = append(ancestors, dir)
	}
	slices.Reverse(ancestors)
	return ancestors
}

func matchSinglePath(
	stack []IgnoreRuleSet,
	cleanFull string,
	isDir bool,
) string {
	ignored := false
	matchedSource := ""

	var lastBase string
	var lastRel string
	var lastValid bool

	for _, rs := range stack {
		if len(rs.rules) == 0 {
			continue
		}
		var slashRel string
		var valid bool
		if rs.BaseDir == lastBase {
			slashRel = lastRel
			valid = lastValid
		} else {
			slashRel, valid = relUnderBaseDir(rs.BaseDir, cleanFull)
			lastBase = rs.BaseDir
			lastRel = slashRel
			lastValid = valid
		}
		if !valid {
			continue
		}

		for _, r := range rs.rules {
			if r.dirOnly && !isDir {
				continue
			}
			if r.regex.MatchString(slashRel) {
				ignored = !r.negate
				if ignored {
					matchedSource = rs.SourceName
				} else {
					matchedSource = ""
				}
			}
		}
	}

	return matchedSource
}

func relUnderBaseDir(baseDir, fullPath string) (string, bool) {
	rel, err := filepath.Rel(baseDir, fullPath)
	if err != nil {
		return "", false
	}
	if rel == "." || rel == ".." ||
		strings.HasPrefix(rel, ".."+string(filepath.Separator)) ||
		filepath.IsAbs(rel) {
		return "", false
	}
	return filepath.ToSlash(rel), true
}
