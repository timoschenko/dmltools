// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"fmt"
	"strings"
	"testing"
)

func TestPrintSessionOrSavedDecl(t *testing.T) {
	t.Parallel()

	// TODO: Pending for Synthetic token injection feature
	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "%s int x;\n",
			expected: "%s int x;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-with-initializer",
			code:     "%s int x = 1;\n",
			expected: "%s int x = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-tuple-basic",
			code:     "%s (int x, int y);\n",
			expected: "%s (int x, int y);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-tuple-with-initializer",
			code:     "%s (int x, int y) = {1, 2};\n",
			expected: "%s (int x, int y) = {1, 2};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-tuple-line-wrapping",
			code:     "%s (int this_is_a_very_long_var_name_exceeding_cap, int b, int c);\n",
			expected: "%s (\n    int this_is_a_very_long_var_name_exceeding_cap,\n    int b,\n    int c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-tuple-with-initializer-line-wrapping",
			code:     "%s (int a, int b, int c) = {this_is_a_very_long_initializer_element_that_exceeds_cap, 2, 3};\n",
			expected: "%s (int a, int b, int c) = {\n    this_is_a_very_long_initializer_element_that_exceeds_cap,\n    2,\n    3\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-mid-size-ident",
			code:     "%s int mid_size_ident_name;\n",
			expected: "%s int mid_size_ident_name;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-long-size-ident",
			code:     "%s int this_is_a_very_long_ident_name_that_exceeds_the_line_cap_completely;\n",
			expected: "%s int this_is_a_very_long_ident_name_that_exceeds_the_line_cap_completely;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-mid-size-ident-with-initializer",
			code:     "%s int mid_size_ident_name = 1;\n",
			expected: "%s int mid_size_ident_name = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-long-size-ident-with-initializer",
			code:     "%s int this_is_a_very_long_ident_name_that_exceeds_the_line_cap_completely = 1;\n",
			expected: "%s int this_is_a_very_long_ident_name_that_exceeds_the_line_cap_completely = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-mid-size-type",
			code:     "%s mid_size_type_name x;\n",
			expected: "%s mid_size_type_name x;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-long-size-type",
			code:     "%s this_is_a_very_long_type_name_that_exceeds_the_line_cap_completely x;\n",
			expected: "%s this_is_a_very_long_type_name_that_exceeds_the_line_cap_completely x;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-mid-size-type-with-initializer",
			code:     "%s mid_size_type_name x = 1;\n",
			expected: "%s mid_size_type_name x = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-long-size-type-with-initializer",
			code:     "%s this_is_a_very_long_type_name_that_exceeds_the_line_cap_completely x = 1;\n",
			expected: "%s this_is_a_very_long_type_name_that_exceeds_the_line_cap_completely x = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-binary-expr-wrapping",
			code:     "%s int x = 1 + 2 + 3 + 4 + 5;\n",
			expected: "%s int x = 1\n    + 2\n    + 3\n    + 4\n    + 5;\n",
			lineCap:  20,
			pc:       func(f func()) { f() },
		},
	}

	baseLineCap := 35
	for _, kind := range []string{"session", "saved"} {
		lineCap := baseLineCap + len(kind)
		for _, tc := range testCases {
			testName := fmt.Sprintf("%s_%s", kind, tc.name)
			t.Run(testName, func(tt *testing.T) {
				// TODO: AI WA: pending printer implementation for cdecl and initializer
				if tc.task006Skip {
					tt.Skip(
						"pending printer implementation for cdecl and initializer",
					)
				}
				tc.pc(func() {
					tt.Helper()

					code := fmt.Sprintf(tc.code, kind)
					expected := fmt.Sprintf(tc.expected, kind)

					tcLineCap := lineCap
					if tc.lineCap > 0 {
						tcLineCap = tc.lineCap
					}

					runPrinterTestCase(
						tt,
						&printerTestCase{
							name:     testName,
							code:     code,
							expected: expected,
							masked:   tc.masked,
							lineCap:  tcLineCap,
						},
					)
				})
			})
		}
	}
}

func TestPrintSessionOrSavedDeclComments(t *testing.T) {
	t.Parallel()

	// TODO: AI WA: pending printer implementation for cdecl and initializer
	testCases := []printerTestCase{
		{
			name:     "01-comments-before-semi",
			code:     "%s int x // comment\n;\n",
			expected: "%s int x // comment\n;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comments-in-parens",
			code:     "%s (int x, // comment\nint y);\n",
			expected: "%s (\n    int x, // comment\n    int y\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-long-comment-before-semi",
			code:     "%s int x // this is a very long comment that exceeds the line cap completely\n;\n",
			expected: "%s int x // this is a very long comment that exceeds the line cap completely\n;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-short-comment-after-semi",
			code:     "%s int x; // short comment\n",
			expected: "%s int x; // short comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-long-comment-after-semi",
			code:     "%s int x; // this is a very long comment that exceeds the line cap completely\n",
			expected: "%s int x; // this is a very long comment that exceeds the line cap completely\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-short-comment-before-semi-with-initializer",
			code:     "%s int x = 1 // short comment\n;\n",
			expected: "%s int x = 1 // short comment\n;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-long-comment-before-semi-with-initializer",
			code:     "%s int x = 1 // this is a very long comment that exceeds the line cap completely\n;\n",
			expected: "%s int x = 1 // this is a very long comment that exceeds the line cap completely\n;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-short-comment-after-semi-with-initializer",
			code:     "%s int x = 1; // short comment\n",
			expected: "%s int x = 1; // short comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-long-comment-after-semi-with-initializer",
			code:     "%s int x = 1; // this is a very long comment that exceeds the line cap completely\n",
			expected: "%s int x = 1; // this is a very long comment that exceeds the line cap completely\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-comments-all-over-the-place",
			code:     "dml 1.4;\n\n\n// c1\n%s /* c2 */ int /* c3 */ x /* c4 */ = /* c5 */ 1 /* c6 */ ; // c7\n",
			expected: "dml 1.4;\n\n\n// c1\n%s /* c2 */ int /* c3 */ x /* c4 */ = /* c5 */ 1 /* c6 */; // c7\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "11-line-comments-all-over-the-place",
			code: "dml 1.4;\n\n\n// c1\n%s // c2\nint // c3\nx // c4\n= // c5\n1 // c6\n;\n// c7\n",
			expected: "dml 1.4;\n\n\n// c1\n%s // c2\n" +
				"    int // c3\n" +
				"    x // c4\n" +
				"    = // c5\n" +
				"    1 // c6\n" +
				";\n// c7\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "12-line-comment-after-assign",
			code:     "%s int x = // comment\n 1 + 2;\n",
			expected: "%s int x = // comment\n    1 + 2;\n",
			pc:       func(f func()) { f() },
		},
	}

	baseLineCap := 35
	for _, kind := range []string{"session", "saved"} {
		lineCap := baseLineCap + len(kind)
		for _, tc := range testCases {
			testName := fmt.Sprintf("%s_%s", kind, tc.name)
			t.Run(testName, func(tt *testing.T) {
				// TODO: AI WA: pending printer implementation for cdecl and initializer
				if tc.task006Skip {
					tt.Skip(
						"pending printer implementation for cdecl and initializer",
					)
				}
				tc.pc(func() {
					tt.Helper()

					code := fmt.Sprintf(tc.code, kind)
					expected := fmt.Sprintf(tc.expected, kind)

					masked := MaskTransform(tc.masked, func(s string) string {
						if strings.Contains(s, "%s") {
							return fmt.Sprintf(s, kind)
						}
						return s
					})

					runPrinterTestCase(
						tt,
						&printerTestCase{
							name:     testName,
							code:     code,
							expected: expected,
							masked:   masked,
							lineCap:  lineCap,
						},
					)
				})
			})
		}
	}
}
