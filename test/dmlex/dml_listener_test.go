// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"os"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

func loadAllSyntaxExample(t *testing.T) string {
	t.Helper()
	path := "../../grammar/dml/examples/example-all-syntax-1.4.dml"
	code, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	return string(code)
}

func countParseTreeNodes(tree antlr.ParseTree) (rules, terminals, errors int) {
	if _, ok := tree.(antlr.ErrorNode); ok {
		return 0, 0, 1
	}
	if _, ok := tree.(antlr.TerminalNode); ok {
		return 0, 1, 0
	}
	rules = 1
	childCount := tree.GetChildCount()
	for i := range childCount {
		child, ok := tree.GetChild(i).(antlr.ParseTree)
		if !ok {
			continue
		}
		cr, ct, ce := countParseTreeNodes(child)
		rules += cr
		terminals += ct
		errors += ce
	}
	return rules, terminals, errors
}

type NodeAuditListener struct {
	*dmlex.BasedmlListener
	t             *testing.T
	parser        antlr.Parser
	EnterCount    int
	ExitCount     int
	TerminalCount int
	ErrorCount    int
	RuleStack     []int
	VisitedRules  map[int]int
	Errors        []string
}

func NewNodeAuditListener(
	t *testing.T,
	parser antlr.Parser,
) *NodeAuditListener {
	return &NodeAuditListener{
		BasedmlListener: new(dmlex.BasedmlListener),
		t:               t,
		parser:          parser,
		VisitedRules:    make(map[int]int),
	}
}

func (l *NodeAuditListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	if ctx == nil {
		l.t.Fatalf("EnterEveryRule called with nil ctx")
	}
	ruleIndex := ctx.GetRuleIndex()
	l.EnterCount++
	l.VisitedRules[ruleIndex]++
	l.RuleStack = append(l.RuleStack, ruleIndex)
}

func (l *NodeAuditListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	if ctx == nil {
		l.t.Fatalf("ExitEveryRule called with nil ctx")
	}
	if len(l.RuleStack) == 0 {
		l.t.Fatalf(
			"ExitEveryRule called with empty stack for rule %d",
			ctx.GetRuleIndex(),
		)
	}
	lastIdx := len(l.RuleStack) - 1
	top := l.RuleStack[lastIdx]
	l.RuleStack = l.RuleStack[:lastIdx]
	if top != ctx.GetRuleIndex() {
		ruleNames := l.parser.GetRuleNames()
		topName := ruleNames[top]
		ctxName := ruleNames[ctx.GetRuleIndex()]
		l.t.Fatalf(
			"Rule stack mismatch on exit: entered %s (%d) but exiting %s (%d)",
			topName,
			top,
			ctxName,
			ctx.GetRuleIndex(),
		)
	}
	l.ExitCount++
}

func (l *NodeAuditListener) VisitTerminal(node antlr.TerminalNode) {
	if node == nil {
		l.t.Fatalf("VisitTerminal called with nil node")
	}
	if node.GetSymbol() == nil {
		l.t.Fatalf("VisitTerminal called with nil symbol")
	}
	l.TerminalCount++
}

func (l *NodeAuditListener) VisitErrorNode(node antlr.ErrorNode) {
	l.ErrorCount++
	if node != nil {
		l.Errors = append(l.Errors, node.GetText())
	}
}

type TypedRuleDispatchListener struct {
	*dmlex.BasedmlListener
	counts map[string]int
}

func NewTypedRuleDispatchListener() *TypedRuleDispatchListener {
	return &TypedRuleDispatchListener{
		BasedmlListener: new(dmlex.BasedmlListener),
		counts:          make(map[string]int),
	}
}

func (l *TypedRuleDispatchListener) record(
	name string,
	ctx antlr.ParserRuleContext,
) {
	if ctx == nil {
		panic("nil context passed to " + name)
	}
	l.counts[name]++
}

// File tags & structure
func (l *TypedRuleDispatchListener) EnterDml(c *dmlex.DmlContext) {
	l.record("Dml", c)

	_ = c.DmlTag()
	_ = c.DmlCode()
	_ = c.EOF()
}

func (l *TypedRuleDispatchListener) EnterDmlCode(c *dmlex.DmlCodeContext) {
	l.record("DmlCode", c)

	_ = c.ProvisionalTag()
	_ = c.DeviceTag()
	_ = c.BitorderTag()
	_ = c.DeviceStatements()
}

func (l *TypedRuleDispatchListener) EnterDmlTag(c *dmlex.DmlTagContext) {
	l.record("DmlTag", c)

	_ = c.DML()
	_ = c.FLOAT_LITERAL()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterBitorderTag(
	c *dmlex.BitorderTagContext,
) {
	l.record("BitorderTag", c)

	_ = c.BITORDER()
	_ = c.ENDIAN()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterProvisionalTag(
	c *dmlex.ProvisionalTagContext,
) {
	l.record("ProvisionalTag", c)

	_ = c.PROVISIONAL()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterDeviceStatements(
	c *dmlex.DeviceStatementsContext,
) {
	l.record("DeviceStatements", c)

	_ = c.AllDeviceStatement()
	_ = c.DeviceStatement(0)
}

// Top-level statements
func (l *TypedRuleDispatchListener) EnterConstantStatement(
	c *dmlex.ConstantStatementContext,
) {
	l.record("ConstantStatement", c)

	_ = c.CONSTANT()
	_ = c.Ident()
	_ = c.ASSIGN()
	_ = c.Expression()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterImportStatement(
	c *dmlex.ImportStatementContext,
) {
	l.record("ImportStatement", c)

	_ = c.IMPORT()
	_ = c.STRING_LITERAL()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterExportStatement(
	c *dmlex.ExportStatementContext,
) {
	l.record("ExportStatement", c)

	_ = c.EXPORT()
	_ = c.Expression()
	_ = c.AS()
	_ = c.ExportAsExpression()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterLogGroupStatement(
	c *dmlex.LogGroupStatementContext,
) {
	l.record("LogGroupStatement", c)

	_ = c.LOGGROUP()
	_ = c.Ident()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterTypedefStatement(
	c *dmlex.TypedefStatementContext,
) {
	l.record("TypedefStatement", c)

	_ = c.TYPEDEF()
	_ = c.Cdecl()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterTemplateStatement(
	c *dmlex.TemplateStatementContext,
) {
	l.record("TemplateStatement", c)

	_ = c.TEMPLATE()
	_ = c.Objident()
	_ = c.LBRACE()
	_ = c.RBRACE()
	if isTemlp := c.IstemplateExpr(); isTemlp != nil {
		_ = isTemlp.IS()
		if isTemlList := isTemlp.IstemplateList(); isTemlList != nil {
			_ = isTemlList.AllCOMMA()
			_ = isTemlList.AllObjident()
			_ = isTemlList.LPAREN()
			_ = isTemlList.RPAREN()
			_ = isTemlList.Objident(0)
			_ = isTemlList.COMMA(0)
		}
	}
	_ = c.AllTemplateBody()
	_ = c.TemplateBody(0)
}

func (l *TypedRuleDispatchListener) EnterExternStatement(
	c *dmlex.ExternStatementContext,
) {
	l.record("ExternStatement", c)

	_ = c.EXTERN()
	_ = c.Cdecl()
	_ = c.SEMI()
	_ = c.TYPEDEF()
}

// Objects
func (l *TypedRuleDispatchListener) EnterObject(c *dmlex.ObjectContext) {
	l.record("Object", c)

	_ = c.Object_decl()
	_ = c.SessionOrSavedStatement()
	_ = c.Hook_decl()
}

func (l *TypedRuleDispatchListener) EnterObjectBody(
	c *dmlex.ObjectBodyContext,
) {
	l.record("ObjectBody", c)

	_ = c.SEMI()
	_ = c.ObjectDescSpec()
	_ = c.LBRACE()
	_ = c.RBRACE()
	_ = c.AllObjectStatement()
	_ = c.ObjectStatement(0)
}

func (l *TypedRuleDispatchListener) EnterObject_decl(
	c *dmlex.Object_declContext,
) {
	l.record("Object_decl", c)

	_ = c.REGISTER()
	_ = c.Objident()
	_ = c.ObjectBody()
	_ = c.Array_list()
	_ = c.SIZE()
	_ = c.SizeSpec()
	_ = c.OffsetSpec()
	_ = c.FIELD()
	_ = c.BitrangeDecl()
	_ = c.Object_type()
	_ = c.IN()
}

func (l *TypedRuleDispatchListener) EnterObjectStatement(
	c *dmlex.ObjectStatementContext,
) {
	l.record("ObjectStatement", c)

	_ = c.Object()
	_ = c.Param()
	_ = c.MethodUnified()
	_ = c.ObjectIf()
	_ = c.ErrorStatement()
	_ = c.IneachStatement()
}

// Methods
func (l *TypedRuleDispatchListener) EnterMethodUnifiedDeclaration(
	c *dmlex.MethodUnifiedDeclarationContext,
) {
	l.record("MethodUnifiedDeclaration", c)

	_ = c.METHOD()
	_ = c.Objident()
	_ = c.LPAREN()
	_ = c.RPAREN()
	_ = c.SEMI()
	_ = c.MethodQualifiers()
	_ = c.MethodParamsMaybeUntyped()
	_ = c.MethodOutParams()
	_ = c.THROWS()
}

func (l *TypedRuleDispatchListener) EnterMethodUnifiedDefinition(
	c *dmlex.MethodUnifiedDefinitionContext,
) {
	l.record("MethodUnifiedDefinition", c)

	_ = c.METHOD()
	_ = c.Objident()
	_ = c.LPAREN()
	_ = c.RPAREN()
	_ = c.RBRACE()
	_ = c.LBRACE()
	_ = c.LBRACE_NEW()
	_ = c.MethodQualifiers()
	_ = c.MethodParamsMaybeUntyped()
	_ = c.MethodOutParams()
	_ = c.THROWS()
	_ = c.AllStatement()
	_ = c.Statement(0)
	_ = c.DEFAULT()
	_ = c.DEFAULT_NEW()
}

// Types & Data declarations
func (l *TypedRuleDispatchListener) EnterStruct(c *dmlex.StructContext) {
	l.record("Struct", c)

	_ = c.STRUCT()
	_ = c.LBRACE()
	_ = c.RBRACE()
	_ = c.Ident()
	_ = c.AllStructDecls()
	_ = c.StructDecls(0)
}

func (l *TypedRuleDispatchListener) EnterUnion(c *dmlex.UnionContext) {
	l.record("Union", c)

	_ = c.UNION()
	_ = c.LBRACE()
	_ = c.RBRACE()
	_ = c.Ident()
	_ = c.AllStructDecls()
	_ = c.StructDecls(0)
}

func (l *TypedRuleDispatchListener) EnterLayout(c *dmlex.LayoutContext) {
	l.record("Layout", c)

	_ = c.LAYOUT()
	_ = c.STRING_LITERAL()
	_ = c.LBRACE()
	_ = c.RBRACE()
	_ = c.AllLayoutDecls()
	_ = c.LayoutDecls(0)
}

func (l *TypedRuleDispatchListener) EnterBitfields(c *dmlex.BitfieldsContext) {
	l.record("Bitfields", c)

	_ = c.BITFIELDS()
	_ = c.INTEGER_LITERAL()
	_ = c.LBRACE()
	_ = c.RBRACE()
	_ = c.AllBitfields_decls()
	_ = c.Bitfields_decls(0)
}

func (l *TypedRuleDispatchListener) EnterVariableStatement(
	c *dmlex.VariableStatementContext,
) {
	l.record("VariableStatement", c)

	_ = c.VariableDecl()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterSessionOrSavedStatement(
	c *dmlex.SessionOrSavedStatementContext,
) {
	l.record("SessionOrSavedStatement", c)

	_ = c.SessionOrSavedDecl()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterHook_decl(c *dmlex.Hook_declContext) {
	l.record("Hook_decl", c)

	_ = c.HOOK()
	_ = c.LPAREN()
	_ = c.RPAREN()
	_ = c.Ident()
	_ = c.SEMI()
	_ = c.Cdecl_list()
	_ = c.AllLBRACK()
	_ = c.LBRACK(0)
	_ = c.AllExpression()
	_ = c.Expression(0)
	_ = c.AllRBRACK()
	_ = c.RBRACK(0)
}

// Control flow & statements
func (l *TypedRuleDispatchListener) EnterCompoundStatement(
	c *dmlex.CompoundStatementContext,
) {
	l.record("CompoundStatement", c)

	_ = c.LBRACE()
	_ = c.RBRACE()
	_ = c.AllStatement()
	_ = c.Statement(0)
}

func (l *TypedRuleDispatchListener) EnterIfStatement(
	c *dmlex.IfStatementContext,
) {
	l.record("IfStatement", c)

	_ = c.IF()
	_ = c.LPAREN()
	_ = c.ConditionExpression()
	_ = c.RPAREN()
	_ = c.Statement()
	_ = c.ELSE()
	_ = c.ElseStatement()
	_ = c.HASH_IF()
	_ = c.HASH_ELSE()
}

func (l *TypedRuleDispatchListener) EnterSwitchStatement(
	c *dmlex.SwitchStatementContext,
) {
	l.record("SwitchStatement", c)

	_ = c.SWITCH()
	_ = c.LPAREN()
	_ = c.Expression()
	_ = c.RPAREN()
	_ = c.LBRACE()
	_ = c.SwitchBodyStatement()
	_ = c.RBRACE()
}

func (l *TypedRuleDispatchListener) EnterForLoopStatement(
	c *dmlex.ForLoopStatementContext,
) {
	l.record("ForLoopStatement", c)

	_ = c.FOR()
	_ = c.LPAREN()
	_ = c.AllSEMI()
	_ = c.SEMI(0)
	_ = c.RPAREN()
	_ = c.LoopStatement()
	_ = c.ForLoopInitStatement()
	_ = c.ConditionExpression()
	_ = c.ForLoopStepStatement()
}

func (l *TypedRuleDispatchListener) EnterWhileLoopStatement(
	c *dmlex.WhileLoopStatementContext,
) {
	l.record("WhileLoopStatement", c)

	_ = c.WHILE()
	_ = c.LPAREN()
	_ = c.ConditionExpression()
	_ = c.RPAREN()
	_ = c.LoopStatement()
}

func (l *TypedRuleDispatchListener) EnterDoWhileLoopStatement(
	c *dmlex.DoWhileLoopStatementContext,
) {
	l.record("DoWhileLoopStatement", c)

	_ = c.DO()
	_ = c.LoopStatement()
	_ = c.WHILE()
	_ = c.LPAREN()
	_ = c.ConditionExpression()
	_ = c.RPAREN()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterForeachStatement(
	c *dmlex.ForeachStatementContext,
) {
	l.record("ForeachStatement", c)

	_ = c.FOREACH()
	_ = c.Ident()
	_ = c.IN()
	_ = c.LPAREN()
	_ = c.Expression()
	_ = c.RPAREN()
	_ = c.LoopStatement()
	_ = c.HASH_FOREACH()
}

func (l *TypedRuleDispatchListener) EnterThrowStatement(
	c *dmlex.ThrowStatementContext,
) {
	l.record("ThrowStatement", c)

	_ = c.THROW()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterTryStatement(
	c *dmlex.TryStatementContext,
) {
	l.record("TryStatement", c)

	_ = c.TRY()
	_ = c.Statement()
	_ = c.CATCH()
	_ = c.CatchStatement()
}

func (l *TypedRuleDispatchListener) EnterReturnStatement(
	c *dmlex.ReturnStatementContext,
) {
	l.record("ReturnStatement", c)

	_ = c.RETURN()
	_ = c.SEMI()
	_ = c.Initializer()
}

func (l *TypedRuleDispatchListener) EnterAssertStatement(
	c *dmlex.AssertStatementContext,
) {
	l.record("AssertStatement", c)

	_ = c.ASSERT()
	_ = c.ConditionExpression()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterLogStatement(
	c *dmlex.LogStatementContext,
) {
	l.record("LogStatement", c)

	_ = c.LOG()
	_ = c.LogKind()
	_ = c.COLON()
	_ = c.LogFormatString()
	_ = c.SEMI()
	_ = c.AllCOMMA()
	_ = c.COMMA(0)
	_ = c.LogLevel()
	_ = c.LogGroup()
	_ = c.Expression_list_ntc_nonempty()
}

func (l *TypedRuleDispatchListener) EnterDeleteStatement(
	c *dmlex.DeleteStatementContext,
) {
	l.record("DeleteStatement", c)

	_ = c.DELETE()
	_ = c.Expression()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterAfterStatement(
	c *dmlex.AfterStatementContext,
) {
	l.record("AfterStatement", c)

	_ = c.AFTER()
	_ = c.AfterDelayCondition()
	_ = c.COLON()
	_ = c.Expression()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterWarningStatement(
	c *dmlex.WarningStatementContext,
) {
	l.record("WarningStatement", c)

	_ = c.WARNING()
	_ = c.Bracketed_string_literal()
	_ = c.SEMI()
}

func (l *TypedRuleDispatchListener) EnterSelectStatement(
	c *dmlex.SelectStatementContext,
) {
	l.record("SelectStatement", c)

	_ = c.Ident()
	_ = c.IN()
	_ = c.AllLPAREN()
	_ = c.LPAREN(0)
	_ = c.Expression()
	_ = c.AllRPAREN()
	_ = c.RPAREN(0)
	_ = c.WHERE()
	_ = c.ConditionExpression()
	_ = c.LoopStatement()
	_ = c.ElseStatement()
}

// Expressions
func (l *TypedRuleDispatchListener) EnterAssignStatement(
	c *dmlex.AssignStatementContext,
) {
	l.record("AssignStatement", c)

	_ = c.Assign_chain()
	_ = c.Tuple_literal()
	_ = c.ASSIGN()
	_ = c.Initializer()
}

func (l *TypedRuleDispatchListener) EnterFunctionCallExpr(
	c *dmlex.FunctionCallExprContext,
) {
	l.record("FunctionCallExpr", c)

	_ = c.Expression()
	_ = c.ArgumentsListEmpty()
	_ = c.ArgumentsList()
}

func (l *TypedRuleDispatchListener) EnterMemberAccessExpr(
	c *dmlex.MemberAccessExprContext,
) {
	l.record("MemberAccessExpr", c)

	_ = c.Expression()
	_ = c.AllExpressionAccessIdent()
	_ = c.ExpressionAccessIdent(0)
	_ = c.AllDOT()
	_ = c.DOT(0)
	_ = c.AllARROW()
	_ = c.ARROW(0)
}

func (l *TypedRuleDispatchListener) EnterArrayAccessExpr(
	c *dmlex.ArrayAccessExprContext,
) {
	l.record("ArrayAccessExpr", c)

	_ = c.AllExpression()
	_ = c.Expression(0)
	_ = c.LBRACK()
	_ = c.RBRACK()
}

func (l *TypedRuleDispatchListener) EnterBitSliceExpr(
	c *dmlex.BitSliceExprContext,
) {
	l.record("BitSliceExpr", c)

	_ = c.AllExpression()
	_ = c.Expression(0)
	_ = c.LBRACK()
	_ = c.BitSlice()
	_ = c.RBRACK()
}

func (l *TypedRuleDispatchListener) EnterCastExpr(c *dmlex.CastExprContext) {
	l.record("CastExpr", c)

	_ = c.CAST()
	_ = c.LPAREN()
	_ = c.Expression()
	_ = c.COMMA()
	_ = c.Ctypedecl()
	_ = c.RPAREN()
}

func (l *TypedRuleDispatchListener) EnterTernaryIfExpr(
	c *dmlex.TernaryIfExprContext,
) {
	l.record("TernaryIfExpr", c)

	_ = c.AllExpression()
	_ = c.Expression(0)
	_ = c.QUESTION()
	_ = c.COLON()
}

func TestDmlListener_VisitAllNodes_SyntaxExample(t *testing.T) {
	t.Parallel()

	code := loadAllSyntaxExample(t)
	file := ParseDmlFileFromText(t, code, "example-all-syntax-1.4.dml", 0)
	tree := file.Parser.Dml()

	audit := NewNodeAuditListener(t, file.Parser)
	antlr.ParseTreeWalkerDefault.Walk(audit, tree)

	expectedRules, expectedTerminals, expectedErrors := countParseTreeNodes(
		tree,
	)

	if audit.ErrorCount != 0 || expectedErrors != 0 {
		t.Fatalf(
			"unexpected error nodes: listener observed %d, tree has %d. Errors: %v",
			audit.ErrorCount,
			expectedErrors,
			audit.Errors,
		)
	}

	if audit.EnterCount != expectedRules {
		t.Errorf(
			"entered rules count mismatch: listener visited %d, expected %d",
			audit.EnterCount,
			expectedRules,
		)
	}

	if audit.ExitCount != expectedRules {
		t.Errorf(
			"exited rules count mismatch: listener exited %d, expected %d",
			audit.ExitCount,
			expectedRules,
		)
	}

	if audit.TerminalCount != expectedTerminals {
		t.Errorf(
			"terminal nodes count mismatch: listener visited %d, expected %d",
			audit.TerminalCount,
			expectedTerminals,
		)
	}

	if len(audit.RuleStack) != 0 {
		t.Errorf(
			"rule stack is not empty at end of walk: %d unclosed rules",
			len(audit.RuleStack),
		)
	}

	if audit.EnterCount == 0 {
		t.Errorf("expected enter count to be positive")
	}

	t.Logf(
		"Successfully visited all nodes in example-all-syntax-1.4.dml: %d rules entered/exited, %d terminal nodes, 0 error nodes",
		audit.EnterCount,
		audit.TerminalCount,
	)
}

func TestDmlListener_BaseListenerEmptyWalk(t *testing.T) {
	t.Parallel()

	code := loadAllSyntaxExample(t)
	file := ParseDmlFileFromText(t, code, "example-all-syntax-1.4.dml", 0)
	tree := file.Parser.Dml()

	base := new(dmlex.BasedmlListener)
	antlr.ParseTreeWalkerDefault.Walk(base, tree)
}

func TestDmlListener_TypedCallbackDispatch(t *testing.T) {
	t.Parallel()

	code := loadAllSyntaxExample(t)
	file := ParseDmlFileFromText(t, code, "example-all-syntax-1.4.dml", 0)
	tree := file.Parser.Dml()

	listener := NewTypedRuleDispatchListener()
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	expectedCallbacks := []string{
		"Dml",
		"DmlCode",
		"DmlTag",
		"BitorderTag",
		"ProvisionalTag",
		"DeviceStatements",
		"ConstantStatement",
		"ImportStatement",
		"ExportStatement",
		"LogGroupStatement",
		"TypedefStatement",
		"TemplateStatement",
		"ExternStatement",
		"Object",
		"ObjectBody",
		"Object_decl",
		"ObjectStatement",
		"MethodUnifiedDeclaration",
		"MethodUnifiedDefinition",
		"Struct",
		"Union",
		"Layout",
		"Bitfields",
		"VariableStatement",
		"SessionOrSavedStatement",
		"Hook_decl",
		"CompoundStatement",
		"IfStatement",
		"SwitchStatement",
		"ForLoopStatement",
		"WhileLoopStatement",
		"DoWhileLoopStatement",
		"ForeachStatement",
		"ThrowStatement",
		"TryStatement",
		"ReturnStatement",
		"AssertStatement",
		"LogStatement",
		"WarningStatement",
		"DeleteStatement",
		"AfterStatement",
		"SelectStatement",
		"AssignStatement",
		"FunctionCallExpr",
		"MemberAccessExpr",
		"ArrayAccessExpr",
		"BitSliceExpr",
		"CastExpr",
		"TernaryIfExpr",
	}

	for _, name := range expectedCallbacks {
		count := listener.counts[name]
		if count == 0 {
			t.Errorf(
				"expected typed callback %s to be invoked at least once, got 0",
				name,
			)
		} else {
			t.Logf("Typed callback %-25s invoked %4d times", name, count)
		}
	}
}

func TestDmlListener_GrammarRuleCoverage(t *testing.T) {
	t.Parallel()

	code := loadAllSyntaxExample(t)
	file := ParseDmlFileFromText(t, code, "example-all-syntax-1.4.dml", 0)
	tree := file.Parser.Dml()

	audit := NewNodeAuditListener(t, file.Parser)
	antlr.ParseTreeWalkerDefault.Walk(audit, tree)

	ruleNames := file.Parser.GetRuleNames()
	totalRules := len(ruleNames)
	visitedCount := len(audit.VisitedRules)

	coveragePct := (float64(visitedCount) / float64(totalRules)) * 100.0
	t.Logf(
		"Grammar rule coverage for example-all-syntax-1.4.dml: %d/%d (%.1f%%)",
		visitedCount,
		totalRules,
		coveragePct,
	)

	// In a comprehensive syntax test file, we expect 100% grammar rule coverage.
	if visitedCount < totalRules {
		var unvisited []string
		for idx, name := range ruleNames {
			if audit.VisitedRules[idx] == 0 {
				unvisited = append(unvisited, name)
			}
		}
		t.Errorf(
			"expected 100%% grammar rule coverage (all %d rules), got %d (%.1f%%). Unvisited rules: %v",
			totalRules,
			visitedCount,
			coveragePct,
			unvisited,
		)
	}

	// Verify key core rules are in visited set
	keyRules := []string{
		"dml",
		"dmlCode",
		"dmlTag",
		"bitorderTag",
		"deviceStatements",
		"constantStatement",
		"importStatement",
		"exportStatement",
		"logGroupStatement",
		"typedefStatement",
		"templateStatement",
		"object",
		"object_decl",
		"objectBody",
		"methodUnified",
		"struct",
		"layout",
		"bitfields",
		"statement",
		"ifStatement",
		"switchStatement",
		"forLoopStatement",
		"assignStatement",
		"expression",
	}

	ruleIndexMap := make(map[string]int, totalRules)
	for idx, name := range ruleNames {
		ruleIndexMap[name] = idx
	}

	for _, rule := range keyRules {
		idx, found := ruleIndexMap[rule]
		if !found {
			t.Errorf("key rule %q not found in parser ruleNames", rule)
			continue
		}
		if audit.VisitedRules[idx] == 0 {
			t.Errorf("key rule %q was not visited", rule)
		}
	}
}

func TestDmlListener_ParseListenerIntegration(t *testing.T) {
	t.Parallel()

	code := loadAllSyntaxExample(t)
	input := antlr.NewInputStream(code)
	errorCounter := NewDmlLexErrorCounterListener(
		t,
		"example-all-syntax-1.4.dml",
		0,
	)

	audit := NewNodeAuditListener(t, nil)
	file := dmlex.ParseStream(input, errorCounter, audit)
	audit.parser = file.Parser

	_ = file.Parser.Dml()

	if audit.ErrorCount != 0 {
		t.Fatalf("unexpected error nodes during parse: %v", audit.Errors)
	}

	if audit.EnterCount == 0 {
		t.Errorf("expected enter count to be positive with AddParseListener")
	}

	if audit.EnterCount != audit.ExitCount {
		t.Errorf(
			"enter and exit counts mismatch: entered %d, exited %d",
			audit.EnterCount,
			audit.ExitCount,
		)
	}

	if len(audit.RuleStack) != 0 {
		t.Errorf(
			"rule stack not empty after parse listening: %d unclosed rules",
			len(audit.RuleStack),
		)
	}

	t.Logf(
		"ParseListener visited %d rules entered/exited, %d terminal nodes",
		audit.EnterCount,
		audit.TerminalCount,
	)
}
