// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"bytes"
	"errors"
	"fmt"
	dmlex "github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer"
	lextest "github.com/timoschenko/dmltools/test/dmlex"
	"math"
	"runtime"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

const isTokenDump = false // set to true to dump tokens to logs for extended logs

// GetMemStats captures a clean snapshot of the current memory state.
func GetMemStats(m *runtime.MemStats) {
	// Optional: runtime.GC() ensures we measure live objects,
	// not just pending garbage.
	runtime.GC()
	runtime.ReadMemStats(m)
}

func diffUint64(end, start uint64) int64 {
	if end >= start {
		diff := end - start
		if diff > math.MaxInt64 {
			return math.MaxInt64
		}
		return int64(diff)
	}
	diff := start - end
	if diff > math.MaxInt64 {
		return math.MinInt64
	}
	return -int64(diff)
}

// LogMemDiff calculates and logs the delta between two MemStats snapshots.
func LogMemDiff(t *testing.T, label string, start, end *runtime.MemStats) {
	t.Helper() // Marks this func as a helper so line numbers in logs point to the caller

	GetMemStats(end)
	// Calculate deltas
	allocDiff := diffUint64(end.Alloc, start.Alloc)
	sysDiff := diffUint64(end.Sys, start.Sys)
	mallocs := end.Mallocs - start.Mallocs
	gcCount := end.NumGC - start.NumGC

	t.Logf(
		"--- Memory Diff: %s { Heap Alloc: %+.2f MiB (%d bytes), "+
			"Total Sys: %+.2f MiB, Mallocs: %d, GC Cycles: %d }\n",
		label, bToMb(allocDiff), allocDiff, bToMb(sysDiff), mallocs, gcCount,
	)
}

func bToMb(b int64) float64 {
	return float64(b) / 1024 / 1024
}

func dumpTokensToLogs(testSuite *testing.T, dml dmlex.IDmlContext) {
	testSuite.Helper()

	ts := dml.GetParser().GetTokenStream()

	start := dml.GetStart()
	stop := dml.GetStop()

	dumpTokensToLogsRaw(testSuite, ts, start, stop)

}

func dumpTokensToLogsRaw(
	testSuite *testing.T, ts antlr.TokenStream, start, stop antlr.Token,
) {
	testSuite.Helper()

	jsonSanitize := func(s string) string {
		s = strings.ReplaceAll(s, "\\", "\\\\")
		for _, esc := range []struct{ old, new string }{
			{"\n", "\\n"},
			{"\r", "\\r"},
			{"\t", "\\t"},
			{"'", "\\'"},
		} {
			s = strings.ReplaceAll(s, esc.old, esc.new)
		}
		return s
	}
	cts, ok := ts.(*antlr.CommonTokenStream)
	if !ok {
		panic("unreachable: " + fmt.Sprintf("%T", ts))
	}

	allTokens := cts.GetAllTokens()

	sb := strings.Builder{}
	_, _ = sb.WriteString("TOKENS: {\n")
	for i := start.GetTokenIndex(); i <= stop.GetTokenIndex(); i++ {
		token := allTokens[i]
		_, _ = fmt.Fprintf(
			&sb,
			"% 3d: {type:% 4d, name: % 16s, lineno:% 3d, "+
				"colno:% 3d, range:[% 4d,% 4d], text:'%s'}\n",
			i,
			token.GetTokenType(),
			dmlex.TokenType(token.GetTokenType()).String(),
			token.GetLine(),
			token.GetColumn(),
			token.GetStart(),
			token.GetStop(),
			jsonSanitize(token.GetText()),
		)
	}

	_, _ = sb.WriteString("}\n")
	testSuite.Log(sb.String())
}

// Wraps code for debbugin and add line numbers, and line lenght
func codePretty(code string) string {
	var sb strings.Builder
	for i, line := range strings.Split(code, "\n") {
		_, _ = fmt.Fprintf(&sb, "% 3d (% 3d)>%s\n", i+1, len(line), line)
	}
	return sb.String()
}

type printerTestCase struct {
	task006Skip          bool
	name, code, expected string

	// Mark if test has known bug (leads to SKIP state)
	masked Masked

	// if 0 then default
	lineCap int

	// helper to let tester point to exact test case
	pc func(f func())
}

type testRunnerFunc func(testSuite *testing.T, ptc *printerTestCase)

func runPrinterTestCase(testSuite *testing.T, ptc *printerTestCase) {
	testSuite.Helper()

	testSuite.Run(ptc.name, func(t *testing.T) {
		t.Helper() // use pc func

		if ptc.task006Skip {
			t.Skip("Skipping due to compiler workarounds task006")
		}

		mask, err := normalizeMask(ptc.masked)
		if err != nil {
			t.Fatalf("invalid mask: %v", err)
		}

		if mask != nil && errors.Is(mask, ErrMaskSkip) {
			var skipMask MaskedSkip
			reason := mask.Error()
			if errors.As(mask, &skipMask) {
				reason = skipMask.Reason()
			}
			if mask.Issue() > 0 {
				t.Skipf("[Issue #%d] %s", mask.Issue(), reason)
			} else {
				t.Skipf("[Uncategorized] %s", reason)
			}
			return
		}

		tc := *ptc
		tc.code = strings.ReplaceAll(tc.code, "`", "\"")
		tc.expected = strings.ReplaceAll(tc.expected, "`", "\"")

		isVerbose := testing.Verbose()
		if isVerbose {
			t.Logf("INPUT (linecap %d): ```\n%s\n```\n",
				tc.lineCap, codePretty(tc.code))
		}

		const dmlHeader = "dml 1.4;\n\n\n"
		hasDmlHeader := (strings.HasPrefix(tc.expected, "dml 1.4;") ||
			strings.Contains(tc.expected, "dml 1.4;") ||
			strings.Contains(tc.expected, "dml /*") ||
			strings.Contains(tc.expected, "*/ dml"))
		code := tc.code
		if !hasDmlHeader {
			code = dmlHeader + code
		}

		isSyntaxMask := mask != nil && errors.Is(mask, ErrMaskSyntaxError)

		var syntaxErrors []dmlex.SyntaxErrorInfo
		var file *dmlex.File

		if isSyntaxMask {
			inputStream := antlr.NewInputStream(code)
			errorListener := dmlex.NewErrorListener("", 1)
			errorListener.Handler = func(e dmlex.SyntaxErrorInfo) {
				syntaxErrors = append(syntaxErrors, e)
			}
			file = dmlex.ParseStream(inputStream, errorListener)
		} else {
			file = lextest.ParseDmlFileFromText(
				t, code, "", 1,
			)
			if t.Failed() {
				// fail to parse file
				return
			}
		}

		dml := file.Parser.Dml()

		if isSyntaxMask {
			if len(syntaxErrors) > 0 {
				if mask.Issue() > 0 {
					t.Skipf(
						"[Issue #%d] test case has known syntax error: %s",
						mask.Issue(),
						syntaxErrors[0].Msg,
					)
				} else {
					t.Skipf(
						"[Uncategorized] test case has known syntax error: %s",
						syntaxErrors[0].Msg,
					)
				}
				return
			}
			t.Fatalf(
				"test case expected syntax error but parsed cleanly [Issue #%d]. Consider removing the mask",
				mask.Issue(),
			)
			return
		}

		if t.Failed() {
			// fail to parse file
			return
		}
		if isVerbose && isTokenDump {
			dumpTokensToLogs(t, dml)
		}
		buffer := bytes.NewBufferString("")
		if err := printer.PrinterWriteTo(
			dml,
			file,
			buffer,
			tc.lineCap,
		); err != nil {
			t.Errorf("PrinterWriteTo failed: %v", err)
			return
		}
		result := buffer.String()
		if !hasDmlHeader {
			if strings.HasPrefix(result, "dml 1.4;\n\n\n") {
				result = result[len("dml 1.4;\n\n\n"):]
			} else if strings.HasPrefix(result, "dml 1.4;\n\n") {
				result = result[len("dml 1.4;\n\n"):]
			} else if strings.HasPrefix(result, "dml 1.4;\n") {
				result = result[len("dml 1.4;\n"):]
			} else {
				t.Errorf("Code has missing dml header")
				return
			}
		}
		if isVerbose {
			t.Logf("EXPECT: (linecap %d)```\n%s\n```\n",
				tc.lineCap, codePretty(tc.expected))

			t.Logf("OUTPUT: (linecap %d)```\n%s\n```\n",
				tc.lineCap, codePretty(result))
		}

		var maskedFormatter MaskedFormatter
		var maskedExpected string
		hasMaskedExpected := false
		if mask != nil && errors.As(mask, &maskedFormatter) {
			maskedExpected = maskedFormatter.Expected()
			hasMaskedExpected = true
		}

		switch {
		case result == tc.expected:
			if mask != nil {
				t.Errorf(
					"test case passed but it is still masked. Consider removing the mask",
				)

			}
		case hasMaskedExpected && result == maskedExpected:
			if mask.Issue() > 0 {
				t.Skipf("[Issue #%d] test case masked", mask.Issue())
			} else {
				t.Skip("test case has been masked")
			}
		default:
			maskedMsg := ""
			if hasMaskedExpected {
				maskedMsg = fmt.Sprintf(
					"\n  Masked: %#v",
					strings.ReplaceAll(maskedExpected, "\"", "`"),
				)
			}
			t.Errorf(
				"Code is not formatted correctly\n"+
					"  Expect: %#v%s\n"+
					"     Got: %#v",
				strings.ReplaceAll(tc.expected, "\"", "`"),
				maskedMsg,
				strings.ReplaceAll(result, "\"", "`"),
			)
		}
	})
}

func runPrinterMethodTestCase(testSuite *testing.T, ptc *printerTestCase) {
	testSuite.Helper()

	wrap := func(val string) string {
		if val == "" {
			return ""
		}
		var sb strings.Builder
		sb.WriteString("dml 1.4;\n\n\nmethod foo() {\n")
		for i, line := range strings.Split(val, "\n") {
			if i > 0 {
				sb.WriteString("\n")
			}
			if line != "" {
				sb.WriteString("    ")
				sb.WriteString(line)
			}
		}
		sb.WriteString("\n}\n")
		return sb.String()
	}

	tc := *ptc
	tc.code = wrap(tc.code)
	tc.expected = wrap(tc.expected)

	tc.masked = MaskTransform(tc.masked, wrap)
	tc.lineCap = tc.lineCap + len("    ") // adjust the indent

	runPrinterTestCase(testSuite, &tc)
}

func runPrinterTestsuite(
	testSuite *testing.T,
	pred testRunnerFunc,
	testCases []printerTestCase,
	defaultLineCap int,
) {
	testSuite.Helper()

	for testIdx := range testCases {
		tc := testCases[testIdx]
		lineCap := tc.lineCap
		if lineCap <= 0 {
			lineCap = defaultLineCap
		}
		tc.lineCap = lineCap

		tc.pc(func() {
			testSuite.Helper()
			pred(testSuite, &tc)
		})
	}
}
