// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintHookDecl(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		// 1. Basic Definitions
		{
			name:     "01-basic-no-args",
			code:     "   hook (  ) myHook  ;\n",
			expected: "hook() myHook;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-with-args",
			code:     "hook   (   int a   ,   float b ) myHook;\n",
			expected: "hook(int a, float b) myHook;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-wrap-args",
			code:     "hook (int    this_is_a_very_very_long_arg_name_1, int this_is_a_very_very_long_arg_name_2)   myHook;\n",
			expected: "hook(\n    int this_is_a_very_very_long_arg_name_1,\n    int this_is_a_very_very_long_arg_name_2\n) myHook;\n",
			pc:       func(f func()) { f() },
		},
		// 3. Array Dimensions
		{
			name:     "05-array-dims-simple",
			code:     "hook () myHook [   2 ] [ 3 ] ;\n",
			expected: "hook() myHook[2][3];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-wrap-array-dims",
			code:     "hook (  ) myHook   [ this_is_a_very_long_expression_exceeding_cap  ]   [   3   ] ;\n",
			expected: "hook() myHook[\n    this_is_a_very_long_expression_exceeding_cap\n][3];\n",
			pc:       func(f func()) { f() },
		},
		// 4. Trailing Comma (Compiler workaround)
		{
			masked: MaskTask006Skip(
				"printHook_decl panics on trailing comma due to unsupported TokenComma in hook_decl",
			),
			name:     "07-trailing-comma",
			code:     "hook (int a,) myHook;\n",
			expected: "hook(int a,) myHook;\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintHookDecl_Comments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-comment-after-semi",
			code:     "hook    (  )    myHook  ;   // comment\n",
			expected: "hook() myHook; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-inline-args",
			code:     "hook   (  /* c1 */   int   a   /* c2 */   )   myHook;\n",
			expected: "hook(/* c1 */ int a /* c2 */) myHook;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-inline-array",
			code:     "hook   (  )   myHook    /* c1 */   [    /* c2 */    2    /* c3 */   ]    /* c4 */    ;\n",
			expected: "hook() myHook /* c1 */[\n    /* c2 */ 2 /* c3 */] /* c4 */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-line-comment-inside-empty-args",
			code:     "hook ( // c1\n) myHook;\n",
			expected: "hook( // c1\n) myHook;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-line-comment-before-arg",
			code:     "hook (    // c1\nint    a   )  myHook;\n",
			expected: "hook( // c1\n    int a\n) myHook;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-line-comment-after-arg",
			code:     "hook (int a // c1\n) myHook;\n",
			expected: "hook(\n    int a // c1\n) myHook;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-line-comment-multiple-args",
			code:     "hook (int a, // c1\nfloat b // c2\n) myHook;\n",
			expected: "hook(\n    int a, // c1\n    float b // c2\n) myHook;\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
