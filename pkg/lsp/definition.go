// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"
)

func isIdentChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '_' || ch == '$'
}

func isIdentOrDot(ch byte) bool {
	return isIdentChar(ch) || ch == '.'
}

func extractChainAtLineCol(line string, byteCol int) ([]string, bool) {
	if len(line) == 0 {
		return nil, false
	}

	col := min(max(byteCol, 0), len(line))
	if col > 0 &&
		(col == len(line) || !isIdentOrDot(line[col])) &&
		isIdentOrDot(line[col-1]) {
		col--
	}

	if col >= len(line) || !isIdentOrDot(line[col]) {
		return nil, false
	}

	start := col
	for start > 0 && isIdentOrDot(line[start-1]) {
		start--
	}

	end := col
	for end < len(line) && isIdentOrDot(line[end]) {
		end++
	}

	rawChain := line[start:end]
	trimmed := strings.Trim(rawChain, ".")
	if trimmed == "" {
		return nil, false
	}

	if trimmed == "is" {
		remainder := strings.TrimLeft(line[end:], " \t(")
		idx := 0
		for idx < len(remainder) && isIdentChar(remainder[idx]) {
			idx++
		}
		if idx > 0 {
			return []string{remainder[:idx]}, true
		}
		return nil, false
	}

	segments := strings.Split(rawChain, ".")
	currOffset := start
	activeIdx := 0
	var result []string
	for _, seg := range segments {
		segLen := len(seg)
		if segLen > 0 {
			result = append(result, seg)
			if col >= currOffset && col <= currOffset+segLen {
				activeIdx = len(result) - 1
			}
		}
		currOffset += segLen + 1
	}

	if len(result) == 0 {
		return nil, false
	}
	if activeIdx >= len(result) {
		activeIdx = len(result) - 1
	}
	return result[:activeIdx+1], true
}

func resolveImportPath(
	fsCache *FsCache,
	currentURIStr string,
	workspaceRoot string,
	importRelPath string,
) (uri.URI, bool) {
	currentPath, err := CanonicalPath(currentURIStr)
	var searchDirs []string
	if err == nil {
		searchDirs = append(searchDirs, filepath.Dir(currentPath))
	}

	if workspaceRoot != "" {
		if wsPath, wsErr := CanonicalPath(workspaceRoot); wsErr == nil {
			if len(searchDirs) == 0 || wsPath != searchDirs[0] {
				searchDirs = append(searchDirs, wsPath)
			}
			if currentPath == "" {
				currentPath = filepath.Join(wsPath, "root.dml")
			}
		}
	}

	if currentPath == "" || len(searchDirs) == 0 {
		return "", false
	}

	lookupCache := fsCache.LookupCache()
	resolved, ok := dmlex.ResolveImportPath(
		currentPath,
		importRelPath,
		searchDirs,
		func(candidate string) bool {
			meta, statErr := lookupCache.Stat(candidate)
			return statErr == nil && meta.Mode.IsRegular()
		},
	)
	if !ok {
		return "", false
	}

	return uri.File(resolved), true
}

func findInTransitiveImports(
	startFile *FileSymbols,
	fsCache *FsCache,
	workspaceRoot string,
	matchFn func(*FileSymbols) (*Symbol, bool),
) (*Symbol, *FileSymbols) {
	if startFile == nil {
		return nil, nil
	}

	visited := make(map[string]struct{})
	if startCanon, err := CanonicalPath(string(startFile.URI)); err == nil {
		visited[startCanon] = struct{}{}
	}

	queue := []*FileSymbols{startFile}
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, imp := range curr.Imports {
			targetURI, ok := resolveImportPath(
				fsCache,
				string(curr.URI),
				workspaceRoot,
				imp.Path,
			)
			if !ok {
				continue
			}

			canon, err := CanonicalPath(string(targetURI))
			if err != nil {
				canon = string(targetURI)
			}
			if _, seen := visited[canon]; seen {
				continue
			}
			visited[canon] = struct{}{}

			_, _, impSyms, err := fsCache.GetOrParse(string(targetURI))
			if err != nil || impSyms == nil {
				continue
			}

			if sym, found := matchFn(impSyms); found {
				return sym, impSyms
			}
			queue = append(queue, impSyms)
		}
	}

	return nil, nil
}

func searchImports(
	name string,
	fileSyms *FileSymbols,
	fsCache *FsCache,
	workspaceRoot string,
) (*Symbol, *FileSymbols) {
	return findInTransitiveImports(
		fileSyms,
		fsCache,
		workspaceRoot,
		func(impSyms *FileSymbols) (*Symbol, bool) {
			if impSyms.GlobalScope == nil {
				return nil, false
			}
			sym, ok := impSyms.GlobalScope.Symbols[name]
			return sym, ok
		},
	)
}

func findTemplateSymbol(
	tmplName string,
	ownerFile *FileSymbols,
	rootFile *FileSymbols,
	fsCache *FsCache,
	workspaceRoot string,
) (*Symbol, *FileSymbols) {
	lookupInFile := func(fsyms *FileSymbols) (*Symbol, *FileSymbols) {
		if fsyms == nil || fsyms.GlobalScope == nil {
			return nil, nil
		}
		if sym, ok := fsyms.GlobalScope.Symbols[tmplName]; ok &&
			sym.Kind == SymbolKindTemplate {
			return sym, fsyms
		}
		return findInTransitiveImports(
			fsyms,
			fsCache,
			workspaceRoot,
			func(impSyms *FileSymbols) (*Symbol, bool) {
				if impSyms.GlobalScope == nil {
					return nil, false
				}
				sym, ok := impSyms.GlobalScope.Symbols[tmplName]
				if ok && sym.Kind == SymbolKindTemplate {
					return sym, true
				}
				return nil, false
			},
		)
	}

	if sym, fsyms := lookupInFile(ownerFile); sym != nil {
		return sym, fsyms
	}
	if rootFile != nil && rootFile != ownerFile {
		return lookupInFile(rootFile)
	}
	return nil, nil
}

func isDeviceMemberKind(kind SymbolKind) bool {
	switch kind {
	case SymbolKindObject,
		SymbolKindParam,
		SymbolKindMethod,
		SymbolKindSession,
		SymbolKindSaved,
		SymbolKindHook:
		return true
	default:
		return false
	}
}

func lookupMemberInSymbol(
	sym *Symbol,
	memberName string,
	ownerFile *FileSymbols,
	rootFile *FileSymbols,
	fsCache *FsCache,
	workspaceRoot string,
	visitedTmpls map[string]struct{},
) (*Symbol, *FileSymbols) {
	if sym == nil {
		return nil, nil
	}
	if child, ok := sym.Children[memberName]; ok {
		return child, ownerFile
	}
	if sym.Scope != nil {
		if child, ok := sym.Scope.Symbols[memberName]; ok {
			return child, ownerFile
		}
	}
	if found, foundFile := lookupMemberInTemplates(
		sym.Templates,
		memberName,
		ownerFile,
		rootFile,
		fsCache,
		workspaceRoot,
		visitedTmpls,
	); found != nil {
		return found, foundFile
	}
	if ownerFile != nil && sym == ownerFile.DeviceSymbol {
		return findInTransitiveImports(
			ownerFile,
			fsCache,
			workspaceRoot,
			func(impSyms *FileSymbols) (*Symbol, bool) {
				if impSyms.GlobalScope == nil {
					return nil, false
				}
				s, ok := impSyms.GlobalScope.Symbols[memberName]
				if !ok || s == impSyms.DeviceSymbol ||
					!isDeviceMemberKind(s.Kind) {
					return nil, false
				}
				return s, true
			},
		)
	}
	return nil, nil
}

func lookupMemberInTemplates(
	templates []string,
	memberName string,
	ownerFile *FileSymbols,
	rootFile *FileSymbols,
	fsCache *FsCache,
	workspaceRoot string,
	visitedTmpls map[string]struct{},
) (*Symbol, *FileSymbols) {
	for _, tmplName := range templates {
		if _, seen := visitedTmpls[tmplName]; seen {
			continue
		}
		visitedTmpls[tmplName] = struct{}{}

		tmplSym, tmplFile := findTemplateSymbol(
			tmplName,
			ownerFile,
			rootFile,
			fsCache,
			workspaceRoot,
		)
		if tmplSym == nil {
			continue
		}
		if found, foundFile := lookupMemberInSymbol(
			tmplSym,
			memberName,
			tmplFile,
			rootFile,
			fsCache,
			workspaceRoot,
			visitedTmpls,
		); found != nil {
			return found, foundFile
		}
	}
	return nil, nil
}

func resolveBaseSymbol(
	base string,
	innermostScope *Scope,
	fileSyms *FileSymbols,
	fsCache *FsCache,
	workspaceRoot string,
) (*Symbol, *FileSymbols) {
	switch base {
	case "this":
		for sc := innermostScope; sc != nil; sc = sc.Parent {
			if sc.Owner != nil && sc.Kind == ScopeObject {
				return sc.Owner, fileSyms
			}
		}
		return nil, nil
	case "parent":
		for sc := innermostScope; sc != nil; sc = sc.Parent {
			if sc.Owner != nil && sc.Kind == ScopeObject {
				if sc.Owner.Parent != nil {
					return sc.Owner.Parent, fileSyms
				}
				return nil, nil
			}
		}
		return nil, nil
	case "dev":
		if fileSyms.DeviceSymbol != nil {
			return fileSyms.DeviceSymbol, fileSyms
		}
		return nil, nil
	}

	for sc := innermostScope; sc != nil; sc = sc.Parent {
		if s, ok := sc.Symbols[base]; ok {
			return s, fileSyms
		}
		if sc.Kind == ScopeObject && sc.Owner != nil &&
			len(sc.Owner.Templates) > 0 {
			if sym, ownerFile := lookupMemberInTemplates(
				sc.Owner.Templates,
				base,
				fileSyms,
				fileSyms,
				fsCache,
				workspaceRoot,
				make(map[string]struct{}),
			); sym != nil {
				return sym, ownerFile
			}
		}
		if sc.Kind == ScopeGlobal && len(fileSyms.DeviceTemplates) > 0 {
			if sym, ownerFile := lookupMemberInTemplates(
				fileSyms.DeviceTemplates,
				base,
				fileSyms,
				fileSyms,
				fsCache,
				workspaceRoot,
				make(map[string]struct{}),
			); sym != nil {
				return sym, ownerFile
			}
		}
	}

	return searchImports(base, fileSyms, fsCache, workspaceRoot)
}

func resolveChain(
	chain []string,
	innermostScope *Scope,
	fileSyms *FileSymbols,
	fsCache *FsCache,
	workspaceRoot string,
) *Symbol {
	if len(chain) == 0 {
		return nil
	}

	currentSym, currentFile := resolveBaseSymbol(
		chain[0],
		innermostScope,
		fileSyms,
		fsCache,
		workspaceRoot,
	)
	if currentSym == nil {
		return nil
	}

	for _, seg := range chain[1:] {
		switch seg {
		case "this":
			// Stays currentSym
		case "parent":
			if currentSym.Parent != nil {
				currentSym = currentSym.Parent
			} else {
				return nil
			}
		default:
			nextSym, nextFile := lookupMemberInSymbol(
				currentSym,
				seg,
				currentFile,
				fileSyms,
				fsCache,
				workspaceRoot,
				make(map[string]struct{}),
			)
			if nextSym == nil {
				return nil
			}
			currentSym = nextSym
			currentFile = nextFile
		}
	}

	return currentSym
}

// ResolveDefinition finds the target location for a definition request.
func ResolveDefinition(
	_ context.Context,
	fsCache *FsCache,
	uriStr string,
	workspaceRoot string,
	pos protocol.Position,
) (*protocol.Location, error) {
	buf, _, fileSyms, err := fsCache.GetOrParse(uriStr)
	if err != nil || fileSyms == nil {
		return nil, nil
	}

	for _, imp := range fileSyms.Imports {
		if rangeContains(imp.Range, pos) {
			targetURI, ok := resolveImportPath(
				fsCache,
				uriStr,
				workspaceRoot,
				imp.Path,
			)
			if ok {
				return &protocol.Location{
					URI: targetURI,
					Range: protocol.Range{
						Start: protocol.Position{Line: 0, Character: 0},
						End:   protocol.Position{Line: 0, Character: 0},
					},
				}, nil
			}
			return nil, nil
		}
	}

	var (
		line    string
		byteCol int
		hasLine bool
	)
	if entry, ok := fsCache.Get(uriStr); ok {
		line, byteCol, hasLine = entry.LineAndByteCol(pos)
	} else {
		line, byteCol, hasLine = LineAndByteCol(buf.GetRawString(), pos)
	}
	if !hasLine {
		return nil, nil
	}

	chain, ok := extractChainAtLineCol(line, byteCol)
	if !ok || len(chain) == 0 {
		return nil, nil
	}

	innermostScope := fileSyms.GlobalScope.EnclosingScope(pos)
	if innermostScope == nil {
		innermostScope = fileSyms.GlobalScope
	}

	sym := resolveChain(
		chain,
		innermostScope,
		fileSyms,
		fsCache,
		workspaceRoot,
	)
	if sym != nil {
		return &protocol.Location{
			URI:   sym.URI,
			Range: sym.NameRange,
		}, nil
	}

	return nil, nil
}
