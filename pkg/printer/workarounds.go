// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"
	"github.com/timoschenko/dmltools/pkg/printer/layout"

	"github.com/antlr4-go/antlr/v4"
)

// Deprecated: This is workaround
func hasNestedArrayLiteral_Workaround(node antlr.Tree) bool {
	if node == nil {
		return false
	}
	for _, child := range node.GetChildren() {
		switch c := child.(type) {
		case *dmlex.ArrayLiteralExprContext:
			return true
		case dmlex.IExpression_listContext,
			dmlex.IExpression_list_ntc_nonemptyContext,
			dmlex.IExpressionContext:
			if hasNestedArrayLiteral_Workaround(c) {
				return true
			}
		}
	}
	return false
}

//go:inline
func formatUnaryPreSpace(flag format.FormatFlag) format.FormatFlag {
	return flag.
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorAfterPriority(format.PriorityHigh)
}

// TODO: AI WA: Prevent adjacent unary operator tokens from fusing
// with operand starting tokens
//
// Deprecated: This is workaround
func ensureSpaceForUnaryExpr_Workaround(
	children []antlr.Tree,
) bool {
	if len(children) < 2 {
		return false
	}

	opNode, ok := children[0].(antlr.TerminalNode)
	if !ok {
		return false
	}
	exprNode, ok := children[1].(dmlex.IExpressionContext)
	if !ok {
		return false
	}

	return layout.WillTokensMalform(opNode.GetSymbol(), exprNode.GetStart())
}
