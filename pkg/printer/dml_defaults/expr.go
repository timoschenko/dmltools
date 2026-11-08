// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmldefaults

import (
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"
)

var exprDefaults staticTokenMap

// Expression returns the baseline FormatFlag for every TokenType
// in the expression context.
func Expression() format.TokenMap {
	exprDefaults.once.Do(initExpressionDefaults)
	return &exprDefaults.data
}

func initExpressionDefaults() {
	tfd := Tokens()

	exprDefaults.data = make(format.TokenMapImpl)
	for i := range dmlex.TokenLastTokenType {
		exprDefaults.data[i] = tfd.GetFormatFlag(i)
	}

	// 1. Structural Tokens (Priority 10)
	// These cling to their adjacent operands or operators.
	const BeatAllExceptComments = format.PriorityAboveNormal //  PriorityHigh)

	// Dot and Arrow: a.b, a->b
	exprDefaults.data[dmlex.TokenDot] = tfd.GetFormatFlag(dmlex.TokenDot).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithBreakBeforePriority(format.PriorityLowest).
		WithSeparatorAfter(format.SepNone).
		WithSeparatorPriority(BeatAllExceptComments)
	exprDefaults.data[dmlex.TokenArrow] = tfd.GetFormatFlag(dmlex.TokenArrow).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithBreakBeforePriority(format.PriorityLowest).
		WithSeparatorAfter(format.SepNone).
		WithSeparatorPriority(BeatAllExceptComments)

	// Parentheses: foo(bar)
	exprDefaults.data[dmlex.TokenLparen] = tfd.GetFormatFlag(dmlex.TokenLparen).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithSeparatorBeforePriority(format.PriorityAboveNormal).
		WithSeparatorAfterPriority(format.PriorityHigh)
	exprDefaults.data[dmlex.TokenRparen] = tfd.GetFormatFlag(dmlex.TokenRparen).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorBeforePriority(format.PriorityHigh).
		WithSeparatorAfterPriority(format.PriorityDefault)

	// Brackets: a[b]
	// Arrays and bitfields in Cdecl need high priority cling
	exprDefaults.data[dmlex.TokenLbrack] = tfd.GetFormatFlag(dmlex.TokenLbrack).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorBeforePriority(format.PriorityAboveNormal).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithSeparatorAfterPriority(format.PriorityHigh).
		WithBreakAfterPriority(format.PriorityHigh) // Should loose to ()
	exprDefaults.data[dmlex.TokenRbrack] = tfd.GetFormatFlag(dmlex.TokenRbrack).
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithSeparatorBeforePriority(format.PriorityHigh).
		WithBreakBeforePriority(format.PriorityHigh). // Should loose to ()
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorAfterPriority(format.PriorityMediumLow)

	// 4. Prefix/Postfix Operators (Baseline)
	ops := []dmlex.TokenType{
		dmlex.TokenInc,
		dmlex.TokenDec,
		dmlex.TokenBang,
		dmlex.TokenTilde,
	}
	for _, op := range ops {
		exprDefaults.data[op] = tfd.GetFormatFlag(op).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorPriority(format.PriorityNormal)
	}

	// {}
	exprDefaults.data[dmlex.TokenLbrace] = tfd.GetFormatFlag(dmlex.TokenLbrace).
		WithIndentIncrease().
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorBeforePriority(format.PriorityLowest).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithSeparatorAfterPriority(format.PriorityMax). // FIXME lowe prio
		WithBreakAfterPriority(format.PriorityVeryHigh)
	exprDefaults.data[dmlex.TokenRbrace] = tfd.GetFormatFlag(dmlex.TokenRbrace).
		WithIndentDecrease().
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithSeparatorBeforePriority(format.PriorityMax). // FIXME lowe prio
		WithBreakBeforePriority(format.PriorityVeryHigh).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorAfterPriority(format.PriorityLowest)

	exprDefaults.data[dmlex.TokenEach] = tfd.GetFormatFlag(dmlex.TokenEach).
		WithSeparatorAfter(format.SepSpace)
	exprDefaults.data[dmlex.TokenIn] = tfd.GetFormatFlag(dmlex.TokenIn).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorAfterPriority(format.PriorityHigh)
	exprDefaults.data[dmlex.TokenNew] = tfd.GetFormatFlag(dmlex.TokenNew).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorAfterPriority(format.PriorityHigh)

	// 5. Ternary and Conditional Operators: ? : #? #:
	ternaryOps := []dmlex.TokenType{
		dmlex.TokenQuestion,
		dmlex.TokenColon,
		dmlex.TokenAltQuestion,
		dmlex.TokenAltColon,
	}
	for _, op := range ternaryOps {
		exprDefaults.data[op] = tfd.GetFormatFlag(op).
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityVeryHigh).
			WithSeparatorBeforePriority(format.PriorityMediumHigh).
			WithSeparatorAfter(format.SepSpace).
			WithBreakAfterPriority(format.PriorityDefault).
			WithSeparatorAfterPriority(format.PriorityMediumHigh)
	}
}
