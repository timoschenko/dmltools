// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

// g4-tag: bracketed_string_literal
func (p *Printer) printBracketed_string_literal(
	node dmlex.IBracketed_string_literalContext,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("BracketedStringLiteral")
	defer endGroup()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenLparen, dmlex.TokenRparen:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IComposed_string_literalContext:
			if err := p.printComposed_string_literal(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: composed_string_literal
func (p *Printer) printComposed_string_literal(
	node dmlex.IComposed_string_literalContext,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("ComposedStringLiteral")
	defer endGroup()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				return dmlex.GetTokenType(child.GetSymbol()) == dmlex.TokenAdd
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IUtf8_sconstContext:
			if err := p.printUtf8_sconst(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: utf8_sconst
func (p *Printer) printUtf8_sconst(
	node dmlex.IUtf8_sconstContext,
) error {
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				return dmlex.GetTokenType(
					child.GetSymbol(),
				) == dmlex.TokenStringLiteral
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
