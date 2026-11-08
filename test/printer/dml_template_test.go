// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintTemplateStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 30
	testCases := []printerTestCase{
		{
			name:     "01-empty",
			code:     "template t {}",
			expected: "template t {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-empty-with-space",
			code:     "    template   t   { \n\n }",
			expected: "template t {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-is-header",
			code:     "template t is (t1) {}",
			expected: "template t is (t1) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-is-header",
			code:     "template very_long_template_name is (t1) {}",
			expected: "template very_long_template_name is (\n    t1\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-is-header-messy",
			code:     "template   t   is   (  t1  )   {  }",
			expected: "template t is (t1) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-is-header-multiple",
			code:     "template t is (t1, t2, t3) {}",
			expected: "template t is (t1, t2, t3) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-is-header-multiple-messy",
			code:     "template   t   is   (  t1  ,  t2  ,  t3  )   {  }",
			expected: "template t is (t1, t2, t3) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-is-header-comments-01",
			code:     "template /* tt */ ta /* l */ is /* r */ (// comment\nt1) {}",
			expected: "template /* tt */ ta /* l */ is /* r */ (\n    // comment\n    t1\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-is-header-comments-02",
			code:     "template t is (t1 // comment\n) {}",
			expected: "template t is (\n    t1 // comment\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-is-header-comments-03",
			code:     "template t is (t1, // comment\nt2) {}",
			expected: "template t is (\n    t1, // comment\n    t2\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-is-header-comments-04",
			code:     "template t is (t1, t2 /* comment */) {}",
			expected: "template t is (\n    t1, t2 /* comment */\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-is-header-comments-05",
			code:     "template t is (t1, t2 /*\n * comment\n */) {}",
			expected: "template t is (\n    t1,\n    t2 /*\n * comment\n */\n) {}\n",
			masked: MaskMissingWrappingForMultilineBlockCommentsInIs(
				"template t is (t1, t2 /*\n * comment\n */) {}\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "13-is-header-comments-06",
			code:     "template t is (t1, t2 /* comment */, t3) {}",
			expected: "template t is (\n    t1, t2 /* comment */, t3\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-is-header-comments-07",
			code:     "template t is (t1, t2 /* comment */, t3 /*\n * comment\n */) {}",
			expected: "template t is (\n    t1,\n    t2 /* comment */,\n    t3 /*\n * comment\n */\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-is-header-comments-08",
			code:     "template t is (t1, t2, t3 // comment\n) {}",
			expected: "template t is (\n    t1,\n    t2,\n    t3 // comment\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "16-is-header-comments-09",
			code:     "template t is (t1, t2, t3 /* comment */) {}",
			expected: "template t is (\n    t1, t2, t3 /* comment */\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "17-is-header-comments-10",
			code:     "template t is (t1, t2, t3\n// comment\n) {}",
			expected: "template t is (\n    t1, t2, t3 // comment\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "18-is-header-comments-11",
			code:     "template t is (t1, t2, t3\n/* comment */\n) {}",
			expected: "template t is (\n    t1, t2, t3 /* comment */\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "20-body-param",
			code:     "template t { param p; }",
			expected: "template t {\n    param p;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "21-body-param-messy",
			code:     "template   t   {   param    p;   }",
			expected: "template t {\n    param p;\n}\n",
			pc:       func(f func()) { f() },
		},
		// {  // FIXME: BUG dml header is corrupted by a comment
		// 	name:     "22-body-param-comments",
		// 	code:     "// before\ntemplate t { // on template line\n    // in body\n    param p; // on param line\n} // after\n",
		// 	expected: "// before\ntemplate t { // on template line\n    // in body\n    param p; // on param line\n} // after\n",
		// 	pc:       func(f func()) { f() },
		// },
		{
			name:     "23-body-shared-hook",
			code:     "template t { shared hook() h; }",
			expected: "template t {\n    shared hook() h;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "24-body-multiple-same-kind",
			code:     "template t {\n    param p1;\n    param p2;\n}",
			expected: "template t {\n    param p1;\n    param p2;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "25-body-multiple-different-kind",
			code:     "template t {\n    param p;\n    shared hook() h;\n}",
			expected: "template t {\n    param p;\n\n    shared hook() h;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "26-body-multiple-different-kind-messy",
			code:     "template t { param p; shared hook() h; }",
			expected: "template t {\n    param p;\n\n    shared hook() h;\n}\n",
			pc:       func(f func()) { f() },
		},
		// { // FIXME: BUG Code has missing dml header
		// 	name:     "27-comments-before-template",
		// 	code:     "/* comment */ template t {}",
		// 	expected: "/* comment */\ntemplate t {}\n",
		// 	pc:       func(f func()) { f() },
		// },
		{
			name:     "28-comments-between-template-objident",
			code:     "template /* comment */ t {}",
			expected: "template /* comment */ t {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "29-comments-between-objident-istemplateExpr",
			code:     "template t /* c */ is (t1) {}",
			expected: "template t /* c */ is (t1) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "30-comments-between-objident-lbrace",
			code:     "template t /* comment */ {}",
			expected: "template t /* comment */ {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "31-comments-between-istemplateExpr-lbrace",
			code:     "template t is (t1) /* c */ {}",
			expected: "template t is (t1) /* c */ {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "32-comments-after-lbrace",
			code:     "template t { /* comment */ }",
			expected: "template t {\n    /* comment */\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "33-comments-before-rbrace",
			code:     "template t { param p; /* comment */ }",
			expected: "template t {\n    param p; /* comment */\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "34-comments-after-rbrace",
			code:     "template t {} /* comment */",
			expected: "template t {} /* comment */\n",
			pc:       func(f func()) { f() },
		},
		// { // FIXME: BUG Code has missing dml header
		// 	name:     "35-line-comment-before-template",
		// 	code:     "// comment\ntemplate t {}",
		// 	expected: "// comment\ntemplate t {}\n",
		// 	pc:       func(f func()) { f() },
		// },
		{
			name:     "36-line-comment-between-template-objident",
			code:     "template // comment\nt {}",
			expected: "template // comment\nt {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "37-line-comment-between-objident-lbrace",
			code:     "template t // comment\n{}",
			expected: "template t // comment\n{}\n",
			pc:       func(f func()) { f() },
		},
		// { // FIXME: Syntax error - mismatched input '/' expecting <EOF>
		// 	name:     "38-line-comment-after-rbrace",
		// 	code:     "template t {}\n// comment",
		// 	expected: "template t {}\n// comment\n",
		// 	pc:       func(f func()) { f() },
		// },
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
