// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestTokenAdjacencyTraps(t *testing.T) {
	t.Parallel()

	const lineCap = 60
	testCases := []printerTestCase{
		// =========================================================================
		// 1. MUL and DIV Interactions with Comments & Dereference
		// =========================================================================
		{
			// Trap: '/' followed by '*' forms '/*' (block comment start).
			// Spacing must be preserved between '/' and '*' in 'b / *buf'.
			name: "01-mul-after-div-deref",
			code: "method m() {\n" +
				"    local int a = b / *buf;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = b / *buf;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: If 'b / *buf' collapsed to 'b/*buf', the '/*' would swallow
			// everything until the later block comment end '*/'.
			name: "02-mul-after-div-with-trailing-comment",
			code: "method m() {\n" +
				"    local int a = b / *buf;\n" +
				"    /* trailing comment */\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = b / *buf; /* trailing comment */\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: '*' followed by '/' forms '*/' (COMMENT_INLINE_END).
			// Grouped expressions ensure separation between '*' and '/'.
			name: "03-mul-before-div-grouped-assignment",
			code: "method m() {\n" +
				"    a = (x * y) / z;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    a = (x * y) / z;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Dereferenced pointer operand before division.
			name: "04-mul-before-div-pointer-deref",
			code: "method m() {\n" +
				"    local int a = (*p) / d;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = (*p) / d;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: '*' followed immediately by '//' line comment.
			// Spacing must ensure '* //' rather than '*//'.
			name: "05-mul-before-line-comment-pointer",
			code: "param p: const int * // line comment\n" +
				"    = 1;\n",
			expected: "param p: const int * // line comment\n" +
				"    = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Binary multiplication operator followed by '//' line comment.
			name: "06-mul-before-line-comment-expr",
			code: "method m() {\n" +
				"    local int a = b * // line comment\n" +
				"        c;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = b * // line comment\n" +
				"        c;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Issue #20 - pointer asterisk followed by block comment '/* ... */'.
			// Dropping space produces '*/*', which contains '*/' (comment end).
			name:     "07-block-comment-after-mul-pointer",
			code:     "param x: const short_t * /* c3 */ = 1;\n",
			expected: "param x: const short_t * /* c3 */ = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			// Trap: Binary multiplication operator followed by block comment '/* ... */'.
			// In cvarDecl, SepNone after '*' causes space drop resulting in '*/*'.
			name: "08-block-comment-after-mul-expr",
			code: "method m() {\n" +
				"    local int a = b * /* comment */ c;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = b * /* comment */ c;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Block comment '/* ... */' followed by pointer asterisk '*'.
			name:     "09-block-comment-before-mul-pointer",
			code:     "param x: const short_t /* c3 */ * = 1;\n",
			expected: "param x: const short_t /* c3 */ * = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			// Trap: Block comment '/* ... */' followed by binary multiplication '*'.
			name: "10-block-comment-before-mul-expr",
			code: "method m() {\n" +
				"    local int a = b /* comment */ * c;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = b /* comment */ * c;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Block comment '/* ... */' followed by unary dereference '*'.
			name: "11-block-comment-before-unary-deref",
			code: "method m() {\n" +
				"    local int a = /* comment */ *buf;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = /* comment */ *buf;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Division operator '/' followed by block comment '/* ... */'.
			name: "12-div-before-block-comment",
			code: "method m() {\n" +
				"    local int a = b / /* comment */ c;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = b / /* comment */ c;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},

		// =========================================================================
		// 2. Compound Colon Tokens
		// =========================================================================
		{
			// Trap: ':' followed by '{' forms ':{' (LBRACE_NEW).
			// LBRACE_NEW is only allowed in method definitions. In caseStatement,
			// the colon must wrap with newline so it doesn't emit 'case 1:{'.
			name: "13-case-with-compound-block",
			code: "method m() {\n" +
				"    switch (x) {\n" +
				"    case 1: {\n" +
				"        break;\n" +
				"    }\n" +
				"    }\n" +
				"}\n",
			expected: "method m() {\n" +
				"    switch (x) {\n" +
				"    case 1:\n" +
				"        {\n" +
				"            break;\n" +
				"        }\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: ':' followed by 'default' forms ':default' (DEFAULT_NEW).
			// DEFAULT_NEW is only allowed in method qualifiers. In switch cases,
			// 'case 1:' must not collapse with 'default:' into 'case 1:default:'.
			name: "14-case-immediately-before-default",
			code: "method m() {\n" +
				"    switch (x) {\n" +
				"    case 1:\n" +
				"    default:\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			expected: "method m() {\n" +
				"    switch (x) {\n" +
				"    case 1:\n" +
				"    default:\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},

		// =========================================================================
		// 3. Unary and Binary Operator Collisions
		// =========================================================================
		{
			// Trap: Binary '+' followed by unary '+' must not collapse into '++' (INC).
			name: "15-binary-plus-with-unary-plus",
			code: "method m() {\n" +
				"    local int a = b + +c;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = b + +c;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Binary '-' followed by unary '-' must not collapse into '--' (DEC).
			name: "16-binary-minus-with-unary-minus",
			code: "method m() {\n" +
				"    local int a = b - -c;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = b - -c;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Nested unary '+' (+ +b) must not collapse into '++b' (pre-increment).
			name: "17-nested-unary-plus",
			code: "method m() {\n" +
				"    local int a = + +b;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = + +b;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Nested unary '-' (- -b) must not collapse into '--b' (pre-decrement).
			name: "18-nested-unary-minus",
			code: "method m() {\n" +
				"    local int a = - -b;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = - -b;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Binary bitwise '&' followed by address-of '&' must not collapse into '&&' (AND).
			// In cvarDecl, SepNone on Bitand causes collapse to 'b&&c'.
			name: "19-binary-and-with-address-of",
			code: "method m() {\n" +
				"    local int a = b & &c;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local int a = b & &c;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: Nested address-of (& &b) must not collapse into '&&' (logical AND).
			name: "20-nested-address-of",
			code: "method m() {\n" +
				"    local void *p = & &b;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    local void *p = & &b;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},

		// =========================================================================
		// 4. Operator and Assignment Collisions
		// =========================================================================
		{
			// Trap: Binary operators in statements must retain spacing around operands
			// to avoid merging into compound assignments (+=, -=, *=, /=, %=).
			name: "21-binary-arithmetic-assignment-spacing",
			code: "method m() {\n" +
				"    a = b + c;\n" +
				"    d = b - c;\n" +
				"    e = b * c;\n" +
				"    f = b / c;\n" +
				"    g = b % c;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    a = b + c;\n" +
				"    d = b - c;\n" +
				"    e = b * c;\n" +
				"    f = b / c;\n" +
				"    g = b % c;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},

		// =========================================================================
		// 5. Keyword & Identifier Separations
		// =========================================================================
		{
			// Trap: Statement keywords followed by identifiers must not collapse
			// into single identifiers (e.g. 'assert x' -> 'assertx').
			name: "22-keyword-ident-spacing",
			code: "method m() {\n" +
				"    assert x;\n" +
				"    delete y;\n" +
				"    return z;\n" +
				"}\n",
			expected: "method m() {\n" +
				"    assert x;\n" +
				"    delete y;\n" +
				"    return z;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			// Trap: 'case' followed by identifier must not collapse into 'casefoo:'.
			name: "23-case-label-ident-spacing",
			code: "method m() {\n" +
				"    switch (x) {\n" +
				"    case foo:\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			expected: "method m() {\n" +
				"    switch (x) {\n" +
				"    case foo:\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
