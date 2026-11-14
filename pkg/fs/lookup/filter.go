// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lookup

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/timoschenko/dmltools/pkg/fs"
)

// Filter evaluates path inclusion and exclusion rules.
type Filter interface {
	IsForceExcluded(path string) bool
	MatchExclude(relToRoot, fullPath string, isDir bool) string
	IsIncluded(relToRoot, fullPath string) bool
	HasInclude() bool
}

type regexFilter struct {
	fsys          FileSystem
	include       *regexp.Regexp
	exclude       *regexp.Regexp
	extendExclude *regexp.Regexp
	forceExclude  *regexp.Regexp
}

// NewFilter creates a standalone [Filter] compiled from regular-expression options.
func NewFilter(opts ...Option) (Filter, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return newFilterFromConfig(&cfg)
}

func newFilterFromConfig(cfg *config) (*regexFilter, error) {
	var err error
	var incRe, excRe, extRe, forceRe *regexp.Regexp

	if cfg.includePattern != "" {
		incRe, err = regexp.Compile(cfg.includePattern)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid regular expression for %s: %w",
				FilterInclude,
				err,
			)
		}
	}
	if cfg.excludePattern != "" {
		excRe, err = regexp.Compile(cfg.excludePattern)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid regular expression for %s: %w",
				FilterExclude,
				err,
			)
		}
	}
	if cfg.extendExcludePattern != "" {
		extRe, err = regexp.Compile(cfg.extendExcludePattern)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid regular expression for %s: %w",
				FilterExtendExclude,
				err,
			)
		}
	}
	if cfg.forceExcludePattern != "" {
		forceRe, err = regexp.Compile(cfg.forceExcludePattern)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid regular expression for %s: %w",
				FilterForceExclude,
				err,
			)
		}
	}

	return &regexFilter{
		fsys:          cfg.fsys,
		include:       incRe,
		exclude:       excRe,
		extendExclude: extRe,
		forceExclude:  forceRe,
	}, nil
}

func (f *regexFilter) HasInclude() bool {
	return f.include != nil
}

func (f *regexFilter) IsForceExcluded(path string) bool {
	if f.forceExclude == nil {
		return false
	}
	normFile := fs.NormalizePath(path, false)
	if f.forceExclude.MatchString(normFile) {
		return true
	}
	normDir := fs.NormalizePath(path, true)
	if f.forceExclude.MatchString(normDir) {
		return true
	}
	if !filepath.IsAbs(path) {
		if abs, err := f.fsys.Abs(path); err == nil {
			if f.forceExclude.MatchString(fs.NormalizePath(abs, false)) {
				return true
			}
			if f.forceExclude.MatchString(fs.NormalizePath(abs, true)) {
				return true
			}
		}
	}
	return false
}

func (f *regexFilter) matchExcludeFilter(normPath string) string {
	if f.forceExclude != nil && f.forceExclude.MatchString(normPath) {
		return FilterForceExclude
	}
	if f.exclude != nil && f.exclude.MatchString(normPath) {
		return FilterExclude
	}
	if f.extendExclude != nil && f.extendExclude.MatchString(normPath) {
		return FilterExtendExclude
	}
	return ""
}

func (f *regexFilter) MatchExclude(
	relToRoot, fullPath string,
	isDir bool,
) string {
	normRel := fs.NormalizePath(relToRoot, isDir)
	if name := f.matchExcludeFilter(normRel); name != "" {
		return name
	}
	normFull := fs.NormalizePath(fullPath, isDir)
	if name := f.matchExcludeFilter(normFull); name != "" {
		return name
	}
	if !filepath.IsAbs(fullPath) {
		if abs, err := f.fsys.Abs(fullPath); err == nil {
			normAbs := fs.NormalizePath(abs, isDir)
			if name := f.matchExcludeFilter(normAbs); name != "" {
				return name
			}
		}
	}
	return ""
}

func (f *regexFilter) IsIncluded(relToRoot, fullPath string) bool {
	if f.include == nil {
		return true
	}
	normRel := fs.NormalizePath(relToRoot, false)
	if f.include.MatchString(normRel) {
		return true
	}
	normFull := fs.NormalizePath(fullPath, false)
	if f.include.MatchString(normFull) {
		return true
	}
	if !filepath.IsAbs(fullPath) {
		if abs, err := f.fsys.Abs(fullPath); err == nil {
			if f.include.MatchString(fs.NormalizePath(abs, false)) {
				return true
			}
		}
	}
	return false
}
