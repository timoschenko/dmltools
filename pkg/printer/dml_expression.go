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

type rootOrNested int8

const (
	rootNode rootOrNested = iota
	nestedNode
)

// g4-tag: expression
func (p *Printer) printExpression(node dmlex.IExpressionContext) error {
	if node == nil {
		return nil
	}

	p.newContext(p.settings.ExpressionDefaults())
	defer p.freeContext()

	return p.printExpressionRaw(node, rootNode)
}

func (p *Printer) printExpressionRaw(
	node dmlex.IExpressionContext,
	rn rootOrNested,
) error {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *dmlex.LiteralAtomContext:
		return p.printLiteralAtom(n, rn)
	case *dmlex.PostIncrementExprContext, *dmlex.PreIncrementExprContext,
		*dmlex.UnaryExprContext, *dmlex.NotExprContext,
		*dmlex.BitwiseNotExprContext, *dmlex.AddressofExprContext,
		*dmlex.DerefExprContext:
		return p.printUnaryExpr(n, rn)
	case *dmlex.MultiplicativeExprContext, *dmlex.AdditiveExprContext,
		*dmlex.ShiftExprContext, *dmlex.RelationalExprContext,
		*dmlex.EqualityExprContext, *dmlex.BitwiseExprContext,
		*dmlex.LogicalExprContext:
		return p.printBinaryExpr(n, rn)
	case *dmlex.TernaryIfExprContext, *dmlex.TernaryConstIfExprContext:
		return p.printTernaryExpr(n, rn)
	case *dmlex.ParenExprContext:
		return p.printParenExpr(n, rn)
	case *dmlex.SizeofExprContext, *dmlex.SizeoftypeExprContext,
		*dmlex.DefinedExprContext, *dmlex.StringifyExprContext:
		return p.printBuiltinUnaryExpr(n, rn)
	case *dmlex.ArrayAccessExprContext:
		return p.printArrayAccessExpr(n, rn)
	case *dmlex.BitSliceExprContext:
		return p.printBitSliceExpr(n, rn)
	case *dmlex.MemberAccessExprContext:
		return p.printMemberAccessExpr(n, rn)
	case *dmlex.FunctionCallExprContext:
		return p.printFunctionCallExpr(n, rn)
	case *dmlex.ArrayLiteralExprContext:
		return p.printArrayLiteralExpr(n, rn)
	case *dmlex.EachExprContext:
		return p.printEachExpr(n, rn)
	case *dmlex.NewExprContext:
		return p.printNewExpr(n, rn)
	case *dmlex.CastExprContext:
		return p.printCastExpr(n, rn)
	default:
		panic("unreachable: " + fmt.Sprintf("%T", n))
	}
}

// g4: literalAtom
// g4: expressionIdent
func (p *Printer) printLiteralAtom(
	node antlr.Tree,
	rn rootOrNested,
) error {
	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(
				func() bool { return false },
				"no terminal nodes expected",
			)
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext: // TODO: Unwrap here to avoid recursion
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
			}

		case dmlex.IExpressionIdentContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

//go:inline
func formatUnaryPost(flag format.FormatFlag) format.FormatFlag {
	return flag.
		WithSeparatorBefore(format.SepNone).
		WithSeparatorBeforePriority(format.PriorityHigh)
}

//go:inline
func formatUnaryPre(flag format.FormatFlag) format.FormatFlag {
	return flag.
		WithSeparatorAfter(format.SepNone).
		WithSeparatorAfterPriority(format.PriorityHigh)
}

func (p *Printer) printUnaryExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("UnaryExpr")
	defer endGroup()

	_, isPost := node.(*dmlex.PostIncrementExprContext)
	formatFn := formatUnaryPre
	if isPost {
		formatFn = formatUnaryPost
	}

	children := node.GetChildren()
	if !isPost && ensureSpaceForUnaryExpr_Workaround(children) {
		formatFn = formatUnaryPreSpace
	}

	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenInc,
					dmlex.TokenDec,
					dmlex.TokenSub,
					dmlex.TokenAdd,
					dmlex.TokenBang,
					dmlex.TokenTilde,
					dmlex.TokenBitand,
					dmlex.TokenMul:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenWithCustomFormat(child, formatFn); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

//go:inline
func isSameBinaryExprType(a, b antlr.Tree) bool {
	switch a.(type) {
	case *dmlex.MultiplicativeExprContext:
		_, ok := b.(*dmlex.MultiplicativeExprContext)
		return ok
	case *dmlex.AdditiveExprContext:
		_, ok := b.(*dmlex.AdditiveExprContext)
		return ok
	case *dmlex.ShiftExprContext:
		_, ok := b.(*dmlex.ShiftExprContext)
		return ok
	case *dmlex.RelationalExprContext:
		_, ok := b.(*dmlex.RelationalExprContext)
		return ok
	case *dmlex.EqualityExprContext:
		_, ok := b.(*dmlex.EqualityExprContext)
		return ok
	case *dmlex.BitwiseExprContext:
		_, ok := b.(*dmlex.BitwiseExprContext)
		return ok
	case *dmlex.LogicalExprContext:
		_, ok := b.(*dmlex.LogicalExprContext)
		return ok
	default:
		return false
	}
}

func (p *Printer) printBinaryExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	if rn == rootNode {
		endGroup := p.LayoutEngine.BeginWrapGroup("BinaryExpr")
		defer endGroup()
	}

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
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenMul,
					dmlex.TokenDiv,
					dmlex.TokenMod,
					dmlex.TokenAdd,
					dmlex.TokenSub,
					dmlex.TokenLshift,
					dmlex.TokenRshift,
					dmlex.TokenLt,
					dmlex.TokenGt,
					dmlex.TokenLe,
					dmlex.TokenGe,
					dmlex.TokenEqual,
					dmlex.TokenNotequal,
					dmlex.TokenBitand,
					dmlex.TokenCaret,
					dmlex.TokenBitor,
					dmlex.TokenAnd,
					dmlex.TokenOr:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if isSameBinaryExprType(node, child) {
				for _, nestedChild := range slices.Backward(child.GetChildren()) {
					stack = append(stack, nestedChild)
				}
			} else {
				if err := p.printExpressionRaw(child, rootNode); err != nil {
					return err
				}
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

func (p *Printer) printTernaryExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	if rn == rootNode {
		endGroup := p.LayoutEngine.BeginWrapGroup("TernaryExpr")
		defer endGroup()
	}

	var endSoftWrap func()
	defer func() {
		if endSoftWrap != nil {
			endSoftWrap()
		}
	}()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenQuestion,
					dmlex.TokenColon,
					dmlex.TokenAltQuestion,
					dmlex.TokenAltColon:
					return true
				}
				return false
			}, "token is not supported")
			if rn == rootNode && endSoftWrap == nil &&
				(tt == dmlex.TokenQuestion || tt == dmlex.TokenAltQuestion) {
				endSoftWrap = p.LayoutEngine.BeginSoftIndent()
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			childRn := rootNode
			switch child.(type) {
			case *dmlex.TernaryIfExprContext, *dmlex.TernaryConstIfExprContext:
				childRn = nestedNode
			}
			if err := p.printExpressionRaw(child, childRn); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

func (p *Printer) printParenExpr(
	node antlr.Tree,
	_ rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("ParenExpr")
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
		case dmlex.IExpressionContext:
			if err := p.printExpressionRaw(child, rootNode); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4-tag: sizeofExpr, sizeoftypeExpr, definedExpr, stringifyExpr
func (p *Printer) printBuiltinUnaryExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	// FIXME: Root detection might be needed here (previously used `isRoot` workaround).
	endGroup := p.LayoutEngine.BeginWrapGroup("BuiltinUnaryExpr")
	defer endGroup()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenSizeof,
					dmlex.TokenSizeoftype,
					dmlex.TokenDefined,
					dmlex.TokenStringify,
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
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
			}
		case dmlex.ITypeopargContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

func (p *Printer) printArrayAccessExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("ArrayAccess")
	defer endGroup()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenLbrack, dmlex.TokenRbrack:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

func (p *Printer) printBitSliceExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("BitSlice")
	defer endGroup()

	overrides := map[dmlex.TokenType]format.FormatFlag{
		dmlex.TokenColon: p.context.GetFormatFlag(dmlex.TokenColon).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveNormal).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityAboveNormal).
			WithBreakBeforePriority(format.PriorityLowest).
			WithBreakAfterPriority(format.PriorityLowest),
		dmlex.TokenComma: p.context.GetFormatFlag(dmlex.TokenComma).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepSpace).
			// no wrapping in slice [12:21, le]
			WithSeparatorPriority(format.PriorityHigh),
		dmlex.TokenLbrack: p.context.GetFormatFlag(dmlex.TokenLbrack).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityAboveHighest),
	}
	p.newContext(format.NewTokenMap(overrides))
	defer p.freeContext()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenLbrack, dmlex.TokenRbrack:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
			}

		case dmlex.IBitSliceContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

func (p *Printer) printMemberAccessExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginGreedyGroup("MemberAccess")
	defer endGroup()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenDot, dmlex.TokenArrow:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
			}

		case dmlex.IExpressionAccessIdentContext:
			// TODO: Use this - if err := p.printLiteralAtom(child, rn); err != nil {
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

func (p *Printer) printFunctionCallExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("FunctionCall")
	defer endGroup()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			// NOTE: comments only
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			switch realType := child.(type) {
			case *dmlex.LiteralAtomContext:
				if err := p.printLiteralAtom(realType, rn); err != nil {
					return err
				}
			case *dmlex.MemberAccessExprContext:
				if err := p.printMemberAccessExpr(realType, rn); err != nil {
					return err
				}
			default:
				if err := p.printExpressionRaw(child, rn); err != nil {
					return err
				}
			}
		case dmlex.IArgumentsListContext:
			if err := p.printArgumentsList(child); err != nil {
				return err
			}
		case dmlex.IArgumentsListEmptyContext:
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

// g4: tuple_literal
// g4: expression_list
// g4: expression_list_ntc_nonempty
func (p *Printer) printExpression_listExpr(
	node antlr.Tree,
	_ rootOrNested,
) error {
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
				case dmlex.TokenComma, dmlex.TokenLparen, dmlex.TokenRparen:
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
		case dmlex.ITuple_literalContext,
			dmlex.IExpression_listContext,
			dmlex.IExpression_list_ntc_nonemptyContext:
			// Iterates in reverse, appending so the first child stays on top
			for _, nestedChild := range slices.Backward(child.GetChildren()) {
				stack = append(stack, nestedChild)
			}
		}
	}

	return nil
}

func (p *Printer) printArrayLiteralExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("ArrayLiteral")
	defer func() {
		if endGroup != nil {
			endGroup()
		}
	}()

	var endElementsGroup func()
	defer func() {
		if endElementsGroup != nil {
			endElementsGroup()
		}
	}()

	hasNested := hasNestedArrayLiteral_Workaround(node)

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
			if tt == dmlex.TokenRbrack {
				if endElementsGroup != nil {
					endElementsGroup()
					endElementsGroup = nil
				}
				rbrackIdx := TokenIndex(child.GetSymbol().GetTokenIndex())
				if err := p.pushTokenRange(rbrackIdx); err != nil {
					return err
				}
				if endGroup != nil {
					endGroup()
					endGroup = nil
				}
				rightiestIdx := p.gatherTokenIndexFromRightiest(
					child.GetSymbol(),
				)
				if rightiestIdx > rbrackIdx {
					if err := p.pushTokenRange(rightiestIdx); err != nil {
						return err
					}
				}
				continue
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if tt == dmlex.TokenLbrack && endElementsGroup == nil {
				endElementsGroup = p.LayoutEngine.BeginWrapGroup(
					"arrayElements",
				)
				if hasNested {
					p.LayoutEngine.TriggerCurrentGroup()
				}
			}
		case dmlex.IExpressionContext:
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
			}

		case dmlex.IExpression_listContext:
			if err := p.printExpression_listExpr(child, rn); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

func (p *Printer) printEachExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("EachExpr")
	defer endGroup()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenEach,
					dmlex.TokenIn,
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
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
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

func (p *Printer) printNewExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("NewExpr")
	defer endGroup()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			assertLazy(func() bool {
				switch dmlex.GetTokenType(child.GetSymbol()) {
				case dmlex.TokenNew, dmlex.TokenLbrack, dmlex.TokenRbrack:
					return true
				}
				return false
			}, "token is not supported")
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
		case dmlex.IExpressionContext:
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
			}

		case dmlex.ICtypedeclContext:
			if err := p.printCtypedecl(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}

func (p *Printer) printCastExpr(
	node antlr.Tree,
	rn rootOrNested,
) error {
	endGroup := p.LayoutEngine.BeginWrapGroup("CastExpr")
	defer endGroup()

	var endArgsGroup func()
	defer func() {
		if endArgsGroup != nil {
			endArgsGroup()
		}
	}()

	children := node.GetChildren()
	for _, rawChild := range children {
		switch child := rawChild.(type) {
		case antlr.TerminalNode:
			tt := dmlex.GetTokenType(child.GetSymbol())
			assertLazy(func() bool {
				switch tt {
				case dmlex.TokenCast, dmlex.TokenLparen,
					dmlex.TokenComma, dmlex.TokenRparen:
					return true
				}
				return false
			}, "token is not supported")
			if tt == dmlex.TokenRparen && endArgsGroup != nil {
				endArgsGroup()
				endArgsGroup = nil
			}
			if err := p.pushTokenRangeFromRightiest(child); err != nil {
				return err
			}
			if tt == dmlex.TokenLparen && endArgsGroup == nil {
				endArgsGroup = p.LayoutEngine.BeginWrapGroup("castArgs")
			}
		case dmlex.IExpressionContext:
			if err := p.printExpressionRaw(child, rn); err != nil {
				return err
			}

		case dmlex.ICtypedeclContext:
			if err := p.printCtypedecl(child); err != nil {
				return err
			}
		default:
			panic("unreachable: " + fmt.Sprintf("%T", child))
		}
	}
	return nil
}
