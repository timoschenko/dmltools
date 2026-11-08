// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintLogGroupStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 30
	testCases := []printerTestCase{
		{
			name:     "01-empty",
			code:     "loggroup Foo;\n",
			expected: "loggroup Foo;\n",
			pc:       func(f func()) { f() },
		},
		// Whitespace and new lines formatting test case
		{
			name:     "02-messy-whitespace",
			code:     "    loggroup    Foo   ;\n",
			expected: "loggroup Foo;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-messy-newlines",
			code:     "\n\n    loggroup\n    \nFoo\n  ;\n\n",
			expected: "loggroup Foo;\n",
			pc:       func(f func()) { f() },
		},
		// Comments test cases
		// {  // FIXME: BUG dml header is corrupted by a comment
		// 	name:     "04-normal-before",
		// 	code:     "// comment\nloggroup Foo;\n",
		// 	expected: "// comment\nloggroup Foo;\n",
		// 	pc:       func(f func()) { f() },
		// },
		{
			name:     "05-normal-sameline-tail",
			code:     "loggroup Foo; // comment\n",
			expected: "loggroup Foo; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-normal-sameline-very-long",
			code:     "loggroup     Foo;     // This is very long comment\n",
			expected: "loggroup Foo; // This is very long comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-normal-after",
			code:     "loggroup Foo;\n// comment\n",
			expected: "loggroup Foo;\n// comment\n",
			pc:       func(f func()) { f() },
		},
		// { // FIXME: BUG dml header is corrupted by a comment
		// 	name:     "08-multiline-before",
		// 	code:     "    /*\n * comment\n */  \n   loggroup    Foo;  \n",
		// 	expected: "/*\n * comment\n */\nloggroup Foo;\n",
		// 	pc:       func(f func()) { f() },
		// },
		{
			name:     "09-multiline-sameline-tail",
			code:     "  loggroup    Foo;     /* comment */   \n",
			expected: "loggroup Foo; /* comment */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-multiline-sameline-tail-very-long",
			code:     "   loggroup    Foo;     /* This is very long comment */   \n",
			expected: "loggroup Foo; /* This is very long comment */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-multiline-sameline-tail-new-line",
			code:     "    loggroup     Foo ;    /*\n * comment\n */    \n",
			expected: "loggroup Foo; /*\n * comment\n */\n",
			pc:       func(f func()) { f() },
		},

		{
			name:     "12-multiline-sameline-tail-very-long-new-line",
			code:     "loggroup Foo; /* \n This is very long comment\n * comment\n */\n",
			expected: "loggroup Foo; /*\nThis is very long comment\n* comment\n*/\n",
			masked: MaskMultilineBlockCommentIndentation(
				"loggroup Foo; /*\n This is very long comment\n * comment\n */\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "13-multiline-after",
			code:     "loggroup Foo;\n/*\n * comment\n */\n",
			expected: "loggroup Foo;\n/*\n * comment\n */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-multiline-inline",
			code:     "   loggroup   /* comment */    Foo   ;   \n",
			expected: "loggroup /* comment */ Foo;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-multiline-inline-semi",
			code:     "   loggroup   Foo   /* comment */   ;\n",
			expected: "loggroup Foo /* comment */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "016-three-statements",
			code:     "loggroup Foo;\nloggroup Bar;\nloggroup Baz;\n",
			expected: "loggroup Foo;\nloggroup Bar;\nloggroup Baz;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "017-three-statements-messy",
			code:     "loggroup Foo;\n\n\nloggroup Bar;\nloggroup Baz;\n",
			expected: "loggroup Foo;\n\n\nloggroup Bar;\nloggroup Baz;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "018-grouping-3-statements-messy-comments",
			code:     "loggroup Fo; // Fo\n\n\nloggroup Baar; /* Baar */\nloggroup Baz; // Baz\n",
			expected: "loggroup Fo; // Fo\n\n\nloggroup Baar; /* Baar */\nloggroup Baz; // Baz\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintImportStatement(t *testing.T) {
	t.Parallel()

	//
	// IMPORTANT: This test using "`" instead of "\"" for better reading and editing
	//

	const lineCap = 30
	testCases := []printerTestCase{
		//
		// IMPORTANT: This test using "`" instead of "\"" for better reading and editing
		//
		{
			name:     "01-empty",
			code:     "import `foo.dml`;\n",
			expected: "import `foo.dml`;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-messy-whitespace",
			code:     "    import    `foo.dml`   ;\n",
			expected: "import `foo.dml`;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-messy-newlines",
			code:     "\n\n    import\n    \n`foo.dml`\n  ;\n\n",
			expected: "import `foo.dml`;\n",
			pc:       func(f func()) { f() },
		},
		// Comments test cases
		// {  // FIXME: BUG dml header is corrupted by a comment
		// 	name:     "04-normal-before",
		// 	code:     "// comment    \n   import   `foo.dml`   ;   \n",
		// 	expected: "// comment\nimport `foo.dml`;\n",
		// 	pc:       func(f func()) { f() },
		// },
		{
			name:     "05-normal-sameline-tail",
			code:     "   import   `foo.dml`   ; // comment\n",
			expected: "import `foo.dml`; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-normal-sameline-very-long",
			code:     "   import    `foo.dml`   ;    // This is very long comment\n",
			expected: "import `foo.dml`; // This is very long comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-normal-after",
			code:     "import `foo.dml`;\n// comment\n",
			expected: "import `foo.dml`;\n// comment\n",
			pc:       func(f func()) { f() },
		},
		// { // FIXME: BUG dml header is corrupted by a comment
		// 	name:     "08-multiline-before",
		// 	code:     "    /*\n * comment\n */   \n  import `foo.dml`;\n",
		// 	expected: "/*\n * comment\n */\nimport `foo.dml`;\n",
		// 	pc:       func(f func()) { f() },
		// },
		{
			name:     "09-multiline-sameline-tail",
			code:     "import `foo.dml`  ;    /* cmnt */  \n",
			expected: "import `foo.dml`; /* cmnt */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-multiline-sameline-tail-very-long",
			code:     "import `foo.dml`;   /* This is very long comment */\n",
			expected: "import `foo.dml`; /* This is very long comment */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-multiline-sameline-tail-new-line",
			code:     "import `foo.dml`;    /*\n * comment\n */\n",
			expected: "import `foo.dml`; /*\n * comment\n */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-multiline-sameline-tail-very-long-new-line",
			code:     "import `foo.dml`; /* \n This is very long comment\n * comment\n */\n",
			expected: "import `foo.dml`; /*\nThis is very long comment\n* comment\n*/\n",
			masked: MaskMultilineBlockCommentIndentation(
				"import `foo.dml`; /*\n This is very long comment\n * comment\n */\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "13-multiline-after",
			code:     "import `foo.dml`;\n/*\n * comment\n */\n",
			expected: "import `foo.dml`;\n/*\n * comment\n */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-multiline-inline",
			code:     "import/* cmnt */`foo.dml`;\n",
			expected: "import /* cmnt */ `foo.dml`;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-multiline-inline-semi",
			code:     "import `foo.dml`/* cmnt */;\n",
			expected: "import `foo.dml` /* cmnt */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "016-three-statements",
			code:     "import `foo.dml`;\nimport `bar.dml`;\nimport `baz.dml`;\n",
			expected: "import `foo.dml`;\nimport `bar.dml`;\nimport `baz.dml`;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "017-three-statements-messy",
			code:     "import `foo.dml`;\n\n\nimport `bar.dml`;\nimport `baz.dml`;\n",
			expected: "import `foo.dml`;\n\n\nimport `bar.dml`;\nimport `baz.dml`;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "018-grouping-3-statements-messy-comments",
			code:     "import `Fo.dml`; // Fo\n\n\nimport `Baar.dml`; /* Baar */\nimport `Baz.dml`; // Baz\n",
			expected: "import `Fo.dml`; // Fo\n\n\nimport `Baar.dml`; /* Baar */\nimport `Baz.dml`; // Baz\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestIsTemplateStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 30
	testCases := []printerTestCase{
		{
			name:     "001-empty",
			code:     "is Foo;\n",
			expected: "is Foo;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-messy-whitespace",
			code:     "    is    Foo   ;\n",
			expected: "is Foo;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "003-messy-newlines",
			code:     "\n\n    is\n    \nFoo\n  ;\n\n",
			expected: "is Foo;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "101-tuple-exact",
			code:     "is (foo, bar, baz);\n",
			expected: "is (foo, bar, baz);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "101-tuple-keep-newline",
			code:     "is (\n    foo, bar, baz\n);\n",
			expected: "is (\n    foo, bar, baz\n);\n",
			masked:   MaskInputNewlinePreservation("is (foo, bar, baz);\n"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "102-tuple-keep-multiline",
			code:     "is (\n    foo,\n    bar,\n    baz\n);\n",
			expected: "is (\n    foo,\n    bar,\n    baz\n);\n",
			masked:   MaskInputNewlinePreservation("is (foo, bar, baz);\n"),
			pc:       func(f func()) { f() },
		},
		{
			name:     "200-tuple-long",
			code:     "is (foo, bar, baz, expect__30);\n",
			expected: "is (\n    foo, bar, baz, expect__30\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "300-tuple-very-long",
			code:     "is (foo123, bar456, baz789, expect_29);\n",
			expected: "is (\n    foo123,\n    bar456,\n    baz789,\n    expect_29\n);\n",
			pc:       func(f func()) { f() },
		},
		// { // TODO: dml header is corrupted
		// 	name:     "407-is-header-comments-01",
		// 	code:     "/* l */ is /* r */ (// comment\nt1) {}",
		// 	expected: "is (\n    // comment\n    t1\n) {}\n",
		// 	pc:       func(f func()) { f() },
		// },
		{
			name:     "408-is-header-comments-02",
			code:     "is (t1 // comment\n);",
			expected: "is (\n    t1 // comment\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "409-is-header-comments-03",
			code:     "is (t1, // comment\nt2);",
			expected: "is (\n    t1, // comment\n    t2\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "410-is-header-comments-04",
			code:     "is (t1, t2 /* comment */);",
			expected: "is (t1, t2 /* comment */);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "411-is-header-comments-05",
			code:     "is (t1, t2 /*\n * comment\n */);",
			expected: "is (\n    t1,\n    t2 /*\n * comment\n */\n);\n",
			masked: MaskMissingWrappingForMultilineBlockCommentsInIs(
				"is (t1, t2 /*\n * comment\n */);\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "412-is-header-comments-06-real",
			code:     "is (t1, t2 /* comment */, t3);",
			expected: "is (t1, t2 /* comment */, t3);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "412-is-header-comments-06",
			code:     "is (t1, t2 /* commment */, t3);",
			expected: "is (\n    t1, t2 /* commment */, t3\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "413-is-header-comments-07",
			code:     "is (t1, t2 /* comment */, t3 /*\n * comment\n */);",
			expected: "is (\n    t1,\n    t2 /* comment */,\n    t3 /*\n * comment\n */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "414-is-header-comments-08",
			code:     "is (t1, t2, t3 // comment\n);",
			expected: "is (\n    t1,\n    t2,\n    t3 // comment\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "415-is-header-comments-09-real",
			code:     "is (t1, t2, t3 /* comment */);",
			expected: "is (t1, t2, t3 /* comment */);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "415-is-header-comments-09",
			code:     "is (t1, t2, t3 /* commment */);",
			expected: "is (\n    t1, t2, t3 /* commment */\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "416-is-header-comments-10",
			code:     "is (t1, t2, t3\n// comment\n);",
			expected: "is (\n    t1, t2, t3 // comment\n);\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "417-is-header-comments-11",
			code:     "is (t1, t2, t3\n/* comment */\n);",
			expected: "is (\n    t1, t2, t3 /* comment */\n);\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestInEachStatement(t *testing.T) {
	t.Parallel()

	const lineCap = 30
	testCases := []printerTestCase{
		{
			name:     "001-empty",
			code:     "in each Foo {}\n",
			expected: "in each Foo {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-messy-whitespace",
			code:     "    in   each    Foo   {   }\n",
			expected: "in each Foo {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "003-messy-newlines",
			code:     "\n\n in \n each \n Foo \n { \n }\n\n",
			expected: "in each Foo {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "101-tuple-exact",
			code:     "in each (foo, bar, baz) {}\n",
			expected: "in each (foo, bar, baz) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "102-tuple-keep-newline",
			code:     "in each (\n    foo, bar, baz\n) {}\n",
			expected: "in each (\n    foo, bar, baz\n) {}\n",
			masked: MaskInputNewlinePreservation(
				"in each (foo, bar, baz) {}\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "103-tuple-keep-multiline",
			code:     "in each (\n    foo,\n    bar,\n    baz\n) {}\n",
			expected: "in each (\n    foo,\n    bar,\n    baz\n) {}\n",
			masked: MaskInputNewlinePreservation(
				"in each (foo, bar, baz) {}\n",
			),
			pc: func(f func()) { f() },
		},
		{
			name:     "200-tuple-long",
			code:     "in each (foo, bar, baz, expect__30) {}\n",
			expected: "in each (\n    foo, bar, baz, expect__30\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "300-tuple-very-long",
			code:     "in each (foo123, bar456, baz789, expect_29) {}\n",
			expected: "in each (\n    foo123,\n    bar456,\n    baz789,\n    expect_29\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "401-comments-before",
			code:     "// comment\nin each Foo {}\n",
			expected: "// comment\nin each Foo {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "402-comments-tuple-inside",
			code:     "in each (t1 /* comment */, t2) {}\n",
			expected: "in each (\n    t1 /* comment */, t2\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "403-comments-head-inline",
			code:     "in /* comment */ each Foo {}\n",
			expected: "in /* comment */ each Foo {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "404-comments-before-body",
			code:     "in each Foo /* comment */ {}\n",
			expected: "in each Foo /* comment */ {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "405-comments-inside-body",
			code:     "in each Foo {\n    // comment\n}\n",
			expected: "in each Foo {\n    // comment\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "406-comments-tuple-newline",
			code:     "in each (t1, // comment\nt2) {}\n",
			expected: "in each (\n    t1, // comment\n    t2\n) {}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "501-populated-body",
			code:     "in each Foo {\n    param p;\n}\n",
			expected: "in each Foo {\n    param p;\n}\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintToplevelIf(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		{
			name:     "01-toplevel-if-basic",
			code:     "dml 1.4;\n\n#if (cond) {\nparam p1;\n}\n",
			expected: "dml 1.4;\n\n\n#if (cond) {\n    param p1;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-toplevel-if-else",
			code:     "dml 1.4;\n\n#if (cond) {\nparam p1;\n} #else {\nparam p2;\n}\n",
			expected: "dml 1.4;\n\n\n#if (cond) {\n    param p1;\n} #else {\n    param p2;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-toplevel-if-else-if-chain",
			code:     "dml 1.4;\n\n#if (c1) {\nparam p1;\n} #else #if (c2) {\nparam p2;\n} #else {\nparam p3;\n}\n",
			expected: "dml 1.4;\n\n\n#if (c1) {\n    param p1;\n} #else #if (c2) {\n    param p2;\n} #else {\n    param p3;\n}\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-toplevel-nested-if-unbraced",
			code:     "dml 1.4;\n\n#if (c1)\n#if (c2)\nparam p;\n",
			expected: "dml 1.4;\n\n\n#if (c1) #if (c2) param p;\n",
			pc:       func(f func()) { f() },
		},
	}

	lineCap := 40
	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
