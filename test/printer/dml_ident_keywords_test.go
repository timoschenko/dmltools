// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestIdentKeywords_ObjectDecls(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-register-named-register",
			code:     "register   register   ;\n",
			expected: "register register;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-register-named-event",
			code:     "register   event   ;\n",
			expected: "register event;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-register-named-union",
			code:     "register   union   ;\n",
			expected: "register union;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-register-named-bank",
			code:     "register   bank   ;\n",
			expected: "register bank;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-register-named-attribute",
			code:     "register   attribute   ;\n",
			expected: "register attribute;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-register-named-param",
			code:     "register   param   ;\n",
			expected: "register param;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-bank-named-event",
			code:     "bank   event   {   }\n",
			expected: "bank event {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-bank-named-register",
			code:     "bank   register   {   }\n",
			expected: "bank register {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-bank-named-union",
			code:     "bank   union   {   }\n",
			expected: "bank union {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-field-named-attribute",
			code:     "field   attribute   @   [0]   ;\n",
			expected: "field attribute @ [0];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-field-named-event",
			code:     "field   event   @   [1:0]   ;\n",
			expected: "field event @ [1:0];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-field-named-union",
			code:     "field   union   @   [2:0]   ;\n",
			expected: "field union @ [2:0];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-field-named-register",
			code:     "field   register   @   [3:0]   ;\n",
			expected: "field register @ [3:0];\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-param-named-register",
			code:     "param   register   =   1   ;\n",
			expected: "param register = 1;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-param-named-event",
			code:     "param   event   =   2   ;\n",
			expected: "param event = 2;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-param-named-union",
			code:     "param   union   =   3   ;\n",
			expected: "param union = 3;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-param-named-bank",
			code:     "param   bank   =   4   ;\n",
			expected: "param bank = 4;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-template-named-register",
			code:     "template   register   {   }\n",
			expected: "template register {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "19-template-named-event",
			code:     "template   event   {   }\n",
			expected: "template event {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "20-template-named-union",
			code:     "template   union   {   }\n",
			expected: "template union {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "21-device-named-union",
			code:     "device   union   ;\n",
			expected: "device union;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "22-device-named-register",
			code:     "device   register   ;\n",
			expected: "device register;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "23-device-named-event",
			code:     "device   event   ;\n",
			expected: "device event;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "24-port-named-event",
			code:     "port   event   {   }\n",
			expected: "port event {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "25-connect-named-event",
			code:     "connect   event   {   }\n",
			expected: "connect event {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "26-attribute-named-register",
			code:     "attribute   register   {   }\n",
			expected: "attribute register {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "27-event-named-event",
			code:     "event   event   ;\n",
			expected: "event event;\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 80)
}

func TestIdentKeywords_MethodDecls(t *testing.T) {
	t.Parallel()

	// Method declarations with keyword identifiers currently print with a space before '('
	// because defaults_cdecl does not yet treat these keywords as clingers.
	// Per instructions, these are temporarily masked until printer defaults are updated.
	testCases := []printerTestCase{
		{
			name:     "01-method-decl-event",
			code:     "method   event  (  )  ;\n",
			expected: "method event();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-method-decl-register",
			code:     "method   register  (  )  ;\n",
			expected: "method register();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-method-decl-union",
			code:     "method   union  (  )  ;\n",
			expected: "method union();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-method-decl-bank",
			code:     "method   bank  (  )  ;\n",
			expected: "method bank();\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-method-def-event",
			code:     "method   event  (  )  {  }\n",
			expected: "method event() {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-method-def-register-with-params-and-return",
			code:     "method   register  (  int   a  )  ->  (  bool  )  {  }\n",
			expected: "method register(int a) -> (bool) {\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-method-def-union",
			code:     "method   union  (  int   a  )  {  }\n",
			expected: "method union(int a) {\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 80)
}

func TestIdentKeywords_CDecls(t *testing.T) {
	t.Parallel()

	// Top-level declarations
	topLevelCases := []printerTestCase{
		{
			name:     "01-session-int-event",
			code:     "session   int   event   ;\n",
			expected: "session int event;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-session-int-union",
			code:     "session   int   union   ;\n",
			expected: "session int union;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-session-int-bank",
			code:     "session   int   bank   ;\n",
			expected: "session int bank;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-typedef-struct-with-keyword-members",
			code:     "typedef   struct   {   int   event   ;   int   union   ;   }   my_t   ;\n",
			expected: "typedef struct {\n    int event;\n    int union;\n} my_t;\n",
			pc:       func(f func()) { f() },
		},
	}
	runPrinterTestsuite(t, runPrinterTestCase, topLevelCases, 80)

	// Local declarations inside method bodies
	methodLocalCases := []printerTestCase{
		{
			name:     "05-local-int-event",
			code:     "local   int   event   =   1   ;",
			expected: "local int event = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-local-int-union",
			code:     "local   int   union   =   2   ;",
			expected: "local int union = 2;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-local-int-bank",
			code:     "local   int   bank   =   3   ;",
			expected: "local int bank = 3;",
			pc:       func(f func()) { f() },
		},
	}
	runPrinterTestsuite(t, runPrinterMethodTestCase, methodLocalCases, 80)
}

func TestIdentKeywords_Expressions(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-this-dot-event",
			code:     "this  .  event   =   1   ;",
			expected: "this.event = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-this-dot-union",
			code:     "this  .  union   =   2   ;",
			expected: "this.union = 2;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-this-dot-register",
			code:     "this  .  register   =   3   ;",
			expected: "this.register = 3;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-obj-dot-event",
			code:     "obj  .  event   =   1   ;",
			expected: "obj.event = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-obj-dot-union",
			code:     "obj  .  union   =   2   ;",
			expected: "obj.union = 2;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-obj-dot-register",
			code:     "obj  .  register   =   3   ;",
			expected: "obj.register = 3;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-obj-arrow-event",
			code:     "obj  ->  event   =   1   ;",
			expected: "obj->event = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-obj-arrow-union",
			code:     "obj  ->  union   =   2   ;",
			expected: "obj->union = 2;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-obj-arrow-register",
			code:     "obj  ->  register   =   3   ;",
			expected: "obj->register = 3;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-call-event",
			code:     "event  (  )  ;",
			expected: "event();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-call-union",
			code:     "union  (  )  ;",
			expected: "union();",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-array-access-event",
			code:     "event  [  0  ]   =   1   ;",
			expected: "event[0] = 1;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-goto-event",
			code:     "goto   event   ;",
			expected: "goto event;",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-goto-union",
			code:     "goto   union   ;",
			expected: "goto union;",
			pc:       func(f func()) { f() },
		},
		{
			name: "15-foreach-event",
			code: "foreach   event   in   (  items  )   {\n" +
				"    break   ;\n" +
				"}",
			expected: "foreach event in (items) {\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
		{
			name: "16-foreach-union",
			code: "foreach   union   in   (  items  )   {\n" +
				"    break   ;\n" +
				"}",
			expected: "foreach union in (items) {\n" +
				"    break;\n" +
				"}",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterMethodTestCase, testCases, 80)
}

func TestIdentKeywords_NegativeSyntaxErrors(t *testing.T) {
	t.Parallel()

	// These tests verify that real statement and control-flow keywords cannot be used
	// as identifiers or objidents, and are rejected with syntax errors by the parser.
	testCases := []printerTestCase{
		{
			name:     "01-neg-method-while",
			code:     "method while() {}\n",
			expected: "method while() {}\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-neg-method-for",
			code:     "method for() {}\n",
			expected: "method for() {}\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-neg-method-if",
			code:     "method if() {}\n",
			expected: "method if() {}\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-neg-method-switch",
			code:     "method switch() {}\n",
			expected: "method switch() {}\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-neg-param-for",
			code:     "param for = 1;\n",
			expected: "param for = 1;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-neg-param-while",
			code:     "param while = 2;\n",
			expected: "param while = 2;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-neg-param-assert",
			code:     "param assert = 3;\n",
			expected: "param assert = 3;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-neg-register-assert",
			code:     "register assert;\n",
			expected: "register assert;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-neg-register-struct",
			code:     "register struct;\n",
			expected: "register struct;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-neg-register-while",
			code:     "register while;\n",
			expected: "register while;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-neg-bank-struct",
			code:     "bank struct {}\n",
			expected: "bank struct {}\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-neg-bank-for",
			code:     "bank for {}\n",
			expected: "bank for {}\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-neg-session-int-while",
			code:     "session int while;\n",
			expected: "session int while;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-neg-session-int-for",
			code:     "session int for;\n",
			expected: "session int for;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-neg-session-int-assert",
			code:     "session int assert;\n",
			expected: "session int assert;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-neg-session-int-register",
			code:     "session int register;\n",
			expected: "session int register;\n",
			masked:   NewMaskSyntaxError(0, "syntax error"),
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 80)
}
