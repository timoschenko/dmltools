// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestReturnStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-return",
			code:     "   return    `John`   ;",
			expected: "return `John`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-empty-return",
			code:     "    return    ;",
			expected: "return;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-line-wrapping-long-function-call",
			code:     "return ThisIsAVeryLongFunctionNameCallThatExceedsFortyCharacters(arg1, arg2);",
			expected: "return ThisIsAVeryLongFunctionNameCallThatExceedsFortyCharacters(\n    arg1, arg2\n);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-line-wrapping-long-function-call-wrap-args",
			code:     "return this_size_to_force_wrap_args(arg1, arg2);",
			expected: "return this_size_to_force_wrap_args(\n    arg1, arg2\n);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-line-wrapping-long-struct-initializer",
			code:     "return {\n       .very_long_property_name   =   `value1`,\n    .another_long_property = `value2`\n};",
			expected: "return {\n    .very_long_property_name = `value1`,\n    .another_long_property = `value2`\n};",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-parenthesized-expression-no-wrap",
			code:     "return (a + b);",
			expected: "return (a + b);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-parenthesized-expression-first-level-wrap",
			code:     "return (not_long_name_a + not_long_name_b);",
			expected: "return (\n    not_long_name_a + not_long_name_b\n);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-parenthesized-expression-full-wrap",
			code:     "return (very_long_variable_name_a + very_long_variable_name_b);",
			expected: "return (\n    very_long_variable_name_a\n    + very_long_variable_name_b\n);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-tuple-return-no-wrap",
			code:     "return    (  a  ,   b  ,   c  );",
			expected: "return (a, b, c);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-tuple-return-wrap",
			code:     "return (very_long_variable_name_a, very_long_variable_name_b, very_long_variable_name_c);",
			expected: "return (\n    very_long_variable_name_a,\n    very_long_variable_name_b,\n    very_long_variable_name_c\n);",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestReturnStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-comment-alignment-inline-block",
			code:     "return /* block */ 123; // end",
			expected: "return /* block */ 123; // end",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-alignment-line-comments",
			code:     "return // ret\n123 // init\n; // semi",
			expected: "return // ret\n    123 // init\n; // semi",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-alignment-block-before-return",
			code:     "/* comment */ return `John`;",
			expected: "/* comment */ return `John`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-alignment-block-between-return-and-expr",
			code:     "return /* comment */ `John`;",
			expected: "return /* comment */ `John`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-alignment-block-before-semi",
			code:     "return `John` /* comment */;",
			expected: "return `John` /* comment */;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comment-alignment-block-after-semi",
			code:     "return `John`; /* comment */",
			expected: "return `John`; /* comment */",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-comment-alignment-multiple-block",
			code:     "/* 1 */ return /* 2 */ `John` /* 3 */ ; /* 4 */",
			expected: "/* 1 */ return /* 2 */ `John` /* 3 */; /* 4 */",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-comment-alignment-line-before-return",
			code:     "// comment\nreturn `John`;",
			expected: "// comment\nreturn `John`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-comment-alignment-line-between-return-and-expr",
			code:     "return // comment\n`John`;",
			expected: "return // comment\n    `John`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-comment-alignment-line-before-semi",
			code:     "return `John` // comment\n;",
			expected: "return `John` // comment\n;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-comment-alignment-line-after-semi",
			code:     "return `John`; // comment",
			expected: "return `John`; // comment",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-comment-alignment-multiple-line",
			code:     "// 1\nreturn // 2\n`John` // 3\n; // 4",
			expected: "// 1\nreturn // 2\n    `John` // 3\n; // 4",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
