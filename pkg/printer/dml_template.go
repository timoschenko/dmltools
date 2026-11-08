// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"

	"github.com/antlr4-go/antlr/v4"
)

// g4-tag: templateStatement
func (p *Printer) printTemplateStatement(
	node dmlex.ITemplateStatementContext,
) error {
	bodyStatements := node.AllTemplateBody()
	hasStatements := len(bodyStatements) > 0
	// FIXME: Deprecased, process all overrides dynamicly within the for-each loop
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

	// FIXME: Deprecased, process all overrides dynamicly within the for-each loop
	defaultCtx := p.settings.TokenFormatDefaults()
	lflag := defaultCtx.GetFormatFlag(dmlex.TokenLbrace).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorPriority(format.PriorityMax)
	rflag := defaultCtx.GetFormatFlag(dmlex.TokenRbrace).
		WithSeparatorPriority(format.PriorityMax)

	// FIXME: Deprecased, process all overrides dynamicly within the for-each loop
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

	var prevKind StatementKind
	var prevStopToken antlr.Token
	hadStatement := false

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenTemplate, dmlex.TokenLbrace, dmlex.TokenRbrace:
					return true
				}
				return false
			}, "token is not supported")

			if tt == dmlex.TokenRbrace && hasStatements {
				// Ensure gap tokens (comments) are pushed with current indentation
				// before we decrease it for the closing brace.
				rbraceStart := TokenIndex(child.GetSymbol().GetTokenIndex())
				if err := p.pushTokenRange(rbraceStart - 1); err != nil {
					return err
				}
			}

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenLbrace && hasStatements {
				if err := p.writeNewLines(1); err != nil {
					return err
				}
			}

		case dmlex.IObjidentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IIstemplateExprContext:
			if err := p.printIsTemplateExpr(child); err != nil {
				return err
			}

		case dmlex.ITemplateBodyContext:
			currKind := categorizeTemplateBody(child)
			startIdx := child.GetStart().GetTokenIndex()

			var nextStartTok antlr.Token
			if p.lastProcessedTokenIndex < TokenIndex(startIdx) {
				nextStartTok = p.extractToken(p.lastProcessedTokenIndex)
			} else {
				nextStartTok = child.GetStart()
			}

			if hadStatement {
				blankLines := countInputBlankLines(prevStopToken, nextStartTok)
				emptyLines := computeNestedSpacing(
					prevKind,
					currKind,
					blankLines,
				)
				if err := p.writeNewLines(emptyLines + 1); err != nil {
					return err
				}
			}

			if err := p.pushTokenRange(TokenIndex(startIdx - 1)); err != nil {
				return err
			}
			if err := p.printTemplateBody(child); err != nil {
				return err
			}

			hadStatement = true
			prevKind = currKind
			prevStopIdx := p.gatherTokenIndexFromRightiestToken(child.GetStop())
			prevStopToken = p.extractToken(prevStopIdx)

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: istemplateExpr
func (p *Printer) printIsTemplateExpr(
	node dmlex.IIstemplateExprContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				return tt == dmlex.TokenIs
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IIstemplateListContext:
			if err := p.printIsTemplateList(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: istemplateList
func (p *Printer) printIsTemplateList(node dmlex.IIstemplateListContext) error {
	if node == nil || node.GetChildCount() == 0 {
		return nil
	}

	startToken := castStartToken(node)
	tokenIndexBegin := getTokenIndexFromToken(startToken)

	p.newContext(p.settings.ExpressionDefaults())
	defer p.freeContext()

	overrides := make(format.TokenMapImpl)

	// FIXME: Deprecased, process all overrides dynamicly within the for-each loop
	// Ensure a space before the first token (whether it's LPAREN or an identifier)
	firstTok := p.extractToken(tokenIndexBegin)
	firstTT := dmlex.GetTokenType(firstTok)
	overrides[firstTT] = p.context.GetFormatFlag(firstTT).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorBeforePriority(format.PriorityMax)
	// make sure it overrides any SepNone

	// FIXME: Deprecased, process all overrides dynamicly within the for-each loop
	if node.LPAREN() != nil {
		overrides[dmlex.TokenLparen] = p.context.GetFormatFlag(dmlex.TokenLparen).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithBreakAfterPriority(format.PriorityAboveNormal).

			// 8
			WithIndentIncrease().
			WithSeparatorBeforePriority(format.PriorityMax).
			WithSeparatorAfterPriority(format.PriorityMax)

		overrides[dmlex.TokenRparen] = p.context.GetFormatFlag(dmlex.TokenRparen).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityAboveNormal).

			// 8
			WithIndentDecrease().
			WithSeparatorBeforePriority(format.PriorityMax)
	}

	overrides[dmlex.TokenComma] = p.context.GetFormatFlag(dmlex.TokenComma).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(format.PriorityDefault). // 7
		WithSeparatorBeforePriority(format.PriorityMax).
		WithSeparatorAfterPriority(format.PriorityMax)

	p.newContext(overrides)
	defer p.freeContext()

	var endGroup func()
	if node.LPAREN() != nil {
		endGroup = p.LayoutEngine.BeginWrapGroup("templateArgs")
		defer endGroup()
	}

	var endObjGroup func()
	defer func() {
		if endObjGroup != nil {
			endObjGroup()
		}
	}()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenLparen, dmlex.TokenComma, dmlex.TokenRparen:
					return true
				}
				return false
			}, "token is not supported")

			if tt == dmlex.TokenRparen && endObjGroup != nil {
				endObjGroup()
				endObjGroup = nil
			}

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenLparen && endObjGroup == nil {
				endObjGroup = p.LayoutEngine.BeginWrapGroup("templateArgsObjs")
			}

		case dmlex.IObjidentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: templateBody
func (p *Printer) printTemplateBody(node dmlex.ITemplateBodyContext) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				return tt == dmlex.TokenShared
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IHook_declContext:
			if err := p.printHook_decl(child); err != nil {
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

// g4-tag: isTemplateStatement
func (p *Printer) printIsTemplateStatement(
	node dmlex.IIstemplateStatementContext,
) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				return tt == dmlex.TokenSemi
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IIstemplateExprContext:
			if err := p.printIsTemplateExpr(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return p.Flush()
}
