// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestStatementSpacing_TopLevelCompact(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-constants-consecutive",
			code:     "constant A = 1;\nconstant B = 2;\nconstant C = 3;\n",
			expected: "constant A = 1;\nconstant B = 2;\nconstant C = 3;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "02-constants-subgroups-preserved",
			code: "constant A = 1;\nconstant B = 2;\n\nconstant X = 10;\nconstant Y = 20;\n",
			expected: "constant A = 1;\nconstant B = 2;\n\n" +
				"constant X = 10;\nconstant Y = 20;\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "03-typedefs-consecutive",
			code:     "typedef int a_t;\ntypedef int b_t;\n",
			expected: "typedef int a_t;\ntypedef int b_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-externs-consecutive",
			code:     "extern void f1();\nextern void f2();\n",
			expected: "extern void f1();\nextern void f2();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-exports-consecutive",
			code:     "export a as b;\nexport c as d;\n",
			expected: "export a as b;\nexport c as d;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-params-toplevel-consecutive",
			code:     "param p1 = 1;\nparam p2 = 2;\n",
			expected: "param p1 = 1;\nparam p2 = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "07-imports-subgroups-preserved",
			code: "import \"simics/api.dml\";\nimport \"simics/devs.dml\";\n\n" +
				"import \"local/reg.dml\";\nimport \"local/util.dml\";\n",
			expected: "import \"simics/api.dml\";\nimport \"simics/devs.dml\";\n\n" +
				"import \"local/reg.dml\";\nimport \"local/util.dml\";\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "08-mixed-compact-import-constant",
			code:     "import \"simics/api.dml\";\nconstant BUFFER = 1024;\n",
			expected: "import \"simics/api.dml\";\n\nconstant BUFFER = 1024;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-constants-two-empty-lines-preserved",
			code:     "constant A = 1;\n\n\nconstant B = 2;\n",
			expected: "constant A = 1;\n\n\nconstant B = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-constants-three-empty-lines-collapsed-to-two",
			code:     "constant A = 1;\n\n\n\nconstant B = 2;\n",
			expected: "constant A = 1;\n\n\nconstant B = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-mixed-compact-two-empty-lines-preserved",
			code:     "import \"simics/api.dml\";\n\n\nconstant BUFFER = 1024;\n",
			expected: "import \"simics/api.dml\";\n\n\nconstant BUFFER = 1024;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-sessions-toplevel-consecutive",
			code:     "session int a = 1;\nsession int b = 2;\n",
			expected: "session int a = 1;\nsession int b = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-sessions-subgroups-preserved",
			code:     "session int a = 1;\n\nsession int b = 2;\n",
			expected: "session int a = 1;\n\nsession int b = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-sessions-two-empty-lines-preserved",
			code:     "session int a = 1;\n\n\nsession int b = 2;\n",
			expected: "session int a = 1;\n\n\nsession int b = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-sessions-three-empty-lines-collapsed-to-two",
			code:     "session int a = 1;\n\n\n\nsession int b = 2;\n",
			expected: "session int a = 1;\n\n\nsession int b = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-saveds-toplevel-consecutive",
			code:     "saved int a = 1;\nsaved int b = 2;\n",
			expected: "saved int a = 1;\nsaved int b = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-mixed-compact-session-saved",
			code:     "session int a = 1;\nsaved int b = 2;\n",
			expected: "session int a = 1;\n\nsaved int b = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-mixed-compact-param-session",
			code:     "param p = 1;\nsession int s = 2;\n",
			expected: "param p = 1;\n\nsession int s = 2;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestStatementSpacing_TopLevelMajor(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-templates-consecutive",
			code:     "template T1 {}\ntemplate T2 {}\n",
			expected: "template T1 {}\n\n\ntemplate T2 {}\n",
			// Note: Current formatter already produces 2 blank lines (\n\n\n) between templates,
			// which matches the expected target output for major statements.
			pc: func(f func()) { f() },
		},
		{
			name:     "02-methods-toplevel-consecutive",
			code:     "method m1() {}\nmethod m2() {}\n",
			expected: "method m1() {\n}\n\n\nmethod m2() {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-mixed-constants-group-then-bank",
			code:     "constant A = 1;\nconstant B = 2;\nbank regs {\n}\n",
			expected: "constant A = 1;\nconstant B = 2;\n\n\nbank regs {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-collapsing-excess-top-level-major",
			code:     "template T1 {}\n\n\n\n\ntemplate T2 {}\n",
			expected: "template T1 {}\n\n\ntemplate T2 {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-collapsing-excess-top-level-compact-subgroups",
			code:     "import \"a.dml\";\n\n\n\n\nimport \"b.dml\";\n",
			expected: "import \"a.dml\";\n\n\nimport \"b.dml\";\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-mixed-session-then-bank",
			code:     "session int a = 1;\nbank b {}\n",
			expected: "session int a = 1;\n\n\nbank b {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-mixed-bank-then-saved",
			code:     "bank b {}\nsaved int a = 1;\n",
			expected: "bank b {\n}\n\n\nsaved int a = 1;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestStatementSpacing_NestedScope(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-methods-in-bank-separated",
			code: "bank b {\n    method m1() {}\n    method m2() {}\n}\n",
			expected: "bank b {\n    method m1() {\n    }\n\n" +
				"    method m2() {\n    }\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "02-methods-in-template-separated",
			code: "template t {\n    method m1() {}\n    method m2() {}\n}\n",
			expected: "template t {\n    method m1() {\n    }\n\n" +
				"    method m2() {\n    }\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "03-bank-params-subgroups-preserved",
			code: "bank b {\n    param p1 = 1;\n    param p2 = 2;\n\n" +
				"    param q1 = 10;\n    param q2 = 20;\n}\n",
			expected: "bank b {\n    param p1 = 1;\n    param p2 = 2;\n\n" +
				"    param q1 = 10;\n    param q2 = 20;\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "04-collapsing-excess-in-bank",
			code:     "bank b {\n    param p1 = 1;\n\n\n\n    param p2 = 2;\n}\n",
			expected: "bank b {\n    param p1 = 1;\n\n    param p2 = 2;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-bank-params-then-method",
			code:     "bank b {\n    param p = 1;\n    method m() {}\n}\n",
			expected: "bank b {\n    param p = 1;\n\n    method m() {\n    }\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "06-bank-sessions-consecutive",
			code: "bank b {\n    session int s1 = 1;\n    session int s2 = 2;\n}\n",
			expected: "bank b {\n    session int s1 = 1;\n" +
				"    session int s2 = 2;\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "07-bank-sessions-one-empty-line-preserved",
			code: "bank b {\n    session int s1 = 1;\n\n    session int s2 = 2;\n}\n",
			expected: "bank b {\n    session int s1 = 1;\n\n" +
				"    session int s2 = 2;\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "08-bank-sessions-collapsing-excess",
			code: "bank b {\n    session int s1 = 1;\n\n\n\n    session int s2 = 2;\n}\n",
			expected: "bank b {\n    session int s1 = 1;\n\n" +
				"    session int s2 = 2;\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "09-bank-saveds-consecutive",
			code: "bank b {\n    saved int s1 = 1;\n    saved int s2 = 2;\n}\n",
			expected: "bank b {\n    saved int s1 = 1;\n" +
				"    saved int s2 = 2;\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "10-bank-mixed-session-and-saved",
			code: "bank b {\n    session int s = 1;\n    saved int v = 2;\n}\n",
			expected: "bank b {\n    session int s = 1;\n\n" +
				"    saved int v = 2;\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "11-bank-mixed-param-and-session",
			code: "bank b {\n    param p = 1;\n    session int s = 2;\n}\n",
			expected: "bank b {\n    param p = 1;\n\n" +
				"    session int s = 2;\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "12-bank-session-then-method",
			code: "bank b {\n    session int s = 1;\n    method m() {}\n}\n",
			expected: "bank b {\n    session int s = 1;\n\n" +
				"    method m() {\n    }\n}\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestStatementSpacing_BracesCleanup(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-bank-leading-trailing-empty-lines",
			code:     "bank b {\n\n    param p = 1;\n\n}\n",
			expected: "bank b {\n    param p = 1;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-template-leading-trailing-empty-lines",
			code:     "template t {\n\n    param p = 1;\n\n}\n",
			expected: "template t {\n    param p = 1;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "03-register-leading-trailing-empty-lines",
			code: "register r size 4 @ 0x0 {\n\n    field f @ [0];\n\n}\n",
			expected: "register r size 4 @ 0x0 {\n" +
				"    field f @ [0];\n}\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "04-bank-session-leading-trailing-empty-lines",
			code:     "bank b {\n\n    session int s = 1;\n\n}\n",
			expected: "bank b {\n    session int s = 1;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-bank-saved-leading-trailing-empty-lines",
			code:     "bank b {\n\n    saved int s = 1;\n\n}\n",
			expected: "bank b {\n    saved int s = 1;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-struct-leading-trailing-empty-lines",
			code:     "typedef struct {\n\n    int a;\n\n} s_t;\n",
			expected: "typedef struct {\n    int a;\n} s_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-union-leading-trailing-empty-lines",
			code:     "typedef union {\n\n    int a;\n\n} u_t;\n",
			expected: "typedef union {\n    int a;\n} u_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-layout-leading-trailing-empty-lines",
			code:     "typedef layout `little-endian` {\n\n    uint8 a;\n\n} l_t;\n",
			expected: "typedef layout `little-endian` {\n    uint8 a;\n} l_t;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-bitfields-leading-trailing-empty-lines",
			code:     "typedef bitfields 32 {\n\n    int a @ [0:4];\n\n} bf_t;\n",
			expected: "typedef bitfields 32 {\n    int a @ [0:4];\n} bf_t;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestStatementSpacing_HeaderBlock(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-compact-header-directives",
			code:     "dml 1.4;\ndevice my_dev;\nbitorder le;\n",
			expected: "dml 1.4;\ndevice my_dev;\nbitorder le;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-header-block-then-import",
			code:     "dml 1.4;\ndevice my_dev;\nbitorder le;\nimport \"api.dml\";\n",
			expected: "dml 1.4;\ndevice my_dev;\nbitorder le;\n\n\nimport \"api.dml\";\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-header-directives-one-empty-line-preserved",
			code:     "dml 1.4;\n\ndevice my_dev;\n\nbitorder le;\n",
			expected: "dml 1.4;\n\ndevice my_dev;\n\nbitorder le;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-header-directives-two-empty-lines-preserved",
			code:     "dml 1.4;\n\n\ndevice my_dev;\n\n\nbitorder le;\n",
			expected: "dml 1.4;\n\n\ndevice my_dev;\n\n\nbitorder le;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-header-directives-excess-collapsed-to-two",
			code:     "dml 1.4;\n\n\n\ndevice my_dev;\n\n\n\nbitorder le;\n",
			expected: "dml 1.4;\n\n\ndevice my_dev;\n\n\nbitorder le;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestStatementSpacing_MethodBody(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-consecutive-statements",
			code: "method m() {\n    local int a = 1;\n    local int b = 2;\n}\n",
			expected: "method m() {\n" +
				"    local int a = 1;\n" +
				"    local int b = 2;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "02-subgroups-preserved",
			code: "method m() {\n    local int a = 1;\n\n    local int b = 2;\n}\n",
			expected: "method m() {\n" +
				"    local int a = 1;\n\n" +
				"    local int b = 2;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "03-collapsing-excess",
			code: "method m() {\n    local int a = 1;\n\n\n\n    local int b = 2;\n}\n",
			expected: "method m() {\n" +
				"    local int a = 1;\n\n" +
				"    local int b = 2;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "04-braces-cleanup",
			code: "method m() {\n\n    local int a = 1;\n\n}\n",
			expected: "method m() {\n" +
				"    local int a = 1;\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
		{
			name: "05-compound-statement-in-if",
			code: "method m() {\n    if (x) {\n        local int a = 1;\n\n        local int b = 2;\n    }\n}\n",
			expected: "method m() {\n" +
				"    if (x) {\n" +
				"        local int a = 1;\n\n" +
				"        local int b = 2;\n" +
				"    }\n" +
				"}\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestStatementSpacing_StructMembers(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-consecutive-members",
			code: "typedef struct {\n    int a;\n    int b;\n} s_t;\n",
			expected: "typedef struct {\n" +
				"    int a;\n" +
				"    int b;\n" +
				"} s_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "02-subgroups-preserved",
			code: "typedef struct {\n    int a;\n\n    int b;\n} s_t;\n",
			expected: "typedef struct {\n" +
				"    int a;\n\n" +
				"    int b;\n" +
				"} s_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "03-collapsing-excess",
			code: "typedef struct {\n    int a;\n\n\n\n    int b;\n} s_t;\n",
			expected: "typedef struct {\n" +
				"    int a;\n\n" +
				"    int b;\n" +
				"} s_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "04-c-bitfields-subgroups-preserved",
			code: "typedef struct {\n    uint32 a : 8;\n\n    uint32 b : 8;\n} s_t;\n",
			expected: "typedef struct {\n" +
				"    uint32 a : 8;\n\n" +
				"    uint32 b : 8;\n" +
				"} s_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "05-nested-struct-subgroups-preserved",
			code: "typedef struct {\n    struct {\n        int x;\n\n        int y;\n    } inner;\n\n    int z;\n} s_t;\n",
			expected: "typedef struct {\n" +
				"    struct {\n" +
				"        int x;\n\n" +
				"        int y;\n" +
				"    } inner;\n\n" +
				"    int z;\n" +
				"} s_t;\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestStatementSpacing_UnionMembers(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-consecutive-members",
			code: "typedef union {\n    int a;\n    float b;\n} u_t;\n",
			expected: "typedef union {\n" +
				"    int a;\n" +
				"    float b;\n" +
				"} u_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "02-subgroups-preserved",
			code: "typedef union {\n    int a;\n\n    float b;\n} u_t;\n",
			expected: "typedef union {\n" +
				"    int a;\n\n" +
				"    float b;\n" +
				"} u_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "03-collapsing-excess",
			code: "typedef union {\n    int a;\n\n\n\n    float b;\n} u_t;\n",
			expected: "typedef union {\n" +
				"    int a;\n\n" +
				"    float b;\n" +
				"} u_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "04-nested-struct-in-union-subgroups-preserved",
			code: "typedef union {\n    struct {\n        uint32 a : 8;\n        uint32 b : 8;\n    };\n\n    uint32 value;\n} u_t;\n",
			expected: "typedef union {\n" +
				"    struct {\n" +
				"        uint32 a : 8;\n" +
				"        uint32 b : 8;\n" +
				"    };\n\n" +
				"    uint32 value;\n" +
				"} u_t;\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestStatementSpacing_LayoutMembers(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-consecutive-members",
			code: "typedef layout `little-endian` {\n    uint8 a;\n    uint8 b;\n} l_t;\n",
			expected: "typedef layout `little-endian` {\n" +
				"    uint8 a;\n" +
				"    uint8 b;\n" +
				"} l_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "02-subgroups-preserved",
			code: "typedef layout `little-endian` {\n    uint8 a;\n\n    uint8 b;\n} l_t;\n",
			expected: "typedef layout `little-endian` {\n" +
				"    uint8 a;\n\n" +
				"    uint8 b;\n" +
				"} l_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "03-collapsing-excess",
			code: "typedef layout `little-endian` {\n    uint8 a;\n\n\n\n    uint8 b;\n} l_t;\n",
			expected: "typedef layout `little-endian` {\n" +
				"    uint8 a;\n\n" +
				"    uint8 b;\n" +
				"} l_t;\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestStatementSpacing_BitfieldMembers(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-consecutive-members",
			code: "typedef bitfields 32 {\n    int a @ [0:4];\n    int b @ [5:8];\n} bf_t;\n",
			expected: "typedef bitfields 32 {\n" +
				"    int a @ [0:4];\n" +
				"    int b @ [5:8];\n" +
				"} bf_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "02-subgroups-preserved",
			code: "typedef bitfields 32 {\n    int a @ [0:4];\n\n    int b @ [5:8];\n} bf_t;\n",
			expected: "typedef bitfields 32 {\n" +
				"    int a @ [0:4];\n\n" +
				"    int b @ [5:8];\n" +
				"} bf_t;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "03-collapsing-excess",
			code: "typedef bitfields 32 {\n    int a @ [0:4];\n\n\n\n    int b @ [5:8];\n} bf_t;\n",
			expected: "typedef bitfields 32 {\n" +
				"    int a @ [0:4];\n\n" +
				"    int b @ [5:8];\n" +
				"} bf_t;\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
