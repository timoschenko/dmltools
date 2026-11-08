// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestForeachStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-basic-foreach",
			code:     "foreach   item   in   (  items  )   {\n  process(item);\n}",
			expected: "foreach item in (items) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-hash-foreach",
			code:     "#foreach   item   in   (  items  )   {\n  process(item);\n}",
			expected: "#foreach item in (items) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-unbraced-foreach",
			code:     "foreach   item   in   (  items  )   process(item);",
			expected: "foreach item in (items) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-unbraced-hash-foreach",
			code:     "#foreach   item   in   (  items  )   process(item);",
			expected: "#foreach item in (items) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "05-line-wrapping-expression",
			code: "foreach item in (very_long_expression_a + very_long_expression_b) {\n" +
				"    process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    very_long_expression_a\n" +
				"    + very_long_expression_b\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "06-nested-foreach",
			code: "foreach item in (items) {\n" +
				"foreach sub in (item.subs) {\n" +
				"process(sub);\n" +
				"}\n" +
				"}",
			expected: "foreach item in (items) {\n" +
				"    foreach sub in (item.subs) {\n" +
				"        process(sub);\n" +
				"    }\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "07-nested-hash-foreach",
			code: "#foreach x in ([1, 2]) {\n" +
				"#foreach y in ([3, 4]) {\n" +
				"process(x, y);\n" +
				"}\n" +
				"}",
			expected: "#foreach x in ([1, 2]) {\n" +
				"    #foreach y in ([3, 4]) {\n" +
				"        process(x, y);\n" +
				"    }\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "08-messy-newlines-collapsed-in-header",
			code: "foreach\nitem\nin\n(\nitems\n)\n{\nprocess(item);\n}",
			expected: "foreach item in (items) {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name:     "09-keyword-ident-discard",
			code:     "foreach   _   in   (  items  )   {\n  x = 1;\n}",
			expected: "foreach _ in (items) {\n    x = 1;\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-keyword-ident-bank",
			code:     "foreach   bank   in   (  banks  )   {\n  process(bank);\n}",
			expected: "foreach bank in (banks) {\n    process(bank);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-keyword-ident-field",
			code:     "foreach   field   in   (  fields  )   {\n  process(field);\n}",
			expected: "foreach field in (fields) {\n    process(field);\n}",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestForeachStatement_ItemSizes(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Short item identifier (1 char): fits longer expressions on a single line
		{
			name:     "01-short-item-fits-single-line",
			code:     "foreach i in (bank_a.registers_list) {\n  process(i);\n}",
			expected: "foreach i in (bank_a.registers_list) {\n    process(i);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "02-short-item-outer-wrap",
			code: "foreach i in (bank_a.all_registers_list) {\n  process(i);\n}",
			expected: "foreach i in (\n" +
				"    bank_a.all_registers_list\n" +
				") {\n" +
				"    process(i);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Medium item identifier: exact lineCap boundary vs forcing (expr) to wrap
		{
			name:     "10-medium-item-fits-at-cap",
			code:     "foreach reg_item in (bank_a.regs_list) {\n  process(reg_item);\n}",
			expected: "foreach reg_item in (bank_a.regs_list) {\n    process(reg_item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "11-medium-item-forces-expr-wrap",
			code: "foreach reg_item in (bank_a.registers_list) {\n  process(reg_item);\n}",
			expected: "foreach reg_item in (\n" +
				"    bank_a.registers_list\n" +
				") {\n" +
				"    process(reg_item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-medium-item-with-exploded-expr",
			code: "foreach register_item in (very_long_expression_a + very_long_expression_b) {\n" +
				"  process(register_item);\n" +
				"}",
			expected: "foreach register_item in (\n" +
				"    very_long_expression_a\n" +
				"    + very_long_expression_b\n" +
				") {\n" +
				"    process(register_item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Long item identifier: forces even a short expression to wrap inside (expr)
		{
			name: "20-long-item-short-expr-fits",
			code: "foreach long_register_item in (items) {\n" +
				"  process(long_register_item);\n" +
				"}",
			expected: "foreach long_register_item in (items) {\n" +
				"    process(long_register_item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-long-item-forces-short-expr-wrap",
			code: "foreach very_long_item_name_a in (items) {\n" +
				"  process(very_long_item_name_a);\n" +
				"}",
			expected: "foreach very_long_item_name_a in (\n" +
				"    items\n" +
				") {\n" +
				"    process(very_long_item_name_a);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-long-item-hash-foreach-forces-wrap",
			code: "#foreach very_long_item_name in (items) {\n" +
				"  process(very_long_item_name);\n" +
				"}",
			expected: "#foreach very_long_item_name in (\n" +
				"    items\n" +
				") {\n" +
				"    process(very_long_item_name);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Super-long item identifier: `foreach <ident> in (` exceeds lineCap by itself,
		// but never wraps before `<ident>`, `in`, or `(`; only `(expr)` wraps.
		{
			name: "30-superlong-item-never-wraps-header",
			code: "foreach super_long_iteration_variable_name in (items) {\n" +
				"  x = 1;\n" +
				"}",
			expected: "foreach super_long_iteration_variable_name in (\n" +
				"    items\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "31-superlong-item-hash-foreach-with-exploded-expr",
			code: "#foreach super_long_iteration_variable_name in (" +
				"very_long_expression_a + very_long_expression_b) {\n" +
				"  x = 1;\n" +
				"}",
			expected: "#foreach super_long_iteration_variable_name in (\n" +
				"    very_long_expression_a\n" +
				"    + very_long_expression_b\n" +
				") {\n" +
				"    x = 1;\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestForeachStatement_EachInExpressions(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - entire `foreach ... in (each ... in (...)) {` fits within lineCap
		{
			name:     "01-level1-each-short",
			code:     "foreach   r   in   (  each   reg   in   (  this  )  )   {\n  r.reset();\n}",
			expected: "foreach r in (each reg in (this)) {\n    r.reset();\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-level1-each-register-keyword",
			code:     "foreach r in (each register in (b)) {\n  r.reset();\n}",
			expected: "foreach r in (each register in (b)) {\n    r.reset();\n}",
			pc:       func(f func()) { f() },
		},

		// Level 2: Outer `foreach (...)` wraps, inner `each <ident> in (<target>)` stays on one line
		{
			name: "10-level2-each-outer-wrap",
			code: "foreach item in (each hard_reset in (regs)) {\n  item.hard_reset();\n}",
			expected: "foreach item in (\n" +
				"    each hard_reset in (regs)\n" +
				") {\n" +
				"    item.hard_reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-each-member-chain-compact",
			code: "foreach r in (each reg in (this.bank_a.grp)) {\n  r.reset();\n}",
			expected: "foreach r in (\n" +
				"    each reg in (this.bank_a.grp)\n" +
				") {\n" +
				"    r.reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-each-long-item-ident",
			code: "foreach register_obj in (each reg in (this)) {\n  register_obj.reset();\n}",
			expected: "foreach register_obj in (\n" +
				"    each reg in (this)\n" +
				") {\n" +
				"    register_obj.reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-nested-each-compact",
			code: "foreach i in (each a in (each b in (this))) {\n  process(i);\n}",
			expected: "foreach i in (\n" +
				"    each a in (each b in (this))\n" +
				") {\n" +
				"    process(i);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Both outer `foreach (...)` and inner `each ... in (...)` wrap
		{
			name: "20-level3-each-target-wraps",
			code: "foreach item in (each hard_reset in (this.main_register_bank)) {\n" +
				"  item.hard_reset();\n" +
				"}",
			expected: "foreach item in (\n" +
				"    each hard_reset in (\n" +
				"        this.main_register_bank\n" +
				"    )\n" +
				") {\n" +
				"    item.hard_reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-each-long-template-ident",
			code: "foreach item in (each very_long_template_name in (this.regs)) {\n" +
				"  process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    each very_long_template_name in (\n" +
				"        this.regs\n" +
				"    )\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-each-superlong-template-ident",
			code: "foreach item in (each super_long_template_identifier_name in (regs)) {\n" +
				"  process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    each super_long_template_identifier_name in (\n" +
				"        regs\n" +
				"    )\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-each-array-index-target",
			code: "foreach item in (each register in (this.banks[bank_index])) {\n" +
				"  process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    each register in (\n" +
				"        this.banks[bank_index]\n" +
				"    )\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "24-level3-each-member-chain-explosion",
			code: "foreach item in (each reg in (dev.subsystem.controller.bank_a)) {\n" +
				"  process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    each reg in (\n" +
				"        dev.subsystem.controller.bank_a\n" +
				"    )\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "25-level3-each-long-member-chain-explosion",
			code: "foreach item in (each reg in (device.subsystem.controller.bank_a)) {\n" +
				"  process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    each reg in (\n" +
				"        device.subsystem.controller\n" +
				"        .bank_a\n" +
				"    )\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 4: Deep nested `each ... in (each ... in (...))`
		{
			name: "30-level4-nested-each-exploded",
			code: "foreach item in (each sub_tmpl in (" +
				"each group_tmpl in (this.main_bank))) {\n" +
				"  process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    each sub_tmpl in (\n" +
				"        each group_tmpl in (\n" +
				"            this.main_bank\n" +
				"        )\n" +
				"    )\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestForeachStatement_CompileTimeArrays(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - entire `#foreach ... in ([...]) {` fits within lineCap
		{
			name:     "01-level1-short-int-array",
			code:     "#foreach   x   in   (  [ 1 , 2 , 3 ]  )   {\n  process(x);\n}",
			expected: "#foreach x in ([1, 2, 3]) {\n    process(x);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-level1-short-ident-array",
			code:     "#foreach x in ([a, b, c, d]) {\n  process(x);\n}",
			expected: "#foreach x in ([a, b, c, d]) {\n    process(x);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-level1-empty-array",
			code:     "#foreach x in ([   ]) {\n  process(x);\n}",
			expected: "#foreach x in ([]) {\n    process(x);\n}",
			pc:       func(f func()) { f() },
		},

		// Level 2: Outer `foreach (...)` wraps, inner `[...]` fits on a single line
		{
			name: "10-level2-outer-wrap-array",
			code: "#foreach item in ([a, b, c, d, e, f, g]) {\n  process(item);\n}",
			expected: "#foreach item in (\n" +
				"    [a, b, c, d, e, f, g]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-outer-wrap-binary-expr-array",
			code: "#foreach item in ([a + b, c * d, e - f]) {\n  process(item);\n}",
			expected: "#foreach item in (\n" +
				"    [a + b, c * d, e - f]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-long-item-ident-with-short-array",
			code: "#foreach very_long_item_name in ([1, 2, 3]) {\n" +
				"  process(very_long_item_name);\n" +
				"}",
			expected: "#foreach very_long_item_name in (\n" +
				"    [1, 2, 3]\n" +
				") {\n" +
				"    process(very_long_item_name);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Multi-line exploded `[...]` inside `#foreach` / `foreach`
		{
			name: "20-level3-exploded-array-all-items",
			code: "#foreach item in ([aaaaaaaa, bbbbbbbb, cccccccc, dddddddd]) {\n" +
				"  process(item);\n" +
				"}",
			expected: "#foreach item in (\n" +
				"    [\n" +
				"        aaaaaaaa,\n" +
				"        bbbbbbbb,\n" +
				"        cccccccc,\n" +
				"        dddddddd\n" +
				"    ]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-exploded-array-one-superlong-item",
			code: "#foreach item in ([aaaa, bbbb, cccccccccccccccccccccccccccc, dddd]) {\n" +
				"  process(item);\n" +
				"}",
			expected: "#foreach item in (\n" +
				"    [\n" +
				"        aaaa,\n" +
				"        bbbb,\n" +
				"        cccccccccccccccccccccccccccc,\n" +
				"        dddd\n" +
				"    ]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-exploded-array-binary-exprs",
			code: "#foreach item in (" +
				"[aaaaaaa + aaaaaaaa, bbbbbbbbbb - bbbbbbbbbb, ccccccccc * cccccccc]) {\n" +
				"  process(item);\n" +
				"}",
			expected: "#foreach item in (\n" +
				"    [\n" +
				"        aaaaaaa + aaaaaaaa,\n" +
				"        bbbbbbbbbb - bbbbbbbbbb,\n" +
				"        ccccccccc * cccccccc\n" +
				"    ]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-nested-compile-time-array",
			code: "#foreach item in ([aaaa, [bbbb, cccc], dddd]) {\n" +
				"  process(item);\n" +
				"}",
			expected: "#foreach item in (\n" +
				"    [\n" +
				"        aaaa,\n" +
				"        [bbbb, cccc],\n" +
				"        dddd\n" +
				"    ]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "24-level4-nested-compile-time-array-2levels",
			code: "#foreach item in ([aaaa, [bbbb, [cccc, dddd]], eeee]) {\n" +
				"  process(item);\n" +
				"}",
			expected: "#foreach item in (\n" +
				"    [\n" +
				"        aaaa,\n" +
				"        [\n" +
				"            bbbb,\n" +
				"            [cccc, dddd]\n" +
				"        ],\n" +
				"        eeee\n" +
				"    ]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestForeachStatement_ComplexExpressions(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line
		{
			name:     "01-level1-member-access",
			code:     "foreach item in (this.bank.regs) {\n  process(item);\n}",
			expected: "foreach item in (this.bank.regs) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-level1-method-call",
			code:     "foreach item in (get_items(dev)) {\n  process(item);\n}",
			expected: "foreach item in (get_items(dev)) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-level1-array-index",
			code:     "foreach item in (banks[i].regs) {\n  process(item);\n}",
			expected: "foreach item in (banks[i].regs) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},

		// Level 2: Outer `foreach (...)` wrap, compact inner expression
		{
			name: "10-level2-member-chain-outer-wrap",
			code: "foreach item in (device.subsystem.active_items) {\n  process(item);\n}",
			expected: "foreach item in (\n" +
				"    device.subsystem.active_items\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-method-call-outer-wrap",
			code: "foreach item in (fetch_active_items(dev, idx)) {\n  process(item);\n}",
			expected: "foreach item in (\n" +
				"    fetch_active_items(dev, idx)\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-ternary-outer-wrap",
			code: "foreach item in (fast ? fast_items : slow_items) {\n  process(item);\n}",
			expected: "foreach item in (\n" +
				"    fast ? fast_items : slow_items\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-const-ternary-array-outer-wrap",
			code: "#foreach x in (flag #? [1, 2, 3] #: [4, 5, 6]) {\n  process(x);\n}",
			expected: "#foreach x in (\n" +
				"    flag #? [1, 2, 3] #: [4, 5, 6]\n" +
				") {\n" +
				"    process(x);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Level 3: Deep expression explosion inside `foreach (...)`
		{
			name: "20-level3-method-call-args-explosion",
			code: "foreach item in (collect_items(source_bank, target_group, filter_mask)) {\n" +
				"  process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    collect_items(\n" +
				"        source_bank,\n" +
				"        target_group,\n" +
				"        filter_mask\n" +
				"    )\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-ternary-explosion",
			code: "foreach item in (" +
				"use_fast_mode ? fast_queue.get_items() : slow_queue.get_items()) {\n" +
				"  process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    use_fast_mode\n" +
				"        ? fast_queue.get_items()\n" +
				"        : slow_queue.get_items()\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-const-ternary-arrays-explosion",
			code: "#foreach item in (" +
				"use_extended_list #? [alpha_item, beta_item] #: [gamma_item, delta_item]) {\n" +
				"  process(item);\n" +
				"}",
			expected: "#foreach item in (\n" +
				"    use_extended_list\n" +
				"        #? [alpha_item, beta_item]\n" +
				"        #: [gamma_item, delta_item]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-cast-expression",
			code: "foreach item in (cast(raw_sequence_ptr, sequence(register))) {\n" +
				"  process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    cast(\n" +
				"        raw_sequence_ptr,\n" +
				"        sequence(register)\n" +
				"    )\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},

		// Unbraced foreach with wrapping expressions (synthetic loop braces expected)
		{
			name: "30-unbraced-level2-each-in",
			code: "foreach item in (each hard_reset in (regs))\n  item.hard_reset();",
			expected: "foreach item in (\n" +
				"    each hard_reset in (regs)\n" +
				") {\n" +
				"    item.hard_reset();\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "31-unbraced-level3-exploded",
			code: "foreach item in (very_long_expression_a + very_long_expression_b)\n" +
				"  process(item);",
			expected: "foreach item in (\n" +
				"    very_long_expression_a\n" +
				"    + very_long_expression_b\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestForeachStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-leading-comment",
			code:     "// c\nforeach item in (items) {\nprocess(item);\n}",
			expected: "// c\nforeach item in (items) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-comment-between-foreach-and-ident",
			code:     "foreach /* c */ item in (items) {\nprocess(item);\n}",
			expected: "foreach /* c */ item in (items) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-around-in",
			code:     "foreach item /* c */ in (items) {\nprocess(item);\n}",
			expected: "foreach item /* c */ in (items) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-inside-paren",
			code:     "foreach item in (/* c */ items) {\nprocess(item);\n}",
			expected: "foreach item in (/* c */ items) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-between-paren-and-brace",
			code:     "foreach item in (items) /* c */ {\nprocess(item);\n}",
			expected: "foreach item in (items) /* c */ {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-trailing-comment",
			code:     "foreach item in (items) {\nprocess(item);\n} // done",
			expected: "foreach item in (items) {\n    process(item);\n} // done",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-comment-between-in-and-lparen",
			code:     "foreach item in /* c */ (items) {\nprocess(item);\n}",
			expected: "foreach item in /* c */ (items) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-comment-before-rparen",
			code:     "foreach item in (items /* c */) {\nprocess(item);\n}",
			expected: "foreach item in (items /* c */) {\n    process(item);\n}",
			pc:       func(f func()) { f() },
		},
		{
			name: "09-hash-foreach-all-inline-comments",
			code: "#foreach /*c1*/ x /*c2*/ in /*c3*/ (/*c4*/ [1, 2] /*c5*/) /*c6*/ {\n" +
				"process(x);\n" +
				"}",
			expected: "#foreach /*c1*/ x /*c2*/ in /*c3*/ (\n" +
				"    /*c4*/ [1, 2] /*c5*/\n" +
				") /*c6*/ {\n" +
				"    process(x);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "10-comments-inside-each-in-expr",
			code: "foreach item in (each /*c1*/ reg /*c2*/ in (/*c3*/ this)) {\n" +
				"process(item);\n" +
				"}",
			expected: "foreach item in (\n" +
				"    each /*c1*/ reg /*c2*/ in (\n" +
				"        /*c3*/ this\n" +
				"    )\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "11-comments-inside-compile-time-array",
			code: "#foreach item in ([a, /* c */ b, c, d]) {\n" +
				"process(item);\n" +
				"}",
			expected: "#foreach item in (\n" +
				"    [a, /* c */ b, c, d]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "12-line-comment-after-lparen",
			code: "foreach item in ( // c\nitems\n) {\nprocess(item);\n}",
			expected: "foreach item in ( // c\n" +
				"    items\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "13-line-comment-inside-array",
			code: "#foreach item in ([\na, // c\nb\n]) {\nprocess(item);\n}",
			expected: "#foreach item in (\n" +
				"    [\n" +
				"        a, // c\n" +
				"        b\n" +
				"    ]\n" +
				") {\n" +
				"    process(item);\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
