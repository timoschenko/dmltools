// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintExpressionCalls(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "param TEST =     func    (a,b,c);\n",
			expected: "param TEST = func(a, b, c);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-no-args",
			code:     "param TEST =func(  )  ;\n",
			expected: "param TEST = func();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-single-arg",
			code:     "param TEST = func(  a  );\n",
			expected: "param TEST = func(a);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-long-args",
			code:     "param TEST = func(a, b, c, d, e, f, g, h, i, j);\n",
			expected: "param TEST = func(\n    a, b, c, d, e, f, g, h, i, j\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-with-comments",
			code:     "param TEST = func(a, /* comment */ b, c);\n",
			expected: "param TEST = func(\n    a, /* comment */ b, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-first-arg-quite-long",
			code:     "param TEST = func(quite_long_arg_name, b, c);\n",
			expected: "param TEST = func(\n    quite_long_arg_name, b, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-middle-arg-quite-long",
			code:     "param TEST = func(a, quite_long_arg_name, c);\n",
			expected: "param TEST = func(\n    a, quite_long_arg_name, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-last-arg-quite-long",
			code:     "param TEST = func(a, b, quite_long_arg_name);\n",
			expected: "param TEST = func(\n    a, b, quite_long_arg_name\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-arg-too-long",
			code:     "param TEST = func(a, this_argument_is_way_too_long_for_the_cap, c);\n",
			expected: "param TEST = func(\n    a,\n    this_argument_is_way_too_long_for_the_cap,\n    c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-func-name-too-long",
			code:     "param TEST = this_function_name_is_way_too_long_for_cap(a, b, c);\n",
			expected: "param TEST = this_function_name_is_way_too_long_for_cap(\n    a, b, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-param-name-too-long",
			code:     "param THIS_PARAM_NAME_IS_WAY_TOO_LONG_FOR_THE_CAP = func(a, b, c);\n",
			expected: "param THIS_PARAM_NAME_IS_WAY_TOO_LONG_FOR_THE_CAP = func(\n    a, b, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-member-access",
			code:     "param TEST = this.obj->arr[123].qwe(a,b,c);\n",
			expected: "param TEST = this.obj->arr[123].qwe(\n    a, b, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-member-access-long",
			code:     "param TEST = this.obj->arr[123].qwe(a, quite_long_arg_name, c);\n",
			expected: "param TEST = this.obj->arr[123].qwe(\n    a, quite_long_arg_name, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-nested-calls",
			code:     "param TEST = func(a, inner(b, c), d);\n",
			expected: "param TEST = func(a, inner(b, c), d);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-nested-calls-long",
			code:     "param TEST = func(a, inner(b, quite_long_arg_name, c), d);\n",
			expected: "param TEST = func(\n    a,\n    inner(b, quite_long_arg_name, c),\n    d\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticTrailingComma(
				"param TEST = func(\n    a,\n    inner(\n        b,\n        this_inner_argument_is_way_too_long_for_the_cap,\n        c\n    ),\n    d\n);\n",
			),
			name:     "16-nested-calls-very-long",
			code:     "param TEST = func(a, inner(b, this_inner_argument_is_way_too_long_for_the_cap, c), d);\n",
			expected: "param TEST = func(\n    a,\n    inner(\n        b,\n        this_inner_argument_is_way_too_long_for_the_cap,\n        c,\n    ),\n    d\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionCalls_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "param TEST = func(a, // comment\n b, c);\n",
			expected: "param TEST = func(\n    a, // comment\n    b,\n    c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-no-args",
			code:     "param TEST = func(   // comment\n);\n",
			expected: "param TEST = func( // comment\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-single-arg",
			code:     "param TEST = func(a // comment\n);\n",
			expected: "param TEST = func(\n    a // comment\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-long-args",
			code:     "param TEST = func(a, b, c, d, e, // comment\n f, g, h, i, j);\n",
			expected: "param TEST = func(\n    a,\n    b,\n    c,\n    d,\n    e, // comment\n    f,\n    g,\n    h,\n    i,\n    j\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-with-comments",
			code:     "param TEST = func(a, /* comment */ b, c // comment\n);\n",
			expected: "param TEST = func(\n    a, /* comment */ b, c // comment\n);\n",
			masked: MaskArgumentOverWrappingWithLineComments(
				"param TEST = func(\n    a,\n    /* comment */ b,\n    c // comment\n);\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "06-first-arg-quite-long",
			code:     "param TEST = func(quite_long_arg_name, // comment\n b, c);\n",
			expected: "param TEST = func(\n    quite_long_arg_name, // comment\n    b,\n    c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-middle-arg-quite-long",
			code:     "param TEST = func(a, quite_long_arg_name, // comment\n c);\n",
			expected: "param TEST = func(\n    a,\n    quite_long_arg_name, // comment\n    c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-last-arg-quite-long",
			code:     "param TEST = func(a, b, quite_long_arg_name // comment\n);\n",
			expected: "param TEST = func(\n    a, b, quite_long_arg_name // comment\n);\n",
			masked: MaskArgumentOverWrappingWithLineComments(
				"param TEST = func(\n    a,\n    b,\n    quite_long_arg_name // comment\n);\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "09-arg-too-long",
			code:     "param TEST = func(a, this_argument_is_way_too_long_for_the_cap, // comment\n c);\n",
			expected: "param TEST = func(\n    a,\n    this_argument_is_way_too_long_for_the_cap, // comment\n    c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-func-name-too-long",
			code:     "param TEST = this_function_name_is_way_too_long_for_cap(a, b, // comment\n c);\n",
			expected: "param TEST = this_function_name_is_way_too_long_for_cap(\n    a,\n    b, // comment\n    c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-param-name-too-long",
			code:     "param THIS_PARAM_NAME_IS_WAY_TOO_LONG_FOR_THE_CAP = func(a, // comment\n b, c);\n",
			expected: "param THIS_PARAM_NAME_IS_WAY_TOO_LONG_FOR_THE_CAP = func(\n    a, // comment\n    b,\n    c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-member-access",
			code:     "param TEST = this.obj->arr[123].qwe(a, // comment\n b, c);\n",
			expected: "param TEST = this.obj->arr[123].qwe(\n    a, // comment\n    b,\n    c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-member-access-long",
			code:     "param TEST = this.obj->arr[123].qwe(a, quite_long_arg_name, // comment\n c);\n",
			expected: "param TEST = this.obj->arr[123].qwe(\n    a,\n    quite_long_arg_name, // comment\n    c\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintExpressionCalls_InlineComments(t *testing.T) {
	t.Parallel()

	const lineCap = 40

	testCases := []printerTestCase{
		{
			name:     "01-basic",
			code:     "param TEST = func(a, /* c1 */ b, c);\n",
			expected: "param TEST = func(a, /* c1 */ b, c);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-no-args",
			code:     "param TEST = func(/* c1 */);\n",
			expected: "param TEST = func(/* c1 */);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-single-arg",
			code:     "param TEST = func(a /* c1 */);\n",
			expected: "param TEST = func(a /* c1 */);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-long-args",
			code:     "param TEST = func(a, b, c, d, e, /* c1 */ f, g, h, i, j);\n",
			expected: "param TEST = func(\n    a,\n    b,\n    c,\n    d,\n    e,\n    /* c1 */ f,\n    g,\n    h,\n    i,\n    j\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-with-comments",
			code:     "param TEST = func(a, /* c1 */ b, c /* c2 */);\n",
			expected: "param TEST = func(\n    a, /* c1 */ b, c /* c2 */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-first-arg-quite-long",
			code:     "param TEST = func(quite_long_arg_name /* c1 */, b, c);\n",
			expected: "param TEST = func(\n    quite_long_arg_name /* c1 */, b, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-middle-arg-quite-long",
			code:     "param TEST = func(a, /* c1 */ quite_long_arg_name, c);\n",
			expected: "param TEST = func(\n    a, /* c1 */ quite_long_arg_name, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-last-arg-quite-long",
			code:     "param TEST = func(a, b, quite_long_arg_name /* c1 */);\n",
			expected: "param TEST = func(\n    a, b, quite_long_arg_name /* c1 */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-arg-too-long",
			code:     "param TEST = func(a, this_argument_is_way_too_long_for_the_cap, /* c1 */ c);\n",
			expected: "param TEST = func(\n    a,\n    this_argument_is_way_too_long_for_the_cap,\n    /* c1 */ c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-func-name-too-long",
			code:     "param TEST = this_function_name_is_way_too_long_for_cap(a, /* c1 */ b, c);\n",
			expected: "param TEST = this_function_name_is_way_too_long_for_cap(\n    a, /* c1 */ b, c\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-param-name-too-long",
			code:     "param THIS_PARAM_NAME_IS_WAY_TOO_LONG_FOR_THE_CAP = func(a, b /* c1 */, c);\n",
			expected: "param THIS_PARAM_NAME_IS_WAY_TOO_LONG_FOR_THE_CAP = func(\n    a, b /* c1 */, c\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
