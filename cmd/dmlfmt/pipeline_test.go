// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/timoschenko/dmltools/pkg/fs/stream"

	"github.com/antlr4-go/antlr/v4"
)

func TestParseDmlStream(t *testing.T) {
	t.Parallel()

	input := antlr.NewInputStream("dml 1.4;\ndevice test_dev;\n")
	file, tree, listener := parseDmlStream(input, "test.dml")

	if file == nil || tree == nil || listener == nil {
		t.Fatalf(
			"expected non-nil parse results, got file=%v, tree=%v, listener=%v",
			file,
			tree,
			listener,
		)
	}
	if len(listener.Errors) > 0 {
		t.Errorf("unexpected syntax errors: %v", listener.Errors)
	}
}

func TestParseInput_Valid(t *testing.T) {
	t.Parallel()

	ui := NewUI(&bytes.Buffer{}, &bytes.Buffer{}, ColorNever)
	stats := &noopStatsCollector{}

	code := "dml 1.4;\ndevice dev;\n"
	buf := stream.NewHashedBufferFromString("dev.dml", code)

	job := InputJob{
		Path:      "dev.dml",
		OrigMode:  0644,
		Size:      len(code),
		Buffer:    buf,
		InputHash: buf.Hash(),
	}

	file, tree, earlyJob := parseInput(job, ui, stats)
	if earlyJob != nil {
		t.Fatalf("expected earlyJob == nil, got %+v", earlyJob)
	}
	if file == nil || tree == nil {
		t.Fatalf("expected non-nil file and tree")
	}
}

func TestParseInput_UnimplementedVersion(t *testing.T) {
	t.Parallel()

	ui := NewUI(&bytes.Buffer{}, &bytes.Buffer{}, ColorNever)
	stats := &noopStatsCollector{}

	code := "dml 1.2;\ndevice dev;\n"
	buf := stream.NewHashedBufferFromString("dev.dml", code)

	job := InputJob{
		Path:      "dev.dml",
		OrigMode:  0644,
		Size:      len(code),
		Buffer:    buf,
		InputHash: buf.Hash(),
	}

	file, tree, earlyJob := parseInput(job, ui, stats)
	if earlyJob == nil {
		t.Fatalf("expected earlyJob != nil for dml 1.2")
	}
	if !earlyJob.Ignored {
		t.Errorf("expected earlyJob.Ignored=true")
	}
	if earlyJob.IgnoreReason != "1.2" {
		t.Errorf("expected IgnoreReason='1.2', got %q", earlyJob.IgnoreReason)
	}
	if file != nil || tree != nil {
		t.Errorf("expected nil file and tree on early exit")
	}
}

func TestParseInput_SyntaxError(t *testing.T) {
	t.Parallel()

	var errBuf bytes.Buffer
	ui := NewUI(&bytes.Buffer{}, &errBuf, ColorNever)
	stats := &noopStatsCollector{}

	code := "dml 1.4;\ninvalid syntax !!! ;"
	buf := stream.NewHashedBufferFromString("bad.dml", code)

	job := InputJob{
		Path:      "bad.dml",
		OrigMode:  0644,
		Size:      len(code),
		Buffer:    buf,
		InputHash: buf.Hash(),
	}

	file, tree, earlyJob := parseInput(job, ui, stats)
	if earlyJob == nil {
		t.Fatalf("expected earlyJob != nil for syntax error")
	}
	if !earlyJob.SyntaxErrors {
		t.Errorf("expected earlyJob.SyntaxErrors=true")
	}
	if file != nil || tree != nil {
		t.Errorf("expected nil file and tree on early exit")
	}
	if !strings.Contains(errBuf.String(), "Syntax error") {
		t.Errorf("expected UI to report syntax error, got: %q", errBuf.String())
	}
}

func TestFormatAST(t *testing.T) {
	t.Parallel()

	input := antlr.NewInputStream("dml 1.4;\ndevice test_dev;\n")
	file, tree, _ := parseDmlStream(input, "test.dml")

	var buf bytes.Buffer
	buf.Grow(64)
	hash, err := formatAST(&buf, tree, file, 99, "printer panic")
	if err != nil {
		t.Fatalf("unexpected formatAST error: %v", err)
	}
	if buf.Len() == 0 {
		t.Errorf("expected non-empty output bytes in buf")
	}
	if hash == (stream.Hash{}) {
		t.Errorf("expected non-zero hash")
	}
}

func TestValidateFormatting(t *testing.T) {
	t.Parallel()

	t.Run("Success_Formatted", func(t *testing.T) {
		t.Parallel()
		input := antlr.NewInputStream("dml 1.4;\n\ndevice test_dev;\n")
		file, tree, listener := parseDmlStream(input, "test.dml")
		if len(listener.Errors) > 0 {
			t.Fatalf("unexpected parse errors: %v", listener.Errors)
		}

		var buf bytes.Buffer
		expectedHash, err := formatAST(
			&buf,
			tree,
			file,
			DefaultMaxLineCap,
			"printer panic",
		)
		if err != nil {
			t.Fatalf("unexpected formatAST error: %v", err)
		}

		err = validateFormatting(
			"test.dml",
			&buf,
			expectedHash,
			DefaultMaxLineCap,
		)
		if err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	})

	t.Run("Success_ComplexStructure", func(t *testing.T) {
		t.Parallel()
		complexDML := `dml 1.4;

device test_dev;

template my_tmpl {
	method get_val() -> (uint32) {
		return 42;
	}
}
`
		input := antlr.NewInputStream(complexDML)
		file, tree, listener := parseDmlStream(input, "complex.dml")
		if len(listener.Errors) > 0 {
			t.Fatalf("unexpected parse errors: %v", listener.Errors)
		}

		var buf bytes.Buffer
		expectedHash, err := formatAST(
			&buf,
			tree,
			file,
			DefaultMaxLineCap,
			"printer panic",
		)
		if err != nil {
			t.Fatalf("unexpected formatAST error: %v", err)
		}

		err = validateFormatting(
			"complex.dml",
			&buf,
			expectedHash,
			DefaultMaxLineCap,
		)
		if err != nil {
			t.Fatalf("unexpected validation error: %v", err)
		}
	})

	t.Run("SyntaxError", func(t *testing.T) {
		t.Parallel()
		formatted := []byte("dml 1.4;\n bad syntax here !!!")
		hasher := stream.NewFileHasher()
		_, _ = hasher.Write(formatted)
		expectedHash := hasher.Sum()

		err := validateFormatting(
			"bad.dml",
			bytes.NewBuffer(formatted),
			expectedHash,
			DefaultMaxLineCap,
		)
		if err == nil {
			t.Fatalf("expected error during syntax check")
		}
		expectedSubstr := "validating file bad.dml: syntax error during safety check"
		if !strings.Contains(err.Error(), expectedSubstr) {
			t.Errorf("expected %q, got %v", expectedSubstr, err)
		}
	})

	t.Run("NonIdempotent_RealReformatting", func(t *testing.T) {
		t.Parallel()
		// Passing unformatted code as the 'formatted' buffer:
		// Re-formatting it during validation will change whitespace and hash
		unformatted := []byte("dml   1.4  ;\n  device   dev  ;\n")
		hasher := stream.NewFileHasher()
		_, _ = hasher.Write(unformatted)
		unformattedHash := hasher.Sum()

		err := validateFormatting(
			"unformatted.dml",
			bytes.NewBuffer(unformatted),
			unformattedHash,
			DefaultMaxLineCap,
		)
		if err == nil {
			t.Fatalf("expected non-idempotent formatting error")
		}
		expectedSubstr := "validating file unformatted.dml: additional formatting detected (formatting is not idempotent)"
		if !strings.Contains(err.Error(), expectedSubstr) {
			t.Errorf("expected %q, got %v", expectedSubstr, err)
		}
	})

	t.Run("NonIdempotent_HashMismatch", func(t *testing.T) {
		t.Parallel()
		formatted := []byte("dml 1.4;\n\ndevice test_dev;\n")
		wrongHash := stream.Hash{FNV64: 1, Size: 2, CRC64: 3}

		err := validateFormatting(
			"mismatch.dml",
			bytes.NewBuffer(formatted),
			wrongHash,
			DefaultMaxLineCap,
		)
		if err == nil {
			t.Fatalf("expected error during idempotency check")
		}
		expectedSubstr := "validating file mismatch.dml: additional formatting detected (formatting is not idempotent)"
		if !strings.Contains(err.Error(), expectedSubstr) {
			t.Errorf("expected %q, got %v", expectedSubstr, err)
		}
	})

	t.Run("CustomLineCap", func(t *testing.T) {
		t.Parallel()
		input := antlr.NewInputStream("dml 1.4;\n\ndevice dev;\n")
		file, tree, _ := parseDmlStream(input, "cap.dml")

		var buf bytes.Buffer
		expectedHash, err := formatAST(&buf, tree, file, 80, "printer panic")
		if err != nil {
			t.Fatalf("unexpected formatAST error: %v", err)
		}

		err = validateFormatting("cap.dml", &buf, expectedHash, 80)
		if err != nil {
			t.Fatalf("unexpected validation error with lineCap=80: %v", err)
		}
	})
}
