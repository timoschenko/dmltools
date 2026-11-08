// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestMixedComments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		// ==========================================
		// Category 1: Issue #75 Motivating Snippets
		// ==========================================
		{
			name: "01-issue75-snippet1-banner-multiline",
			code: "// blabla\n\n" +
				"dml 1.4;\n\n" +
				"/********************************\n" +
				"//Some comments which has a syntax error on this line\n" +
				"//rest of the comment\n" +
				"//more info\n" +
				"//\n" +
				"********************************/\n",
			expected: "// blabla\n" +
				"dml 1.4;\n" +
				"/********************************\n" +
				"//Some comments which has a syntax error on this line\n" +
				"//rest of the comment\n" +
				"//more info\n" +
				"//\n" +
				"********************************/\n",
			pc: func(f func()) { f() },
		},
		{
			name: "02-issue75-snippet2-license-header",
			code: "/* -*- coding: utf-8 -*-\n\n" +
				"  multiline license text\n" +
				"  here\n  \n" +
				"*/\n\n" +
				"dml 1.4;\n",
			expected: "/* -*- coding: utf-8 -*-\n\n" +
				"  multiline license text\n" +
				"  here\n\n" +
				"*/\n" +
				"dml 1.4;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "03-issue75-snippet2-license-header-even-asterisks",
			code: "/* -*- coding: utf-8 -*-\n\n" +
				"  multiline license text\n" +
				"  here\n  \n" +
				"**/\n\n" +
				"dml 1.4;\n",
			expected: "/* -*- coding: utf-8 -*-\n\n" +
				"  multiline license text\n" +
				"  here\n\n" +
				"**/\n" +
				"dml 1.4;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "04-issue75-snippet3-attribute-6-asterisks",
			code: "dml 1.4;\n\n" +
				"/******Some text\"******/\n" +
				"attribute foo is (bool_attr) { }\n",
			expected: "dml 1.4;\n\n\n" +
				"/******Some text\"******/\n" +
				"attribute foo is (bool_attr) {\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "05-issue75-snippet3-attribute-4-asterisks",
			code: "dml 1.4;\n\n" +
				"/****Some text\"****/\n" +
				"attribute foo is (bool_attr) { }\n",
			expected: "dml 1.4;\n\n\n" +
				"/****Some text\"****/\n" +
				"attribute foo is (bool_attr) {\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "06-issue75-snippet3-attribute-2-asterisks",
			code: "dml 1.4;\n\n" +
				"/**Some text\"**/\n" +
				"attribute foo is (bool_attr) { }\n",
			expected: "dml 1.4;\n\n\n" +
				"/**Some text\"**/\n" +
				"attribute foo is (bool_attr) {\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "07-issue75-snippet3-attribute-5-asterisks-clean",
			code: "dml 1.4;\n\n" +
				"/*****Some text\"*****/\n" +
				"attribute foo is (bool_attr) { }\n",
			expected: "dml 1.4;\n\n\n" +
				"/*****Some text\"*****/\n" +
				"attribute foo is (bool_attr) {\n" +
				"}\n",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 2: Single-line Block Comment Asterisk Boundaries
		// ==========================================
		{
			name:     "10-singleline-0-asterisks",
			code:     "/* hello */\n",
			expected: "/* hello */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-singleline-1-asterisk-open-1-close",
			code:     "/** hello */\n",
			expected: "/** hello */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-singleline-2-asterisks-close",
			code:     "/* hello **/\n",
			expected: "/* hello **/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-singleline-3-asterisks-close",
			code:     "/* hello ***/\n",
			expected: "/* hello ***/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-singleline-4-asterisks-close",
			code:     "/* hello ****/\n",
			expected: "/* hello ****/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-singleline-5-asterisks-close",
			code:     "/* hello *****/\n",
			expected: "/* hello *****/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-singleline-internal-quotes-and-asterisks",
			code:     "/* \"hello\" *** 'world' */\n",
			expected: "/* \"hello\" *** 'world' */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-singleline-internal-quotes-2-asterisks-close",
			code:     "/* \"hello\" **/\n",
			expected: "/* \"hello\" **/\n",
			pc:       func(f func()) { f() },
		},

		// ==========================================
		// Category 3: Multi-line Block Comment Boundaries & Embedded Comments
		// ==========================================
		{
			name:     "20-multiline-clean-asterisks",
			code:     "/*\n * banner\n */\n",
			expected: "/*\n * banner\n */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "21-multiline-2-asterisks-close",
			code:     "/*\n * banner\n **/\n",
			expected: "/*\n * banner\n **/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "22-multiline-3-asterisks-close",
			code:     "/*\n * banner\n ***/\n",
			expected: "/*\n * banner\n ***/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "23-multiline-4-asterisks-close",
			code:     "/*\n * banner\n ****/\n",
			expected: "/*\n * banner\n ****/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "24-multiline-embedded-line-comments",
			code:     "/*\n// line comment\n*/\n",
			expected: "/*\n// line comment\n*/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "25-multiline-embedded-line-comments-2-asterisks-close",
			code:     "/*\n// line comment\n**/\n",
			expected: "/*\n// line comment\n**/\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "26-multiline-embedded-line-comments-banner",
			code: "/*\n" +
				"////////////////////////////////\n" +
				"// header\n" +
				"////////////////////////////////\n" +
				"*/\n",
			expected: "/*\n" +
				"////////////////////////////////\n" +
				"// header\n" +
				"////////////////////////////////\n" +
				"*/\n",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 4: Method Body and Statement Contexts
		// ==========================================
		{
			name: "30-method-banner-comment-even-asterisks",
			code: "method foo() {\n" +
				"    /********************************\n" +
				"    // internal method comment\n" +
				"    ********************************/\n" +
				"    bar();\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    /********************************\n" +
				"    // internal method comment\n" +
				"    ********************************/\n" +
				"    bar();\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "31-method-banner-comment-odd-asterisks",
			code: "method foo() {\n" +
				"    /*********************************\n" +
				"    // internal method comment\n" +
				"    *********************************/\n" +
				"    bar();\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    /*********************************\n" +
				"    // internal method comment\n" +
				"    *********************************/\n" +
				"    bar();\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "32-multiple-block-comments-no-code-swallowing",
			code: "/* comment 1 **/\n" +
				"dml 1.4;\n\n\n" +
				"device dev;\n" +
				"/* comment 2 */\n",
			expected: "/* comment 1 **/\n" +
				"dml 1.4;\n\n\n" +
				"device dev;\n" +
				"/* comment 2 */\n",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 5: Regex Content in Comments
		// ==========================================
		{
			name:     "40-regex-char-class",
			code:     "/* [a-zA-Z0-9_]+ */\n",
			expected: "/* [a-zA-Z0-9_]+ */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "41-regex-slashes-and-flags",
			code:     "/* /^[a-z]+/[gim]* */\n",
			expected: "/* /^[a-z]+/[gim]* */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "42-regex-escapes",
			code:     "/* \\d+\\.\\d+([eE][+-]?\\d+)? */\n",
			expected: "/* \\d+\\.\\d+([eE][+-]?\\d+)? */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "43-regex-negated-class",
			code:     "/* [^/\\*]+ */\n",
			expected: "/* [^/\\*]+ */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "44-regex-lookaround",
			code:     "/* (?<=foo)(?<!bar)(?:baz|qux)+ */\n",
			expected: "/* (?<=foo)(?<!bar)(?:baz|qux)+ */\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "45-regex-full-complex",
			code: "/* ^(?:[a-zA-Z_]\\w*|\\d+(?:\\.\\d+)?|" +
				"/.*?/|[!@#$%^&*()-+=[\\]{}|;:,.<>?/~]+)$ */\n",
			expected: "/* ^(?:[a-zA-Z_]\\w*|\\d+(?:\\.\\d+)?|" +
				"/.*?/|[!@#$%^&*()-+=[\\]{}|;:,.<>?/~]+)$ */\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "46-regex-line-comment",
			code:     "// ^[a-zA-Z_][a-zA-Z0-9_]*$\n",
			expected: "// ^[a-zA-Z_][a-zA-Z0-9_]*$\n",
			pc:       func(f func()) { f() },
		},

		// ==========================================
		// Category 6: Escape Sequences in Comments
		// ==========================================
		{
			name:     "50-escape-newlines-and-tabs",
			code:     "/* line1\\nline2\\tline3 */\n",
			expected: "/* line1\\nline2\\tline3 */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "51-escape-backslashes",
			code:     "/* \\ \\\\ \\\\\\ \\\\\\\\ */\n",
			expected: "/* \\ \\\\ \\\\\\ \\\\\\\\ */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "52-escape-quotes",
			code:     "/* \\\" \\' */\n",
			expected: "/* \\\" \\' */\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "53-escape-at-eol",
			code: "/* escape at line end \\\n" +
				"next line */\n",
			expected: "/* escape at line end \\\n" +
				"next line */\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "54-escape-slash-star",
			code:     "/* \\/* not opening \\*/\n",
			expected: "/* \\/* not opening \\*/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "55-escape-hex-unicode",
			code:     "/* \\x1b[31m \\u001b \\U0000001b */\n",
			expected: "/* \\x1b[31m \\u001b \\U0000001b */\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "56-escape-line-comment-continuation",
			code: "// backslash at end of line \\\n" +
				"// next line\n",
			expected: "// backslash at end of line \\\n" +
				"// next line\n",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 7: Quotes in Comments
		// ==========================================
		{
			name:     "60-quotes-unmatched-double",
			code:     "/* \" */\n",
			expected: "/* \" */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "61-quotes-unmatched-single",
			code:     "/* ' */\n",
			expected: "/* ' */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "62-quotes-mixed",
			code:     "/* \" ' \" ' */\n",
			expected: "/* \" ' \" ' */\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "63-quotes-in-multiline",
			code: "/*\n" +
				"\"hello world\"\n" +
				"'c'\n" +
				"*/\n",
			expected: "/*\n" +
				"\"hello world\"\n" +
				"'c'\n" +
				"*/\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "64-quotes-triple",
			code:     "/* \"\"\" ''' */\n",
			expected: "/* \"\"\" ''' */\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "65-quotes-with-code",
			code: "/* \"unclosed string in comment */\n" +
				"attribute foo is (bool_attr) {}\n",
			expected: "/* \"unclosed string in comment */\n" +
				"attribute foo is (bool_attr) {\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "66-quotes-single-with-code",
			code: "/* 'unclosed char in comment */\n" +
				"attribute foo is (bool_attr) {}\n",
			expected: "/* 'unclosed char in comment */\n" +
				"attribute foo is (bool_attr) {\n" +
				"}\n",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 8: Delimiter Edge Cases, Adjacency & Nesting
		// ==========================================
		{
			name:     "70-nesting-open-inside-open",
			code:     "/* /* */\n",
			expected: "/* /* */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "71-nesting-pseudo-nested",
			code:     "/* /* nested? */\n",
			expected: "/* /* nested? */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "72-multiple-nested-openers",
			code:     "/* /* /* /* */\n",
			expected: "/* /* /* /* */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "73-adjacent-comments-no-space",
			code:     "/***//**/\n",
			expected: "/***/ /**/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "74-adjacent-multiple-slashes",
			code:     "/* a *//* b */\n",
			expected: "/* a */ /* b */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "75-separated-slashes-and-stars",
			code:     "/* * / * / */\n",
			expected: "/* * / * / */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "76-line-comment-containing-block",
			code:     "// /* not a block comment */\n",
			expected: "// /* not a block comment */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "77-block-comment-containing-line-comment-chars",
			code:     "/* // // // */\n",
			expected: "/* // // // */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "78-slashes-and-stars-run",
			code:     "/* //*//** */\n",
			expected: "/* //*/ /** */\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "82-method-inline-complex-comments",
			code: "method foo() {\n" +
				"    /* /^[0-9]+$/ */\n" +
				"    /* \"unclosed quote */\n" +
				"    /***//**/\n" +
				"    local int x = 1;\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    /* /^[0-9]+$/ */\n" +
				"    /* \"unclosed quote */\n" +
				"    /***/ /**/\n" +
				"    local int x = 1;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "83-nested-comment-ignored-like-c-valid-code-follows",
			code: "/* comment /* still comment */\n" +
				"attribute a is (bool_attr) {}\n",
			expected: "/* comment /* still comment */\n" +
				"attribute a is (bool_attr) {\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "85-line-comment-at-eof-no-newline",
			code: "attribute a is (bool_attr) {}\n" +
				"// trailing comment at EOF",
			expected: "attribute a is (bool_attr) {\n" +
				"}\n" +
				"// trailing comment at EOF",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 80)
}
