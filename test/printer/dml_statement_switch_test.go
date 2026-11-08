// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestSwitchStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 60
	testCases := []printerTestCase{
		{
			name:     "01-empty-switch",
			code:     "switch (x) {}",
			expected: "switch (x) {\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "02-single-case-break",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"break;\n" +
				"}",
			expected: "switch (x) {\n" +
				"case 1:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "03-case-with-default",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"call();\n" +
				"break;\n" +
				"default:\n" +
				"other();\n" +
				"break;\n" +
				"}",
			expected: "switch (x) {\n" +
				"case 1:\n" +
				"    call();\n" +
				"    break;\n" +
				"default:\n" +
				"    other();\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "04-fallthrough-cases",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"case 2:\n" +
				"call();\n" +
				"break;\n" +
				"case 3:\n" +
				"default:\n" +
				"other();\n" +
				"}",
			expected: "switch (x) {\n" +
				"case 1:\n" +
				"case 2:\n" +
				"    call();\n" +
				"    break;\n" +
				"case 3:\n" +
				"default:\n" +
				"    other();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "05-multiple-statements-in-case",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"x = 1;\n" +
				"y = 2;\n" +
				"break;\n" +
				"default:\n" +
				"x = 0;\n" +
				"}",
			expected: "switch (x) {\n" +
				"case 1:\n" +
				"    x = 1;\n" +
				"    y = 2;\n" +
				"    break;\n" +
				"default:\n" +
				"    x = 0;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "06-compound-statement-in-case",
			code: "switch (x) {\n" +
				"case 1: {\n" +
				"local int val;\n" +
				"return;\n" +
				"}\n" +
				"default:\n" +
				"break;\n" +
				"}",
			expected: "switch (x) {\n" +
				"case 1:\n" +
				"    {\n" +
				"        local int val;\n" +
				"        return;\n" +
				"    }\n" +
				"default:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "07-nested-switch",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"switch (y) {\n" +
				"case 2:\n" +
				"break;\n" +
				"}\n" +
				"break;\n" +
				"default:\n" +
				"break;\n" +
				"}",
			expected: "switch (x) {\n" +
				"case 1:\n" +
				"    switch (y) {\n" +
				"    case 2:\n" +
				"        break;\n" +
				"    }\n" +
				"    break;\n" +
				"default:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "08-hash-if-in-switch",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"#if (cond) {\n" +
				"case 2:\n" +
				"call();\n" +
				"break;\n" +
				"}\n" +
				"default:\n" +
				"break;\n" +
				"}",
			expected: "switch (x) {\n" +
				"case 1:\n" +
				"#if (cond) {\n" +
				"case 2:\n" +
				"    call();\n" +
				"    break;\n" +
				"}\n" +
				"default:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "08a-hash-if-else-in-switch",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"#if (cond) {\n" +
				"case 2:\n" +
				"call();\n" +
				"break;\n" +
				"} #else {\n" +
				"case 3:\n" +
				"other();\n" +
				"break;\n" +
				"}\n" +
				"default:\n" +
				"break;\n" +
				"}",
			expected: "switch (x) {\n" +
				"case 1:\n" +
				"#if (cond) {\n" +
				"case 2:\n" +
				"    call();\n" +
				"    break;\n" +
				"} #else {\n" +
				"case 3:\n" +
				"    other();\n" +
				"    break;\n" +
				"}\n" +
				"default:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "09-whitespace-normalization",
			code: "switch   (  val  )   {\n" +
				"  case   0  :   \n" +
				"  handle_zero();\n" +
				"  break;\n" +
				"  default  :   \n" +
				"  handle_default();\n" +
				"  break;\n" +
				"}",
			expected: "switch (val) {\n" +
				"case 0:\n" +
				"    handle_zero();\n" +
				"    break;\n" +
				"default:\n" +
				"    handle_default();\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name:    "10-line-wrapping-expression",
			lineCap: 40,
			code: "switch (very_long_variable_name_a + very_long_variable_name_b) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			expected: "switch (\n" +
				"    very_long_variable_name_a\n" +
				"    + very_long_variable_name_b\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-break-after-hash-if-in-switch",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"test();\n" +
				"break;\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"test2();\n" +
				"}\n" +
				"break;\n" +
				"default:\n" +
				"log error, 1, log_group_ident: `message %x`, foo;\n" +
				"}",
			expected: "switch (foo) {\n" +
				"case A:\n" +
				"    test();\n" +
				"    break;\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"    test2();\n" +
				"}\n" +
				"    break;\n" +
				"default:\n" +
				"    log error, 1, log_group_ident: `message %x`,\n" +
				"        foo;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-case-labels-inside-hash-if",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"test();\n" +
				"break;\n" +
				"case B:\n" +
				"#if (FEAT) {\n" +
				"case CA:\n" +
				"}\n" +
				"#if (FEAT_2) {\n" +
				"case CB:\n" +
				"}\n" +
				"val = test2();\n" +
				"break;\n" +
				"default:\n" +
				"log error, 1, log_group_ident: `message %x`, foo;\n" +
				"}",
			expected: "switch (foo) {\n" +
				"case A:\n" +
				"    test();\n" +
				"    break;\n" +
				"case B:\n" +
				"#if (FEAT) {\n" +
				"case CA:\n" +
				"}\n" +
				"#if (FEAT_2) {\n" +
				"case CB:\n" +
				"}\n" +
				"    val = test2();\n" +
				"    break;\n" +
				"default:\n" +
				"    log error, 1, log_group_ident: `message %x`,\n" +
				"        foo;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-statements-between-hash-if-blocks",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"test();\n" +
				"}\n" +
				"val = test2();\n" +
				"break;\n" +
				"}",
			expected: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"    test();\n" +
				"}\n" +
				"    val = test2();\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "14-hash-if-case-fallthrough",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"}\n" +
				"common_code();\n" +
				"break;\n" +
				"}",
			expected: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"}\n" +
				"    common_code();\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "15-break-after-hash-if-else-in-switch",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"test();\n" +
				"break;\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"test2();\n" +
				"} #else {\n" +
				"case C:\n" +
				"test3();\n" +
				"}\n" +
				"break;\n" +
				"default:\n" +
				"log error, 1, log_group_ident: `message %x`, foo;\n" +
				"}",
			expected: "switch (foo) {\n" +
				"case A:\n" +
				"    test();\n" +
				"    break;\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"    test2();\n" +
				"} #else {\n" +
				"case C:\n" +
				"    test3();\n" +
				"}\n" +
				"    break;\n" +
				"default:\n" +
				"    log error, 1, log_group_ident: `message %x`,\n" +
				"        foo;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "16-case-labels-inside-hash-if-else",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"test();\n" +
				"break;\n" +
				"case B:\n" +
				"#if (FEAT) {\n" +
				"case CA:\n" +
				"} #else {\n" +
				"case CB:\n" +
				"}\n" +
				"val = test2();\n" +
				"break;\n" +
				"default:\n" +
				"log error, 1, log_group_ident: `message %x`, foo;\n" +
				"}",
			expected: "switch (foo) {\n" +
				"case A:\n" +
				"    test();\n" +
				"    break;\n" +
				"case B:\n" +
				"#if (FEAT) {\n" +
				"case CA:\n" +
				"} #else {\n" +
				"case CB:\n" +
				"}\n" +
				"    val = test2();\n" +
				"    break;\n" +
				"default:\n" +
				"    log error, 1, log_group_ident: `message %x`,\n" +
				"        foo;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "17-statements-between-hash-if-else-blocks",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"test();\n" +
				"} #else {\n" +
				"case C:\n" +
				"test2();\n" +
				"}\n" +
				"val = test3();\n" +
				"break;\n" +
				"}",
			expected: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"    test();\n" +
				"} #else {\n" +
				"case C:\n" +
				"    test2();\n" +
				"}\n" +
				"    val = test3();\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "18-hash-if-else-case-fallthrough",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"} #else {\n" +
				"case C:\n" +
				"}\n" +
				"common_code();\n" +
				"break;\n" +
				"}",
			expected: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"} #else {\n" +
				"case C:\n" +
				"}\n" +
				"    common_code();\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestSwitchStatement_ConditionWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - entire statement and condition fit within lineCap
		{
			name: "01-level1-simple-variable",
			code: "switch (val) {\ncase 0:\n  break;\n}",
			expected: "switch (val) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "02-level1-simple-arithmetic",
			code: "switch (count + offset) {\ncase 0:\n  break;\n}",
			expected: "switch (count + offset) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "03-level1-simple-bitwise",
			code: "switch (status & MASK) {\ncase 0:\n  break;\n}",
			expected: "switch (status & MASK) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "04-level1-simple-relational",
			code: "switch (count < limit) {\ncase 0:\n  break;\n}",
			expected: "switch (count < limit) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "05-level1-simple-ternary",
			code: "switch (flag ? a : b) {\ncase 0:\n  break;\n}",
			expected: "switch (flag ? a : b) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "06-level1-simple-method-call",
			code: "switch (get_val(dev, id)) {\ncase 0:\n  break;\n}",
			expected: "switch (get_val(dev, id)) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "07-level1-simple-bit-slice",
			code: "switch (data[3:0]) {\ncase 0:\n  break;\n}",
			expected: "switch (data[3:0]) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "08-level1-simple-cast",
			code: "switch (cast(val, uint32)) {\ncase 0:\n  break;\n}",
			expected: "switch (cast(val, uint32)) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 2: Outer Parenthesized Wrap - statement exceeds lineCap,
		// but inner condition fits on one line
		{
			name: "10-level2-outer-wrap-single-variable",
			code: "switch (status_code_from_device_register) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    status_code_from_device_register\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-outer-wrap-compound-arithmetic",
			code: "switch (base_offset + index_offset + 4) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    base_offset + index_offset + 4\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-outer-wrap-compound-logical",
			code: "switch (flag_alpha && flag_beta && ready) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    flag_alpha && flag_beta && ready\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-nested-parens-both-compact",
			code: "switch ((cond_a && cond_b) || (cond_c && cond_d)) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    (cond_a && cond_b)\n" +
				"    || (cond_c && cond_d)\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "14-level2-bitwise-slice-compact",
			code: "switch ((data[15:0] & STATUS_MASK) >> 4) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    (data[15:0] & STATUS_MASK) >> 4\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "15-level2-member-access-compact",
			code: "switch (device.subsystem.get_state_now()) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    device.subsystem.get_state_now()\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "16-level2-cast-compact",
			code: "switch (cast(buffer_data[15:0], uint32)) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    cast(buffer_data[15:0], uint32)\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Deep Expression Explosion - both statement and inner condition
		// exceed lineCap, so inner expressions explode into multi-line
		{
			name: "20-level3-logical-and-explosion",
			code: "switch (very_long_condition_name_a && very_long_condition_name_b) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-multi-operand-chain",
			code: "switch (flag_alpha && flag_beta && flag_gamma && flag_delta) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    flag_alpha\n" +
				"    && flag_beta\n" +
				"    && flag_gamma\n" +
				"    && flag_delta\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-asymmetric-nested-explosion",
			code: "switch ((very_long_cond_a && very_long_cond_b) || (cond_c && cond_d)) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"    || (cond_c && cond_d)\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-symmetric-nested-explosion",
			code: "switch ((very_long_cond_a && very_long_cond_b) || " +
				"(very_long_cond_c && very_long_cond_d)) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"    || (\n" +
				"        very_long_cond_c\n" +
				"        && very_long_cond_d\n" +
				"    )\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "24-level3-relational-arithmetic-explosion",
			code: "switch (current_offset + transfer_bytes <= max_buffer_capacity) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    current_offset + transfer_bytes\n" +
				"    <= max_buffer_capacity\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "25-level3-deep-arithmetic-explosion",
			code: "switch (very_long_current_offset + very_long_transfer_bytes <= max_buffer_capacity) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    very_long_current_offset\n" +
				"    + very_long_transfer_bytes\n" +
				"    <= max_buffer_capacity\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "26-level3-ternary-explosion",
			code: "switch (use_fast_mode ? fast_queue.has_items() : slow_queue.has_items()) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    use_fast_mode\n" +
				"        ? fast_queue.has_items()\n" +
				"        : slow_queue.has_items()\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "27-level3-ternary-with-nested-exploded-cond",
			code: "switch ((very_long_cond_a && very_long_cond_b) ? check_fast() : check_slow()) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"        ? check_fast()\n" +
				"        : check_slow()\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "28-level3-func-call-args-explosion",
			code: "switch (validate_transfer(source_addr, destination_addr, transfer_len)) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    validate_transfer(\n" +
				"        source_addr,\n" +
				"        destination_addr,\n" +
				"        transfer_len\n" +
				"    )\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "29-level3-member-chain-explosion",
			code: "switch (bus.controller.subsystem.peripheral.is_ready()) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    bus.controller.subsystem.peripheral\n" +
				"    .is_ready()\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "30-level3-bitwise-slice-explosion",
			code: "switch ((buffer_data[31:16] & STATUS_MASK) != EXPECTED_PATTERN) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    (buffer_data[31:16] & STATUS_MASK)\n" +
				"    != EXPECTED_PATTERN\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "31-level3-unary-negation-explosion",
			code: "switch (!(device_is_initialized && is_transfer_complete)) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    !(\n" +
				"        device_is_initialized\n" +
				"        && is_transfer_complete\n" +
				"    )\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 4: Deep Cascading Hierarchy - 3 levels of nesting
		{
			name: "40-level4-three-level-nesting",
			code: "switch (((deep_operand_a + deep_operand_b) > threshold && deep_cond_c) || " +
				"fallback_cond) {\ncase 0:\n  break;\n}",
			expected: "switch (\n" +
				"    (\n" +
				"        (\n" +
				"            deep_operand_a\n" +
				"            + deep_operand_b\n" +
				"        ) > threshold\n" +
				"        && deep_cond_c\n" +
				"    )\n" +
				"    || fallback_cond\n" +
				") {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			masked: MaskBinaryOperatorWrappingAfterMultilineParen(
				"switch (\n" +
					"    (\n" +
					"        (\n" +
					"            deep_operand_a\n" +
					"            + deep_operand_b\n" +
					"        )\n" +
					"        > threshold\n" +
					"        && deep_cond_c\n" +
					"    )\n" +
					"    || fallback_cond\n" +
					") {\n" +
					"case 0:\n" +
					"    break;\n" +
					"}",
			),
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestSwitchStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name: "01-leading-comment",
			code: "// c\n" +
				"switch (x) {\n" +
				"case 1:\n" +
				"break;\n" +
				"}",
			expected: "// c\n" +
				"switch (x) {\n" +
				"case 1:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "02-comment-between-switch-and-paren",
			code: "switch /* c */ (val) {\n" +
				"case 0:\n" +
				"break;\n" +
				"}",
			expected: "switch /* c */ (val) {\n" +
				"case 0:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "03-comment-after-switch-lparen",
			code: "switch (/* c */ x) {\n" +
				"case 1:\n" +
				"break;\n" +
				"}",
			expected: "switch (/* c */ x) {\n" +
				"case 1:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "04-comment-before-rparen",
			code: "switch (x /* c */) {\n" +
				"case 1:\n" +
				"break;\n" +
				"}",
			expected: "switch (x /* c */) {\n" +
				"case 1:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "04a-comment-between-paren-and-brace",
			code: "switch (x) /* c */ {\n" +
				"case 1:\n" +
				"break;\n" +
				"}",
			expected: "switch (x) /* c */ {\n" +
				"case 1:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "04b-comment-between-operands-in-condition",
			code: "switch (a /* c */ + b) {\n" +
				"case 1:\n" +
				"break;\n" +
				"}",
			expected: "switch (a /* c */ + b) {\n" +
				"case 1:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "05-comment-between-case-and-expr",
			code: "switch (x) {\n" +
				"case /* c */ 1:\n" +
				"break;\n" +
				"}",
			expected: "switch (x) {\n" +
				"case /* c */ 1:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "06-comment-before-colon",
			code: "switch (val) {\n" +
				"case 0 /* c */ :\n" +
				"break;\n" +
				"}",
			expected: "switch (val) {\n" +
				"case 0 /* c */:\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "07-comment-after-colon",
			code: "switch (val) {\n" +
				"case 0: /* c */\n" +
				"break;\n" +
				"}",
			expected: "switch (val) {\n" +
				"case 0:\n" +
				"    /* c */\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "08-comment-on-colon-line",
			code: "switch (x) {\n" +
				"case 1: // c\n" +
				"break;\n" +
				"}",
			expected: "switch (x) {\n" +
				"case 1: // c\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "09-trailing-comment",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"break;\n" +
				"} // done",
			expected: "switch (x) {\n" +
				"case 1:\n" +
				"    break;\n" +
				"} // done",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
