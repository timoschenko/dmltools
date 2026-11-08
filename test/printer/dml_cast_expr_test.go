// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintCastExpr(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "param TEST =     cast    (a,const    foo_t    *  );\n",
			expected: "param TEST = cast(a, const foo_t *);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-complex-expr",
			code:     "param TEST = cast(a+b,int);\n",
			expected: "param TEST = cast(a + b, int);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-long-expr",
			code:     "param TEST = cast(this_is_a_very_long_variable_name_that_forces_wrap, int);\n",
			expected: "param TEST = cast(\n    this_is_a_very_long_variable_name_that_forces_wrap,\n    int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-long-type",
			code:     "param TEST = cast(a, const struct_with_a_very_long_name_that_forces_wrap *);\n",
			expected: "param TEST = cast(\n    a,\n    const struct_with_a_very_long_name_that_forces_wrap *\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-nested-cast",
			code:     "param TEST = cast(cast(a, int), float);\n",
			expected: "param TEST = cast(cast(a, int), float);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-nested-long-cast",
			code:     "param TEST = cast(cast(this_is_a_very_long_variable_name_that_forces_wrap, int), float);\n",
			expected: "param TEST = cast(\n    cast(\n        this_is_a_very_long_variable_name_that_forces_wrap,\n        int\n    ),\n    float\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintCastExpr_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-with-comments",
			code:     "param TEST = cast(a, /* comment */ int);\n",
			expected: "param TEST = cast(a, /* comment */ int);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-first-arg-quite-long",
			code:     "param TEST = cast(quite_long_arg_name, // comment\n int);\n",
			expected: "param TEST = cast(\n    quite_long_arg_name, // comment\n    int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-last-arg-quite-long",
			code:     "param TEST = cast(a, quite_long_type_name // comment\n);\n",
			expected: "param TEST = cast(\n    a, quite_long_type_name // comment\n);\n",
			masked: MaskArgumentOverWrappingWithLineComments(
				"param TEST = cast(\n    a,\n    quite_long_type_name // comment\n);\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "04-first-arg-quite-long-block-comment",
			code:     "param TEST = cast(quite_long_arg_name /* c1 */, int);\n",
			expected: "param TEST = cast(\n    quite_long_arg_name /* c1 */, int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-line-comment-before-cast",
			code:     "param TEST = // comment\ncast(a, int);\n",
			expected: "param TEST = // comment\n    cast(a, int);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-block-comment-before-cast",
			code:     "param TEST = /* comment */ cast(a, int);\n",
			expected: "param TEST = /* comment */ cast(a, int);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-line-comment-after-semi",
			code:     "param TEST = cast(a, int); // comment\n",
			expected: "param TEST = cast(a, int); // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-block-comment-after-semi",
			code:     "param TEST = cast(a, int); /* comment */\n",
			expected: "param TEST = cast(a, int); /* comment */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-long-comment-forces-wrap",
			code:     "param TEST = cast(a, /* this is a somewhat long comment */ int);\n",
			expected: "param TEST = cast(\n    a,\n    /* this is a somewhat long comment */ int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-very-long-comment-exceeds-cap",
			code:     "param TEST = cast(a, /* this is a very extremely super incredibly long comment that exceeds line cap */ int);\n",
			expected: "param TEST = cast(\n    a,\n    /* this is a very extremely super incredibly long comment that exceeds line cap */ int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-wrap-comment-before-cast",
			code:     "param TEST = /* this is a very long block comment that forces wrapping */ cast(a, int);\n",
			expected: "param TEST = /* this is a very long block comment that forces wrapping */ cast(\n    a, int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-wrap-comment-between-cast-and-lparen",
			code:     "param TEST = cast /* this is a very long block comment that forces wrapping */ (a, int);\n",
			expected: "param TEST = cast /* this is a very long block comment that forces wrapping */ (\n    a, int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-wrap-comment-after-lparen",
			code:     "param TEST = cast( /* this is a very long block comment that forces wrapping */ a, int);\n",
			expected: "param TEST = cast(\n    /* this is a very long block comment that forces wrapping */ a,\n    int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-wrap-comment-before-comma",
			code:     "param TEST = cast(a /* this is a very long block comment that forces wrapping */, int);\n",
			expected: "param TEST = cast(\n    a /* this is a very long block comment that forces wrapping */,\n    int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-wrap-comment-after-comma",
			code:     "param TEST = cast(a, /* this is a very long block comment that forces wrapping */ int);\n",
			expected: "param TEST = cast(\n    a,\n    /* this is a very long block comment that forces wrapping */ int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-wrap-comment-before-rparen",
			code:     "param TEST = cast(a, int /* this is a very long block comment that forces wrapping */);\n",
			expected: "param TEST = cast(\n    a,\n    int /* this is a very long block comment that forces wrapping */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-wrap-comment-after-rparen",
			code:     "param TEST = cast(a, int) /* this is a very long block comment that forces wrapping */;\n",
			expected: "param TEST = cast(\n    a, int\n) /* this is a very long block comment that forces wrapping */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-wrap-comment-after-semi",
			code:     "param TEST = cast(a, int); /* this is a very long block comment that forces wrapping */\n",
			expected: "param TEST = cast(\n    a, int\n); /* this is a very long block comment that forces wrapping */\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
