// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestCCodeBlock(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-header-single-line",
			code:     "header %{ int x = 0; %}\n",
			expected: "header %{\n int x = 0;\n%}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-header-multi-line",
			code:     "header %{\n    int x = 0;   \n    int y = 1;\n%}\n",
			expected: "header %{\n    int x = 0;\n    int y = 1;\n%}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-header2-multi-line",
			code:     "_header %{\n    int x = 0;\n%}\n",
			expected: "_header %{\n    int x = 0;\n%}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-footer-multi-line",
			code:     "footer %{\n    int x = 0;\n%}\n",
			expected: "footer %{\n    int x = 0;\n%}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-strip-trailing-whitespace",
			code:     "header %{\n    int x = 0;   \n    int y = 1;\t\n%}\n",
			expected: "header %{\n    int x = 0;\n    int y = 1;\n%}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comments-around-block",
			code:     "// start\nheader %{\n    int x = 0;\n%}\n// end\n",
			expected: "// start\nheader %{\n    int x = 0;\n%}\n// end\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
