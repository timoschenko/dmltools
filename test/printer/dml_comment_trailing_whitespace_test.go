// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"testing"
)

func TestCommentTrailingWhitespace(t *testing.T) {
	t.Parallel()

	testCases := []printerTestCase{
		// ==========================================
		// Category 1: Single-Line (//) Comments
		// ==========================================
		{
			name:     "01-line-comment-trailing-spaces",
			code:     "param x = 1; // comment with spaces   \n",
			expected: "param x = 1; // comment with spaces\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "02-line-comment-trailing-tabs",
			code:     "param x = 1; // comment with tabs\t\t\n",
			expected: "param x = 1; // comment with tabs\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "03-line-comment-trailing-mixed",
			code:     "param x = 1; // mixed spaces and tabs \t \n",
			expected: "param x = 1; // mixed spaces and tabs\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "04-line-comment-empty-spaces",
			code:     "param x = 1; //   \n",
			expected: "param x = 1; //\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "05-line-comment-empty-tabs",
			code:     "param x = 1; //\t\t\n",
			expected: "param x = 1; //\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "06-line-comment-empty-bare",
			code:     "param x = 1; //\n",
			expected: "param x = 1; //\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "07-line-comment-consecutive-trailing-spaces",
			code: "// line one   \n" +
				"// line two    \n" +
				"param x = 1;\n",
			expected: "// line one\n" +
				"// line two\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "08-line-comment-at-eof-no-newline-spaces",
			code: "param x = 1;\n" +
				"// trailing comment at EOF   ",
			expected: "param x = 1;\n" +
				"// trailing comment at EOF",
			pc: func(f func()) { f() },
		},
		{
			name: "09-line-comment-at-eof-no-newline-empty-spaces",
			code: "param x = 1;\n" +
				"//   ",
			expected: "param x = 1;\n" +
				"//",
			pc: func(f func()) { f() },
		},
		{
			name: "10-line-comment-triple-slash-banner",
			code: "/// documentation banner   \n" +
				"param x = 1;\n",
			expected: "/// documentation banner\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "11-line-comment-dashes-banner",
			code: "//---------------------------------   \n" +
				"param x = 1;\n",
			expected: "//---------------------------------\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "12-line-comment-inside-method",
			code: "method run() {\n" +
				"// step 1   \n" +
				"local int a = 1; // assign   \n" +
				"// empty comment with spaces   \n" +
				"//   \n" +
				"}\n",
			expected: "method run() {\n" +
				"    // step 1\n" +
				"    local int a = 1; // assign\n" +
				"    // empty comment with spaces\n" +
				"    //\n" +
				"}\n",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 2: Single-Line Block (/* ... */) Comments
		// ==========================================
		{
			name: "20-singleline-block-trailing-spaces-outside",
			code: "/* comment */   \n" +
				"param x = 1;\n",
			expected: "/* comment */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "21-singleline-block-trailing-tabs-outside",
			code: "/* comment */\t\t\n" +
				"param x = 1;\n",
			expected: "/* comment */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "22-singleline-block-trailing-spaces-inside",
			code: "/* comment   */\n" +
				"param x = 1;\n",
			expected: "/* comment */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "23-singleline-block-trailing-tabs-inside",
			code: "/* comment\t\t*/\n" +
				"param x = 1;\n",
			expected: "/* comment */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "24-singleline-block-trailing-inside-and-outside",
			code: "/* comment   */   \n" +
				"param x = 1;\n",
			expected: "/* comment */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "25-singleline-block-empty-spaces",
			code: "/*   */\n" +
				"param x = 1;\n",
			expected: "/* */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "26-singleline-block-empty-tabs",
			code: "/*\t\t*/\n" +
				"param x = 1;\n",
			expected: "/* */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "27-singleline-block-empty-bare",
			code: "/**/\n" +
				"param x = 1;\n",
			expected: "/**/\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "28-singleline-block-empty-single-space",
			code: "/* */\n" +
				"param x = 1;\n",
			expected: "/* */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "29-singleline-block-asterisks-banner-with-spaces",
			code: "/**** text   ****/\n" +
				"param x = 1;\n",
			expected: "/**** text ****/\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "30-singleline-block-asterisks-banner-spaces-only",
			code: "/****   ****/\n" +
				"param x = 1;\n",
			expected: "/**** ****/\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "31-singleline-block-asterisks-banner-bare",
			code: "/******/\n" +
				"param x = 1;\n",
			expected: "/******/\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name:     "32-singleline-block-inline-in-code",
			code:     "param x = /* inline   */ 10;\n",
			expected: "param x = /* inline */ 10;\n",
			pc:       func(f func()) { f() },
		},
		{
			name:     "33-singleline-block-after-code-statement",
			code:     "param x = 1; /* trailing annotation   */   \n",
			expected: "param x = 1; /* trailing annotation */\n",
			pc:       func(f func()) { f() },
		},
		{
			name: "34-singleline-block-eof-no-newline",
			code: "param x = 1;\n" +
				"/* trailing eof block   */",
			expected: "param x = 1;\n" +
				"/* trailing eof block */",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 3: Multi-Line Block (/* ... \n ... */) Comments
		// ==========================================
		{
			name: "40-multiline-block-trailing-on-first-line",
			code: "/* first line   \n" +
				" * second line\n" +
				" */\n" +
				"param x = 1;\n",
			expected: "/* first line\n" +
				" * second line\n" +
				" */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "41-multiline-block-first-line-empty-spaces",
			code: "/*   \n" +
				" * second line\n" +
				" */\n" +
				"param x = 1;\n",
			expected: "/*\n" +
				" * second line\n" +
				" */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "42-multiline-block-middle-lines-trailing-spaces",
			code: "/*\n" +
				" * line 1   \n" +
				" * line 2    \n" +
				" */\n" +
				"param x = 1;\n",
			expected: "/*\n" +
				" * line 1\n" +
				" * line 2\n" +
				" */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "43-multiline-block-middle-lines-trailing-tabs",
			code: "/*\n" +
				" * line 1\t\t\n" +
				" * line 2\t \n" +
				" */\n" +
				"param x = 1;\n",
			expected: "/*\n" +
				" * line 1\n" +
				" * line 2\n" +
				" */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "44-multiline-block-whitespace-only-intermediate-lines",
			code: "/*\n" +
				" * line 1\n" +
				"   \n" +
				"\t\t\n" +
				" * line 2\n" +
				" */\n" +
				"param x = 1;\n",
			expected: "/*\n" +
				" * line 1\n\n\n" +
				" * line 2\n" +
				" */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "45-multiline-block-closing-line-text-with-spaces",
			code: "/*\n" +
				" * line 1\n" +
				" * line 2   */\n" +
				"param x = 1;\n",
			expected: "/*\n" +
				" * line 1\n" +
				" * line 2 */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "46-multiline-block-closing-line-only-indentation",
			code: "/*\n" +
				" * line 1\n" +
				"   */\n" +
				"param x = 1;\n",
			expected: "/*\n" +
				" * line 1\n" +
				"   */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "47-multiline-block-closing-line-trailing-spaces-outside",
			code: "/*\n" +
				" * line 1\n" +
				" */   \n" +
				"param x = 1;\n",
			expected: "/*\n" +
				" * line 1\n" +
				" */\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "48-multiline-block-banner-asterisks",
			code: "/************************   \n" +
				" * banner text   \n" +
				" ************************/   \n" +
				"param x = 1;\n",
			expected: "/************************\n" +
				" * banner text\n" +
				" ************************/\n" +
				"param x = 1;\n",
			pc: func(f func()) { f() },
		},
		{
			name: "49-multiline-block-eof-no-newline",
			code: "param x = 1;\n" +
				"/*\n" +
				" * line 1   \n" +
				" */",
			expected: "param x = 1;\n" +
				"/*\n" +
				" * line 1\n" +
				" */",
			pc: func(f func()) { f() },
		},

		// ==========================================
		// Category 4: Mixed Comments and Complex Structures
		// ==========================================
		{
			name: "60-mixed-line-and-block-comments",
			code: "// header line   \n" +
				"/* multiline header   \n" +
				" * line 2   \n" +
				" */   \n" +
				"param val = 42; // end comment   \n",
			expected: "// header line\n" +
				"/* multiline header\n" +
				" * line 2\n" +
				" */\n" +
				"param val = 42; // end comment\n",
			pc: func(f func()) { f() },
		},
		{
			name: "61-comment-in-list-delimiters",
			code: "constant arr = [ // item comment   \n" +
				"1, // first   \n" +
				"2, // second   \n" +
				"]; // array end   \n",
			expected: "constant arr = [ // item comment\n" +
				"    1, // first\n" +
				"    2, // second\n" +
				"]; // array end\n",
			pc: func(f func()) { f() },
		},
		{
			name: "62-comment-in-bank-and-register",
			code: "bank regs { // bank comment   \n" +
				"register ctrl size 4 @ 0x00 { // reg comment   \n" +
				"/* field doc   */   \n" +
				"field en @ [0]; // bit 0   \n" +
				"} // reg close   \n" +
				"} // bank close   \n",
			expected: "bank regs { // bank comment\n" +
				"    register ctrl size 4 @ 0x00 { // reg comment\n" +
				"        /* field doc */\n" +
				"        field en @ [0]; // bit 0\n" +
				"    } // reg close\n" +
				"} // bank close\n",
			pc: func(f func()) { f() },
		},
	}

	runPrinterTestsuite(t, runPrinterTestCase, testCases, 80)
}
