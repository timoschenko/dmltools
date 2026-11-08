// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestLogStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-kinds-basic",
			code: "log   info  :   `msg`   ;\n" +
				"log   error  :   `msg`   ;\n" +
				"log   warning  :   `msg`   ;\n" +
				"log   spec_viol  :   `msg`   ;\n" +
				"log   unimpl  :   `msg`   ;\n" +
				"log   critical  :   `msg`   ;\n" +
				"log   custom_kind  :   `msg`   ;",
			expected: "log info: `msg`;\n" +
				"log error: `msg`;\n" +
				"log warning: `msg`;\n" +
				"log spec_viol: `msg`;\n" +
				"log unimpl: `msg`;\n" +
				"log critical: `msg`;\n" +
				"log custom_kind: `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name:     "02-with-level",
			code:     "log   info  ,   1  :   `msg`   ;",
			expected: "log info, 1: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-with-then-level",
			code:     "log   info  ,   1   then   2  :   `msg`   ;",
			expected: "log info, 1 then 2: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-with-group",
			code:     "log   info  ,   2  ,   my_group  :   `msg`   ;",
			expected: "log info, 2, my_group: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04a-with-multiple-groups",
			code:     "log   info  ,   1  ,   Read   |   Write  :   `msg`   ;",
			expected: "log info, 1, Read | Write: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04b-with-three-groups",
			code:     "log   info  ,   1  ,   R   |   W   |   X  :   `msg`   ;",
			expected: "log info, 1, R | W | X: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04c-with-then-level-and-multiple-groups",
			code:     "log   info  ,   1   then   2  ,   Read   |   Write  :   `msg`   ;",
			expected: "log info, 1 then 2, Read | Write: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "05-with-format-arg",
			code: "log   info  :   `val: %d`  ,   val   ;",
			expected: "log info: `val: %d`,\n" +
				"    val;",
			pc: func(f func()) { f() },
		},
		{
			name: "06-with-multiple-format-args",
			code: "log   info  ,   1  :   `a=%d, b=%s`  ,   a  ,   b   ;",
			expected: "log info, 1: `a=%d, b=%s`,\n" +
				"    a, b;",
			pc: func(f func()) { f() },
		},
		{
			name: "07-line-wrapping-long-args",
			code: "log   info  ,   1  :   `val1=%d, val2=%s`  ,   " +
				"very_long_argument_1  ,   very_long_argument_2   ;",
			expected: "log info, 1: `val1=%d, val2=%s`,\n" +
				"    very_long_argument_1,\n" +
				"    very_long_argument_2;",
			pc: func(f func()) { f() },
		},
		{
			name:     "08-long-level-identifier-no-wrap",
			code:     "log   info  ,   very_long_level_identifier_that_exceeds_line_cap  :   `msg`   ;",
			expected: "log info, very_long_level_identifier_that_exceeds_line_cap: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "09-long-group-identifier-no-wrap",
			code: "log   info  ,   1  ,   " +
				"very_long_group_identifier_that_exceeds_line_cap  :   `msg`   ;",
			expected: "log info, 1, very_long_group_identifier_that_exceeds_line_cap: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "10-long-custom-kind-no-wrap",
			code: "log   very_long_custom_log_kind_identifier_that_exceeds_line_cap  :   `msg`   ;",
			expected: "log very_long_custom_log_kind_identifier_that_exceeds_line_cap: " +
				"`msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "11-user-example-single-arg",
			code: "log info, 1, group: `test %s`, `by test`;",
			expected: "log info, 1, group: `test %s`,\n" +
				"    `by test`;",
			pc: func(f func()) { f() },
		},
		{
			name: "12-user-example-two-args",
			code: "log info, 1, group: `test %s - %d`, `by test`, 42;",
			expected: "log info, 1, group: `test %s - %d`,\n" +
				"    `by test`, 42;",
			pc: func(f func()) { f() },
		},
		{
			name: "13-user-example-three-args",
			code: "log info, 1, group: `test %s - %d %#x`, `by test`, 42, 432;",
			expected: "log info, 1, group: `test %s - %d %#x`,\n" +
				"    `by test`, 42, 432;",
			pc: func(f func()) { f() },
		},
		{
			name: "14-string-literal-kind-basic",
			code: "log   `info`  :   `msg`   ;\n" +
				"log   `custom_kind`  :   `msg`   ;",
			expected: "log `info`: `msg`;\n" +
				"log `custom_kind`: `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name:     "15-string-literal-kind-with-level-and-group",
			code:     "log   `info`  ,   1  ,   my_group  :   `msg`   ;",
			expected: "log `info`, 1, my_group: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "16-string-literal-kind-with-args",
			code: "log   `info`  :   `val: %d`  ,   val   ;",
			expected: "log `info`: `val: %d`,\n" +
				"    val;",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestLogStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-leading-comment",
			code:     "// comment\nlog   info  :   `msg`   ;",
			expected: "// comment\nlog info: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-between-log-and-kind",
			code:     "log   /* c */   info  :   `msg`   ;",
			expected: "log /* c */ info: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-around-commas",
			code:     "log   info   /* c */   ,   1  :   `msg`   ;",
			expected: "log info /* c */, 1: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-before-colon",
			code:     "log   info   /* c */   :   `msg`   ;",
			expected: "log info /* c */: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-after-colon",
			code:     "log   info  :   /* c */   `msg`   ;",
			expected: "log info: /* c */ `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "06-comment-between-format-and-args",
			code: "log   info  :   `val: %d`   /* c */   ,   val   ;",
			expected: "log info: `val: %d` /* c */,\n" +
				"    val;",
			pc: func(f func()) { f() },
		},
		{
			name:     "07-trailing-comment",
			code:     "log   info  :   `msg`   ;   // done",
			expected: "log info: `msg`; // done",
			pc:       func(f func()) { f() },
		},
		{
			name: "08-comment-around-group-bitor",
			code: "log   info  ,   1  ,   Read   /* c1 */   |   /* c2 */   Write  :   `msg`   ;",
			expected: "log info, 1, Read /* c1 */\n" +
				"    | /* c2 */ Write: `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name:     "09-comment-before-string-literal-kind",
			code:     "log   /* c */   `info`  :   `msg`   ;",
			expected: "log /* c */ `info`: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-comment-after-string-literal-kind",
			code:     "log   `info`   /* c */   :   `msg`   ;",
			expected: "log `info` /* c */: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-comment-between-string-literal-kind-and-comma",
			code:     "log   `info`   /* c */   ,   1  :   `msg`   ;",
			expected: "log `info` /* c */, 1: `msg`;",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestLogStatement_ConditionWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - statement and condition fit within lineCap
		{
			name:     "01-level1-simple-relational",
			code:     "log info, (lvl < 10): `msg`;",
			expected: "log info, (lvl < 10): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-level1-simple-logical",
			code:     "log info, (ready && count < 10): `msg`;",
			expected: "log info, (ready && count < 10): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-level1-ternary",
			code:     "log info, (flag ? 1 : 2): `msg`;",
			expected: "log info, (flag ? 1 : 2): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-level1-method-call",
			code:     "log info, (get_level(dev)): `msg`;",
			expected: "log info, (get_level(dev)): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-level1-then-level",
			code:     "log info, 1 then (flag ? 2 : 3): `msg`;",
			expected: "log info, 1 then (flag ? 2 : 3): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-level1-log-group-ternary",
			code:     "log info, 1, (f ? a : b): `m`;",
			expected: "log info, 1, (f ? a : b): `m`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-level1-log-group-bitwise",
			code:     "log info, 1, ((fl & M) != 0): `m`;",
			expected: "log info, 1, ((fl & M) != 0): `m`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-level1-bitwise-slice",
			code:     "log info, ((val[3:0] & M) != 0): `msg`;",
			expected: "log info, ((val[3:0] & M) != 0): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-level1-nested-parens",
			code:     "log info, ((a && b) || (c && d)): `m`;",
			expected: "log info, ((a && b) || (c && d)): `m`;",
			pc:       func(f func()) { f() },
		},

		// Level 2: Outer Parenthesized Wrap - statement exceeds lineCap, but inner condition fits on one line
		{
			name: "10-level2-outer-wrap-single-relational",
			code: "log info, (status == STATUS_WAITING_FOR_DEVICE): `msg`;",
			expected: "log info, (\n" +
				"    status == STATUS_WAITING_FOR_DEVICE\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-outer-wrap-compound-logical",
			code: "log info, (flag_alpha && flag_beta && is_ready): `msg`;",
			expected: "log info, (\n" +
				"    flag_alpha && flag_beta && is_ready\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-outer-wrap-ternary",
			code: "log info, (flag ? LEVEL_FAST : LEVEL_SLOW): `msg`;",
			expected: "log info, (\n" +
				"    flag ? LEVEL_FAST : LEVEL_SLOW\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-outer-wrap-method-call",
			code: "log info, (get_custom_level(dev, ctx)): `msg`;",
			expected: "log info, (\n" +
				"    get_custom_level(dev, ctx)\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "14-level2-outer-wrap-then-level",
			code: "log info, 1 then (flag ? LEVEL_FAST : LEVEL_SLOW): `msg`;",
			expected: "log info, 1 then (\n" +
				"    flag ? LEVEL_FAST : LEVEL_SLOW\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "15-level2-outer-wrap-log-group",
			code: "log info, 1, (flag ? GROUP_ALPHA : GROUP_BETA): `msg`;",
			expected: "log info, 1, (\n" +
				"    flag ? GROUP_ALPHA : GROUP_BETA\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "16-level2-outer-wrap-nested-parens",
			code: "log info, ((cond_a && cond_b) || (cond_c && cond_d)): `msg`;",
			expected: "log info, (\n" +
				"    (cond_a && cond_b)\n" +
				"    || (cond_c && cond_d)\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},

		// Level 3: Deep Expression Explosion - both statement and condition exceed lineCap
		{
			name: "20-level3-deep-logical-explosion",
			code: "log info, (very_long_condition_name_a && very_long_condition_name_b): `msg`;",
			expected: "log info, (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-deep-arithmetic-explosion",
			code: "log info, (very_long_level_offset_expression_a + very_long_level_offset_expression_b): `msg`;",
			expected: "log info, (\n" +
				"    very_long_level_offset_expression_a\n" +
				"    + very_long_level_offset_expression_b\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-deep-bitwise-log-group-explosion",
			code: "log info, 1, (very_long_group_bitmask_flag_a | very_long_group_bitmask_flag_b): `msg`;",
			expected: "log info, 1, (\n" +
				"    very_long_group_bitmask_flag_a\n" +
				"    | very_long_group_bitmask_flag_b\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-deep-ternary-explosion",
			code: "log info, (use_extended_debug ? LEVEL_TRACE : LEVEL_DEBUG): `msg`;",
			expected: "log info, (\n" +
				"    use_extended_debug\n" +
				"        ? LEVEL_TRACE\n" +
				"        : LEVEL_DEBUG\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "24-level3-deep-method-call-args-explosion",
			code: "log info, (calculate_log_level_for_device(dev, ctx)): `msg`;",
			expected: "log info, (\n" +
				"    calculate_log_level_for_device(\n" +
				"        dev, ctx\n" +
				"    )\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestLogStatement_LogGroups(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Multiple log groups without parens: fits on single line or wraps only after : fmtstr
		{
			name: "01-no-parens-two-groups-with-args-wrap-after-fmtstr",
			code: "log info, 1, Read | Write: `Access to %#x`, addr;",
			expected: "log info, 1, Read\n" +
				"    | Write: `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},
		{
			name: "02-no-parens-two-groups-multiple-args-wrap-after-fmtstr",
			code: "log info, 1, Read | Write: `Access to %#x, val=%d`, addr, value;",
			expected: "log info, 1, Read\n" +
				"    | Write: `Access to %#x, val=%d`,\n" +
				"    addr, value;",
			pc: func(f func()) { f() },
		},
		{
			name: "03-no-parens-long-groups-no-wrap",
			code: "log info, 1, Register_Read | Register_Write: `Access to %#x`, addr;",
			expected: "log info, 1, Register_Read\n" +
				"    | Register_Write: `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},
		{
			name: "04-no-parens-three-long-groups-no-wrap",
			code: "log info, 1, Register_Read | Register_Write | Register_Execute: `msg`;",
			expected: "log info, 1, Register_Read\n" +
				"    | Register_Write\n" +
				"    | Register_Execute: `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "05-no-parens-long-level-and-groups-no-wrap",
			code: "log info, base_level + offset_level, Register_Read | Register_Write: `msg`;",
			expected: "log info, base_level\n" +
				"    + offset_level, Register_Read\n" +
				"    | Register_Write: `msg`;",
			pc: func(f func()) { f() },
		},

		// Level 1: Multiple log groups in parens () fitting on a single line
		{
			name:     "10-parens-level1-two-groups",
			code:     "log info, 1, (Read | Write): `msg`;",
			expected: "log info, 1, (Read | Write): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-parens-level1-three-groups",
			code:     "log info, 1, (R | W | X): `msg`;",
			expected: "log info, 1, (R | W | X): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "12-parens-level1-groups-fit-args-wrap",
			code: "log info, 1, (R | W): `v=%d`, very_long_argument_1, very_long_argument_2;",
			expected: "log info, 1, (R | W): `v=%d`,\n" +
				"    very_long_argument_1,\n" +
				"    very_long_argument_2;",
			pc: func(f func()) { f() },
		},

		// Level 2: Outer parenthesized wrap - statement exceeds lineCap, inner groups fit on one line
		{
			name: "19-parens-level2-short-groups-outer-wrap-due-to-fmtstr",
			code: "log info, 1, (Read | Write): `Access to %#x`, addr;",
			expected: "log info, 1, (\n" +
				"    Read | Write\n" +
				"): `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},
		{
			name: "20-parens-level2-two-groups-outer-wrap",
			code: "log info, 1, (Register_Read | Register_Write): `msg`;",
			expected: "log info, 1, (\n" +
				"    Register_Read | Register_Write\n" +
				"): `msg`;",
			pc: func(f func()) { f() },
		},
		{
			name: "21-parens-level2-two-groups-outer-wrap-with-args",
			code: "log info, 1, (Register_Read | Register_Write): `Access to %#x`, addr;",
			expected: "log info, 1, (\n" +
				"    Register_Read | Register_Write\n" +
				"): `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},
		{
			name: "22-parens-level2-three-groups-outer-wrap",
			code: "log info, 1, (Grp_Read | Grp_Write | Grp_Exec): `Access to %#x`, addr;",
			expected: "log info, 1, (\n" +
				"    Grp_Read | Grp_Write | Grp_Exec\n" +
				"): `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},
		{
			name: "23-parens-level2-then-level-with-groups-outer-wrap",
			code: "log info, 1 then 2, (Register_Read | Register_Write): `Access to %#x`, addr;",
			expected: "log info, 1 then 2, (\n" +
				"    Register_Read | Register_Write\n" +
				"): `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},

		// Level 3: Deep expression explosion - inner group expression also exceeds lineCap
		{
			name: "30-parens-level3-two-long-groups-explosion",
			code: "log info, 1, (Register_Read_Group | Register_Write_Group): `Access to %#x`, addr;",
			expected: "log info, 1, (\n" +
				"    Register_Read_Group\n" +
				"    | Register_Write_Group\n" +
				"): `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},
		{
			name: "31-parens-level3-three-long-groups-explosion",
			code: "log info, 1, (Register_Read_Group | Register_Write_Group | Register_Execute_Group): " +
				"`Access to %#x`, addr;",
			expected: "log info, 1, (\n" +
				"    Register_Read_Group\n" +
				"    | Register_Write_Group\n" +
				"    | Register_Execute_Group\n" +
				"): `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},
		{
			name: "32-parens-level3-three-long-groups-and-long-args-explosion",
			code: "log info, 1, (Register_Read_Group | Register_Write_Group | Register_Execute_Group): " +
				"`Access to %#x, val=%d`, target_register_address, target_register_value;",
			expected: "log info, 1, (\n" +
				"    Register_Read_Group\n" +
				"    | Register_Write_Group\n" +
				"    | Register_Execute_Group\n" +
				"): `Access to %#x, val=%d`,\n" +
				"    target_register_address,\n" +
				"    target_register_value;",
			pc: func(f func()) { f() },
		},
		{
			name: "33-parens-level3-nested-mask-and-groups-explosion",
			code: "log info, 1, ((Register_Read_Group | Register_Write_Group) & Active_Group_Mask): " +
				"`Access to %#x`, addr;",
			expected: "log info, 1, (\n" +
				"    (\n" +
				"        Register_Read_Group\n" +
				"        | Register_Write_Group\n" +
				"    )\n" +
				"    & Active_Group_Mask\n" +
				"): `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},
		{
			name: "34-parens-level3-conditional-group-with-bitor-explosion",
			code: "log info, 1, (Common_Log_Group | (is_write ? Register_Write : Register_Read)): " +
				"`Access to %#x`, addr;",
			expected: "log info, 1, (\n" +
				"    Common_Log_Group\n" +
				"    | (\n" +
				"        is_write\n" +
				"            ? Register_Write\n" +
				"            : Register_Read\n" +
				"    )\n" +
				"): `Access to %#x`,\n" +
				"    addr;",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestLogStatement_FormatStringWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Format string concatenation without parens
		{
			name:     "01-no-parens-short-concat-fits",
			code:     "log info: `hello ` + `world`;",
			expected: "log info: `hello ` + `world`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "02-no-parens-multi-concat-wrapping",
			code: "log info: `first string part that is long ` + `second string part that is long`;",
			expected: "log info: `first string part that is long `\n" +
				"    + `second string part that is long`;",
			pc: func(f func()) { f() },
		},
		{
			name: "02a-no-parens-three-part-concat-wrapping",
			code: "log info: `first part of message ` + `second part of message ` " +
				"+ `third part of message`;",
			expected: "log info: `first part of message `\n" +
				"    + `second part of message `\n" +
				"    + `third part of message`;",
			pc: func(f func()) { f() },
		},
		{
			name: "03-no-parens-with-level-and-group-concat-wrapping",
			code: "log info, 1, my_group: `first part of message ` + `second part of message`;",
			expected: "log info, 1, my_group: `first part of message `\n" +
				"    + `second part of message`;",
			pc: func(f func()) { f() },
		},

		// Format string concatenation with parens (bracketed string literal)
		{
			name:     "04-parens-level1-fits",
			code:     "log info: (`hello ` + `world`);",
			expected: "log info: (`hello ` + `world`);",
			pc:       func(f func()) { f() },
		},
		{
			name: "05-parens-level2-outer-wrap",
			code: "log very_long_custom_kind: (`first part ` + `second part`);",
			expected: "log very_long_custom_kind: (\n" +
				"    `first part ` + `second part`\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "06-parens-level3-deep-explosion",
			code: "log info: (`first part of long message ` + `second part of long message`);",
			expected: "log info: (\n" +
				"    `first part of long message `\n" +
				"    + `second part of long message`\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "07-parens-level3-with-args",
			code: "log info: (`first part of long message ` + `second part of long message`), " +
				"very_long_argument_1, very_long_argument_2;",
			expected: "log info: (\n" +
				"    `first part of long message `\n" +
				"    + `second part of long message`\n" +
				"),\n" +
				"    very_long_argument_1,\n" +
				"    very_long_argument_2;",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestLogStatement_CombinedWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level in () wraps, format string fits
		{
			name: "01-level-wraps-format-fits",
			code: "log info, (status == STATUS_WAITING_FOR_DEVICE): `device ready`;",
			expected: "log info, (\n" +
				"    status == STATUS_WAITING_FOR_DEVICE\n" +
				"): `device ready`;",
			pc: func(f func()) { f() },
		},

		// Level in () wraps, format string wraps without parens
		{
			name: "02-level-wraps-format-wraps-no-parens",
			code: "log info, (status == STATUS_WAITING_FOR_DEVICE): `device status changed to ` " +
				"+ `waiting for device state`;",
			expected: "log info, (\n" +
				"    status == STATUS_WAITING_FOR_DEVICE\n" +
				"): `device status changed to `\n" +
				"    + `waiting for device state`;",
			pc: func(f func()) { f() },
		},

		// Level fits, group in () wraps
		{
			name: "03-level-fits-group-wraps",
			code: "log info, 1, (flag ? GROUP_ALPHA : GROUP_BETA): `status ok`;",
			expected: "log info, 1, (\n" +
				"    flag ? GROUP_ALPHA : GROUP_BETA\n" +
				"): `status ok`;",
			pc: func(f func()) { f() },
		},

		// Both level and group in () wrap
		{
			name: "04-both-level-and-group-wrap",
			code: "log info, (status == STATUS_WAITING_FOR_DEVICE), " +
				"(flag ? GROUP_ALPHA : GROUP_BETA): `device ready`;",
			expected: "log info, (\n" +
				"    status == STATUS_WAITING_FOR_DEVICE\n" +
				"), (\n" +
				"    flag ? GROUP_ALPHA : GROUP_BETA\n" +
				"): `device ready`;",
			pc: func(f func()) { f() },
		},

		// Level in () wraps, format arguments wrap after : fmtstr
		{
			name: "05-level-wraps-format-args-wrap",
			code: "log info, (status == STATUS_WAITING_FOR_DEVICE): `status: %d, %s`, " +
				"very_long_device_status_arg, very_long_status_string_arg;",
			expected: "log info, (\n" +
				"    status == STATUS_WAITING_FOR_DEVICE\n" +
				"): `status: %d, %s`,\n" +
				"    very_long_device_status_arg,\n" +
				"    very_long_status_string_arg;",
			pc: func(f func()) { f() },
		},

		// Format arguments deep explosion after : fmtstr
		{
			name: "06-format-args-deep-explosion",
			code: "log info, 1, my_group: `val1=%d, val2=%s, val3=%x`, " +
				"very_long_argument_alpha_1, very_long_argument_beta_2, " +
				"very_long_argument_gamma_3;",
			expected: "log info, 1, my_group: `val1=%d, val2=%s, val3=%x`,\n" +
				"    very_long_argument_alpha_1,\n" +
				"    very_long_argument_beta_2,\n" +
				"    very_long_argument_gamma_3;",
			pc: func(f func()) { f() },
		},

		// Format string wraps without parens and format arguments wrap
		{
			name: "07-format-wraps-no-parens-with-args",
			code: "log info: `device error code: %d, ` + `additional info: %s`, " +
				"very_long_device_error_code, very_long_device_error_details;",
			expected: "log info: `device error code: %d, `\n" +
				"    + `additional info: %s`,\n" +
				"    very_long_device_error_code,\n" +
				"    very_long_device_error_details;",
			pc: func(f func()) { f() },
		},

		// Level in (), multiple groups with | in (), format string concat, and format args all wrap
		{
			name: "08-level-and-multi-group-and-format-and-args-wrap",
			code: "log info, (status == STATUS_WAITING_FOR_DEVICE), " +
				"(Register_Read | Register_Write): " +
				"`Access to register %#x ` + `with device status %d`, " +
				"target_register_address_val, current_device_status_val;",
			expected: "log info, (\n" +
				"    status == STATUS_WAITING_FOR_DEVICE\n" +
				"), (\n" +
				"    Register_Read | Register_Write\n" +
				"): `Access to register %#x `\n" +
				"    + `with device status %d`,\n" +
				"    target_register_address_val,\n" +
				"    current_device_status_val;",
			pc: func(f func()) { f() },
		},

		{
			name:    "09-hello-world-format-and-args-wrap",
			lineCap: 12,
			code:    "log info: `Hello` + `World %d %s`, 1234, `!`;",
			expected: "log info: `Hello`\n" +
				"    + `World %d %s`,\n" +
				"    1234,\n" +
				"    `!`;",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestLogStatement_RichExpressions(t *testing.T) {
	t.Parallel()

	const lineCap = 80
	testCases := []printerTestCase{
		// Section 1: Member Access & Field Navigation in Format Arguments
		{
			name: "01-format-arg-member-dot",
			code: "log info: `val: %d`,   this  .  val   ;",
			expected: "log info: `val: %d`,\n" +
				"    this.val;",
			pc: func(f func()) { f() },
		},
		{
			name: "02-format-arg-member-arrow",
			code: "log info: `flag: %d`,   ctx  ->  flags   ;",
			expected: "log info: `flag: %d`,\n" +
				"    ctx->flags;",
			pc: func(f func()) { f() },
		},
		{
			name: "03-format-arg-member-chained",
			code: "log info: `name: %s`,   dev  .  bank  .  reg  .  name   ;",
			expected: "log info: `name: %s`,\n" +
				"    dev.bank.reg.name;",
			pc: func(f func()) { f() },
		},
		{
			name: "04-format-arg-member-mixed-chain",
			code: "log info: `status: %d`,   this  .  bank  ->  status  .  code   ;",
			expected: "log info: `status: %d`,\n" +
				"    this.bank->status.code;",
			pc: func(f func()) { f() },
		},

		// Section 2: Array Subscripting & Bit Slicing in Format Arguments
		{
			name: "05-format-arg-array-subscript",
			code: "log info: `buf[%d]: %d`,   index  ,   buffer  [  index  ]   ;",
			expected: "log info: `buf[%d]: %d`,\n" +
				"    index, buffer[index];",
			pc: func(f func()) { f() },
		},
		{
			name: "06-format-arg-array-multidimensional",
			code: "log info: `cell: %d`,   matrix  [  row  ]  [  col  ]   ;",
			expected: "log info: `cell: %d`,\n" +
				"    matrix[row][col];",
			pc: func(f func()) { f() },
		},
		{
			name: "07-format-arg-bit-slice-range",
			code: "log info: `slice: %#x`,   data  [  31  :  0  ]   ;",
			expected: "log info: `slice: %#x`,\n" +
				"    data[31:0];",
			pc: func(f func()) { f() },
		},
		{
			name: "08-format-arg-member-bit-slice",
			code: "log info: `reg bits: %#x`,   this  .  status  [  7  :  0  ]   ;",
			expected: "log info: `reg bits: %#x`,\n" +
				"    this.status[7:0];",
			pc: func(f func()) { f() },
		},

		// Section 3: Function & Method Calls in Format Arguments
		{
			name: "09-format-arg-call-no-args",
			code: "log info: `status: %d`,   get_status  (  )   ;",
			expected: "log info: `status: %d`,\n" +
				"    get_status();",
			pc: func(f func()) { f() },
		},
		{
			name: "10-format-arg-call-with-args",
			code: "log info: `crc: %#x`,   compute_crc  (  buffer  ,  len  )   ;",
			expected: "log info: `crc: %#x`,\n" +
				"    compute_crc(buffer, len);",
			pc: func(f func()) { f() },
		},
		{
			name: "11-format-arg-call-member-method",
			code: "log info: `reg val: %#x`,   dev  .  read_register  (  offset  )   ;",
			expected: "log info: `reg val: %#x`,\n" +
				"    dev.read_register(offset);",
			pc: func(f func()) { f() },
		},
		{
			name: "12-format-arg-call-nested",
			code: "log info: `desc: %s`,   get_desc  (  dev  .  get_target  (  )  )   ;",
			expected: "log info: `desc: %s`,\n" +
				"    get_desc(dev.get_target());",
			pc: func(f func()) { f() },
		},

		// Section 4: Unary & Address Operations in Format Arguments
		{
			name: "13-format-arg-unary-arithmetic",
			code: "log info: `delta: %d, val: %d`,   -  offset  ,   +  val   ;",
			expected: "log info: `delta: %d, val: %d`,\n" +
				"    -offset, +val;",
			pc: func(f func()) { f() },
		},
		{
			name: "14-format-arg-unary-logical-and-bitwise-not",
			code: "log info: `not: %d, inv: %#x`,   !  is_ready  ,   ~  mask   ;",
			expected: "log info: `not: %d, inv: %#x`,\n" +
				"    !is_ready, ~mask;",
			pc: func(f func()) { f() },
		},
		{
			name: "15-format-arg-unary-address-and-deref",
			code: "log info: `ptr: %p, val: %d`,   &  var  ,   *  ptr   ;",
			expected: "log info: `ptr: %p, val: %d`,\n" +
				"    &var, *ptr;",
			pc: func(f func()) { f() },
		},
		{
			name: "16-format-arg-unary-increment-decrement",
			code: "log info: `cnt: %d, step: %d`,   count  ++  ,   ++  total   ;",
			expected: "log info: `cnt: %d, step: %d`,\n" +
				"    count++, ++total;",
			pc: func(f func()) { f() },
		},

		// Section 5: Binary Operations in Format Arguments
		{
			name: "17-format-arg-binary-arithmetic",
			code: "log info: `offset: %#x`,   base_addr  +  index  *  4   ;",
			expected: "log info: `offset: %#x`,\n" +
				"    base_addr + index * 4;",
			pc: func(f func()) { f() },
		},
		{
			name: "18-format-arg-binary-bitwise-and-shift",
			code: "log info: `byte: %#x`,   (  val  &  0xff00  )  >>  8   ;",
			expected: "log info: `byte: %#x`,\n" +
				"    (val & 0xff00) >> 8;",
			pc: func(f func()) { f() },
		},
		{
			name: "19-format-arg-binary-relational-and-equality",
			code: "log info: `cmp: %d, eq: %d`,   count  >  0  ,   status  ==  OK   ;",
			expected: "log info: `cmp: %d, eq: %d`,\n" +
				"    count > 0, status == OK;",
			pc: func(f func()) { f() },
		},
		{
			name: "20-format-arg-binary-logical",
			code: "log info: `valid: %d`,   is_valid  &&  !  has_error   ;",
			expected: "log info: `valid: %d`,\n" +
				"    is_valid && !has_error;",
			pc: func(f func()) { f() },
		},

		// Section 6: Ternary Conditional Expressions in Format Arguments
		{
			name: "21-format-arg-ternary-simple",
			code: "log info: `state: %s`,   is_enabled  ?  `enabled`  :  `disabled`   ;",
			expected: "log info: `state: %s`,\n" +
				"    is_enabled ? `enabled` : `disabled`;",
			pc: func(f func()) { f() },
		},
		{
			name: "22-format-arg-ternary-parenthesized",
			code: "log info: `result: %d`,   (  val  >  10  ?  val  *  2  :  0  )   ;",
			expected: "log info: `result: %d`,\n" +
				"    (val > 10 ? val * 2 : 0);",
			pc: func(f func()) { f() },
		},
		{
			name: "23-format-arg-ternary-multiple",
			code: "log info: `a=%s, b=%s`,   ok  ?  `yes`  :  `no`  ,   " +
				"flag  ?  `high`  :  `low`   ;",
			expected: "log info: `a=%s, b=%s`,\n" +
				"    ok ? `yes` : `no`, flag ? `high` : `low`;",
			pc: func(f func()) { f() },
		},

		// Section 7: DML Built-in Expressions in Format Arguments
		{
			name: "24-format-arg-cast-expression",
			code: "log info: `ptr: %p, val: %u`,   cast  (  buf  ,  void  *  )  ,   " +
				"cast  (  signed_val  ,  uint32  )   ;",
			expected: "log info: `ptr: %p, val: %u`,\n" +
				"    cast(buf, void *), cast(signed_val, uint32);",
			pc: func(f func()) { f() },
		},
		{
			name: "25-format-arg-sizeof-and-sizeoftype",
			code: "log info: `buf_sz: %d, type_sz: %d`,   sizeof  (  buffer  )  ,   " +
				"sizeoftype  (  uint64  )   ;",
			expected: "log info: `buf_sz: %d, type_sz: %d`,\n" +
				"    sizeof(buffer), sizeoftype(uint64);",
			pc: func(f func()) { f() },
		},
		{
			name: "26-format-arg-defined-expression",
			code: "log info: `defined: %d`,   defined  (  feature_flag  )   ;",
			expected: "log info: `defined: %d`,\n" +
				"    defined(feature_flag);",
			pc: func(f func()) { f() },
		},
		{
			name: "27-format-arg-stringify-expression",
			code: "log info: `name: %s`,   stringify  (  device_id  )   ;",
			expected: "log info: `name: %s`,\n" +
				"    stringify(device_id);",
			pc: func(f func()) { f() },
		},

		// Section 8: Array Literals in Format Arguments
		{
			name: "28-format-arg-array-literal",
			code: "log info: `list: %s`,   [  1  ,  2  ,  3  ]   ;",
			expected: "log info: `list: %s`,\n" +
				"    [1, 2, 3];",
			pc: func(f func()) { f() },
		},

		// Section 9: Rich Expressions in Log Level & Then-Level
		{
			name:     "29-level-member-access",
			code:     "log info  ,   dev  .  log_level  :   `msg`   ;",
			expected: "log info, dev.log_level: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "30-level-function-call",
			code:     "log info  ,   get_log_level  (  )  :   `msg`   ;",
			expected: "log info, get_log_level(): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "31-level-array-access",
			code:     "log info  ,   levels  [  idx  ]  :   `msg`   ;",
			expected: "log info, levels[idx]: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "32-level-cast",
			code:     "log info  ,   cast  (  lvl  ,  uint32  )  :   `msg`   ;",
			expected: "log info, cast(lvl, uint32): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "33-then-level-member-access",
			code: "log info  ,   this  .  min_level   then   " +
				"this  .  max_level  :   `msg`   ;",
			expected: "log info, this.min_level then this.max_level: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "34-then-level-function-call",
			code: "log info  ,   1   then   " +
				"calc_next_level  (  dev  )  :   `msg`   ;",
			expected: "log info, 1 then calc_next_level(dev): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "35-then-level-arithmetic",
			code:     "log info  ,   1   then   base_level  +  1  :   `msg`   ;",
			expected: "log info, 1 then base_level + 1: `msg`;",
			pc:       func(f func()) { f() },
		},

		// Section 10: Rich Expressions in Log Group
		{
			name:     "36-group-member-access",
			code:     "log info  ,   1  ,   dev  .  groups  .  dma  :   `msg`   ;",
			expected: "log info, 1, dev.groups.dma: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "37-group-array-access",
			code:     "log info  ,   1  ,   groups  [  i  ]  :   `msg`   ;",
			expected: "log info, 1, groups[i]: `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "38-group-function-call",
			code:     "log info  ,   1  ,   get_active_group  (  )  :   `msg`   ;",
			expected: "log info, 1, get_active_group(): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "39-group-cast",
			code:     "log info  ,   1  ,   cast  (  grp  ,  uint32  )  :   `msg`   ;",
			expected: "log info, 1, cast(grp, uint32): `msg`;",
			pc:       func(f func()) { f() },
		},
		{
			name: "40-group-bitwise-with-not",
			code: "log info  ,   1  ,   " +
				"(  All_Groups  &  ~  Disabled_Group  )  :   `msg`   ;",
			expected: "log info, 1, (All_Groups & ~Disabled_Group): `msg`;",
			pc:       func(f func()) { f() },
		},

		// Section 11: Complex Combinations
		{
			name: "41-complex-all-parts-rich",
			code: "log info  ,   dev  .  get_level  (  )  ,   " +
				"(  Group_A  |  Group_B  )  :   " +
				"`access %s[%d] = %#x`  ,   this  .  name  ,   index  ,   " +
				"buffer  [  index  ]   ;",
			expected: "log info, dev.get_level(), (Group_A | Group_B): `access %s[%d] = %#x`,\n" +
				"    this.name, index, buffer[index];",
			pc: func(f func()) { f() },
		},
		{
			name: "42-complex-deep-subscript-cast-and-call",
			code: "log info: `entry: %#x`,   " +
				"dev  .  banks  [  bank_id  ]  .  regs  [  " +
				"cast  (  addr  [  7  :  0  ]  ,  uint32  )  ]  .  " +
				"read  (  ctx  )   ;",
			expected: "log info: `entry: %#x`,\n" +
				"    dev.banks[bank_id].regs[cast(addr[7:0], uint32)].read(ctx);",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestLogStatement_RichExpressionsWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// 01: Function call inside format argument wrapping its arguments
		{
			name: "01-format-arg-call-arguments-wrap",
			code: "log info: `res: %d`, compute_value(very_long_argument_alpha, very_long_argument_beta);",
			expected: "log info: `res: %d`,\n" +
				"    compute_value(\n" +
				"        very_long_argument_alpha,\n" +
				"        very_long_argument_beta\n" +
				"    );",
			pc: func(f func()) { f() },
		},
		// 02: Binary arithmetic expression in format argument wrapping
		{
			name: "02-format-arg-binary-arithmetic-wrap",
			code: "log info: `sum: %d`, very_long_operand_alpha + very_long_operand_beta;",
			expected: "log info: `sum: %d`,\n" +
				"    very_long_operand_alpha\n" +
				"    + very_long_operand_beta;",
			pc: func(f func()) { f() },
		},
		// 03: Ternary expression in format argument wrapping
		{
			name: "03-format-arg-ternary-wrap",
			code: "log info: `status: %s`, " +
				"(is_device_ready ? LEVEL_SUCCESS_REPORT : LEVEL_FAILURE_REPORT);",
			expected: "log info: `status: %s`,\n" +
				"    (\n" +
				"        is_device_ready\n" +
				"            ? LEVEL_SUCCESS_REPORT\n" +
				"            : LEVEL_FAILURE_REPORT\n" +
				"    );",
			pc: func(f func()) { f() },
		},
		// 04: Cast expression in format argument wrapping
		{
			name: "04-format-arg-cast-wrap",
			code: "log info: `cast: %p`, " +
				"cast(very_long_source_pointer_expression, const void *);",
			expected: "log info: `cast: %p`,\n" +
				"    cast(\n" +
				"        very_long_source_pointer_expression,\n" +
				"        const void *\n" +
				"    );",
			pc: func(f func()) { f() },
		},
		// 05: Multiple rich format arguments wrapping after commas
		{
			name: "05-multiple-rich-args-wrapping",
			code: "log info: `dev=%s, addr=%#x, val=%#x`, this.get_device_name(), " +
				"base_address + offset_value, buffer[current_index];",
			expected: "log info: `dev=%s, addr=%#x, val=%#x`,\n" +
				"    this.get_device_name(),\n" +
				"    base_address + offset_value,\n" +
				"    buffer[current_index];",
			pc: func(f func()) { f() },
		},
		// 06: Rich level, rich group, and rich format args wrapping together
		{
			name: "06-rich-level-group-and-args-wrap",
			code: "log info, (dev.calculate_level(ctx)), (Register_Read | Register_Write): " +
				"`val: %#x`, this.read_current_value();",
			expected: "log info, (dev.calculate_level(ctx)), (\n" +
				"    Register_Read | Register_Write\n" +
				"): `val: %#x`,\n" +
				"    this.read_current_value();",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestLogStatement_IndentationLeakage(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-format-args-wrap-no-leak",
			code: "log info: `hello, %s`, `world`;\n" +
				"next_statement = keep_normal_indentation;",
			expected: "log info: `hello, %s`,\n" +
				"    `world`;\n" +
				"next_statement = keep_normal_indentation;",
			pc: func(f func()) { f() },
		},
		{
			name: "02-parens-format-args-wrap-no-leak",
			code: "log info: (`hello, %s`), `world`;\n" +
				"next_statement = keep_normal_indentation;",
			expected: "log info: (`hello, %s`),\n" +
				"    `world`;\n" +
				"next_statement = keep_normal_indentation;",
			pc: func(f func()) { f() },
		},
		{
			name: "03-parens-wrapped-format-no-args-no-leak",
			code: "log info: (`very long format string part one ` + `very long format string part two`);\n" +
				"next_statement = keep_normal_indentation;",
			expected: "log info: (\n" +
				"    `very long format string part one `\n" +
				"    + `very long format string part two`\n" +
				");\n" +
				"next_statement = keep_normal_indentation;",
			pc: func(f func()) { f() },
		},
		{
			name: "04-rich-format-args-call-wrap-no-leak",
			code: "log info: `res: %d`, compute_value(very_long_argument_alpha, very_long_argument_beta);\n" +
				"next_statement = keep_normal_indentation;",
			expected: "log info: `res: %d`,\n" +
				"    compute_value(\n" +
				"        very_long_argument_alpha,\n" +
				"        very_long_argument_beta\n" +
				"    );\n" +
				"next_statement = keep_normal_indentation;",
			pc: func(f func()) { f() },
		},
		{
			name: "05-rich-format-args-parenthesized-ternary-no-leak",
			code: "log info: `status: %s`, (is_device_ready ? LEVEL_SUCCESS_REPORT : LEVEL_FAILURE_REPORT);\n" +
				"next_statement = keep_normal_indentation;",
			expected: "log info: `status: %s`,\n" +
				"    (\n" +
				"        is_device_ready\n" +
				"            ? LEVEL_SUCCESS_REPORT\n" +
				"            : LEVEL_FAILURE_REPORT\n" +
				"    );\n" +
				"next_statement = keep_normal_indentation;",
			pc: func(f func()) { f() },
		},
		{
			name: "06-format-arg-bitwise-shift-no-leak",
			code: "log info: `byte: %#x`, (val & 0xff00) >> 8;\n" +
				"next_statement = keep_normal_indentation;",
			expected: "log info: `byte: %#x`,\n" +
				"    (val & 0xff00) >> 8;\n" +
				"next_statement = keep_normal_indentation;",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
