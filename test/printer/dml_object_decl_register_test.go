// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintObjectDecl_Register(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-basic-semi",
			code:     "register     my_reg    ;\n",
			expected: "register my_reg;\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-block",
			code:     "register    my_reg    {   }\n",
			expected: "register my_reg {\n}\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-with-array-list",
			code:     "register    my_reg     [   i   <   10   ];\n",
			expected: "register my_reg[i < 10];\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-with-size",
			code:     "register    my_reg   size   4  ;\n",
			expected: "register my_reg size 4;\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-with-offset",
			code:     "register    my_reg    @    0x100   ;\n",
			expected: "register my_reg @ 0x100;\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-with-size-and-offset",
			code:     "register    my_reg    size   4    @    0x100   ;\n",
			expected: "register my_reg size 4 @ 0x100;\n",
			lineCap:  40,
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-full-signature",
			code:     "register my_reg    [    i   < 10    ]   size   4   @   (  0x100   +   i   *   4)   {   field   f   @  [ 0 ] ;  }\n",
			expected: "register my_reg[i < 10] size 4 @ (0x100 + i * 4) {\n    field f @ [0];\n}\n",
			lineCap:  80,
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 40)
}

func TestPrintObjectDecl_Register_Comments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-comment-after-size",
			code:     "register   my_reg   size   /* c */   4;\n",
			expected: "register my_reg size /* c */ 4;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-after-offset",
			code:     "register   my_reg   @  /* c */    0x100 ;\n",
			expected: "register my_reg @ /* c */ 0x100;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 40)
}
