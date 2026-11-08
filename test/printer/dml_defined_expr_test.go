// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintDefinedExpr(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		// Suite A: Native Parentheses (Active Tests)
		{
			name:     "A01-basic",
			code:     "param TEST =     defined    (x  );\n",
			expected: "param TEST = defined(x);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A02-member-access",
			code:     "param TEST = defined(this.prop);\n",
			expected: "param TEST = defined(this.prop);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A03-array-access",
			code:     "param TEST = defined(my_array[0]);\n",
			expected: "param TEST = defined(my_array[0]);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A04-long-expr",
			code:     "param TEST = defined(this_is_a_very_long_variable_name_that_forces_wrap);\n",
			expected: "param TEST = defined(\n    this_is_a_very_long_variable_name_that_forces_wrap\n);\n",
			pc:       func(f func()) { f() },
		},

		// Suite B: Missing Parentheses (Skipped Tests)
		{
			masked:   MaskSyntheticSizeofParens("param TEST = defined x;\n"),
			name:     "B01-basic",
			code:     "param TEST =     defined    x  ;\n",
			expected: "param TEST = defined(x);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticSizeofParens(
				"param TEST = defined this.prop;\n",
			),
			name:     "B02-member-access",
			code:     "param TEST = defined this.prop;\n",
			expected: "param TEST = defined(this.prop);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticSizeofParens(
				"param TEST = defined my_array[0];\n",
			),
			name:     "B03-array-access",
			code:     "param TEST = defined my_array[0];\n",
			expected: "param TEST = defined(my_array[0]);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticSizeofParens(
				"param TEST = defined this_is_a_very_long_variable_name_that_forces_wrap;\n",
			),
			name:     "B04-long-expr",
			code:     "param TEST = defined this_is_a_very_long_variable_name_that_forces_wrap;\n",
			expected: "param TEST = defined(\n    this_is_a_very_long_variable_name_that_forces_wrap\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintDefinedExpr_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		// Suite A: Native Parentheses Comments
		{
			name:     "A01-comment-before-defined",
			code:     "param TEST = /* comment */ defined(x);\n",
			expected: "param TEST = /* comment */ defined(x);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A02-comment-between-defined-and-lparen",
			code:     "param TEST = defined /* comment */ (x);\n",
			expected: "param TEST = defined /* comment */ (x);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A03-comment-between-lparen-and-x",
			code:     "param TEST = defined( /* comment */ x);\n",
			expected: "param TEST = defined(/* comment */ x);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A04-comment-between-x-and-rparen",
			code:     "param TEST = defined(x /* comment */);\n",
			expected: "param TEST = defined(x /* comment */);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A05-comment-after-rparen",
			code:     "param TEST = defined(x) /* comment */;\n",
			expected: "param TEST = defined(x) /* comment */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A06-long-comment-forces-wrap",
			code:     "param TEST = defined( /* this is a somewhat long comment */ x);\n",
			expected: "param TEST = defined(\n    /* this is a somewhat long comment */ x\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A07-wrap-comment-before-defined",
			code:     "param TEST = /* this is a very long block comment that forces wrapping */ defined(x);\n",
			expected: "param TEST = /* this is a very long block comment that forces wrapping */ defined(\n    x\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A08-wrap-comment-between-defined-and-lparen",
			code:     "param TEST = defined /* this is a very long block comment that forces wrapping */ (x);\n",
			expected: "param TEST = defined /* this is a very long block comment that forces wrapping */ (\n    x\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A09-wrap-comment-between-lparen-and-x",
			code:     "param TEST = defined( /* this is a very long block comment that forces wrapping */ x);\n",
			expected: "param TEST = defined(\n    /* this is a very long block comment that forces wrapping */ x\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A10-wrap-comment-between-x-and-rparen",
			code:     "param TEST = defined(x /* this is a very long block comment that forces wrapping */);\n",
			expected: "param TEST = defined(\n    x /* this is a very long block comment that forces wrapping */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "A11-wrap-comment-after-rparen",
			code:     "param TEST = defined(x) /* this is a very long block comment that forces wrapping */;\n",
			expected: "param TEST = defined(\n    x\n) /* this is a very long block comment that forces wrapping */;\n",
			pc:       func(f func()) { f() },
		},

		// Suite B: Missing Parentheses Comments (Skipped)
		{
			masked: MaskSyntheticSizeofParens(
				"param TEST = /* comment */ defined x;\n",
			),
			name:     "B01-comment-before-defined",
			code:     "param TEST = /* comment */ defined x;\n",
			expected: "param TEST = /* comment */ defined(x);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticSizeofParens(
				"param TEST = defined /* comment */ x;\n",
			),
			name:     "B02-comment-between-defined-and-x",
			code:     "param TEST = defined /* comment */ x;\n",
			expected: "param TEST = defined(/* comment */ x);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticSizeofParens(
				"param TEST = defined x /* comment */;\n",
			),
			name:     "B03-comment-after-x",
			code:     "param TEST = defined x /* comment */;\n",
			expected: "param TEST = defined(x) /* comment */;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
