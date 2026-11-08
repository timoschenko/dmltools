// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmldefaults

import (
	"sync"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"
)

type staticTokenMap struct {
	once sync.Once
	data format.TokenMapImpl
}

var tokenFormatDefaults staticTokenMap

// Tokens returns the baseline FormatFlag for every TokenType
// in the expression context.
func Tokens() format.TokenMap {
	tokenFormatDefaults.once.Do(initTokenFormatDefaults)
	return &tokenFormatDefaults.data
}

func initTokenFormatDefaults() {
	tokenFormatDefaults.data = make(format.TokenMapImpl)

	defaultsFlags := format.FormatFlag(0).
		WithSeparatorBefore(format.SepSpace).
		WithBreakBeforePriority(format.PriorityDefault).
		WithSeparatorAfter(format.SepSpace).
		WithBreakAfterPriority(format.PriorityDefault).
		WithSeparatorPriority(format.PriorityDefault)
	for i := range dmlex.TokenLastTokenType {
		tokenFormatDefaults.data[i] = defaultsFlags
	}

	// 2. Separators
	// Comma: a, b
	tokenFormatDefaults.data[dmlex.TokenComma] = tokenFormatDefaults.data[dmlex.TokenComma].
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithBreakAfterPriority(format.PriorityAboveNormal).
		WithSeparatorBeforePriority(format.PriorityHigh).
		WithSeparatorAfterPriority(format.PriorityHigh)

	// Semicolon: stmt;
	tokenFormatDefaults.data[dmlex.TokenSemi] = tokenFormatDefaults.data[dmlex.TokenSemi].
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepMandatoryWrap).

		// NOTE: Using higher token priority for semicolon due to inline comments
		// to force following behavior: `/**/ ;`
		WithSeparatorBeforePriority(format.PriorityAboveHighest)

	// AS: export expr as expr;
	tokenFormatDefaults.data[dmlex.TokenAs] = tokenFormatDefaults.data[dmlex.TokenAs].
		WithSeparatorBefore(format.SepSpaceMaybeWrap).
		WithBreakBeforePriority(format.PriorityLowest).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorPriority(format.PriorityHigh)

	// Delete statement: delete expr;
	tokenFormatDefaults.data[dmlex.TokenDelete] = tokenFormatDefaults.data[dmlex.TokenDelete].
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorPriority(format.PriorityHigh)

	// Assert statement: assert expr;
	tokenFormatDefaults.data[dmlex.TokenAssert] = tokenFormatDefaults.data[dmlex.TokenAssert].
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorPriority(format.PriorityHigh)

	// Goto statement: goto ident;
	tokenFormatDefaults.data[dmlex.TokenGoto] = tokenFormatDefaults.data[dmlex.TokenGoto].
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorAfterPriority(format.PriorityHigh)

	// Then in log statement: log info, 1 then 2: "msg";
	tokenFormatDefaults.data[dmlex.TokenThen] = tokenFormatDefaults.data[dmlex.TokenThen].
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorPriority(format.PriorityHigh)

	// Brackets and parentheses: a[0], foo(bar), { ... }
	tokenFormatDefaults.data[dmlex.TokenLparen] = tokenFormatDefaults.data[dmlex.TokenLparen].
		WithIndentIncrease().
		WithSeparatorBefore(format.SepNone).
		WithSeparatorBeforePriority(format.PriorityMediumLow).
		WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
		WithBreakAfterPriority(format.PriorityHigh).
		WithSeparatorAfterPriority(format.PriorityHigh)
	tokenFormatDefaults.data[dmlex.TokenRparen] = tokenFormatDefaults.data[dmlex.TokenRparen].
		WithIndentDecrease().
		WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
		WithBreakBeforePriority(format.PriorityDefault).
		WithSeparatorBeforePriority(format.PriorityMediumHigh)

		// above the /**/ comment

	//
	tokenFormatDefaults.data[dmlex.TokenLbrace] = tokenFormatDefaults.data[dmlex.TokenLbrace].
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepMandatoryWrap).
		WithIndentIncrease()
	tokenFormatDefaults.data[dmlex.TokenRbrace] = tokenFormatDefaults.data[dmlex.TokenRbrace].
		WithSeparatorBefore(format.SepMandatoryWrap).
		WithIndentDecrease()
	// Else and #else
	tokenFormatDefaults.data[dmlex.TokenElse] = tokenFormatDefaults.data[dmlex.TokenElse].
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorBeforePriority(format.PriorityAboveNormal)
	tokenFormatDefaults.data[dmlex.TokenHashElse] = tokenFormatDefaults.data[dmlex.TokenHashElse].
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorBeforePriority(format.PriorityAboveNormal)

	tokenFormatDefaults.data[dmlex.TokenLbrack] = tokenFormatDefaults.data[dmlex.TokenLbrack].
		WithIndentIncrease()
	tokenFormatDefaults.data[dmlex.TokenRbrack] = tokenFormatDefaults.data[dmlex.TokenRbrack].
		WithIndentDecrease()

	// 3. Binary/Logical/Relational Operators
	// These must beat LPAREN's PriorityAboveNormal (8) to preserve spaces in grouped expressions like `a * (b)`
	binOps := []dmlex.TokenType{
		dmlex.TokenMul,
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
		dmlex.TokenOr,
	}
	for _, op := range binOps {
		tokenFormatDefaults.data[op] = tokenFormatDefaults.data[op].
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityDefault).
			WithSeparatorAfter(format.SepSpace).
			// has to be more than LPAREN
			WithSeparatorPriority(format.PriorityMediumHigh) // (9) > LPAREN (8)
	}

	assignOps := []dmlex.TokenType{
		dmlex.TokenAssign,
		dmlex.TokenAddAssign,
		dmlex.TokenSubAssign,
		dmlex.TokenMulAssign,
		dmlex.TokenDivAssign,
		dmlex.TokenModAssign,
		dmlex.TokenLshiftAssign,
		dmlex.TokenRshiftAssign,
		dmlex.TokenAndAssign,
		dmlex.TokenXorAssign,
		dmlex.TokenOrAssign,
	}
	for _, op := range assignOps {
		tokenFormatDefaults.data[op] = tokenFormatDefaults.data[op].
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorAfter(format.SepSpace).
			WithSeparatorPriority(format.PriorityMediumHigh) // (9) > LPAREN (8)
	}

	tokenFormatDefaults.data[dmlex.TokenAt] = tokenFormatDefaults.data[dmlex.TokenAt].
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorPriority(format.PriorityMediumHigh)

	// 4. Ternary and Conditional Operators: ? : #? #:
	ternaryOps := []dmlex.TokenType{
		dmlex.TokenQuestion,
		dmlex.TokenColon,
		dmlex.TokenAltQuestion,
		dmlex.TokenAltColon,
	}
	for _, op := range ternaryOps {
		tokenFormatDefaults.data[op] = tokenFormatDefaults.data[op].
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityVeryHigh).
			WithSeparatorBeforePriority(format.PriorityMediumHigh).
			WithSeparatorAfter(format.SepSpace).
			WithBreakAfterPriority(format.PriorityDefault).
			WithSeparatorAfterPriority(format.PriorityMediumHigh)
	}

	// 5. Comments

	// Line comment: // ...
	tokenFormatDefaults.data[dmlex.TokenComment] = tokenFormatDefaults.data[dmlex.TokenComment].
		WithSeparatorBefore(format.SepSpace).

		// NOTE: Comment shall beat the comma to avoid going to the next line
		WithSeparatorBeforePriority(format.PriorityMax).
		WithSeparatorAfter(format.SepMandatoryWrap).

		// NOTE: comment should always end with a new line
		//       otherwise next token will become a comment as well
		WithSeparatorAfterPriority(format.PriorityMax)

	// Single-line comment: /* \n ... */ (COMMENT_INLINE in g4)
	tokenFormatDefaults.data[dmlex.TokenCommentInline] = tokenFormatDefaults.data[dmlex.TokenCommentInline].
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorPriority(format.PriorityHigh)

	tokenFormatDefaults.data[dmlex.TokenCommentInlineEol] = tokenFormatDefaults.data[dmlex.TokenCommentInline].
		WithSeparatorAfter(format.SepMandatoryWrap).
		WithSeparatorAfterPriority(format.PriorityMax).
		WithSeparatorBeforePriority(format.PriorityMediumHigh)

	// Multi-line block comment: /* ... */ (COMMENT_BLOCK in g4)
	tokenFormatDefaults.data[dmlex.TokenCommentBlock] = tokenFormatDefaults.data[dmlex.TokenCommentBlock].
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorPriority(format.PriorityMediumHigh)

		// PriorityBelowNormal PriorityHighest
	tokenFormatDefaults.data[dmlex.TokenCommentBlockEol] = tokenFormatDefaults.data[dmlex.TokenCommentBlock].
		WithSeparatorAfter(format.SepMandatoryWrap).
		WithSeparatorAfterPriority(format.PriorityMax).
		WithSeparatorBeforePriority(format.PriorityMediumHigh)

	// 5. Special cases from porting.go
	// TokenWs and TokenEOF should not have separators or priority impact
	tokenFormatDefaults.data[dmlex.TokenWs] = format.FormatFlag(0)
	tokenFormatDefaults.data[dmlex.TokenEOF] = format.FormatFlag(0)
}
