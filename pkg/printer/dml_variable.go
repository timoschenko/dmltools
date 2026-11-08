// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"

	"github.com/antlr4-go/antlr/v4"
)

// g4-tag: cvarDecl
func (p *Printer) printCvarDecl(
	node dmlex.ICvarDeclContext,
) error {
	p.newContext(p.settings.CdeclDefaults())
	defer p.freeContext()

	// FIXME: Deprecased, process all overrides dynamicly within the for-each loop
	if node.LPAREN() != nil {
		parenOverrides := map[dmlex.TokenType]format.FormatFlag{
			dmlex.TokenLparen: p.context.GetFormatFlag(dmlex.TokenLparen).
				WithSeparatorBefore(format.SepSpace).
				WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
				WithBreakAfterPriority(format.PriorityAboveNormal).
				WithIndentIncrease().
				WithSeparatorAfterPriority(format.PriorityMax),
			dmlex.TokenRparen: p.context.GetFormatFlag(dmlex.TokenRparen).
				WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
				WithBreakBeforePriority(format.PriorityAboveNormal).
				WithIndentDecrease().
				WithSeparatorBeforePriority(format.PriorityMax),
		}
		p.newContext(format.NewTokenMap(parenOverrides))
		defer p.freeContext()
	}

	var endGroup func()
	if node.LPAREN() != nil {
		endGroup = p.LayoutEngine.BeginWrapGroup("cdecl_list")
	}
	defer func() {
		if endGroup != nil {
			endGroup()
		}
	}()

	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenLparen, dmlex.TokenRparen, dmlex.TokenAssign:
					return true
				}
				return false
			}, "token is not supported")

			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

			if tt == dmlex.TokenRparen && endGroup != nil {
				endGroup()
				endGroup = nil
			}

		case dmlex.ICdeclContext:
			if err := p.printCdecl(child); err != nil {
				return err
			}

		case dmlex.ICdecl_listContext:
			if err := p.printCdecl_list(child); err != nil {
				return err
			}

		case dmlex.IInitializerContext:
			if err := p.printInitializer(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	stopIdx := p.gatherTokenIndexFromRightiestToken(node.GetStop())
	return p.pushTokenRange(stopIdx)
}

// g4-tag: sessionOrSavedDecl
func (p *Printer) printSessionOrSavedDecl(
	node dmlex.ISessionOrSavedDeclContext,
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
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IGlobal_decl_kindContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			assert(endSoftIndentHandle == nil, "unexpected")
			endSoftIndentHandle = p.LayoutEngine.BeginSoftIndent()

		case dmlex.ICvarDeclContext:
			if err := p.printCvarDecl(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: variableDecl
func (p *Printer) printVariableDecl(
	node dmlex.IVariableDeclContext,
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
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.ILocal_decl_kindContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			assert(endSoftIndentHandle == nil, "unexpected")
			endSoftIndentHandle = p.LayoutEngine.BeginSoftIndent()

		case dmlex.ICvarDeclContext:
			if err := p.printCvarDecl(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: variableStatement
func (p *Printer) printVariableStatement(
	node dmlex.IVariableStatementContext,
) error {
	p.newContext(p.settings.CdeclDefaults())
	defer p.freeContext()

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

		case dmlex.IVariableDeclContext:
			if err := p.printVariableDecl(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: sessionOrSavedStatement
func (p *Printer) printSessionOrSavedStatement(
	node dmlex.ISessionOrSavedStatementContext,
) error {
	p.newContext(p.settings.CdeclDefaults())
	defer p.freeContext()

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

		case dmlex.ISessionOrSavedDeclContext:
			if err := p.printSessionOrSavedDecl(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}
