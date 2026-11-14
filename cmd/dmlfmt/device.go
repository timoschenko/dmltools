// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs/lookup"
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

// ImportRef describes an import statement location and target string.
type ImportRef struct {
	Name   string
	Line   int
	Column int
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

// extractDmlFileMeta inspects a parsed DML AST for a device declaration tag
// and all top-level import statements (including inside #if / #else).
func extractDmlFileMeta(
	tree antlr.ParseTree,
) (bool, int, int, []ImportRef) {
	dmlCtx, ok := tree.(dmlex.IDmlContext)
	if !ok || dmlCtx == nil {
		return false, 0, 0, nil
	}

	codeCtx := dmlCtx.DmlCode()
	if codeCtx == nil {
		return false, 0, 0, nil
	}

	var (
		hasDevice bool
		devLine   int
		devCol    int
	)
	if devTag := codeCtx.DeviceTag(); devTag != nil {
		hasDevice = true
		if tok := devTag.GetStart(); tok != nil {
			devLine = tok.GetLine()
			devCol = tok.GetColumn()
		}
	}

	var imports []ImportRef
	collectImportsFromStatements(codeCtx.DeviceStatements(), &imports)

	return hasDevice, devLine, devCol, imports
}

func collectImportsFromStatements(
	stmts dmlex.IDeviceStatementsContext,
	out *[]ImportRef,
) {
	if stmts == nil {
		return
	}
	for _, stmt := range stmts.AllDeviceStatement() {
		if stmt == nil {
			continue
		}
		if imp := stmt.ImportStatement(); imp != nil {
			if ref, ok := extractImportRef(imp); ok {
				*out = append(*out, ref)
			}
		}
		if ifCtx := stmt.ToplevelIf(); ifCtx != nil {
			collectImportsFromToplevelIf(ifCtx, out)
		}
	}
}

func collectImportsFromToplevelIf(
	ifCtx dmlex.IToplevelIfContext,
	out *[]ImportRef,
) {
	if ifCtx == nil {
		return
	}
	if ifStmt := ifCtx.ToplevelIfStatement(); ifStmt != nil {
		collectImportsFromStatements(ifStmt.DeviceStatements(), out)
	}
	if elseStmt := ifCtx.ToplevelElseStatement(); elseStmt != nil {
		collectImportsFromStatements(elseStmt.DeviceStatements(), out)
		collectImportsFromToplevelIf(elseStmt.ToplevelIf(), out)
	}
}

func extractImportRef(imp dmlex.IImportStatementContext) (ImportRef, bool) {
	lit := imp.STRING_LITERAL()
	if lit == nil {
		return ImportRef{}, false
	}
	raw := lit.GetText()
	if raw == "" {
		return ImportRef{}, false
	}
	var line, col int
	if sym := lit.GetSymbol(); sym != nil {
		line = sym.GetLine()
		col = sym.GetColumn()
	}
	return ImportRef{
		Name:   unquoteDmlString(raw),
		Line:   line,
		Column: col,
	}, true
}

func unquoteDmlString(raw string) string {
	if unquoted, err := strconv.Unquote(raw); err == nil {
		return unquoted
	}
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		return raw[1 : len(raw)-1]
	}
	return raw
}

// deviceImportResolver resolves imports and enforces workspace boundaries
// and ignore/exclude/include filters in --device mode.
type deviceImportResolver struct {
	cache           lookup.Cache
	finder          lookup.Finder
	ui              *UI
	includeSet      bool
	verbose         bool
	rootDisplayPath string
	deviceDir       string
	boundaryRoot    string
	includeDirs     []string
	visited         map[string]struct{}
}

func newDeviceImportResolver(
	cfg *CLIConfig,
	finder lookup.Finder,
	ui *UI,
) (*deviceImportResolver, error) {
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

	return &deviceImportResolver{
		cache:           cache,
		finder:          finder,
		ui:              ui,
		includeSet:      cfg.IncludeSet,
		verbose:         cfg.Verbose,
		rootDisplayPath: displayDevice,
		deviceDir:       filepath.Dir(displayDevice),
		boundaryRoot:    cache.EvalBoundaryRoot(boundary.RootDir),
		includeDirs:     includeDirs,
		visited:         make(map[string]struct{}),
	}, nil
}

func (r *deviceImportResolver) initialPaths() []string {
	if !r.shouldFormatAndFollow(r.rootDisplayPath, true) {
		return nil
	}
	return []string{r.rootDisplayPath}
}

func (r *deviceImportResolver) shouldFormatAndFollow(
	targetPath string,
	isRoot bool,
) bool {
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
			(evalPath != absPath &&
				r.finder.Filter().IsForceExcluded(evalPath)) {
			if r.verbose {
				r.ui.IgnoredFilter(
					displayPath,
					false,
					lookup.FilterForceExclude,
				)
			}
			return false
		}
		return true
	}

	inBoundary := r.cache.IsWithinBoundary(r.boundaryRoot, absPath)
	if !inBoundary {
		if !r.includeSet {
			if r.verbose {
				r.ui.Verbose(
					"%s ignored: outside workspace boundary\n",
					displayPath,
				)
			}
			return false
		}
		if !r.isOutsideIncluded(absPath, evalPath, displayPath) {
			if r.verbose {
				r.ui.Verbose(
					"%s ignored: outside workspace boundary and does not match --include\n",
					displayPath,
				)
			}
			return false
		}
	}

	if filterName := r.matchExcludeWithAncestors(
		absPath,
		evalPath,
		displayPath,
	); filterName != "" {
		if r.verbose {
			r.ui.IgnoredFilter(displayPath, false, filterName)
		}
		return false
	}

	if ignoreFile, err := r.finder.MatchIgnore(absPath); err == nil &&
		ignoreFile != "" {
		if r.verbose {
			r.ui.IgnoredFilter(displayPath, false, ignoreFile)
		}
		return false
	}

	return true
}

func (r *deviceImportResolver) canonicalizePath(absPath string) string {
	if resolved, err := filepath.EvalSymlinks(absPath); err == nil {
		return filepath.Clean(resolved)
	}
	return filepath.Clean(absPath)
}

func (r *deviceImportResolver) relUnderBoundary(
	evalPath string,
) (string, bool) {
	rel, err := filepath.Rel(r.boundaryRoot, evalPath)
	if err != nil {
		return "", false
	}
	if rel == "." || rel == ".." ||
		strings.HasPrefix(rel, ".."+string(filepath.Separator)) ||
		filepath.IsAbs(rel) {
		return "", false
	}
	return rel, true
}

func (r *deviceImportResolver) matchExcludeWithAncestors(
	absPath, evalPath, displayPath string,
) string {
	filter := r.finder.Filter()
	rel, underBoundary := r.relUnderBoundary(evalPath)
	if underBoundary {
		parts := strings.Split(filepath.ToSlash(rel), "/")
		for i := range len(parts) - 1 {
			subRel := strings.Join(parts[:i+1], "/")
			subFull := filepath.Join(r.boundaryRoot, filepath.FromSlash(subRel))
			if name := filter.MatchExclude(subRel, subFull, true); name != "" {
				return name
			}
		}
		return filter.MatchExclude(rel, displayPath, false)
	}

	if name := filter.MatchExclude(displayPath, absPath, false); name != "" {
		return name
	}
	if evalPath != absPath {
		return filter.MatchExclude(evalPath, evalPath, false)
	}
	return ""
}

func (r *deviceImportResolver) isOutsideIncluded(
	absPath, evalPath, displayPath string,
) bool {
	filter := r.finder.Filter()
	if filter.IsIncluded(displayPath, absPath) {
		return true
	}
	if evalPath != absPath && filter.IsIncluded(evalPath, evalPath) {
		return true
	}
	return false
}

func (r *deviceImportResolver) resolveImportPath(
	importingPath, impName string,
) (string, bool) {
	importingDir := filepath.Dir(filepath.Clean(importingPath))
	nativeImp := filepath.FromSlash(impName)

	if strings.HasPrefix(impName, "./") {
		candidate := filepath.Join(importingDir, nativeImp)
		if r.isRegularFile(candidate) {
			return candidate, true
		}
		return "", false
	}

	if strings.HasPrefix(impName, "../") {
		candidate := filepath.Join(importingDir, nativeImp)
		if r.isRegularFile(candidate) {
			return candidate, true
		}
	}

	searchDirs := make([]string, 0, len(r.includeDirs)+2)
	searchDirs = append(searchDirs, r.includeDirs...)
	searchDirs = append(searchDirs, r.deviceDir, importingDir)

	for _, dir := range searchDirs {
		verCandidate := filepath.Join(dir, dmlex.SupportedDmlVersion, nativeImp)
		if r.isRegularFile(verCandidate) {
			return verCandidate, true
		}
		directCandidate := filepath.Join(dir, nativeImp)
		if r.isRegularFile(directCandidate) {
			return directCandidate, true
		}
	}

	return "", false
}

func (r *deviceImportResolver) isRegularFile(path string) bool {
	meta, err := r.cache.Stat(filepath.Clean(path))
	return err == nil && meta.Mode.IsRegular()
}

func handleDeviceJobResult(
	job *OutputJob,
	resolver *deviceImportResolver,
	cfg *CLIConfig,
	ui *UI,
	paths *[]string,
) bool {
	isRoot := job.Path == resolver.rootDisplayPath
	if isRoot && !job.HasDeviceTag {
		ui.SyntaxError(
			job.Path,
			1,
			0,
			"EDEVICE: missing device declaration",
		)
		return true
	}
	if !isRoot && job.HasDeviceTag {
		ui.SyntaxError(
			job.Path,
			job.DeviceLine,
			job.DeviceColumn,
			"EDEVIMP: cannot import a device file",
		)
		return true
	}

	if isRoot {
		if resolved, ok := resolver.resolveImportPath(
			job.Path,
			implicitBuiltinsImport,
		); ok && resolver.shouldFormatAndFollow(resolved, false) {
			*paths = append(*paths, resolved)
		}
	}

	importFailed := false
	for _, imp := range job.Imports {
		resolved, ok := resolver.resolveImportPath(job.Path, imp.Name)
		if !ok {
			if cfg.FollowAllImport {
				ui.SyntaxError(
					job.Path,
					imp.Line,
					imp.Column,
					fmt.Sprintf(
						"EIMPORT: cannot find file to import: %s",
						imp.Name,
					),
				)
				importFailed = true
			} else if cfg.Verbose {
				ui.Verbose(
					"%s:%d: imported file not found, skipping: %s\n",
					job.Path,
					imp.Line,
					imp.Name,
				)
			}
			continue
		}
		if resolver.shouldFormatAndFollow(resolved, false) {
			*paths = append(*paths, resolved)
		}
	}

	return importFailed
}
