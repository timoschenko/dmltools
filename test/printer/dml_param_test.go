// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintParamDecl(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		// 1. Basic Definitions
		{
			name:     "01-basic",
			code:     "   param    x  ;\n",
			expected: "param x;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-typed",
			code:     "param   x  :   int  ;\n",
			expected: "param x: int;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-assigned",
			code:     "param   x   =   1  ;\n",
			expected: "param x = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-default",
			code:     "param   x   default    1  ;\n",
			expected: "param x default 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-typed-and-assigned",
			code:     "param x :   int    =   1;\n",
			expected: "param x: int = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-typed-and-default",
			code:     "param x : int default 1;\n",
			expected: "param x: int default 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-auto",
			code:     "param x auto;\n",
			expected: "param x auto;\n",
			pc:       func(f func()) { f() },
		},
		// 2. Line Wrapping (Line Cap = 40)
		{
			name:     "08-long-identifier",
			code:     "param this_is_a_very_long_param_name_exceeding_cap;\n",
			expected: "param this_is_a_very_long_param_name_exceeding_cap;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-long-type-decl",
			code:     "param x : this_is_a_very_long_type_name_exceeding_cap;\n",
			expected: "param x\n    : this_is_a_very_long_type_name_exceeding_cap;\n",
			pc:       func(f func()) { f() },
		},
		{
			// TODO: add Synthetick ()?
			name:     "10-long-expression",
			code:     "param x = this_is_a_very_long_expression_exceeding_cap;\n",
			expected: "param x = this_is_a_very_long_expression_exceeding_cap;\n",
			pc:       func(f func()) { f() },
		},
		{
			// TODO: add Synthetick ()?
			name:     "11-long-default-expression",
			code:     "param x default this_is_a_very_long_expression_exceeding_cap;\n",
			expected: "param x default this_is_a_very_long_expression_exceeding_cap;\n",
			pc:       func(f func()) { f() },
		},
		{
			// TODO: add Synthetick ()?
			name:     "12-long-everything",
			code:     "param this_is_a_very_long_param_name : this_is_a_very_long_type_name = this_is_a_very_long_expression_exceeding_cap;\n",
			expected: "param this_is_a_very_long_param_name\n    : this_is_a_very_long_type_name = this_is_a_very_long_expression_exceeding_cap;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-long-auto",
			code:     "param this_is_a_very_long_param_name_exceeding_cap auto;\n",
			expected: "param this_is_a_very_long_param_name_exceeding_cap\n    auto;\n",
			pc:       func(f func()) { f() },
		},
		// 3. Complex Type Annotations
		{
			name:     "14-complex-type-short",
			code:     "param x : const short_t *;\n",
			expected: "param x: const short_t *;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-complex-type-long",
			code:     "param x : const very_long_typename_exceeding_line_cap_t *;\n",
			expected: "param x\n    : const very_long_typename_exceeding_line_cap_t *;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-complex-type-long-with-assign",
			code:     "param x : const very_long_typename_exceeding_line_cap_t * = 1;\n",
			expected: "param x\n    : const very_long_typename_exceeding_line_cap_t * = 1;\n",
			pc:       func(f func()) { f() },
		},
		// 4. Complex Expressions
		{
			name:     "17-complex-expr-short",
			code:     "param x = cast(a, int);\n",
			expected: "param x = cast(a, int);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-complex-expr-long",
			code:     "param x = cast(this_is_a_very_long_variable_name_exceeding_cap, int);\n",
			expected: "param x = cast(\n    this_is_a_very_long_variable_name_exceeding_cap,\n    int\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "19-default-parenthesized",
			code:     "param name default (defined x) #? true #: false;\n",
			expected: "param name default (defined x) #? true #: false;\n",
			lineCap:  80,
			pc:       func(f func()) { f() },
		},
		{
			name:     "20-default-new-parenthesized",
			code:     "param name :default (defined x) #? true #: false;\n",
			expected: "param name :default (defined x) #? true #: false;\n",
			lineCap:  80,
			pc:       func(f func()) { f() },
		},
		{
			name:     "21-default-simple-parenthesized",
			code:     "param x default (1 + 2);\n",
			expected: "param x default (1 + 2);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "22-default-new-simple-parenthesized",
			code:     "param x :default (1 + 2);\n",
			expected: "param x :default (1 + 2);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "23-default-new-basic",
			code:     "param   x   :default    1  ;\n",
			expected: "param x :default 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "24-typed-and-default-new",
			code:     "param x : int :default 1;\n",
			expected: "param x: int :default 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "25-long-default-new-expression",
			code:     "param x :default this_is_a_very_long_expression_exceeding_cap;\n",
			expected: "param x :default this_is_a_very_long_expression_exceeding_cap;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "26-long-typed-default-new",
			code: "param this_is_a_very_long_param_name : " +
				"this_is_a_very_long_type_name :default " +
				"this_is_a_very_long_expression_exceeding_cap;\n",
			expected: "param this_is_a_very_long_param_name\n" +
				"    : this_is_a_very_long_type_name :default " +
				"this_is_a_very_long_expression_exceeding_cap;\n",
			pc: func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintParamDeclComments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		// 1. End-of-line Comments
		{
			name:     "01-short-comment-after-semi",
			code:     "param x; // short comment\n",
			expected: "param x; // short comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-long-comment-after-semi",
			code:     "param x; // this is a very long comment that exceeds the line cap completely\n",
			expected: "param x; // this is a very long comment that exceeds the line cap completely\n",
			pc:       func(f func()) { f() },
		},
		// 2. Inline Block Comments
		{
			name:     "03-inline-before-ident",
			code:     "param /* c1 */ x;\n",
			expected: "param /* c1 */ x;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-inline-after-ident",
			code:     "param x /* c2 */ : int;\n",
			expected: "param x /* c2 */ : int;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-inline-after-colon",
			code:     "param x : /* c3 */ int;\n",
			expected: "param x: /* c3 */ int;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-inline-after-type",
			code:     "param x : int /* c4 */ = 1;\n",
			expected: "param x: int /* c4 */ = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-inline-after-assign",
			code:     "param x : int = /* c5 */ 1;\n",
			expected: "param x: int = /* c5 */ 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-inline-before-auto",
			code:     "param x /* c6 */ auto;\n",
			expected: "param x /* c6 */ auto;\n",
			pc:       func(f func()) { f() },
		},
		// 3. Multiline/Interleaved Comments
		{
			name:     "09-line-comments-interleaved",
			code:     "param // c2\nx // c3\n: // c4\nint // c5\n= // c6\n1 // c7\n;\n",
			expected: "param // c2\n    x // c3\n    : // c4\n    int // c5\n    = // c6\n    1 // c7\n    ;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-block-comments-interleaved",
			code:     "param /* c1 */ x /* c2 */ : /* c3 */ int /* c4 */ = /* c5 */ 1 /* c6 */ ;\n",
			expected: "param /* c1 */ x /* c2 */ : /* c3 */ int /* c4 */ = /* c5 */ 1 /* c6 */;\n",
			pc:       func(f func()) { f() },
		},
		// 4. Complex Type Annotations with Comments
		{
			name:     "11-complex-type-comment-before-colon",
			code:     "param x /* c1 */ : const short_t *;\n",
			expected: "param x /* c1 */ : const short_t *;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-complex-type-comment-after-colon",
			code:     "param x : /* c2 */ const short_t *;\n",
			expected: "param x: /* c2 */ const short_t *;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-complex-type-comment-after-type",
			code:     "param x : const short_t * /* c3 */ = 1;\n",
			expected: "param x: const short_t * /* c3 */ = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-complex-type-long-interleaved-comments",
			code:     "param /* c1 */ x /* c2 */ : /* c3 */ const very_long_typename_exceeding_line_cap_t * /* c4 */ = /* c5 */ 1 /* c6 */ ;\n",
			expected: "param /* c1 */ x /* c2 */ : /* c3 */ const very_long_typename_exceeding_line_cap_t * /* c4 */ = /* c5 */ 1 /* c6 */;\n",
			pc:       func(f func()) { f() },
		},
		// 5. Complex Expressions with Comments
		{
			name:     "15-complex-expr-comment-inline",
			code:     "param x = cast(a /* c1 */, int);\n",
			expected: "param x = cast(a /* c1 */, int);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-complex-expr-comment-long",
			code:     "param x = cast(this_is_a_very_long_variable_name_exceeding_cap /* c1 */, int);\n",
			expected: "param x = cast(\n    this_is_a_very_long_variable_name_exceeding_cap /* c1 */,\n    int\n);\n",
			pc:       func(f func()) { f() },
		},
		// 6. Default and DefaultNew with Comments
		{
			name:     "17-inline-comment-default",
			code:     "param x /* c1 */ default /* c2 */ 1;\n",
			expected: "param x /* c1 */ default /* c2 */ 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-inline-comment-default-new",
			code:     "param x /* c1 */ :default /* c2 */ 1;\n",
			expected: "param x /* c1 */ :default /* c2 */ 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "19-inline-comment-typed-default",
			code:     "param x : int /* c */ default 1;\n",
			expected: "param x: int /* c */ default 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "20-inline-comment-typed-default-new",
			code:     "param x : int /* c */ :default 1;\n",
			expected: "param x: int /* c */ :default 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "21-line-comment-default",
			code: "param // c1\nx // c2\n: // c3\nint // c4\ndefault // c5\n1 // c6\n;\n",
			expected: "param // c1\n    x // c2\n    : // c3\n" +
				"    int // c4\n    default // c5\n    1 // c6\n    ;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "22-line-comment-default-new",
			code: "param // c1\nx // c2\n: // c3\nint // c4\n:default // c5\n1 // c6\n;\n",
			expected: "param // c1\n    x // c2\n    : // c3\n" +
				"    int // c4\n    :default // c5\n    1 // c6\n    ;\n",
			pc: func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
