// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"strings"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

type switchParserTestCase struct {
	name string
	code string
}

func runSwitchParserTestCase(t *testing.T, tc *switchParserTestCase) {
	t.Helper()

	wrap := func(val string) string {
		var sb strings.Builder
		sb.WriteString("dml 1.4;\n\nmethod foo() {\n")
		for line := range strings.SplitSeq(val, "\n") {
			if line != "" {
				sb.WriteString("    ")
				sb.WriteString(line)
			}
			sb.WriteString("\n")
		}
		sb.WriteString("}\n")
		return sb.String()
	}

	code := tc.code
	if !strings.HasPrefix(code, "dml 1.4;") {
		code = wrap(code)
	}

	inputStream := antlr.NewInputStream(code)
	errorListener := dmlex.NewErrorListener("test.dml", 0)
	var syntaxErrors []dmlex.SyntaxErrorInfo
	errorListener.Handler = func(e dmlex.SyntaxErrorInfo) {
		syntaxErrors = append(syntaxErrors, e)
	}

	file := dmlex.ParseStream(inputStream, errorListener)
	_ = file.Parser.Dml()

	if len(syntaxErrors) > 0 {
		t.Fatalf(
			"unexpected syntax errors (%d): %v",
			len(syntaxErrors),
			syntaxErrors,
		)
	}
}

func TestParseSwitchStatement(t *testing.T) {
	t.Parallel()

	testCases := []switchParserTestCase{
		{
			name: "01-empty-switch",
			code: "switch (x) {}",
		},
		{
			name: "02-single-case-break",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"    break;\n" +
				"}",
		},
		{
			name: "03-case-with-default",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"    call();\n" +
				"    break;\n" +
				"default:\n" +
				"    other();\n" +
				"    break;\n" +
				"}",
		},
		{
			name: "04-hash-if-enclosing-cases",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"#if (cond) {\n" +
				"case 2:\n" +
				"    call();\n" +
				"    break;\n" +
				"}\n" +
				"default:\n" +
				"    break;\n" +
				"}",
		},
		{
			name: "05-break-after-hash-if",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"    test();\n" +
				"    break;\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"    test2();\n" +
				"}\n" +
				"    break;\n" +
				"default:\n" +
				"    log error, 1, log_group_ident: \"message %x\", foo;\n" +
				"}",
		},
		{
			name: "06-case-labels-inside-hash-if-with-fallthrough",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"    test();\n" +
				"    break;\n" +
				"case B:\n" +
				"#if (FEAT) {\n" +
				"case CA:\n" +
				"}\n" +
				"#if (FEAT_2) {\n" +
				"case CB:\n" +
				"}\n" +
				"    val = test2();\n" +
				"    break;\n" +
				"default:\n" +
				"    log error, 1, log_group_ident: \"message %x\", foo;\n" +
				"}",
		},
		{
			name: "07-multiple-statements-after-hash-if",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"    test();\n" +
				"}\n" +
				"    val = test2();\n" +
				"    break;\n" +
				"}",
		},
		{
			name: "08-hash-if-case-label-fallthrough-to-statement",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"}\n" +
				"    common_code();\n" +
				"    break;\n" +
				"}",
		},
		{
			name: "09-hash-if-else-enclosing-cases",
			code: "switch (x) {\n" +
				"case 1:\n" +
				"#if (cond) {\n" +
				"case 2:\n" +
				"    call();\n" +
				"    break;\n" +
				"} #else {\n" +
				"case 3:\n" +
				"    other();\n" +
				"    break;\n" +
				"}\n" +
				"default:\n" +
				"    break;\n" +
				"}",
		},
		{
			name: "10-break-after-hash-if-else",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"    test();\n" +
				"    break;\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"    test2();\n" +
				"} #else {\n" +
				"case C:\n" +
				"    test3();\n" +
				"}\n" +
				"    break;\n" +
				"default:\n" +
				"    log error, 1, log_group_ident: \"message %x\", foo;\n" +
				"}",
		},
		{
			name: "11-case-labels-inside-hash-if-else-with-fallthrough",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"    test();\n" +
				"    break;\n" +
				"case B:\n" +
				"#if (FEAT) {\n" +
				"case CA:\n" +
				"} #else {\n" +
				"case CB:\n" +
				"}\n" +
				"    val = test2();\n" +
				"    break;\n" +
				"default:\n" +
				"    log error, 1, log_group_ident: \"message %x\", foo;\n" +
				"}",
		},
		{
			name: "12-multiple-statements-after-hash-if-else",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"    test();\n" +
				"} #else {\n" +
				"case C:\n" +
				"    test2();\n" +
				"}\n" +
				"    val = test3();\n" +
				"    break;\n" +
				"}",
		},
		{
			name: "13-hash-if-else-case-label-fallthrough-to-statement",
			code: "switch (foo) {\n" +
				"case A:\n" +
				"#if (FEAT) {\n" +
				"case B:\n" +
				"} #else {\n" +
				"case C:\n" +
				"}\n" +
				"    common_code();\n" +
				"    break;\n" +
				"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runSwitchParserTestCase(t, &tc)
		})
	}
}
