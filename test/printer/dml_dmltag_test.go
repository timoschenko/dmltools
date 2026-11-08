// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestPrintDmlTag(t *testing.T) {
	t.Parallel()

	const lineCap = 30
	const EXPECTED = "dml 1.4;\n"
	testCases := []printerTestCase{
		{
			name:     "01-1:1",
			code:     EXPECTED,
			expected: EXPECTED,
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-wrong-ident",
			code:     "    dml 1.4;\n",
			expected: EXPECTED,
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-extra-newlines",
			code:     "\n\ndml 1.4;\n\n\n",
			expected: EXPECTED,
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-extra-whitespace",
			code:     "dml    1.4   ;",
			expected: EXPECTED,
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-extra-new-lines-inside",
			code:     "dml  \n  1.4  \n\n ;\n",
			expected: EXPECTED,
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintDmlTag_comments(t *testing.T) {
	t.Parallel()

	const lineCap = 30
	testCases := []printerTestCase{
		{
			name:     "01-normal-before",
			code:     "   // comment\ndml 1.4;\n",
			expected: "// comment\ndml 1.4;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-normal-sameline-tail",
			code:     "dml 1.4;       // comment\n",
			expected: "dml 1.4; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-normal-after",
			code:     "dml 1.4;\n         // comment\n",
			expected: "dml 1.4;\n// comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-multiline-before",
			code:     "    /*\n * comment\n */     \n   dml 1.4;\n",
			expected: "/*\n * comment\n */\ndml 1.4;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-multiline-sameline-tail",
			code:     "dml 1.4;    /* comment */\n",
			expected: "dml 1.4; /* comment */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-multiline-after",
			code:     "   dml 1.4;  \n   /*\n * comment\n */    \n",
			expected: "dml 1.4;\n/*\n * comment\n */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-multiline-inline",
			code:     "dml     /* comment */    1.4;\n",
			expected: "dml /* comment */ 1.4;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-multiline-inline",
			code:     "dml    /*\n * comment\n */    1.4   ;\n",
			expected: "dml /*\n * comment\n */ 1.4;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "09-all-comments",
			code: `// Before dml
/*
 * Multi-line comment
 */

// Normal comment

/* line start */    dml    /* inline multi-line */    1.4   /* semi */  ;       // sl comment

// After dml


/* Another multi-line */
`,
			// TODO: Keep at max 1 empty line between comments
			// 		 that have empty lines between them
			expected: `// Before dml
/*
 * Multi-line comment
 */
// Normal comment
/* line start */ dml /* inline multi-line */ 1.4 /* semi */; // sl comment
// After dml
/* Another multi-line */
`,
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintDeviceTag(t *testing.T) {
	t.Parallel()

	const lineCap = 30
	const EXPECTED = "dml 1.4;\n\n\ndevice name;\n"
	testCases := []printerTestCase{
		{
			name:     "01-1:1",
			code:     EXPECTED,
			expected: EXPECTED,
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-wrong-ident",
			code:     "    dml 1.4;  \n   device name;  \n",
			expected: "dml 1.4;\ndevice name;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-extra-newlines",
			code:     "\n\n\n\ndml 1.4;\n\n\n\ndevice name;\n\n\n\n",
			expected: EXPECTED,
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-extra-whitespace",
			code:     "dml    1.4   ;\ndevice      name   ;",
			expected: "dml 1.4;\ndevice name;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-extra-new-lines-inside",
			code:     "dml  \n  1.4  \n\n ;\ndevice \n\n  name;\n",
			expected: "dml 1.4;\ndevice name;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-normal-before",
			code:     "dml 1.4;\n\n   // comment\n\n\ndevice name;\n",
			expected: "dml 1.4;\n\n// comment\ndevice name;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-normal-sameline-tail",
			code:     "dml 1.4;\n\n\ndevice name;   // comment\n",
			expected: "dml 1.4;\n\n\ndevice name; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-normal-after",
			code:     "dml 1.4;\n\n\ndevice name;\n\n\n// comment\n\n\n",
			expected: "dml 1.4;\n\n\ndevice name;\n// comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-multiline-before",
			code:     "dml 1.4;\n\n\n/*\n * comment\n */\n\n\ndevice name;\n",
			expected: "dml 1.4;\n\n\n/*\n * comment\n */\ndevice name;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-multiline-sameline-tail",
			code:     "dml 1.4;\n\n\ndevice name;    /* comment */\n",
			expected: "dml 1.4;\n\n\ndevice name; /* comment */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-multiline-tail",
			code:     "dml 1.4;\n\n\n  device    name;      /*\n * comment\n */\n\n\n",
			expected: "dml 1.4;\n\n\ndevice name; /*\n * comment\n */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "12-multiline-tail-long",
			code:     "dml 1.4;\n\n\n  device    name;      /*\n * comment that more than 30 charactes but on next line\n */\n\n\n",
			expected: "dml 1.4;\n\n\ndevice name; /*\n * comment that more than 30 charactes but on next line\n */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "13-multiline-after",
			code:     "dml 1.4;\n\n\ndevice name;\n\n/*\n * comment\n */\n\n\n",
			expected: "dml 1.4;\n\n\ndevice name;\n/*\n * comment\n */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "14-multiline-inline",
			code:     "dml 1.4;\n\n\ndevice  \n  /* comment */     name \n;\n",
			expected: "dml 1.4;\n\n\ndevice /* comment */ name;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "15-multiline-inline-semi",
			code:     "dml 1.4;\n\n\ndevice name    \n /* comment */   ;\n",
			expected: "dml 1.4;\n\n\ndevice name /* comment */;\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "100-all-comments",
			code: `dml 1.4;
// Before device
/*
 * Multi-line comment
 */

// Normal comment

    /* line start */    device    /* inline multi-line */    dev_name     /* semi */    ;        // slcomment

// After device


/* Another multi-line */
`,
			expected: `dml 1.4;
// Before device
/*
 * Multi-line comment
 */
// Normal comment
/* line start */ device /* inline multi-line */ dev_name /* semi */; // slcomment
// After device
/* Another multi-line */
`,
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintBitorderTag(t *testing.T) {
	t.Parallel()

	const lineCap = 30
	const EXPECTED = "device name;\n\n\nbitorder le;\n"
	testCases := []printerTestCase{
		{
			name:     "01-1:1",
			code:     EXPECTED,
			expected: EXPECTED,
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-wrong-ident",
			code:     "  \n   device name;  \n     bitorder le;",
			expected: "device name;\nbitorder le;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-extra-newlines",
			code:     "\n\n\n\n\ndevice name;\n\n\n\n\n\nbitorder le;\n\n\n\n",
			expected: EXPECTED,
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-extra-whitespace",
			code:     "\ndevice      name   ;\nbitorder     le     ;",
			expected: "device name;\nbitorder le;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-extra-new-lines-inside",
			code:     "device \n\n  name;\nbitorder\n\n    le\n;",
			expected: "device name;\nbitorder le;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-normal-before",
			code:     "\n\n\ndevice name;\n\n\n   // comment\n\n\nbitorder le;\n",
			expected: "device name;\n\n\n// comment\nbitorder le;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "07-normal-sameline-tail",
			code:     "device name;\n\n\nbitorder le;   // comment\n",
			expected: "device name;\n\n\nbitorder le; // comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "08-normal-after",
			code:     "device name;\n\n\nbitorder le;\n     // comment\n",
			expected: "device name;\n\n\nbitorder le;\n// comment\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "09-multiline-before",
			code:     "/*\n * comment\n */\ndevice name;\n\n\nbitorder le;\n",
			expected: "/*\n * comment\n */\ndevice name;\n\n\nbitorder le;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "10-multiline-sameline-tail",
			code:     "device name;\n\n\nbitorder le;    /* comment */\n",
			expected: "device name;\n\n\nbitorder le; /* comment */\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "11-multiline-after",
			code:     "device name;\n\n\nbitorder le;\n/*\n * comment\n */\n",
			expected: "device name;\n\n\nbitorder le;\n/*\n * comment\n */\n",
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}

func TestPrintProvisionalTag(t *testing.T) {
	t.Parallel()

	const (
		lineCap     = 30
		lineCapLong = 50
	)

	const (
		EXPECTED_1IDENT           = "dml 1.4;\n\n\nprovisional name;\n"
		EXPECTED_3ident_short     = "dml 1.4;\n\n\nprovisional name, other_name, another;\n"
		EXPECTED_3ident_long      = "dml 1.4;\n\n\nprovisional (\n    name, other_name, another\n);\n"
		EXPECTED_3ident_very_long = "dml 1.4;\n\n\nprovisional (\n    name,\n    other_name,\n    yet_another_very_very_long_name\n);\n"
	)
	testCases := []printerTestCase{
		{
			name:     "001-1ident-1:1",
			code:     EXPECTED_1IDENT,
			expected: EXPECTED_1IDENT,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "002-1ident-wrong-ident",
			code:     "    dml 1.4;  \n   provisional name;  \n",
			expected: "dml 1.4;\nprovisional name;\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "003-1ident-extra-newlines",
			code:     "\n\ndml 1.4;\n\n\n\n\n\nprovisional name;\n\n\n\n\n\n",
			expected: EXPECTED_1IDENT,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "004-1ident-extra-whitespace",
			code:     "dml    1.4   ;\nprovisional      name   ;",
			expected: "dml 1.4;\nprovisional name;\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "005-1ident-extra-new-lines-inside",
			code:     "dml  \n  1.4  \n\n ;\nprovisional \n\n  name;\n",
			expected: "dml 1.4;\nprovisional name;\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\nprovisional (name);\n",
			),
			name:     "006-1ident-redundent-multiline",
			code:     "dml  \n  1.4  \n\n ;\nprovisional (\n    name \n ) ;\n",
			expected: EXPECTED_1IDENT,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "007-1ident-normal-before",
			code:     "\ndml 1.4;\n\n   // comment\n\n\nprovisional name;\n",
			expected: "dml 1.4;\n\n// comment\nprovisional name;\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "008-1ident-normal-sameline-tail",
			code:     "dml 1.4;\n\n\nprovisional name;   // comment\n",
			expected: "dml 1.4;\n\n\nprovisional name; // comment\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "009-1ident-normal-after",
			code:     "dml 1.4;\n\n\nprovisional name;\n\n// comment\n\n\n",
			expected: "dml 1.4;\n\n\nprovisional name;\n// comment\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "010-1ident-multiline-before",
			code:     "dml 1.4;\n\n/*\n * comment\n */\n\nprovisional name;\n",
			expected: "dml 1.4;\n\n/*\n * comment\n */\nprovisional name;\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "011-1ident-multiline-sameline-tail",
			code:     "dml 1.4;\n\n\nprovisional name;    /* cmnt */\n",
			expected: "dml 1.4;\n\n\nprovisional name; /* cmnt */\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "012-1ident-multiline-after",
			code:     "dml 1.4;\n\n\nprovisional name;\n\n/*\n * comment\n */\n\n",
			expected: "dml 1.4;\n\n\nprovisional name;\n/*\n * comment\n */\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "013-1ident-multiline-inline",
			code:     "dml 1.4;\n\n\nprovisional    /* cmnt */    name; \n",
			expected: "dml 1.4;\n\n\nprovisional /* cmnt */ name;\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "100-3ident_short-1:1",
			code:     EXPECTED_3ident_short,
			expected: EXPECTED_3ident_short,
			lineCap:  lineCapLong,
			pc:       func(f func()) { f() },
		},
		{
			name:     "101-3ident_short-wrong-ident",
			code:     "    dml 1.4;  \n   provisional name, other_name, another;  \n",
			expected: "dml 1.4;\nprovisional name, other_name, another;\n",
			lineCap:  lineCapLong,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\n\n\nprovisional (name, other_name, another);\n",
			),
			name:     "102-3ident_short-extra-newlines",
			code:     "\n\ndml 1.4;\n\n\n\n\n\nprovisional (\n\n\nname,\n\nother_name,\n\nanother\n\n);\n\n\n\n\n\n",
			expected: EXPECTED_3ident_short,
			lineCap:  lineCapLong,
			pc:       func(f func()) { f() },
		},
		{
			name:     "103-3ident_short-extra-whitespace",
			code:     "dml    1.4   ;\nprovisional      name,         other_name, \n\n\n  another   ;",
			expected: "dml 1.4;\nprovisional name, other_name, another;\n",
			lineCap:  lineCapLong,
			pc:       func(f func()) { f() },
		},
		{
			name:     "104-3ident_short-extra-new-lines-inside",
			code:     "dml  \n  1.4  \n\n ;\nprovisional \n\n  name, other_name, another;\n",
			expected: "dml 1.4;\nprovisional name, other_name, another;\n",
			lineCap:  lineCapLong,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\nprovisional (name, other_name, another);\n",
			),
			name:     "105-3ident_short-redundent-multiline",
			code:     "dml  \n  1.4  \n\n ;\nprovisional (\n    name,\n    other_name,\n    another\n);\n",
			expected: EXPECTED_3ident_short,
			lineCap:  lineCapLong,
			pc:       func(f func()) { f() },
		},
		{
			name:     "200-3ident_long-1:1",
			code:     EXPECTED_3ident_long,
			expected: EXPECTED_3ident_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\nprovisional name,\n    other_name,\n    another;\n",
			),
			name:     "201-3ident_long-wrong-ident",
			code:     "    dml 1.4;  \n   provisional name, other_name, another;  \n",
			expected: EXPECTED_3ident_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "202-3ident_long-extra-newlines",
			code:     "\n\ndml 1.4;\n\n\n\n\n\nprovisional (\n\n\nname,\n\nother_name,\n\nanother\n\n);\n\n\n\n\n\n",
			expected: EXPECTED_3ident_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\nprovisional name,\n    other_name,\n    another;\n",
			),
			name:     "203-3ident_long-extra-whitespace",
			code:     "dml    1.4   ;\nprovisional      name,         other_name, \n\n\n  another   ;",
			expected: EXPECTED_3ident_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\nprovisional name,\n    other_name,\n    another;\n",
			),
			name:     "204-3ident_long-extra-new-lines-inside",
			code:     "dml  \n  1.4  \n\n ;\nprovisional \n\n  name, other_name, another;\n",
			expected: EXPECTED_3ident_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "205-3ident_long-redundent-multiline",
			code:     "dml  \n  1.4  \n\n ;\nprovisional (\n    name,\n    other_name,\n    another\n);\n",
			expected: "dml 1.4;\nprovisional (\n    name, other_name, another\n);\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "300-3ident_very_long-1:1",
			code:     EXPECTED_3ident_very_long,
			expected: EXPECTED_3ident_very_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\nprovisional name,\n    other_name,\n    yet_another_very_very_long_name;\n",
			),
			name:     "301-3ident_very_long-wrong-ident",
			code:     "    dml 1.4;  \n   provisional name, other_name, yet_another_very_very_long_name;  \n",
			expected: EXPECTED_3ident_very_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "302-3ident_very_long-extra-newlines",
			code:     "\n\ndml 1.4;\n\n\n\n\n\nprovisional (\n\n\nname,\n\nother_name,\n\nyet_another_very_very_long_name\n\n);\n\n\n\n\n\n",
			expected: EXPECTED_3ident_very_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\nprovisional name,\n    other_name,\n    yet_another_very_very_long_name;\n",
			),
			name:     "303-3ident_very_long-extra-whitespace",
			code:     "dml    1.4   ;\nprovisional      name,         other_name, \n\n\n  yet_another_very_very_long_name   ;",
			expected: EXPECTED_3ident_very_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\nprovisional name,\n    other_name,\n    yet_another_very_very_long_name;\n",
			),
			name:     "304-3ident_very_long-extra-new-lines-inside",
			code:     "dml  \n  1.4  \n\n ;\nprovisional \n\n  name, other_name, yet_another_very_very_long_name;\n",
			expected: EXPECTED_3ident_very_long,
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			name:     "305-3ident_very_long-redundent-multiline",
			code:     "dml  \n  1.4  \n\n ;\nprovisional (\n    name,\n    other_name,\n    yet_another_very_very_long_name\n);\n",
			expected: "dml 1.4;\nprovisional (\n    name,\n    other_name,\n    yet_another_very_very_long_name\n);\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\nprovisional yet_another_very_very_long_name;\n",
			),
			name:     "400-1ident_very_long",
			code:     "dml  1.4;\nprovisional yet_another_very_very_long_name;\n",
			expected: "dml 1.4;\n\n\nprovisional (\n    yet_another_very_very_long_name\n);\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\n\n\nprovisional /*a*/ name /*x*/,\n    /*b*/ other_name /*y*/,\n    /*c*/ another /*z*/;\n",
			),
			name:     "500-3ident-comments",
			code:     "dml 1.4;\n\n\nprovisional  /*a*/   name   /*x*/  ,  /*b*/   other_name   /*y*/  ,   /*c*/   another   /*z*/  ;\n",
			expected: "dml 1.4;\n\n\nprovisional (\n    /*a*/ name /*x*/,\n    /*b*/ other_name /*y*/,\n    /*c*/ another /*z*/\n);\n",
			lineCap:  lineCap,
			pc:       func(f func()) { f() },
		},
		{
			masked: MaskSyntheticProvisionalParens(
				"dml 1.4;\n\n\nprovisional /*a*/ a /*x*/,\n    /*b*/ b /*y*/,\n    /*c*/ c /*z*/;\n",
			),
			name:     "501-3ident-comments-shorts",
			code:     "dml 1.4;\n\n\nprovisional   /*a*/   a   /*x*/  ,   /*b*/   b   /*y*/  ,   /*c*/   c   /*z*/  ;\n",
			expected: "dml 1.4;\n\n\nprovisional (\n    /*a*/ a /*x*/, /*b*/ b /*y*/, /*c*/ c /*z*/\n);\n",
			lineCap:  50,
			pc:       func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, lineCap)
}
