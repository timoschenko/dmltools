// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

//go:build !dmltools.ndebug

package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/timoschenko/dmltools/pkg/printer/layout"
)

func TestFormatDML_Success(t *testing.T) {
	input := `dml 1.4;
device test_dev;
`
	res := FormatDML(input, 80)
	if res.Error != "" {
		t.Fatalf("expected no error, got: %s", res.Error)
	}
	if res.SyntaxError {
		t.Fatalf("expected SyntaxError to be false")
	}
	if !strings.Contains(res.Result, "dml 1.4;") ||
		!strings.Contains(res.Result, "device test_dev;") {
		t.Fatalf("unexpected format result:\n%s", res.Result)
	}
}

func TestFormatDML_Trace(t *testing.T) {
	input := `dml 1.4;
device test_dev;
`
	res := FormatDML(input, 80)
	if res.Error != "" {
		t.Fatalf("expected no error, got: %s", res.Error)
	}
	if res.Trace == "" {
		t.Fatalf("expected res.Trace to be non-empty")
	}
	if res.TraceCount <= 0 {
		t.Fatalf("expected res.TraceCount > 0, got %d", res.TraceCount)
	}

	var entries []layout.TraceEntry
	if err := json.Unmarshal([]byte(res.Trace), &entries); err != nil {
		t.Fatalf("failed to unmarshal trace JSON: %v", err)
	}
	if len(entries) != res.TraceCount {
		t.Errorf("expected %d entries, got %d", res.TraceCount, len(entries))
	}
	if len(entries) > 0 {
		first := entries[0]
		if first.Name != "TokenFormatDecision" {
			t.Errorf(
				"expected first entry name 'TokenFormatDecision', got %q",
				first.Name,
			)
		}
		if first.Category != "formatter" {
			t.Errorf("expected category 'formatter', got %q", first.Category)
		}
		if first.Args == nil {
			t.Errorf("expected Args to be non-nil")
		}
	}
}

func TestFormatDML_TraceDisabled(t *testing.T) {
	input := `dml 1.4;
device test_dev;
`
	res := FormatDML(input, 80, false)
	if res.Error != "" {
		t.Fatalf("expected no error, got: %s", res.Error)
	}
	if res.Trace != "" {
		t.Fatalf(
			"expected res.Trace to be empty when disabled, got: %s",
			res.Trace,
		)
	}
	if res.TraceCount != 0 {
		t.Fatalf("expected res.TraceCount to be 0, got %d", res.TraceCount)
	}
}

func TestFormatDML_TraceOnSyntaxError(t *testing.T) {
	input := `dml 1.4;
invalid statement;
`
	res := FormatDML(input, 80)
	if res.Error == "" {
		t.Fatalf("expected syntax error, got none")
	}
	if res.Trace != "" {
		t.Fatalf("expected no trace on syntax error, got: %s", res.Trace)
	}
	if res.TraceCount != 0 {
		t.Fatalf("expected res.TraceCount 0, got %d", res.TraceCount)
	}
}

func TestFormatDML_SyntaxError(t *testing.T) {
	input := `dml 1.4;
device test_dev;
invalid statement here;
`
	res := FormatDML(input, 80)
	if res.Error == "" {
		t.Fatalf("expected an error, got none")
	}
	if !res.SyntaxError {
		t.Fatalf("expected SyntaxError to be true")
	}
	if res.Line != 3 {
		t.Errorf("expected error line 3, got: %d", res.Line)
	}
	if res.Column != 0 {
		t.Errorf("expected error column 0, got: %d", res.Column)
	}
	if !strings.HasPrefix(res.Error, "line 3:0:") {
		t.Errorf("expected error to start with 'line 3:0:', got: %s", res.Error)
	}
}

func TestFormatDML_SyntaxError_MidLine(t *testing.T) {
	input := `dml 1.4;
device dev {
    bad token;
}
`
	res := FormatDML(input, 80)
	if res.Error == "" {
		t.Fatalf("expected error, got none")
	}
	if !res.SyntaxError {
		t.Fatalf("expected SyntaxError true")
	}
	if res.Line != 2 {
		t.Errorf("expected error line 2, got: %d", res.Line)
	}
	if res.Column <= 0 {
		t.Errorf("expected error column > 0, got: %d", res.Column)
	}
}

func TestFormatDML_LexerError(t *testing.T) {
	input := `dml 1.4;
@
`
	res := FormatDML(input, 80)
	if res.Error == "" {
		t.Fatalf("expected error for '@', got none")
	}
	if !res.SyntaxError {
		t.Fatalf("expected SyntaxError true")
	}
	if res.Line != 2 {
		t.Errorf("expected error line 2, got: %d", res.Line)
	}
}

func TestFormatDML_Groups(t *testing.T) {
	input := `dml 1.4;
// comment 1
device dev;
/* comment 2 */
method foo(int a, int b) {
}
`
	res := FormatDML(input, 80)
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	t.Logf("Tokens count: %d, Trace count: %d", len(res.Tokens), res.TraceCount)
	hasGroup := false
	for i, tok := range res.Tokens {
		t.Logf(
			"Token[%d]: text=%q type=%d group=%q",
			i,
			tok.Text,
			tok.Type,
			tok.GroupPath,
		)
		if tok.GroupPath != "" {
			hasGroup = true
		}
	}
	var entries []layout.TraceEntry
	if err := json.Unmarshal([]byte(res.Trace), &entries); err != nil {
		t.Fatalf("json.Unmarshal() - %v", err)
	}
	for i, e := range entries {
		t.Logf(
			"Trace[%d]: Right=%q RightGroup=%q LeftGroup=%q",
			i,
			e.Args.Right.Text,
			e.Args.Right.GroupPath,
			e.Args.Left.GroupPath,
		)
	}
	if !hasGroup {
		t.Errorf("expected at least one token to have GroupPath")
	}
}
