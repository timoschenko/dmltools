// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"strings"
	"testing"
)

func TestErrorAndWarningStatement(t *testing.T) {
	const lineCap = 40

	testCases := []struct {
		name, code, expected, kwOnly string
		masked                       Masked
	}{
		{
			name:     "01-basic-statement",
			code:     "KW `test message`;",
			expected: "KW `test message`;",
		},
		{
			name:     "02-messy-whitespace",
			code:     "    KW    `test message`   ;",
			expected: "KW `test message`;",
		},
		{
			name:     "03-long-string-literal-wrapping",
			code:     "KW `this is a very long error or warning message exceeding line cap size limit`;",
			expected: "KW `this is a very long error or warning message exceeding line cap size limit`;",
		},
		{
			name:     "04-comment-alignment-inline-block",
			code:     "KW   /* block */   `test message`; // end",
			expected: "KW /* block */ `test message`; // end",
		},
		{
			name:     "05-comment-alignment-line-before-after",
			code:     "// comment before\nKW `test message`;\n// comment after",
			expected: "// comment before\nKW `test message`;\n// comment after",
		},
		{
			name:     "06-comment-alignment-block-before-semi",
			code:     "KW   `test message`   /* comment */  ;",
			expected: "KW `test message` /* comment */;",
		},
		{
			name:     "07-error-without-string-literal",
			code:     "error   ;",
			expected: "error;",
			kwOnly:   "error",
		},
		{
			name:     "08-string-concatenation",
			code:     "KW `hello ` + `world`;",
			expected: "KW `hello ` + `world`;",
		},
		{
			name:     "09-string-concatenation-messy-whitespace",
			code:     "KW    `hello `   +    `world`   ;",
			expected: "KW `hello ` + `world`;",
		},
		{
			name:     "10-multi-string-concatenation",
			code:     "KW `a ` + `b ` + `c`;",
			expected: "KW `a ` + `b ` + `c`;",
		},
		{
			name:     "11-parenthesized-single-string",
			code:     "KW ( `test message` );",
			expected: "KW (`test message`);",
		},
		{
			name:     "12-parenthesized-concatenation",
			code:     "KW ( `hello ` + `world` );",
			expected: "KW (`hello ` + `world`);",
		},
		{
			name:     "13-parenthesized-multi-concatenation",
			code:     "KW ( `a ` + `b ` + `c` );",
			expected: "KW (`a ` + `b ` + `c`);",
		},
		{
			name:     "14-parenthesized-with-inline-comments",
			code:     "KW ( /* start */ `hello ` + /* mid */ `world` /* end */ );",
			expected: "KW (\n    /* start */ `hello `\n    + /* mid */ `world` /* end */\n);",
		},
	}

	var suiteCases []printerTestCase
	for _, kw := range []string{"error", "warning"} {
		for _, tc := range testCases {
			if tc.kwOnly != "" && tc.kwOnly != kw {
				continue
			}
			name := kw + "-" + tc.name
			code := strings.ReplaceAll(tc.code, "KW", kw)
			expected := strings.ReplaceAll(tc.expected, "KW", kw)
			masked := MaskReplaceAll(tc.masked, "KW", kw)

			suiteCases = append(suiteCases, printerTestCase{
				name:     name,
				code:     code,
				expected: expected,
				masked:   masked,
				lineCap:  lineCap,
				pc:       func(f func()) { f() },
			})
		}
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, suiteCases, lineCap)
}
