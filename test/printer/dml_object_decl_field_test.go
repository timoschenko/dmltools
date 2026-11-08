// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintObjectDecl_Field(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-basic-semi",
			code:     "   field my_field  ;\n",
			expected: "field my_field;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-block",
			code:     "field   my_field   {\n}\n",
			expected: "field my_field {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-with-array-list",
			code:     "field   my_field [ i < 10 ]  ;\n",
			expected: "field my_field[i < 10];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-bitrange-single",
			code:     "field   my_field  @ [ \n 10 ];\n",
			expected: "field my_field @ [10];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-bitrange-range",
			code:     "field my_field @ [ 31    :  16   ];\n",
			expected: "field my_field @ [31:16];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-array-and-bitrange",
			code:     "field my_field [ i <   10 ]@[   31:   16 ]   {param p; }\n",
			expected: "field my_field[i < 10] @ [31:16] {\n    param p;\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintObjectDecl_Field_Nested(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-nested-field",
			code:     "register my_reg {\nfield in_field   @   [   10 ];\n}\n",
			expected: "register my_reg {\n    field in_field @ [10];\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-nested-field-array",
			code:     "register my_reg {\nfield    in_field [ i <   10   ] @ [   31 :   16 ];\n}\n",
			expected: "register my_reg {\n    field in_field[i < 10] @ [31:16];\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintObjectDecl_Field_Comments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-comment-after-semi",
			code:     "field   my_field;     // trailing comment\n",
			expected: "field my_field; // trailing comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-inline-bitrange",
			code:     "field my_field   @ [   /* c1 */    31    /* c2 */     :  /* c3 */    16   /* c4 */   ];\n",
			expected: "field my_field @ [ /* c1 */ 31 /* c2 */:/* c3 */ 16 /* c4 */ ];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-trailing-bitrange",
			code:     "field my_field   @   [   10   ];    // comment\n",
			expected: "field my_field @ [10]; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-before-block",
			code:     "field my_field @   [   10 ]    /* c */ {\n}\n",
			expected: "field my_field @ [10] /* c */ {\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
