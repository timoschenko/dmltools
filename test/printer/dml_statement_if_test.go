// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestIfStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-if",
			code:     "if   (  cond  )   {\n  x  =  1;\n}",
			expected: "if (cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-if-else",
			code:     "if   (  cond  )   {\n  x  =  1;\n}   else   {\n  x  =  2;\n}",
			expected: "if (cond) {\n    x = 1;\n} else {\n    x = 2;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "03-if-else-if-chain",
			code: "if (c1) {\nx = 1;\n} else if (c2) {\nx = 2;\n} else {\nx = 3;\n}",
			expected: "if (c1) {\n    x = 1;\n} else if (c2) {\n" +
				"    x = 2;\n} else {\n    x = 3;\n}",
			pc: func(f func()) { f() },
		},
		{
			name:     "04-hash-if",
			code:     "#if   (  cond  )   {\n  x  =  1;\n}",
			expected: "#if (cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-hash-if-hash-else",
			code:     "#if (cond) {\nx = 1;\n} #else {\nx = 2;\n}",
			expected: "#if (cond) {\n    x = 1;\n} #else {\n    x = 2;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "06-hash-if-hash-else-hash-if-chain",
			code: "#if (c1) {\nx = 1;\n} #else #if (c2) {\nx = 2;\n} #else {\nx = 3;\n}",
			expected: "#if (c1) {\n    x = 1;\n} #else #if (c2) {\n" +
				"    x = 2;\n} #else {\n    x = 3;\n}",
			pc: func(f func()) { f() },
		},
		{
			name:     "07-unbraced-if",
			code:     "if   (  cond  )   x  =  1;",
			expected: "if (cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-unbraced-if-else",
			code:     "if (cond) x = 1; else x = 2;",
			expected: "if (cond) {\n    x = 1;\n} else {\n    x = 2;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-unbraced-hash-if-hash-else",
			code:     "#if (cond) x = 1; #else x = 2;",
			expected: "#if (cond) {\n    x = 1;\n} #else {\n    x = 2;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "10-multiline-condition-wrapping",
			code: "if (very_long_condition_variable_a && " +
				"very_long_condition_variable_b) {\n    x = 1;\n}",
			expected: "if (\n    very_long_condition_variable_a\n" +
				"    && very_long_condition_variable_b\n) {\n    x = 1;\n}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-nested-if",
			code: "if (c1) {\nif (c2) {\nx = 1;\n} else {\nx = 2;\n}\n}",
			expected: "if (c1) {\n    if (c2) {\n        x = 1;\n" +
				"    } else {\n        x = 2;\n    }\n}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestIfStatement_ConditionWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - entire statement and condition fit within lineCap
		{
			name:     "01-level1-simple-relational",
			code:     "if (count < limit) {\n  x = 1;\n}",
			expected: "if (count < limit) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-level1-simple-logical",
			code:     "if (ready && count < 10) {\n  x = 1;\n}",
			expected: "if (ready && count < 10) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-level1-ternary",
			code:     "if (flag ? a < 5 : b < 5) {\n  x = 1;\n}",
			expected: "if (flag ? a < 5 : b < 5) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-level1-method-call",
			code:     "if (is_valid(dev, id)) {\n  x = 1;\n}",
			expected: "if (is_valid(dev, id)) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "05-level1-else-if-simple-logical",
			code: "if (c1) {\n  x = 1;\n} else if (ready && count < 10) {\n  x = 2;\n}",
			expected: "if (c1) {\n" +
				"    x = 1;\n" +
				"} else if (ready && count < 10) {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name:     "06-level1-hash-if-defined",
			code:     "#if (defined(FEAT_A) && FEAT_A) {\n  x = 1;\n}",
			expected: "#if (defined(FEAT_A) && FEAT_A) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "07-level1-hash-if-const-ternary",
			code: "#if (flag #? a < 5 #: b < 5) {\n  x = 1;\n} #else {\n  x = 2;\n}",
			expected: "#if (flag #? a < 5 #: b < 5) {\n" +
				"    x = 1;\n" +
				"} #else {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "08-level1-hash-else-hash-if-simple",
			code: "#if (c1) {\n  x = 1;\n} #else #if (ready && count < 10) {\n  x = 2;\n}",
			expected: "#if (c1) {\n" +
				"    x = 1;\n" +
				"} #else #if (ready && count < 10) {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 2: Outer Parenthesized Wrap - statement exceeds lineCap,
		// but inner condition fits on one line
		{
			name: "10-level2-outer-wrap-single-relational",
			code: "if (status == STATUS_WAITING_FOR_DEVICE) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    status == STATUS_WAITING_FOR_DEVICE\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-outer-wrap-compound-logical",
			code: "if (flag_alpha && flag_beta && is_ready) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    flag_alpha && flag_beta && is_ready\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-nested-parens-both-compact",
			code: "if ((cond_a && cond_b) || (cond_c && cond_d)) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    (cond_a && cond_b)\n" +
				"    || (cond_c && cond_d)\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-bitwise-slice-compact",
			code: "if ((reg_data[15:0] & STATUS_MASK) != 0) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    (reg_data[15:0] & STATUS_MASK) != 0\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "14-level2-member-access-compact",
			code: "if (device.bus.subsystem.is_active_now()) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    device.bus.subsystem.is_active_now()\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "15-level2-else-if-outer-wrap",
			code: "if (c1) {\n  x = 1;\n} else if (status == STATUS_WAITING_FOR_DEV) {\n  x = 2;\n}",
			expected: "if (c1) {\n" +
				"    x = 1;\n" +
				"} else if (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "16-level2-hash-if-outer-wrap",
			code: "#if (status == STATUS_WAITING_FOR_DEVICE) {\n  x = 1;\n} #else {\n  x = 2;\n}",
			expected: "#if (\n" +
				"    status == STATUS_WAITING_FOR_DEVICE\n" +
				") {\n" +
				"    x = 1;\n" +
				"} #else {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "17-level2-hash-else-hash-if-outer-wrap",
			code: "#if (c1) {\n  x = 1;\n} #else #if (status == STATUS_WAITING_FOR_DEV) {\n  x = 2;\n}",
			expected: "#if (c1) {\n" +
				"    x = 1;\n" +
				"} #else #if (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Deep Expression Explosion - both statement and inner condition
		// exceed lineCap, so inner expressions explode into multi-line
		{
			name: "20-level3-logical-and-explosion",
			code: "if (very_long_condition_name_a && very_long_condition_name_b) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-multi-operand-chain",
			code: "if (flag_alpha && flag_beta && flag_gamma && flag_delta) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    flag_alpha\n" +
				"    && flag_beta\n" +
				"    && flag_gamma\n" +
				"    && flag_delta\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-asymmetric-nested-explosion",
			code: "if ((very_long_cond_a && very_long_cond_b) || (cond_c && cond_d)) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"    || (cond_c && cond_d)\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-symmetric-nested-explosion",
			code: "if ((very_long_cond_a && very_long_cond_b) || " +
				"(very_long_cond_c && very_long_cond_d)) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"    || (\n" +
				"        very_long_cond_c\n" +
				"        && very_long_cond_d\n" +
				"    )\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "24-level3-relational-arithmetic-explosion",
			code: "if (current_offset + transfer_bytes <= max_buffer_capacity) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    current_offset + transfer_bytes\n" +
				"    <= max_buffer_capacity\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "25-level3-deep-arithmetic-explosion",
			code: "if (very_long_current_offset + very_long_transfer_bytes <= max_buffer_capacity) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    very_long_current_offset\n" +
				"    + very_long_transfer_bytes\n" +
				"    <= max_buffer_capacity\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "26-level3-ternary-explosion",
			code: "if (use_fast_mode ? fast_queue.has_items() : slow_queue.has_items()) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    use_fast_mode\n" +
				"        ? fast_queue.has_items()\n" +
				"        : slow_queue.has_items()\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "27-level3-ternary-with-nested-exploded-cond",
			code: "if ((very_long_cond_a && very_long_cond_b) ? check_fast() : check_slow()) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"        ? check_fast()\n" +
				"        : check_slow()\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "28-level3-func-call-args-explosion",
			code: "if (validate_transfer(source_addr, destination_addr, transfer_len)) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    validate_transfer(\n" +
				"        source_addr,\n" +
				"        destination_addr,\n" +
				"        transfer_len\n" +
				"    )\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "29-level3-member-chain-explosion",
			code: "if (bus.controller.subsystem.peripheral.is_ready()) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    bus.controller.subsystem.peripheral\n" +
				"    .is_ready()\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "30-level3-bitwise-slice-explosion",
			code: "if ((buffer_data[31:16] & STATUS_MASK) != EXPECTED_PATTERN) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    (buffer_data[31:16] & STATUS_MASK)\n" +
				"    != EXPECTED_PATTERN\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "31-level3-unary-negation-explosion",
			code: "if (!(device_is_initialized && is_transfer_complete)) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    !(\n" +
				"        device_is_initialized\n" +
				"        && is_transfer_complete\n" +
				"    )\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "32-level3-else-if-logical-explosion",
			code: "if (very_long_condition_name_a && very_long_condition_name_b) {\n" +
				"  x = 1;\n" +
				"} else if (very_long_condition_name_c && very_long_condition_name_d) {\n" +
				"  x = 2;\n" +
				"} else {\n" +
				"  x = 3;\n" +
				"}",
			expected: "if (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"} else if (\n" +
				"    very_long_condition_name_c\n" +
				"    && very_long_condition_name_d\n" +
				") {\n" +
				"    x = 2;\n" +
				"} else {\n" +
				"    x = 3;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "33-level3-else-if-asymmetric-nested-explosion",
			code: "if (c1) {\n" +
				"  x = 1;\n" +
				"} else if ((very_long_cond_a && very_long_cond_b) || (cond_c && cond_d)) {\n" +
				"  x = 2;\n" +
				"} else {\n" +
				"  x = 3;\n" +
				"}",
			expected: "if (c1) {\n" +
				"    x = 1;\n" +
				"} else if (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"    || (cond_c && cond_d)\n" +
				") {\n" +
				"    x = 2;\n" +
				"} else {\n" +
				"    x = 3;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "34-level3-hash-if-logical-explosion",
			code: "#if (very_long_condition_name_a && very_long_condition_name_b) {\n" +
				"  x = 1;\n" +
				"} #else #if (very_long_condition_name_c && very_long_condition_name_d) {\n" +
				"  x = 2;\n" +
				"} #else {\n" +
				"  x = 3;\n" +
				"}",
			expected: "#if (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"} #else #if (\n" +
				"    very_long_condition_name_c\n" +
				"    && very_long_condition_name_d\n" +
				") {\n" +
				"    x = 2;\n" +
				"} #else {\n" +
				"    x = 3;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "35-level3-hash-if-const-ternary-explosion",
			code: "#if (use_fast_mode #? fast_queue_is_enabled #: slow_queue_is_enabled) {\n" +
				"  x = 1;\n" +
				"} #else {\n" +
				"  x = 2;\n" +
				"}",
			expected: "#if (\n" +
				"    use_fast_mode\n" +
				"        #? fast_queue_is_enabled\n" +
				"        #: slow_queue_is_enabled\n" +
				") {\n" +
				"    x = 1;\n" +
				"} #else {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 4: Deep Cascading Hierarchy - 3 levels of nesting
		{
			name: "40-level4-three-level-nesting",
			code: "if (((deep_operand_a + deep_operand_b) > threshold && deep_cond_c) || " +
				"fallback_cond) {\n  x = 1;\n}",
			expected: "if (\n" +
				"    (\n" +
				"        (\n" +
				"            deep_operand_a\n" +
				"            + deep_operand_b\n" +
				"        ) > threshold\n" +
				"        && deep_cond_c\n" +
				"    )\n" +
				"    || fallback_cond\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			masked: MaskBinaryOperatorWrappingAfterMultilineParen(
				"if (\n" +
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
					"    x = 1;\n" +
					"}",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "41-level4-else-if-three-level-nesting",
			code: "if (c1) {\n  x = 1;\n} else if (((deep_operand_a + deep_operand_b) > " +
				"threshold && deep_cond_c) || fallback_cond) {\n  x = 2;\n}",
			expected: "if (c1) {\n" +
				"    x = 1;\n" +
				"} else if (\n" +
				"    (\n" +
				"        (\n" +
				"            deep_operand_a\n" +
				"            + deep_operand_b\n" +
				"        ) > threshold\n" +
				"        && deep_cond_c\n" +
				"    )\n" +
				"    || fallback_cond\n" +
				") {\n" +
				"    x = 2;\n" +
				"}",
			masked: MaskBinaryOperatorWrappingAfterMultilineParen(
				"if (c1) {\n" +
					"    x = 1;\n" +
					"} else if (\n" +
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
					"    x = 2;\n" +
					"}",
			),
			pc: func(f func()) { f() },
		},

		// Unbraced If / Else / #if / #else with Wrapping Conditions
		{
			name: "50-unbraced-level2-outer-wrap",
			code: "if (status == STATUS_WAITING_FOR_DEVICE)\n  x = 1;\nelse\n  x = 2;",
			expected: "if (\n" +
				"    status == STATUS_WAITING_FOR_DEVICE\n" +
				") {\n" +
				"    x = 1;\n" +
				"} else {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "51-unbraced-level3-exploded",
			code: "if (very_long_condition_name_a && very_long_condition_name_b)\n" +
				"  x = 1;\n" +
				"else if (very_long_condition_name_c && very_long_condition_name_d)\n" +
				"  x = 2;\n" +
				"else\n" +
				"  x = 3;",
			expected: "if (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"} else if (\n" +
				"    very_long_condition_name_c\n" +
				"    && very_long_condition_name_d\n" +
				") {\n" +
				"    x = 2;\n" +
				"} else {\n" +
				"    x = 3;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "52-unbraced-level3-nested",
			code: "if ((very_long_cond_a && very_long_cond_b) || (cond_c && cond_d))\n  x = 1;",
			expected: "if (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"    || (cond_c && cond_d)\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "53-unbraced-hash-if-level3-exploded",
			code: "#if (very_long_condition_name_a && very_long_condition_name_b)\n" +
				"  x = 1;\n" +
				"#else\n" +
				"  x = 2;",
			expected: "#if (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"} #else {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestIfStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-leading-comment",
			code:     "// comment\nif (cond) {\nx = 1;\n}",
			expected: "// comment\nif (cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-between-if-and-paren",
			code:     "if /* c */ (cond) {\nx = 1;\n}",
			expected: "if /* c */ (cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-inside-condition",
			code:     "if (/* c */ cond) {\nx = 1;\n}",
			expected: "if (/* c */ cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-between-paren-and-brace",
			code:     "if (cond) /* c */ {\nx = 1;\n}",
			expected: "if (cond) /* c */ {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "05-comment-between-brace-and-else",
			code: "if (cond) {\nx = 1;\n} /* c */ else {\nx = 2;\n}",
			expected: "if (cond) {\n    x = 1;\n} /* c */ else {\n" +
				"    x = 2;\n}",
			pc: func(f func()) { f() },
		},
		{
			name: "06-comment-between-else-and-brace",
			code: "if (cond) {\nx = 1;\n} else /* c */ {\nx = 2;\n}",
			expected: "if (cond) {\n    x = 1;\n} else /* c */ {\n" +
				"    x = 2;\n}",
			pc: func(f func()) { f() },
		},
		{
			name: "07-comment-between-else-and-if",
			code: "if (c1) {\nx = 1;\n} else /* c */ if (c2) {\nx = 2;\n}",
			expected: "if (c1) {\n    x = 1;\n} else /* c */ if (c2) {\n" +
				"    x = 2;\n}",
			pc: func(f func()) { f() },
		},
		{
			name: "08-comment-hash-if-hash-else",
			code: "#if /* c1 */ (cond) {\nx = 1;\n} /* c2 */ #else /* c3 */ {\nx = 2;\n}",
			expected: "#if /* c1 */ (cond) {\n    x = 1;\n} /* c2 */ #else /* c3 */ {\n" +
				"    x = 2;\n}",
			pc: func(f func()) { f() },
		},
		{
			name:     "09-trailing-comment",
			code:     "if (cond) {\nx = 1;\n} // done",
			expected: "if (cond) {\n    x = 1;\n} // done",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
