// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"fmt"
	"strings"
	"testing"
)

func TestPrintObjectDecl(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		// 1. Basic Definitions
		{
			name:     "01-basic-semi",
			code:     "   bank my_bank  ;\n",
			expected: "bank my_bank;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-block",
			code:     "bank   my_bank   {\n}\n",
			expected: "bank my_bank {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-basic-block-longname",
			code:     "bank   very_very_long_name_for_the_object_to_force_wrapping   {\n}\n",
			expected: "bank very_very_long_name_for_the_object_to_force_wrapping {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-with-array-list",
			code:     "bank   my_bank [ i < 10 ]  ;\n",
			expected: "bank my_bank[i < 10];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-wrap-array-list",
			code:     "bank my_bank [ i < this_is_a_very_long_expression_exceeding_cap ] {   }\n",
			expected: "bank my_bank[\n    i < this_is_a_very_long_expression_exceeding_cap\n] {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-one-parameter",
			code:     "bank my_bank {\nparam p;\n}\n",
			expected: "bank my_bank {\n    param p;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-three-parameters",
			code:     "bank my_bank {\nparam p1;\nparam p2;\nparam p3;\n}\n",
			expected: "bank my_bank {\n    param p1;\n    param p2;\n    param p3;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-desc",
			code:     "bank    my_bank    `desc`   {\nparam p;\n}\n",
			expected: "bank my_bank `desc` {\n    param p;\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	objectTypes := []string{
		"attribute",
		"bank",
		"connect",
		"event",
		"group",
		"implement",
		"interface",
		"port",
		"subdevice",
	}
	for _, tc := range testCases {
		for _, objType := range objectTypes {
			testName := fmt.Sprintf("%s-%s", tc.name, objType)
			t.Run(testName, func(tt *testing.T) {
				tc.pc(func() {
					if tc.task006Skip {
						tt.Skip(
							"pending printer implementation for array_list/objectBody",
						)
					}
					tt.Helper()

					code := strings.ReplaceAll(tc.code, "bank", objType)
					expected := strings.ReplaceAll(tc.expected, "bank", objType)
					masked := MaskReplaceAll(tc.masked, "bank", objType)
					masked = MaskTransform(masked, func(s string) string {
						if strings.Contains(s, "%s") {
							return fmt.Sprintf(s, "")
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

func TestPrintObjectDecl_Nested(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-nested-object",
			code:     "group my_group {\n%s in_%s{ param p;}\n}\n",
			expected: "group my_group {\n    %s in_%s {\n        param p;\n    }\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-nested-with-params",
			code:     "group my_group {\n%s in_%s{ param p1; param p2;}\n}\n",
			expected: "group my_group {\n    %s in_%s {\n        param p1;\n        param p2;\n    }\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-nested-array-ref",
			code:     "group my_group {\n%s    in_%s [    i <   10   ] {   param    p;}\n}\n",
			expected: "group my_group {\n    %s in_%s[i < 10] {\n        param p;\n    }\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-nested-istemplate",
			code:     "group my_group {\n %s   in_%s   is   my_template  {param   p;}\n}\n",
			expected: "group my_group {\n    %s in_%s is my_template {\n        param p;\n    }\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-nested-istemplate-multi",
			code:     "group my_group {\n%s in_%s is ( t1, t2 ) { param p;}\n}\n",
			expected: "group my_group {\n    %s in_%s is (t1, t2) {\n        param p;\n    }\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	objectTypes := []string{
		"attribute",
		"bank",
		"connect",
		"event",
		"group",
		"implement",
		"interface",
		"port",
		"subdevice",
	}
	for _, tc := range testCases {
		for _, objType := range objectTypes {
			testName := fmt.Sprintf("%s-%s", tc.name, objType)
			t.Run(testName, func(tt *testing.T) {
				tc.pc(func() {
					if tc.task006Skip {
						tt.Skip("pending printer implementation for objectBody")
					}
					tt.Helper()

					code := fmt.Sprintf(tc.code, objType, objType)
					expected := fmt.Sprintf(tc.expected, objType, objType)
					masked := MaskTransform(tc.masked, func(s string) string {
						s = fmt.Sprintf(s, objType, objType)
						if strings.Contains(s, "%s") {
							s = fmt.Sprintf(
								s,
								"",
							) // cleanup any remaining %s just in case
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

func TestPrintObjectDecl_Comments(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-comment-after-semi",
			code:     "bank    my_bank  ;   // comment\n",
			expected: "bank my_bank; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-inline-array",
			code:     "bank   my_bank    /* c1 */   [    /* c2 */    i < 10    /* c3 */   ]    /* c4 */    ;\n",
			expected: "bank my_bank /* c1 */ [\n    /* c2 */ i < 10 /* c3 */] /* c4 */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-before-block",
			code:     "bank my_bank /* c1 */ {\n}\n",
			expected: "bank my_bank /* c1 */ {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-trailing-brace",
			code:     "bank my_bank {      // opening comment\n param p;\n}    // closing comment\n",
			expected: "bank my_bank { // opening comment\n    param p;\n} // closing comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-inside-empty-block",
			code:     "bank my_bank {\n// inner comment\n}\n",
			expected: "bank my_bank {\n    // inner comment\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comment-standalone-line",
			code:     "bank my_bank {\nparam p1;\n// separator comment\nparam p2;\n}\n",
			expected: "bank my_bank {\n    param p1;\n    // separator comment\n    param p2;\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	objectTypes := []string{
		"attribute",
		"bank",
		"connect",
		"event",
		"group",
		"implement",
		"interface",
		"port",
		"subdevice",
	}
	for _, tc := range testCases {
		for _, objType := range objectTypes {
			testName := fmt.Sprintf("%s-%s", tc.name, objType)
			t.Run(testName, func(tt *testing.T) {
				tc.pc(func() {
					if tc.task006Skip {
						tt.Skip(
							"pending printer implementation for array_list/objectBody",
						)
					}
					tt.Helper()

					code := strings.ReplaceAll(tc.code, "bank", objType)
					expected := strings.ReplaceAll(tc.expected, "bank", objType)
					masked := MaskReplaceAll(tc.masked, "bank", objType)
					masked = MaskTransform(masked, func(s string) string {
						if strings.Contains(s, "%s") {
							return fmt.Sprintf(s, "")
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

func TestPrintObjectDecl_ObjectIf(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-if-basic",
			code:     "bank my_bank {\n#if (cond) {\nparam p1;\n}\n}\n",
			expected: "bank my_bank {\n    #if (cond) {\n        param p1;\n    }\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-if-else",
			code:     "bank my_bank {\n#if (cond) {\nparam p1;\n} #else {\nparam p2;\n}\n}\n",
			expected: "bank my_bank {\n    #if (cond) {\n        param p1;\n    } #else {\n        param p2;\n    }\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-if-else-if-chain",
			code:     "bank my_bank {\n#if (c1) {\nparam p1;\n} #else #if (c2) {\nparam p2;\n} #else {\nparam p3;\n}\n}\n",
			expected: "bank my_bank {\n    #if (c1) {\n        param p1;\n    } #else #if (c2) {\n        param p2;\n    } #else {\n        param p3;\n    }\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-nested-if-unbraced",
			code:     "bank my_bank {\n#if (c1)\n#if (c2)\nparam p;\n}\n",
			expected: "bank my_bank {\n    #if (c1) #if (c2) param p;\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
