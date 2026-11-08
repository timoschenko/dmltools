// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout_test

import (
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/layout"

	"github.com/antlr4-go/antlr/v4"
)

type malformMockToken struct {
	antlr.Token
	tokenType int
	text      string
}

func (m malformMockToken) GetTokenType() int { return m.tokenType }
func (m malformMockToken) GetText() string   { return m.text }

func tok(tt dmlex.TokenType, text string) antlr.Token {
	return malformMockToken{
		tokenType: int(tt),
		text:      text,
	}
}

func TestWillTokensMalform(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		t1       antlr.Token
		t2       antlr.Token
		expected bool
	}{
		// Nil or empty tokens
		{"nil-t1", nil, tok(dmlex.TokenAdd, "+"), false},
		{"nil-t2", tok(dmlex.TokenAdd, "+"), nil, false},
		{"empty-t1", tok(dmlex.TokenAdd, ""), tok(dmlex.TokenAdd, "+"), false},
		{"empty-t2", tok(dmlex.TokenAdd, "+"), tok(dmlex.TokenAdd, ""), false},

		// 1. Comment & Code Block starters
		{"div-mul", tok(dmlex.TokenDiv, "/"), tok(dmlex.TokenMul, "*"), true},
		{"div-div", tok(dmlex.TokenDiv, "/"), tok(dmlex.TokenDiv, "/"), true},
		{
			"div-assign",
			tok(dmlex.TokenDiv, "/"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{"mul-div", tok(dmlex.TokenMul, "*"), tok(dmlex.TokenDiv, "/"), true},
		{
			"mul-assign",
			tok(dmlex.TokenMul, "*"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{
			"mul-block-comment",
			tok(dmlex.TokenMul, "*"),
			tok(dmlex.TokenCommentBlock, "/* c */"),
			false,
		},
		{
			"mod-lbrace",
			tok(dmlex.TokenMod, "%"),
			tok(dmlex.TokenLbrace, "{"),
			true,
		},
		{
			"mod-assign",
			tok(dmlex.TokenMod, "%"),
			tok(dmlex.TokenAssign, "="),
			true,
		},

		// 2. Colon compound tokens
		{
			"colon-lbrace",
			tok(dmlex.TokenColon, ":"),
			tok(dmlex.TokenLbrace, "{"),
			true,
		},
		{
			"colon-assign",
			tok(dmlex.TokenColon, ":"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{
			"colon-default",
			tok(dmlex.TokenColon, ":"),
			tok(dmlex.TokenDefault, "default"),
			true,
		},

		// 3. Unary and binary operator collisions
		{"add-add", tok(dmlex.TokenAdd, "+"), tok(dmlex.TokenAdd, "+"), true},
		{
			"add-assign",
			tok(dmlex.TokenAdd, "+"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{"sub-sub", tok(dmlex.TokenSub, "-"), tok(dmlex.TokenSub, "-"), true},
		{
			"sub-assign",
			tok(dmlex.TokenSub, "-"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{"sub-gt", tok(dmlex.TokenSub, "-"), tok(dmlex.TokenGt, ">"), true},
		{
			"bitand-bitand",
			tok(dmlex.TokenBitand, "&"),
			tok(dmlex.TokenBitand, "&"),
			true,
		},
		{
			"bitand-assign",
			tok(dmlex.TokenBitand, "&"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{
			"bitor-bitor",
			tok(dmlex.TokenBitor, "|"),
			tok(dmlex.TokenBitor, "|"),
			true,
		},
		{
			"bitor-assign",
			tok(dmlex.TokenBitor, "|"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{
			"caret-assign",
			tok(dmlex.TokenCaret, "^"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{
			"bang-assign",
			tok(dmlex.TokenBang, "!"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{
			"assign-assign",
			tok(dmlex.TokenAssign, "="),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{"lt-lt", tok(dmlex.TokenLt, "<"), tok(dmlex.TokenLt, "<"), true},
		{
			"lt-assign",
			tok(dmlex.TokenLt, "<"),
			tok(dmlex.TokenAssign, "="),
			true,
		},
		{"gt-gt", tok(dmlex.TokenGt, ">"), tok(dmlex.TokenGt, ">"), true},
		{
			"gt-assign",
			tok(dmlex.TokenGt, ">"),
			tok(dmlex.TokenAssign, "="),
			true,
		},

		// 4. Dot and numerics
		{
			"dot-int",
			tok(dmlex.TokenDot, "."),
			tok(dmlex.TokenIntegerLiteral, "123"),
			true,
		},
		{
			"int-dot",
			tok(dmlex.TokenIntegerLiteral, "123"),
			tok(dmlex.TokenDot, "."),
			true,
		},
		{"dot-dot", tok(dmlex.TokenDot, "."), tok(dmlex.TokenDot, "."), true},

		// 5. Unknown / mock tokens without token type (returns false)
		{
			"unknown-t1",
			tok(dmlex.TokenUnknown, "a"),
			tok(dmlex.TokenIdentifier, "b"),
			false,
		},
		{
			"unknown-t2",
			tok(dmlex.TokenIdentifier, "a"),
			tok(dmlex.TokenUnknown, "b"),
			false,
		},

		// 6. Identifier / Keyword merges
		{
			"ident-ident",
			tok(dmlex.TokenIdentifier, "foo"),
			tok(dmlex.TokenIdentifier, "bar"),
			true,
		},
		{
			"keyword-ident",
			tok(dmlex.TokenReturn, "return"),
			tok(dmlex.TokenIdentifier, "x"),
			true,
		},
		{
			"keyword-keyword",
			tok(dmlex.TokenFor, "for"),
			tok(dmlex.TokenEach, "each"),
			true,
		},
		{
			"ident-int",
			tok(dmlex.TokenIdentifier, "val"),
			tok(dmlex.TokenIntegerLiteral, "1"),
			true,
		},

		// 7. Legitimate SepNone pairs (must NOT be treated as malforming)
		{
			"lparen-rparen",
			tok(dmlex.TokenLparen, "("),
			tok(dmlex.TokenRparen, ")"),
			false,
		},
		{
			"ident-lparen",
			tok(dmlex.TokenIdentifier, "foo"),
			tok(dmlex.TokenLparen, "("),
			false,
		},
		{
			"dot-ident",
			tok(dmlex.TokenDot, "."),
			tok(dmlex.TokenIdentifier, "field"),
			false,
		},
		{
			"arrow-ident",
			tok(dmlex.TokenArrow, "->"),
			tok(dmlex.TokenIdentifier, "field"),
			false,
		},
		{
			"mul-ident",
			tok(dmlex.TokenMul, "*"),
			tok(dmlex.TokenIdentifier, "ptr"),
			false,
		},
		{
			"sub-int",
			tok(dmlex.TokenSub, "-"),
			tok(dmlex.TokenIntegerLiteral, "1"),
			false,
		},
		{
			"add-int",
			tok(dmlex.TokenAdd, "+"),
			tok(dmlex.TokenIntegerLiteral, "1"),
			false,
		},
		{
			"bang-ident",
			tok(dmlex.TokenBang, "!"),
			tok(dmlex.TokenIdentifier, "flag"),
			false,
		},
		{
			"tilde-ident",
			tok(dmlex.TokenTilde, "~"),
			tok(dmlex.TokenIdentifier, "mask"),
			false,
		},
		{
			"ident-lbrack",
			tok(dmlex.TokenIdentifier, "arr"),
			tok(dmlex.TokenLbrack, "["),
			false,
		},
		{
			"lbrack-int",
			tok(dmlex.TokenLbrack, "["),
			tok(dmlex.TokenIntegerLiteral, "0"),
			false,
		},
		{
			"rbrack-semi",
			tok(dmlex.TokenRbrack, "]"),
			tok(dmlex.TokenSemi, ";"),
			false,
		},
		{
			"rparen-semi",
			tok(dmlex.TokenRparen, ")"),
			tok(dmlex.TokenSemi, ";"),
			false,
		},
		{
			"string-semi",
			tok(dmlex.TokenStringLiteral, "\"str\""),
			tok(dmlex.TokenSemi, ";"),
			false,
		},
		{
			"char-semi",
			tok(dmlex.TokenCharLiteral, "'c'"),
			tok(dmlex.TokenSemi, ";"),
			false,
		},
		{
			"comment-block-semi",
			tok(dmlex.TokenCommentBlock, "/* c */"),
			tok(dmlex.TokenSemi, ";"),
			false,
		},
		{
			"comment-block-ident",
			tok(dmlex.TokenCommentBlock, "/* c */"),
			tok(dmlex.TokenIdentifier, "x"),
			false,
		},
		{
			"comment-block-mul",
			tok(dmlex.TokenCommentBlock, "/* c */"),
			tok(dmlex.TokenMul, "*"),
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := layout.WillTokensMalform(tc.t1, tc.t2)
			if got != tc.expected {
				t.Errorf(
					"WillTokensMalform(%q (%s), %q (%s)) = %v, want %v",
					tc.t1.GetText(), dmlex.GetTokenType(tc.t1),
					tc.t2.GetText(), dmlex.GetTokenType(tc.t2),
					got, tc.expected,
				)
			}
		})
	}
}
