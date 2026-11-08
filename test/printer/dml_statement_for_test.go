// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestForLoopStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-for",
			code:     "for   (  i  =  0  ;   i  <  10  ;   i  ++  )   {\n  x  =  1;\n}",
			expected: "for (i = 0; i < 10; i++) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-empty-for",
			code:     "for   (  ;  ;  )   {\n  x  =  1;\n}",
			expected: "for (;;) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-for-decl-init",
			code:     "for (local int i = 0; i < 10; i++) {\n  x = 1;\n}",
			expected: "for (local int i = 0; i < 10; i++) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-for-multiple-init-and-step",
			code:     "for (i = 0, j = 10; i < j; i++, j--) {\n  x = 1;\n}",
			expected: "for (i = 0, j = 10; i < j; i++, j--) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-for-assignop",
			code:     "for (i = 0; i < 10; i += 2) {\n  x = 1;\n}",
			expected: "for (i = 0; i < 10; i += 2) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-for-no-init",
			code:     "for (; i < 10; i++) {\n  x = 1;\n}",
			expected: "for (; i < 10; i++) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-for-no-step",
			code:     "for (i = 0; i < 10;) {\n  x = 1;\n}",
			expected: "for (i = 0; i < 10;) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-for-no-cond",
			code:     "for (i = 0;; i++) {\n  x = 1;\n}",
			expected: "for (i = 0;; i++) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-unbraced-for",
			code:     "for (i = 0; i < 10; i++)\n  x = 1;",
			expected: "for (i = 0; i < 10; i++) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "10-nested-for",
			code: "for (i = 0; i < 10; i++) {\n  for (j = 0; j < 5; j++) {\n    x = 1;\n  }\n}",
			expected: "for (i = 0; i < 10; i++) {\n" +
				"    for (j = 0; j < 5; j++) {\n" +
				"        x = 1;\n" +
				"    }\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name:     "11-prefix-step",
			code:     "for   (  local   int   i  =  0  ;   i  <  10  ;   ++  i  )   {\n  x  +=  i;\n}",
			expected: "for (local int i = 0; i < 10; ++i) {\n    x += i;\n}",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestForLoopStatement_ConditionWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-level2-wrapping-condition",
			code: "for (i = 0; very_long_condition_name_a && very_long_condition_name_b; i++) {\n  x = 1;\n}",
			expected: "for (\n" +
				"    i = 0;\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b;\n" +
				"    i++\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "02-line-wrapping-long-header",
			code: "for (local int index = 0; index < limit_value; ++index) {\n    x += index;\n}",
			expected: "for (\n" +
				"    local int index = 0;\n" +
				"    index < limit_value;\n" +
				"    ++index\n" +
				") {\n" +
				"    x += index;\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestForLoopStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-leading-comment",
			code:     "// c\nfor (i = 0; i < 10; i++) {\n  x = 1;\n}",
			expected: "// c\nfor (i = 0; i < 10; i++) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-between-for-and-paren",
			code:     "for /* c */ (i = 0; i < 10; i++) {\n  x = 1;\n}",
			expected: "for /* c */ (i = 0; i < 10; i++) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-before-first-semi",
			code:     "for (i = 0 /* c */; i < 10; i++) {\n  x = 1;\n}",
			expected: "for (i = 0 /* c */; i < 10; i++) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-after-first-semi",
			code:     "for (i = 0; /* c */ i < 10; i++) {\n  x = 1;\n}",
			expected: "for (i = 0; /* c */ i < 10; i++) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-after-second-semi",
			code:     "for (i = 0; i < 10; /* c */ ++i) {\n  x += i;\n}",
			expected: "for (i = 0; i < 10; /* c */ ++i) {\n    x += i;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comment-between-paren-and-brace",
			code:     "for (i = 0; i < 10; i++) /* c */ {\n  x = 1;\n}",
			expected: "for (i = 0; i < 10; i++) /* c */ {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-trailing-comment",
			code:     "for (i = 0; i < 10; i++) {\n  x = 1;\n} // done",
			expected: "for (i = 0; i < 10; i++) {\n    x = 1;\n} // done",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
