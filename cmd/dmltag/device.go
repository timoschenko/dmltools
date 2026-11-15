// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs/lookup"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
)

const implicitBuiltinsImport = "dml-builtins.dml"

var dmlcValueFlags = map[string]bool{
	"-D":                true,
	"-P":                true,
	"--dep":             true,
	"--dep-target":      true,
	"--warn":            true,
	"--nowarn":          true,
	"--simics-api":      true,
	"--max-errors":      true,
	"--no-compat":       true,
	"--breaking-change": true,
	"--split-c-file":    true,
}

func nextDmlcArg(
	tail []string,
	flagName string,
	requireNonEmpty bool,
) (string, string) {
	if len(tail) == 0 {
		if requireNonEmpty {
			return "", fmt.Sprintf(
				"flag '%s' requires a directory argument",
				flagName,
			)
		}
		return "", fmt.Sprintf("flag '%s' requires an argument", flagName)
	}
	val := tail[0]
	if requireNonEmpty && val == "" {
		return "", fmt.Sprintf(
			"flag '%s' requires a directory argument",
			flagName,
		)
	}
	return val, ""
}

// parseDmlcFlags extracts include directories from dmlc flags after '--'
// and rejects unexpected positional arguments.
func parseDmlcFlags(args []string) ([]string, string) {
	var includeDirs []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-I" || arg == "--include-dir":
			val, errMsg := nextDmlcArg(args[i+1:], arg, true)
			if errMsg != "" {
				return nil, errMsg
			}
			i++
			includeDirs = append(includeDirs, val)
		case strings.HasPrefix(arg, "-I="):
			val := arg[len("-I="):]
			if val == "" {
				return nil, "flag '-I=' requires a directory argument"
			}
			includeDirs = append(includeDirs, val)
		case strings.HasPrefix(arg, "-I") && len(arg) > 2:
			includeDirs = append(includeDirs, arg[2:])
		case strings.HasPrefix(arg, "--include-dir="):
			val := arg[len("--include-dir="):]
			if val == "" {
				return nil, "flag '--include-dir=' requires a directory argument"
			}
			includeDirs = append(includeDirs, val)
		case dmlcValueFlags[arg]:
			if _, errMsg := nextDmlcArg(args[i+1:], arg, false); errMsg != "" {
				return nil, errMsg
			}
			i++
		case strings.HasPrefix(arg, "-") && arg != "-":
			// Ignore other dmlc flags.
		default:
			return nil, fmt.Sprintf(
				"unexpected argument '%s' after '--' in device mode",
				arg,
			)
		}
	}

	return includeDirs, ""
}

type deviceResolver struct {
	cache           lookup.Cache
	finder          lookup.Finder
	includeSet      bool
	verbose         bool
	followAllImport bool
	rootDisplayPath string
	deviceDir       string
	boundaryRoot    string
	includeDirs     []string
	visited         map[string]struct{}
	stderr          io.Writer
}

func newDeviceResolver(
	cfg *CLIConfig,
	finder lookup.Finder,
	stderr io.Writer,
) (*deviceResolver, error) {
	cache := finder.Cache()
	displayDevice := filepath.Clean(cfg.Device)
	absDevice, err := cache.Abs(displayDevice)
	if err != nil {
		return nil, fmt.Errorf("resolving device file %s: %w", cfg.Device, err)
	}

	meta, err := cache.Stat(absDevice)
	if err != nil {
		return nil, fmt.Errorf("reading device file %s: %w", cfg.Device, err)
	}
	if meta.Mode.IsDir() {
		return nil, fmt.Errorf(
			"device path %s is a directory, expected a DML file",
			cfg.Device,
		)
	}

	boundary, err := cache.ResolveBoundary(absDevice)
	if err != nil {
		return nil, fmt.Errorf(
			"resolving workspace boundary for %s: %w",
			cfg.Device,
			err,
		)
	}

	includeDirs := make([]string, 0, len(cfg.DmlcIncludeDirs))
	for _, dir := range cfg.DmlcIncludeDirs {
		includeDirs = append(includeDirs, filepath.Clean(dir))
	}

	return &deviceResolver{
		cache:           cache,
		finder:          finder,
		includeSet:      cfg.IncludeSet,
		verbose:         cfg.Verbose,
		followAllImport: cfg.FollowAllImport,
		rootDisplayPath: displayDevice,
		deviceDir:       filepath.Dir(displayDevice),
		boundaryRoot:    cache.EvalBoundaryRoot(boundary.RootDir),
		includeDirs:     includeDirs,
		visited:         make(map[string]struct{}),
		stderr:          stderr,
	}, nil
}

func (r *deviceResolver) canonicalizePath(absPath string) string {
	if resolved, err := filepath.EvalSymlinks(absPath); err == nil {
		return filepath.Clean(resolved)
	}
	return filepath.Clean(absPath)
}

func (r *deviceResolver) isRegularFile(path string) bool {
	meta, err := r.cache.Stat(filepath.Clean(path))
	return err == nil && meta.Mode.IsRegular()
}

func (r *deviceResolver) resolveImportPath(
	importingPath, impName string,
) (string, bool) {
	importingDir := filepath.Dir(filepath.Clean(importingPath))
	searchDirs := make([]string, 0, len(r.includeDirs)+2)
	searchDirs = append(searchDirs, r.includeDirs...)
	searchDirs = append(searchDirs, r.deviceDir, importingDir)
	return dmlex.ResolveImportPath(
		importingPath,
		impName,
		searchDirs,
		r.isRegularFile,
	)
}

func (r *deviceResolver) shouldFollow(targetPath string, isRoot bool) bool {
	displayPath := filepath.Clean(targetPath)
	absPath, err := r.cache.Abs(displayPath)
	if err != nil {
		return false
	}
	evalPath := r.canonicalizePath(absPath)

	if _, seen := r.visited[evalPath]; seen {
		return false
	}
	r.visited[evalPath] = struct{}{}

	if isRoot {
		if r.finder.Filter().IsForceExcluded(displayPath) ||
			(evalPath != absPath && r.finder.Filter().IsForceExcluded(evalPath)) {
			return false
		}
		return true
	}

	inBoundary := r.cache.IsWithinBoundary(r.boundaryRoot, absPath)
	if !inBoundary {
		if !r.includeSet {
			return false
		}
		if !r.finder.Filter().IsIncluded(displayPath, absPath) &&
			(evalPath == absPath || !r.finder.Filter().IsIncluded(evalPath, evalPath)) {
			return false
		}
	}

	if filterName := r.finder.Filter().MatchExclude(
		displayPath,
		absPath,
		false,
	); filterName != "" {
		return false
	}

	if ignoreFile, err := r.finder.MatchIgnore(absPath); err == nil &&
		ignoreFile != "" {
		return false
	}

	return true
}

func collectDeviceFiles(
	cfg *CLIConfig,
	finder lookup.Finder,
	stderr io.Writer,
) ([]string, error) {
	resolver, err := newDeviceResolver(cfg, finder, stderr)
	if err != nil {
		return nil, err
	}

	if !resolver.shouldFollow(resolver.rootDisplayPath, true) {
		return nil, nil
	}

	var results []string
	queue := []string{resolver.rootDisplayPath}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		results = append(results, curr)

		cleanCurr := filepath.Clean(curr)
		//nolint:gosec // user-specified path
		data, err := os.ReadFile(cleanCurr)
		if err != nil {
			if cfg.Verbose {
				_, _ = fmt.Fprintf(
					stderr,
					"dmltag: reading %s: %v\n",
					curr,
					err,
				)
			}
			continue
		}

		buf := stream.NewHashedBufferFromString(curr, string(data))
		file := dmlex.ParseStream(buf)
		tree := file.Parser.Dml()
		hasDevice, _, _, imports := dmlex.ExtractDmlFileMeta(tree)

		isRoot := curr == resolver.rootDisplayPath
		if isRoot && !hasDevice {
			return nil, fmt.Errorf(
				"device file %s is missing device declaration",
				curr,
			)
		}
		if !isRoot && hasDevice {
			return nil, fmt.Errorf(
				"imported file %s cannot contain device declaration",
				curr,
			)
		}

		if isRoot {
			if resolved, ok := resolver.resolveImportPath(
				curr,
				implicitBuiltinsImport,
			); ok {
				if resolver.shouldFollow(resolved, false) {
					queue = append(queue, resolved)
				}
			}
		}

		for _, imp := range imports {
			resolved, ok := resolver.resolveImportPath(curr, imp.Name)
			if !ok {
				if resolver.followAllImport {
					return nil, fmt.Errorf(
						"imported file not found: %s (imported by %s:%d)",
						imp.Name,
						curr,
						imp.Line,
					)
				}
				if resolver.verbose {
					_, _ = fmt.Fprintf(
						stderr,
						"dmltag: imported file %q not found in %s:%d\n",
						imp.Name,
						curr,
						imp.Line,
					)
				}
				continue
			}
			if resolver.shouldFollow(resolved, false) {
				queue = append(queue, resolved)
			}
		}
	}

	return results, nil
}
