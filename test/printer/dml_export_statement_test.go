// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintExportStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "export foo as bar;\n",
			expected: "export foo as bar;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-messy",
			code:     "  export   foo   as   bar  ;  \n",
			expected: "export foo as bar;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-wrap-long-name",
			code:     "export a_very_long_expression_name_that_exceeds_forty_characters as bar;\n",
			expected: "export a_very_long_expression_name_that_exceeds_forty_characters\n    as bar;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-wrap-long-alias",
			code:     "export foo as a_very_long_expression_name_that_exceeds_forty_characters;\n",
			expected: "export foo\n    as a_very_long_expression_name_that_exceeds_forty_characters;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-wrap-long-both",
			code:     "export a_very_long_expression_name_that_exceeds_forty_characters as another_very_long_expression_name_that_exceeds_forty_characters;\n",
			expected: "export a_very_long_expression_name_that_exceeds_forty_characters\n    as another_very_long_expression_name_that_exceeds_forty_characters;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExportStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-comment-tail",
			code:     "export foo as bar; // comment\n",
			expected: "export foo as bar; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-inline",
			code:     "export /* foo */ foo as bar;\n",
			expected: "export /* foo */ foo as bar;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-scattered-basic",
			code:     "export // c1\nfoo // c2\nas // c3\nbar; // c4\n",
			expected: "export // c1\nfoo // c2\n    as // c3\n    bar; // c4\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-scattered-inline",
			code:     "export /* c1 */ foo /* c2 */ as /* c3 */ bar /* c4 */;\n",
			expected: "export /* c1 */ foo /* c2 */\n    as /* c3 */ bar /* c4 */;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
