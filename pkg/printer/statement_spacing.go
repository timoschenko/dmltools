// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type StatementKind uint8

const (
	StmtKindUnknown StatementKind = iota
	StmtKindHeader
	// Compact Top-Level / Nested declarations
	StmtKindImport
	StmtKindLogGroup
	StmtKindConstant
	StmtKindTypedef
	StmtKindExtern
	StmtKindExport
	StmtKindParam
	StmtKindFieldOrRegCompact
	StmtKindIsTemplate
	StmtKindSession
	StmtKindSaved
	// Major Top-Level / Nested declarations
	StmtKindTemplate
	StmtKindObjectMajor
	StmtKindMethod
	StmtKindCCode
	StmtKindToplevelIf
	StmtKindObjectIf
	StmtKindIneach
	StmtKindError
	StmtKindSharedHook
)

//go:inline
func (k StatementKind) IsMajor() bool {
	switch k {
	case StmtKindTemplate, StmtKindObjectMajor, StmtKindMethod,
		StmtKindCCode, StmtKindToplevelIf, StmtKindObjectIf,
		StmtKindIneach, StmtKindError, StmtKindSharedHook:
		return true
	default:
		return false
	}
}

//go:inline
func (k StatementKind) IsCompact() bool {
	return !k.IsMajor() && k != StmtKindHeader && k != StmtKindUnknown
}

//go:inline
func getSessionOrSavedKind(node dmlex.IObjectContext) StatementKind {
	if node == nil {
		return StmtKindUnknown
	}
	stmt := node.SessionOrSavedStatement()
	if stmt == nil {
		return StmtKindUnknown
	}
	decl := stmt.SessionOrSavedDecl()
	if decl == nil {
		return StmtKindUnknown
	}
	kindCtx := decl.Global_decl_kind()
	if kindCtx == nil {
		return StmtKindUnknown
	}
	if kindCtx.SESSION() != nil {
		return StmtKindSession
	}
	if kindCtx.SAVED() != nil {
		return StmtKindSaved
	}
	return StmtKindUnknown
}

//go:inline
func categorizeDeviceStatement(
	child dmlex.IDeviceStatementContext,
) StatementKind {
	if child == nil || child.GetChildCount() != 1 {
		return StmtKindUnknown
	}
	switch node := child.GetChild(0).(type) {
	case dmlex.IImportStatementContext:
		return StmtKindImport
	case dmlex.ILogGroupStatementContext:
		return StmtKindLogGroup
	case dmlex.IConstantStatementContext:
		return StmtKindConstant
	case dmlex.ITypedefStatementContext:
		return StmtKindTypedef
	case dmlex.IExternStatementContext:
		return StmtKindExtern
	case dmlex.IExportStatementContext:
		return StmtKindExport
	case dmlex.IParamContext:
		return StmtKindParam
	case dmlex.ITemplateStatementContext:
		return StmtKindTemplate
	case dmlex.ICCodeBlockContext:
		return StmtKindCCode
	case dmlex.IObjectContext:
		if kind := getSessionOrSavedKind(node); kind != StmtKindUnknown {
			return kind
		}
		return StmtKindObjectMajor
	case dmlex.IMethodUnifiedContext:
		return StmtKindMethod
	case dmlex.IToplevelIfContext:
		return StmtKindToplevelIf
	case dmlex.IIneachStatementContext:
		return StmtKindIneach
	case dmlex.IErrorStatementContext:
		return StmtKindError
	case dmlex.IIstemplateStatementContext:
		return StmtKindIsTemplate
	default:
		return StmtKindUnknown
	}
}

//go:inline
func categorizeObjectStatement(
	child dmlex.IObjectStatementContext,
) StatementKind {
	if child == nil || child.GetChildCount() != 1 {
		return StmtKindUnknown
	}
	switch node := child.GetChild(0).(type) {
	case dmlex.IParamContext:
		return StmtKindParam
	case dmlex.IMethodUnifiedContext:
		return StmtKindMethod
	case dmlex.IObjectContext:
		if kind := getSessionOrSavedKind(node); kind != StmtKindUnknown {
			return kind
		}
		if node.GetChildCount() == 1 {
			if objDecl, ok := node.GetChild(0).(dmlex.IObject_declContext); ok {
				if body := objDecl.ObjectBody(); body != nil &&
					body.SEMI() != nil {
					return StmtKindFieldOrRegCompact
				}
			}
		}
		return StmtKindObjectMajor
	case dmlex.IIstemplateStatementContext:
		return StmtKindIsTemplate
	case dmlex.IObjectIfContext:
		return StmtKindObjectIf
	case dmlex.IIneachStatementContext:
		return StmtKindIneach
	case dmlex.IErrorStatementContext:
		return StmtKindError
	default:
		return StmtKindUnknown
	}
}

//go:inline
func categorizeTemplateBody(child dmlex.ITemplateBodyContext) StatementKind {
	if child == nil || child.GetChildCount() == 0 {
		return StmtKindUnknown
	}
	switch first := child.GetChild(0).(type) {
	case dmlex.IObjectStatementContext:
		if child.GetChildCount() == 1 {
			return categorizeObjectStatement(first)
		}
	case antlr.TerminalNode:
		if dmlex.GetTokenType(first.GetSymbol()) == dmlex.TokenShared {
			return StmtKindSharedHook
		}
	}
	return StmtKindUnknown
}

const (
	// Top-level statement spacing rules (number of empty lines).

	// Always 2 empty lines after header block.
	topLevelHeaderBlankLines = 2
	// Always 2 empty lines between major statements or major <-> compact.
	topLevelMajorBlankLines = 2
	// Different declaration categories must not be grouped together: range [1, 2].
	topLevelDiffCategoryMinBlankLines = 1
	topLevelDiffCategoryMaxBlankLines = 2
	// Same category: range [0, 2].
	topLevelSameCategoryMaxBlankLines = 2

	// Nested statement spacing rules (number of empty lines).

	// Always 1 empty line between methods or major blocks (1 line less than toplevel 2).
	nestedMajorBlankLines = 1
	// Different categories: range [1, 1] = 1 empty line.
	nestedDiffCategoryBlankLines = 1
	// Same category: range [0, 1] (1 line less than toplevel 2).
	nestedSameCategoryMaxBlankLines = 1

	// Spacing rules for method body statements (number of empty lines).

	// Same as nested: range [0, 1].
	methodBodyMaxBlankLines = 1

	// Spacing rules for type member declarations (struct, union, layout, bitfields).
	// Same as nested/method body: range [0, 1].
	memberMaxBlankLines = 1
)

//go:inline
func computeTopLevelSpacing(
	prevKind, currKind StatementKind,
	inputBlankLines int,
) int {
	if prevKind == StmtKindHeader {
		return topLevelHeaderBlankLines
	}
	if prevKind.IsMajor() || currKind.IsMajor() {
		return topLevelMajorBlankLines
	}
	if prevKind != currKind {
		return min(
			max(inputBlankLines, topLevelDiffCategoryMinBlankLines),
			topLevelDiffCategoryMaxBlankLines,
		)
	}
	return min(inputBlankLines, topLevelSameCategoryMaxBlankLines)
}

//go:inline
func computeNestedSpacing(
	prevKind, currKind StatementKind,
	inputBlankLines int,
) int {
	if prevKind.IsMajor() || currKind.IsMajor() {
		return nestedMajorBlankLines
	}
	if prevKind != currKind {
		return nestedDiffCategoryBlankLines
	}
	return min(inputBlankLines, nestedSameCategoryMaxBlankLines)
}

//go:inline
func computeMethodBodySpacing(inputBlankLines int) int {
	return min(inputBlankLines, methodBodyMaxBlankLines)
}

//go:inline
func computeMemberSpacing(inputBlankLines int) int {
	return min(inputBlankLines, memberMaxBlankLines)
}

//go:inline
func countInputBlankLines(prevStopTok, nextStartTok antlr.Token) int {
	if prevStopTok == nil || nextStartTok == nil {
		return 0
	}
	endLine := prevStopTok.GetLine() + strings.Count(
		prevStopTok.GetText(),
		"\n",
	)
	if !strings.HasSuffix(prevStopTok.GetText(), "\n") {
		endLine++
	}
	return max(0, nextStartTok.GetLine()-endLine)
}
