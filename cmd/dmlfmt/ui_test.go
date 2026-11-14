// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"
)

func TestUIRouting(t *testing.T) {
	var out, err bytes.Buffer
	ui := NewUI(&out, &err, ColorNever)

	// Print -> out
	ui.Print("hello %s\n", "world")
	if out.String() != "hello world\n" || err.Len() != 0 {
		t.Errorf(
			"Print should write to out only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// Write -> out
	_, _ = ui.Write([]byte("raw code"))
	if out.String() != "raw code" || err.Len() != 0 {
		t.Errorf(
			"Write should write to out only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// Reformatted -> out
	ui.Reformatted("test.dml")
	if out.Len() == 0 || err.Len() != 0 {
		t.Errorf(
			"Reformatted should write to out only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// PrintError -> err
	ui.PrintError("Directory is not found")
	if err.Len() == 0 || out.Len() != 0 {
		t.Errorf(
			"PrintError should write to err only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// ReformattedStdin -> err
	ui.ReformattedStdin("-")
	if err.Len() == 0 || out.Len() != 0 {
		t.Errorf(
			"ReformattedStdin should write to err only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// WouldReformat -> out
	ui.WouldReformat("test.dml")
	if out.Len() == 0 || err.Len() != 0 {
		t.Errorf(
			"WouldReformat should write to out only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// Ignored -> out
	ui.Ignored("test.dml", "1.2")
	if out.Len() == 0 || err.Len() != 0 {
		t.Errorf(
			"Ignored should write to out only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// IgnoredReason -> out
	ui.IgnoredReason("test.dml", "file size exceeds 50 MiB limit")
	if out.Len() == 0 || err.Len() != 0 {
		t.Errorf(
			"IgnoredReason should write to out only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// SyntaxError -> err
	ui.SyntaxError("test.dml", 1, 2, "bad")
	if err.Len() == 0 || out.Len() != 0 {
		t.Errorf(
			"SyntaxError should write to err only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// PrintUsage -> out
	ui.PrintUsage()
	if out.Len() == 0 || err.Len() != 0 {
		t.Errorf(
			"PrintUsage should write to out only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// PrintErrorUsage -> err
	ui.PrintErrorUsage()
	if err.Len() == 0 || out.Len() != 0 {
		t.Errorf(
			"PrintErrorUsage should write to err only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// PrintVersion -> out
	ui.PrintVersion()
	if out.Len() == 0 || err.Len() != 0 {
		t.Errorf(
			"PrintVersion should write to out only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()
	err.Reset()

	// PrintNoSourcesError -> err
	ui.PrintNoSourcesError()
	expectedNoSrc := "Usage: dmlfmt [OPTIONS] SRC ...\n\nOne of 'SRC' or 'code' is required.\n"
	if err.String() != expectedNoSrc || out.Len() != 0 {
		t.Errorf(
			"PrintNoSourcesError should write to err only, got out=%q, err=%q",
			out.String(),
			err.String(),
		)
	}
}

func TestUIUncoloredOutput(t *testing.T) {
	var out, err bytes.Buffer
	ui := NewUI(&out, &err, ColorNever)

	ui.Reformatted("dev.dml")
	if out.String() != "reformatted dev.dml\n" {
		t.Errorf("unexpected output: %q", out.String())
	}
	out.Reset()

	ui.ReformattedStdin("-")
	if err.String() != "reformatted -\n" {
		t.Errorf("unexpected output: %q", err.String())
	}
	err.Reset()

	ui.WouldReformat("dev.dml")
	if out.String() != "would reformat dev.dml\n" || err.Len() != 0 {
		t.Errorf("unexpected output: out=%q err=%q", out.String(), err.String())
	}
	out.Reset()

	ui.Ignored("dev.dml", "1.2")
	if out.String() != "ignored dev.dml (dml 1.2)\n" || err.Len() != 0 {
		t.Errorf("unexpected output: out=%q err=%q", out.String(), err.String())
	}
	out.Reset()

	ui.IgnoredReason("dev.dml", "file size exceeds 50 MiB limit")
	if out.String() != "ignored dev.dml (file size exceeds 50 MiB limit)\n" ||
		err.Len() != 0 {
		t.Errorf(
			"unexpected output: out=%q err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()

	ui.IgnoredReason("dev.dml", "")
	if out.String() != "ignored dev.dml\n" || err.Len() != 0 {
		t.Errorf(
			"unexpected output: out=%q err=%q",
			out.String(),
			err.String(),
		)
	}
	out.Reset()

	ui.SyntaxError("dev.dml", 5, 10, "mismatched token")
	if err.String() != "dev.dml:5 column 10: Syntax error - mismatched token\n" {
		t.Errorf("unexpected output: %q", err.String())
	}
	err.Reset()

	ui.PrintError("Error reading file %s: %s", "dev.dml", "permission denied")
	if err.String() != "Error reading file dev.dml: permission denied\n" {
		t.Errorf("unexpected output: %q", err.String())
	}
	err.Reset()

	ui.PrintError("Error: %s is a directory", "src")
	if err.String() != "Error: src is a directory\n" {
		t.Errorf("unexpected output: %q", err.String())
	}
	err.Reset()

	ui.PrintError("Error: flag -%s is not implemented yet.", "diff")
	if err.String() != "Error: flag -diff is not implemented yet.\n" {
		t.Errorf("unexpected output: %q", err.String())
	}
	err.Reset()

	ui.PrintError("'SRC' and 'code' cannot be passed simultaneously.")
	if err.String() != "'SRC' and 'code' cannot be passed simultaneously.\n" {
		t.Errorf("unexpected output: %q", err.String())
	}
}

func TestUIColoredOutput(t *testing.T) {
	var out, err bytes.Buffer
	ui := NewUI(&out, &err, ColorAlways)

	if !ui.Out().IsColorEnabled() || !ui.Err().IsColorEnabled() {
		t.Errorf("expected IsColorEnabled() to be true for both streams")
	}

	ui.Reformatted("dev.dml")
	outStr := out.String()
	if !strings.Contains(outStr, ansiBoldGreen) ||
		!strings.Contains(outStr, "reformatted") ||
		!strings.Contains(outStr, ansiReset) {
		t.Errorf("expected colored reformatted message, got %q", outStr)
	}
	out.Reset()

	ui.WouldReformat("dev.dml")
	outStr = out.String()
	if !strings.Contains(outStr, ansiBoldBlue) ||
		!strings.Contains(outStr, "would reformat") ||
		!strings.Contains(outStr, ansiReset) {
		t.Errorf("expected colored would reformat message, got %q", outStr)
	}
	if err.Len() != 0 {
		t.Errorf("expected err to be empty, got %q", err.String())
	}
	out.Reset()

	ui.Ignored("dev.dml", "1.2")
	outStr = out.String()
	if !strings.Contains(outStr, ansiBoldYellow) ||
		!strings.Contains(outStr, "ignored") ||
		!strings.Contains(outStr, "(dml 1.2)") ||
		!strings.Contains(outStr, ansiReset) {
		t.Errorf("expected colored ignored message, got %q", outStr)
	}
	if err.Len() != 0 {
		t.Errorf("expected err to be empty, got %q", err.String())
	}
	out.Reset()

	ui.IgnoredReason("dev.dml", "file size exceeds 50 MiB limit")
	outStr = out.String()
	if !strings.Contains(outStr, ansiBoldYellow) ||
		!strings.Contains(outStr, "ignored") ||
		!strings.Contains(outStr, "(file size exceeds 50 MiB limit)") ||
		!strings.Contains(outStr, ansiReset) {
		t.Errorf("expected colored ignored reason message, got %q", outStr)
	}
	if err.Len() != 0 {
		t.Errorf("expected err to be empty, got %q", err.String())
	}
	out.Reset()

	ui.IgnoredReason("dev.dml", "")
	outStr = out.String()
	if !strings.Contains(outStr, ansiBoldYellow) ||
		!strings.Contains(outStr, "ignored") ||
		!strings.Contains(outStr, ansiReset) {
		t.Errorf(
			"expected colored ignored empty reason message, got %q",
			outStr,
		)
	}
	if err.Len() != 0 {
		t.Errorf("expected err to be empty, got %q", err.String())
	}
	out.Reset()

	ui.SyntaxError("dev.dml", 1, 2, "error msg")
	errStr := err.String()
	if !strings.Contains(errStr, ansiBoldRed) ||
		!strings.Contains(errStr, "Syntax error") ||
		!strings.Contains(errStr, ansiReset) {
		t.Errorf("expected colored syntax error message, got %q", errStr)
	}
	err.Reset()

	ui.PrintError("Error reading file %s: %s", "dev.dml", "boom")
	errStr = err.String()
	if !strings.Contains(errStr, ansiBoldRed) ||
		!strings.Contains(errStr, "Error") {
		t.Errorf("expected colored file read error message, got %q", errStr)
	}
	err.Reset()

	ui.PrintError("Error: %s is a directory", "dir")
	errStr = err.String()
	if !strings.Contains(errStr, ansiBoldRed) ||
		!strings.Contains(errStr, "Error:") {
		t.Errorf("expected colored directory error message, got %q", errStr)
	}
	err.Reset()

	ui.PrintError("Custom error message")
	errStr = err.String()
	if !strings.Contains(errStr, ansiBoldRed) ||
		!strings.Contains(errStr, "Custom error message") {
		t.Errorf("expected colored custom error message, got %q", errStr)
	}
	err.Reset()

	ui.PrintNoSourcesError()
	if strings.Contains(err.String(), "\033") {
		t.Errorf(
			"expected PrintNoSourcesError to be uncolored, got %q",
			err.String(),
		)
	}
}

func TestUISetColor(t *testing.T) {
	var out, err bytes.Buffer
	ui := NewUI(&out, &err, ColorNever)

	ui.Reformatted("file.dml")
	if strings.Contains(out.String(), "\033") {
		t.Errorf(
			"expected no escape codes when ColorNever, got %q",
			out.String(),
		)
	}
	out.Reset()

	ui.SetColor(true)
	if !ui.Out().IsColorEnabled() || !ui.Err().IsColorEnabled() {
		t.Errorf(
			"expected IsColorEnabled() to be true on both streams after SetColor(true)",
		)
	}
	ui.Reformatted("file.dml")
	if !strings.Contains(out.String(), "\033") {
		t.Errorf(
			"expected escape codes after SetColor(true), got %q",
			out.String(),
		)
	}
	out.Reset()

	ui.SetColor(false)
	if ui.Out().IsColorEnabled() || ui.Err().IsColorEnabled() {
		t.Errorf(
			"expected IsColorEnabled() to be false on both streams after SetColor(false)",
		)
	}
	ui.Reformatted("file.dml")
	if strings.Contains(out.String(), "\033") {
		t.Errorf(
			"expected no escape codes after SetColor(false), got %q",
			out.String(),
		)
	}
}

func TestUIColorAutoWithBuffers(t *testing.T) {
	var out, err bytes.Buffer
	ui := NewUI(&out, &err, ColorAuto)

	// Since bytes.Buffer is not a character device (TTY), ColorAuto should disable color
	if ui.Out().IsColorEnabled() || ui.Err().IsColorEnabled() {
		t.Errorf(
			"expected IsColorEnabled() to be false for non-terminal buffers in ColorAuto mode",
		)
	}

	ui.Reformatted("test.dml")
	if strings.Contains(out.String(), "\033") {
		t.Errorf(
			"expected plain text for bytes.Buffer in ColorAuto mode, got %q",
			out.String(),
		)
	}

	ui.PrintError("Error: something failed")
	if strings.Contains(err.String(), "\033") {
		t.Errorf(
			"expected plain text for bytes.Buffer in ColorAuto mode, got %q",
			err.String(),
		)
	}

	// Also verify SetColorMode updates streams
	ui.SetColorMode(ColorAlways)
	if !ui.Out().IsColorEnabled() || !ui.Err().IsColorEnabled() {
		t.Errorf(
			"expected IsColorEnabled() to be true after SetColorMode(ColorAlways)",
		)
	}
	ui.SetColorMode(ColorAuto)
	if ui.Out().IsColorEnabled() || ui.Err().IsColorEnabled() {
		t.Errorf(
			"expected IsColorEnabled() to be false after SetColorMode(ColorAuto) on non-terminal buffers",
		)
	}
}

func TestIsColorSupportedEnvironment(t *testing.T) {
	// A non-file writer is never a supported terminal
	var buf bytes.Buffer
	if isColorSupported(&buf) {
		t.Errorf("expected isColorSupported(&bytes.Buffer) == false")
	}

	// When NO_COLOR is set, color is always disabled
	t.Setenv("NO_COLOR", "1")
	if isColorSupported(os.Stdout) {
		t.Errorf("expected isColorSupported to be false when NO_COLOR is set")
	}

	// When TERM=dumb, color is disabled
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "dumb")
	if isColorSupported(os.Stdout) {
		t.Errorf("expected isColorSupported to be false when TERM=dumb")
	}
}

func TestUISyntaxErrorListener(t *testing.T) {
	var out, err bytes.Buffer
	ui := NewUI(&out, &err, ColorNever)

	listener := ui.NewSyntaxErrorListener("sample.dml")
	if listener == nil || listener.Handler == nil {
		t.Fatalf("expected initialized listener with Handler")
	}

	listener.Handler(dmlex.SyntaxErrorInfo{
		Line:   10,
		Column: 4,
		Msg:    "unexpected token",
	})

	expected := "sample.dml:10 column 4: Syntax error - unexpected token\n"
	if err.String() != expected {
		t.Errorf("expected %q, got %q", expected, err.String())
	}
}

func TestDefaultUIStreams(t *testing.T) {
	ui := NewUI(nil, nil, ColorNever)
	if ui.Out().Unwrap() != os.Stdout {
		t.Errorf("expected default Out to be os.Stdout")
	}
	if ui.Err().Unwrap() != os.Stderr {
		t.Errorf("expected default Err to be os.Stderr")
	}
}

func TestUIConcurrentPrinting(t *testing.T) {
	var out, err bytes.Buffer
	ui := NewUI(&out, &err, ColorAlways)

	var wg sync.WaitGroup
	workers := 10
	iterations := 10

	for i := range workers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range iterations {
				ui.Print("worker %d print %d\n", id, j)
				ui.PrintError("worker %d error %d", id, j)
				ui.Reformatted(fmt.Sprintf("file_%d_%d.dml", id, j))
				ui.ReformattedStdin(fmt.Sprintf("stdin_%d_%d", id, j))
				ui.WouldReformat(fmt.Sprintf("check_%d_%d.dml", id, j))
				ui.SyntaxError(
					fmt.Sprintf("syntax_%d_%d.dml", id, j),
					j,
					j,
					"unexpected token",
				)
				rawMsg := fmt.Sprintf("raw_%d_%d\n", id, j)
				_, _ = ui.Write([]byte(rawMsg))
			}
		}(i)
	}

	wg.Wait()

	if out.Len() == 0 {
		t.Errorf("expected out buffer to have content")
	}
	if err.Len() == 0 {
		t.Errorf("expected err buffer to have content")
	}
}

func TestUIVerbose(t *testing.T) {
	t.Parallel()

	t.Run("NoColor", func(t *testing.T) {
		var out, err bytes.Buffer
		ui := NewUI(&out, &err, ColorNever)

		ui.Processing("dev.dml")
		expected := "Processing dev.dml...\n"
		if out.String() != expected || err.Len() != 0 {
			t.Errorf(
				"Processing: got out=%q, err=%q, expected out=%q",
				out.String(),
				err.String(),
				expected,
			)
		}
		out.Reset()

		ui.IgnoredFilter("sub", true, "--exclude")
		expected = "sub directory ignored: matches the --exclude regular expression\n"
		if out.String() != expected || err.Len() != 0 {
			t.Errorf(
				"IgnoredFilter (dir): got out=%q, expected %q",
				out.String(),
				expected,
			)
		}
		out.Reset()

		ui.IgnoredFilter("dev.dml", false, "--extend-exclude")
		expected = "dev.dml ignored: matches the --extend-exclude regular expression\n"
		if out.String() != expected || err.Len() != 0 {
			t.Errorf(
				"IgnoredFilter (file): got out=%q, expected %q",
				out.String(),
				expected,
			)
		}
		out.Reset()

		ui.Unchanged("dev.dml")
		expected = "dev.dml already well formatted, good job.\n"
		if out.String() != expected || err.Len() != 0 {
			t.Errorf(
				"Unchanged: got out=%q, err=%q, expected out=%q",
				out.String(),
				err.String(),
				expected,
			)
		}
	})

	t.Run("ColorAlways", func(t *testing.T) {
		var out, err bytes.Buffer
		ui := NewUI(&out, &err, ColorAlways)

		ui.Processing("dev.dml")
		expected := "\033[90mProcessing dev.dml...\033[0m\n"
		if out.String() != expected || err.Len() != 0 {
			t.Errorf(
				"Processing colored: got out=%q, expected %q",
				out.String(),
				expected,
			)
		}
		out.Reset()

		ui.IgnoredFilter("sub", true, "--exclude")
		expected = "\033[90msub directory ignored: matches the --exclude regular expression\033[0m\n"
		if out.String() != expected || err.Len() != 0 {
			t.Errorf(
				"IgnoredFilter colored (dir): got out=%q, expected %q",
				out.String(),
				expected,
			)
		}
		out.Reset()

		ui.IgnoredFilter("dev.dml", false, "--force-exclude")
		expected = "\033[90mdev.dml ignored: matches the --force-exclude regular expression\033[0m\n"
		if out.String() != expected || err.Len() != 0 {
			t.Errorf(
				"IgnoredFilter colored (file): got out=%q, expected %q",
				out.String(),
				expected,
			)
		}
		out.Reset()

		ui.Unchanged("dev.dml")
		expected = "\033[90mdev.dml already well formatted, good job.\033[0m\n"
		if out.String() != expected || err.Len() != 0 {
			t.Errorf(
				"Unchanged colored: got out=%q, expected %q",
				out.String(),
				expected,
			)
		}
	})
}
