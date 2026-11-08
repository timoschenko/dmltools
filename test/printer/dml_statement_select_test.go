// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestSelectStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 80
	testCases := []printerTestCase{
		{
			name: "01-basic-select",
			code: "select   item   in   (  items  )   where   (  item  .  valid  )   {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "02-select-with-else",
			code: "select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"} else {\n" +
				"    fallback();\n" +
				"}",
			expected: "select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"} else {\n" +
				"    fallback();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "03-select-else-select-chain",
			code: "select item in (fast_items) where (item.valid) {\n" +
				"    process(item);\n" +
				"} else select item in (slow_items) where (item.valid) {\n" +
				"    process(item);\n" +
				"} else {\n" +
				"    fallback();\n" +
				"}",
			expected: "select item in (fast_items) where (item.valid) {\n" +
				"    process(item);\n" +
				"} else select item in (slow_items) where (item.valid) {\n" +
				"    process(item);\n" +
				"} else {\n" +
				"    fallback();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "04-hash-select",
			code: "#select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			expected: "#select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "05-hash-select-with-hash-else",
			code: "#select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"} #else {\n" +
				"    fallback();\n" +
				"}",
			expected: "#select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"} #else {\n" +
				"    fallback();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "06-hash-select-hash-else-hash-select-chain",
			code: "#select item in (items_a) where (item.valid) {\n" +
				"    process(item);\n" +
				"} #else #select item in (items_b) where (item.valid) {\n" +
				"    process(item);\n" +
				"} #else {\n" +
				"    fallback();\n" +
				"}",
			expected: "#select item in (items_a) where (item.valid) {\n" +
				"    process(item);\n" +
				"} #else #select item in (items_b) where (item.valid) {\n" +
				"    process(item);\n" +
				"} #else {\n" +
				"    fallback();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "07-unbraced-select",
			code: "select item in (items) where (item.valid) process(item);",
			expected: "select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "08-unbraced-hash-select",
			code: "#select item in (items) where (item.valid) process(item);",
			expected: "#select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "09-unbraced-select-with-else",
			code: "select item in (items) where (item.valid) process(item); else fallback();",
			expected: "select item in (items) where (item.valid) {\n" +
				"    process(item);\n" +
				"} else {\n" +
				"    fallback();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "10-nested-select",
			code: "select item in (items) where (item.valid) {\n" +
				"select sub in (item.subs) where (sub.enabled) {\n" +
				"process(sub);\n" +
				"}\n" +
				"}",
			expected: "select item in (items) where (item.valid) {\n" +
				"    select sub in (item.subs) where (sub.enabled) {\n" +
				"        process(sub);\n" +
				"    }\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-nested-hash-select",
			code: "#select x in ([1, 2]) where (x > 0) {\n" +
				"#select y in ([3, 4]) where (y > 2) {\n" +
				"process(x, y);\n" +
				"}\n" +
				"}",
			expected: "#select x in ([1, 2]) where (x > 0) {\n" +
				"    #select y in ([3, 4]) where (y > 2) {\n" +
				"        process(x, y);\n" +
				"    }\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-messy-newlines-collapsed-in-header",
			code: "select\nitem\nin\n(\nitems\n)\nwhere\n(\ncond\n)\n{\nprocess(item);\n}",
			expected: "select item in (items) where (cond) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestSelectStatement_ItemSizes(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Short item identifier (1 char: i, r, x)
		{
			name:     "01-short-item-fits-single-line",
			code:     "select i in (bank.regs) where (i.valid) {\n  process(i);\n}",
			expected: "select i in (bank.regs) where (i.valid) {\n    process(i);\n}",
			lineCap:  41,
			pc:       func(f func()) { f() },
		},
		{
			name: "02-short-item-outer-wrap-in",
			code: "select i in (bank_a.all_registers_list) where (i.valid) {\n  process(i);\n}",
			expected: "select i in (\n" +
				"    bank_a.all_registers_list\n" +
				") where (i.valid) {\n" +
				"    process(i);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "03-short-item-outer-wrap-where",
			code: "select i in (items) where (status == STATUS_READY_TO_PROCESS) {\n  process(i);\n}",
			expected: "select i in (items) where (\n" +
				"    status == STATUS_READY_TO_PROCESS\n" +
				") {\n" +
				"    process(i);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Medium item identifier: reg_item, field_item
		{
			name:     "10-medium-item-fits-at-cap",
			code:     "select reg in (bank.regs) where (reg.v) {\n  process(reg);\n}",
			expected: "select reg in (bank.regs) where (reg.v) {\n    process(reg);\n}",
			lineCap:  41,
			pc:       func(f func()) { f() },
		},
		{
			name: "11-medium-item-forces-in-wrap",
			code: "select reg_item in (bank_a.registers_list) where (reg_item.valid) {\n" +
				"  process(reg_item);\n" +
				"}",
			expected: "select reg_item in (\n" +
				"    bank_a.registers_list\n" +
				") where (reg_item.valid) {\n" +
				"    process(reg_item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-medium-item-forces-where-wrap",
			code: "select reg_item in (regs) where (reg_item.status == STATUS_ACTIVE) {\n" +
				"  process(reg_item);\n" +
				"}",
			expected: "select reg_item in (regs) where (\n" +
				"    reg_item.status == STATUS_ACTIVE\n" +
				") {\n" +
				"    process(reg_item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Long item identifier: long_register_item, very_long_item_name
		{
			name: "20-long-item-short-exprs-fit",
			code: "select long_reg in (items) where (c) {\n  process(long_reg);\n}",
			expected: "select long_reg in (items) where (c) {\n" +
				"    process(long_reg);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-long-item-forces-in-wrap",
			code: "select very_long_item_name_a in (items) where (item.valid) {\n" +
				"  process(very_long_item_name_a);\n" +
				"}",
			expected: "select very_long_item_name_a in (\n" +
				"    items\n" +
				") where (item.valid) {\n" +
				"    process(very_long_item_name_a);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-long-item-forces-where-wrap",
			code: "select long_item in (items) where (long_item.is_ready_now) {\n" +
				"  process(long_item);\n" +
				"}",
			expected: "select long_item in (items) where (\n" +
				"    long_item.is_ready_now\n" +
				") {\n" +
				"    process(long_item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "23-long-item-hash-select-forces-wrap",
			code: "#select very_long_item_name in (items) where (item.valid) {\n" +
				"  process(very_long_item_name);\n" +
				"}",
			expected: "#select very_long_item_name in (\n" +
				"    items\n" +
				") where (item.valid) {\n" +
				"    process(very_long_item_name);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Super-long item identifier: never wraps header before ident, in, or (
		{
			name: "30-superlong-item-never-wraps-header",
			code: "select super_long_selection_variable_name in (items) where (cond) {\n" +
				"  x = 1;\n" +
				"}",
			expected: "select super_long_selection_variable_name in (\n" +
				"    items\n" +
				") where (cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "31-superlong-item-with-where-wrap",
			code: "select super_long_selection_variable_name in (items) where (" +
				"status == STATUS_READY_TO_PROCESS) {\n  x = 1;\n}",
			expected: "select super_long_selection_variable_name in (\n" +
				"    items\n" +
				") where (\n" +
				"    status == STATUS_READY_TO_PROCESS\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Keyword identifiers as items
		{
			name:     "40-keyword-ident-discard",
			code:     "select   _   in   (  items  )   where   (  true  )   {\n  x = 1;\n}",
			expected: "select _ in (items) where (true) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "41-keyword-ident-bank",
			code:     "select   bank   in   (  banks  )   where   (  bank.v  )   {\n  process(bank);\n}",
			expected: "select bank in (banks) where (bank.v) {\n    process(bank);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "42-keyword-ident-field",
			code: "select field in (fields) where (field.valid) {\n" +
				"  process(field);\n" +
				"}",
			expected: "select field in (fields) where (field.valid) {\n" +
				"    process(field);\n" +
				"}",
			lineCap: 50,
			pc:      func(f func()) { f() },
		},
		{
			name: "43-keyword-ident-port",
			code: "select port in (ports) where (port.valid) {\n" +
				"  process(port);\n" +
				"}",
			expected: "select port in (ports) where (port.valid) {\n" +
				"    process(port);\n" +
				"}",
			lineCap: 50,
			pc:      func(f func()) { f() },
		},
		{
			name: "44-keyword-ident-device",
			code: "select device in (devs) where (device.id == 1) {\n" +
				"  process(device);\n" +
				"}",
			expected: "select device in (devs) where (device.id == 1) {\n" +
				"    process(device);\n" +
				"}",
			lineCap: 50,
			pc:      func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestSelectStatement_ConditionWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - entire statement and condition fit within lineCap
		{
			name:     "01-level1-simple-relational",
			code:     "select i in (items) where (i.count < 10) {\n  x = 1;\n}",
			expected: "select i in (items) where (i.count < 10) {\n    x = 1;\n}",
			lineCap:  45,
			pc:       func(f func()) { f() },
		},
		{
			name: "02-level1-simple-logical",
			code: "select i in (items) where (i.ready && i.ok) {\n  x = 1;\n}",
			expected: "select i in (items) where (i.ready && i.ok) {\n" +
				"    x = 1;\n" +
				"}",
			lineCap: 45,
			pc:      func(f func()) { f() },
		},
		{
			name: "03-level1-ternary",
			code: "select i in (items) where (flag ? a : b) {\n  x = 1;\n}",
			expected: "select i in (items) where (flag ? a : b) {\n" +
				"    x = 1;\n" +
				"}",
			lineCap: 45,
			pc:      func(f func()) { f() },
		},
		{
			name: "04-level1-method-call",
			code: "select i in (items) where (is_valid(dev, i)) {\n  x = 1;\n}",
			expected: "select i in (items) where (is_valid(dev, i)) {\n" +
				"    x = 1;\n" +
				"}",
			lineCap: 50,
			pc:      func(f func()) { f() },
		},
		{
			name: "05-level1-bitwise-slice",
			code: "select i in (items) where ((i.val[3:0] & M) != 0) {\n  x = 1;\n}",
			expected: "select i in (items) where ((i.val[3:0] & M) != 0) {\n" +
				"    x = 1;\n" +
				"}",
			lineCap: 55,
			pc:      func(f func()) { f() },
		},

		// Level 2: Outer Parenthesized Wrap - statement exceeds lineCap,
		// but inner condition fits on one line
		{
			name: "10-level2-outer-wrap-single-relational",
			code: "select item in (items) where (status == STATUS_WAITING_FOR_DEV) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-outer-wrap-compound-logical",
			code: "select item in (items) where (flag_alpha && flag_beta && ready) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
				"    flag_alpha && flag_beta && ready\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-nested-parens-both-compact",
			code: "select item in (items) where ((cond_a && cond_b) || (cond_c && cond_d)) {\n" +
				"  x = 1;\n" +
				"}",
			expected: "select item in (items) where (\n" +
				"    (cond_a && cond_b)\n" +
				"    || (cond_c && cond_d)\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-bitwise-slice-compact",
			code: "select item in (items) where ((item.data[15:0] & STATUS_MASK) != 0) {\n" +
				"  x = 1;\n" +
				"}",
			expected: "select item in (items) where (\n" +
				"    (item.data[15:0] & STATUS_MASK) != 0\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "14-level2-member-access-compact",
			code: "select item in (items) where (device.subsystem.is_active_now()) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
				"    device.subsystem.is_active_now()\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Deep Expression Explosion - inner expressions explode into multi-line
		{
			name: "20-level3-logical-and-explosion",
			code: "select item in (items) where (" +
				"very_long_condition_name_a && very_long_condition_name_b) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-multi-operand-chain",
			code: "select item in (items) where (" +
				"flag_alpha && flag_beta && flag_gamma && flag_delta) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
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
			code: "select item in (items) where (" +
				"(very_long_cond_a && very_long_cond_b) || (cond_c && cond_d)) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
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
			code: "select item in (items) where ((very_long_cond_a && very_long_cond_b) || " +
				"(very_long_cond_c && very_long_cond_d)) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
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
			code: "select item in (items) where (" +
				"current_offset + transfer_bytes <= max_buffer_capacity) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
				"    current_offset + transfer_bytes\n" +
				"    <= max_buffer_capacity\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "25-level3-deep-arithmetic-explosion",
			code: "select item in (items) where (" +
				"very_long_current_offset + very_long_transfer_bytes <= max_buffer_capacity) {\n" +
				"  x = 1;\n" +
				"}",
			expected: "select item in (items) where (\n" +
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
			code: "select item in (items) where (" +
				"use_fast_mode ? fast_queue.has_items() : slow_queue.has_items()) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
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
			code: "select item in (items) where (" +
				"(very_long_cond_a && very_long_cond_b) ? check_fast() : check_slow()) {\n" +
				"  x = 1;\n" +
				"}",
			expected: "select item in (items) where (\n" +
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
			code: "select item in (items) where (" +
				"validate_transfer(source_addr, destination_addr, transfer_len)) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
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
			code: "select item in (items) where (" +
				"bus.controller.subsystem.peripheral.is_ready()) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
				"    bus.controller.subsystem.peripheral\n" +
				"    .is_ready()\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "30-level3-bitwise-slice-explosion",
			code: "select item in (items) where (" +
				"(buffer_data[31:16] & STATUS_MASK) != EXPECTED_PATTERN) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
				"    (buffer_data[31:16] & STATUS_MASK)\n" +
				"    != EXPECTED_PATTERN\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "31-level3-unary-negation-explosion",
			code: "select item in (items) where (" +
				"!(device_is_initialized && is_transfer_complete)) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
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
			code: "select item in (items) where (((deep_operand_a + deep_operand_b) > threshold && " +
				"deep_cond_c) || fallback_cond) {\n  x = 1;\n}",
			expected: "select item in (items) where (\n" +
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
				"select item in (items) where (\n" +
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

		// Unbraced Select with Wrapping Conditions
		{
			name: "50-unbraced-level2-outer-wrap",
			code: "select item in (items) where (status == STATUS_WAITING_FOR_DEV)\n  x = 1;",
			expected: "select item in (items) where (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "51-unbraced-level3-exploded",
			code: "select item in (items) where (" +
				"very_long_condition_name_a && very_long_condition_name_b)\n  x = 1;",
			expected: "select item in (items) where (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "52-unbraced-level3-nested",
			code: "select item in (items) where (" +
				"(very_long_cond_a && very_long_cond_b) || (cond_c && cond_d))\n  x = 1;",
			expected: "select item in (items) where (\n" +
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

func TestSelectStatement_EachInExpressions(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - fits within lineCap
		{
			name: "01-level1-each-short",
			code: "select   r   in   (  each   reg   in   (  this  )  )   where   " +
				"(  r  .  valid  )   {\n  r.reset();\n}",
			expected: "select r in (each reg in (this)) where (r.valid) {\n" +
				"    r.reset();\n" +
				"}",
			lineCap: 50,
			pc:      func(f func()) { f() },
		},
		{
			name: "02-level1-each-register-keyword",
			code: "select r in (each register in (b)) where (r.valid) {\n  r.reset();\n}",
			expected: "select r in (each register in (b)) where (r.valid) {\n" +
				"    r.reset();\n" +
				"}",
			lineCap: 55,
			pc:      func(f func()) { f() },
		},

		// Level 2: Outer select in (...) wraps, inner each stays compact
		{
			name: "10-level2-each-outer-wrap",
			code: "select item in (each hard_reset in (regs)) where (item.valid) {\n" +
				"  item.hard_reset();\n" +
				"}",
			expected: "select item in (\n" +
				"    each hard_reset in (regs)\n" +
				") where (item.valid) {\n" +
				"    item.hard_reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-each-member-chain-compact",
			code: "select r in (each reg in (this.bank_a.grp)) where (r.valid) {\n" +
				"  r.reset();\n" +
				"}",
			expected: "select r in (\n" +
				"    each reg in (this.bank_a.grp)\n" +
				") where (r.valid) {\n" +
				"    r.reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-each-long-item-ident",
			code: "select register_obj in (each reg in (this)) where (register_obj.valid) {\n" +
				"  register_obj.reset();\n" +
				"}",
			expected: "select register_obj in (\n" +
				"    each reg in (this)\n" +
				") where (register_obj.valid) {\n" +
				"    register_obj.reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-nested-each-compact",
			code: "select i in (each a in (each b in (this))) where (i.valid) {\n" +
				"  process(i);\n" +
				"}",
			expected: "select i in (\n" +
				"    each a in (each b in (this))\n" +
				") where (i.valid) {\n" +
				"    process(i);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Both outer select in (...) and inner each ... in (...) wrap
		{
			name: "20-level3-each-target-wraps",
			code: "select item in (each hard_reset in (this.main_register_bank)) where (item.valid) {\n" +
				"  item.hard_reset();\n" +
				"}",
			expected: "select item in (\n" +
				"    each hard_reset in (\n" +
				"        this.main_register_bank\n" +
				"    )\n" +
				") where (item.valid) {\n" +
				"    item.hard_reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-each-long-template-ident",
			code: "select item in (each very_long_template_name in (this.regs)) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (\n" +
				"    each very_long_template_name in (\n" +
				"        this.regs\n" +
				"    )\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-each-superlong-template-ident",
			code: "select item in (each super_long_template_identifier_name in (regs)) " +
				"where (item.valid) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    each super_long_template_identifier_name in (\n" +
				"        regs\n" +
				"    )\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-each-array-index-target",
			code: "select item in (each register in (this.banks[bank_index])) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (\n" +
				"    each register in (\n" +
				"        this.banks[bank_index]\n" +
				"    )\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "24-level3-each-member-chain-explosion",
			code: "select item in (each reg in (dev.subsystem.controller.bank_a)) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (\n" +
				"    each reg in (\n" +
				"        dev.subsystem.controller.bank_a\n" +
				"    )\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 4: Deep nested each ... in (each ... in (...))
		{
			name: "30-level4-nested-each-exploded",
			code: "select item in (each sub_tmpl in (" +
				"each group_tmpl in (this.main_bank))) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (\n" +
				"    each sub_tmpl in (\n" +
				"        each group_tmpl in (\n" +
				"            this.main_bank\n" +
				"        )\n" +
				"    )\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Unbraced select with each in
		{
			name: "40-unbraced-level2-each-in",
			code: "select item in (each hard_reset in (regs)) where (item.valid)\n" +
				"  item.hard_reset();",
			expected: "select item in (\n" +
				"    each hard_reset in (regs)\n" +
				") where (item.valid) {\n" +
				"    item.hard_reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestSelectStatement_CompileTimeArrays(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - fits within lineCap
		{
			name: "01-level1-short-int-array",
			code: "#select   x   in   (  [ 1 , 2 , 3 ]  )   where   (  x  >  1  )   {\n" +
				"  process(x);\n" +
				"}",
			expected: "#select x in ([1, 2, 3]) where (x > 1) {\n" +
				"    process(x);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "02-level1-short-ident-array",
			code: "#select x in ([a, b, c, d]) where (x.valid) {\n  process(x);\n}",
			expected: "#select x in ([a, b, c, d]) where (x.valid) {\n" +
				"    process(x);\n" +
				"}",
			lineCap: 45,
			pc:      func(f func()) { f() },
		},
		{
			name:     "03-level1-empty-array",
			code:     "#select x in ([   ]) where (true) {\n  process(x);\n}",
			expected: "#select x in ([]) where (true) {\n    process(x);\n}",
			pc:       func(f func()) { f() },
		},

		// Level 2: Outer #select in (...) wraps, inner [...] fits on a single line
		{
			name: "10-level2-outer-wrap-array",
			code: "#select item in ([a, b, c, d, e, f, g]) where (item.valid) {\n  process(item);\n}",
			expected: "#select item in (\n" +
				"    [a, b, c, d, e, f, g]\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-outer-wrap-binary-expr-array",
			code: "#select item in ([a + b, c * d, e - f]) where (item > 0) {\n  process(item);\n}",
			expected: "#select item in (\n" +
				"    [a + b, c * d, e - f]\n" +
				") where (item > 0) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-long-item-ident-with-short-array",
			code: "#select very_long_item_name in ([1, 2, 3]) where (item > 0) {\n" +
				"  process(very_long_item_name);\n" +
				"}",
			expected: "#select very_long_item_name in (\n" +
				"    [1, 2, 3]\n" +
				") where (item > 0) {\n" +
				"    process(very_long_item_name);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Multi-line exploded [...] inside #select
		{
			name: "20-level3-exploded-array-all-items",
			code: "#select item in ([aaaaaaaa, bbbbbbbb, cccccccc, dddddddd]) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "#select item in (\n" +
				"    [\n" +
				"        aaaaaaaa,\n" +
				"        bbbbbbbb,\n" +
				"        cccccccc,\n" +
				"        dddddddd\n" +
				"    ]\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-exploded-array-one-superlong-item",
			code: "#select item in ([aaaa, bbbb, cccccccccccccccccccccccccccc, dddd]) " +
				"where (item.valid) {\n  process(item);\n}",
			expected: "#select item in (\n" +
				"    [\n" +
				"        aaaa,\n" +
				"        bbbb,\n" +
				"        cccccccccccccccccccccccccccc,\n" +
				"        dddd\n" +
				"    ]\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-exploded-array-binary-exprs",
			code: "#select item in (" +
				"[aaaaaaa + aaaaaaaa, bbbbbbbbbb - bbbbbbbbbb, ccccccccc * cccccccc]) " +
				"where (item > 0) {\n  process(item);\n}",
			expected: "#select item in (\n" +
				"    [\n" +
				"        aaaaaaa + aaaaaaaa,\n" +
				"        bbbbbbbbbb - bbbbbbbbbb,\n" +
				"        ccccccccc * cccccccc\n" +
				"    ]\n" +
				") where (item > 0) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-nested-compile-time-array",
			code: "#select item in ([aaaa, [bbbb, cccc], dddd]) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "#select item in (\n" +
				"    [\n" +
				"        aaaa,\n" +
				"        [bbbb, cccc],\n" +
				"        dddd\n" +
				"    ]\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 4: Deep nested compile-time array (2 levels)
		{
			name: "24-level4-nested-compile-time-array-2levels",
			code: "#select item in ([aaaa, [bbbb, [cccc, dddd]], eeee]) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "#select item in (\n" +
				"    [\n" +
				"        aaaa,\n" +
				"        [\n" +
				"            bbbb,\n" +
				"            [cccc, dddd]\n" +
				"        ],\n" +
				"        eeee\n" +
				"    ]\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Unbraced #select with compile-time array
		{
			name: "30-unbraced-compile-time-array",
			code: "#select item in ([a, b, c, d, e, f, g]) where (item.valid)\n  process(item);",
			expected: "#select item in (\n" +
				"    [a, b, c, d, e, f, g]\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestSelectStatement_ComplexExpressions(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line
		{
			name:     "01-level1-member-access",
			code:     "select item in (this.bank.regs) where (item.valid) {\n  process(item);\n}",
			expected: "select item in (this.bank.regs) where (item.valid) {\n    process(item);\n}",
			lineCap:  55,
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-level1-method-call",
			code:     "select item in (get_items(dev)) where (item.valid) {\n  process(item);\n}",
			expected: "select item in (get_items(dev)) where (item.valid) {\n    process(item);\n}",
			lineCap:  55,
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-level1-array-index",
			code:     "select item in (banks[i].regs) where (item.valid) {\n  process(item);\n}",
			expected: "select item in (banks[i].regs) where (item.valid) {\n    process(item);\n}",
			lineCap:  55,
			pc:       func(f func()) { f() },
		},

		// Level 2: Outer select in (...) wrap, compact inner expression
		{
			name: "10-level2-member-chain-outer-wrap",
			code: "select item in (device.subsystem.active_items) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (\n" +
				"    device.subsystem.active_items\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-method-call-outer-wrap",
			code: "select item in (fetch_active_items(dev, idx)) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (\n" +
				"    fetch_active_items(dev, idx)\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-ternary-outer-wrap",
			code: "select item in (fast ? fast_items : slow_items) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (\n" +
				"    fast ? fast_items : slow_items\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-const-ternary-array-outer-wrap",
			code: "#select x in (flag #? [1, 2, 3] #: [4, 5, 6]) where (x > 0) {\n" +
				"  process(x);\n" +
				"}",
			expected: "#select x in (\n" +
				"    flag #? [1, 2, 3] #: [4, 5, 6]\n" +
				") where (x > 0) {\n" +
				"    process(x);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Deep expression explosion inside select in (...)
		{
			name: "20-level3-method-call-args-explosion",
			code: "select item in (collect_items(source_bank, target_group, filter_mask)) " +
				"where (item.valid) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    collect_items(\n" +
				"        source_bank,\n" +
				"        target_group,\n" +
				"        filter_mask\n" +
				"    )\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-ternary-explosion",
			code: "select item in (" +
				"use_fast_mode ? fast_queue.get_items() : slow_queue.get_items()) " +
				"where (item.valid) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    use_fast_mode\n" +
				"        ? fast_queue.get_items()\n" +
				"        : slow_queue.get_items()\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-const-ternary-arrays-explosion",
			code: "#select item in (" +
				"use_extended_list #? [alpha_item, beta_item] #: [gamma_item, delta_item]) " +
				"where (item.valid) {\n  process(item);\n}",
			expected: "#select item in (\n" +
				"    use_extended_list\n" +
				"        #? [alpha_item, beta_item]\n" +
				"        #: [gamma_item, delta_item]\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-cast-expression",
			code: "select item in (cast(raw_sequence_ptr, sequence(register))) " +
				"where (item.valid) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    cast(\n" +
				"        raw_sequence_ptr,\n" +
				"        sequence(register)\n" +
				"    )\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Unbraced select with wrapping expressions
		{
			name: "30-unbraced-level2-member-chain",
			code: "select item in (device.subsystem.active_items) where (item.valid)\n" +
				"  process(item);",
			expected: "select item in (\n" +
				"    device.subsystem.active_items\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "31-unbraced-level3-exploded",
			code: "select item in (very_long_expression_a + very_long_expression_b) " +
				"where (item.valid)\n  process(item);",
			expected: "select item in (\n" +
				"    very_long_expression_a\n" +
				"    + very_long_expression_b\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestSelectStatement_CombinedWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// in fits, where wraps (outer wrap)
		{
			name: "01-in-fits-where-outer-wrap",
			code: "select item in (items) where (status == STATUS_READY_TO_PROCESS) {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (items) where (\n" +
				"    status == STATUS_READY_TO_PROCESS\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// in fits, where explodes
		{
			name: "02-in-fits-where-exploded",
			code: "select item in (items) where (" +
				"very_long_condition_name_a && very_long_condition_name_b) {\n  process(item);\n}",
			expected: "select item in (items) where (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// in wraps (outer wrap), where fits
		{
			name: "03-in-outer-wrap-where-fits",
			code: "select item in (very_long_expression_collection_list) where (item.valid) {\n" +
				"  process(item);\n" +
				"}",
			expected: "select item in (\n" +
				"    very_long_expression_collection_list\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// in explodes, where fits
		{
			name: "04-in-exploded-where-fits",
			code: "select item in (very_long_expression_a + very_long_expression_b) " +
				"where (item.valid) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    very_long_expression_a\n" +
				"    + very_long_expression_b\n" +
				") where (item.valid) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Both in and where outer wrap
		{
			name: "05-both-outer-wrap",
			code: "select item in (very_long_expression_collection_list) where (" +
				"status == STATUS_WAITING_FOR_DEV) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    very_long_expression_collection_list\n" +
				") where (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// in outer wrap, where exploded
		{
			name: "06-in-outer-wrap-where-exploded",
			code: "select item in (very_long_expression_collection_list) where (" +
				"very_long_condition_name_a && very_long_condition_name_b) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    very_long_expression_collection_list\n" +
				") where (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// in exploded, where outer wrap
		{
			name: "07-in-exploded-where-outer-wrap",
			code: "select item in (very_long_expression_a + very_long_expression_b) where (" +
				"status == STATUS_WAITING_FOR_DEV) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    very_long_expression_a\n" +
				"    + very_long_expression_b\n" +
				") where (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Both in and where exploded
		{
			name: "08-both-exploded",
			code: "select item in (very_long_expression_a + very_long_expression_b) where (" +
				"very_long_condition_name_a && very_long_condition_name_b) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    very_long_expression_a\n" +
				"    + very_long_expression_b\n" +
				") where (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// in each-in wrapped, where exploded
		{
			name: "09-in-each-in-wrapped-where-exploded",
			code: "select item in (each hard_reset in (this.main_register_bank)) where (" +
				"very_long_condition_name_a && very_long_condition_name_b) {\n  process(item);\n}",
			expected: "select item in (\n" +
				"    each hard_reset in (\n" +
				"        this.main_register_bank\n" +
				"    )\n" +
				") where (\n" +
				"    very_long_condition_name_a\n" +
				"    && very_long_condition_name_b\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// in compile-time array exploded, where outer wrap
		{
			name: "10-in-array-exploded-where-outer-wrap",
			code: "#select item in ([aaaaaaaa, bbbbbbbb, cccccccc, dddddddd]) where (" +
				"status == STATUS_WAITING_FOR_DEV) {\n  process(item);\n}",
			expected: "#select item in (\n" +
				"    [\n" +
				"        aaaaaaaa,\n" +
				"        bbbbbbbb,\n" +
				"        cccccccc,\n" +
				"        dddddddd\n" +
				"    ]\n" +
				") where (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// hash-select with wrapped in and wrapped where with hash-else
		{
			name: "11-hash-select-hash-else-both-wrapped",
			code: "#select item in (very_long_expression_collection_list) where (" +
				"status == STATUS_WAITING_FOR_DEV) {\n" +
				"  process(item);\n" +
				"} #else {\n" +
				"  fallback();\n" +
				"}",
			expected: "#select item in (\n" +
				"    very_long_expression_collection_list\n" +
				") where (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    process(item);\n" +
				"} #else {\n" +
				"    fallback();\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Chained else select with combined wrapping in both branches
		{
			name: "12-chained-else-select-both-wrapped",
			code: "select item in (collection_a_very_long_list) where (" +
				"status == STATUS_WAITING_FOR_DEV) {\n" +
				"  process(item);\n" +
				"} else select item in (collection_b_very_long_list) where (" +
				"status == STATUS_WAITING_FOR_DEV) {\n" +
				"  process(item);\n" +
				"} else {\n" +
				"  fallback();\n" +
				"}",
			expected: "select item in (\n" +
				"    collection_a_very_long_list\n" +
				") where (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    process(item);\n" +
				"} else select item in (\n" +
				"    collection_b_very_long_list\n" +
				") where (\n" +
				"    status == STATUS_WAITING_FOR_DEV\n" +
				") {\n" +
				"    process(item);\n" +
				"} else {\n" +
				"    fallback();\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestSelectStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 50
	testCases := []printerTestCase{
		{
			name: "01-leading-comment",
			code: "// c\n" +
				"select item in (items) where (cond) {\n" +
				"x = 1;\n" +
				"}",
			expected: "// c\n" +
				"select item in (items) where (cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "02-comment-between-select-and-ident",
			code: "select /* c */ item in (items) where (cond) {\n" +
				"x = 1;\n" +
				"}",
			expected: "select /* c */ item in (items) where (cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "03-comment-between-ident-and-in",
			code: "select item /* c */ in (items) where (cond) {\n" +
				"x = 1;\n" +
				"}",
			expected: "select item /* c */ in (items) where (cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "04-comment-between-in-and-lparen",
			code: "select item in /* c */ (items) where (cond) {\n" +
				"x = 1;\n" +
				"}",
			expected: "select item in /* c */ (items) where (cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "05-comment-inside-in-paren",
			code: "select item in (/* c */ items) where (cond) {\n" +
				"x = 1;\n" +
				"}",
			expected: "select item in (/* c */ items) where (cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "06-comment-before-in-rparen",
			code: "select item in (items /* c */) where (cond) {\n" +
				"x = 1;\n" +
				"}",
			expected: "select item in (items /* c */) where (cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "07-comment-between-in-rparen-and-where",
			code: "select item in (items) /* c */ where (cond) {\n" +
				"x = 1;\n" +
				"}",
			expected: "select item in (items) /* c */ where (cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "08-comment-between-where-and-lparen",
			code: "select item in (items) where /* c */ (cond) {\n" +
				"x = 1;\n" +
				"}",
			expected: "select item in (items) where /* c */ (cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "09-comment-inside-where-condition",
			code: "select item in (items) where (/* c */ cond) {\n" +
				"x = 1;\n" +
				"}",
			expected: "select item in (items) where (/* c */ cond) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "10-comment-before-where-rparen",
			code: "select item in (items) where (cond /* c */) {\n" +
				"x = 1;\n" +
				"}",
			expected: "select item in (items) where (cond /* c */) {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-comment-between-where-rparen-and-brace",
			code: "select item in (items) where (cond) /* c */ {\n" +
				"x = 1;\n" +
				"}",
			expected: "select item in (items) where (cond) /* c */ {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-comment-between-brace-and-else",
			code: "select item in (items) where (cond) {\n" +
				"x = 1;\n" +
				"} /* c */ else {\n" +
				"x = 2;\n" +
				"}",
			expected: "select item in (items) where (cond) {\n" +
				"    x = 1;\n" +
				"} /* c */ else {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-comment-between-brace-and-hash-else",
			code: "#select item in (items) where (cond) {\n" +
				"x = 1;\n" +
				"} /* c */ #else {\n" +
				"x = 2;\n" +
				"}",
			expected: "#select item in (items) where (cond) {\n" +
				"    x = 1;\n" +
				"} /* c */ #else {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "14-comment-after-else",
			code: "select item in (items) where (cond) {\n" +
				"x = 1;\n" +
				"} else /* c */ {\n" +
				"x = 2;\n" +
				"}",
			expected: "select item in (items) where (cond) {\n" +
				"    x = 1;\n" +
				"} else /* c */ {\n" +
				"    x = 2;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "15-trailing-comment",
			code: "select item in (items) where (cond) {\n" +
				"x = 1;\n" +
				"} // done",
			expected: "select item in (items) where (cond) {\n" +
				"    x = 1;\n" +
				"} // done",
			pc: func(f func()) { f() },
		},
		{
			name: "16-comments-inside-each-in-expr",
			code: "select item in (each /*c1*/ reg /*c2*/ in (/*c3*/ this)) " +
				"where (/*c4*/ item.valid /*c5*/) {\n" +
				"process(item);\n" +
				"}",
			expected: "select item in (\n" +
				"    each /*c1*/ reg /*c2*/ in (/*c3*/ this)\n" +
				") where (/*c4*/ item.valid /*c5*/) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "17-comments-inside-compile-time-array",
			code: "#select item in ([a, /* c */ b, c, d]) where (item.valid) {\n" +
				"process(item);\n" +
				"}",
			expected: "#select item in ([a, /* c */ b, c, d]) where (\n" +
				"    item.valid\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "18-line-comment-after-lparen",
			code: "select item in ( // c1\n" +
				"items\n" +
				") where ( // c2\n" +
				"cond\n" +
				") {\n" +
				"process(item);\n" +
				"}",
			expected: "select item in ( // c1\n" +
				"    items\n" +
				") where ( // c2\n" +
				"    cond\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "19-line-comment-inside-array",
			code: "#select item in ([\n" +
				"a, // c\n" +
				"b\n" +
				"]) where (item > 0) {\n" +
				"process(item);\n" +
				"}",
			expected: "#select item in (\n" +
				"    [\n" +
				"        a, // c\n" +
				"        b\n" +
				"    ]\n" +
				") where (item > 0) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "20-hash-select-all-inline-comments",
			code: "#select /*c1*/ x /*c2*/ in /*c3*/ (/*c4*/ [1, 2] /*c5*/) /*c6*/ " +
				"where /*c7*/ (/*c8*/ x > 0 /*c9*/) /*c10*/ {\n" +
				"process(x);\n" +
				"}",
			expected: "#select /*c1*/ x /*c2*/ in /*c3*/ (\n" +
				"    /*c4*/ [1, 2] /*c5*/\n" +
				") /*c6*/ where /*c7*/ (\n" +
				"    /*c8*/ x > 0 /*c9*/\n" +
				") /*c10*/ {\n" +
				"    process(x);\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
