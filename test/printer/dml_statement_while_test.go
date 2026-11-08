// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestWhileLoopStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-while",
			code:     "while   (  cond  )   {\n  x  =  1;\n}",
			expected: "while (cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-unbraced-while",
			code:     "while   (  cond  )   x  =  1;",
			expected: "while (cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "03-line-wrapping-condition",
			code: "while (very_long_condition_name_a && very_long_condition_name_b) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "04-nested-while",
			code: "while   (  c1  )   {\n  while   (  c2  )   {\n    x  =  1;\n  }\n}",
			expected: "while (c1) {\n" +
				"    while (c2) {\n" +
				"        x = 1;\n" +
				"    }\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestWhileLoopStatement_ConditionWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - entire statement and condition fit within lineCap
		{
			name:     "01-level1-simple-relational",
			code:     "while (count < limit) {\n  x = 1;\n}",
			expected: "while (count < limit) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-level1-simple-logical",
			code:     "while (ready && count < 10) {\n  x = 1;\n}",
			expected: "while (ready && count < 10) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-level1-ternary",
			code:     "while (flag ? a < 5 : b < 5) {\n  x = 1;\n}",
			expected: "while (flag ? a < 5 : b < 5) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-level1-method-call",
			code:     "while (is_valid(dev, id)) {\n  x = 1;\n}",
			expected: "while (is_valid(dev, id)) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},

		// Level 2: Outer Parenthesized Wrap - statement exceeds lineCap,
		// but inner condition fits on one line
		{
			name: "10-level2-outer-wrap-single-relational",
			code: "while (status == STATUS_WAITING_FOR_DEV) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-outer-wrap-compound-logical",
			code: "while (flag_alpha && flag_beta && ready) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    flag_alpha && flag_beta && ready\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-nested-parens-both-compact",
			code: "while ((cond_a && cond_b) || (cond_c && cond_d)) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    (cond_a && cond_b)\n" +
				"    || (cond_c && cond_d)\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-bitwise-slice-compact",
			code: "while ((data[15:0] & STATUS_MASK) != 0) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    (data[15:0] & STATUS_MASK) != 0\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "14-level2-member-access-compact",
			code: "while (device.subsystem.is_active_now()) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    device.subsystem.is_active_now()\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Deep Expression Explosion - both statement and inner condition
		// exceed lineCap, so inner expressions explode into multi-line
		{
			name: "20-level3-logical-and-explosion",
			code: "while (very_long_condition_name_a && very_long_condition_name_b) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-multi-operand-chain",
			code: "while (flag_alpha && flag_beta && flag_gamma && flag_delta) {\n  x = 1;\n}",
			expected: "while (\n" +
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
			code: "while ((very_long_cond_a && very_long_cond_b) || (cond_c && cond_d)) {\n  x = 1;\n}",
			expected: "while (\n" +
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
			code: "while ((very_long_cond_a && very_long_cond_b) || " +
				"(very_long_cond_c && very_long_cond_d)) {\n  x = 1;\n}",
			expected: "while (\n" +
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
			code: "while (current_offset + transfer_bytes <= max_buffer_capacity) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    current_offset + transfer_bytes\n" +
				"    <= max_buffer_capacity\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "25-level3-deep-arithmetic-explosion",
			code: "while (very_long_current_offset + very_long_transfer_bytes <= max_buffer_capacity) {\n  x = 1;\n}",
			expected: "while (\n" +
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
			code: "while (use_fast_mode ? fast_queue.has_items() : slow_queue.has_items()) {\n  x = 1;\n}",
			expected: "while (\n" +
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
			code: "while ((very_long_cond_a && very_long_cond_b) ? check_fast() : check_slow()) {\n  x = 1;\n}",
			expected: "while (\n" +
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
			code: "while (validate_transfer(source_addr, destination_addr, transfer_len)) {\n  x = 1;\n}",
			expected: "while (\n" +
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
			code: "while (bus.controller.subsystem.peripheral.is_ready()) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    bus.controller.subsystem.peripheral\n" +
				"    .is_ready()\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "30-level3-bitwise-slice-explosion",
			code: "while ((buffer_data[31:16] & STATUS_MASK) != EXPECTED_PATTERN) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    (buffer_data[31:16] & STATUS_MASK)\n" +
				"    != EXPECTED_PATTERN\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "31-level3-unary-negation-explosion",
			code: "while (!(device_is_initialized && is_transfer_complete)) {\n  x = 1;\n}",
			expected: "while (\n" +
				"    !(\n" +
				"        device_is_initialized\n" +
				"        && is_transfer_complete\n" +
				"    )\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 4: Deep Cascading Hierarchy - 3 levels of nesting
		{
			name: "40-level4-three-level-nesting",
			code: "while (((deep_operand_a + deep_operand_b) > threshold && deep_cond_c) || " +
				"fallback_cond) {\n  x = 1;\n}",
			expected: "while (\n" +
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
				"while (\n" +
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

		// Unbraced While Loops with Wrapping Conditions
		{
			name: "50-unbraced-level2-outer-wrap",
			code: "while (status == STATUS_WAITING_FOR_DEV)\n  x = 1;",
			expected: "while (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "51-unbraced-level3-exploded",
			code: "while (very_long_condition_name_a && very_long_condition_name_b)\n  x = 1;",
			expected: "while (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "52-unbraced-level3-nested",
			code: "while ((very_long_cond_a && very_long_cond_b) || (cond_c && cond_d))\n  x = 1;",
			expected: "while (\n" +
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
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestWhileLoopStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-leading-comment",
			code:     "// c\nwhile (cond) {\nx = 1;\n}",
			expected: "// c\nwhile (cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-between-while-and-paren",
			code:     "while /* c */ (cond) {\nx = 1;\n}",
			expected: "while /* c */ (cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-inside-condition",
			code:     "while (/* c */ cond) {\nx = 1;\n}",
			expected: "while (/* c */ cond) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-between-paren-and-brace",
			code:     "while (cond) /* c */ {\nx = 1;\n}",
			expected: "while (cond) /* c */ {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-trailing-comment",
			code:     "while (cond) {\nx = 1;\n} // done",
			expected: "while (cond) {\n    x = 1;\n} // done",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestDoWhileLoopStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-do-while",
			code:     "do   {\n  x  =  1;\n}   while   (  cond  )   ;",
			expected: "do {\n    x = 1;\n} while (cond);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-unbraced-do-while",
			code:     "do   x  =  1;   while   (  cond  )   ;",
			expected: "do {\n    x = 1;\n} while (cond);",
			pc:       func(f func()) { f() },
		},
		{
			name: "03-line-wrapping-condition",
			code: "do {\n  x = 1;\n} while (very_long_condition_name_a && very_long_condition_name_b);",
			expected: "do {\n" +
				"    x = 1;\n" +
				"} while (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				");",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestDoWhileLoopStatement_ConditionWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - entire statement and condition fit within lineCap
		{
			name: "01-level1-simple-logical",
			code: "do {\n  x = 1;\n} while (ready && count < 10);",
			expected: "do {\n" +
				"    x = 1;\n" +
				"} while (ready && count < 10);",
			pc: func(f func()) { f() },
		},

		// Level 2: Outer Parenthesized Wrap - statement exceeds lineCap,
		// but inner condition fits on one line
		{
			name: "10-level2-outer-wrap-single-relational",
			code: "do {\n  x = 1;\n} while (status == STATUS_WAITING_FOR_DEV);",
			expected: "do {\n" +
				"    x = 1;\n" +
				"} while (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-nested-parens-both-compact",
			code: "do {\n  x = 1;\n} while ((cond_a && cond_b) || (cond_c && cond_d));",
			expected: "do {\n" +
				"    x = 1;\n" +
				"} while (\n" +
				"    (cond_a && cond_b)\n" +
				"    || (cond_c && cond_d)\n" +
				");",
			pc: func(f func()) { f() },
		},

		// Level 3: Deep Expression Explosion - both statement and inner condition
		// exceed lineCap, so inner expressions explode into multi-line
		{
			name: "20-level3-logical-and-explosion",
			code: "do {\n  x = 1;\n} while (very_long_condition_name_a && very_long_condition_name_b);",
			expected: "do {\n" +
				"    x = 1;\n" +
				"} while (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-asymmetric-nested-explosion",
			code: "do {\n  x = 1;\n} while ((very_long_cond_a && very_long_cond_b) || " +
				"(cond_c && cond_d));",
			expected: "do {\n" +
				"    x = 1;\n" +
				"} while (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"    || (cond_c && cond_d)\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-symmetric-nested-explosion",
			code: "do {\n  x = 1;\n} while ((very_long_cond_a && very_long_cond_b) || " +
				"(very_long_cond_c && very_long_cond_d));",
			expected: "do {\n" +
				"    x = 1;\n" +
				"} while (\n" +
				"    (\n" +
				"        very_long_cond_a\n" +
				"        && very_long_cond_b\n" +
				"    )\n" +
				"    || (\n" +
				"        very_long_cond_c\n" +
				"        && very_long_cond_d\n" +
				"    )\n" +
				");",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestDoWhileLoopStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-leading-comment",
			code:     "// c\ndo {\nx = 1;\n} while (cond);",
			expected: "// c\ndo {\n    x = 1;\n} while (cond);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-between-do-and-brace",
			code:     "do /* c */ {\nx = 1;\n} while (cond);",
			expected: "do /* c */ {\n    x = 1;\n} while (cond);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-between-brace-and-while",
			code:     "do {\nx = 1;\n} /* c */ while (cond);",
			expected: "do {\n    x = 1;\n} /* c */ while (cond);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-after-while",
			code:     "do {\nx = 1;\n} while /* c */ (cond);",
			expected: "do {\n    x = 1;\n} while /* c */ (cond);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-trailing-comment-after-semi",
			code:     "do {\nx = 1;\n} while (cond); // done",
			expected: "do {\n    x = 1;\n} while (cond); // done",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
