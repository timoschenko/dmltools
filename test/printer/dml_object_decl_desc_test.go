// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintObjectDecl_Desc(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-basic-desc-block",
			code:     "bank my_bank `Desc` {}\n",
			expected: "bank my_bank `Desc` {\n}\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-desc-semi",
			code:     "bank my_bank `Desc`;\n",
			expected: "bank my_bank `Desc`;\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-with-array-list",
			code:     "bank my_bank [ i < 10 ] `Desc` {}\n",
			expected: "bank my_bank[i < 10] `Desc` {\n}\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-register-size-offset",
			code:     "register my_reg size 4 @ 0x100 `Desc` {}\n",
			expected: "register my_reg size 4 @ 0x100 `Desc` {\n}\n",
			lineCap:  60,
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-field-bitrange",
			code:     "field my_field @ [ 10 ] `Desc` {}\n",
			expected: "field my_field @ [10] `Desc` {\n}\n",
			lineCap:  60,
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-with-template",
			code:     "bank my_bank is (T) `Desc` {}\n",
			expected: "bank my_bank is (T) `Desc` {\n}\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-composed-string",
			code:     "bank my_bank `Desc1` + `Desc2` {}\n",
			expected: "bank my_bank `Desc1` + `Desc2` {\n}\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 40)
}

func TestPrintObjectDecl_Desc_Comments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-comment-before-desc",
			code:     "bank my_bank /* c1 */ `Desc` {}\n",
			expected: "bank my_bank /* c1 */ `Desc` {\n}\n",
			lineCap:  60,
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-after-desc",
			code:     "bank my_bank `Desc` /* c2 */ {}\n",
			expected: "bank my_bank `Desc` /* c2 */ {\n}\n",
			lineCap:  60,
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-inside-composed",
			code:     "bank my_bank `Desc1` + /* c3 */ `Desc2` {}\n",
			expected: "bank my_bank `Desc1` + /* c3 */ `Desc2` {\n}\n",
			lineCap:  60,
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 40)
}
