// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"slices"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"

	"github.com/antlr4-go/antlr/v4"
)

func isCompoundStatement(node antlr.Tree) bool {
	if _, ok := node.(dmlex.ICompoundStatementContext); ok {
		return true
	}

	for _, child := range node.GetChildren() {
		if _, ok := child.(dmlex.ICompoundStatementContext); ok {
			return true
		}
	}

	return false
}

const (
	nestedStatement      = false
	addCompoundIfNeeeded = true
)

// g4-tag: statement
// g4-tag: compoundStatement
func (p *Printer) printStatement(
	node antlr.Tree,
) error {
	return p.printStatement2(node, nestedStatement)
}

func (p *Printer) printStatement2(
	node antlr.Tree, ensureCompound bool,
) error {
	p.newContext(p.settings.TokenFormatDefaults())
	defer p.freeContext()

	stack := make([]antlr.Tree, 0, 16)
	for _, nestedChild := range slices.Backward(node.GetChildren()) {
		stack = append(stack, nestedChild)
		if ensureCompound && isCompoundStatement(nestedChild) {
			ensureCompound = false
		}
	}

	if ensureCompound {
		if err := p.pushSynthecic(
			dmlex.NewFakeLBrace(p.source),
		); err != nil {
			return err
		}
	}

	for len(stack) > 0 {
		lastIdx := len(stack) - 1
		rawChild := stack[lastIdx]
		stack = stack[:lastIdx]

		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenSemi,
					dmlex.TokenLbrace,
					dmlex.TokenRbrace:
					return true
				}
				return false
			}, "token is not supported")
			if tt == dmlex.TokenRbrace {
				if err := p.pushTokenWithCustomFormat(
					child,
					func(ff format.FormatFlag) format.FormatFlag {
						return ff.WithSeparatorAfter(format.SepMandatoryWrap)
					},
				); err != nil {
					return err
				}
			} else {
				if err := p.pushTokenRangeFromRightiest(child); err != nil {
					return err
				}
			}
		case dmlex.ICompoundStatementContext:
			if err := p.printCompoundStatement(child); err != nil {
				return err
			}
		case dmlex.IStatementContext:
			for _, nestedChild := range slices.Backward(child.GetChildren()) {
				stack = append(stack, nestedChild)
			}
		case dmlex.IAfterStatementContext:
			if err := p.printAfterStatement(child); err != nil {
				return err
			}
		case dmlex.IAssertStatementContext:
			if err := p.printAssertStatement(child); err != nil {
				return err
			}
		case dmlex.IAssignopContext:
			if err := p.printAssignop(child); err != nil {
				return err
			}
		case dmlex.IAssignStatementContext:
			if err := p.printAssignStatement(child); err != nil {
				return err
			}
		case dmlex.IBreakStatementContext:
			if err := p.printBreakStatement(child); err != nil {
				return err
			}
		case dmlex.IContinueStatementContext:
			if err := p.printContinueStatement(child); err != nil {
				return err
			}
		case dmlex.IDeleteStatementContext:
			if err := p.printDeleteStatement(child); err != nil {
				return err
			}
		case dmlex.IDoWhileLoopStatementContext:
			if err := p.printDoWhileLoopStatement(child); err != nil {
				return err
			}
		case dmlex.IErrorStatementContext, dmlex.IWarningStatementContext:
			if err := p.printErrorOrWarningStatement(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		case dmlex.IForeachStatementContext:
			if err := p.printForeachStatement(child); err != nil {
				return err
			}
		case dmlex.IForLoopStatementContext:
			if err := p.printForLoopStatement(child); err != nil {
				return err
			}
		case dmlex.IGotoStatementContext:
			if err := p.printGotoStatement(child); err != nil {
				return err
			}
		case dmlex.IIfStatementContext:
			if err := p.printIfStatement(child); err != nil {
				return err
			}
		case dmlex.ILogStatementContext:
			if err := p.printLogStatement(child); err != nil {
				return err
			}
		case dmlex.IReturnStatementContext:
			if err := p.printReturnStatement(child); err != nil {
				return err
			}
		case dmlex.ISelectStatementContext:
			if err := p.printSelectStatement(child); err != nil {
				return err
			}
		case dmlex.ISwitchStatementContext:
			if err := p.printSwitchStatement(child); err != nil {
				return err
			}
		case dmlex.IThrowStatementContext:
			if err := p.printThrowStatement(child); err != nil {
				return err
			}
		case dmlex.ITryStatementContext:
			if err := p.printTryStatement(child); err != nil {
				return err
			}
		case dmlex.IVariableStatementContext:
			if err := p.printVariableStatement(child); err != nil {
				return err
			}
		case dmlex.IWhileLoopStatementContext:
			if err := p.printWhileLoopStatement(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	if ensureCompound {
		if err := p.pushSynthecic(
			dmlex.NewFakeRBrace(p.source),
		); err != nil {
			return err
		}
	}

	return nil
}

// g4-tag: compoundStatement
func (p *Printer) printCompoundStatement(
	node dmlex.ICompoundStatementContext,
) error {
	if node.LBRACE() != nil {
		if err := p.pushTokenRangeFromRightiest(node.LBRACE()); err != nil {
			return err
		}
	}
	hadStatement := false
	var prevStopToken antlr.Token
	for _, stmt := range node.AllStatement() {
		startIdx := stmt.GetStart().GetTokenIndex()
		var nextStartTok antlr.Token
		if p.lastProcessedTokenIndex < TokenIndex(startIdx) {
			nextStartTok = p.extractToken(p.lastProcessedTokenIndex)
		} else {
			nextStartTok = stmt.GetStart()
		}
		if hadStatement {
			blankLines := countInputBlankLines(prevStopToken, nextStartTok)
			emptyLines := computeMethodBodySpacing(blankLines)
			if err := p.writeNewLines(emptyLines + 1); err != nil {
				return err
			}
		}
		if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
			return err
		}
		if err := p.printStatement(stmt); err != nil {
			return err
		}
		hadStatement = true
		prevStopIdx := p.gatherTokenIndexFromRightiestToken(stmt.GetStop())
		prevStopToken = p.extractToken(prevStopIdx)
	}
	if node.RBRACE() != nil {
		if err := p.pushTokenWithCustomFormat(
			node.RBRACE(),
			func(ff format.FormatFlag) format.FormatFlag {
				return ff.WithSeparatorAfter(format.SepMandatoryWrap)
			},
		); err != nil {
			return err
		}
	}
	return nil
}

// g4-tag: identList
func (p *Printer) printIdentList(node dmlex.IIdentListContext) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				return tt == dmlex.TokenComma
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IIdentContext:
			if err := p.printIdent(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: afterDelayCondition
func (p *Printer) printAfterDelayCondition(
	node dmlex.IAfterDelayConditionContext,
) error {
	endGroup := p.LayoutEngine.BeginGreedyGroup("AfterDelayCondition")
	defer endGroup()

	endSoftIndent := p.LayoutEngine.BeginSoftIndent()
	defer endSoftIndent()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenIdentifier,
					dmlex.TokenArrow,
					dmlex.TokenLparen,
					dmlex.TokenRparen:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}

		case dmlex.IIdentListContext:
			if err := p.printIdentList(child); err != nil {
				return err
			}

		case dmlex.IIdentContext:
			if err := p.printIdent(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: afterStatement
func (p *Printer) printAfterStatement(node dmlex.IAfterStatementContext) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenColon: p.context.GetFormatFlag(dmlex.TokenColon).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityHigh),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenAfter,
					dmlex.TokenColon,
					dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IAfterDelayConditionContext:
			if err := p.printAfterDelayCondition(child); err != nil {
				return err
			}

		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: assignop
func (p *Printer) printAssignop(node dmlex.IAssignopContext) error {
	var endSoftIndentHandle func()
	endSoftIndent := func() {
		if endSoftIndentHandle != nil {
			endSoftIndentHandle()
			endSoftIndentHandle = nil
		}
	}
	defer endSoftIndent()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			assert(endSoftIndentHandle == nil, "unexpected")
			endSoftIndentHandle = p.LayoutEngine.BeginSoftIndent()
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: assignStatement
// g4-tag: assign_chain
// g4-tag: tuple_literal
func (p *Printer) printAssignStatement(
	node dmlex.IAssignStatementContext,
) error {
	var endSoftIndentHandle func()
	endSoftIndent := func() {
		if endSoftIndentHandle != nil {
			endSoftIndentHandle()
			endSoftIndentHandle = nil
		}
	}
	defer endSoftIndent()

	stack := make([]antlr.Tree, 0, 16)
	for _, nestedChild := range slices.Backward(node.GetChildren()) {
		stack = append(stack, nestedChild)
	}

	for len(stack) > 0 {
		lastIdx := len(stack) - 1
		rawChild := stack[lastIdx]
		stack = stack[:lastIdx]

		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if tt == dmlex.TokenAssign && endSoftIndentHandle == nil {
				endSoftIndentHandle = p.LayoutEngine.BeginSoftIndent()
			}
		case dmlex.IAssign_chainContext:
			for _, nestedChild := range slices.Backward(child.GetChildren()) {
				stack = append(stack, nestedChild)
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		case dmlex.IInitializerContext:
			if err := p.printInitializer(child); err != nil {
				return err
			}
		case dmlex.ITuple_literalContext:
			if err := p.printExpression_listExpr(child, rootNode); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: breakStatement
func (p *Printer) printBreakStatement(node dmlex.IBreakStatementContext) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenBreak, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: continueStatement
func (p *Printer) printContinueStatement(
	node dmlex.IContinueStatementContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenContinue, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: loopStatement
func (p *Printer) printLoopStatement(node dmlex.ILoopStatementContext) error {
	assertLazy(
		func() bool {
			return node.GetChildCount() == 1
		}, "loopStatement expects exactly 1 child",
	)

	if err := p.printStatement2(
		node.GetChild(0), addCompoundIfNeeeded,
	); err != nil {
		return err
	}

	return nil
}

// g4-tag: elseStatement
func (p *Printer) printElseStatement(node dmlex.IElseStatementContext) error {
	assertLazy(
		func() bool {
			return node.GetChildCount() == 1
		}, "loopStatement expects exactly 1 child",
	)

	child := node.GetChild(0)
	for _, grandChild := range child.GetChildren() {
		switch grandChild.(type) {
		case dmlex.IIfStatementContext, dmlex.ISelectStatementContext:
			// print as is
			return p.printStatement(child)
		}
	}

	if err := p.printStatement2(
		child, addCompoundIfNeeeded,
	); err != nil {
		return err
	}

	return nil
}

// g4-tag: doWhileLoopStatement
func (p *Printer) printDoWhileLoopStatement(
	node dmlex.IDoWhileLoopStatementContext,
) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenWhile: p.context.GetFormatFlag(dmlex.TokenWhile).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	return p.printWhileLoopStatement(node)
}

// g4-tag: foreachStatement
func (p *Printer) printForeachStatement(
	node dmlex.IForeachStatementContext,
) error {
	var endExprGroupHandle func()
	endExprGroup := func() {
		if endExprGroupHandle != nil {
			endExprGroupHandle()
			endExprGroupHandle = nil
		}
	}
	defer endExprGroup()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenForeach,
					dmlex.TokenHashForeach,
					dmlex.TokenIn,
					dmlex.TokenLparen,
					dmlex.TokenRparen:
					return true
				}
				return false
			}, "token is not supported")

			if tt == dmlex.TokenLparen {
				assert(
					endExprGroupHandle == nil,
					"foreachExpression must be nil",
				)
				endExprGroupHandle = p.LayoutEngine.BeginWrapGroup(
					"foreachExpression",
				)
			}

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenRparen {
				endExprGroup()
			}

		case dmlex.IIdentContext:
			if err := p.printIdent(child); err != nil {
				return err
			}

		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}

		case dmlex.ILoopStatementContext:
			if err := p.printLoopStatement(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: innerForLoopStatement
func (p *Printer) printInnerForLoopStatement(
	node dmlex.IInnerForLoopStatementContext,
) error {
	if node.AssignStatement() != nil {
		return p.printAssignStatement(node.AssignStatement())
	}
	if node.Assignop() != nil {
		return p.printAssignop(node.Assignop())
	}
	if node.Expression() != nil {
		return p.printExpression(node.Expression())
	}
	return nil
}

// g4-tag: forLoopInitStatement
func (p *Printer) printForLoopInitStatement(
	node dmlex.IForLoopInitStatementContext,
) error {
	if node.VariableDecl() != nil {
		return p.printVariableDecl(node.VariableDecl())
	}
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				return tt == dmlex.TokenComma
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IInnerForLoopStatementContext:
			if err := p.printInnerForLoopStatement(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: forLoopStepStatement
func (p *Printer) printForLoopStepStatement(
	node dmlex.IForLoopStepStatementContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				return tt == dmlex.TokenComma
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IInnerForLoopStatementContext:
			if err := p.printInnerForLoopStatement(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: forLoopStatement
func (p *Printer) printForLoopStatement(
	node dmlex.IForLoopStatementContext,
) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenSemi: p.context.GetFormatFlag(dmlex.TokenSemi).
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithBreakAfterPriority(format.PriorityDefault).
			WithSeparatorAfterPriority(format.PriorityAboveNormal),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	var endHeaderGroupHandle func()
	endHeaderGroup := func() {
		if endHeaderGroupHandle != nil {
			endHeaderGroupHandle()
			endHeaderGroupHandle = nil
		}
	}
	defer endHeaderGroup()

	semiCount := 0
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenFor,
					dmlex.TokenLparen,
					dmlex.TokenRparen,
					dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")

			if tt == dmlex.TokenLparen {
				assert(endHeaderGroupHandle == nil, "forHeader must be nil")
				endHeaderGroupHandle = p.LayoutEngine.BeginWrapGroup(
					"forHeader",
				)
			}

			if tt == dmlex.TokenSemi {
				semiCount++
				if semiCount == 1 && node.ConditionExpression() == nil {
					if err := p.pushTokenWithCustomFormat(
						child,
						func(ff format.FormatFlag) format.FormatFlag {
							return ff.WithSeparatorAfter(format.SepNone)
						},
					); err != nil {
						return err
					}
					continue
				}
			}

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenRparen {
				endHeaderGroup()
			}

		case dmlex.IForLoopInitStatementContext:
			if err := p.printForLoopInitStatement(child); err != nil {
				return err
			}

		case dmlex.IConditionExpressionContext:
			if err := p.printConditionExpression(child); err != nil {
				return err
			}

		case dmlex.IForLoopStepStatementContext:
			if err := p.printForLoopStepStatement(child); err != nil {
				return err
			}

		case dmlex.ILoopStatementContext:
			if err := p.printLoopStatement(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: gotoStatement
func (p *Printer) printGotoStatement(node dmlex.IGotoStatementContext) error {
	var endSoftIndentHandle func()
	endSoftIndent := func() {
		if endSoftIndentHandle != nil {
			endSoftIndentHandle()
			endSoftIndentHandle = nil
		}
	}
	defer endSoftIndent()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenGoto, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if tt == dmlex.TokenSemi {
				endSoftIndent()
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if tt == dmlex.TokenGoto && node.Ident() != nil {
				assert(endSoftIndentHandle == nil, "unexpected")
				endSoftIndentHandle = p.LayoutEngine.BeginSoftIndent()
			}
		case dmlex.IIdentContext:
			if err := p.printIdent(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: ifStatement
// g4-tag: elseStatement
func (p *Printer) printIfStatement(node dmlex.IIfStatementContext) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenElse: p.context.GetFormatFlag(dmlex.TokenElse).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal),
		dmlex.TokenHashElse: p.context.GetFormatFlag(dmlex.TokenHashElse).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	var endCondGroupHandle func()
	endCondGroup := func() {
		if endCondGroupHandle != nil {
			endCondGroupHandle()
			endCondGroupHandle = nil
		}
	}
	defer endCondGroup()

	stack := make([]antlr.Tree, 0, 16)
	for _, nestedChild := range slices.Backward(node.GetChildren()) {
		stack = append(stack, nestedChild)
	}

	for len(stack) > 0 {
		lastIdx := len(stack) - 1
		rawChild := stack[lastIdx]
		stack = stack[:lastIdx]

		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenIf,
					dmlex.TokenHashIf,
					dmlex.TokenLparen,
					dmlex.TokenRparen,
					dmlex.TokenElse,
					dmlex.TokenHashElse:
					return true
				}
				return false
			}, "token is not supported")

			if tt == dmlex.TokenLparen {
				assert(endCondGroupHandle == nil, "ifCondition must be nil")
				endCondGroupHandle = p.LayoutEngine.BeginWrapGroup(
					"ifCondition",
				)
			}

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenRparen {
				endCondGroup()
			}

		case dmlex.IConditionExpressionContext:
			if err := p.printConditionExpression(child); err != nil {
				return err
			}

		case dmlex.IStatementContext:
			if innerIf := child.IfStatement(); innerIf != nil {
				for _, nestedChild := range slices.Backward(innerIf.GetChildren()) {
					stack = append(stack, nestedChild)
				}
			} else {
				if err := p.printStatement2(
					child, addCompoundIfNeeeded,
				); err != nil {
					return err
				}
			}

		case dmlex.IElseStatementContext:
			if err := p.printElseStatement(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: logLevel
// g4-tag: thenLevel
func (p *Printer) printLogLevel(child dmlex.ILogLevelContext) error {
	endGroup := p.LayoutEngine.BeginGreedyGroup("LogLevel")
	defer endGroup()
	endSoftIndent := p.LayoutEngine.BeginSoftIndent()
	defer endSoftIndent()

	for _, rawSubChild := range child.GetChildren() {
		switch subChild := rawSubChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				return dmlex.GetTokenType(
					subChild.GetSymbol(),
				) == dmlex.TokenThen
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(
				subChild,
			); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(subChild); err != nil {
				return err
			}
		case dmlex.IThenLevelContext:
			if subChild.Expression() != nil {
				if err := p.printExpression(
					subChild.Expression(),
				); err != nil {
					return err
				}
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", subChild))
		}
	}
	return nil
}

// g4-tag: logGroup
func (p *Printer) printLogGroup(child dmlex.ILogGroupContext) error {
	endGroup := p.LayoutEngine.BeginGreedyGroup("LogGroup")
	defer endGroup()
	endSoftIndent := p.LayoutEngine.BeginSoftIndent()
	defer endSoftIndent()

	if child.Expression() != nil {
		if err := p.printExpression(child.Expression()); err != nil {
			return err
		}
	}
	return nil
}

// g4-tag: logStatement
// g4-tag: logKind
// g4-tag: logFormatString
func (p *Printer) printLogStatement(node dmlex.ILogStatementContext) error {
	hasExprList := node.Expression_list_ntc_nonempty() != nil

	colonFlag := p.context.GetFormatFlag(dmlex.TokenColon).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorBeforePriority(format.PriorityHigh)

	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenColon: colonFlag,
		dmlex.TokenComma: p.context.GetFormatFlag(dmlex.TokenComma).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorPriority(format.PriorityHigh),
	}
	p.newContext(format.NewTokenMap(overrides))
	contextFreed := false
	defer func() {
		if !contextFreed {
			p.freeContext()
		}
	}()

	var endSoftIndent func()
	defer func() {
		if endSoftIndent != nil {
			endSoftIndent()
		}
	}()

	colonSeen := false
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenLog,
					dmlex.TokenColon,
					dmlex.TokenComma,
					dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")

			switch {
			case tt == dmlex.TokenComma && colonSeen && hasExprList:
				if err := p.pushTokenWithCustomFormat(
					child,
					func(ff format.FormatFlag) format.FormatFlag {
						return ff.WithSeparatorAfter(format.SepMandatoryWrap).
							WithSeparatorAfterPriority(format.PriorityHigh)
					},
				); err != nil {
					return err
				}
			case tt == dmlex.TokenSemi:
				if endSoftIndent != nil {
					endSoftIndent()
					endSoftIndent = nil
				}
				if err := p.pushTokenRangeFromRightiest(child); err != nil {
					return err
				}
			default:
				if err := p.pushTokenRangeFromRightiest(child); err != nil {
					return err
				}
			}

			if tt == dmlex.TokenColon {
				colonSeen = true
				if !contextFreed {
					p.freeContext()
					contextFreed = true
				}
				endSoftIndent = p.LayoutEngine.BeginSoftIndent()
			}
		case dmlex.ILogKindContext:
			for _, rawChild := range child.GetChildren() {
				switch term := rawChild.(type) {
				case antlr.TerminalNode:
					assertLazy(func() bool {
						switch dmlex.GetTokenType(term.GetSymbol()) {
						case dmlex.TokenInfo,
							dmlex.TokenError,
							dmlex.TokenSpecViol,
							dmlex.TokenWarning,
							dmlex.TokenUnimpl,
							dmlex.TokenCritical,
							dmlex.TokenIdentifier:
							return true
						}
						return false
					}, "token is not supported")
					if err := p.pushTokenRangeFromRightiest(term); err != nil {
						return err
					}
				case dmlex.IUtf8_sconstContext:
					if err := p.printUtf8_sconst(term); err != nil {
						return err
					}
				default:
					panic("unreachable: " + fmt.Sprintf("%T", term))
				}
			}
		case dmlex.ILogLevelContext:
			if err := p.printLogLevel(child); err != nil {
				return err
			}
		case dmlex.ILogGroupContext:
			if err := p.printLogGroup(child); err != nil {
				return err
			}
		case dmlex.ILogFormatStringContext:
			if child.Bracketed_string_literal() != nil {
				if err := p.printBracketed_string_literal(
					child.Bracketed_string_literal(),
				); err != nil {
					return err
				}
			}
		case dmlex.IExpression_list_ntc_nonemptyContext:
			if err := p.printExpression_listExpr(child, rootNode); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: selectStatement
func (p *Printer) printSelectStatement(
	node dmlex.ISelectStatementContext,
) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenElse: p.context.GetFormatFlag(dmlex.TokenElse).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal),
		dmlex.TokenHashElse: p.context.GetFormatFlag(dmlex.TokenHashElse).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	var endInExprGroupHandle func()
	endInExprGroup := func() {
		if endInExprGroupHandle != nil {
			endInExprGroupHandle()
			endInExprGroupHandle = nil
		}
	}
	defer endInExprGroup()

	var endWhereCondGroupHandle func()
	endWhereCondGroup := func() {
		if endWhereCondGroupHandle != nil {
			endWhereCondGroupHandle()
			endWhereCondGroupHandle = nil
		}
	}
	defer endWhereCondGroup()

	stack := make([]antlr.Tree, 0, 16)
	for _, nestedChild := range slices.Backward(node.GetChildren()) {
		stack = append(stack, nestedChild)
	}

	seenIn := false
	seenWhere := false

	for len(stack) > 0 {
		lastIdx := len(stack) - 1
		rawChild := stack[lastIdx]
		stack = stack[:lastIdx]

		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenSelect,
					dmlex.TokenHashSelect,
					dmlex.TokenIn,
					dmlex.TokenWhere,
					dmlex.TokenLparen,
					dmlex.TokenRparen,
					dmlex.TokenElse,
					dmlex.TokenHashElse:
					return true
				}
				return false
			}, "token is not supported")

			switch tt {
			case dmlex.TokenIn:
				seenIn = true
			case dmlex.TokenWhere:
				seenWhere = true
			case dmlex.TokenLparen:
				if seenWhere {
					assert(
						endWhereCondGroupHandle == nil,
						"selectWhereCondition must be nil",
					)
					endWhereCondGroupHandle = p.LayoutEngine.BeginWrapGroup(
						"selectCondition",
					)
				} else if seenIn {
					assert(
						endInExprGroupHandle == nil,
						"selectInExpression must be nil",
					)
					endInExprGroupHandle = p.LayoutEngine.BeginWrapGroup(
						"selectExpression",
					)
				}
			}

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenRparen {
				if endWhereCondGroupHandle != nil {
					endWhereCondGroup()
					seenWhere = false
				} else if endInExprGroupHandle != nil {
					endInExprGroup()
					seenIn = false
				}
			}

		case dmlex.IIdentContext:
			if err := p.printIdent(child); err != nil {
				return err
			}

		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}

		case dmlex.IConditionExpressionContext:
			if err := p.printConditionExpression(child); err != nil {
				return err
			}

		case dmlex.ILoopStatementContext:
			if err := p.printLoopStatement(child); err != nil {
				return err
			}

		case dmlex.IElseStatementContext:
			if err := p.printElseStatement(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: switchStatement
func (p *Printer) printSwitchStatement(
	node dmlex.ISwitchStatementContext,
) error {
	var endCondGroupHandle func()
	endCondGroup := func() {
		if endCondGroupHandle != nil {
			endCondGroupHandle()
			endCondGroupHandle = nil
		}
	}
	defer endCondGroup()

	var lastCaseHadStatements bool

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenSwitch,
					dmlex.TokenLparen,
					dmlex.TokenRparen,
					dmlex.TokenLbrace,
					dmlex.TokenRbrace:
					return true
				}
				return false
			}, "token is not supported")

			if tt == dmlex.TokenLparen {
				assert(endCondGroupHandle == nil, "switchCondition must be nil")
				endCondGroupHandle = p.LayoutEngine.BeginWrapGroup(
					"switchCondition",
				)
			}

			switch tt {
			case dmlex.TokenLbrace:
				if err := p.pushTokenWithCustomFormat(
					child,
					func(ff format.FormatFlag) format.FormatFlag {
						return ff.WithoutIndentIncrease().
							WithSeparatorAfter(format.SepMandatoryWrap)
					},
				); err != nil {
					return err
				}
			case dmlex.TokenRbrace:
				if err := p.pushTokenWithCustomFormat(
					child,
					func(ff format.FormatFlag) format.FormatFlag {
						ff = ff.WithSeparatorBefore(format.SepMandatoryWrap).
							WithSeparatorAfter(format.SepMandatoryWrap)
						if lastCaseHadStatements {
							return ff.WithIndentDecrease()
						}
						return ff.WithoutIndentDecrease()
					},
				); err != nil {
					return err
				}
			default:
				if err := p.pushTokenRangeFromRightiest(child); err != nil {
					return err
				}
			}

			if tt == dmlex.TokenRparen {
				endCondGroup()
			}

		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}

		case dmlex.ISwitchBodyStatementContext:
			var err error
			lastCaseHadStatements, err = p.printSwitchBodyStatement(child)
			if err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: switchBodyStatement
func (p *Printer) printSwitchBodyStatement(
	node dmlex.ISwitchBodyStatementContext,
) (lastCaseHadStatements bool, err error) {
	prevHadStatements := false
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(
				func() bool { return false },
				"no terminal nodes expected",
			)
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return false, err
			}
		case dmlex.ICaseStatementContext:
			prevHadStatements, err = p.printCaseStatement(
				child,
				prevHadStatements,
			)
			if err != nil {
				return false, err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return prevHadStatements, nil
}

// g4-tag: caseStatement
func (p *Printer) printCaseStatement(
	node dmlex.ICaseStatementContext,
	prevHadStatements bool,
) (hasStatements bool, err error) {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				return tt == dmlex.TokenColon
			}, "token is not supported")

			hasStmts := len(node.AllStatement()) > 0
			if err := p.pushTokenWithCustomFormat(
				child,
				func(ff format.FormatFlag) format.FormatFlag {
					ff = ff.WithSeparatorBefore(format.SepNone).
						WithSeparatorBeforePriority(format.PriorityHigh).
						WithSeparatorAfter(format.SepMandatoryWrap).
						WithSeparatorAfterPriority(format.PriorityHigh)
					if hasStmts {
						return ff.WithIndentIncrease()
					}
					return ff.WithoutIndentIncrease()
				},
			); err != nil {
				return false, err
			}

		case dmlex.ICaseOrDefaultStatementContext:
			if err := p.printCaseOrDefaultStatement(
				child,
				prevHadStatements,
			); err != nil {
				return false, err
			}

		case dmlex.IStatementContext:
			if err := p.printStatement(child); err != nil {
				return false, err
			}

		case dmlex.ISwitchBodyHashIfStatementContext:
			hasTrailingStmts := len(node.AllStatement()) > 0
			if _, err := p.printSwitchBodyHashIfStatement(
				child,
				prevHadStatements,
				hasTrailingStmts,
			); err != nil {
				return false, err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return len(node.AllStatement()) > 0, nil
}

// g4-tag: caseOrDefaultStatement
func (p *Printer) printCaseOrDefaultStatement(
	node dmlex.ICaseOrDefaultStatementContext,
	prevHadStatements bool,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenCase, dmlex.TokenDefault:
					return true
				}
				return false
			}, "token is not supported")

			if err := p.pushTokenWithCustomFormat(
				child,
				func(ff format.FormatFlag) format.FormatFlag {
					ff = ff.WithSeparatorBefore(format.SepMandatoryWrap)
					if prevHadStatements {
						return ff.WithIndentDecrease()
					}
					return ff.WithoutIndentDecrease()
				},
			); err != nil {
				return err
			}

		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: switchBodyHashIfStatement
func (p *Printer) printSwitchBodyHashIfStatement(
	node dmlex.ISwitchBodyHashIfStatementContext,
	prevHadStatements bool,
	hasTrailingStatements bool,
) (hasStatements bool, err error) {
	var endCondGroupHandle func()
	endCondGroup := func() {
		if endCondGroupHandle != nil {
			endCondGroupHandle()
			endCondGroupHandle = nil
		}
	}
	defer endCondGroup()

	currentHadStatements := prevHadStatements

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenHashIf,
					dmlex.TokenLparen,
					dmlex.TokenRparen,
					dmlex.TokenLbrace,
					dmlex.TokenRbrace,
					dmlex.TokenHashElse:
					return true
				}
				return false
			}, "token is not supported")

			switch tt {
			case dmlex.TokenHashIf:
				if err := p.pushTokenWithCustomFormat(
					child,
					func(ff format.FormatFlag) format.FormatFlag {
						ff = ff.WithSeparatorBefore(format.SepMandatoryWrap)
						if currentHadStatements {
							return ff.WithIndentDecrease()
						}
						return ff.WithoutIndentDecrease()
					},
				); err != nil {
					return false, err
				}
				currentHadStatements = false

			case dmlex.TokenLparen:
				assert(endCondGroupHandle == nil, "hashIfCondition must be nil")
				endCondGroupHandle = p.LayoutEngine.BeginWrapGroup(
					"hashIfCondition",
				)
				if err := p.pushTokenRangeFromRightiest(child); err != nil {
					return false, err
				}

			case dmlex.TokenRparen:
				if err := p.pushTokenRangeFromRightiest(child); err != nil {
					return false, err
				}
				endCondGroup()

			case dmlex.TokenLbrace:
				if err := p.pushTokenWithCustomFormat(
					child,
					func(ff format.FormatFlag) format.FormatFlag {
						return ff.WithoutIndentIncrease().
							WithSeparatorAfter(format.SepMandatoryWrap)
					},
				); err != nil {
					return false, err
				}

			case dmlex.TokenRbrace:
				rbraces := node.AllRBRACE()
				isLastRbrace := len(rbraces) > 0 &&
					child == rbraces[len(rbraces)-1]
				if err := p.pushTokenWithCustomFormat(
					child,
					func(ff format.FormatFlag) format.FormatFlag {
						ff = ff.WithSeparatorBefore(format.SepMandatoryWrap).
							WithSeparatorAfter(format.SepMandatoryWrap)
						if currentHadStatements {
							ff = ff.WithIndentDecrease()
						} else {
							ff = ff.WithoutIndentDecrease()
						}
						if isLastRbrace && hasTrailingStatements {
							ff = ff.WithIndentIncrease()
						}
						return ff
					},
				); err != nil {
					return false, err
				}
				currentHadStatements = false

			case dmlex.TokenHashElse:
				if err := p.pushTokenWithCustomFormat(
					child,
					func(ff format.FormatFlag) format.FormatFlag {
						ff = ff.WithSeparatorBefore(format.SepSpace).
							WithSeparatorBeforePriority(format.PriorityAboveNormal)
						if currentHadStatements {
							return ff.WithIndentDecrease()
						}
						return ff.WithoutIndentDecrease()
					},
				); err != nil {
					return false, err
				}
				currentHadStatements = false
			}

		case dmlex.IConditionExpressionContext:
			if err := p.printConditionExpression(child); err != nil {
				return false, err
			}

		case dmlex.ISwitchBodyStatementContext:
			currentHadStatements, err = p.printSwitchBodyStatement(child)
			if err != nil {
				return false, err
			}

		case dmlex.ISwitchBodyHashElseStatementContext:
			currentHadStatements, err = p.printSwitchBodyHashElseStatement(
				child,
			)
			if err != nil {
				return false, err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return currentHadStatements, nil
}

// g4-tag: switchBodyHashElseStatement
func (p *Printer) printSwitchBodyHashElseStatement(
	node dmlex.ISwitchBodyHashElseStatementContext,
) (bool, error) {
	if node.SwitchBodyStatement() != nil {
		return p.printSwitchBodyStatement(node.SwitchBodyStatement())
	}
	return false, nil
}

// g4-tag: throwStatement
func (p *Printer) printThrowStatement(node dmlex.IThrowStatementContext) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenThrow, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: tryStatement
// g4-tag: catchStatement
func (p *Printer) printTryStatement(node dmlex.ITryStatementContext) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenCatch: p.context.GetFormatFlag(dmlex.TokenCatch).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityAboveNormal),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	stack := make([]antlr.Tree, 0, 8)
	for _, nestedChild := range slices.Backward(node.GetChildren()) {
		stack = append(stack, nestedChild)
	}

	for len(stack) > 0 {
		lastIdx := len(stack) - 1
		rawChild := stack[lastIdx]
		stack = stack[:lastIdx]

		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenTry,
					dmlex.TokenCatch:
					return true
				}
				return false
			}, "token is not supported")

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IStatementContext:
			if err := p.printStatement2(
				child, addCompoundIfNeeeded,
			); err != nil {
				return err
			}

		case dmlex.ICatchStatementContext:
			for _, nestedChild := range slices.Backward(child.GetChildren()) {
				stack = append(stack, nestedChild)
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: whileLoopStatement
func (p *Printer) printWhileLoopStatement(node antlr.Tree) error {
	var endCondGroupHandle func()
	endCondGroup := func() {
		if endCondGroupHandle != nil {
			endCondGroupHandle()
			endCondGroupHandle = nil
		}
	}
	defer endCondGroup()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenDo,
					dmlex.TokenWhile,
					dmlex.TokenLparen,
					dmlex.TokenRparen,
					dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")

			if tt == dmlex.TokenLparen {
				assert(endCondGroupHandle == nil, "whileCondition must be nil")
				endCondGroupHandle = p.LayoutEngine.BeginWrapGroup(
					"whileCondition",
				)
			}

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenRparen {
				endCondGroup()
			}

		case dmlex.IConditionExpressionContext:
			if err := p.printConditionExpression(child); err != nil {
				return err
			}

		case dmlex.ILoopStatementContext:
			if err := p.printLoopStatement(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: deviceStatement
func (p *Printer) printDeviceStatement(
	node dmlex.IDeviceStatementContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(
				func() bool { return false },
				"no terminal nodes expected",
			)
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ITemplateStatementContext:
			if err := p.printTemplateStatement(child); err != nil {
				return err
			}
		case dmlex.IParamContext:
			if err := p.printParamStatement(child); err != nil {
				return err
			}
		case dmlex.ICCodeBlockContext:
			if err := p.printCCodeBlock(child); err != nil {
				return err
			}
		case dmlex.ILogGroupStatementContext:
			if err := p.printLogGroupStatement(child); err != nil {
				return err
			}
		case dmlex.IConstantStatementContext:
			if err := p.printConstantStatement(child); err != nil {
				return err
			}
		case dmlex.ITypedefStatementContext:
			if err := p.printTypedefStatement(child); err != nil {
				return err
			}
		case dmlex.IExternStatementContext:
			if err := p.printExternStatement(child); err != nil {
				return err
			}
		case dmlex.IImportStatementContext:
			if err := p.printImportStatement(child); err != nil {
				return err
			}
		case dmlex.IExportStatementContext:
			if err := p.printExportStatement(child); err != nil {
				return err
			}
		case dmlex.IObjectContext:
			if err := p.printObject(child); err != nil {
				return err
			}
		case dmlex.IMethodUnifiedContext:
			if err := p.printMethodUnified(child); err != nil {
				return err
			}
		case dmlex.IIstemplateStatementContext:
			if err := p.printIsTemplateStatement(child); err != nil {
				return err
			}
		case dmlex.IToplevelIfContext:
			if err := p.printToplevelIf(child); err != nil {
				return err
			}
		case dmlex.IErrorStatementContext:
			if err := p.printErrorOrWarningStatement(child); err != nil {
				return err
			}
		case dmlex.IIneachStatementContext:
			if err := p.printIneachStatement(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return p.Flush()
}

// g4-tag: objectStatement
func (p *Printer) printObjectStatement(
	node dmlex.IObjectStatementContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(
				func() bool { return false },
				"no terminal nodes expected",
			)
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IObjectContext:
			if err := p.printObject(child); err != nil {
				return err
			}
		case dmlex.IParamContext:
			if err := p.printParamStatement(child); err != nil {
				return err
			}
		case dmlex.IMethodUnifiedContext:
			if err := p.printMethodUnified(child); err != nil {
				return err
			}
		case dmlex.IIstemplateStatementContext:
			if err := p.printIsTemplateStatement(child); err != nil {
				return err
			}
		case dmlex.IObjectIfContext:
			if err := p.printObjectIf(child); err != nil {
				return err
			}
		case dmlex.IErrorStatementContext:
			if err := p.printErrorOrWarningStatement(child); err != nil {
				return err
			}
		case dmlex.IIneachStatementContext:
			if err := p.printIneachStatement(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return p.writeNewLines(1)
}

// g4-tag: logGroupStatement
func (p *Printer) printLogGroupStatement(
	node dmlex.ILogGroupStatementContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenLoggroup, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IIdentContext:
			if err := p.printIdent(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return p.Flush()
}

// g4-tag: externStatement
func (p *Printer) printExternStatement(
	node dmlex.IExternStatementContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenTypedef, dmlex.TokenSemi, dmlex.TokenExtern:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.ICdeclContext:
			if err := p.printCdecl(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: importStatement
func (p *Printer) printImportStatement(
	node dmlex.IImportStatementContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenImport,
					dmlex.TokenStringLiteral,
					dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return p.Flush()
}

// g4-tag: exportStatement
func (p *Printer) printExportStatement(
	node dmlex.IExportStatementContext,
) error {
	var endSoftIndentHandle func()
	endSoftIndent := func() {
		if endSoftIndentHandle != nil {
			endSoftIndentHandle()
			endSoftIndentHandle = nil
		}
	}
	defer endSoftIndent()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenExport, dmlex.TokenAs, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if tt == dmlex.TokenAs {
				assert(endSoftIndentHandle == nil, "unexpected token")
				endSoftIndentHandle = p.LayoutEngine.BeginSoftIndent()
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if tt == dmlex.TokenSemi && endSoftIndentHandle != nil {
				endSoftIndent()
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		case dmlex.IExportAsExpressionContext:
			if err := p.printExportAsExpression(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return p.Flush()
}

func (p *Printer) printExportAsExpression(
	node dmlex.IExportAsExpressionContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(
				func() bool { return false },
				"no terminal nodes expected",
			)
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: toplevelIf
//
//go:inline
func (p *Printer) printToplevelIf(node dmlex.IToplevelIfContext) error {
	return p.printObjectIfStatement(node)
}

// g4-tag: ineachStatement
func (p *Printer) printIneachStatement(
	node dmlex.IIneachStatementContext,
) error {
	hasStatements := len(node.AllObjectStatement()) > 0
	if !hasStatements && node.LBRACE() != nil && node.RBRACE() != nil {
		lbraceIdx := getTokenIndexFromToken(node.LBRACE().GetSymbol())
		rbraceIdx := getTokenIndexFromToken(node.RBRACE().GetSymbol())
		for i := lbraceIdx + 1; i < rbraceIdx; i++ {
			if isSignificantToken(p.extractToken(i)) {
				hasStatements = true
				break
			}
		}
	}

	defaultCtx := p.settings.TokenFormatDefaults()
	lflag := defaultCtx.GetFormatFlag(dmlex.TokenLbrace).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorPriority(format.PriorityMax)
	rflag := defaultCtx.GetFormatFlag(dmlex.TokenRbrace).
		WithSeparatorPriority(format.PriorityMax)

	if hasStatements {
		lflag = lflag.WithSeparatorAfter(format.SepMandatoryWrap).
			WithIndentIncrease()
		rflag = rflag.WithSeparatorBefore(format.SepMandatoryWrap).
			WithSeparatorAfter(format.SepNone).WithIndentDecrease()
	} else {
		lflag = lflag.WithSeparatorAfter(format.SepNone)
		rflag = rflag.WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepSpace)
	}

	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenLbrace: lflag,
		dmlex.TokenRbrace: rflag,
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenIn,
					dmlex.TokenEach,
					dmlex.TokenLbrace,
					dmlex.TokenRbrace:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if tt == dmlex.TokenLbrace && hasStatements {
				if err := p.writeNewLines(1); err != nil {
					return err
				}
			}
		case dmlex.IIstemplateListContext:
			if err := p.printIsTemplateList(child); err != nil {
				return err
			}
		case dmlex.IObjectStatementContext:
			if err := p.printObjectStatement(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: object
func (p *Printer) printObject(node dmlex.IObjectContext) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(
				func() bool { return false },
				"no terminal nodes expected",
			)
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IObject_declContext:
			if err := p.printObject_decl(child); err != nil {
				return err
			}
		case dmlex.ISessionOrSavedStatementContext:
			if err := p.printSessionOrSavedStatement(child); err != nil {
				return err
			}
		case dmlex.IHook_declContext:
			if err := p.printHook_decl(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: returnStatement
func (p *Printer) printReturnStatement(
	node dmlex.IReturnStatementContext,
) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenReturn: p.context.GetFormatFlag(dmlex.TokenReturn).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityHigh),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	var endSoftIndent func()
	defer func() {
		if endSoftIndent != nil {
			endSoftIndent()
		}
	}()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenReturn, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if tt == dmlex.TokenSemi && endSoftIndent != nil {
				endSoftIndent()
				endSoftIndent = nil
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if tt == dmlex.TokenReturn && node.Initializer() != nil {
				endSoftIndent = p.LayoutEngine.BeginSoftIndent()
			}
		case dmlex.IInitializerContext:
			if err := p.printInitializer(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: deleteStatement
func (p *Printer) printDeleteStatement(
	node dmlex.IDeleteStatementContext,
) error {
	var endSoftIndentHandle func()
	endSoftIndent := func() {
		if endSoftIndentHandle != nil {
			endSoftIndentHandle()
			endSoftIndentHandle = nil
		}
	}
	defer endSoftIndent()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenDelete, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if tt == dmlex.TokenSemi {
				endSoftIndent()
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if tt == dmlex.TokenDelete && node.Expression() != nil {
				assert(endSoftIndentHandle == nil, "unexpeted")
				endSoftIndentHandle = p.LayoutEngine.BeginSoftIndent()
			}
		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: assertStatement
func (p *Printer) printAssertStatement(
	node dmlex.IAssertStatementContext,
) error {
	var endSoftIndentHandle func()
	endSoftIndent := func() {
		if endSoftIndentHandle != nil {
			endSoftIndentHandle()
			endSoftIndentHandle = nil
		}
	}
	defer endSoftIndent()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenAssert, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if tt == dmlex.TokenSemi {
				endSoftIndent()
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if tt == dmlex.TokenAssert && node.ConditionExpression() != nil {
				assert(endSoftIndentHandle == nil, "unexpected")
				endSoftIndentHandle = p.LayoutEngine.BeginSoftIndent()
			}
		case dmlex.IConditionExpressionContext:
			if err := p.printConditionExpression(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: errorStatement
// g4-tag: warningStatement
func (p *Printer) printErrorOrWarningStatement(
	node antlr.Tree,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenError, dmlex.TokenWarning, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IBracketed_string_literalContext:
			if err := p.printBracketed_string_literal(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: hook_decl
func (p *Printer) printHook_decl(node dmlex.IHook_declContext) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{ // TODO: Try to remove override
		dmlex.TokenLparen: DefaultFlagFor(dmlex.TokenLparen).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorPriority(format.PriorityVeryHigh),
		// Lowered from PriorityMax to allow comments to win

		dmlex.TokenRparen: DefaultFlagFor(dmlex.TokenRparen).
			WithSeparatorPriority(format.PriorityVeryHigh),
		// Lowered from PriorityMax to allow comments to win

		dmlex.TokenLbrack: DefaultFlagFor(dmlex.TokenLbrack).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(format.PriorityAboveNormal).
			WithSeparatorPriority(format.PriorityVeryHigh),
		// Upgraded to beat default tokens, but still lose to comments

		dmlex.TokenRbrack: DefaultFlagFor(dmlex.TokenRbrack).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityAboveNormal).
			WithSeparatorPriority(format.PriorityVeryHigh),
	}

	// For empty params, force no-space between parentheses to prevent wrapping
	if node.Cdecl_list() == nil { // TODO: Workaround
		overrides[dmlex.TokenLparen] = overrides[dmlex.TokenLparen].WithSeparatorAfter(
			format.SepNone,
		)
		overrides[dmlex.TokenRparen] = overrides[dmlex.TokenRparen].WithSeparatorBefore(
			format.SepNone,
		)
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	var endHookArgsGroupHandle func()
	endHookArgsGroup := func() {
		if endHookArgsGroupHandle != nil {
			endHookArgsGroupHandle()
			endHookArgsGroupHandle = nil
		}
	}
	defer endHookArgsGroup()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			killArgsWrap := false
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenHook,
					dmlex.TokenLparen,
					dmlex.TokenRparen,
					dmlex.TokenLbrack,
					dmlex.TokenRbrack,
					dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")

			switch tt {
			case dmlex.TokenLparen:
				if node.Cdecl_list() != nil {
					assert(
						endHookArgsGroupHandle == nil,
						"hookArgs must be nil",
					)
					endHookArgsGroupHandle = p.LayoutEngine.BeginWrapGroup(
						"hookArgs",
					)
				}
			case dmlex.TokenRparen:
				killArgsWrap = true
			}

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if killArgsWrap && endHookArgsGroupHandle != nil {
				endHookArgsGroup()
			}

		case dmlex.ICdecl_listContext:
			if err := p.printCdecl_list(child); err != nil {
				return err
			}

		case dmlex.IIdentContext:
			if err := p.printIdent(child); err != nil {
				return err
			}

		case dmlex.IExpressionContext:
			if err := p.printExpression(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return p.Flush()
}
