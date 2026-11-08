// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintStringifyExpr(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "param TEST =     stringify    (a  );\n",
			expected: "param TEST = stringify(a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-complex-expr",
			code:     "param TEST = stringify(a+b);\n",
			expected: "param TEST = stringify(a + b);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-long-expr",
			code:     "param TEST = stringify(this_is_a_very_long_variable_name_that_forces_wrap);\n",
			expected: "param TEST = stringify(\n    this_is_a_very_long_variable_name_that_forces_wrap\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-nested-stringify",
			code:     "param TEST = stringify(stringify(a));\n",
			expected: "param TEST = stringify(stringify(a));\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-nested-long-stringify",
			code:     "param TEST = stringify(stringify(this_is_a_very_long_variable_name_that_forces_wrap));\n",
			expected: "param TEST = stringify(\n    stringify(\n        this_is_a_very_long_variable_name_that_forces_wrap\n    )\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintStringifyExpr_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-with-comments",
			code:     "param TEST = stringify(a /* comment */);\n",
			expected: "param TEST = stringify(a /* comment */);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-first-arg-quite-long",
			code:     "param TEST = stringify(quite_long_arg_name // comment\n);\n",
			expected: "param TEST = stringify(\n    quite_long_arg_name // comment\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-first-arg-quite-long-block-comment",
			code:     "param TEST = stringify(quite_long_arg_name /* c1 */);\n",
			expected: "param TEST = stringify(\n    quite_long_arg_name /* c1 */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-line-comment-before-stringify",
			code:     "param TEST = // comment\nstringify(a);\n",
			expected: "param TEST = // comment\n    stringify(a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-block-comment-before-stringify",
			code:     "param TEST = /* comment */ stringify(a);\n",
			expected: "param TEST = /* comment */ stringify(a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-line-comment-after-semi",
			code:     "param TEST = stringify(a); // comment\n",
			expected: "param TEST = stringify(a); // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-block-comment-after-semi",
			code:     "param TEST = stringify(a); /* comment */\n",
			expected: "param TEST = stringify(a); /* comment */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-long-comment-forces-wrap",
			code:     "param TEST = stringify(a /* this is a somewhat long comment */);\n",
			expected: "param TEST = stringify(\n    a /* this is a somewhat long comment */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-very-long-comment-exceeds-cap",
			code:     "param TEST = stringify(a /* this is a very extremely super incredibly long comment that exceeds line cap */);\n",
			expected: "param TEST = stringify(\n    a /* this is a very extremely super incredibly long comment that exceeds line cap */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-wrap-comment-before-stringify",
			code:     "param TEST = /* this is a very long block comment that forces wrapping */ stringify(a);\n",
			expected: "param TEST = /* this is a very long block comment that forces wrapping */ stringify(\n    a\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-wrap-comment-between-stringify-and-lparen",
			code:     "param TEST = stringify /* this is a very long block comment that forces wrapping */ (a);\n",
			expected: "param TEST = stringify /* this is a very long block comment that forces wrapping */ (\n    a\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-wrap-comment-after-lparen",
			code:     "param TEST = stringify( /* this is a very long block comment that forces wrapping */ a);\n",
			expected: "param TEST = stringify(\n    /* this is a very long block comment that forces wrapping */ a\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-wrap-comment-before-rparen",
			code:     "param TEST = stringify(a /* this is a very long block comment that forces wrapping */);\n",
			expected: "param TEST = stringify(\n    a /* this is a very long block comment that forces wrapping */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-wrap-comment-after-rparen",
			code:     "param TEST = stringify(a) /* this is a very long block comment that forces wrapping */;\n",
			expected: "param TEST = stringify(\n    a\n) /* this is a very long block comment that forces wrapping */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-wrap-comment-after-semi",
			code:     "param TEST = stringify(a); /* this is a very long block comment that forces wrapping */\n",
			expected: "param TEST = stringify(\n    a\n); /* this is a very long block comment that forces wrapping */\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
