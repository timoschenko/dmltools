// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestTryStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-try-catch-compound",
			code:     "try   {\n  x  =  1;\n}   catch   {\n  x  =  0;\n}",
			expected: "try {\n    x = 1;\n} catch {\n    x = 0;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-single-statement-unbraced",
			code:     "try   x  =  1;   catch   x  =  0;",
			expected: "try {\n    x = 1;\n} catch {\n    x = 0;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "03-nested-try-catch",
			code: "try {\ntry {\nx = 1;\n} catch {\nthrow;\n}\n} catch {\nx = 0;\n}",
			expected: "try {\n" +
				"    try {\n" +
				"        x = 1;\n" +
				"    } catch {\n" +
				"        throw;\n" +
				"    }\n" +
				"} catch {\n" +
				"    x = 0;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name:     "04-try-with-throw",
			code:     "try {\nthrow;\n} catch {\nlog info: \"err\";\n}",
			expected: "try {\n    throw;\n} catch {\n    log info: \"err\";\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "05-line-wrapping-multiline-body",
			code: "try {\n" +
				"    x = (very_long_variable_name_a + very_long_variable_name_b);\n" +
				"} catch {\n" +
				"    y = (very_long_variable_name_a + very_long_variable_name_b);\n" +
				"}",
			expected: "try {\n" +
				"    x = (\n" +
				"        very_long_variable_name_a\n" +
				"        + very_long_variable_name_b\n" +
				"    );\n" +
				"} catch {\n" +
				"    y = (\n" +
				"        very_long_variable_name_a\n" +
				"        + very_long_variable_name_b\n" +
				"    );\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestTryStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-leading-comment",
			code:     "// comment\ntry {\nx = 1;\n} catch {\nx = 0;\n}",
			expected: "// comment\ntry {\n    x = 1;\n} catch {\n    x = 0;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-between-try-and-brace",
			code:     "try /* c */ {\nx = 1;\n} catch {\nx = 0;\n}",
			expected: "try /* c */ {\n    x = 1;\n} catch {\n    x = 0;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-inside-try",
			code:     "try {\n// in try\nx = 1;\n} catch {\nx = 0;\n}",
			expected: "try {\n    // in try\n    x = 1;\n} catch {\n    x = 0;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-between-brace-and-catch",
			code:     "try {\nx = 1;\n} /* c */ catch {\nx = 0;\n}",
			expected: "try {\n    x = 1;\n} /* c */ catch {\n    x = 0;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-between-catch-and-brace",
			code:     "try {\nx = 1;\n} catch /* c */ {\nx = 0;\n}",
			expected: "try {\n    x = 1;\n} catch /* c */ {\n    x = 0;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comment-inside-catch",
			code:     "try {\nx = 1;\n} catch {\n// in catch\nx = 0;\n}",
			expected: "try {\n    x = 1;\n} catch {\n    // in catch\n    x = 0;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-trailing-comment",
			code:     "try {\nx = 1;\n} catch {\nx = 0;\n} // done",
			expected: "try {\n    x = 1;\n} catch {\n    x = 0;\n} // done",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
