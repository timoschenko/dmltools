// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"strings"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

type stringLiteralValidTestCase struct {
	name   string
	code   string
	masked string
}

type stringLiteralSyntaxErrorTestCase struct {
	name   string
	code   string
	masked string
}

func parseDmlCode(code string) ([]dmlex.SyntaxErrorInfo, *dmlex.File) {
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
	return syntaxErrors, file
}

func runStringLiteralValidTestCase(
	t *testing.T,
	tc *stringLiteralValidTestCase,
) {
	t.Helper()

	syntaxErrors, _ := parseDmlCode(tc.code)

	if len(syntaxErrors) > 0 {
		if tc.masked != "" {
			t.Skipf("masked: %s (syntax errors: %v)", tc.masked, syntaxErrors)
			return
		}
		t.Fatalf(
			"unexpected syntax errors (%d): %v",
			len(syntaxErrors),
			syntaxErrors,
		)
	}

	if tc.masked != "" {
		t.Fatalf(
			"test case %q passed but is still masked: %s. Consider removing the mask",
			tc.name,
			tc.masked,
		)
	}
}

func runStringLiteralSyntaxErrorTestCase(
	t *testing.T,
	tc *stringLiteralSyntaxErrorTestCase,
) {
	t.Helper()

	syntaxErrors, _ := parseDmlCode(tc.code)

	if len(syntaxErrors) == 0 {
		if tc.masked != "" {
			t.Skipf(
				"masked: %s (expected syntax error, but parsed cleanly)",
				tc.masked,
			)
			return
		}
		t.Fatalf("expected syntax error, but parsed cleanly without errors")
	}

	if tc.masked != "" {
		t.Fatalf(
			"test case %q failed with syntax error as expected but is still masked: %s. Consider removing the mask",
			tc.name,
			tc.masked,
		)
	}
}

func TestStringLiteral_Valid(t *testing.T) {
	t.Parallel()

	testCases := []stringLiteralValidTestCase{
		// Basic fine cases from user specification
		{
			name: "fine-hello-world",
			code: `param p = "hello world";`,
		},
		{
			name: "fine-single-quote-inside",
			code: `param p = "i can't";`,
		},
		{
			name: "fine-escaped-quotes-with-suffix",
			code: `param p = "hello \"world\" - 123";`,
		},
		{
			name: "fine-escaped-quotes-enclosing-word",
			code: `param p = "hello \"world\"";`,
		},
		{
			name: "fine-escaped-quote-in-middle",
			code: `param p = "hello \" quote";`,
		},
		{
			name: "fine-escaped-quote-before-closing-quote",
			code: `param p = "hello quote\"";`,
		},
		{
			name: "fine-single-quoted-word-inside",
			code: `param p = "hello 'quote'";`,
		},

		// Additional escape sequences
		{
			name: "fine-c-standard-escapes",
			code: `param p = "line1\nline2\ttab\r\b\a\f\v";`,
		},
		{
			name: "fine-escaped-backslash",
			code: `param p = "backslash\\\\test";`,
		},
		{
			name: "fine-trailing-escaped-backslash",
			code: `param p = "trailing\\\\";`,
		},
		{
			name: "fine-hex-escapes",
			code: `param p = "\x00\x1f\xFF\xab";`,
		},
		{
			name: "fine-octal-escapes",
			code: `param p = "\0\000\123";`,
		},
		{
			name: "fine-escaped-single-quote-in-string",
			code: `param p = "it\'s fine";`,
		},
		{
			name: "fine-empty-string",
			code: `param p = "";`,
		},
		{
			name: "fine-only-escaped-quote",
			code: `param p = "\"";`,
		},
		{
			name: "fine-multiple-escaped-quotes",
			code: `param p = "\"\"\"";`,
		},
		{
			name: "fine-utf8-unicode",
			code: `param p = "Привет мир! 🌍 日本語";`,
		},
		{
			name: "fine-multiline-literal-newline",
			code: "param p = \"line1\nline2\";",
		},

		// Contexts: log statements
		{
			name: "fine-log-info-simple",
			code: `method m() { log info: "hello world"; }`,
		},
		{
			name: "fine-log-info-with-escape",
			code: `method m() { log info: "hello \"world\""; }`,
		},
		{
			name: "fine-log-info-with-level-and-group",
			code: `method m() { log info, 1, mygroup: "hello \"world\""; }`,
		},
		{
			name: "fine-log-info-format-with-args",
			code: `method m() { log info: "val=%d, s=%s", 42, "str"; }`,
		},

		// Contexts: warning and error statements
		{
			name: "fine-warning-simple",
			code: `method m() { warning "hello world"; }`,
		},
		{
			name: "fine-warning-with-escape",
			code: `method m() { warning "hello \"world\""; }`,
		},
		{
			name: "fine-warning-bracketed-concat",
			code: `method m() { warning ("hello " + "world"); }`,
		},
		{
			name: "fine-error-with-escape",
			code: `method m() { error "fatal \"failure\""; }`,
		},

		// Contexts: string concatenation
		{
			name: "fine-string-concat",
			code: `param p = "hello " + "world";`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runStringLiteralValidTestCase(t, &tc)
		})
	}
}

func TestStringLiteral_SyntaxErrors(t *testing.T) {
	t.Parallel()

	testCases := []stringLiteralSyntaxErrorTestCase{
		// Syntax error cases from user specification
		{
			name: "syntax-err-unescaped-quote-unclosed",
			code: `param p = "hello "world";`,
		},
		{
			name: "syntax-err-unescaped-quote-with-trailing-string",
			code: `param p = "hello "world"";`,
		},
		{
			name: "syntax-err-single-quoted-multi-char-string",
			code: `param p = 'hello "world"';`,
		},

		// Additional syntax error cases
		{
			name: "syntax-err-unclosed-double-quote-at-eof",
			code: `param p = "unclosed;`,
		},
		{
			name: "syntax-err-unclosed-trailing-escaped-quote",
			code: `param p = "trailing escape\";`,
		},
		{
			name: "syntax-err-multi-char-single-quote-literal",
			code: `param p = 'multi-char';`,
		},
		{
			name: "syntax-err-unclosed-single-quote",
			code: `param p = 'unclosed;`,
		},
		{
			name: "syntax-err-warning-char-literal",
			code: `method m() { warning 'a'; }`,
		},
		{
			name: "syntax-err-log-unclosed-string",
			code: `method m() { log info: "unclosed; }`,
		},
		{
			name: "syntax-err-log-unescaped-inner-quote",
			code: `method m() { log info: "hello "world""; }`,
		},
		{
			name: "syntax-err-warning-unescaped-inner-quote",
			code: `method m() { warning "hello "world""; }`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runStringLiteralSyntaxErrorTestCase(t, &tc)
		})
	}
}

func TestStringLiteral_LexerTokens(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		input        string
		expectedText string
	}{
		{
			name:         "simple-string",
			input:        `"hello world"`,
			expectedText: `"hello world"`,
		},
		{
			name:         "inner-single-quote",
			input:        `"i can't"`,
			expectedText: `"i can't"`,
		},
		{
			name:         "escaped-double-quotes",
			input:        `"hello \"world\""`,
			expectedText: `"hello \"world\""`,
		},
		{
			name:         "escaped-quote-at-end",
			input:        `"hello quote\""`,
			expectedText: `"hello quote\""`,
		},
		{
			name:         "c-escapes",
			input:        `"line1\nline2\t"`,
			expectedText: `"line1\nline2\t"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputStream := antlr.NewInputStream(tc.input)
			lexer := dmlex.NewdmlLexer(inputStream)
			tokens := lexer.GetAllTokens()

			if len(tokens) != 1 {
				t.Fatalf(
					"expected exactly 1 token, got %d: %v",
					len(tokens),
					tokens,
				)
			}

			tok := tokens[0]
			if dmlex.TokenType(tok.GetTokenType()) != dmlex.TokenStringLiteral {
				t.Errorf(
					"expected TokenStringLiteral, got %v",
					dmlex.TokenType(tok.GetTokenType()),
				)
			}
			if tok.GetChannel() != antlr.TokenDefaultChannel {
				t.Errorf(
					"expected channel %d, got %d",
					antlr.TokenDefaultChannel,
					tok.GetChannel(),
				)
			}
			if tok.GetText() != tc.expectedText {
				t.Errorf(
					"expected text %q, got %q",
					tc.expectedText,
					tok.GetText(),
				)
			}
		})
	}
}
