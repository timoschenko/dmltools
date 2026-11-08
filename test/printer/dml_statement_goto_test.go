// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestGotoStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-goto",
			code:     "goto   target   ;",
			expected: "goto target;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-messy-whitespace",
			code:     "    goto     my_label    ;\n",
			expected: "goto my_label;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-line-wrapping-long-identifier",
			code:     "goto    very_long_label_name_exceeding_forty_characters;",
			expected: "goto very_long_label_name_exceeding_forty_characters;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-sequential-gotos",
			code:     "goto   step1   ;\ngoto   step2   ;",
			expected: "goto step1;\ngoto step2;",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestGotoStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-block-before",
			code:     "/* c */   goto   target;",
			expected: "/* c */ goto target;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-block-between-goto-and-ident",
			code:     "goto   /* c */   target;",
			expected: "goto /* c */ target;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-block-before-semi",
			code:     "goto   target   /* c */  ;",
			expected: "goto target /* c */;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-block-after-semi",
			code:     "goto target; /* c */",
			expected: "goto target; /* c */",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-line-before-goto",
			code:     "// c\n   goto   target;",
			expected: "// c\ngoto target;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-line-between-goto-and-ident",
			code:     "goto   // c\n   target;",
			expected: "goto // c\n    target;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-line-before-semi",
			code:     "goto   target   // c\n;",
			expected: "goto target // c\n;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-line-after-semi",
			code:     "goto target;   // c",
			expected: "goto target; // c",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
