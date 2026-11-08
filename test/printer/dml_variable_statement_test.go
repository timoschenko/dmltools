// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"fmt"
	"testing"
)

func TestVariableStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	kinds := []string{"local", "session", "static", "saved"}

	testCases := []printerTestCase{
		{
			name:     "basic-declaration",
			code:     "   %s    int   x   ;",
			expected: "%s int x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "declaration-with-initializer",
			code:     "%s   int   x   =   42;",
			expected: "%s int x = 42;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "line-wrapping",
			code:     "%s   int   very_long_variable_name_exceeding_forty_characters   =   (100);",
			expected: "%s int very_long_variable_name_exceeding_forty_characters = (\n    100\n);",
			pc:       func(f func()) { f() },
		},
	}

	suiteCases := make([]printerTestCase, 0, len(kinds)*len(testCases))
	for _, kind := range kinds {
		for _, tc := range testCases {
			name := fmt.Sprintf("%s-%s", kind, tc.name)
			code := fmt.Sprintf(tc.code, kind)
			expected := fmt.Sprintf(tc.expected, kind)
			masked := MaskTransform(tc.masked, func(s string) string {
				return fmt.Sprintf(s, kind)
			})

			suiteCases = append(suiteCases, printerTestCase{
				task006Skip: tc.task006Skip,
				name:        name,
				code:        code,
				expected:    expected,
				masked:      masked,
				lineCap:     tc.lineCap,
				pc:          tc.pc,
			})
		}
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, suiteCases, lineCap)
}

func TestVariableStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	kinds := []string{"local", "session", "static", "saved"}

	testCases := []printerTestCase{
		{
			name:     "block-comment",
			code:     "%s   /* block */   int   x  =  1;",
			expected: "%s /* block */ int x = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "leading-comment",
			code:     "// comment\n %s   int   x   = 1;",
			expected: "// comment\n%s int x = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "inline-comment",
			code:     "%s // inline\n int   x   = 1;",
			expected: "%s // inline\n    int x = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "inline-comment-after-assign",
			code:     "%s int x = // comment\n 1 + 2;",
			expected: "%s int x = // comment\n    1 + 2;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "trailing-comment",
			code:     "%s int  x  = 1;    // trailing",
			expected: "%s int x = 1; // trailing",
			pc:       func(f func()) { f() },
		},
		{
			name:     "middle-block-comment",
			code:     "%s int x    /* comment */= 1;",
			expected: "%s int x /* comment */ = 1;",
			pc:       func(f func()) { f() },
		},
	}

	suiteCases := make([]printerTestCase, 0, len(kinds)*len(testCases))
	for _, kind := range kinds {
		for _, tc := range testCases {
			name := fmt.Sprintf("%s-%s", kind, tc.name)
			code := fmt.Sprintf(tc.code, kind)
			expected := fmt.Sprintf(tc.expected, kind)
			masked := MaskTransform(tc.masked, func(s string) string {
				return fmt.Sprintf(s, kind)
			})

			suiteCases = append(suiteCases, printerTestCase{
				task006Skip: tc.task006Skip,
				name:        name,
				code:        code,
				expected:    expected,
				masked:      masked,
				lineCap:     tc.lineCap,
				pc:          tc.pc,
			})
		}
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, suiteCases, lineCap)
}

func TestVariableStatement_CdeclList(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	kinds := []string{"local", "session", "static", "saved"}

	testCases := []printerTestCase{
		{
			name:     "01-tuple-single",
			code:     "%s   (  int   x  );",
			expected: "%s int x;",
			masked: MaskSingleElementTupleParenthesesRemoval(
				"%s (int x);",
			), // OK, to pospone until
			pc: func(f func()) { f() },
		},
		{
			name:     "02-tuple-multiple",
			code:     "%s (int   x,   double y);",
			expected: "%s (int x, double y);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-tuple-with-initializer",
			code:     "%s (   int   x  ,double   y) = get_pair();",
			expected: "%s (int x, double y) = get_pair();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-tuple-line-wrapping",
			code:     "%s (int  very_long_variable_name_a  , int very_long_variable_name_b);",
			expected: "%s (\n    int very_long_variable_name_a,\n    int very_long_variable_name_b\n);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-tuple-with-initializer-line-wrapping",
			code:     "%s (int long_variable_a, int long_variable_b) = very_long_initializer_expression_call();",
			expected: "%s (\n    int long_variable_a,\n    int long_variable_b\n) = very_long_initializer_expression_call();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-tuple-mid-line-wrapping",
			code:     "%s (    int   variable_a ,  bool  variable_b) = foo();",
			expected: "%s (\n    int variable_a, bool variable_b\n) = foo();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-tuple-discard",
			code:     "%s (   bool   a ,  int _) = m();",
			expected: "%s (bool a, int _) = m();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-tuple-multi-discard",
			code:     "%s (bool   a,   int  _  ,  float   _) = returns_three();",
			expected: "%s (\n    bool a, int _, float _\n) = returns_three();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-tuple-mixed-types",
			code:     "%s (const    int * const   ptr,    char buf[  10]);",
			expected: "%s (const int *const ptr, char buf[10]);",
			lineCap:  80,
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-tuple-literal-initializer",
			code:     "%s (  int  x ,  bool  b) = (1, true);",
			expected: "%s (int x, bool b) = (1, true);",
			pc:       func(f func()) { f() },
		},
	}

	suiteCases := make([]printerTestCase, 0, len(kinds)*len(testCases))
	for _, kind := range kinds {
		for _, tc := range testCases {
			name := fmt.Sprintf("%s-%s", kind, tc.name)
			code := fmt.Sprintf(tc.code, kind)
			expected := fmt.Sprintf(tc.expected, kind)
			masked := MaskTransform(tc.masked, func(s string) string {
				return fmt.Sprintf(s, kind)
			})

			suiteCases = append(suiteCases, printerTestCase{
				task006Skip: tc.task006Skip,
				name:        name,
				code:        code,
				expected:    expected,
				masked:      masked,
				lineCap:     tc.lineCap,
				pc:          tc.pc,
			})
		}
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, suiteCases, lineCap)
}

func TestVariableStatement_CdeclList_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	kinds := []string{"local", "session", "static", "saved"}

	testCases := []printerTestCase{
		{
			name:     "01-comments-in-parens",
			code:     "%s (int   x,   // comment\n   int y);",
			expected: "%s (\n    int x, // comment\n    int y\n);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comments-around-lparen",
			code:     "%s   /* comment */   (int   x,  int y);",
			expected: "%s /* comment */ (int x, int y);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comments-before-rparen",
			code:     "%s (  int  x , int y /* comment */);",
			expected: "%s (int x, int y /* comment */);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comments-after-rparen",
			code:     "%s (int    x,   int   y) /* note */ = 1;",
			expected: "%s (int x, int y) /* note */ = 1;",
			pc:       func(f func()) { f() },
		},
	}

	suiteCases := make([]printerTestCase, 0, len(kinds)*len(testCases))
	for _, kind := range kinds {
		for _, tc := range testCases {
			name := fmt.Sprintf("%s-%s", kind, tc.name)
			code := fmt.Sprintf(tc.code, kind)
			expected := fmt.Sprintf(tc.expected, kind)
			masked := MaskTransform(tc.masked, func(s string) string {
				return fmt.Sprintf(s, kind)
			})

			suiteCases = append(suiteCases, printerTestCase{
				task006Skip: tc.task006Skip,
				name:        name,
				code:        code,
				expected:    expected,
				masked:      masked,
				lineCap:     tc.lineCap,
				pc:          tc.pc,
			})
		}
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, suiteCases, lineCap)
}

func TestVariableStatement_InlineTypes(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	kinds := []string{"local", "session", "static", "saved"}

	testCases := []printerTestCase{
		{
			name:     "bitfields-no-init",
			code:     "%s bitfields 32 { uint8 x @ [7:0]; } data;",
			expected: "%s bitfields 32 {\n    uint8 x @ [7:0];\n} data;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "bitfields-with-init",
			code:     "%s bitfields 32 { uint8 x @ [7:0]; } data = 1;",
			expected: "%s bitfields 32 {\n    uint8 x @ [7:0];\n} data = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name: "bitfields-multi-field",
			code: "%s bitfields 32 { int a @ [0:4]; int b @ [5]; } bf = 0;",
			expected: "%s bitfields 32 {\n" +
				"    int a @ [0:4];\n" +
				"    int b @ [5];\n" +
				"} bf = 0;",
			pc: func(f func()) { f() },
		},
		{
			name:     "bitfields-messy-whitespace",
			code:     "%s   bitfields   32  {   uint8   x  @  [  7  :  0  ]  ;  }   data  =  1  ;",
			expected: "%s bitfields 32 {\n    uint8 x @ [7:0];\n} data = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name: "bitfields-comments",
			code: "%s bitfields 32 {\n    // comment\n    uint8 x @ [7:0]; // trailing\n} data = 1;",
			expected: "%s bitfields 32 {\n" +
				"    // comment\n" +
				"    uint8 x @ [7:0]; // trailing\n" +
				"} data = 1;",
			pc: func(f func()) { f() },
		},
		{
			name: "struct-no-init",
			code: "%s struct { int a; float b; } s;",
			expected: "%s struct {\n" +
				"    int a;\n" +
				"    float b;\n" +
				"} s;",
			pc: func(f func()) { f() },
		},
		{
			name: "struct-with-cdecl-init",
			code: "%s struct { int a; float b; } s = {1, 2.0};",
			expected: "%s struct {\n" +
				"    int a;\n" +
				"    float b;\n" +
				"} s = {1, 2.0};",
			pc: func(f func()) { f() },
		},
		{
			name:     "struct-with-designated-init",
			code:     "%s struct { int a; } s = {.a = 1};",
			expected: "%s struct {\n    int a;\n} s = {.a = 1};",
			pc:       func(f func()) { f() },
		},
		{
			name: "struct-with-partial-designated-init",
			code: "%s struct { int a; int b; } s = {.a = 1, ...};",
			expected: "%s struct {\n" +
				"    int a;\n" +
				"    int b;\n" +
				"} s = {.a = 1, ...};",
			pc: func(f func()) { f() },
		},
		{
			name:     "layout-no-init",
			code:     "%s layout `big-endian` { uint16 a; } l;",
			expected: "%s layout `big-endian` {\n    uint16 a;\n} l;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "layout-with-init",
			code:     "%s layout `big-endian` { uint16 a; } l = {0};",
			expected: "%s layout `big-endian` {\n    uint16 a;\n} l = {0};",
			pc:       func(f func()) { f() },
		},
	}

	suiteCases := make([]printerTestCase, 0, len(kinds)*len(testCases))
	for _, kind := range kinds {
		for _, tc := range testCases {
			name := fmt.Sprintf("%s-%s", kind, tc.name)
			code := fmt.Sprintf(tc.code, kind)
			expected := fmt.Sprintf(tc.expected, kind)
			masked := MaskTransform(tc.masked, func(s string) string {
				return fmt.Sprintf(s, kind)
			})

			suiteCases = append(suiteCases, printerTestCase{
				task006Skip: tc.task006Skip,
				name:        name,
				code:        code,
				expected:    expected,
				masked:      masked,
				lineCap:     tc.lineCap,
				pc:          tc.pc,
			})
		}
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, suiteCases, lineCap)
}

func TestVariableStatement_ComplexTypes(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	kinds := []string{"local", "session", "static", "saved"}

	testCases := []printerTestCase{
		{
			name:     "pointer",
			code:     "%s  int  * ptr;",
			expected: "%s int *ptr;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "pointer-with-initializer",
			code:     "%s int *ptr = NULL;",
			expected: "%s int *ptr = NULL;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "pointer-const",
			code:     "%s const   int * const ptr = NULL;",
			expected: "%s const int *const ptr = NULL;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "pointer-double",
			code:     "%s  int    *   * double_ptr;",
			expected: "%s int **double_ptr;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "array",
			code:     "%s int   arr[ 10];",
			expected: "%s int arr[10];",
			pc:       func(f func()) { f() },
		},
		{
			name:     "array-literal-initializer",
			code:     "%s    int arr[3] = [1, 2, 3];",
			expected: "%s int arr[3] = [1, 2, 3];",
			pc:       func(f func()) { f() },
		},
		{
			name:     "array-cdecl-initializer",
			code:     "%s int arr  [3] = {1, 2, 3};",
			expected: "%s int arr[3] = {1, 2, 3};",
			pc:       func(f func()) { f() },
		},
		{
			name:     "array-2d",
			code:     "%s  int  grid[ 2]   [ 3];",
			expected: "%s int grid[2][3];",
			pc:       func(f func()) { f() },
		},
		{
			name:     "func-ptr",
			code:     "%s   void  ( * callback )(  int, char);",
			expected: "%s void (*callback)(int, char);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "func-ptr-with-initializer",
			code:     "%s   int  ( * cb ) ( void) = my_cb;",
			expected: "%s int (*cb)(void) = my_cb;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "typeof",
			code:     "%s typeof(  x)   y;",
			expected: "%s typeof(x) y;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "typeof-with-initializer",
			code:     "%s   typeof(  x)   y = x;",
			expected: "%s typeof(x) y = x;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "sequence",
			code:     "%s sequence(  int)  seq;",
			expected: "%s sequence(int) seq;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "sequence-with-initializer",
			code:     "%s   sequence(  int)   seq   = [1, 2, 3];",
			expected: "%s sequence(int) seq = [1, 2, 3];",
			pc:       func(f func()) { f() },
		},
		{
			name:     "hook",
			code:     "%s hook(   int  ,  char) h;",
			expected: "%s hook(int, char) h;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "discard-ident",
			code:     "%s   int   _ = 5;",
			expected: "%s int _ = 5;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "discard-ident-no-init",
			code:     "%s   int   _;",
			expected: "%s int _;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "const-qualifier",
			code:     "%s   const   int   c = 10;",
			expected: "%s const int c = 10;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "unsigned-qualifier",
			code:     "%s   unsigned   int  u = 42;",
			expected: "%s unsigned int u = 42;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "signed-qualifier",
			code:     "%s  signed   int  s  = -1;",
			expected: "%s signed int s = -1;",
			pc:       func(f func()) { f() },
		},
	}

	suiteCases := make([]printerTestCase, 0, len(kinds)*len(testCases))
	for _, kind := range kinds {
		for _, tc := range testCases {
			name := fmt.Sprintf("%s-%s", kind, tc.name)
			code := fmt.Sprintf(tc.code, kind)
			expected := fmt.Sprintf(tc.expected, kind)
			masked := MaskTransform(tc.masked, func(s string) string {
				return fmt.Sprintf(s, kind)
			})

			suiteCases = append(suiteCases, printerTestCase{
				task006Skip: tc.task006Skip,
				name:        name,
				code:        code,
				expected:    expected,
				masked:      masked,
				lineCap:     tc.lineCap,
				pc:          tc.pc,
			})
		}
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, suiteCases, lineCap)
}

func TestVariableStatement_Wrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 25

	kinds := []string{"local", "session", "static", "saved"}

	testCases := []printerTestCase{
		{
			name: "binary-expr-wrapping",
			code: "%s int x = 1 + 2 + 3 + 4 + 5;",
			expected: "%s int x = 1\n" +
				"    + 2\n" +
				"    + 3\n" +
				"    + 4\n" +
				"    + 5;",
			pc: func(f func()) { f() },
		},
		{
			name: "parenthesized-binary-expr-wrapping",
			code: "%s int x = (1 + 2 + 3 + 4 + 5);",
			expected: "%s int x = (\n" +
				"    1 + 2 + 3 + 4 + 5\n" +
				");",
			pc: func(f func()) { f() },
		},
	}

	suiteCases := make([]printerTestCase, 0, len(kinds)*len(testCases))
	for _, kind := range kinds {
		for _, tc := range testCases {
			name := fmt.Sprintf("%s-%s", kind, tc.name)
			code := fmt.Sprintf(tc.code, kind)
			expected := fmt.Sprintf(tc.expected, kind)

			suiteCases = append(suiteCases, printerTestCase{
				task006Skip: tc.task006Skip,
				name:        name,
				code:        code,
				expected:    expected,
				lineCap:     tc.lineCap,
				pc:          tc.pc,
			})
		}
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, suiteCases, lineCap)
}
