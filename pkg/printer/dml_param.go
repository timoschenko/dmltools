// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"

	"github.com/antlr4-go/antlr/v4"
)

// g4-tag: param
func (p *Printer) printParamStatement(node dmlex.IParamContext) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenSemi: format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorPriority(format.PriorityMax),
		dmlex.TokenColon: format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityLowest).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorPriority(format.PriorityAboveNormal),
		dmlex.TokenAssign: format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityMax).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityMediumHigh),
		dmlex.TokenNewAssign: format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityMax).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityMediumHigh),
		dmlex.TokenDefault: format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityMax).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityMediumHigh),
		dmlex.TokenDefaultNew: format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityMax).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorAfterPriority(format.PriorityMediumHigh),
		dmlex.TokenAuto: format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityLowest).
			WithSeparatorPriority(format.PriorityAboveNormal),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	var endSoftIdent func()
	defer func() {
		if endSoftIdent != nil {
			endSoftIdent()
		}
	}()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenParam,
					dmlex.TokenColon,
					dmlex.TokenAssign,
					dmlex.TokenDefault,
					dmlex.TokenNewAssign,
					dmlex.TokenDefaultNew,
					dmlex.TokenAuto,
					dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenParam && endSoftIdent == nil {
				endSoftIdent = p.LayoutEngine.BeginSoftIndent()
			}

			if tt == dmlex.TokenSemi && endSoftIdent != nil {
				endSoftIdent()
				endSoftIdent = nil
			}

		case dmlex.IObjidentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.ICtypedeclContext:
			if err := p.printCtypedecl(child); err != nil {
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

// g4-tag: constantStatement
func (p *Printer) printConstantStatement(
	node dmlex.IConstantStatementContext,
) error {
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenConstant: format.FormatFlag(0).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorPriority(format.PriorityHigh),
		dmlex.TokenSemi: format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorPriority(format.PriorityMax),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	var endSoftIdent func()
	defer func() {
		if endSoftIdent != nil {
			endSoftIdent()
		}
	}()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenConstant, dmlex.TokenAssign, dmlex.TokenSemi:
					return true
				}
				return false
			}, "token is not supported")

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenAssign && endSoftIdent == nil {
				endSoftIdent = p.LayoutEngine.BeginSoftIndent()
			}

			if tt == dmlex.TokenSemi && endSoftIdent != nil {
				endSoftIdent()
				endSoftIdent = nil
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
