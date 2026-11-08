// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestAssertStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-assert-variable",
			code:     "assert   x   ;",
			expected: "assert x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-assert-condition",
			code:     "assert    x  >  0;",
			expected: "assert x > 0;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-basic-assert-property",
			code:     "assert   obj  .  prop;",
			expected: "assert obj.prop;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-line-wrapping-long-variable",
			code:     "assert    very_long_variable_name_exceeding_forty_characters;",
			expected: "assert very_long_variable_name_exceeding_forty_characters;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-line-wrapping-parenthesized-expression",
			code:     "assert (very_long_variable_name_a + very_long_variable_name_b);",
			expected: "assert (\n    very_long_variable_name_a\n    + very_long_variable_name_b\n);",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAssertStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-comment-alignment-inline-block",
			code:     "assert   /* block */   x; // end",
			expected: "assert /* block */ x; // end",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-alignment-line-comments",
			code:     "assert // ret\nx // init\n; // semi",
			expected: "assert // ret\n    x // init\n; // semi",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-alignment-block-before-assert",
			code:     "/* comment */    assert   x;",
			expected: "/* comment */ assert x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-alignment-block-between-assert-and-expr",
			code:     "assert   /* comment */   x;",
			expected: "assert /* comment */ x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-alignment-block-before-semi",
			code:     "assert   x   /* comment */  ;",
			expected: "assert x /* comment */;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comment-alignment-block-after-semi",
			code:     "assert x; /* comment */",
			expected: "assert x; /* comment */",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-comment-alignment-line-before-assert",
			code:     "// comment\n   assert    x;",
			expected: "// comment\nassert x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-comment-alignment-line-between-assert-and-expr",
			code:     "assert   // comment\n  x;",
			expected: "assert // comment\n    x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-comment-alignment-line-before-semi",
			code:     "assert   x   // comment\n;",
			expected: "assert x // comment\n;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-comment-alignment-line-after-semi",
			code:     "assert    x  ;   // comment",
			expected: "assert x; // comment",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
