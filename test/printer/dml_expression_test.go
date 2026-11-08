// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestPrintExpressionLiterals(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-num-int",
			code:     "constant   num_int   =   123;\n",
			expected: "constant num_int = 123;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-num-hex",
			code:     "constant   num_hex   =   0xFF  ;\n",
			expected: "constant num_hex = 0xFF;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-num-bin",
			code:     "constant   TESTCASE   =   0b1010  ;\n",
			expected: "constant TESTCASE = 0b1010;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-num-float",
			code:     "constant   TESTCASE   =   3.14  ;\n",
			expected: "constant TESTCASE = 3.14;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-char",
			code:     "constant   TESTCASE   =   'c'  ;\n",
			expected: "constant TESTCASE = 'c';\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-string",
			code:     "constant   TESTCASE   =   \"string\"   ;\n",
			expected: "constant TESTCASE = \"string\";\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-undefined",
			code:     "constant    TESTCASE    =   undefined   ;\n",
			expected: "constant TESTCASE = undefined;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-default",
			code:     "constant    TESTCASE   =    default   ;\n",
			expected: "constant TESTCASE = default;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-this",
			code:     "constant    TESTCASE    =    this   ;  \n",
			expected: "constant TESTCASE = this;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-ident",
			code:     "constant   TESTCASE    =    ident   ;\n",
			expected: "constant TESTCASE = ident;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-messy-whitespace-num",
			code:     "   constant   TESTCASE   =   123  ;\n",
			expected: "constant TESTCASE = 123;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-messy-newlines",
			code:     "\nconstant\nTESTCASE\n=\n123\n;\n",
			expected: "constant TESTCASE = 123;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-num-int-comment-before",
			code:     "constant TESTCASE = // comment\n123;\n",
			expected: "constant TESTCASE = // comment\n    123;\n", // TODO: Use synthetic paren injection
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant T = 123 // comment\n    ;\n",
			),
			name:     "14-num-int-comment-after",
			code:     "constant T = 123 // comment\n;\n",
			expected: "constant T = (\n    123 // comment\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-num-int-inline-comment-before",
			code:     "constant   TESTCASE   =   /* comment */    123;\n",
			expected: "constant TESTCASE = /* comment */ 123;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-num-int-inline-comment-after",
			code:     "constant TESTCASE = 123    /* comment */   ;\n",
			expected: "constant TESTCASE = 123 /* comment */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-num-int-multiple-comments",
			code:     "constant TESTCASE = // one\n// two\n123;\n",
			expected: "constant TESTCASE = // one\n    // two\n    123;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-num-int-messy-with-comments",
			code:     "   constant   /* c1 */   TESTCASE   /* c2 */   =   /* c3 */   123   /* c4 */   ;   // c5\n",
			expected: "constant /* c1 */ TESTCASE /* c2 */ = /* c3 */ 123 /* c4 */; // c5\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskTask006(
				"constant TESTCASE = \"string\" // comment\n    ;\n",
			),
			name:     "19-string-with-comment",
			code:     "constant TESTCASE = \"string\" // comment\n;\n",
			expected: "constant TESTCASE = \"string\" // comment\n;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "20-ident-with-comment",
			code:     "constant TESTCASE = ident /* comment */;\n",
			expected: "constant TESTCASE = ident /* comment */;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionArithmetic(t *testing.T) {
	t.Parallel()

	const lineCap = 50

	testCases := []printerTestCase{
		// Standard Unary Operators
		{
			pc:       func(f func()) { f() },
			name:     "01-unary-minus",
			code:     "constant TESTCASE = - a;\n",
			expected: "constant TESTCASE = -a;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "02-unary-plus",
			code:     "constant TESTCASE = + b;\n",
			expected: "constant TESTCASE = +b;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "03-pre-inc",
			code:     "constant TESTCASE = ++ a;\n",
			expected: "constant TESTCASE = ++a;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "04-pre-dec",
			code:     "constant TESTCASE = -- a;\n",
			expected: "constant TESTCASE = --a;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "05-post-inc",
			code:     "constant TESTCASE = a ++;\n",
			expected: "constant TESTCASE = a++;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "06-post-dec",
			code:     "constant TESTCASE = a --;\n",
			expected: "constant TESTCASE = a--;\n",
		},

		// Standard Binary Operators
		{
			pc:       func(f func()) { f() },
			name:     "07-add",
			code:     "constant TESTCASE = a+b;\n",
			expected: "constant TESTCASE = a + b;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "08-sub",
			code:     "constant TESTCASE = a-b;\n",
			expected: "constant TESTCASE = a - b;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "09-mul",
			code:     "constant TESTCASE = a*b;\n",
			expected: "constant TESTCASE = a * b;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "10-div",
			code:     "constant TESTCASE = a/b;\n",
			expected: "constant TESTCASE = a / b;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "11-mod",
			code:     "constant TESTCASE = a%b;\n",
			expected: "constant TESTCASE = a % b;\n",
		},

		// Messy Whitespace & Newlines
		{
			pc:       func(f func()) { f() },
			name:     "12-messy-binary",
			code:     "constant TESTCASE =   a   \n  +  \n  b  ;\n",
			expected: "constant TESTCASE = a + b;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "13-messy-unary",
			code:     "constant TESTCASE =  \n  - \n   a \n ;\n",
			expected: "constant TESTCASE = -a;\n",
		},

		// Expressions with Comments
		{
			pc:       func(f func()) { f() },
			name:     "14-binary-comment-inline",
			code:     "constant TC =/* c0 */   a\n     /* c1 */   %/* c2 */b   ;\n",
			expected: "constant TC = /* c0 */ a /* c1 */ % /* c2 */ b;\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "16-unary-comment",
			code:     "constant T =/*a*/++   /*b*/ \n a;\n",
			expected: "constant T = /*a*/ ++/*b*/ a;\n",
			masked: MaskBlockCommentWrapping(
				"constant T = /*a*/ ++/*b*/\n    a;\n",
			),
		},
		{
			pc:       func(f func()) { f() },
			name:     "17-post-unary-comment",
			code:     "constant T =a/*b*/ ++ /*c*/;\n",
			expected: "constant T = a /*b*/++ /*c*/;\n",
		},
		{
			pc: func(f func()) { f() },
			masked: MaskSyntheticParens(
				"constant TESTCASE = a + // c1\n    b;\n",
			),
			name:     "18-binary-comment-line",
			code:     "constant TESTCASE = a + // c1\nb;\n",
			expected: "constant TESTCASE = (\n    a + // c1\n    b\n);\n",
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionBitwiseAndLogical(t *testing.T) {
	t.Parallel()

	const lineCap = 50

	testCases := []printerTestCase{
		// Basic unary operators
		{
			name:     "01-bitwise-not",
			code:     "constant TESTCASE = ~ a;\n",
			expected: "constant TESTCASE = ~a;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-logical-not",
			code:     "constant TESTCASE = ! a;\n",
			expected: "constant TESTCASE = !a;\n",
			pc:       func(f func()) { f() },
		},

		// Basic binary operators
		{
			name:     "03-bitwise-and",
			code:     "constant TESTCASE = a&b;\n",
			expected: "constant TESTCASE = a & b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-bitwise-or",
			code:     "constant TESTCASE = a|b;\n",
			expected: "constant TESTCASE = a | b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-bitwise-xor",
			code:     "constant TESTCASE = a^b;\n",
			expected: "constant TESTCASE = a ^ b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-left-shift",
			code:     "constant TESTCASE = a<<2;\n",
			expected: "constant TESTCASE = a << 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-right-shift",
			code:     "constant TESTCASE = a>>2;\n",
			expected: "constant TESTCASE = a >> 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-logical-and",
			code:     "constant TESTCASE = a&&b;\n",
			expected: "constant TESTCASE = a && b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-logical-or",
			code:     "constant TESTCASE = a||b;\n",
			expected: "constant TESTCASE = a || b;\n",
			pc:       func(f func()) { f() },
		},

		// Messy whitespace and comments
		{
			name:     "10-messy-binary",
			code:     "constant TESTCASE =   a   \n  &  \n  b  ;\n",
			expected: "constant TESTCASE = a & b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-messy-unary",
			code:     "constant TESTCASE =  \n  ~ \n   a \n ;\n",
			expected: "constant TESTCASE = ~a;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-binary-comment-inline",
			code:     "constant TC =/* c0 */  a   \n    /* c1 */   ||/* c2 */b   ;\n",
			expected: "constant TC = /* c0 */ a /* c1 */ || /* c2 */ b;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionRelationalAndEquality(t *testing.T) {
	t.Parallel()

	const lineCap = 50

	testCases := []printerTestCase{
		{
			name:     "01-eq",
			code:     "constant TESTCASE = a==b;\n",
			expected: "constant TESTCASE = a == b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-neq",
			code:     "constant TESTCASE = a!=b;\n",
			expected: "constant TESTCASE = a != b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-lt",
			code:     "constant TESTCASE = a<b;\n",
			expected: "constant TESTCASE = a < b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-gt",
			code:     "constant TESTCASE = a>b;\n",
			expected: "constant TESTCASE = a > b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-lte",
			code:     "constant TESTCASE = a<=b;\n",
			expected: "constant TESTCASE = a <= b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-gte",
			code:     "constant TESTCASE = a>=b;\n",
			expected: "constant TESTCASE = a >= b;\n",
			pc:       func(f func()) { f() },
		},

		// Messy whitespace and comments
		{
			name:     "07-messy-relational",
			code:     "constant TESTCASE =   a   \n  >=  \n  b  ;\n",
			expected: "constant TESTCASE = a >= b;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-comment-inline",
			code:     "constant TC =/* c0 */  a \n    /* c1 */   ==/* c2 */b   ;\n",
			expected: "constant TC = /* c0 */ a /* c1 */ == /* c2 */ b;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionMemoryAndTypes(t *testing.T) {
	t.Parallel()

	const lineCap = 50

	testCases := []printerTestCase{
		// Memory Operations (Unary)
		{
			name:     "01-addressof",
			code:     "constant TESTCASE =   & a;\n",
			expected: "constant TESTCASE = &a;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-deref",
			code:     "constant TESTCASE =   *   a;\n",
			expected: "constant TESTCASE = *a;\n",
			pc:       func(f func()) { f() },
		},

		// Builtin unary-like operators
		{
			masked: MaskSyntheticSizeofParens(
				"constant TESTCASE = sizeof a;\n",
			),
			name:     "03-sizeof-no-paren",
			code:     "constant TESTCASE = sizeof a;\n",
			expected: "constant TESTCASE = sizeof(a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-sizeof-paren",
			code:     "constant TESTCASE =   sizeof( a );\n",
			expected: "constant TESTCASE = sizeof(a);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticSizeofParens(
				"constant TESTCASE = sizeoftype int;\n",
			),
			name:     "05-sizeoftype",
			code:     "constant TESTCASE =   sizeoftype  int  ;\n",
			expected: "constant TESTCASE = sizeoftype(int);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskTask006(
				"constant T = /*0*/ sizeoftype(\n    /*1*/ int /*2*/\n) /*3*/; // 4\n",
			),
			name:     "06-sizeoftype-comments",
			code:     "constant T =  /*0*/ sizeoftype  (/*1*/ int  /*2*/ ) /*3*/ ; // 4\n",
			expected: "constant T = /*0*/ sizeoftype(/*1*/ int /*2*/) /*3*/; // 4\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-defined",
			code:     "constant TESTCASE =   defined  a  ;\n",
			expected: "constant TESTCASE = defined a;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-defined-paren",
			code:     "constant TESTCASE =   defined ( a ) ;\n",
			expected: "constant TESTCASE = defined(a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-stringify",
			code:     "constant TESTCASE =   stringify(  a  );\n",
			expected: "constant TESTCASE = stringify(a);\n",
			pc:       func(f func()) { f() },
		},

		// New Operator
		{
			name:     "10-new-simple",
			code:     "constant TESTCASE =   new   int   ;\n",
			expected: "constant TESTCASE = new int;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-new-array",
			code:     "constant TESTCASE =   new   int [ 10 ] ;\n",
			expected: "constant TESTCASE = new int[10];\n",
			pc:       func(f func()) { f() },
		},

		// Messy testing
		{
			masked: MaskSyntheticSizeofParens(
				"constant TESTCASE = sizeof a;\n",
			),
			name:     "12-messy-sizeof",
			code:     "constant TESTCASE =   sizeof  \n  a  ;\n",
			expected: "constant TESTCASE = sizeof(a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-messy-new",
			code:     "constant TESTCASE =   new \n  int  [ \n  10 \n  ]  ;\n",
			expected: "constant TESTCASE = new int[10];\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticSizeofParens(
				"constant TESTCASE = /* c0 */ sizeof /* c1 */ a;\n",
			),
			name:     "14-comment-inline",
			code:     "constant TESTCASE =/* c0 */     sizeof   /* c1 */   a   ;\n",
			expected: "constant TESTCASE = /* c0 */ sizeof(/* c1 */ a);\n", // TODO: Synthetic token support
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-comment-inline",
			code:     "constant TESTCASE =/* c0 */     sizeof (  /* c1 */   a )  ;\n",
			expected: "constant TESTCASE = /* c0 */ sizeof(/* c1 */ a);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionBinaryOp_Newlines(t *testing.T) {
	t.Parallel()

	const lineCap = 30

	testCases := []printerTestCase{
		{
			name:     "001-clean",
			code:     "constant TE = lv %% mv %% rv;\n",
			expected: "constant TE = lv %% mv %% rv;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-clean-paren",
			code:     "constant T = (lv %% mv %% rv);\n",
			expected: "constant T = (lv %% mv %% rv);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "101-paren-l",
			code:     "constant TESTCASE = (l_name_that_doesnt_fit %% mval %% rval);\n",
			expected: "constant TESTCASE = (\n    l_name_that_doesnt_fit\n    %% mval\n    %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "102-paren-m",
			code:     "constant TESTCASE = (lval %% m_name_that_doesnt_fit %% rval);\n",
			expected: "constant TESTCASE = (\n    lval\n    %% m_name_that_doesnt_fit\n    %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "103-paren-r",
			code:     "constant TESTCASE = (lval %% rval %% r_name_that_doesnt_fit);\n",
			expected: "constant TESTCASE = (\n    lval\n    %% rval\n    %% r_name_that_doesnt_fit\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "201-paren-l-super",
			code:     "constant TESTCASE = (l_super_long_name_that_doesnt_fit %% mval %% rval);\n",
			expected: "constant TESTCASE = (\n    l_super_long_name_that_doesnt_fit\n    %% mval\n    %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "202-paren-m-super",
			code:     "constant TESTCASE = (lval %% m_super_long_name_that_doesnt_fit %% rval);\n",
			expected: "constant TESTCASE = (\n    lval\n    %% m_super_long_name_that_doesnt_fit\n    %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "203-paren-r-super",
			code:     "constant TESTCASE = (lval %% rval %% r_super_long_name_that_doesnt_fit);\n",
			expected: "constant TESTCASE = (\n    lval\n    %% rval\n    %% r_super_long_name_that_doesnt_fit\n);\n",
			pc:       func(f func()) { f() },
		},

		// If expression is too long and require newlines then add parens
		// to visially highlight that this is one expression
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = l_name_that_doesnt_fit\n    %% mval\n    %% rval;\n",
			),
			name:     "301-clean-l",
			code:     "constant TESTCASE = l_name_that_doesnt_fit %% mval %% rval;\n",
			expected: "constant TESTCASE = (\n    l_name_that_doesnt_fit\n    %% mval\n    %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% m_name_that_doesnt_fit\n    %% rval;\n",
			),
			name:     "302-clean-m",
			code:     "constant TESTCASE = lval %% m_name_that_doesnt_fit %% rval;\n",
			expected: "constant TESTCASE = (\n    lval\n    %% m_name_that_doesnt_fit\n    %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% rval\n    %% r_name_that_doesnt_fit;\n",
			),
			name:     "303-clean-r",
			code:     "constant TESTCASE = lval %% rval %% r_name_that_doesnt_fit;\n",
			expected: "constant TESTCASE = (\n    lval\n    %% rval\n    %% r_name_that_doesnt_fit\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = l_super_long_name_that_doesnt_fit\n    %% mval\n    %% rval;\n",
			),
			name:     "401-clean-l-super",
			code:     "constant TESTCASE = l_super_long_name_that_doesnt_fit %% mval %% rval;\n",
			expected: "constant TESTCASE = (\n    l_super_long_name_that_doesnt_fit\n    %% mval %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% m_super_long_name_that_doesnt_fit\n    %% rval;\n",
			),
			name:     "402-clean-m-super",
			code:     "constant TESTCASE = lval %% m_super_long_name_that_doesnt_fit %% rval;\n",
			expected: "constant TESTCASE = (\n    lval\n    %% m_super_long_name_that_doesnt_fit\n    %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% rval\n    %% r_super_long_name_that_doesnt_fit;\n",
			),
			name:     "403-clean-r-super",
			code:     "constant TESTCASE = lval %% rval %% r_super_long_name_that_doesnt_fit;\n",
			expected: "constant TESTCASE = (\n    lval\n    %% rval\n    %% r_super_long_name_that_doesnt_fit\n);\n",
			pc:       func(f func()) { f() },
		},

		// Check that if line long to require a new line but not long enouht to be multiline
		{
			name:     "501-paren-l-oneline",
			code:     "constant TESTCASE = (l_lg_name %% mval %% rval);\n",
			expected: "constant TESTCASE = (\n    l_lg_name %% mval %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "502-paren-m-oneline",
			code:     "constant TESTCASE = (lval %% m_lg_name %% rval);\n",
			expected: "constant TESTCASE = (\n    lval %% m_lg_name %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "503-paren-r-oneline",
			code:     "constant TESTCASE = (lval %% rval %% r_lg_name);\n",
			expected: "constant TESTCASE = (\n    lval %% rval %% r_lg_name\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = l_long_name\n    %% mval\n    %% rval;\n",
			),
			name:     "601-clean-l-oneline",
			code:     "constant TESTCASE = l_long_name %% mval %% rval;\n",
			expected: "constant TESTCASE = (\n    l_long_name %% mval %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% m_long_name\n    %% rval;\n",
			),
			name:     "602-clean-m-oneline",
			code:     "constant TESTCASE = lval %% m_long_name %% rval;\n",
			expected: "constant TESTCASE = (\n    lval %% m_long_name %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% rval\n    %% r_long_name;\n",
			),
			name:     "603-clean-r-oneline",
			code:     "constant TESTCASE = lval %% rval %% r_long_name;\n",
			expected: "constant TESTCASE = (\n    lval %% rval %% r_long_name\n);\n",
			pc:       func(f func()) { f() },
		},

		// Comments test cases
		{
			name:     "701-clean-l-comment-line",
			code:     "constant TESTCASE = // coment\nlval %% mval %% rval;\n",
			expected: "constant TESTCASE = // coment\n    lval %% mval %% rval;\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% // comment\n    mval\n    %% rval;\n",
			),
			name:     "702-clean-m-comment-line",
			code:     "constant TESTCASE = lval %% // comment\n mval %% rval;\n",
			expected: "constant TESTCASE = (\n    lval %% // comment\n    mval %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% mval\n    %% // comment\n    rval;\n",
			),
			name:     "703-clean-r-comment-line",
			code:     "constant TESTCASE = lval %% mval %% // comment\n rval;\n",
			expected: "constant TESTCASE = (\n    lval %% mval %% // comment\n    rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval // comment\n    %% mval %% rval;\n",
			),
			name:     "801-clean-l-comment-line",
			code:     "constant TESTCASE = lval // comment\n %% mval %% rval;\n",
			expected: "constant TESTCASE = (\n    lval // comment\n    %% mval %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% mval // comment\n    %% rval;\n",
			),
			name:     "802-clean-m-comment-line",
			code:     "constant TESTCASE = lval %% mval  // comment\n %% rval;\n",
			expected: "constant TESTCASE = (\n    lval %% mval // comment\n    %% rval\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticParens(
				"constant TESTCASE = lval\n    %% mval\n    %% rval // comment\n    ;\n",
			),
			name:     "803-clean-r-comment-line",
			code:     "constant TESTCASE = lval %% mval %% rval // comment\n;\n",
			expected: "constant TESTCASE = (\n    lval %% mval %% rval // comment\n);\n",
			pc:       func(f func()) { f() },
		},
	}
	opcodes := []struct {
		name, op string
	}{
		{"bitwise-and", "&"},
		{"bitwise-or", "|"},
		{"bitwise-xor", "^"},
		{"left-shift", "<<"},
		{"right-shift", ">>"},
		{"logical-and", "&&"},
		{"logical-or", "||"},
		{"eq", "=="},
		{"neq", "!="},
		{"lt", "<"},
		{"gt", ">"},
		{"lte", "<="},
		{"gte", ">="},
		{"add", "+"},
		{"sub", "-"},
		{"mul", "*"},
		{"div", "/"},
		{"modulo", "%"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			for _, op := range opcodes {
				runPrinterTestCase(
					tt,
					&printerTestCase{
						name:     op.name,
						code:     strings.ReplaceAll(tc.code, "%%", op.op),
						expected: strings.ReplaceAll(tc.expected, "%%", op.op),
						masked:   MaskReplaceAll(tc.masked, "%%", op.op),
						lineCap:  lineCap,
					},
				)
			}
		})
	}
}

func resolveUnaryMask(m Masked, op string, isPost bool) Masked {
	if unaryMask, ok := m.(MaskUnaryIncorrectFormatting); ok {
		return unaryMask.Resolve(op, isPost)
	}
	if isPost {
		return MaskReplaceAll(MaskReplaceAll(m, "%b", ""), "%p", op)
	}
	return MaskReplaceAll(MaskReplaceAll(m, "%p", ""), "%b", op)
}

func TestPrintExpressionUnaryOp_Newlines(t *testing.T) {
	t.Parallel()

	const lineCap = 30

	testCases := []struct {
		name, code, expected    string
		masked                  Masked
		prefixOnly, postfixOnly bool
	}{
		{
			name:     "001-clean",
			code:     "constant TES    =    %b   expr  %p  ;\n",
			expected: "constant TES = %bexpr%p;\n",
		},
		{
			name:     "002-clean-paren",
			code:     "constant TESTCASE    =   (   %b   expr   %p  )   ;\n",
			expected: "constant TESTCASE = (%bexpr%p);\n",
		},
		{
			name:     "101-paren-long",
			code:     "constant TESTCASE = (%bname_that_doesnt_fit%p);\n",
			expected: "constant TESTCASE = (\n    %bname_that_doesnt_fit%p\n);\n",
		},
		{
			name:     "201-paren-super",
			code:     "constant TESTCASE = (%bsuper_long_name_that_doesnt_fit%p);\n",
			expected: "constant TESTCASE = (\n    %bsuper_long_name_that_doesnt_fit%p\n);\n",
		},
		{
			masked: MaskSyntheticParensUnary(
				"constant TESTCASE =\n    %bname_that_doesnt_fit;\n",
				"constant TESTCASE = %bname_that_doesnt_fit;\n",
				"constant TESTCASE = name_that_doesnt_fit%p;\n",
			),
			name:     "301-clean-long",
			code:     "constant TESTCASE = %bname_that_doesnt_fit%p;\n",
			expected: "constant TESTCASE = (\n    %bname_that_doesnt_fit%p\n);\n",
		},
		{
			masked: MaskSyntheticParensUnary(
				"constant TESTCASE =\n    %bsuper_long_name_that_doesnt_fit;\n",
				"constant TESTCASE = %bsuper_long_name_that_doesnt_fit;\n",
				"constant TESTCASE = super_long_name_that_doesnt_fit%p;\n",
			),
			name:     "401-clean-super",
			code:     "constant TESTCASE = %bsuper_long_name_that_doesnt_fit%p;\n",
			expected: "constant TESTCASE = (\n    %bsuper_long_name_that_doesnt_fit%p\n);\n",
		},
		{
			name:     "501-paren-lpartlong",
			code:     "constant TESTCASE = (%blong_name%p);\n",
			expected: "constant TESTCASE = (\n    %blong_name%p\n);\n",
		},
		// { // BUG: Commented out because most of the %b are pass except "&" and "!"
		// 	name:     "601-clean-lpartlong",
		// 	code:     "constant TESTCASE = %blong_name%p;\n",
		// 	expected: "constant TESTCASE = (\n    %blong_name%p\n);\n",
		// 	masked:   MaskUncategorizedFormatting("constant TESTCASE = %blong_name%p;\n"), // BUG(UFF-REGRESSION)(UFF-REGRESSION): Overflow not detected due to statement trailing punctuation
		// },

		// Comments test
		{
			name:     "701-comments-short",
			code:     "constant T =     /*a*/   %b    op   %p   /*b*/   ;\n",
			expected: "constant T = /*a*/ %bop%p /*b*/;\n",
		},
		{
			name:     "702-comments-short-in",
			code:     "constant T =    %b    /*a*/    op    /*b*/    %p   ;\n",
			expected: "constant T = %b/*a*/ op /*b*/%p;\n",
		},
		{
			name:     "801-comments-slash-l",
			code:     "constant T =     // comment\n   %b    op   %p    ;\n",
			expected: "constant T = // comment\n    %bop%p;\n", // TODO: Use synthetic paren inject
		},
		{
			masked: MaskSyntheticParensUnary(
				"constant T =\n    %bop // comment\n    ;\n",
				"constant T = %bop // comment\n    ;\n",
				"constant T = op // comment\n    %p;\n",
			),
			name:     "802-comments-slash-r",
			code:     "constant T =      %b    op  // comment\n   %p    ;\n",
			expected: "constant T = (\n    %bop%p // comment\n);\n",
		},
		{
			masked: MaskSyntheticParensUnary(
				"constant T =\n    %b // comment\n    op;\n",
				"constant T = %b // comment\n    op;\n",
				"",
			),
			name:       "803-comments-slash-l-op",
			code:       "constant T =      %b  // comment\n  op    %p    ;\n",
			expected:   "constant T = (\n    %b// comment\n    op%p\n);\n",
			prefixOnly: true,
		},
		{
			masked: MaskSyntheticParensUnary(
				"constant T =\n    %bop // comment\n    ;\n",
				"constant T = %bop // comment\n    ;\n",
				"constant T = op%p // comment\n    ;\n",
			),
			name:     "804-comments-slash-r-op",
			code:     "constant T =      %b   op    %p  // comment\n   ;\n",
			expected: "constant T = (\n    %bop%p // comment\n);\n",
		},
	}

	preOpcodes := []struct {
		name, op string
	}{
		{"unary-minus", "-"},
		{"unary-plus", "+"},
		{"pre-inc", "++"},
		{"pre-dec", "--"},
		{"bitwise-not", "~"},
		{"logical-not", "!"},
		{"addressof", "&"},
		{"deref", "*"},
	}

	for _, tc := range testCases {
		if tc.postfixOnly {
			continue
		}
		t.Run(tc.name, func(tt *testing.T) {
			for _, op := range preOpcodes {
				runPrinterTestCase(
					tt,
					&printerTestCase{
						name: op.name,
						code: strings.ReplaceAll(
							strings.ReplaceAll(tc.code, "%p", ""),
							"%b",
							op.op,
						),
						expected: strings.ReplaceAll(
							strings.ReplaceAll(tc.expected, "%p", ""),
							"%b",
							op.op,
						),
						masked:  resolveUnaryMask(tc.masked, op.op, false),
						lineCap: lineCap,
					},
				)
			}
		})
	}

	postOpcodes := []struct {
		name, op string
	}{
		{"post-inc", "++"},
		{"post-dec", "--"},
	}

	for _, tc := range testCases {
		if tc.prefixOnly {
			continue
		}
		t.Run(tc.name+"-post", func(tt *testing.T) {
			for _, op := range postOpcodes {
				runPrinterTestCase(
					tt,
					&printerTestCase{
						name: op.name,
						code: strings.ReplaceAll(
							strings.ReplaceAll(tc.code, "%p", op.op),
							"%b",
							"",
						),
						expected: strings.ReplaceAll(
							strings.ReplaceAll(tc.expected, "%p", op.op),
							"%b",
							"",
						),
						masked:  resolveUnaryMask(tc.masked, op.op, true),
						lineCap: lineCap,
					},
				)
			}
		})
	}
}

func TestPrintExpressionArray_Newlines(t *testing.T) {
	t.Parallel()

	const lineCap = 30

	testCases := []printerTestCase{
		{
			name:     "001-clean",
			code:     "constant TESTCASE =     arr   [  v  ]   ;\n",
			expected: "constant TESTCASE = arr[v];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-long-name",
			code:     "constant TESTCASE =     array_30    [  v  ]   ;\n",
			expected: "constant TESTCASE = array_30[\n    v\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "003-long-prefix",
			code:     "constant TESTCASE_1234 =     arr   [  v  ]   ;\n",
			expected: "constant TESTCASE_1234 = arr[\n    v\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "004-long-value",
			code:     "constant T =     arr   [  very_long_name_it_is  ]   ;\n",
			expected: "constant T = arr[\n    very_long_name_it_is\n];\n",
			pc:       func(f func()) { f() },
		},
		// Comments
		{
			name:     "101-comments-short",
			code:     "constant T =   /*c*/  arr    [   v   ]   ;\n",
			expected: "constant T = /*c*/ arr[v];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "102-comments-short",
			code:     "constant T =     arr  /*c*/  [   v   ]   ;\n",
			expected: "constant T = arr /*c*/ [v];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "103-comments-short",
			code:     "constant T =     arr   [  /*c*/  v   ]   ;\n",
			expected: "constant T = arr[/*c*/ v];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "104-comments-short",
			code:     "constant T =     arr   [    v   /*c*/   ]   ;\n",
			expected: "constant T = arr[v /*c*/];\n",
			pc:       func(f func()) { f() },
		},

		{
			name:     "105-comments-short",
			code:     "constant T =     arr   [    v     ]  /*c*/  ;\n",
			expected: "constant T = arr[v] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		// Comments - long name
		{
			name:     "111-comments-long-name",
			code:     "constant T =   /*c*/  arr_long_name    [   v   ]   ;\n",
			expected: "constant T = /*c*/ arr_long_name[\n    v\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "112-comments-long-name",
			code:     "constant T =     arr_long_name  /*c*/  [   v   ]   ;\n",
			expected: "constant T = arr_long_name /*c*/ [\n    v\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "113-comments-long-name",
			code:     "constant T =     arr_long_name   [  /*c*/  v   ]   ;\n",
			expected: "constant T = arr_long_name[\n    /*c*/ v\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "114-comments-long-name",
			code:     "constant T =     arr_long_name   [    v   /*c*/   ]   ;\n",
			expected: "constant T = arr_long_name[\n    v /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "115-comments-long-name",
			code:     "constant T =     arr_long_name   [    v     ]  /*c*/  ;\n",
			expected: "constant T = arr_long_name[\n    v\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		// Comments - long prefix
		{
			name:     "121-comments-long-prefix",
			code:     "constant TESTSUITE =   /*c*/  a    [   v   ]   ;\n",
			expected: "constant TESTSUITE = /*c*/ a[\n    v\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "122-comments-long-prefix",
			code:     "constant TESTSUITE =     a  /*c*/  [   v   ]   ;\n",
			expected: "constant TESTSUITE = a /*c*/ [\n    v\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "123-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [ /*c*/   v   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    /*c*/ v\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "124-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    v   /*c*/  ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    v /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "125-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    v     ]  /*c*/ ;\n",
			expected: "constant TESTSUITE123456 = a[\n    v\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		// Comments - long value
		{
			name:     "131-comments-long-value",
			code:     "constant T =   /*c*/  a    [   very_long_value_it_is   ]   ;\n",
			expected: "constant T = /*c*/ a[\n    very_long_value_it_is\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "132-comments-long-value",
			code:     "constant T =     a  /*c*/  [   very_long_value_it_is   ]   ;\n",
			expected: "constant T = a /*c*/ [\n    very_long_value_it_is\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "133-comments-long-value",
			code:     "constant T =     a   [  /*c*/  very_long_value_it_is   ]   ;\n",
			expected: "constant T = a[\n    /*c*/ very_long_value_it_is\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "134-comments-long-value",
			code:     "constant T =     a   [    very_long_value_it_is   /*c*/   ]   ;\n",
			expected: "constant T = a[\n    very_long_value_it_is /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "135-comments-long-value",
			code:     "constant T =     a   [    very_long_value_it_is     ]  /*c*/  ;\n",
			expected: "constant T = a[\n    very_long_value_it_is\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintBitSlice_Endiannes(t *testing.T) {
	t.Parallel()

	const lineCap = 31

	testCases := []printerTestCase{
		{
			name:     "001-clean",
			code:     "constant TESTCA =     arr   [  v , %e ]   ;\n",
			expected: "constant TESTCA = arr[v, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-long-name",
			code:     "constant TESTCASE =     array_30    [  v ,    %e ]   ;\n",
			expected: "constant TESTCASE = array_30[\n    v, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "003-long-prefix",
			code:     "constant TESTCASE_1234 =     arr   [  v , %e ]   ;\n",
			expected: "constant TESTCASE_1234 = arr[\n    v, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "004-long-value",
			code:     "constant T =     arr   [  very_long_name_it_is , %e ]   ;\n",
			expected: "constant T = arr[\n    very_long_name_it_is, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "005-long-supervalue",
			code:     "constant T =     arr   [  super_very_long_name_it_is , %e ]   ;\n",
			expected: "constant T = arr[\n    super_very_long_name_it_is, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		// Comments
		{
			name:     "101-comments-short",
			code:     "constant T =   /*c*/  arr    [   v , %e ]   ;\n",
			expected: "constant T = /*c*/ arr[v, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "102-comments-short",
			code:     "constant T =     arr  /*c*/  [   v , %e ]   ;\n",
			expected: "constant T = arr /*c*/[v, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "103-comments-short",
			code:     "constant T =     arr   [  /*c*/  v , %e ]   ;\n",
			expected: "constant T = arr[/*c*/ v, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "104-comments-short",
			code:     "constant T =     arr   [    v   /*c*/ , %e ]   ;\n",
			expected: "constant T = arr[v /*c*/, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "105-comments-short",
			code:     "constant T =     arr   [    v   ,   /*c*/   %e  ]   ;\n",
			expected: "constant T = arr[v, /*c*/ %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "106-comments-short",
			code:     "constant T =     arr   [    v   ,     %e   /*c*/  ]   ;\n",
			expected: "constant T = arr[v, %e /*c*/];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "107-comments-short",
			code:     "constant T =     arr   [    v   , %e  ]  /*c*/  ;\n",
			expected: "constant T = arr[v, %e] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		// Comments - long name
		{
			name:     "111-comments-long-name",
			code:     "constant T =   /*c*/  arr_long_name    [   v ,   %e  ]   ;\n",
			expected: "constant T = /*c*/ arr_long_name[\n    v, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "112-comments-long-name",
			code:     "constant T =     arr_long_name  /*c*/  [   v ,   %e  ]   ;\n",
			expected: "constant T = arr_long_name /*c*/[\n    v, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "113-comments-long-name",
			code:     "constant T =     arr_long_name   [  /*c*/  v ,   %e  ]   ;\n",
			expected: "constant T = arr_long_name[\n    /*c*/ v, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "114-comments-long-name",
			code:     "constant T =     arr_long_name   [    v   /*c*/   ,   %e  ]   ;\n",
			expected: "constant T = arr_long_name[\n    v /*c*/, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "115-comments-long-name",
			code:     "constant T =     arr_long_name   [    v   ,   /*c*/   %e   ]    ;\n",
			expected: "constant T = arr_long_name[\n    v, /*c*/ %e\n];\n",
			pc:       func(f func()) { f() },
		},

		{
			name:     "116-comments-long-name",
			code:     "constant T =     arr_long_name   [    v   ,     %e   /*c*/ ]    ;\n",
			expected: "constant T = arr_long_name[\n    v, %e /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "117-comments-long-name",
			code:     "constant T =     arr_long_name   [    v   ,   %e   ]  /*c*/  ;\n",
			expected: "constant T = arr_long_name[\n    v, %e\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		// Comments - long prefix
		{
			name:     "121-comments-long-prefix",
			code:     "constant TESTSUITE =   /*c*/  a    [   v   ,     %e   ]   ;\n",
			expected: "constant TESTSUITE = /*c*/ a[\n    v, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "122-comments-long-prefix",
			code:     "constant TESTSUI =     a  /*c*/  [   v   ,     %e   ]   ;\n",
			expected: "constant TESTSUI = a /*c*/[\n    v, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "123-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [ /*c*/   v   ,     %e   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    /*c*/ v, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "124-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    v   /*c*/   ,     %e   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    v /*c*/, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "125-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    v   ,   /*c*/   %e   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    v, /*c*/ %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "126-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    v   ,     %e   /*c*/ ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    v, %e /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "127-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    v   ,   %e   ]  /*c*/ ;\n",
			expected: "constant TESTSUITE123456 = a[\n    v, %e\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		// Comments - long value
		{
			name:     "131-comments-long-value",
			code:     "constant T =   /*c*/  a    [   very_long_value_it_is ,   %e  ]   ;\n",
			expected: "constant T = /*c*/ a[\n    very_long_value_it_is, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "132-comments-long-value",
			code:     "constant T =     a  /*c*/  [   very_long_value_it_is ,   %e  ]   ;\n",
			expected: "constant T = a /*c*/[\n    very_long_value_it_is, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "133-comments-long-value",
			code:     "constant T =     a   [  /*c*/  very_long_value_it_is ,   %e  ]   ;\n",
			expected: "constant T = a[\n    /*c*/ very_long_value_it_is, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "134-comments-long-value",
			code:     "constant T =     a   [    very_long_value_it_is   /*c*/   ,   %e  ]   ;\n",
			expected: "constant T = a[\n    very_long_value_it_is /*c*/, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "135-comments-long-value",
			code:     "constant T =     a   [    very_long_value_it_is   ,   /*c*/   %e   ]   ;\n",
			expected: "constant T = a[\n    very_long_value_it_is, /*c*/ %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "136-comments-long-value",
			code:     "constant T =     a   [    very_long_value_it_is   ,     %e   /*c*/ ]   ;\n",
			expected: "constant T = a[\n    very_long_value_it_is, %e /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "137-comments-long-value",
			code:     "constant T =     a   [    very_long_value_it_is ,   %e   ]  /*c*/  ;\n",
			expected: "constant T = a[\n    very_long_value_it_is, %e\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
	}

	for _, be := range []string{"be", "le"} {
		for _, tc := range testCases {
			runPrinterTestCase(
				t,
				&printerTestCase{
					name:     tc.name + "-" + be,
					code:     strings.ReplaceAll(tc.code, "%e", be),
					expected: strings.ReplaceAll(tc.expected, "%e", be),
					masked:   MaskReplaceAll(tc.masked, "%e", be),
					lineCap:  lineCap,
				},
			)
		}
	}
}

func TestPrintBitsSlice_Endiannes(t *testing.T) {
	t.Parallel()

	const lineCap = 31

	testCases := []printerTestCase{
		{
			name:     "001-clean",
			code:     "constant TEST =     arr   [  l  :   m , %e ]   ;\n",
			expected: "constant TEST = arr[l:m, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-long-name",
			code:     "constant TESTCASE =     array_30    [  l  :   m ,    %e ]   ;\n",
			expected: "constant TESTCASE = array_30[\n    l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "003-long-prefix",
			code:     "constant TESTCASE_1234 =     arr   [  l  :   m , %e ]   ;\n",
			expected: "constant TESTCASE_1234 = arr[\n    l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "004-long-value",
			code:     "constant T =     arr   [  very_long_name_it_l:m , %e ]   ;\n",
			expected: "constant T = arr[\n    very_long_name_it_l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "005-long-value-r",
			code:     "constant T =     arr   [  l  :  very_long_name_it_r , %e ]   ;\n",
			expected: "constant T = arr[\n    l:very_long_name_it_r, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "006-long-supervalue",
			code:     "constant T =     arr   [  super_very_long_name_l  :  m , %e ]   ;\n",
			expected: "constant T = arr[\n    super_very_long_name_l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "007-long-supervalue-r",
			code:     "constant T =     arr   [  l  :  super_very_long_name_r , %e ]   ;\n",
			expected: "constant T = arr[\n    l:super_very_long_name_r, %e\n];\n",
			pc:       func(f func()) { f() },
		},

		{
			name:     "008-long-supervalue-lr",
			code:     "constant T =     arr   [  super_long_very_long_name_l  :  super_long_very_long_name_r , %e ]   ;\n",
			expected: "constant T = arr[\n    super_long_very_long_name_l:super_long_very_long_name_r, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		// Comments
		{
			name:     "101-comments-short",
			code:     "constant T =   /*c*/  a    [   l  :   m , %e ]   ;\n",
			expected: "constant T = /*c*/ a[l:m, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "102-comments-short",
			code:     "constant T =     a  /*c*/  [   l  :   m , %e ]   ;\n",
			expected: "constant T = a /*c*/[l:m, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "103-comments-short",
			code:     "constant T =     a   [  /*c*/  l  :   m , %e ]   ;\n",
			expected: "constant T = a[/*c*/ l:m, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "104-comments-short",
			code:     "constant T =     a   [    l   /*c*/    :   m  , %e ]   ;\n",
			expected: "constant T = a[l /*c*/ :m, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "105-comments-short",
			code:     "constant T =     a   [    l   :   /*c*/   m , %e ]   ;\n",
			expected: "constant T = a[l: /*c*/ m, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "106-comments-short",
			code:     "constant T =     a   [    l  :   m   /*c*/ , %e ]   ;\n",
			expected: "constant T = a[l:m /*c*/, %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "107-comments-short",
			code:     "constant T =     a   [    l  :   m   ,   /*c*/   %e  ]   ;\n",
			expected: "constant T = a[l:m, /*c*/ %e];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "108-comments-short",
			code:     "constant T =     a   [    l  :   m   ,     %e   /*c*/  ]   ;\n",
			expected: "constant T = a[l:m, %e /*c*/];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "109-comments-short",
			code:     "constant T =     a   [    l  :   m   , %e  ]  /*c*/  ;\n",
			expected: "constant T = a[l:m, %e] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		// Comments - long name
		{
			name:     "111-comments-long-name",
			code:     "constant T =   /*c*/  arr_long_name    [   l  :   m ,   %e  ]   ;\n",
			expected: "constant T = /*c*/ arr_long_name[\n    l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "112-comments-long-name",
			code:     "constant T =     arr_long_name  /*c*/  [   l  :   m ,   %e  ]   ;\n",
			expected: "constant T = arr_long_name /*c*/[\n    l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "113-comments-long-name",
			code:     "constant T =     arr_long_name   [  /*c*/  l  :   m ,   %e  ]   ;\n",
			expected: "constant T = arr_long_name[\n    /*c*/ l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},

		{
			name:     "114-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  /*c*/  :      m ,   %e  ]   ;\n",
			expected: "constant T = arr_long_name[\n    l /*c*/ :m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "115-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  :   /*c*/   m ,   %e  ]   ;\n",
			expected: "constant T = arr_long_name[\n    l: /*c*/ m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "116-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  :   m   /*c*/   ,   %e  ]   ;\n",
			expected: "constant T = arr_long_name[\n    l:m /*c*/, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "117-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  :   m   ,   /*c*/   %e   ]    ;\n",
			expected: "constant T = arr_long_name[\n    l:m, /*c*/ %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "118-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  :   m   ,     %e   /*c*/ ]    ;\n",
			expected: "constant T = arr_long_name[\n    l:m, %e /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "119-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  :   m   ,   %e   ]  /*c*/  ;\n",
			expected: "constant T = arr_long_name[\n    l:m, %e\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		// Comments - long prefix
		{
			name:     "121-comments-long-prefix",
			code:     "constant TESTSUITE =   /*c*/  a    [   l  :   m ,     %e   ]   ;\n",
			expected: "constant TESTSUITE = /*c*/ a[\n    l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "122-comments-long-prefix",
			code:     "constant TESTSUITE =     a  /*c*/  [   l  :   m ,     %e   ]   ;\n",
			expected: "constant TESTSUITE = a /*c*/[\n    l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "123-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [ /*c*/   l  :   m ,     %e   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    /*c*/ l:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "124-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [   l  /*c*/    :   m ,     %e   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l /*c*/ :m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "125-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [   l    :  /*c*/   m ,     %e   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l: /*c*/ m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "126-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    l  :   m   /*c*/   ,     %e   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l:m /*c*/, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "127-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    l  :   m   ,   /*c*/   %e   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l:m, /*c*/ %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "128-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    l  :   m   ,     %e   /*c*/ ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l:m, %e /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "129-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    l  :   m   ,   %e   ]  /*c*/ ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l:m, %e\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		// Comments - long value
		{
			name:     "131-comments-long-value",
			code:     "constant T =   /*c*/  a    [   very_long_value_it:m ,   %e  ]   ;\n",
			expected: "constant T = /*c*/ a[\n    very_long_value_it:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "132-comments-long-value",
			code:     "constant T =     a  /*c*/  [   very_long_value_it:m ,   %e  ]   ;\n",
			expected: "constant T = a /*c*/[\n    very_long_value_it:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "133-comments-long-value",
			code:     "constant T =     a   [  /*c*/  very_long_value:m ,   %e  ]   ;\n",
			expected: "constant T = a[\n    /*c*/ very_long_value:m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "134-comments-long-value",
			code:     "constant T =     a   [  very_long_value    /*c*/   :  m ,   %e  ]   ;\n",
			expected: "constant T = a[\n    very_long_value /*c*/ :m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "135-comments-long-value",
			code:     "constant T =     a   [    very_long_value   :  /*c*/  m      ,   %e  ]   ;\n",
			expected: "constant T = a[\n    very_long_value: /*c*/ m, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "136-comments-long-value",
			code:     "constant T =     a   [    very_long_value:m   /*c*/   ,   %e  ]   ;\n",
			expected: "constant T = a[\n    very_long_value:m /*c*/, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "137-comments-long-value",
			code:     "constant T =     a   [    very_long_value  :  m   ,   /*c*/   %e   ]   ;\n",
			expected: "constant T = a[\n    very_long_value:m, /*c*/ %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "138-comments-long-value",
			code:     "constant T =     a   [    very_long_value   :  m,     %e   /*c*/ ]   ;\n",
			expected: "constant T = a[\n    very_long_value:m, %e /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "139-comments-long-value",
			code:     "constant T =     a   [    very_long_value   : m,   %e   ]  /*c*/  ;\n",
			expected: "constant T = a[\n    very_long_value:m, %e\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},

		{
			name:     "141-comments-long-value-l",
			code:     "constant T =   /*c*/  a    [   l:very_long_value_it ,   %e  ]   ;\n",
			expected: "constant T = /*c*/ a[\n    l:very_long_value_it, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "142-comments-long-value-l",
			code:     "constant T =     a  /*c*/  [   l:very_long_value_it ,   %e  ]   ;\n",
			expected: "constant T = a /*c*/[\n    l:very_long_value_it, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "143-comments-long-value-l",
			code:     "constant T =     a   [  /*c*/  l:very_long_value ,   %e  ]   ;\n",
			expected: "constant T = a[\n    /*c*/ l:very_long_value, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "144-comments-long-value-l",
			code:     "constant T =     a   [  l    /*c*/   :  very_long_value ,   %e  ]   ;\n",
			expected: "constant T = a[\n    l /*c*/ :very_long_value, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "145-comments-long-value-l",
			code:     "constant T =     a   [    l   :  /*c*/  very_long_value      ,   %e  ]   ;\n",
			expected: "constant T = a[\n    l: /*c*/ very_long_value, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "146-comments-long-value-l",
			code:     "constant T =     a   [    l:very_long_value   /*c*/   ,   %e  ]   ;\n",
			expected: "constant T = a[\n    l:very_long_value /*c*/, %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "147-comments-long-value-l",
			code:     "constant T =     a   [    l  :  very_long_value   ,   /*c*/   %e   ]   ;\n",
			expected: "constant T = a[\n    l:very_long_value, /*c*/ %e\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "148-comments-long-value-l",
			code:     "constant T =     a   [    l   :  very_long_value,     %e   /*c*/ ]   ;\n",
			expected: "constant T = a[\n    l:very_long_value, %e /*c*/\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "149-comments-long-value-l",
			code:     "constant T =     a   [    l   : very_long_value,   %e   ]  /*c*/  ;\n",
			expected: "constant T = a[\n    l:very_long_value, %e\n] /*c*/;\n",
			pc:       func(f func()) { f() },
		},
	}

	for _, be := range []string{"be", "le"} {
		for _, tc := range testCases {
			runPrinterTestCase(
				t,
				&printerTestCase{
					name:     tc.name + "-" + be,
					code:     strings.ReplaceAll(tc.code, "%e", be),
					expected: strings.ReplaceAll(tc.expected, "%e", be),
					masked:   MaskReplaceAll(tc.masked, "%e", be),
					lineCap:  lineCap,
				},
			)
		}
	}
}

func TestPrintBitsSlice(t *testing.T) {
	t.Parallel()

	const lineCap = 30

	testCases := []printerTestCase{
		{
			pc:       func(f func()) { f() },
			name:     "001-clean",
			code:     "constant TEST =     arr   [  l  :   m ]   ;\n",
			expected: "constant TEST = arr[l:m];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "002-long-name",
			code:     "constant TESTCASE =     array_30    [  l  :   m ]   ;\n",
			expected: "constant TESTCASE = array_30[\n    l:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "003-long-prefix",
			code:     "constant TESTCASE_1234 =     arr   [  l  :   m ]   ;\n",
			expected: "constant TESTCASE_1234 = arr[\n    l:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "004-long-value",
			code:     "constant T =     arr   [  very_long_name_it_l:m ]  ;\n",
			expected: "constant T = arr[\n    very_long_name_it_l:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "005-long-value-r",
			code:     "constant T =     arr   [  l  :  very_long_name_it_r ]  ;\n",
			expected: "constant T = arr[\n    l:very_long_name_it_r\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "006-long-supervalue-lr",
			code:     "constant T =     arr   [  super_long_very_long_name_l  :  super_long_very_long_name_r ]   ;\n",
			expected: "constant T = arr[\n    super_long_very_long_name_l:super_long_very_long_name_r\n];\n",
		},
		// Comments
		{
			pc:       func(f func()) { f() },
			name:     "101-comments-short",
			code:     "constant T =   /*c*/  arr    [   l  :   m  ]   ;\n",
			expected: "constant T = /*c*/ arr[l:m];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "102-comments-short",
			code:     "constant T =     arr  /*c*/  [   l  :   m ]   ;\n",
			expected: "constant T = arr /*c*/[l:m];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "103-comments-short",
			code:     "constant T =     arr   [  /*c*/  l  :   m ]   ;\n",
			expected: "constant T = arr[/*c*/ l:m];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "104-comments-short",
			code:     "constant T =     arr   [    l   /*c*/    :   m ]   ;\n",
			expected: "constant T = arr[l /*c*/ :m];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "105-comments-short",
			code:     "constant T =     arr   [    l   :   /*c*/   m ]   ;\n",
			expected: "constant T = arr[l: /*c*/ m];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "106-comments-short",
			code:     "constant T =     arr   [    l  :   m   /*c*/ ]   ;\n",
			expected: "constant T = arr[l:m /*c*/];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "107-comments-short",
			code:     "constant T =     arr   [    l  :   m   ]  /*c*/  ;\n",
			expected: "constant T = arr[l:m] /*c*/;\n",
		},
		// Comments - long name
		{
			pc:       func(f func()) { f() },
			name:     "111-comments-long-name",
			code:     "constant T =   /*c*/  arr_long_name    [   l  :   m   ]   ;\n",
			expected: "constant T = /*c*/ arr_long_name[\n    l:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "112-comments-long-name",
			code:     "constant T =     arr_long_name  /*c*/  [   l  :   m   ]   ;\n",
			expected: "constant T = arr_long_name /*c*/[\n    l:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "113-comments-long-name",
			code:     "constant T =     arr_long_name   [  /*c*/  l  :   m   ]   ;\n",
			expected: "constant T = arr_long_name[\n    /*c*/ l:m\n];\n",
		},

		{
			pc:       func(f func()) { f() },
			name:     "114-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  /*c*/  :      m  ]   ;\n",
			expected: "constant T = arr_long_name[\n    l /*c*/ :m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "115-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  :   /*c*/   m  ]   ;\n",
			expected: "constant T = arr_long_name[\n    l: /*c*/ m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "116-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  :   m   /*c*/   ]   ;\n",
			expected: "constant T = arr_long_name[\n    l:m /*c*/\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "117-comments-long-name",
			code:     "constant T =     arr_long_name   [    l  :   m   ]  /*c*/  ;\n",
			expected: "constant T = arr_long_name[\n    l:m\n] /*c*/;\n",
		},
		// Comments - long prefix
		{
			pc:       func(f func()) { f() },
			name:     "121-comments-long-prefix",
			code:     "constant TESTSUITE =   /*c*/  a    [   l  :   m   ]   ;\n",
			expected: "constant TESTSUITE = /*c*/ a[\n    l:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "122-comments-long-prefix",
			code:     "constant TESTSUITE =     a  /*c*/  [   l  :   m   ]   ;\n",
			expected: "constant TESTSUITE = a /*c*/[\n    l:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "123-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [ /*c*/   l  :   m   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    /*c*/ l:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "124-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [   l  /*c*/    :   m   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l /*c*/ :m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "125-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [   l    :  /*c*/   m   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l: /*c*/ m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "126-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    l  :   m   /*c*/   ]   ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l:m /*c*/\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "127-comments-long-prefix",
			code:     "constant TESTSUITE123456 =    a    [    l  :   m   ]  /*c*/ ;\n",
			expected: "constant TESTSUITE123456 = a[\n    l:m\n] /*c*/;\n",
		},
		// Comments - long value
		{
			pc:       func(f func()) { f() },
			name:     "131-comments-long-value",
			code:     "constant T =   /*c*/  a    [   very_long_value_it:m   ]   ;\n",
			expected: "constant T = /*c*/ a[\n    very_long_value_it:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "132-comments-long-value",
			code:     "constant T =     a  /*c*/  [   very_long_value_it:m   ]   ;\n",
			expected: "constant T = a /*c*/[\n    very_long_value_it:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "133-comments-long-value",
			code:     "constant T =     a   [  /*c*/  very_long_value:m   ]   ;\n",
			expected: "constant T = a[\n    /*c*/ very_long_value:m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "134-comments-long-value",
			code:     "constant T =     a   [  very_long_value    /*c*/   :  m   ]   ;\n",
			expected: "constant T = a[\n    very_long_value /*c*/ :m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "135-comments-long-value",
			code:     "constant T =     a   [    very_long_value   :  /*c*/  m   ]   ;\n",
			expected: "constant T = a[\n    very_long_value: /*c*/ m\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "136-comments-long-value",
			code:     "constant T =     a   [    very_long_value:m   /*c*/   ]   ;\n",
			expected: "constant T = a[\n    very_long_value:m /*c*/\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "137-comments-long-value",
			code:     "constant T =     a   [    very_long_value   : m  ]  /*c*/  ;\n",
			expected: "constant T = a[\n    very_long_value:m\n] /*c*/;\n",
		},

		{
			pc:       func(f func()) { f() },
			name:     "141-comments-long-value-l",
			code:     "constant T =   /*c*/  a    [   l:very_long_value_it   ]   ;\n",
			expected: "constant T = /*c*/ a[\n    l:very_long_value_it\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "142-comments-long-value-l",
			code:     "constant T =     a  /*c*/  [   l:very_long_value_it   ]   ;\n",
			expected: "constant T = a /*c*/[\n    l:very_long_value_it\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "143-comments-long-value-l",
			code:     "constant T =     a   [  /*c*/  l:very_long_value   ]   ;\n",
			expected: "constant T = a[\n    /*c*/ l:very_long_value\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "144-comments-long-value-l",
			code:     "constant T =     a   [  l    /*c*/   :  very_long_value  ]   ;\n",
			expected: "constant T = a[\n    l /*c*/ :very_long_value\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "145-comments-long-value-l",
			code:     "constant T =     a   [    l   :  /*c*/  very_long_value     ]   ;\n",
			expected: "constant T = a[\n    l: /*c*/ very_long_value\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "146-comments-long-value-l",
			code:     "constant T =     a   [    l:very_long_value   /*c*/    ]   ;\n",
			expected: "constant T = a[\n    l:very_long_value /*c*/\n];\n",
		},
		{
			pc:       func(f func()) { f() },
			name:     "147-comments-long-value-l",
			code:     "constant T =     a   [    l   : very_long_value  ]  /*c*/  ;\n",
			expected: "constant T = a[\n    l:very_long_value\n] /*c*/;\n",
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintCompiletimeArray_basic(t *testing.T) {

	t.Parallel()

	const lineCap = 30

	testCases := []printerTestCase{
		// TODO: Use `param` instead of `constant` to have DML 1.4 true syntax
		{
			name:     "001-list-short",
			code:     "constant T = [  a,b,c,d ]   ;\n",
			expected: "constant T = [a, b, c, d];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-list-long",
			code:     "constant T =    [     aaaa,  bbbb  ,     cccc,  dddd ]   ;\n",
			expected: "constant T = [\n    aaaa, bbbb, cccc, dddd\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "003-list-superlong-all-items",
			code:     "constant T =    [     aaaaaaaaaaaaaaa,  bbbbbbbbbbbbbbbbbbbb  ,     ccccccccccccccccc,  dddddddddddddddddddd ]   ;\n",
			expected: "constant T = [\n    aaaaaaaaaaaaaaa,\n    bbbbbbbbbbbbbbbbbbbb,\n    ccccccccccccccccc,\n    dddddddddddddddddddd\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "004-list-superlong-any-item",
			code:     "constant T =    [     aaaa,  bbbb  ,     ccccccccccccccccc,  dddd ]   ;\n",
			expected: "constant T = [\n    aaaa,\n    bbbb,\n    ccccccccccccccccc,\n    dddd\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "101-list-short-keep-new-line",
			code:     "constant T = [  \n a,b,c,d ]   ;\n",
			expected: "constant T = [\n    a, b, c, d\n];\n",
			masked: MaskInputNewlinePreservation(
				"constant T = [a, b, c, d];\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "102-list-short-keep-new-lines",
			// more than 3 new lines then use multiline
			code:     "constant T = [  \n a,\n b,\nc,d\n]   ;\n",
			expected: "constant T = [\n    a,\n    b,\n    c,\n    d\n];\n",
			masked: MaskInputNewlinePreservation(
				"constant T = [a, b, c, d];\n",
			),
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintCompiletimeArray_innerExpr(t *testing.T) {

	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		// TODO: Use `param` instead of `constant` to have DML 1.4 true syntax
		{
			name:     "001-list-short",
			code:     "constant T = [  a+b,  c*   d, e-f ]   ;\n",
			expected: "constant T = [a + b, c * d, e - f];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-list-long",
			code:     "constant T =    [     aaa +	  bb  ,     cc -  dd, ee    *   fff ]   ;\n",
			expected: "constant T = [\n    aaa + bb, cc - dd, ee * fff\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "003-list-superlong-all-items",
			code:     "constant T =    [     aaaaaaa + aaaaaaaa,  bbbbbbbbbb- bbbbbbbbbb  ,     ccccccccc*cccccccc,  dddddddddd+dddddddd ]   ;\n",
			expected: "constant T = [\n    aaaaaaa + aaaaaaaa,\n    bbbbbbbbbb - bbbbbbbbbb,\n    ccccccccc * cccccccc,\n    dddddddddd + dddddddd\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "004-list-superlong-any-item",
			code:     "constant T =    [     aaaa + bbbb  ,cccccccccccccccccccccccccccc,  dd + dddd ]   ;\n",
			expected: "constant T = [\n    aaaa + bbbb,\n    cccccccccccccccccccccccccccc,\n    dd + dddd\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "101-list-short-keep-new-line",
			code:     "constant T = [  aa + bb, \ncc - dd, ee * ff]   ;\n",
			expected: "constant T = [\n    aa + bb, cc - dd, ee * ff\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "102-list-short-keep-new-lines",
			// more than 2 new lines then use multiline
			code:     "constant T = [  aa + bb, cc \n - dd, ee \n * ff \n]   ;\n",
			expected: "constant T = [\n    aa + bb,\n    cc - dd,\n    ee * ff\n];\n",
			masked: MaskInputNewlinePreservation(
				"constant T = [\n    aa + bb, cc - dd, ee * ff\n];\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "201-list-superlong-all-items-paren",
			code:     "constant T =    [     aaaaaaa + aaaaaaaa ^  bbbbbbbbbb- bbbbbbbbbb  ,     (ccc*ccc     +   ccc),  dddddddddd+dddddddd ]   ;\n",
			expected: "constant T = [\n    aaaaaaa + aaaaaaaa\n    ^ bbbbbbbbbb - bbbbbbbbbb,\n    (ccc * ccc + ccc),\n    dddddddddd + dddddddd\n];\n",
			pc:       func(f func()) { f() },
		},

		{
			name: "301-nested-list",
			// Always split when nested list
			code:     "constant T = [aaaa, [bbbb, cccc], dddd];\n",
			expected: "constant T = [\n    aaaa,\n    [bbbb, cccc],\n    dddd\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "301-nested-list-2levels",
			// Always split when nested list
			code:     "constant T = [aaaa, [bbbb, [cccc, dddd]], eeee];\n",
			expected: "constant T = [\n    aaaa,\n    [\n        bbbb,\n        [cccc, dddd]\n    ],\n    eeee\n];\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintCompiletimeArray_comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "001-comment-before",
			code:     "constant T = // comment\n [a, b, c];\n",
			expected: "constant T = // comment\n    [a, b, c];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-comment-after-open",
			code:     "constant T  =  [   // comment\n  a,   b,c];\n",
			expected: "constant T = [ // comment\n    a, b, c\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "003-comment-between-elements",
			code:     "constant T = [a, // comment\n b, c];\n",
			expected: "constant T = [\n    a, // comment\n    b,\n    c\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "004-comment-after-last-element",
			code:     "constant T = [a, b, c // comment\n];\n",
			expected: "constant T = [\n    a, b, c // comment\n];\n",
			masked: MaskArgumentOverWrappingWithLineComments(
				"constant T = [\n    a,\n    b,\n    c // comment\n];\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "005-comment-after-close",
			code:     "constant T = [a, b, c] // comment\n;\n",
			expected: "constant T = [a, b, c] // comment\n    ;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "101-block-comment-inline",
			code:     "constant T = [a, /* comment */ b, c];\n",
			expected: "constant T = [a, /* comment */ b, c];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "201-comment-inside-expression",
			code:     "constant T = [a // comment\n + b, c];\n",
			expected: "constant T = [\n    a // comment\n    + b,\n    c\n];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "202-block-comment-inside-expression",
			code:     "constant T = [a + /* comment */ b, c];\n",
			expected: "constant T = [a + /* comment */ b, c];\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionBinaryOp_DeepNesting(t *testing.T) {
	t.Parallel()

	const n = 500
	var inputBuilder, expectedBuilder strings.Builder
	inputBuilder.WriteString("constant TESTCASE = x0")
	expectedBuilder.WriteString("constant TESTCASE = x0")
	for i := 1; i < n; i++ {
		fmt.Fprintf(&inputBuilder, "+x%d", i)
		fmt.Fprintf(&expectedBuilder, " + x%d", i)
	}
	inputBuilder.WriteString(";\n")
	expectedBuilder.WriteString(";\n")

	var mem runtime.MemStats
	GetMemStats(&mem)
	runPrinterTestCase(
		t,
		&printerTestCase{
			name:     "deep-nesting-500-additions",
			code:     inputBuilder.String(),
			expected: expectedBuilder.String(),
			lineCap:  100000,
		},
	)
	var memFinal runtime.MemStats
	LogMemDiff(t, "end of the test", &mem, &memFinal)
}
