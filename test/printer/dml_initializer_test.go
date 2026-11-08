// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintSingleInitializers(t *testing.T) {
	t.Parallel()
	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "saved tp x = {  a   ,  b  ,  c  } ;\n",
			expected: "saved tp x = {a, b, c};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-single-item",
			code:     "saved tp x =   {  a  }  ;\n",
			expected: "saved tp x = {a};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-many-items",
			code:     "saved tp x = {  a  ,   b  ,  c ,  d  , e ,  f  };\n",
			expected: "saved tp x = {a, b, c, d, e, f};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-nested-items",
			code:     "saved tp x = {  a ,  { b ,c} ,  d };\n",
			expected: "saved tp x = {a, {b, c}, d};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-deeply-nested",
			code:     "saved tp x = {  a,{b,{c,d}},   e};\n",
			expected: "saved tp x = {a, {b, {c, d}}, e};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-line-wrapping",
			code:     "saved tp x =  { this_is_a_very_long_element_that_exceeds_cap , b ,  c };\n",
			expected: "saved tp x = {\n    this_is_a_very_long_element_that_exceeds_cap,\n    b,\n    c\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			masked:   MaskSyntheticTrailingComma("saved tp x = {a, b, c,};\n"),
			name:     "07-trailing-comma",
			code:     "saved tp x = {a, b, c,};\n",
			expected: "saved tp x = {a, b, c};\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticTrailingComma(
				"saved tp x = {a, {b, c,}, d};\n",
			),
			name:     "08-trailing-comma-nested",
			code:     "saved tp x = {a, {b, c,}, d};\n",
			expected: "saved tp x = {a, {b, c}, d};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-wrap-one-exceeds",
			code:     "saved tp x = {a, a_very_long_element_that_will_surely_exceed_forty_characters, c};\n",
			expected: "saved tp x = {\n    a,\n    a_very_long_element_that_will_surely_exceed_forty_characters,\n    c\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-wrap-all-together-exceeds",
			code:     "saved tp x = {medium_item_1, medium_item_2, medium_item_3};\n",
			expected: "saved tp x = {\n    medium_item_1,\n    medium_item_2,\n    medium_item_3\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-wrap-prefix-exceeds",
			code:     "saved tp a_long_variable_name_here={a,b,c};\n",
			expected: "saved tp a_long_variable_name_here = {\n    a, b, c\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-wrap-multiple-exceed",
			code:     "saved tp x = {very_long_item_number_one_that_exceeds_the_cap, short_item, very_long_item_number_two_that_exceeds_the_cap};\n",
			expected: "saved tp x = {\n    very_long_item_number_one_that_exceeds_the_cap,\n    short_item,\n    very_long_item_number_two_that_exceeds_the_cap\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-wrap-nested-fits-on-new-line",
			code:     "saved tp a_long_variable_name_here = {a,{b,c},d};\n",
			expected: "saved tp a_long_variable_name_here = {\n    a, {b, c}, d\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-wrap-nested-inner-fits-on-new-line",
			code:     "saved tp x ={a_very_long_element_name_here , { b , c , d }};\n",
			expected: "saved tp x = {\n    a_very_long_element_name_here,\n    {b, c, d}\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-wrap-nested-inner-exceeds",
			code:     "saved tp x = {a, {very_long_element_1, very_long_element_2 }   };\n",
			expected: "saved tp x = {\n    a,\n    {\n        very_long_element_1,\n        very_long_element_2\n    }\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-wrap-deeply-nested-fits",
			code:     "saved tp a_long_variable_name_here = {a,  { b, {c,d} }, e   };\n",
			expected: "saved tp a_long_variable_name_here = {\n    a, {b, {c, d}}, e\n};\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintSingleInitializers_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-beginning",
			code:     "saved tp x = {// comment\n a, b, c};\n",
			expected: "saved tp x = { // comment\n    a, b, c\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-middle",
			code:     "saved tp x = {a, // comment\n b, c};\n",
			expected: "saved tp x = {\n    a, // comment\n    b,\n    c\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-end",
			code:     "saved tp x = {a, b, c // comment\n};\n",
			expected: "saved tp x = {\n    a, b, c // comment\n};\n",
			masked: MaskArgumentOverWrappingWithLineComments(
				"saved tp x = {\n    a,\n    b,\n    c // comment\n};\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "04-before-lbrace",
			code:     "saved tp x = // comment\n{a, b, c};\n",
			expected: "saved tp x = // comment\n{a, b, c};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-after-rbrace",
			code:     "saved tp x = {a, b, c}; // comment\n",
			expected: "saved tp x = {a, b, c}; // comment\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintSingleInitializers_InlineComments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-beginning",
			code:     "saved tp x = {/* comment */ a, b, c};\n",
			expected: "saved tp x = {/* comment */ a, b, c};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-middle",
			code:     "saved tp x = {a, /* comment */ b, c};\n",
			expected: "saved tp x = {a, /* comment */ b, c};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-end",
			code:     "saved tp x = {a, b, c /* comment */};\n",
			expected: "saved tp x = {a, b, c /* comment */};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-within-nested",
			code:     "saved tp x = {a, {b, /* comment */ c}, d};\n",
			expected: "saved tp x = {\n    a, {b, /* comment */ c}, d\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-before-lbrace",
			code:     "saved tp x = /* comment */ {a, b, c};\n",
			expected: "saved tp x = /* comment */ {a, b, c};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-after-rbrace",
			code:     "saved tp x = {a, b, c}; /* comment */\n",
			expected: "saved tp x = {a, b, c}; /* comment */\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintDesignatedInitializers(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "saved tp x = {.foo = a, .bar = b};\n",
			expected: "saved tp x = {.foo = a, .bar = b};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-single-item",
			code:     "saved tp x = {.foo = a};\n",
			expected: "saved tp x = {.foo = a};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-many-items",
			code:     "saved tp x = {.a = 1, .b = 2, .c = 3, .d = 4, .e = 5, .f = 6};\n",
			expected: "saved tp x = {\n    .a = 1,\n    .b = 2,\n    .c = 3,\n    .d = 4,\n    .e = 5,\n    .f = 6\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-nested-items",
			code:     "saved tp x = {.a = 1, .b = {.c = 2, .d = 3}, .e = 4};\n",
			expected: "saved tp x = {\n    .a = 1,\n    .b = {.c = 2, .d = 3},\n    .e = 4\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-deeply-nested",
			code:     "saved tp x = {.a = 1, .b = {.c = 2, .d = {.e = 3, .ff = 4}}, .g = 5};\n",
			expected: "saved tp x = {\n    .a = 1,\n    .b = {\n        .c = 2, .d = {.e = 3, .ff = 4}\n    },\n    .g = 5\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-line-wrapping",
			code:     "saved tp x = {.this_is_a_very_long_element_that_exceeds_cap = 1, .b = 2, .c = 3};\n",
			expected: "saved tp x = {\n    .this_is_a_very_long_element_that_exceeds_cap = 1,\n    .b = 2,\n    .c = 3\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			masked:   MaskSyntheticTrailingComma("saved tp x = {.foo = a,};\n"),
			name:     "07-trailing-comma",
			code:     "saved tp x = {.foo = a,};\n",
			expected: "saved tp x = {.foo = a};\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticTrailingComma(
				"saved tp x = {\n    .a = 1,\n    .b = {.c = 2, .d = 3,},\n    .e = 4\n};\n",
			),
			name:     "08-trailing-comma-nested",
			code:     "saved tp x = {.a = 1, .b = {.c = 2, .d = 3,}, .e = 4};\n",
			expected: "saved tp x = {.a = 1, .b = {.c = 2, .d = 3}, .e = 4};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-wrap-one-exceeds",
			code:     "saved tp x = {.a = 1, .a_very_long_element_that_will_surely_exceed_forty_characters = 2, .c = 3};\n",
			expected: "saved tp x = {\n    .a = 1,\n    .a_very_long_element_that_will_surely_exceed_forty_characters = 2,\n    .c = 3\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-wrap-all-together-exceeds",
			code:     "saved tp x = {.medium_item_1 = 1, .medium_item_2 = 2, .medium_item_3 = 3};\n",
			expected: "saved tp x = {\n    .medium_item_1 = 1,\n    .medium_item_2 = 2,\n    .medium_item_3 = 3\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-wrap-prefix-exceeds",
			code:     "saved tp a_long_variable_name_here = {.a = 1};\n",
			expected: "saved tp a_long_variable_name_here = {\n    .a = 1\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-wrap-multiple-exceed",
			code:     "saved tp x = {.very_long_item_number_one_that_exceeds_the_cap = 1, .short_item = 2, .very_long_item_number_two_that_exceeds_the_cap = 3};\n",
			expected: "saved tp x = {\n    .very_long_item_number_one_that_exceeds_the_cap = 1,\n    .short_item = 2,\n    .very_long_item_number_two_that_exceeds_the_cap = 3\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-wrap-nested-fits-on-new-line",
			code:     "saved tp a_long_variable_name_here = {.a = 1, .nested = {.b = 2, .c = 3}, .d = 4};\n",
			expected: "saved tp a_long_variable_name_here = {\n    .a = 1,\n    .nested = {.b = 2, .c = 3},\n    .d = 4\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-wrap-nested-inner-fits-on-new-line",
			code:     "saved tp x = {.a_very_long_element_name_here = 1, .nested = {.b = 2, .c = 3, .d = 4}};\n",
			expected: "saved tp x = {\n    .a_very_long_element_name_here = 1,\n    .nested = {.b = 2, .c = 3, .d = 4}\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-wrap-nested-inner-exceeds",
			code:     "saved tp x = {.a = 1, .nested = {.very_long_element_1 = 2, .very_long_element_2 = 3}};\n",
			expected: "saved tp x = {\n    .a = 1,\n    .nested = {\n        .very_long_element_1 = 2,\n        .very_long_element_2 = 3\n    }\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-wrap-deeply-nested-fits",
			code:     "saved tp a_long_variable_name_here = {.a = 1, .nested = {.b = 2, .inner = {.c = 3, .d = 4}}, .e = 5};\n",
			expected: "saved tp a_long_variable_name_here = {\n    .a = 1,\n    .nested = {\n        .b = 2,\n        .inner = {.c = 3, .d = 4}\n    },\n    .e = 5\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-trailing-dots",
			code:     "saved tp x = {.foo = a, .bar = b, ...};\n",
			expected: "saved tp x = {.foo = a, .bar = b, ...};\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticTrailingComma(
				"saved tp x = {.foo = a, ...,};\n",
			),
			name:     "18-trailing-dots-comma",
			code:     "saved tp x = {.foo = a, ...,};\n",
			expected: "saved tp x = {.foo = a, ...};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "19-mix-designated-inside-single",
			code:     "saved tp x = {a, {.b = 2, .c = 3}, d};\n",
			expected: "saved tp x = {a, {.b = 2, .c = 3}, d};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "20-mix-single-inside-designated",
			code:     "saved tp x = {.a = 1, .nested = {2, 3}, .d = 4};\n",
			expected: "saved tp x = {\n    .a = 1, .nested = {2, 3}, .d = 4\n};\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintDesignatedInitializers_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-beginning",
			code:     "saved tp x = { // comment\n.a = 1};\n",
			expected: "saved tp x = { // comment\n    .a = 1\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-middle",
			code:     "saved tp x = {.a = 1, // comment\n.b = 2, .c = 3};\n",
			expected: "saved tp x = {\n    .a = 1, // comment\n    .b = 2,\n    .c = 3\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-end",
			code:     "saved tp x = {.a = 1 // comment\n};\n",
			expected: "saved tp x = {\n    .a = 1 // comment\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-before-lbrace",
			code:     "saved tp x = // comment\n{.a = 1};\n",
			expected: "saved tp x = // comment\n{.a = 1};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-after-rbrace",
			code:     "saved tp x = {.a = 1}; // comment\n",
			expected: "saved tp x = {.a = 1}; // comment\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintDesignatedInitializers_InlineComments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-beginning",
			code:     "saved tp x = {/* comment */ .a = 1};\n",
			expected: "saved tp x = {/* comment */ .a = 1};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-middle",
			code:     "saved tp x = {.a = 1, /* comment */ .b = 2, .c = 3};\n",
			expected: "saved tp x = {\n    .a = 1, /* comment */ .b = 2, .c = 3\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-end",
			code:     "saved tp x = {.a = 1 /* comment */};\n",
			expected: "saved tp x = {.a = 1 /* comment */};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-within-nested",
			code:     "saved tp x = {.a = 1, .nested = {.b = 2, /* comment */ .c = 3}, .d = 4};\n",
			expected: "saved tp x = {\n    .a = 1,\n    .nested = {\n        .b = 2, /* comment */ .c = 3\n    },\n    .d = 4\n};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-before-lbrace",
			code:     "saved tp x = /* comment */ {.a = 1};\n",
			expected: "saved tp x = /* comment */ {.a = 1};\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-after-rbrace",
			code:     "saved tp x = {.a = 1}; /* comment */\n",
			expected: "saved tp x = {.a = 1}; /* comment */\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
