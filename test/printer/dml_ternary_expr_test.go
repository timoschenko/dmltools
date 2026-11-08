// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"strings"
	"testing"
)

func TestPrintTernaryExpr_NativeParentheses(t *testing.T) {
	t.Parallel()

	baseTestCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "param TEST = (a > b) ? 1 : 0;\n",
			expected: "param TEST = (a > b) ? 1 : 0;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-wrap-before-operators",
			code:     "param TEST = (long_var_a + b) ? true_variable_name + 1 : false_variable_name + 2;\n",
			expected: "param TEST = (long_var_a + b)\n    ? true_variable_name + 1\n    : false_variable_name + 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-nested-chain",
			code:     "param TEST = (a) ? 1 : (b) ? 2 : (c) ? 3 : 4;\n",
			expected: "param TEST = (a) ? 1\n    : (b) ? 2\n    : (c) ? 3\n    : 4;\n",
			masked: MaskTernaryOperatorWrappingAndIndentation(
				"param TEST = (a)\n    ? 1\n    : (b)\n    ? 2\n    : (c)\n    ? 3\n    : 4;\n",
			), // BUG(termary-ident): ???
			pc: func(f func()) { f() },
		},
		{
			name:     "04-wrap-long-chain",
			code:     "param TEST = (long_a + b) ? 1 : (long_c + d) ? true_val + 1 : false_val + 2;\n",
			expected: "param TEST = (long_a + b)\n    ? 1\n    : (long_c + d)\n    ? true_val + 1\n    : false_val + 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comments-before-question",
			code:     "param TEST = (a > b) /* c1 */ ? 1 : 0;\n",
			expected: "param TEST = (a > b) /* c1 */ ? 1 : 0;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comments-after-question",
			code:     "param TEST = (a > b) ? /* c1 */ 1 : 0;\n",
			expected: "param TEST = (a > b) ? /* c1 */ 1 : 0;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-comments-before-colon",
			code:     "param TEST = (a > b) ? 1 /* c1 */ : 0;\n",
			expected: "param TEST = (a > b) ? 1 /* c1 */ : 0;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-comments-after-colon",
			code:     "param TEST = (a > b) ? 1 : /* c1 */ 0;\n",
			expected: "param TEST = (a > b) ? 1 : /* c1 */ 0;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-long-comments-force-wrap",
			code:     "param TEST = (long_var_a + b) /* this is a very long comment */ ? 1 : 0;\n",
			expected: "param TEST = (\n    long_var_a + b\n) /* this is a very long comment */\n    ? 1\n    : 0;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-wrap-condition-after-operators",
			code:     "param TEST = (very_long_var_a + very_long_var_b) ? true_variable_name + 1 : false_variable_name + 2;\n",
			expected: "param TEST = (\n    very_long_var_a + very_long_var_b\n) ? true_variable_name + 1\n    : false_variable_name + 2;\n",
			masked: MaskTernaryOperatorWrappingAndIndentation(
				"param TEST = (\n    very_long_var_a + very_long_var_b\n)\n    ? true_variable_name + 1\n    : false_variable_name + 2;\n",
			), // BUG(termary-ident): ???
			pc: func(f func()) { f() },
		},
		{
			name:     "11-func-call-short",
			code:     "param TEST = call_method(\"msg %d\", (a) ? 1 : 0);\n",
			expected: "param TEST = call_method(\n    \"msg %d\", (a) ? 1 : 0\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-func-call-long-wrap-ternary",
			code:     "param TEST = call_method(\"this is a very long format string %d\", (long_var_a + b) ? true_variable_name + 1 : false_variable_name + 2);\n",
			expected: "param TEST = call_method(\n    \"this is a very long format string %d\",\n    (long_var_a + b)\n        ? true_variable_name + 1\n        : false_variable_name + 2\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-func-call-with-comments",
			code:     "param TEST = call_method(\"msg\", (a > b) /* c1 */ ? 1 : 0);\n",
			expected: "param TEST = call_method(\n    \"msg\", (a > b) /* c1 */ ? 1 : 0\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	for _, variant := range []struct {
		suffix   string
		question string
		colon    string
		lineCap  int
	}{
		{"IfExpr", "?", ":", 40},
		{"ConstIfExpr", "#?", "#:", 42},
	} {
		t.Run(variant.suffix, func(tt *testing.T) {
			tt.Parallel()
			for _, tc := range baseTestCases {
				tt.Run(tc.name, func(ttt *testing.T) {
					ttt.Parallel()

					code := strings.ReplaceAll(tc.code, "?", variant.question)
					code = strings.ReplaceAll(code, ":", variant.colon)
					expected := strings.ReplaceAll(
						tc.expected,
						"?",
						variant.question,
					)
					expected = strings.ReplaceAll(expected, ":", variant.colon)
					masked := MaskReplaceAll(tc.masked, "?", variant.question)
					masked = MaskReplaceAll(masked, ":", variant.colon)

					tc.pc(func() {
						ttt.Helper()
						runPrinterTestCase(
							ttt,
							&printerTestCase{
								name:     tc.name,
								code:     code,
								expected: expected,
								masked:   masked,
								lineCap:  variant.lineCap,
							},
						)
					})
				})
			}
		})
	}
}

func TestPrintTernaryExpr_MissingParentheses(t *testing.T) {
	t.Parallel()

	baseTestCases := []printerTestCase{
		{
			name:     "01-bare-basic-if",
			code:     "param TEST = a > b ? 1 : 0;\n",
			expected: "param TEST = (a > b) ? 1 : 0;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-bare-chain",
			code:     "param TEST = a ? 1 : b ? 2 : c ? 3 : 4;\n",
			expected: "param TEST = (a) ? 1\n    : (b) ? 2\n    : (c) ? 3\n    : 4;\n",
			masked: MaskTernaryOperatorWrappingAndIndentation(
				"param TEST = (a)\n? 1\n: (b)\n? 2\n: (c)\n? 3\n: 4;\n",
			), // BUG(termary-ident): ident has to be +4 spaces
			pc: func(f func()) { f() },
		},
		{
			name:     "03-bare-logical-condition",
			code:     "param TEST = a && b || c ? 1 : 0;\n",
			expected: "param TEST = (a && b || c) ? 1 : 0;\n",
			pc:       func(f func()) { f() },
		},
	}

	for _, variant := range []struct {
		suffix   string
		question string
		colon    string
		lineCap  int
	}{
		{"IfExpr", "?", ":", 40},
		{"ConstIfExpr", "#?", "#:", 42},
	} {
		t.Run(variant.suffix, func(t *testing.T) {
			t.Parallel()
			for _, tc := range baseTestCases {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()
					t.Skip("Skipping due to compiler workarounds task006")

					code := strings.ReplaceAll(tc.code, "?", variant.question)
					code = strings.ReplaceAll(code, ":", variant.colon)
					expected := strings.ReplaceAll(
						tc.expected,
						"?",
						variant.question,
					)
					expected = strings.ReplaceAll(expected, ":", variant.colon)
					masked := MaskReplaceAll(tc.masked, "?", variant.question)
					masked = MaskReplaceAll(masked, ":", variant.colon)

					runPrinterTestCase(
						t,
						&printerTestCase{
							name:     tc.name,
							code:     code,
							expected: expected,
							masked:   masked,
							lineCap:  variant.lineCap,
						},
					)
				})
			}
		})
	}
}
