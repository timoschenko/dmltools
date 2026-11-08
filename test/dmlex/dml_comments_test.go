// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"strings"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

type commentSyntaxErrorTestCase struct {
	name           string
	code           string
	expectedSymbol string
	expectedSubstr string
}

func runCommentSyntaxErrorTestCase(
	t *testing.T,
	tc *commentSyntaxErrorTestCase,
) {
	t.Helper()

	code := tc.code
	if !strings.HasPrefix(code, "dml 1.4;") {
		code = "dml 1.4;\n" + code
	}

	inputStream := antlr.NewInputStream(code)
	errorListener := dmlex.NewErrorListener("test.dml", 0)
	var syntaxErrors []dmlex.SyntaxErrorInfo
	errorListener.Handler = func(e dmlex.SyntaxErrorInfo) {
		syntaxErrors = append(syntaxErrors, e)
	}

	file := dmlex.ParseStream(inputStream, errorListener)
	_ = file.Parser.Dml()

	if len(syntaxErrors) == 0 {
		t.Fatalf("expected syntax errors, got none")
	}

	firstErr := syntaxErrors[0]
	if tc.expectedSymbol != "" && firstErr.Symbol != tc.expectedSymbol {
		t.Errorf(
			"expected offending symbol %q, got %q",
			tc.expectedSymbol,
			firstErr.Symbol,
		)
	}
	if tc.expectedSubstr != "" &&
		!strings.Contains(firstErr.Msg, tc.expectedSubstr) {
		t.Errorf(
			"expected error message to contain %q, got %q",
			tc.expectedSubstr,
			firstErr.Msg,
		)
	}
}

func TestNonNestingBlockComments_SyntaxErrors(t *testing.T) {
	t.Parallel()

	testCases := []commentSyntaxErrorTestCase{
		{
			name: "01-multiline-nested-pseudo-close",
			code: "/*\n" +
				" * /* nested inner */\n" +
				" */\n",
			expectedSymbol: "*",
			expectedSubstr: "mismatched input '*'",
		},
		{
			name:           "02-slash-star-slash-nested",
			code:           "/* /*/ */\n",
			expectedSymbol: "*",
			expectedSubstr: "mismatched input '*'",
		},
		{
			name:           "03-double-close",
			code:           "/* outer /* inner */ */\n",
			expectedSymbol: "*",
			expectedSubstr: "mismatched input '*'",
		},
		{
			name:           "04-comment-ignored-like-c",
			code:           "/* comment /* still comment */ not a comment */\n",
			expectedSymbol: "not",
			expectedSubstr: "mismatched input 'not'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runCommentSyntaxErrorTestCase(t, &tc)
		})
	}
}

func TestNonNestingBlockComments_LexerTokens(t *testing.T) {
	t.Parallel()

	input := "/* comment /* still comment */ not a comment */"
	inputStream := antlr.NewInputStream(input)
	lexer := dmlex.NewdmlLexer(inputStream)
	tokens := lexer.GetAllTokens()

	// Filter out whitespace tokens if any for cleaner comparison
	var nonWsTokens []antlr.Token
	for _, tok := range tokens {
		if dmlex.TokenType(tok.GetTokenType()) != dmlex.TokenWs {
			nonWsTokens = append(nonWsTokens, tok)
		}
	}

	if len(nonWsTokens) < 6 {
		t.Fatalf(
			"expected at least 6 tokens, got %d: %v",
			len(nonWsTokens),
			nonWsTokens,
		)
	}

	// Token 0: the block comment terminated at the first */
	tok0 := nonWsTokens[0]
	if tok0.GetText() != "/* comment /* still comment */" {
		t.Errorf(
			"expected comment token text '/* comment /* still comment */', got %q",
			tok0.GetText(),
		)
	}
	if tok0.GetChannel() != antlr.TokenHiddenChannel {
		t.Errorf(
			"expected comment token on HIDDEN channel, got %d",
			tok0.GetChannel(),
		)
	}
	tt0 := dmlex.TokenType(tok0.GetTokenType())
	if tt0 != dmlex.TokenCommentInline && tt0 != dmlex.TokenCommentBlock {
		t.Errorf("expected comment token type, got %v", tt0)
	}

	// Subsequent tokens: parsed as source code on DEFAULT channel
	expectedTokens := []struct {
		text string
		tt   dmlex.TokenType
	}{
		{"not", dmlex.TokenIdentifier},
		{"a", dmlex.TokenIdentifier},
		{"comment", dmlex.TokenIdentifier},
		{"*", dmlex.TokenMul},
		{"/", dmlex.TokenDiv},
	}

	for i, expected := range expectedTokens {
		tok := nonWsTokens[i+1]
		if tok.GetText() != expected.text {
			t.Errorf(
				"token %d: expected text %q, got %q",
				i+1,
				expected.text,
				tok.GetText(),
			)
		}
		if tok.GetChannel() != antlr.TokenDefaultChannel {
			t.Errorf(
				"token %d (%q): expected DEFAULT channel, got %d",
				i+1,
				tok.GetText(),
				tok.GetChannel(),
			)
		}
		if dmlex.TokenType(tok.GetTokenType()) != expected.tt {
			t.Errorf(
				"token %d (%q): expected type %v, got %v",
				i+1,
				tok.GetText(),
				expected.tt,
				dmlex.TokenType(tok.GetTokenType()),
			)
		}
	}
}

func TestLineCommentAtEOF_WithoutTrailingNewline(t *testing.T) {
	t.Parallel()

	// Control case: a line comment ending with '\n' parses cleanly as a HIDDEN comment.
	withNewline := "dml 1.4;\n// comment\n"
	errListener1 := dmlex.NewErrorListener("test.dml", 0)
	file1 := dmlex.ParseStream(
		antlr.NewInputStream(withNewline),
		errListener1,
	)
	_ = file1.Parser.Dml()
	if len(errListener1.Errors) != 0 {
		t.Fatalf(
			"unexpected syntax errors with newline: %v",
			errListener1.Errors,
		)
	}

	// Verified fix: a line comment at EOF WITHOUT a trailing newline.
	// With `COMMENT : '//' .*? ('\n' | EOF) -> channel (HIDDEN) ;`,
	// the comment is recognized as TokenComment on HIDDEN channel, no spurious
	// TokenDiv tokens are emitted, and parsing completes with 0 syntax errors.
	withoutNewline := "dml 1.4;\n// trailing comment at EOF without newline"
	lexer := dmlex.NewdmlLexer(antlr.NewInputStream(withoutNewline))
	tokens := lexer.GetAllTokens()

	var divCount int
	var foundComment bool
	for _, tok := range tokens {
		tt := dmlex.TokenType(tok.GetTokenType())
		if tt == dmlex.TokenDiv {
			divCount++
		}
		if tt == dmlex.TokenComment {
			foundComment = true
			if tok.GetChannel() != antlr.TokenHiddenChannel {
				t.Errorf(
					"expected TokenComment on HIDDEN channel, got %d",
					tok.GetChannel(),
				)
			}
			if tok.GetText() != "// trailing comment at EOF without newline" {
				t.Errorf(
					"expected token text '// trailing comment at EOF without newline', got %q",
					tok.GetText(),
				)
			}
		}
	}
	if divCount != 0 {
		t.Errorf(
			"expected 0 TokenDiv tokens from '//' without newline, got %d",
			divCount,
		)
	}
	if !foundComment {
		t.Errorf("expected TokenComment token to be emitted")
	}

	errListener2 := dmlex.NewErrorListener("test.dml", 0)
	file2 := dmlex.ParseStream(
		antlr.NewInputStream(withoutNewline),
		errListener2,
	)
	_ = file2.Parser.Dml()
	if len(errListener2.Errors) != 0 {
		t.Fatalf(
			"unexpected syntax errors for '//' at EOF without newline: %v",
			errListener2.Errors,
		)
	}

	// Verify solo line comment at EOF without trailing newline
	soloComment := "// solo comment at EOF without newline"
	errListener3 := dmlex.NewErrorListener("test.dml", 0)
	file3 := dmlex.ParseStream(
		antlr.NewInputStream(soloComment),
		errListener3,
	)
	_ = file3.Parser.Dml()
	if len(errListener3.Errors) != 0 {
		t.Fatalf(
			"unexpected syntax errors for solo comment at EOF: %v",
			errListener3.Errors,
		)
	}
}
