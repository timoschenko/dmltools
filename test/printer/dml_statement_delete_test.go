// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestDeleteStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-delete-variable",
			code:     "delete   x   ;",
			expected: "delete x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-delete-property",
			code:     "delete    obj  .  prop;",
			expected: "delete obj.prop;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-basic-delete-array-index",
			code:     "delete   arr [  i  ];",
			expected: "delete arr[i];",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-line-wrapping-long-variable",
			code:     "delete    very_long_variable_name_exceeding_forty_characters;",
			expected: "delete very_long_variable_name_exceeding_forty_characters;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-line-wrapping-function-call",
			code:     "delete this_size_to_force_wrap_args(arg1, arg2);",
			expected: "delete this_size_to_force_wrap_args(\n    arg1, arg2\n);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-line-wrapping-parenthesized-expression",
			code:     "delete (very_long_variable_name_a + very_long_variable_name_b);",
			expected: "delete (\n    very_long_variable_name_a\n    + very_long_variable_name_b\n);",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestDeleteStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-comment-alignment-inline-block",
			code:     "delete   /* block */   x; // end",
			expected: "delete /* block */ x; // end",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-alignment-line-comments",
			code:     "delete // ret\nx // init\n; // semi",
			expected: "delete // ret\n    x // init\n; // semi",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-alignment-block-before-delete",
			code:     "/* comment */    delete   x;",
			expected: "/* comment */ delete x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-alignment-block-between-delete-and-expr",
			code:     "delete   /* comment */   x;",
			expected: "delete /* comment */ x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-alignment-block-before-semi",
			code:     "delete   x   /* comment */  ;",
			expected: "delete x /* comment */;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comment-alignment-block-after-semi",
			code:     "delete x; /* comment */",
			expected: "delete x; /* comment */",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-comment-alignment-line-before-delete",
			code:     "// comment\n   delete    x;",
			expected: "// comment\ndelete x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-comment-alignment-line-between-delete-and-expr",
			code:     "delete   // comment\n  x;",
			expected: "delete // comment\n    x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-comment-alignment-line-before-semi",
			code:     "delete   x   // comment\n;",
			expected: "delete x // comment\n;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-comment-alignment-line-after-semi",
			code:     "delete    x  ;   // comment",
			expected: "delete x; // comment",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
