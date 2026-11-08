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

// g4-tag: object_decl
func (p *Printer) printObject_decl(node dmlex.IObject_declContext) error {
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenRegister,
					dmlex.TokenField,
					dmlex.TokenSize,
					dmlex.TokenIn:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IObject_typeContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IObjidentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IArray_listContext:
			if err := p.printArray_list(child); err != nil {
				return err
			}

		case dmlex.ISizeSpecContext:
			if err := p.printSizeSpec(child); err != nil {
				return err
			}

		case dmlex.IOffsetSpecContext:
			if err := p.printOffsetSpec(child); err != nil {
				return err
			}

		case dmlex.IBitrangeDeclContext:
			if err := p.printBitrangeDecl(child); err != nil {
				return err
			}

		case dmlex.IObjectBodyContext:
			if err := p.printObjectBody(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}

// g4-tag: array_list
func (p *Printer) printArray_list(node dmlex.IArray_listContext) error {
	p.newContext(p.settings.ExpressionDefaults())
	defer p.freeContext()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenLbrack, dmlex.TokenRbrack:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IArraydefContext:
			if err := p.printArraydef(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: arraydef
func (p *Printer) printArraydef(node dmlex.IArraydefContext) error {
	p.newContext(p.settings.ExpressionDefaults())
	defer p.freeContext()

	// NOTE: Forbit line wrapping for `<`
	//       `array[X < len]` -> `array[\n\tX < len\n]` {{{
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenLt: format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorPriority(format.PriorityMediumHigh),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()
	// }}}

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenLt, dmlex.TokenEllipsis:
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

// g4-tag: ident
func (p *Printer) printIdent(node dmlex.IIdentContext) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: sizeSpec
func (p *Printer) printSizeSpec(node dmlex.ISizeSpecContext) error {
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

// g4-tag: offsetSpec
func (p *Printer) printOffsetSpec(node dmlex.IOffsetSpecContext) error {
	for _, rawChild := range node.GetChildren() {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				return tt == dmlex.TokenAt
			}, "token is not supported")
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

// g4-tag: bitrangeDecl
func (p *Printer) printBitrangeDecl(node dmlex.IBitrangeDeclContext) error {
	// NOTE: Ensure `arr[12:34]` but not `arr[ 12 : 34 ]` {{{
	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenLbrack: DefaultFlagFor(dmlex.TokenLbrack).
			WithSeparatorAfter(format.SepNone),
		dmlex.TokenColon: DefaultFlagFor(dmlex.TokenColon).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepNone),
		dmlex.TokenRbrack: DefaultFlagFor(dmlex.TokenRbrack).
			WithSeparatorBefore(format.SepNone),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()
	// }}}

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenAt,
					dmlex.TokenLbrack,
					dmlex.TokenColon,
					dmlex.TokenRbrack:
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

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: objectBody
func (p *Printer) printObjectBody(node dmlex.IObjectBodyContext) error {
	hadStatement := false
	var prevKind StatementKind
	var prevStopToken antlr.Token

	for _, rawChild := range node.GetChildren() {
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
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IIstemplateExprContext:
			if err := p.printIsTemplateExpr(child); err != nil {
				return err
			}

		case dmlex.IObjectDescSpecContext:
			if err := p.printObjectDescSpec(child); err != nil {
				return err
			}

		case dmlex.IObjectStatementContext:
			currKind := categorizeObjectStatement(child)
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

			if err := p.printObjectStatement(child); err != nil {
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
	return p.Flush()
}

// g4-tag: objectDescSpec
func (p *Printer) printObjectDescSpec(node dmlex.IObjectDescSpecContext) error {
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

// g4-tag: objectIf
//
//go:inline
func (p *Printer) printObjectIf(node dmlex.IObjectIfContext) error {
	return p.printObjectIfStatement(node)
}

// g4-tag: conditionExpression
func (p *Printer) printConditionExpression(
	node dmlex.IConditionExpressionContext,
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
// g4-tag: toplevelIfStatement
// g4-tag: toplevelElseStatement
// g4-tag: objectIfStatement
// g4-tag: objectElseStatement
func (p *Printer) printObjectIfStatement(node antlr.Tree) error {
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
				case dmlex.TokenHashIf,
					dmlex.TokenLparen,
					dmlex.TokenRparen,
					dmlex.TokenHashElse,
					dmlex.TokenLbrace,
					dmlex.TokenRbrace:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}

		case dmlex.IConditionExpressionContext:
			if err := p.printConditionExpression(child); err != nil {
				return err
			}

		case dmlex.IObjectIfContext,
			dmlex.IObjectIfStatementContext,
			dmlex.IObjectElseStatementContext,
			dmlex.IToplevelIfContext,
			dmlex.IToplevelIfStatementContext,
			dmlex.IToplevelElseStatementContext,
			dmlex.IDeviceStatementsContext:
			for _, nestedChild := range slices.Backward(child.GetChildren()) {
				stack = append(stack, nestedChild)
			}

		case dmlex.IObjectStatementContext:
			if innerIf := child.ObjectIf(); innerIf != nil {
				for _, nestedChild := range slices.Backward(innerIf.GetChildren()) {
					stack = append(stack, nestedChild)
				}
			} else {
				if err := p.printObjectStatement(child); err != nil {
					return err
				}
			}

		case dmlex.IDeviceStatementContext:
			if innerIf := child.ToplevelIf(); innerIf != nil {
				for _, nestedChild := range slices.Backward(innerIf.GetChildren()) {
					stack = append(stack, nestedChild)
				}
			} else {
				if err := p.printDeviceStatement(child); err != nil {
					return err
				}
				if err := p.writeNewLines(1); err != nil {
					return err
				}
			}

		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}

	return nil
}
