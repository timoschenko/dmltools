// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"fmt"
	iofs "io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/timoschenko/dmltools/pkg/fs"
)

// IgnoreCallback is called when a file or directory is ignored due to exclusion/inclusion filters.
type IgnoreCallback func(path string, isDir bool, filterName string)

// FileCollector handles file discovery, recursive directory traversal,
// and filtering based on include, exclude, extend-exclude, and force-exclude regular expressions.
type FileCollector struct {
	include       *regexp.Regexp
	exclude       *regexp.Regexp
	extendExclude *regexp.Regexp
	forceExclude  *regexp.Regexp
	onIgnore      IgnoreCallback
}

// SetOnIgnore sets the callback invoked when a file or directory is skipped.
func (fc *FileCollector) SetOnIgnore(cb IgnoreCallback) {
	fc.onIgnore = cb
}

// NewFileCollector compiles regular expressions from the CLI configuration.
func NewFileCollector(cfg *CLIConfig) (*FileCollector, error) {
	// TODO: AI slop
	fc := &FileCollector{}
	var err error

	if cfg.Include != "" {
		fc.include, err = regexp.Compile(cfg.Include)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid regular expression for --include: %w",
				err,
			)
		}
	}
	if cfg.Exclude != "" {
		fc.exclude, err = regexp.Compile(cfg.Exclude)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid regular expression for --exclude: %w",
				err,
			)
		}
	}
	if cfg.ExtendExclude != "" {
		fc.extendExclude, err = regexp.Compile(cfg.ExtendExclude)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid regular expression for --extend-exclude: %w",
				err,
			)
		}
	}
	if cfg.ForceExclude != "" {
		fc.forceExclude, err = regexp.Compile(cfg.ForceExclude)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid regular expression for --force-exclude: %w",
				err,
			)
		}
	}
	return fc, nil
}

// IsForceExcluded reports whether a file or directory matches the --force-exclude regex.
func (fc *FileCollector) IsForceExcluded(path string) bool {
	// TODO: AI slop
	if fc.forceExclude == nil {
		return false
	}
	normFile := fs.NormalizePath(path, false)
	if fc.forceExclude.MatchString(normFile) {
		return true
	}
	normDir := fs.NormalizePath(path, true)
	if fc.forceExclude.MatchString(normDir) {
		return true
	}
	if !filepath.IsAbs(path) {
		if abs, err := filepath.Abs(path); err == nil {
			if fc.forceExclude.MatchString(fs.NormalizePath(abs, false)) {
				return true
			}
			if fc.forceExclude.MatchString(fs.NormalizePath(abs, true)) {
				return true
			}
		}
	}
	return false
}

func (fc *FileCollector) matchExcludeFilter(normPath string) string {
	if fc.forceExclude != nil && fc.forceExclude.MatchString(normPath) {
		return "--force-exclude"
	}
	if fc.exclude != nil && fc.exclude.MatchString(normPath) {
		return "--exclude"
	}
	if fc.extendExclude != nil && fc.extendExclude.MatchString(normPath) {
		return "--extend-exclude"
	}
	return ""
}

func (fc *FileCollector) pathExcludeFilter(
	relToRoot, fullPath string,
	isDir bool,
) string {
	normRel := fs.NormalizePath(relToRoot, isDir)
	if f := fc.matchExcludeFilter(normRel); f != "" {
		return f
	}
	normFull := fs.NormalizePath(fullPath, isDir)
	if f := fc.matchExcludeFilter(normFull); f != "" {
		return f
	}
	if !filepath.IsAbs(fullPath) {
		if abs, err := filepath.Abs(fullPath); err == nil {
			normAbs := fs.NormalizePath(abs, isDir)
			if f := fc.matchExcludeFilter(normAbs); f != "" {
				return f
			}
		}
	}
	return ""
}

func (fc *FileCollector) isPathIncluded(relToRoot, fullPath string) bool {
	// TODO: AI slop
	if fc.include == nil {
		return true
	}
	normRel := fs.NormalizePath(relToRoot, false)
	if fc.include.MatchString(normRel) {
		return true
	}
	normFull := fs.NormalizePath(fullPath, false)
	if fc.include.MatchString(normFull) {
		return true
	}
	if !filepath.IsAbs(fullPath) {
		if abs, err := filepath.Abs(fullPath); err == nil {
			if fc.include.MatchString(fs.NormalizePath(abs, false)) {
				return true
			}
		}
	}
	return false
}

// Collect discovers, filters, deduplicates, and sorts file paths specified by sources.
func (fc *FileCollector) Collect(sources []string) ([]string, error) {
	// TODO: AI slop
	seen := make(map[string]bool)
	var files []string

	for _, src := range sources {
		if fc.IsForceExcluded(src) {
			if fc.onIgnore != nil {
				isDir := false
				if info, err := os.Stat(src); err == nil && info.IsDir() {
					isDir = true
				}
				cleanSrc := filepath.Clean(src)
				if isDir || strings.HasSuffix(cleanSrc, ".dml") {
					fc.onIgnore(cleanSrc, isDir, "--force-exclude")
				}
			}
			continue
		}

		info, err := os.Stat(src)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", src, err)
		}

		if info.IsDir() {
			cleanSrc := filepath.Clean(src)
			filter := fc.pathExcludeFilter(src, src, true)
			if filter != "" {
				if fc.onIgnore != nil {
					fc.onIgnore(cleanSrc, true, filter)
				}
				continue
			}

			err := filepath.WalkDir(
				src,
				func(path string, d iofs.DirEntry, walkErr error) error {
					if walkErr != nil {
						return fmt.Errorf(
							"reading file %s: %w",
							path,
							walkErr,
						)
					}

					rel, err := filepath.Rel(src, path)
					if err != nil {
						rel = path
					}

					cleanPath := filepath.Clean(path)

					if d.IsDir() {
						if path == src {
							return nil
						}
						filter := fc.pathExcludeFilter(rel, path, true)
						if filter != "" {
							if fc.onIgnore != nil {
								fc.onIgnore(cleanPath, true, filter)
							}
							return filepath.SkipDir
						}
						return nil
					}

					if d.Type()&iofs.ModeSymlink != 0 {
						targetInfo, err := os.Stat(path)
						if err != nil {
							// Broken symlink; skip
							return nil
						}
						if targetInfo.IsDir() {
							// Do not traverse into directory symlinks
							return nil
						}
					}

					isDml := strings.HasSuffix(cleanPath, ".dml")
					filter := fc.pathExcludeFilter(rel, path, false)
					if filter != "" {
						if isDml && fc.onIgnore != nil {
							fc.onIgnore(cleanPath, false, filter)
						}
						return nil
					}

					if !fc.isPathIncluded(rel, path) {
						if isDml && fc.onIgnore != nil {
							fc.onIgnore(cleanPath, false, "--include")
						}
						return nil
					}

					if !seen[cleanPath] {
						seen[cleanPath] = true
						files = append(files, cleanPath)
					}
					return nil
				},
			)
			if err != nil {
				return nil, err
			}
		} else {
			// Explicit file argument: bypasses --include, --exclude, and --extend-exclude.
			// Only --force-exclude applies.
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
