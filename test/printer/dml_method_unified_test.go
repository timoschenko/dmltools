// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintMethodUnified(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-simple-method",
			code:     "method    foo (  )  ;\n",
			expected: "method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-method-with-body",
			code:     "method   foo(  )   {  }  \n",
			expected: "method foo() {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-method-with-qualifiers",
			code:     "shared    method   foo  (  ) ; \n",
			expected: "shared method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-method-with-params-and-throws",
			code:     "method    foo  (   int   a  )   throws  ;\n",
			expected: "method foo(int a) throws;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-method-with-return-type",
			code:     "method   foo( )   ->   ( bool  );\n",
			expected: "method foo() -> (bool);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-method-with-return-type-and-body",
			code:     "method foo()->( bool) {}\n",
			expected: "method foo() -> (bool) {\n}\n",
			pc:       func(f func()) { f() },
		},
		// Wrapping Strategy for Long Qualifiers:
		// When methodQualifiers + method name exceeds the line cap (e.g. 40 chars),
		// it wraps BEFORE the 'method' keyword WITHOUT any extra indentation,
		// similar to old C/C++ template formatting.
		// Example: "shared independent startup memoized\nmethod foo();"
		{
			name:     "07-independent-method",
			code:     "  independent   method foo();\n",
			expected: "independent method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-independent-startup-method",
			code:     "  independent   startup   method   foo();\n",
			expected: "independent startup method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-shared-independent-startup-method",
			code:     "shared    independent    startup    method foo();\n",
			expected: "shared independent startup method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-independent-startup-memoized-method",
			code:     "independent    startup   memoized  method foo();\n",
			expected: "independent startup memoized\nmethod foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-shared-independent-startup-memoized-method",
			code:     "shared   independent    startup   memoized    method foo();\n",
			expected: "shared independent startup memoized\nmethod foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-inline-method",
			code:     "   inline   method foo();\n",
			expected: "inline method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-long-method-independent-startup",
			code:     "independent startup method foo_with_a_very_long_name();\n",
			expected: "independent startup\nmethod foo_with_a_very_long_name();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-long-method-shared",
			code:     "shared method another_extremely_long_function_name();\n",
			expected: "shared\nmethod another_extremely_long_function_name();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-method-default",
			code:     "method foo() default {}\n",
			expected: "method foo() default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15a-method-default-new",
			code:     "method foo() :default {}\n",
			expected: "method foo() :default {\n}\n",
			pc:       func(f func()) { f() },
		},
		// Brace Placement Strategy:
		// Use K&R brace placement only if the 'method' token and the opening '{'
		// are on the same line. If the 'method' token gets wrapped relative to '{'
		// (e.g. due to long return types or throws), the outermost braces must
		// use the Allman style (opening brace on a new line).
		{
			name:     "16-method-with-return-default",
			code:     "method foo() -> (bar) default {}\n",
			expected: "method foo() -> (bar) default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16a-method-with-return-default-new",
			code:     "method foo() -> (bar) :default {}\n",
			expected: "method foo() -> (bar) :default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-long-method-throws-default",
			code:     "method foo_with_a_very_long_name() throws default {}\n",
			expected: "method foo_with_a_very_long_name() throws default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17a-long-method-throws-default-new",
			code:     "method foo_with_a_very_long_name() throws :default {}\n",
			expected: "method foo_with_a_very_long_name() throws :default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-long-method-throws",
			code:     "method another_extremely_long_function_name() throws;\n",
			expected: "method another_extremely_long_function_name() throws;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "19-method-with-long-return-type-throws-default",
			code: "method function_with_long_return_type() -> (very_long_return_type) throws default {}\n",
			expected: "method function_with_long_return_type() -> (\n" +
				"    very_long_return_type\n" +
				") throws default {\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "19a-method-with-long-return-type-throws-default-new",
			code: "method function_with_long_return_type() -> (very_long_return_type) throws :default {}\n",
			expected: "method function_with_long_return_type() -> (\n" +
				"    very_long_return_type\n" +
				") throws :default {\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "20-simple-method-return",
			code:     "method foo() -> (bool);\n",
			expected: "method foo() -> (bool);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "21-method-return-ptr",
			code:     "method foo()    -> (const    foo_t  **    );\n",
			expected: "method foo() -> (const foo_t **);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "22-long-method-return",
			code:     "method foo_with_a_very_long_name() -> (bool);\n",
			expected: "method foo_with_a_very_long_name() -> (\n    bool\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "23-method-return-wrap-elements-inline",
			code:     "method foo_long_name() -> (a, b, c, d) throws {}\n",
			expected: "method foo_long_name() -> (\n    a, b, c, d\n) throws {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "24-method-return-wrap-elements-wrap",
			code: "method foo() -> (very_long_return_type_1, very_long_return_type_2, very_long_return_type_3) throws default {}\n",
			expected: "method foo() -> (\n" +
				"    very_long_return_type_1,\n" +
				"    very_long_return_type_2,\n" +
				"    very_long_return_type_3\n" +
				") throws default {\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "25-method-params-no-wrap",
			code:     "method foo(bool a);\n",
			expected: "method foo(bool a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "26-method-params-top-level-explode",
			code:     "method foo(bool a, foo_t b, foo_t ** const c);\n",
			expected: "method foo(\n    bool a, foo_t b, foo_t **const c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "26-method-params-all-level-explode",
			code:     "method foo(bool a, foo_t b, const foo_t ** const c);\n",
			expected: "method foo(\n    bool a,\n    foo_t b,\n    const foo_t **const c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "27-method-params-all-explode",
			code:     "method foo(bool a, foo_t b, void (*cb_with_a_very_long_name)(int, bool, string));\n",
			expected: "method foo(\n    bool a,\n    foo_t b,\n    void (*cb_with_a_very_long_name)(\n        int, bool, string\n    )\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "28-method-params-default",
			code:     "method foo(int x) default {}\n",
			expected: "method foo(int x) default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "29-method-params-default-new",
			code:     "method foo(int x) :default {}\n",
			expected: "method foo(int x) :default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "30-shared-method-default",
			code:     "shared method foo() default {}\n",
			expected: "shared method foo() default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "31-shared-method-default-new",
			code:     "shared method foo() :default {}\n",
			expected: "shared method foo() :default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "32-method-with-body-default",
			code:     "method foo() default {\nx = 1;\n}\n",
			expected: "method foo() default {\n    x = 1;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "33-method-with-body-default-new",
			code:     "method foo() :default {\nx = 1;\n}\n",
			expected: "method foo() :default {\n    x = 1;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "34-method-with-super-call-default-no-args",
			code:     "method foo() {\ndefault();\n}\n",
			expected: "method foo() {\n    default();\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "35-method-with-super-call-default-with-args",
			code:     "method foo(int x) {\ndefault(x);\n}\n",
			expected: "method foo(int x) {\n    default(x);\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "36-method-with-super-call-default-multi-args-assign",
			code:     "method foo(int a, int b) -> (int) {\nx = default(a, b);\n}\n",
			expected: "method foo(int a, int b) -> (int) {\n    x = default(a, b);\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "37-method-default-with-super-call-default",
			code:     "method foo() default {\ndefault();\n}\n",
			expected: "method foo() default {\n    default();\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "38-method-default-new-with-super-call-default",
			code:     "method foo() :default {\ndefault();\n}\n",
			expected: "method foo() :default {\n    default();\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "39-method-with-super-call-default-return",
			code:     "method foo() -> (int) {\nreturn default();\n}\n",
			expected: "method foo() -> (int) {\n    return default();\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "40-method-with-super-call-default-assign",
			code:     "method foo() -> (int) {\nx = default();\n}\n",
			expected: "method foo() -> (int) {\n    x = default();\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintMethodUnified_InlineComments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-before-method",
			code:     "dml 1.4;\n\n\n/*c*/ method foo();\n",
			expected: "dml 1.4;\n\n\n/*c*/ method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-after-method",
			code:     "method /*c*/ foo();\n",
			expected: "method /*c*/ foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-before-qualifier",
			code:     "dml 1.4;\n\n\n/*c*/ shared method foo();\n",
			expected: "dml 1.4;\n\n\n/*c*/ shared method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-after-qualifier",
			code:     "shared /*c*/ method foo();\n",
			expected: "shared /*c*/ method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-between-name-and-params",
			code:     "method foo /*c*/ ();\n",
			expected: "method foo /*c*/ ();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-inside-params-before-type",
			code:     "method foo( /*c*/ int a );\n",
			expected: "method foo(/*c*/ int a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-inside-params-after-type",
			code:     "method foo( int /*c*/ a );\n",
			expected: "method foo(int /*c*/ a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-inside-params-after-name",
			code:     "method foo( int a /*c*/ );\n",
			expected: "method foo(int a /*c*/);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-inside-params-around-comma",
			code:     "method foo( int a /*c*/, /*d*/ int b );\n",
			expected: "method foo(int a /*c*/, /*d*/ int b);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-before-throws",
			code:     "method foo() /*c*/ throws;\n",
			expected: "method foo() /*c*/ throws;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-after-throws",
			code:     "method foo() throws /*c*/ ;\n",
			expected: "method foo() throws /*c*/;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-before-default",
			code:     "method foo() /*c*/ default {}\n",
			expected: "method foo() /*c*/ default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12a-before-default-new",
			code:     "method foo() /*c*/ :default {}\n",
			expected: "method foo() /*c*/ :default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-after-default",
			code:     "method foo() default /*c*/ {}\n",
			expected: "method foo() default /*c*/ {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13a-after-default-new",
			code:     "method foo() :default /*c*/ {}\n",
			expected: "method foo() :default /*c*/ {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-before-return-arrow",
			code:     "method foo() /*c*/ -> (bool);\n",
			expected: "method foo() /*c*/ -> (bool);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-after-return-arrow",
			code:     "method foo() -> /*c*/ (bool);\n",
			expected: "method foo() -> /*c*/ (bool);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-inside-return-list",
			code:     "method foo() -> ( /*c*/ bool /*d*/ );\n",
			expected: "method foo() -> (/*c*/ bool /*d*/);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-before-body-brace",
			code:     "method foo() /*c*/ {}\n",
			expected: "method foo() /*c*/ {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-inside-body-braces",
			code:     "method foo() { /*c*/ }\n",
			expected: "method foo() {\n    /*c*/\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "19-after-body-brace",
			code:     "method foo() {} /*c*/\n",
			expected: "method foo() {\n} /*c*/\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "20-wrapped-long-method",
			code:     "dml 1.4;\n\n\n/*c*/ independent startup /*d*/ method foo_with_a_very_long_name();\n",
			expected: "dml 1.4;\n\n\n/*c*/ independent startup /*d*/\nmethod foo_with_a_very_long_name();\n",
			masked: MaskBlockCommentWrapping(
				"dml 1.4;\n\n\n/*c*/ independent startup /*d*/ method foo_with_a_very_long_name();\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "21-wrapped-long-return",
			code:     "method foo_with_a_very_long_name() /*c*/ -> (bool) /*d*/ ;\n",
			expected: "method foo_with_a_very_long_name() /*c*/ -> (\n    bool\n) /*d*/;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "22-wrapped-throws-default",
			code:     "method foo_with_a_very_long_name() /*c*/ throws /*d*/ default /*e*/ {}\n",
			expected: "method foo_with_a_very_long_name() /*c*/ throws /*d*/ default /*e*/ {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "22a-wrapped-throws-default-new",
			code:     "method foo_with_a_very_long_name() /*c*/ throws /*d*/ :default /*e*/ {}\n",
			expected: "method foo_with_a_very_long_name() /*c*/ throws /*d*/ :default /*e*/ {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "23-wrapped-params",
			code:     "method foo( /*a*/ bool a /*b*/, /*c*/ foo_t b /*d*/, /*e*/ foo_t ** const c /*f*/ );\n",
			expected: "method foo(\n    /*a*/ bool a /*b*/,\n    /*c*/ foo_t b /*d*/,\n    /*e*/ foo_t **const c /*f*/\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "24-wrapped-return-list",
			code: "method foo() -> ( /*a*/ very_long_return_type_1 /*b*/, /*c*/ very_long_return_type_2 /*d*/ ) throws default {}\n",
			expected: "method foo() -> (\n" +
				"    /*a*/ very_long_return_type_1 /*b*/,\n" +
				"    /*c*/ very_long_return_type_2 /*d*/\n" +
				") throws default {\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "24a-wrapped-return-list-default-new",
			code: "method foo() -> ( /*a*/ very_long_return_type_1 /*b*/, /*c*/ very_long_return_type_2 /*d*/ ) throws :default {}\n",
			expected: "method foo() -> (\n" +
				"    /*a*/ very_long_return_type_1 /*b*/,\n" +
				"    /*c*/ very_long_return_type_2 /*d*/\n" +
				") throws :default {\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "25-super-call-default-inline-comments",
			code:     "method foo() {\n/*a*/ default( /*b*/ x /*c*/ ) /*d*/ ;\n}\n",
			expected: "method foo() {\n    /*a*/ default(/*b*/ x /*c*/) /*d*/;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "25a-super-call-default-no-args-inline-comments",
			code:     "method foo() {\n/*a*/ default( /*b*/ ) /*c*/ ;\n}\n",
			expected: "method foo() {\n    /*a*/ default(/*b*/) /*c*/;\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintMethodUnified_LineComments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-before-method",
			code:     "dml 1.4;\n\n\n//c\n method foo();\n",
			expected: "dml 1.4;\n\n\n//c\nmethod foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-after-method",
			code:     "method //c\n foo();\n",
			expected: "method //c\nfoo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-before-qualifier",
			code:     "dml 1.4;\n\n\n//c\n shared method foo();\n",
			expected: "dml 1.4;\n\n\n//c\nshared method foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-after-qualifier",
			code:     "shared //c\n method foo();\n",
			expected: "shared //c\nmethod foo();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-between-name-and-params",
			code:     "method foo //c\n ();\n",
			expected: "method foo //c\n();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-inside-params-before-type",
			code:     "method foo( //c\n int a );\n",
			expected: "method foo( //c\n    int a\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-inside-params-after-type",
			code:     "method foo( int //c\n a );\n",
			expected: "method foo(\n    int //c\n    a\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-inside-params-after-name",
			code:     "method foo( int a //c\n );\n",
			expected: "method foo(\n    int a //c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-inside-params-around-comma",
			code:     "method foo( int a //c\n, //d\n int b );\n",
			expected: "method foo(\n    int a //c\n    , //d\n    int b\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-before-throws",
			code:     "method foo() //c\n throws;\n",
			expected: "method foo() //c\nthrows;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-after-throws",
			code:     "method foo() throws //c\n ;\n",
			expected: "method foo() throws //c\n;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-before-default",
			code:     "method foo() //c\n default {}\n",
			expected: "method foo() //c\ndefault {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12a-before-default-new",
			code:     "method foo() //c\n :default {}\n",
			expected: "method foo() //c\n:default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-after-default",
			code:     "method foo() default //c\n {}\n",
			expected: "method foo() default //c\n{\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13a-after-default-new",
			code:     "method foo() :default //c\n {}\n",
			expected: "method foo() :default //c\n{\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-before-return-arrow",
			code:     "method foo() //c\n -> (bool);\n",
			expected: "method foo() //c\n-> (bool);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-after-return-arrow",
			code:     "method foo() -> //c\n (bool);\n",
			expected: "method foo() -> //c\n(bool);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-inside-return-list",
			code:     "method foo() -> ( //c\n bool //d\n );\n",
			expected: "method foo() -> ( //c\n    bool //d\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-before-body-brace",
			code:     "method foo() //c\n {}\n",
			expected: "method foo() //c\n{\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-inside-body-braces",
			code:     "method  foo( )  {    //c\n }\n",
			expected: "method foo() { //c\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "19-after-body-brace",
			code:     "method foo() {} //c\n\n",
			expected: "method foo() {\n} //c\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "20-wrapped-long-method",
			code:     "dml 1.4;\n\n\n//c\n independent startup //d\n method foo_with_a_very_long_name();\n",
			expected: "dml 1.4;\n\n\n//c\nindependent startup //d\nmethod foo_with_a_very_long_name();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "21-wrapped-long-return",
			code:     "method foo_with_a_very_long_name() //c\n -> (bool) //d\n ;\n",
			expected: "method foo_with_a_very_long_name() //c\n-> (bool) //d\n;\n",
			masked: MaskReturnTypeOverWrapping(
				"method foo_with_a_very_long_name() //c\n-> (\n    bool\n) //d\n;\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "22-wrapped-throws-default",
			code:     "method foo_with_a_very_long_name() //c\n throws //d\n default //e\n {}\n",
			expected: "method foo_with_a_very_long_name() //c\nthrows //d\ndefault //e\n{\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "22a-wrapped-throws-default-new",
			code:     "method foo_with_a_very_long_name() //c\n throws //d\n :default //e\n {}\n",
			expected: "method foo_with_a_very_long_name() //c\nthrows //d\n:default //e\n{\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "23-wrapped-params",
			code:     "method foo( //a\n bool a //b\n, //c\n foo_t b //d\n, //e\n foo_t ** const c //f\n );\n",
			expected: "method foo( //a\n    bool a //b\n    , //c\n    foo_t b //d\n    , //e\n    foo_t **const c //f\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "24-wrapped-return-list",
			code:     "method foo() -> ( //a\n very_long_return_type_1 //b\n, //c\n very_long_return_type_2 //d\n ) throws default {}\n",
			expected: "method foo() -> ( //a\n    very_long_return_type_1 //b\n    , //c\n    very_long_return_type_2 //d\n) throws default {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "24a-wrapped-return-list-default-new",
			code: "method foo() -> ( //a\n very_long_return_type_1 //b\n, //c\n very_long_return_type_2 //d\n ) throws :default {}\n",
			expected: "method foo() -> ( //a\n" +
				"    very_long_return_type_1 //b\n" +
				"    , //c\n" +
				"    very_long_return_type_2 //d\n" +
				") throws :default {\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "25-super-call-default-line-comments",
			code: "method foo() {\n//a\ndefault(\n//b\nx\n//c\n);\n}\n",
			expected: "method foo() {\n" +
				"    //a\n" +
				"    default(\n" +
				"        //b\n" +
				"        x\n" +
				"        //c\n" +
				"    );\n" +
				"}\n",
			masked: MaskUncategorizedFormatting(
				"method foo() {\n    //a\n    default(\n        //b\n        x //c\n    );\n}\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "25a-super-call-default-no-args-line-comments",
			code:     "method foo() {\n//a\ndefault();\n}\n",
			expected: "method foo() {\n    //a\n    default();\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintMethodUnified_WrappingPriorityCascade(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		// Group 1: 3-Stage Cascade (Qualifiers + Arguments + Returns)
		// Code: "shared independent method foo(int a, int b) -> (int, bool);\n"
		// Length of full statement is 61 chars.
		{
			name:     "01-all-fits-single-line",
			code:     "shared independent method foo(int a, int b) -> (int, bool);\n",
			expected: "shared independent method foo(int a, int b) -> (int, bool);\n",
			lineCap:  70,
			pc:       func(f func()) { f() },
		},
		{
			// LineCap 45: Prefix "shared independent method foo(" is 30 chars.
			// Full line is 61 chars (> 45).
			// Arguments wrap first (outer paren). "int a, int b" fits on one line.
			// Return arguments "-> (int, bool);" is 18 chars (<= 45), so stays single-line.
			name: "02-arguments-wrap-first",
			code: "shared independent method foo(int a, int b) -> (int, bool);\n",
			expected: "shared independent method foo(\n" +
				"    int a, int b\n" +
				") -> (int, bool);\n",
			lineCap: 45,
			pc:      func(f func()) { f() },
		},
		{
			// LineCap 36: Prefix "shared independent method foo(" is 30 chars (<= 36).
			// Arguments wrap first, and return suffix ") -> (first_status_t, second_status_t);"
			// is 40 chars (> 36), so return arguments wrap second.
			// Qualifiers stay single-line because prefix 30 <= 36.
			name: "03-arguments-and-returns-wrap",
			code: "shared independent method foo(int a, int b) -> (first_status_t, second_status_t);\n",
			expected: "shared independent method foo(\n" +
				"    int a, int b\n" +
				") -> (\n" +
				"    first_status_t, second_status_t\n" +
				");\n",
			lineCap: 36,
			pc:      func(f func()) { f() },
		},
		{
			// LineCap 25: Prefix "shared independent method foo(" is 30 chars (> 25).
			// Qualifiers wrap third before "method". Inner return elements also explode.
			name: "04-all-three-wrap-qualifiers-third",
			code: "shared independent method foo(int a, int b) -> (first_status_t, second_status_t);\n",
			expected: "shared independent\n" +
				"method foo(\n" +
				"    int a, int b\n" +
				") -> (\n" +
				"    first_status_t,\n" +
				"    second_status_t\n" +
				");\n",
			lineCap: 25,
			pc:      func(f func()) { f() },
		},

		// Group 2: Qualifiers + Arguments (No Returns)
		// Code: "independent startup method calculate_offset(int x, int y);\n"
		// Length of full statement is 59 chars.
		// Prefix "independent startup method calculate_offset(" is 44 chars.
		{
			name: "05-qualifiers-and-args-args-wrap-first",
			code: "independent startup method calculate_offset(int x, int y);\n",
			expected: "independent startup method calculate_offset(\n" +
				"    int x, int y\n" +
				");\n",
			lineCap: 46,
			pc:      func(f func()) { f() },
		},
		{
			name: "06-qualifiers-and-args-qualifiers-wrap-second",
			code: "independent startup method calculate_offset(int x, int y);\n",
			expected: "independent startup\n" +
				"method calculate_offset(\n" +
				"    int x, int y\n" +
				");\n",
			lineCap: 35,
			pc:      func(f func()) { f() },
		},

		// Group 3: Qualifiers + Returns (No Arguments)
		// Code: "independent startup method get_status() -> (bool, uint32);\n"
		// Length of full statement is 59 chars.
		// Prefix "independent startup method get_status() -> (" is 44 chars.
		{
			name: "07-qualifiers-and-returns-returns-wrap-first",
			code: "independent startup method get_status() -> (bool, uint32);\n",
			expected: "independent startup method get_status() -> (\n" +
				"    bool, uint32\n" +
				");\n",
			lineCap: 46,
			pc:      func(f func()) { f() },
		},
		{
			name: "08-qualifiers-and-returns-qualifiers-wrap-second",
			code: "independent startup method get_status() -> (bool, uint32);\n",
			expected: "independent startup\n" +
				"method get_status() -> (\n" +
				"    bool, uint32\n" +
				");\n",
			lineCap: 35,
			pc:      func(f func()) { f() },
		},

		// Group 4: Method Definition with Body, throws, default
		// Code: "shared independent method process_data(int a, int b) -> (status_t, result_t) throws default {}\n"
		{
			name: "09-definition-args-wrap-first",
			code: "shared independent method process_data(int a, int b) -> (status_t, result_t) throws default {}\n",
			expected: "shared independent method process_data(\n" +
				"    int a, int b\n" +
				") -> (status_t, result_t) throws default {\n" +
				"}\n",
			lineCap: 60,
			pc:      func(f func()) { f() },
		},
		{
			name: "10-definition-args-and-returns-wrap",
			code: "shared independent method process_data(int a, int b) -> (status_t, result_t) throws default {}\n",
			expected: "shared independent method process_data(\n" +
				"    int a, int b\n" +
				") -> (\n" +
				"    status_t, result_t\n" +
				") throws default {\n" +
				"}\n",
			lineCap: 40,
			pc:      func(f func()) { f() },
		},
		{
			name: "11-definition-all-wrap",
			code: "shared independent method process_data(int a, int b) -> (status_t, result_t) throws default {}\n",
			expected: "shared independent\n" +
				"method process_data(\n" +
				"    int a, int b\n" +
				") -> (\n" +
				"    status_t, result_t\n" +
				") throws default {\n" +
				"}\n",
			lineCap: 30,
			pc:      func(f func()) { f() },
		},

		// Group 5: Inner Element Explosion (when elements themselves exceed lineCap)
		{
			name: "12-inner-elements-explode-when-exceeding-line",
			code: "shared independent method foo(very_long_argument_type_1 a, very_long_argument_type_2 b) -> (very_long_return_type_1, very_long_return_type_2);\n",
			expected: "shared independent method foo(\n" +
				"    very_long_argument_type_1 a,\n" +
				"    very_long_argument_type_2 b\n" +
				") -> (\n" +
				"    very_long_return_type_1,\n" +
				"    very_long_return_type_2\n" +
				");\n",
			lineCap: 45,
			pc:      func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintMethodUnified_WrappingPriorityCascade_EdgeCases(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		// Case 1: Empty args with multiple returns
		{
			name: "edge-empty-args-multi-returns",
			code: "independent startup method get_status() -> (bool, uint32, int64);\n",
			expected: "independent startup method get_status() -> (\n" +
				"    bool, uint32, int64\n" +
				");\n",
			lineCap: 45,
			pc:      func(f func()) { f() },
		},
		// Case 2: Sequence parameter type
		{
			name: "edge-sequence-param",
			code: "shared independent method transmit(sequence(uint8) payload, int flags) -> (bool);\n",
			expected: "shared independent method transmit(\n" +
				"    sequence(uint8) payload, int flags\n" +
				") -> (bool);\n",
			lineCap: 45,
			pc:      func(f func()) { f() },
		},
		// Case 3: Function pointer parameter type
		{
			name: "edge-func-ptr-param",
			code: "shared independent method register_handler(void (*handler)(int status), int priority) -> (bool);\n",
			expected: "shared independent method register_handler(\n" +
				"    void (*handler)(int status), int priority\n" +
				") -> (bool);\n",
			lineCap: 45,
			pc:      func(f func()) { f() },
		},
		// Case 4: Comments inside argument list
		{
			name: "edge-comments-inside-args",
			code: "shared independent method process(/* in */ int a, /* in */ int b) -> (bool);\n",
			expected: "shared independent method process(\n" +
				"    /* in */ int a, /* in */ int b\n" +
				") -> (bool);\n",
			lineCap: 40,
			pc:      func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintMethodUnified_WrappingPriorityCascade_Idempotence(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "idempotent-01",
			code:     "shared independent method foo(int a, int b) -> (int, bool);\n",
			expected: "shared independent method foo(int a, int b) -> (int, bool);\n",
			lineCap:  70,
			pc:       func(f func()) { f() },
		},
		{
			name: "idempotent-02",
			code: "shared independent method foo(\n" +
				"    int a, int b\n" +
				") -> (int, bool);\n",
			expected: "shared independent method foo(\n" +
				"    int a, int b\n" +
				") -> (int, bool);\n",
			lineCap: 45,
			pc:      func(f func()) { f() },
		},
		{
			name: "idempotent-03",
			code: "shared independent method foo(\n" +
				"    int a, int b\n" +
				") -> (\n" +
				"    first_status_t, second_status_t\n" +
				");\n",
			expected: "shared independent method foo(\n" +
				"    int a, int b\n" +
				") -> (\n" +
				"    first_status_t, second_status_t\n" +
				");\n",
			lineCap: 36,
			pc:      func(f func()) { f() },
		},
		{
			name: "idempotent-04",
			code: "shared independent\n" +
				"method foo(\n" +
				"    int a, int b\n" +
				") -> (\n" +
				"    first_status_t,\n" +
				"    second_status_t\n" +
				");\n",
			expected: "shared independent\n" +
				"method foo(\n" +
				"    int a, int b\n" +
				") -> (\n" +
				"    first_status_t,\n" +
				"    second_status_t\n" +
				");\n",
			lineCap: 25,
			pc:      func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
