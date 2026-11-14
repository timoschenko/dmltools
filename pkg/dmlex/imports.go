// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmlex

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// ImportRef describes an import statement location and target string.
type ImportRef struct {
	Name   string
	Line   int
	Column int
}

// UnquoteDmlString unquotes a DML string literal token text.
func UnquoteDmlString(raw string) string {
	if unquoted, err := strconv.Unquote(raw); err == nil {
		return unquoted
	}
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		return raw[1 : len(raw)-1]
	}
	return strings.Trim(raw, `"'`)
}

// ExtractDmlFileMeta inspects a parsed DML AST for a device declaration tag
// and all top-level import statements (including inside #if / #else).
func ExtractDmlFileMeta(
	tree antlr.ParseTree,
) (bool, int, int, []ImportRef) {
	dmlCtx, ok := tree.(IDmlContext)
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
	stmts IDeviceStatementsContext,
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
	ifCtx IToplevelIfContext,
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

func extractImportRef(imp IImportStatementContext) (ImportRef, bool) {
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
		Name:   UnquoteDmlString(raw),
		Line:   line,
		Column: col,
	}, true
}

// ResolveImportPath resolves a DML import path according to dmlc precedence:
//   - "./..." is resolved strictly relative to filepath.Dir(importingPath).
//   - "../..." is checked relative to filepath.Dir(importingPath) first, then falls through to searchDirs.
//   - Each dir in searchDirs is probed for "<dir>/<SupportedDmlVersion>/<impName>" first, then "<dir>/<impName>".
func ResolveImportPath(
	importingPath string,
	impName string,
	searchDirs []string,
	isRegularFile func(string) bool,
) (string, bool) {
	if impName == "" || isRegularFile == nil {
		return "", false
	}

	importingDir := filepath.Dir(filepath.Clean(importingPath))
	nativeImp := filepath.FromSlash(impName)

	if strings.HasPrefix(impName, "./") {
		candidate := filepath.Join(importingDir, nativeImp)
		if isRegularFile(candidate) {
			return candidate, true
		}
		return "", false
	}

	if strings.HasPrefix(impName, "../") {
		candidate := filepath.Join(importingDir, nativeImp)
		if isRegularFile(candidate) {
			return candidate, true
		}
	}

	for _, dir := range searchDirs {
		if dir == "" {
			continue
		}
		verCandidate := filepath.Join(dir, SupportedDmlVersion, nativeImp)
		if isRegularFile(verCandidate) {
			return verCandidate, true
		}
		directCandidate := filepath.Join(dir, nativeImp)
		if isRegularFile(directCandidate) {
			return directCandidate, true
		}
	}

	return "", false
}
