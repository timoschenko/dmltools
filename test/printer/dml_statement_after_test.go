// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestAfterStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Immediate after: AFTER COLON expression SEMI
		{
			name:     "01-immediate-after",
			code:     "after   :   my_method  (  )   ;",
			expected: "after: my_method();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-immediate-after-with-args",
			code:     "after   :   my_method  (  arg1  ,   arg2  )   ;",
			expected: "after: my_method(arg1, arg2);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-immediate-after-member-access",
			code:     "after   :   this  .  bank_a  .  reg_b  .  on_event  (  arg1  )   ;",
			expected: "after: this.bank_a.reg_b.on_event(arg1);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-immediate-after-indexed-object",
			code:     "after   :   port_arr  [  i  ]  .  raise  (  )   ;",
			expected: "after: port_arr[i].raise();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-immediate-after-hook-send-now",
			code:     "after   :   my_hook  .  send_now  (  1  ,   false  )   ;",
			expected: "after: my_hook.send_now(1, false);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-immediate-after-messy-newlines",
			code:     "after\n:\ncallback\n(\narg1\n,\narg2\n)\n;",
			expected: "after: callback(arg1, arg2);",
			pc:       func(f func()) { f() },
		},

		// Delay with unit: expression IDENTIFIER
		{
			name:     "10-delay-with-unit-seconds",
			code:     "after   0.1   s   :   callback  (  )   ;",
			expected: "after 0.1 s: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-delay-with-unit-ps",
			code:     "after   100   ps   :   callback  (  )   ;",
			expected: "after 100 ps: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-delay-with-unit-cycles",
			code:     "after   5   cycles   :   callback  (  )   ;",
			expected: "after 5 cycles: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-delay-with-hex-literal-cycles",
			code:     "after   0x10   cycles   :   callback  (  )   ;",
			expected: "after 0x10 cycles: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-delay-with-binary-literal-ps",
			code:     "after   0b1010   ps   :   callback  (  )   ;",
			expected: "after 0b1010 ps: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-delay-with-identifier-unit",
			code:     "after   delay_val   s   :   callback  (  1  ,   false  )   ;",
			expected: "after delay_val s: callback(1, false);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-delay-with-member-access-unit",
			code:     "after   this  .  cfg  .  delay   cycles   :   cb  (  )   ;",
			expected: "after this.cfg.delay cycles: cb();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-delay-with-method-call-unit",
			code:     "after   get_delay  (  )   ps   :   callback  (  )   ;",
			expected: "after get_delay() ps: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-delay-messy-newlines",
			code:     "after\n0.1\ns\n:\ncallback\n(\n)\n;",
			expected: "after 0.1 s: callback();",
			pc:       func(f func()) { f() },
		},

		// Delay / Hook expression only: expression
		{
			name:     "20-delay-expression-only",
			code:     "after   10   :   callback  (  )   ;",
			expected: "after 10: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "21-hook-bound-no-binding",
			code:     "after   h   :   callback  (  )   ;",
			expected: "after h: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "22-hook-bound-member-access-no-binding",
			code:     "after   this  .  bank_a  .  my_hook   :   callback  (  )   ;",
			expected: "after this.bank_a.my_hook: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "23-hook-bound-indexed-no-binding",
			code:     "after   hooks  [  i  ]   :   callback  (  i  )   ;",
			expected: "after hooks[i]: callback(i);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "24-hook-bound-2d-indexed-no-binding",
			code:     "after   this  .  hooks  [  i  ]  [  j  ]   :   callback  (  )   ;",
			expected: "after this.hooks[i][j]: callback();",
			pc:       func(f func()) { f() },
		},

		// Hook-bound single identifier: expression ARROW ident
		{
			name:     "30-hook-bound-single-ident",
			code:     "after   h   ->   x   :   callback  (  x  )   ;",
			expected: "after h -> x: callback(x);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "31-hook-bound-indexed-single-ident",
			code:     "after   this  .  hooks  [  i  ]   ->   msg   :   on_msg  (  msg  )   ;",
			expected: "after this.hooks[i] -> msg: on_msg(msg);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "32-hook-bound-single-ident-discard",
			code:     "after   h   ->   _   :   callback  (  )   ;",
			expected: "after h -> _: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name: "33-hook-bound-single-keyword-idents",
			code: "after   h   ->   bank   :   cb  (  bank  )   ;\n" +
				"after   h   ->   field   :   cb  (  field  )   ;\n" +
				"after   h   ->   port   :   cb  (  port  )   ;\n" +
				"after   h   ->   device   :   cb  (  device  )   ;\n" +
				"after   h   ->   data   :   cb  (  data  )   ;\n" +
				"after   h   ->   event   :   cb  (  event  )   ;\n" +
				"after   h   ->   param   :   cb  (  param  )   ;\n" +
				"after   h   ->   info   :   cb  (  info  )   ;\n" +
				"after   h   ->   size   :   cb  (  size  )   ;",
			expected: "after h -> bank: cb(bank);\n" +
				"after h -> field: cb(field);\n" +
				"after h -> port: cb(port);\n" +
				"after h -> device: cb(device);\n" +
				"after h -> data: cb(data);\n" +
				"after h -> event: cb(event);\n" +
				"after h -> param: cb(param);\n" +
				"after h -> info: cb(info);\n" +
				"after h -> size: cb(size);",
			pc: func(f func()) { f() },
		},
		{
			name:     "34-hook-bound-single-ident-messy-newlines",
			code:     "after\nh\n->\nx\n:\ncallback\n(\nx\n)\n;",
			expected: "after h -> x: callback(x);",
			pc:       func(f func()) { f() },
		},

		// Hook-bound tuple: expression ARROW LPAREN (identList)? RPAREN
		{
			name:     "40-hook-bound-empty-tuple",
			code:     "after   h   ->   (  )   :   callback  (  )   ;",
			expected: "after h -> (): callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "41-hook-bound-single-element-tuple",
			code:     "after   h   ->   (  x  )   :   callback  (  x  )   ;",
			expected: "after h -> (x): callback(x);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "42-hook-bound-two-element-tuple",
			code:     "after   h   ->   (  x  ,   y  )   :   callback  (  x  ,   y  )   ;",
			expected: "after h -> (x, y): callback(x, y);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "43-hook-bound-three-element-tuple",
			code:     "after   h   ->   (  x  ,   y  ,   z  )   :   cb  (  x  ,   y  ,   z  )   ;",
			expected: "after h -> (x, y, z): cb(x, y, z);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "44-hook-bound-tuple-repeated-and-unused-params",
			code:     "after   h   ->   (  x  ,   y  )   :   my_callback  (  x  ,   x  ,   false  )   ;",
			expected: "after h -> (x, y): my_callback(x, x, false);",
			lineCap:  50,
			pc:       func(f func()) { f() },
		},
		{
			name:     "45-hook-bound-tuple-keyword-idents",
			code:     "after   h   ->   (  _  ,   data  ,   size  )   :   cb  (  data  ,   size  )   ;",
			expected: "after h -> (_, data, size): cb(data, size);",
			lineCap:  50,
			pc:       func(f func()) { f() },
		},
		{
			name:     "46-hook-bound-tuple-messy-newlines",
			code:     "after\nh\n->\n(\nx\n,\ny\n)\n:\ncallback\n(\nx\n,\ny\n)\n;",
			expected: "after h -> (x, y): callback(x, y);",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAfterStatement_NoStatementWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-immediate-long-callback-name-no-wrap",
			code:     "after   :   very_long_callback_method_name_exceeding_cap  (  )   ;",
			expected: "after: very_long_callback_method_name_exceeding_cap();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-delay-long-identifier-with-unit-no-wrap",
			code:     "after   very_long_delay_constant_exceeding_line_cap   s   :   callback  (  )   ;",
			expected: "after very_long_delay_constant_exceeding_line_cap s: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-delay-long-identifier-with-cycles-no-wrap",
			code:     "after   very_long_delay_constant_exceeding_line_cap   cycles   :   cb  (  )   ;",
			expected: "after very_long_delay_constant_exceeding_line_cap cycles: cb();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-hook-long-identifier-no-binding-no-wrap",
			code:     "after   very_long_hook_identifier_exceeding_line_cap   :   callback  (  )   ;",
			expected: "after very_long_hook_identifier_exceeding_line_cap: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name: "05-hook-long-identifier-with-single-ident-no-wrap",
			code: "after   very_long_hook_identifier_exceeding_line_cap   ->   x   :   " +
				"callback  (  )   ;",
			expected: "after very_long_hook_identifier_exceeding_line_cap -> x: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name: "06-hook-long-bound-ident-no-wrap",
			code: "after   h   ->   very_long_bound_parameter_name_exceeding_cap   :   " +
				"callback  (  )   ;",
			expected: "after h -> very_long_bound_parameter_name_exceeding_cap: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name: "07-hook-long-identifier-with-empty-tuple-no-wrap",
			code: "after   very_long_hook_identifier_exceeding_line_cap   ->   (  )   :   " +
				"callback  (  )   ;",
			expected: "after very_long_hook_identifier_exceeding_line_cap -> (): callback();",
			pc:       func(f func()) { f() },
		},
		{
			name: "08-delay-unparenthesized-binary-expr-no-wrap",
			code: "after   base_delay_value   +   extra_delay_offset   s   :   callback  (  )   ;",
			expected: "after base_delay_value\n" +
				"    + extra_delay_offset s: callback();",
			pc: func(f func()) { f() },
		},
		{
			name: "09-delay-unparenthesized-member-access-no-wrap",
			code: "after   this  .  cfg  .  timeout   cycles   :   callback  (  )   ;",
			expected: "after this.cfg\n" +
				"    .timeout cycles: callback();",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAfterStatement_ConditionWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Level 1: Single Line - parenthesized expression fits within lineCap
		{
			name:     "01-level1-delay-binary-fits",
			code:     "after   (  a   +   b  )   s   :   callback  (  )   ;",
			expected: "after (a + b) s: callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(a + b) s: callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "02-level1-delay-ternary-fits",
			code:     "after   (  c   ?   1   :   2  )   cycles   :   cb  (  )   ;",
			expected: "after (c ? 1 : 2) cycles: cb();",
			masked: MaskIncorrectParenWhitespace(
				"after(c ? 1 : 2) cycles: cb();",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "03-level1-delay-method-call-fits",
			code:     "after   (  get_d  (  x  )  )   ps   :   cb  (  )   ;",
			expected: "after (get_d(x)) ps: cb();",
			masked:   MaskIncorrectParenWhitespace("after(get_d(x)) ps: cb();"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-level1-hook-expr-no-binding-fits",
			code:     "after   (  c   ?   h1   :   h2  )   :   cb  (  )   ;",
			expected: "after (c ? h1 : h2): cb();",
			masked:   MaskIncorrectParenWhitespace("after(c ? h1 : h2): cb();"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-level1-hook-expr-single-ident-fits",
			code:     "after   (  c   ?   h1   :   h2  )   ->   x   :   cb  (  x  )   ;",
			expected: "after (c ? h1 : h2) -> x: cb(x);",
			masked: MaskIncorrectParenWhitespace(
				"after(c ? h1 : h2) -> x: cb(x);",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "06-level1-hook-expr-tuple-fits",
			code:     "after   (  h  )   ->   (  x  ,   y  )   :   cb  (  x  ,   y  )   ;",
			expected: "after (h) -> (x, y): cb(x, y);",
			masked: MaskIncorrectParenWhitespace(
				"after(h) -> (x, y): cb(x, y);",
			),
			pc: func(f func()) { f() },
		},

		// Level 2: Outer Parenthesized Wrap - statement exceeds lineCap, inner expression fits on one line
		{
			name: "10-level2-delay-binary-outer-wrap",
			code: "after (base_delay + extra_delay) s: callback();",
			expected: "after (\n" +
				"    base_delay + extra_delay\n" +
				") s: callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    base_delay + extra_delay\n" +
					") s: callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "11-level2-delay-ternary-outer-wrap",
			code: "after (fast_mode ? FAST_DELAY : SLOW_DELAY) cycles: callback();",
			expected: "after (\n" +
				"    fast_mode ? FAST_DELAY : SLOW_DELAY\n" +
				") cycles: callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    fast_mode ? FAST_DELAY : SLOW_DELAY\n" +
					") cycles: callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "12-level2-delay-method-call-outer-wrap",
			code: "after (compute_delay(dev, mode)) ps: callback();",
			expected: "after (\n" +
				"    compute_delay(dev, mode)\n" +
				") ps: callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    compute_delay(dev, mode)\n" +
					") ps: callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "13-level2-delay-nested-parens-outer-wrap",
			code: "after ((base_a + base_b) * scale_factor) cycles: callback();",
			expected: "after (\n" +
				"    (base_a + base_b) * scale_factor\n" +
				") cycles: callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    (base_a + base_b) * scale_factor\n" +
					") cycles: callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "14-level2-hook-expr-no-binding-outer-wrap",
			code: "after (use_primary ? hook_alpha : hook_beta): callback();",
			expected: "after (\n" +
				"    use_primary ? hook_alpha : hook_beta\n" +
				"): callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    use_primary ? hook_alpha : hook_beta\n" +
					"): callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "15-level2-hook-expr-single-ident-outer-wrap",
			code: "after (use_primary ? hook_alpha : hook_beta) -> msg: callback(msg);",
			expected: "after (\n" +
				"    use_primary ? hook_alpha : hook_beta\n" +
				") -> msg: callback(msg);",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    use_primary ? hook_alpha : hook_beta\n" +
					") -> msg: callback(msg);",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "16-level2-hook-expr-tuple-outer-wrap",
			code: "after (use_primary ? hook_alpha : hook_beta) -> (x, y): cb(x, y);",
			expected: "after (\n" +
				"    use_primary ? hook_alpha : hook_beta\n" +
				") -> (x, y): cb(x, y);",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    use_primary ? hook_alpha : hook_beta\n" +
					") -> (x, y): cb(x, y);",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "17-level2-hook-tuple-params-wrap",
			code: "after h -> (msg_param_a, msg_param_b): callback();",
			expected: "after h -> (\n" +
				"    msg_param_a, msg_param_b\n" +
				"): callback();",
			masked: MaskUncategorizedFormatting(
				"after h -> (\n" +
					"    msg_param_a,\n" +
					"    msg_param_b): callback();",
			),
			pc: func(f func()) { f() },
		},

		// Level 3: Deep Expression Explosion - both statement and inner expression exceed lineCap
		{
			name: "20-level3-delay-arithmetic-explosion",
			code: "after (very_long_base_delay + very_long_extra_delay) s: callback();",
			expected: "after (\n" +
				"    very_long_base_delay\n" +
				"    + very_long_extra_delay\n" +
				") s: callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    very_long_base_delay\n" +
					"    + very_long_extra_delay\n" +
					") s: callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "21-level3-delay-multi-operand-arithmetic-explosion",
			code: "after (delay_part_alpha + delay_part_beta + delay_part_gamma) cycles: cb();",
			expected: "after (\n" +
				"    delay_part_alpha\n" +
				"    + delay_part_beta\n" +
				"    + delay_part_gamma\n" +
				") cycles: cb();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    delay_part_alpha\n" +
					"    + delay_part_beta\n" +
					"    + delay_part_gamma\n" +
					") cycles: cb();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "22-level3-delay-ternary-explosion",
			code: "after (use_extended_timeout ? EXTENDED_DELAY_CYCLES : DEFAULT_DELAY_CYCLES) " +
				"cycles: callback();",
			expected: "after (\n" +
				"    use_extended_timeout\n" +
				"        ? EXTENDED_DELAY_CYCLES\n" +
				"        : DEFAULT_DELAY_CYCLES\n" +
				") cycles: callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    use_extended_timeout\n" +
					"        ? EXTENDED_DELAY_CYCLES\n" +
					"        : DEFAULT_DELAY_CYCLES\n" +
					") cycles: callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "23-level3-delay-nested-parens-explosion",
			code: "after ((very_long_base_delay_a + very_long_base_delay_b) * scale_factor) " +
				"cycles: callback();",
			expected: "after (\n" +
				"    (\n" +
				"        very_long_base_delay_a\n" +
				"        + very_long_base_delay_b\n" +
				"    )\n" +
				"    * scale_factor\n" +
				") cycles: callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    (\n" +
					"        very_long_base_delay_a\n" +
					"        + very_long_base_delay_b\n" +
					"    )\n" +
					"    * scale_factor\n" +
					") cycles: callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "24-level3-delay-method-call-args-explosion",
			code: "after (calculate_timeout_cycles(device_context, operation_mode)) " +
				"cycles: callback();",
			expected: "after (\n" +
				"    calculate_timeout_cycles(\n" +
				"        device_context, operation_mode\n" +
				"    )\n" +
				") cycles: callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    calculate_timeout_cycles(\n" +
					"        device_context, operation_mode\n" +
					"    )\n" +
					") cycles: callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "25-level3-hook-expr-no-binding-ternary-explosion",
			code: "after (use_primary_channel ? primary_event_hook : secondary_event_hook): " +
				"callback();",
			expected: "after (\n" +
				"    use_primary_channel\n" +
				"        ? primary_event_hook\n" +
				"        : secondary_event_hook\n" +
				"): callback();",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    use_primary_channel\n" +
					"        ? primary_event_hook\n" +
					"        : secondary_event_hook\n" +
					"): callback();",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "26-level3-hook-expr-single-ident-ternary-explosion",
			code: "after (use_primary_channel ? primary_event_hook : secondary_event_hook) " +
				"-> msg: callback(msg);",
			expected: "after (\n" +
				"    use_primary_channel\n" +
				"        ? primary_event_hook\n" +
				"        : secondary_event_hook\n" +
				") -> msg: callback(msg);",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    use_primary_channel\n" +
					"        ? primary_event_hook\n" +
					"        : secondary_event_hook\n" +
					") -> msg: callback(msg);",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "27-level3-hook-expr-tuple-ternary-explosion",
			code: "after (use_primary_channel ? primary_event_hook : secondary_event_hook) " +
				"-> (x, y): cb(x, y);",
			expected: "after (\n" +
				"    use_primary_channel\n" +
				"        ? primary_event_hook\n" +
				"        : secondary_event_hook\n" +
				") -> (x, y): cb(x, y);",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    use_primary_channel\n" +
					"        ? primary_event_hook\n" +
					"        : secondary_event_hook\n" +
					") -> (x, y): cb(x, y);",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "28-level3-hook-tuple-params-explosion",
			code: "after h -> (very_long_msg_param_a, very_long_msg_param_b, " +
				"very_long_msg_param_c): callback();",
			expected: "after h -> (\n" +
				"    very_long_msg_param_a,\n" +
				"    very_long_msg_param_b,\n" +
				"    very_long_msg_param_c\n" +
				"): callback();",
			masked: MaskUncategorizedFormatting(
				"after h -> (\n" +
					"    very_long_msg_param_a,\n" +
					"    very_long_msg_param_b,\n" +
					"    very_long_msg_param_c): callback();",
			),
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAfterStatement_CallbackAndCombinedWrapping(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		// Callback wrapping across all after forms
		{
			name: "01-immediate-after-long-callback-args",
			code: "after: callback(very_long_argument_name_one, very_long_argument_name_two);",
			expected: "after: callback(\n" +
				"    very_long_argument_name_one,\n" +
				"    very_long_argument_name_two\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "02-delay-unit-long-callback-args",
			code: "after 0.1 s: callback(very_long_argument_name_one, very_long_argument_name_two);",
			expected: "after 0.1 s: callback(\n" +
				"    very_long_argument_name_one,\n" +
				"    very_long_argument_name_two\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "03-hook-no-binding-long-callback-args",
			code: "after h: callback(very_long_argument_name_one, very_long_argument_name_two);",
			expected: "after h: callback(\n" +
				"    very_long_argument_name_one,\n" +
				"    very_long_argument_name_two\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "04-hook-single-ident-long-callback-args",
			code: "after h -> x: callback(x, very_long_argument_name_one, " +
				"very_long_argument_name_two);",
			expected: "after h -> x: callback(\n" +
				"    x,\n" +
				"    very_long_argument_name_one,\n" +
				"    very_long_argument_name_two\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "05-hook-tuple-long-callback-args",
			code: "after h -> (x, y): callback(x, y, very_long_extra_argument_name);",
			expected: "after h -> (x, y): callback(\n" +
				"    x, y, very_long_extra_argument_name\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "06-callback-with-exploded-binary-arg",
			code: "after 0.1 s: callback(very_long_operand_alpha + very_long_operand_beta, flag);",
			expected: "after 0.1 s: callback(\n" +
				"    very_long_operand_alpha\n" +
				"    + very_long_operand_beta,\n" +
				"    flag\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "07-callback-with-nested-call-arg",
			code: "after: on_complete(compute_result(source_bank_id, target_register_id), false);",
			expected: "after: on_complete(\n" +
				"    compute_result(\n" +
				"        source_bank_id,\n" +
				"        target_register_id\n" +
				"    ),\n" +
				"    false\n" +
				");",
			pc: func(f func()) { f() },
		},
		{
			name: "08-long-hook-header-wraps-single-callback-arg",
			code: "after very_long_hook_identifier_exceeding_line_cap -> x: callback(x);",
			expected: "after very_long_hook_identifier_exceeding_line_cap -> x: callback(\n" +
				"    x\n" +
				");",
			pc: func(f func()) { f() },
		},

		// Combined wrapping: both parenthesized delay/hook condition and callback arguments wrap
		{
			name: "10-combined-level2-delay-and-callback-args-wrap",
			code: "after (base_delay + extra_delay) s: callback(" +
				"very_long_argument_name_one, very_long_argument_name_two);",
			expected: "after (\n" +
				"    base_delay + extra_delay\n" +
				") s: callback(\n" +
				"    very_long_argument_name_one,\n" +
				"    very_long_argument_name_two\n" +
				");",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    base_delay + extra_delay\n" +
					") s: callback(\n" +
					"    very_long_argument_name_one,\n" +
					"    very_long_argument_name_two\n" +
					");",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "11-combined-level3-delay-explosion-and-callback-args-wrap",
			code: "after (very_long_base_delay + very_long_extra_delay) cycles: callback(" +
				"very_long_argument_name_one, very_long_argument_name_two);",
			expected: "after (\n" +
				"    very_long_base_delay\n" +
				"    + very_long_extra_delay\n" +
				") cycles: callback(\n" +
				"    very_long_argument_name_one,\n" +
				"    very_long_argument_name_two\n" +
				");",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    very_long_base_delay\n" +
					"    + very_long_extra_delay\n" +
					") cycles: callback(\n" +
					"    very_long_argument_name_one,\n" +
					"    very_long_argument_name_two\n" +
					");",
			),
			pc: func(f func()) { f() },
		},
		{
			name: "12-combined-hook-expr-and-callback-args-wrap",
			code: "after (use_primary ? hook_alpha : hook_beta) -> (x, y): callback(" +
				"x, y, very_long_extra_argument_name);",
			expected: "after (\n" +
				"    use_primary ? hook_alpha : hook_beta\n" +
				") -> (x, y): callback(\n" +
				"    x, y, very_long_extra_argument_name\n" +
				");",
			masked: MaskIncorrectParenWhitespace(
				"after(\n" +
					"    use_primary ? hook_alpha : hook_beta\n" +
					") -> (x, y): callback(\n" +
					"    x, y, very_long_extra_argument_name\n" +
					");",
			),
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}

func TestAfterStatement_Comments(t *testing.T) {
	t.Parallel()

	const lineCap = 40
	testCases := []printerTestCase{
		{
			name:     "01-leading-line-comment",
			code:     "// comment\nafter: callback();",
			expected: "// comment\nafter: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-leading-block-comment",
			code:     "/* comment */   after   :   callback  (  )   ;",
			expected: "/* comment */ after: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-comment-between-after-and-colon",
			code:     "after /* c */ : callback();",
			expected: "after /* c */: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-comment-between-after-and-delay",
			code:     "after /* c */ 0.1 s: callback();",
			expected: "after /* c */ 0.1 s: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-comment-between-delay-expr-and-unit",
			code:     "after 0.1 /* c */ s: callback();",
			expected: "after 0.1 /* c */ s: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-comment-between-delay-unit-and-colon",
			code:     "after 0.1 s /* c */ : callback();",
			expected: "after 0.1 s /* c */: callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-comment-after-colon-immediate",
			code:     "after: /* c */ callback();",
			expected: "after: /* c */ callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-comment-after-colon-delay",
			code:     "after 5 cycles: /* c */ callback();",
			expected: "after 5 cycles: /* c */ callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-comment-before-arrow",
			code:     "after h /* c */ -> (x): callback(x);",
			expected: "after h /* c */ -> (x): callback(x);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-comment-after-arrow-single-ident",
			code:     "after h -> /* c */ x: callback(x);",
			expected: "after h -> /* c */ x: callback(x);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-comment-after-arrow-tuple",
			code:     "after h -> /* c */ (x, y): cb(x, y);",
			expected: "after h -> /* c */ (x, y): cb(x, y);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-comment-inside-empty-tuple",
			code:     "after h -> ( /* c */ ): callback();",
			expected: "after h -> (/* c */): callback();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-comment-inside-tuple-params",
			code:     "after h -> (/* c1 */ x /* c2 */, /* c3 */ y /* c4 */): cb(x, y);",
			expected: "after h -> (/* c1 */ x /* c2 */, /* c3 */ y /* c4 */): cb(x, y);",
			lineCap:  70,
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-comment-between-tuple-and-colon",
			code:     "after h -> (x) /* c */ : callback(x);",
			expected: "after h -> (x) /* c */: callback(x);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-comment-between-single-ident-and-colon",
			code:     "after h -> x /* c */ : callback(x);",
			expected: "after h -> x /* c */: callback(x);",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-comment-inside-callback-args",
			code:     "after: cb(/* c1 */ a /* c2 */, /* c3 */ b /* c4 */);",
			expected: "after: cb(/* c1 */ a /* c2 */, /* c3 */ b /* c4 */);",
			lineCap:  60,
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-block-comment-before-semi",
			code:     "after: callback() /* c */ ;",
			expected: "after: callback() /* c */;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-line-comment-before-semi",
			code:     "after: callback() // c\n;",
			expected: "after: callback() // c\n;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "19-trailing-block-comment",
			code:     "after: callback(); /* done */",
			expected: "after: callback(); /* done */",
			pc:       func(f func()) { f() },
		},
		{
			name:     "20-trailing-line-comment",
			code:     "after: callback(); // done",
			expected: "after: callback(); // done",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, lineCap)
}
