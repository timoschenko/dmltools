// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

type testMockToken struct {
	antlr.Token
	tokenType int
	text      string
}

func (t *testMockToken) GetTokenType() int {
	return t.tokenType
}

func (t *testMockToken) GetText() string {
	return t.text
}

func makeSessionOrSavedCtx(kindTok int) *dmlex.ObjectContext {
	objCtx := dmlex.NewEmptyObjectContext()
	stmtCtx := dmlex.NewEmptySessionOrSavedStatementContext()
	declCtx := dmlex.NewEmptySessionOrSavedDeclContext()
	kindCtx := dmlex.NewEmptyGlobal_decl_kindContext()
	kindCtx.AddTokenNode(&testMockToken{tokenType: kindTok})
	declCtx.AddChild(kindCtx)
	stmtCtx.AddChild(declCtx)
	objCtx.AddChild(stmtCtx)
	return objCtx
}

func TestCategorizeDeviceStatement(t *testing.T) {
	if got := categorizeDeviceStatement(nil); got != StmtKindUnknown {
		t.Errorf(
			"categorizeDeviceStatement(nil) = %v, want %v",
			got,
			StmtKindUnknown,
		)
	}

	empty := dmlex.NewEmptyDeviceStatementContext()
	if got := categorizeDeviceStatement(empty); got != StmtKindUnknown {
		t.Errorf(
			"categorizeDeviceStatement(empty) = %v, want %v",
			got,
			StmtKindUnknown,
		)
	}

	multi := dmlex.NewEmptyDeviceStatementContext()
	multi.AddChild(dmlex.NewEmptyImportStatementContext())
	multi.AddChild(dmlex.NewEmptyParamContext())
	if got := categorizeDeviceStatement(multi); got != StmtKindUnknown {
		t.Errorf(
			"categorizeDeviceStatement(multi) = %v, want %v",
			got,
			StmtKindUnknown,
		)
	}

	cases := []struct {
		name     string
		child    antlr.RuleContext
		expected StatementKind
	}{
		{"Import", dmlex.NewEmptyImportStatementContext(), StmtKindImport},
		{
			"LogGroup",
			dmlex.NewEmptyLogGroupStatementContext(),
			StmtKindLogGroup,
		},
		{
			"Constant",
			dmlex.NewEmptyConstantStatementContext(),
			StmtKindConstant,
		},
		{"Typedef", dmlex.NewEmptyTypedefStatementContext(), StmtKindTypedef},
		{"Extern", dmlex.NewEmptyExternStatementContext(), StmtKindExtern},
		{"Export", dmlex.NewEmptyExportStatementContext(), StmtKindExport},
		{"Param", dmlex.NewEmptyParamContext(), StmtKindParam},
		{
			"Template",
			dmlex.NewEmptyTemplateStatementContext(),
			StmtKindTemplate,
		},
		{"CCodeBlock", dmlex.NewEmptyCCodeBlockContext(), StmtKindCCode},
		{"Object", dmlex.NewEmptyObjectContext(), StmtKindObjectMajor},
		{"MethodUnified", dmlex.NewEmptyMethodUnifiedContext(), StmtKindMethod},
		{"ToplevelIf", dmlex.NewEmptyToplevelIfContext(), StmtKindToplevelIf},
		{"Ineach", dmlex.NewEmptyIneachStatementContext(), StmtKindIneach},
		{"Error", dmlex.NewEmptyErrorStatementContext(), StmtKindError},
		{"IsTemplate",
			dmlex.NewEmptyIstemplateStatementContext(),
			StmtKindIsTemplate,
		},
		{
			"Session",
			makeSessionOrSavedCtx(int(dmlex.TokenSession)),
			StmtKindSession,
		},
		{"Saved", makeSessionOrSavedCtx(int(dmlex.TokenSaved)), StmtKindSaved},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmt := dmlex.NewEmptyDeviceStatementContext()
			stmt.AddChild(tc.child)
			if got := categorizeDeviceStatement(stmt); got != tc.expected {
				t.Errorf(
					"categorizeDeviceStatement(%s) = %v, want %v",
					tc.name,
					got,
					tc.expected,
				)
			}
		})
	}
}

func TestCategorizeObjectStatement(t *testing.T) {
	if got := categorizeObjectStatement(nil); got != StmtKindUnknown {
		t.Errorf(
			"categorizeObjectStatement(nil) = %v, want %v",
			got,
			StmtKindUnknown,
		)
	}

	empty := dmlex.NewEmptyObjectStatementContext()
	if got := categorizeObjectStatement(empty); got != StmtKindUnknown {
		t.Errorf(
			"categorizeObjectStatement(empty) = %v, want %v",
			got,
			StmtKindUnknown,
		)
	}

	multi := dmlex.NewEmptyObjectStatementContext()
	multi.AddChild(dmlex.NewEmptyParamContext())
	multi.AddChild(dmlex.NewEmptyMethodUnifiedContext())
	if got := categorizeObjectStatement(multi); got != StmtKindUnknown {
		t.Errorf(
			"categorizeObjectStatement(multi) = %v, want %v",
			got,
			StmtKindUnknown,
		)
	}

	// Compact object (register/field with SEMI in object body)
	compactObjCtx := dmlex.NewEmptyObjectContext()
	objDeclCtx := dmlex.NewEmptyObject_declContext()
	bodyCtx := dmlex.NewEmptyObjectBodyContext()
	bodyCtx.AddTokenNode(&testMockToken{tokenType: int(dmlex.TokenSemi)})
	objDeclCtx.AddChild(bodyCtx)
	compactObjCtx.AddChild(objDeclCtx)

	cases := []struct {
		name     string
		child    antlr.RuleContext
		expected StatementKind
	}{
		{"Param", dmlex.NewEmptyParamContext(), StmtKindParam},
		{"MethodUnified", dmlex.NewEmptyMethodUnifiedContext(), StmtKindMethod},
		{
			"ObjectMajorEmpty",
			dmlex.NewEmptyObjectContext(),
			StmtKindObjectMajor,
		},
		{"ObjectCompact", compactObjCtx, StmtKindFieldOrRegCompact},
		{
			"IsTemplate",
			dmlex.NewEmptyIstemplateStatementContext(),
			StmtKindIsTemplate,
		},
		{"ObjectIf", dmlex.NewEmptyObjectIfContext(), StmtKindObjectIf},
		{"Ineach", dmlex.NewEmptyIneachStatementContext(), StmtKindIneach},
		{"Error", dmlex.NewEmptyErrorStatementContext(), StmtKindError},
		{
			"Session",
			makeSessionOrSavedCtx(int(dmlex.TokenSession)),
			StmtKindSession,
		},
		{"Saved", makeSessionOrSavedCtx(int(dmlex.TokenSaved)), StmtKindSaved},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stmt := dmlex.NewEmptyObjectStatementContext()
			stmt.AddChild(tc.child)
			if got := categorizeObjectStatement(stmt); got != tc.expected {
				t.Errorf(
					"categorizeObjectStatement(%s) = %v, want %v",
					tc.name,
					got,
					tc.expected,
				)
			}
		})
	}
}

func TestCategorizeTemplateBody(t *testing.T) {
	if got := categorizeTemplateBody(nil); got != StmtKindUnknown {
		t.Errorf(
			"categorizeTemplateBody(nil) = %v, want %v",
			got,
			StmtKindUnknown,
		)
	}

	empty := dmlex.NewEmptyTemplateBodyContext()
	if got := categorizeTemplateBody(empty); got != StmtKindUnknown {
		t.Errorf(
			"categorizeTemplateBody(empty) = %v, want %v",
			got,
			StmtKindUnknown,
		)
	}

	// TemplateBody wrapping ObjectStatement (e.g. Param)
	objStmt := dmlex.NewEmptyObjectStatementContext()
	objStmt.AddChild(dmlex.NewEmptyParamContext())
	bodyWithObj := dmlex.NewEmptyTemplateBodyContext()
	bodyWithObj.AddChild(objStmt)

	if got := categorizeTemplateBody(bodyWithObj); got != StmtKindParam {
		t.Errorf(
			"categorizeTemplateBody(bodyWithObj) = %v, want %v",
			got,
			StmtKindParam,
		)
	}

	// TemplateBody with SHARED hook
	bodyWithShared := dmlex.NewEmptyTemplateBodyContext()
	bodyWithShared.AddTokenNode(
		&testMockToken{tokenType: int(dmlex.TokenShared)},
	)
	if got := categorizeTemplateBody(
		bodyWithShared,
	); got != StmtKindSharedHook {
		t.Errorf(
			"categorizeTemplateBody(bodyWithShared) = %v, want %v",
			got,
			StmtKindSharedHook,
		)
	}
}

func TestComputeMemberSpacing(t *testing.T) {
	cases := []struct {
		input    int
		expected int
	}{
		{input: 0, expected: 0},
		{input: 1, expected: 1},
		{input: 2, expected: 1},
		{input: 5, expected: 1},
	}

	for _, tc := range cases {
		if got := computeMemberSpacing(tc.input); got != tc.expected {
			t.Errorf(
				"computeMemberSpacing(%d) = %d, want %d",
				tc.input,
				got,
				tc.expected,
			)
		}
	}
}
