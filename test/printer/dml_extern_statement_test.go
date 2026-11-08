// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintExternStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "extern int x;\n",
			expected: "extern int x;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-messy",
			code:     "  extern   int   x  ;  \n",
			expected: "extern int x;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-typedef",
			code:     "extern typedef void (*func)(int);\n",
			expected: "extern typedef void (*func)(int);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-typedef-messy",
			code:     " extern    typedef     void (*func)(int)  ;\n",
			expected: "extern typedef void (*func)(int);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-multiline-cdecl",
			code:     "extern typedef struct {\n    int a;\n} b;\n",
			expected: "extern typedef struct {\n    int a;\n} b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-wrap-long-name",
			code:     "extern int a_very_long_variable_name_that_exceeds_forty_characters;\n",
			expected: "extern int a_very_long_variable_name_that_exceeds_forty_characters;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-wrap-long-func-name",
			code:     "extern void a_very_long_function_name_exceeding_line_cap(int a, int b);\n",
			expected: "extern void a_very_long_function_name_exceeding_line_cap(\n    int a, int b\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-wrap-long-func-args",
			code:     "extern void f(int a_very_long_argument_name_1, int a_very_long_argument_name_2);\n",
			expected: "extern void f(\n    int a_very_long_argument_name_1,\n    int a_very_long_argument_name_2\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "09-extern-typedef-union-issue-77",
			code: "extern typedef union {\n" +
				"    struct {\n" +
				"        uint32 a   :8;\n" +
				"        uint32 b   :8;\n" +
				"        uint32 ca  :4;\n" +
				"        uint32 cb  :4;\n" +
				"        uint32 d   :8;\n" +
				"    } field;\n" +
				"    uint32 value;\n" +
				"} foo_t;\n",
			expected: "extern typedef union {\n" +
				"    struct {\n" +
				"        uint32 a : 8;\n" +
				"        uint32 b : 8;\n" +
				"        uint32 ca : 4;\n" +
				"        uint32 cb : 4;\n" +
				"        uint32 d : 8;\n" +
				"    } field;\n" +
				"    uint32 value;\n" +
				"} foo_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "10-extern-union-variable",
			code:     "extern   union   foo_u   bar  ;\n",
			expected: "extern union foo_u bar;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExternStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-comment-tail",
			code:     "extern int x; // comment\n",
			expected: "extern int x; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-inline",
			code:     "extern /* foo */ int x;\n",
			expected: "extern /* foo */ int x;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-typedef-inline",
			code:     "extern /* foo */ typedef int x;\n",
			expected: "extern /* foo */ typedef int x;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-scattered-basic",
			code:     "extern // c1\nint // c2\nx; // c3\n",
			expected: "extern // c1\nint // c2\nx; // c3\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-scattered-typedef",
			code:     "extern // c1\ntypedef // c2\nvoid // c3\n(*func)(int); // c4\n",
			expected: "extern // c1\ntypedef // c2\nvoid // c3\n(*func)(int); // c4\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comment-scattered-struct",
			code:     "extern // c1\ntypedef // c2\nstruct   {  // c3\n    int a; // c4\n} b; // c5\n",
			expected: "extern // c1\ntypedef // c2\nstruct { // c3\n    int a; // c4\n} b; // c5\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-comment-scattered-inline-basic",
			code:     "extern /* c1 */ int /* c2 */ x /* c3 */;\n",
			expected: "extern /* c1 */ int /* c2 */ x /* c3 */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-comment-scattered-inline-typedef",
			code:     "extern /* c1 */ typedef /* c2 */ void /* c3 */ (*func)(int) /* c4 */;\n",
			expected: "extern /* c1 */ typedef /* c2 */ void /* c3 */ (*func)(\n    int\n) /* c4 */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-comment-scattered-inline-struct",
			code:     "extern /* c1 */ typedef /* c2 */ struct { /* c3 */ int a; /* c4 */ } b /* c5 */;\n",
			expected: "extern /* c1 */ typedef /* c2 */ struct {\n    /* c3 */ int a; /* c4 */\n} b /* c5 */;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
