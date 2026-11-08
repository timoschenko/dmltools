// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestOpenDelimiterComments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		// ==========================================
		// Category 1: Opening Braces `{`
		// ==========================================
		{
			name: "01-brace-if-body",
			code: "method foo() {\n" +
				"if (qwe) {\n" +
				"// comment\n" +
				"qwe();\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    if (qwe) {\n" +
				"        // comment\n" +
				"        qwe();\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "02-brace-else-body",
			code: "method foo() {\n" +
				"if (qwe) {\n" +
				"qwe();\n" +
				"} else {\n" +
				"// comment\n" +
				"other();\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    if (qwe) {\n" +
				"        qwe();\n" +
				"    } else {\n" +
				"        // comment\n" +
				"        other();\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "03-brace-while-body",
			code: "method foo() {\n" +
				"while (qwe) {\n" +
				"// comment\n" +
				"qwe();\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    while (qwe) {\n" +
				"        // comment\n" +
				"        qwe();\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "04-brace-for-body",
			code: "method foo() {\n" +
				"local int i;\n" +
				"for (i = 0; i < 10; ++i) {\n" +
				"// comment\n" +
				"qwe();\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    local int i;\n" +
				"    for (i = 0; i < 10; ++i) {\n" +
				"        // comment\n" +
				"        qwe();\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "05-brace-foreach-body",
			code: "method foo() {\n" +
				"foreach item in (items) {\n" +
				"// comment\n" +
				"qwe();\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    foreach item in (items) {\n" +
				"        // comment\n" +
				"        qwe();\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "06-brace-switch-body",
			code: "method foo() {\n" +
				"switch (x) {\n" +
				"// comment\n" +
				"case 1:\n" +
				"break;\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    switch (x) {\n" +
				"    // comment\n" +
				"    case 1:\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "07-brace-compound-block",
			code: "method foo() {\n" +
				"{\n" +
				"// comment\n" +
				"qwe();\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    {\n" +
				"        // comment\n" +
				"        qwe();\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "08-brace-bank-decl",
			code: "bank b {\n" +
				"// comment\n" +
				"register r;\n" +
				"}\n",
			expected: "bank b {\n" +
				"    // comment\n" +
				"    register r;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "09-brace-register-decl",
			code: "bank b {\n" +
				"register r {\n" +
				"// comment\n" +
				"field f;\n" +
				"}\n" +
				"}\n",
			expected: "bank b {\n" +
				"    register r {\n" +
				"        // comment\n" +
				"        field f;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "10-brace-field-decl",
			code: "bank b {\n" +
				"register r {\n" +
				"field f {\n" +
				"// comment\n" +
				"param val = 1;\n" +
				"}\n" +
				"}\n" +
				"}\n",
			expected: "bank b {\n" +
				"    register r {\n" +
				"        field f {\n" +
				"            // comment\n" +
				"            param val = 1;\n" +
				"        }\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "11-brace-template-decl",
			code: "template t {\n" +
				"// comment\n" +
				"param val = 1;\n" +
				"}\n",
			expected: "template t {\n" +
				"    // comment\n" +
				"    param val = 1;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "12-brace-method-next-line",
			code: "method foo() {\n" +
				"// comment\n" +
				"qwe();\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    // comment\n" +
				"    qwe();\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "13-brace-method-trailing",
			code: "method foo() { // trailing\n" +
				"qwe();\n" +
				"}\n",
			expected: "method foo() { // trailing\n" +
				"    qwe();\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "14-brace-if-empty-body",
			code: "method foo() {\n" +
				"if (qwe) {\n" +
				"// comment\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    if (qwe) {\n" +
				"        // comment\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "15-brace-struct-initializer",
			code: "method foo() {\n" +
				"local s_t s = {\n" +
				"// comment\n" +
				".a = 1\n" +
				"};\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    local s_t s = {\n" +
				"        // comment\n" +
				"        .a = 1\n" +
				"    };\n" +
				"}\n",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 2: Opening Brackets `[`
		// ==========================================
		{
			name: "16-bracket-constant-array",
			code: "constant T = [\n" +
				"// comment\n" +
				"aaaaaaa + aaaaaaaa,\n" +
				"bbbbbbbbbb - bbbbbbbbbb\n" +
				"];\n",
			expected: "constant T = [\n" +
				"    // comment\n" +
				"    aaaaaaa + aaaaaaaa,\n" +
				"    bbbbbbbbbb - bbbbbbbbbb\n" +
				"];\n",
			pc: func(f func()) { f() },
		},
		{
			name: "17-bracket-nested-array",
			code: "constant matrix = [\n" +
				"[\n" +
				"// comment\n" +
				"1,\n" +
				"2\n" +
				"]\n" +
				"];\n",
			expected: "constant matrix = [\n" +
				"    [\n" +
				"        // comment\n" +
				"        1,\n" +
				"        2\n" +
				"    ]\n" +
				"];\n",
			pc: func(f func()) { f() },
		},
		{
			name: "18-bracket-call-arg-array",
			code: "method foo() {\n" +
				"call([\n" +
				"// comment\n" +
				"1,\n" +
				"2\n" +
				"]);\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    call(\n" +
				"        [\n" +
				"            // comment\n" +
				"            1,\n" +
				"            2\n" +
				"        ]\n" +
				"    );\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "19-bracket-subscript-expr",
			code: "method foo() {\n" +
				"local int x = arr[\n" +
				"// comment\n" +
				"index\n" +
				"];\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    local int x = arr[\n" +
				"        // comment\n" +
				"        index\n" +
				"    ];\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "20-bracket-bit-slice-expr",
			code: "method foo() {\n" +
				"local int x = val[\n" +
				"// comment\n" +
				"7:0\n" +
				"];\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    local int x = val[\n" +
				"        // comment\n" +
				"        7:0\n" +
				"    ];\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "21-bracket-trailing-comment",
			code: "constant T = [ // trailing\n" +
				"1,\n" +
				"2\n" +
				"];\n",
			expected: "constant T = [ // trailing\n" +
				"    1, 2\n" +
				"];\n",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 3: Opening Parentheses `(`
		// ==========================================
		{
			name: "22-paren-if-condition",
			code: "method foo() {\n" +
				"if (\n" +
				"// comment\n" +
				"qwe\n" +
				") {\n" +
				"qwe();\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    if (\n" +
				"        // comment\n" +
				"        qwe\n" +
				"    ) {\n" +
				"        qwe();\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "23-paren-while-condition",
			code: "method foo() {\n" +
				"while (\n" +
				"// comment\n" +
				"qwe\n" +
				") {\n" +
				"qwe();\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    while (\n" +
				"        // comment\n" +
				"        qwe\n" +
				"    ) {\n" +
				"        qwe();\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "24-paren-switch-condition",
			code: "method foo() {\n" +
				"switch (\n" +
				"// comment\n" +
				"x\n" +
				") {\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    switch (\n" +
				"        // comment\n" +
				"        x\n" +
				"    ) {\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "25-paren-for-header",
			code: "method foo() {\n" +
				"local int i;\n" +
				"for (\n" +
				"// comment\n" +
				"i = 0; i < 10; ++i\n" +
				") {\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    local int i;\n" +
				"    for (\n" +
				"        // comment\n" +
				"        i = 0;\n" +
				"        i < 10;\n" +
				"        ++i\n" +
				"    ) {\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "26-paren-func-call",
			code: "method foo() {\n" +
				"bar(\n" +
				"// comment\n" +
				"arg1,\n" +
				"arg2\n" +
				");\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    bar(\n" +
				"        // comment\n" +
				"        arg1,\n" +
				"        arg2\n" +
				"    );\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "27-paren-method-params",
			code: "method foo(\n" +
				"// comment\n" +
				"int a\n" +
				") {\n" +
				"}\n",
			expected: "method foo(\n" +
				"    // comment\n" +
				"    int a\n" +
				") {\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "28-paren-method-return",
			code: "method foo() -> (\n" +
				"// comment\n" +
				"bool\n" +
				") {\n" +
				"}\n",
			expected: "method foo() -> (\n" +
				"    // comment\n" +
				"    bool\n" +
				") {\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "29-paren-grouped-expr",
			code: "method foo() {\n" +
				"local int x = (\n" +
				"// comment\n" +
				"a + b\n" +
				") * c;\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    local int x = (\n" +
				"        // comment\n" +
				"        a + b\n" +
				"    )\n" +
				"        * c;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "30-paren-cast-expr",
			code: "method foo() {\n" +
				"local int x = cast(\n" +
				"// comment\n" +
				"val,\n" +
				"int\n" +
				");\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    local int x = cast(\n" +
				"        // comment\n" +
				"        val,\n" +
				"        int\n" +
				"    );\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "31-paren-foreach-container",
			code: "method foo() {\n" +
				"foreach item in (\n" +
				"// comment\n" +
				"items\n" +
				") {\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    foreach item in (\n" +
				"        // comment\n" +
				"        items\n" +
				"    ) {\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "32-paren-trailing-comment",
			code: "method foo() {\n" +
				"bar( // trailing\n" +
				"arg1\n" +
				");\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    bar( // trailing\n" +
				"        arg1\n" +
				"    );\n" +
				"}\n",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 4: Switch Case & Default Labels
		// ==========================================
		{
			name: "33-switch-case-next-line",
			code: "method foo() {\n" +
				"switch (x) {\n" +
				"case 1:\n" +
				"// comment\n" +
				"break;\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    switch (x) {\n" +
				"    case 1:\n" +
				"        // comment\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "34-switch-default-next-line",
			code: "method foo() {\n" +
				"switch (x) {\n" +
				"default:\n" +
				"// comment\n" +
				"break;\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    switch (x) {\n" +
				"    default:\n" +
				"        // comment\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "35-switch-case-fallthrough-next-line",
			code: "method foo() {\n" +
				"switch (x) {\n" +
				"case 1:\n" +
				"// comment\n" +
				"case 2:\n" +
				"break;\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    switch (x) {\n" +
				"    case 1:\n" +
				"    // comment\n" +
				"    case 2:\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "36-switch-case-multiple-stmts",
			code: "method foo() {\n" +
				"switch (x) {\n" +
				"case 1:\n" +
				"// comment\n" +
				"x = 1;\n" +
				"break;\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    switch (x) {\n" +
				"    case 1:\n" +
				"        // comment\n" +
				"        x = 1;\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "37-switch-case-trailing",
			code: "method foo() {\n" +
				"switch (x) {\n" +
				"case 1: // trailing\n" +
				"break;\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    switch (x) {\n" +
				"    case 1: // trailing\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "38-switch-default-trailing",
			code: "method foo() {\n" +
				"switch (x) {\n" +
				"default: // trailing\n" +
				"break;\n" +
				"}\n" +
				"}\n",
			expected: "method foo() {\n" +
				"    switch (x) {\n" +
				"    default: // trailing\n" +
				"        break;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 109)
}
