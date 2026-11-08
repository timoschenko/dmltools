// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestAssignStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 80
	testCases := []printerTestCase{
		{
			name:     "01-simple-assign",
			code:     "x   =   1   ;",
			expected: "x = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-chained-assign",
			code:     "a   =   b   =   c   =   10   ;",
			expected: "a = b = c = 10;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-tuple-literal-basic",
			code:     "(  a  ,   b  )   =   (  1  ,   2  )  ;",
			expected: "(a, b) = (1, 2);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-tuple-literal-three-elements",
			code:     "(a, b, c) = (1, 2, 3);",
			expected: "(a, b, c) = (1, 2, 3);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-tuple-literal-discard",
			code:     "(_, status) = get_status();",
			expected: "(_, status) = get_status();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-tuple-literal-method-call",
			code:     "(memop, msb1, lsb) = dev._make_read_memop(enabled_bits, &value);",
			expected: "(memop, msb1, lsb) = dev._make_read_memop(enabled_bits, &value);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-tuple-literal-complex-expressions",
			code:     "(a[0], b.field) = (x + 1, y * 2);",
			expected: "(a[0], b.field) = (x + 1, y * 2);",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAssignStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 80
	testCases := []printerTestCase{
		{
			name:     "01-block-before",
			code:     "/* comment */ (a, b) = (1, 2);",
			expected: "/* comment */ (a, b) = (1, 2);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-block-between-tuple-elements",
			code:     "(a, /* comment */ b) = (1, 2);",
			expected: "(a, /* comment */ b) = (1, 2);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-block-before-assign",
			code:     "(a, b) /* comment */ = (1, 2);",
			expected: "(a, b) /* comment */ = (1, 2);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-trailing-line-comment",
			code:     "(a, b) = (1, 2); // trailing",
			expected: "(a, b) = (1, 2); // trailing",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-line-comment-after-assign",
			code:     "x = // comment\n 1 + 2;",
			expected: "x = // comment\n    1 + 2;",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAssignStatement_Wrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 20
	testCases := []printerTestCase{
		{
			name: "01-binary-expr-wrapping",
			code: "x = 1 + 2 + 3 + 4 + 5;",
			expected: "x = 1\n" +
				"    + 2\n" +
				"    + 3\n" +
				"    + 4\n" +
				"    + 5;",
			pc: func(f func()) { f() },
		},
		{
			name: "02-chained-assign-wrapping",
			code: "a = b = 1 + 2 + 3 + 4;",
			expected: "a = b = 1\n" +
				"    + 2\n" +
				"    + 3\n" +
				"    + 4;",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAssignop(t *testing.T) {
	t.Parallel()

	const lineCap = 80
	testCases := []printerTestCase{
		{
			name:     "01-add-assign",
			code:     "x   +=   1   ;",
			expected: "x += 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-sub-assign",
			code:     "x   -=   2   ;",
			expected: "x -= 2;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-mul-assign",
			code:     "x   *=   3   ;",
			expected: "x *= 3;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-div-assign",
			code:     "x   /=   4   ;",
			expected: "x /= 4;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-mod-assign",
			code:     "x   %=   5   ;",
			expected: "x %= 5;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-and-assign",
			code:     "x   &=   6   ;",
			expected: "x &= 6;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-or-assign",
			code:     "x   |=   7   ;",
			expected: "x |= 7;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-xor-assign",
			code:     "x   ^=   8   ;",
			expected: "x ^= 8;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-lshift-assign",
			code:     "x   <<=   9   ;",
			expected: "x <<= 9;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-rshift-assign",
			code:     "x   >>=   10   ;",
			expected: "x >>= 10;",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAssignop_Wrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 20
	testCases := []printerTestCase{
		{
			name: "01-add-assign-wrapping",
			code: "x += 1 + 2 + 3 + 4 + 5;",
			expected: "x += 1\n" +
				"    + 2\n" +
				"    + 3\n" +
				"    + 4\n" +
				"    + 5;",
			pc: func(f func()) { f() },
		},
		{
			name: "02-sub-assign-wrapping",
			code: "x -= 1 + 2 + 3 + 4 + 5;",
			expected: "x -= 1\n" +
				"    + 2\n" +
				"    + 3\n" +
				"    + 4\n" +
				"    + 5;",
			pc: func(f func()) { f() },
		},
		{
			name: "03-mul-assign-wrapping",
			code: "x *= 1 + 2 + 3 + 4 + 5;",
			expected: "x *= 1\n" +
				"    + 2\n" +
				"    + 3\n" +
				"    + 4\n" +
				"    + 5;",
			pc: func(f func()) { f() },
		},
		{
			name: "04-or-assign-wrapping",
			code: "x |= FLAG_A | FLAG_B | FLAG_C;",
			expected: "x |= FLAG_A\n" +
				"    | FLAG_B\n" +
				"    | FLAG_C;",
			pc: func(f func()) { f() },
		},
		{
			name: "05-and-assign-wrapping",
			code: "x &= MASK_A & MASK_B & MASK_C;",
			expected: "x &= MASK_A\n" +
				"    & MASK_B\n" +
				"    & MASK_C;",
			pc: func(f func()) { f() },
		},
		{
			name: "06-lshift-assign-wrapping",
			code: "x <<= shift_a + shift_b + shift_c;",
			expected: "x <<= shift_a\n" +
				"    + shift_b\n" +
				"    + shift_c;",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAssignop_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-line-comment-after-add-assign",
			code:     "x += // comment\n 1 + 2;",
			expected: "x += // comment\n    1 + 2;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-line-comment-after-sub-assign",
			code:     "x -= // comment\n 1 + 2;",
			expected: "x -= // comment\n    1 + 2;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-inline-block-around-op",
			code:     "x /* c1 */ += /* c2 */ 1;",
			expected: "x /* c1 */ += /* c2 */ 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-trailing-line-comment",
			code:     "x += 1; // trailing",
			expected: "x += 1; // trailing",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
