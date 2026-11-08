// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestSimpleControlStatements(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-break",
			code:     "break;",
			expected: "break;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-continue",
			code:     "continue;",
			expected: "continue;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-basic-throw",
			code:     "throw;",
			expected: "throw;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-whitespace-break",
			code:     "   break   ;",
			expected: "break;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-whitespace-continue",
			code:     "   continue   ;",
			expected: "continue;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-whitespace-throw",
			code:     "   throw   ;",
			expected: "throw;",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestSimpleControlStatements_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-leading-comment-break",
			code:     "// c\nbreak;",
			expected: "// c\nbreak;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-leading-comment-continue",
			code:     "// c\ncontinue;",
			expected: "// c\ncontinue;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-leading-comment-throw",
			code:     "// c\nthrow;",
			expected: "// c\nthrow;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-block-comment-before-semi-break",
			code:     "break /* c */ ;",
			expected: "break /* c */;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-block-comment-before-semi-continue",
			code:     "continue /* c */ ;",
			expected: "continue /* c */;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-block-comment-before-semi-throw",
			code:     "throw /* c */ ;",
			expected: "throw /* c */;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-line-comment-before-semi-break",
			code:     "break // c\n;",
			expected: "break // c\n;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-line-comment-before-semi-continue",
			code:     "continue // c\n;",
			expected: "continue // c\n;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-line-comment-before-semi-throw",
			code:     "throw // c\n;",
			expected: "throw // c\n;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-trailing-block-comment-break",
			code:     "break; /* c */",
			expected: "break; /* c */",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-trailing-block-comment-continue",
			code:     "continue; /* c */",
			expected: "continue; /* c */",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-trailing-block-comment-throw",
			code:     "throw; /* c */",
			expected: "throw; /* c */",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-trailing-line-comment-break",
			code:     "break; // c",
			expected: "break; // c",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-trailing-line-comment-continue",
			code:     "continue; // c",
			expected: "continue; // c",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-trailing-line-comment-throw",
			code:     "throw; // c",
			expected: "throw; // c",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
