// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintExpressionEach(t *testing.T) {
	t.Parallel()

	const lineCap = 80

	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "param TESTCASE =   each   i    in   (   expr   );\n",
			expected: "param TESTCASE = each i in (expr);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-messy-whitespace",
			code:     "param TESTCASE = each   i  \n in   (   expr   )  ;\n",
			expected: "param TESTCASE = each i in (expr);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comments",
			code:     "param TESTCASE = each    /*c1*/    i   /*c2*/   in   /*c3*/   (    /*c4*/   expr   /*c5*/   )  ;\n",
			expected: "param TESTCASE = each /*c1*/ i /*c2*/ in /*c3*/ (/*c4*/ expr /*c5*/);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-register-ident",
			code:     "param TESTCASE =   each   register   in  ( expr ) ;\n",
			expected: "param TESTCASE = each register in (expr);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-nested",
			code:     "param TESTCASE =   each   i   in   (  each  j   in   (  expr  )  );\n",
			expected: "param TESTCASE = each i in (each j in (expr));\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-chain",
			code:     "param TESTCASE = each   i in   ( obj  .member->   list);\n",
			expected: "param TESTCASE = each i in (obj.member->list);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionEach_Wrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 30

	testCases := []printerTestCase{
		{
			name:     "01-wrap-list",
			code:     "param TESTCASE = each i in (long_list_expression);\n",
			expected: "param TESTCASE = each i in (\n    long_list_expression\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"param TESTCASE = each long_ident_name in (\n    expr\n);\n",
			),
			name:     "02-wrap-each",
			code:     "param TESTCASE = each long_ident_name in (expr);\n",
			expected: "param TESTCASE = (\n    each long_ident_name in (expr)\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"param T = each long_ident_name in (\n    long_list_expression\n);\n",
			),
			name:     "03-wrap-list-and-each",
			code:     "param T = each long_ident_name in (long_list_expression);\n",
			expected: "param T = (\n    each long_ident_name in (\n        long_list_expression\n    )\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"param T = each long_ident_name in (\n    long_list_expression\n);\n",
			),
			name:     "04-both-wrap",
			code:     "param T = each long_ident_name in (long_list_expression);\n",
			expected: "param T = (\n    each long_ident_name in (\n        long_list_expression\n    )\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-wrap-chain",
			code:     "param T = each i in (obj.long_access_chain);\n",
			expected: "param T = each i in (\n    obj.long_access_chain\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-wrap-long-chain",
			code:     "param T = each i in (obj.long_member_access_chain);\n",
			expected: "param T = each i in (\n    obj\n    .long_member_access_chain\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-wrap-inner-list", // optional/future
			code:     "param T = each i in (a[long_index_here]);\n",
			expected: "param T = each i in (\n    a[long_index_here]\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-wrap-each-simple",
			code:     "param T = each i in (\n    a_long_list_expression\n);\n",
			expected: "param T = each i in (\n    a_long_list_expression\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionEach_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 50

	testCases := []printerTestCase{
		{
			name:     "01-wrap-list-block-comments",
			code:     "param TESTCASE = each /*c1*/ i /*c2*/ in /*c3*/ ( /*c4*/ long_list_expression /*c5*/ );\n",
			expected: "param TESTCASE = each /*c1*/ i /*c2*/ in /*c3*/ (\n    /*c4*/ long_list_expression /*c5*/\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"param TESTCASE = each /*c1*/ long_ident /*c2*/ in /*c3*/ (\n    /*c4*/ e /*c5*/\n);\n",
			),
			name:     "02-wrap-each-block-comments",
			code:     "param TESTCASE = each /*c1*/ long_ident /*c2*/ in /*c3*/ ( /*c4*/ e /*c5*/ );\n",
			expected: "param TESTCASE = (\n    each /*c1*/ long_ident /*c2*/ in /*c3*/ (\n        /*c4*/ e\n    /*c5*/)\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"param T = each /*c1*/ long_ident_name /*c2*/ in /*c3*/ (\n    /*c4*/ long_list_expression /*c5*/\n);\n",
			),
			name:     "03-wrap-list-and-each-block-comments",
			code:     "param T = each /*c1*/ long_ident_name /*c2*/ in /*c3*/ ( /*c4*/ long_list_expression /*c5*/ );\n",
			expected: "param T = (\n    each /*c1*/ long_ident_name /*c2*/ in /*c3*/ (\n        /*c4*/ long_list_expression\n    /*c5*/)\n);\n",
			pc:       func(f func()) { f() },
		},
		// { // TODO: BUG: prints malformed code
		// 	name:     "04-wrap-list-line-comments",
		// 	code:     "param TESTCASE = each i in ( //c1\n long_list_expression //c2\n );\n",
		// 	expected: "param TESTCASE = each i in ( //c1\n    long_list_expression //c2\n);\n",
		// 	masked:   MaskUncategorizedFormatting("param TESTCASE = each i in (\n    //c1\n    long_list_expression //c2\n);\n"),
		// 	pc:       func(f func()) { f() },
		// },
		{
			masked: MaskSyntheticParens(
				"param TESTCASE = each //c1\n    long_ident_name //c2\n    in //c3\n    (\n        expr\n    );\n",
			),
			name:     "05-wrap-each-line-comments",
			code:     "param TESTCASE = each //c1\n long_ident_name //c2\n in //c3\n (expr);\n",
			expected: "param TESTCASE = (\n    each //c1\n    long_ident_name //c2\n    in //c3\n    (expr)\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-wrap-inner-list-comments",
			code:     "param T = each i in (a[ /*c1*/ long_index_here /*c2*/ ]);\n",
			expected: "param T = each i in (\n    a[/*c1*/ long_index_here /*c2*/]\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
