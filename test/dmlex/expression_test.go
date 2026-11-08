// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"fmt"
	"strings"
	"testing"

	lex "github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

type expressionCollector struct {
	*lex.BasedmlListener
	expr lex.IExpressionContext
}

func (v *expressionCollector) EnterConstantStatement(
	ctx *lex.ConstantStatementContext,
) {
	v.expr = ctx.Expression()
}

func stringifyExpr(expr lex.ISimplifiedExpr) string {
	if expr == nil {
		return "<nil>"
	}

	switch e := expr.(type) {
	case *lex.LiteralExpr:
		return e.Value.GetText()
	case *lex.BinaryExpr:
		return fmt.Sprintf(
			"(%s %s %s)",
			e.Op.GetText(),
			stringifyExpr(e.Left),
			stringifyExpr(e.Right),
		)
	case *lex.UnaryExpr:
		if e.IsPost {
			return fmt.Sprintf(
				"(post-%s %s)",
				e.Op.GetText(),
				stringifyExpr(e.Expr),
			)
		}
		return fmt.Sprintf("(%s %s)", e.Op.GetText(), stringifyExpr(e.Expr))
	case *lex.TernaryExpr:
		op := "if"
		if e.IsConst {
			op = "#if"
		}
		return fmt.Sprintf(
			"(%s %s %s %s)",
			op,
			stringifyExpr(e.Cond),
			stringifyExpr(e.TrueExpr),
			stringifyExpr(e.FalseExpr),
		)
	case *lex.FunctionCallExpr:
		args := make([]string, len(e.Args))
		for i, arg := range e.Args {
			args[i] = stringifyExpr(arg)
		}
		return fmt.Sprintf(
			"(call %s %s)",
			stringifyExpr(e.Expr),
			strings.Join(args, " "),
		)
	case *lex.ArrayAccessExpr:
		return fmt.Sprintf(
			"(index %s %s)",
			stringifyExpr(e.Expr),
			stringifyExpr(e.Index),
		)
	case *lex.BitSliceExpr:
		endian := "<no-endian>"
		if e.Endian != nil {
			endian = e.Endian.GetText()
		}
		right := "<nil>"
		if e.Right != nil {
			right = stringifyExpr(e.Right)
		}
		return fmt.Sprintf(
			"(slice %s %s %s %s)",
			stringifyExpr(e.Expr),
			stringifyExpr(e.Left),
			right,
			endian,
		)
	case *lex.MemberAccessExpr:
		// Build nested representation for chained member accesses.
		// a.b->c produces: (member -> (member . a b) c)
		result := stringifyExpr(e.Expr)
		for i := range e.Members {
			op := "."
			if e.IsArrow[i] {
				op = "->"
			}
			result = fmt.Sprintf(
				"(member %s %s %s)",
				op,
				result,
				e.Members[i].GetText(),
			)
		}
		return result
	case *lex.CastExpr:
		return fmt.Sprintf(
			"(cast %s %s)",
			e.Type.GetText(),
			stringifyExpr(e.Expr),
		)
	case *lex.SizeofExpr:
		return fmt.Sprintf("(sizeof %s)", stringifyExpr(e.Expr))
	case *lex.DefinedExpr:
		return fmt.Sprintf("(defined %s)", stringifyExpr(e.Expr))
	case *lex.StringifyExpr:
		return fmt.Sprintf("(stringify %s)", stringifyExpr(e.Expr))
	case *lex.NewExpr:
		if e.Size != nil {
			return fmt.Sprintf(
				"(new-array %s %s)",
				e.Type.GetText(),
				stringifyExpr(e.Size),
			)
		}
		return fmt.Sprintf("(new %s)", e.Type.GetText())
	case *lex.EachExpr:
		return fmt.Sprintf(
			"(each %s %s)",
			e.Ident.GetText(),
			stringifyExpr(e.List),
		)
	case *lex.SizeoftypeExpr:
		return fmt.Sprintf("(sizeoftype %s)", e.Type.GetText())
	case *lex.ParenExpr:
		return fmt.Sprintf("(paren %s)", stringifyExpr(e.Expr))
	case *lex.ArrayLiteralExpr:
		items := make([]string, len(e.Exprs))
		for i, item := range e.Exprs {
			items[i] = stringifyExpr(item)
		}
		return fmt.Sprintf("(array %s)", strings.Join(items, " "))
	default:
		return fmt.Sprintf("<unknown %T>", expr)
	}
}

func runASTVerificationTest(t *testing.T, exprStr string, expectedAST string) {
	code := fmt.Sprintf("dml 1.4;\nconstant TESTCASE = %s;\n", exprStr)

	// Parse stream
	input := antlr.NewInputStream(code)
	errorCounter := NewDmlLexErrorCounterListener(t, "test.dml", 0)
	file := lex.ParseStream(input, errorCounter)
	tree := file.Parser.Dml()

	if len(errorCounter.Errors) > 0 {
		t.Fatalf("Syntax errors encountered while parsing: %s", exprStr)
	}

	collector := &expressionCollector{}
	antlr.ParseTreeWalkerDefault.Walk(collector, tree)

	if collector.expr == nil {
		t.Fatalf(
			"Failed to find constant statement expression for: %s",
			exprStr,
		)
	}

	simplified := lex.WrapExpression(collector.expr)
	actualAST := stringifyExpr(simplified)

	if actualAST != expectedAST {
		t.Errorf(
			"AST Mismatch for '%s'\nExpected: %s\nActual:   %s",
			exprStr,
			expectedAST,
			actualAST,
		)
	}
}

func TestExpressionLiterals(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"123", "123"},
		{"0xFF", "0xFF"},
		{"0b1010", "0b1010"},
		{"3.14", "3.14"},
		{"'c'", "'c'"},
		{"\"string\"", "\"string\""},
		{"undefined", "undefined"},
		{"default", "default"},
		{"this", "this"},
		{"ident", "ident"},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}

func TestExpressionArithmetic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"-a", "(- a)"},
		{"+b", "(+ b)"},
		{"a + b", "(+ a b)"},
		{"a - b", "(- a b)"},
		{"a * b", "(* a b)"},
		{"a / b", "(/ a b)"},
		{"a % b", "(% a b)"},
		{"++a", "(++ a)"},
		{"--a", "(-- a)"},
		{"a++", "(post-++ a)"},
		{"a--", "(post--- a)"},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}

func TestExpressionBitwiseAndLogical(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"~a", "(~ a)"},
		{"!a", "(! a)"},
		{"a & b", "(& a b)"},
		{"a | b", "(| a b)"},
		{"a ^ b", "(^ a b)"},
		{"a << 2", "(<< a 2)"},
		{"a >> 2", "(>> a 2)"},
		{"a && b", "(&& a b)"},
		{"a || b", "(|| a b)"},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}

func TestExpressionRelationalAndEquality(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"a == b", "(== a b)"},
		{"a != b", "(!= a b)"},
		{"a < b", "(< a b)"},
		{"a > b", "(> a b)"},
		{"a <= b", "(<= a b)"},
		{"a >= b", "(>= a b)"},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}

func TestExpressionMemoryAndTypes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"&a", "(& a)"},
		{"*a", "(* a)"},
		{"sizeof a", "(sizeof a)"},
		{"sizeof(a)", "(sizeof (paren a))"},
		{"sizeoftype(int)", "(sizeoftype (int))"},
		{"defined a", "(defined a)"},
		{"stringify(a)", "(stringify a)"},
		{"new int", "(new int)"},
		{"new int[10]", "(new-array int 10)"},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}

func TestExpressionTernary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"a ? b : c", "(if a b c)"},
		{"a #? b #: c", "(#if a b c)"},
		{"a ? b : c ? d : e", "(if a b (if c d e))"}, // Right-associative
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}

func TestExpressionAccess(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"array[i]", "(index array i)"},
		{"array[i, le]", "(slice array i <nil> le)"},
		{"array[a : b]", "(slice array a b <no-endian>)"},
		{"array[a : b, be]", "(slice array a b be)"},
		{"obj.member", "(member . obj member)"},
		{"obj->member", "(member -> obj member)"},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}

func TestExpressionCalls(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"func()", "(call func )"},
		{"func(a, b)", "(call func a b)"},
		{"obj.func(a)", "(call (member . obj func) a)"},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}

func TestExpressionSpecialConstructs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"cast(a, int)", "(cast int a)"},
		{"cast(a, const int*)", "(cast constint* a)"},
		{"each i in (list)", "(each i list)"},
		{"[1, 2, 3]", "(array 1 2 3)"},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}

func TestExpressionPrecedence(t *testing.T) {
	t.Parallel()
	tests := []struct {
		expr     string
		expected string
	}{
		{"a + b * c", "(+ a (* b c))"},
		{"(a + b) * c", "(* (paren (+ a b)) c)"},
		{"a - b - c", "(- (- a b) c)"},
		{"a && b || c", "(|| (&& a b) c)"},
		{"a == b && c != d", "(&& (== a b) (!= c d))"},
	}

	for _, tc := range tests {
		t.Run(tc.expr, func(t *testing.T) {
			runASTVerificationTest(t, tc.expr, tc.expected)
		})
	}
}
