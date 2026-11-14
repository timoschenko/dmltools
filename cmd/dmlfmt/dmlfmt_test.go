// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"fmt"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestParseArgsDefaults(t *testing.T) {
	var stderr bytes.Buffer
	cfg, err := parseArgs([]string{"dev.dml"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.LineLength != DefaultMaxLineCap {
		t.Errorf(
			"expected LineLength %d, got %d",
			DefaultMaxLineCap,
			cfg.LineLength,
		)
	}
	if !cfg.Safe {
		t.Errorf("expected Safe=true by default")
	}
	if cfg.Fast {
		t.Errorf("expected Fast=false by default")
	}
	if cfg.Include != DefaultIncludePattern {
		t.Errorf(
			"expected Include %q, got %q",
			DefaultIncludePattern,
			cfg.Include,
		)
	}
	if cfg.Exclude != DefaultExcludePattern {
		t.Errorf(
			"expected Exclude %q, got %q",
			DefaultExcludePattern,
			cfg.Exclude,
		)
	}
	expectedDefaultWorkers := min(2, runtime.NumCPU())
	if cfg.Workers != expectedDefaultWorkers {
		t.Errorf(
			"expected Workers %d, got %d",
			expectedDefaultWorkers,
			cfg.Workers,
		)
	}
	if len(cfg.Sources) != 1 || cfg.Sources[0] != "dev.dml" {
		t.Errorf("expected Sources [\"dev.dml\"], got %v", cfg.Sources)
	}
}

func TestParseArgsShortAndLongFlags(t *testing.T) {
	var stderr bytes.Buffer
	args := []string{
		"-l", "80",
		"-t", "1.4",
		"-x",
		"--unstable",
		"--enable-unstable-feature", "experimental_wrapping",
		"--check",
		"--diff",
		"--color",
		"--line-ranges", "10-20",
		"--fast",
		"--exclude", "custom_exclude",
		"--extend-exclude", "more_exclude",
		"--force-exclude", "forced_exclude",
		"--stdin-filename", "stdin.dml",
		"--include", `\.dml$`,
		"-W", "4",
		"-v",
		"--no-cache",
		"--nerd-stats",
		"pos1.dml",
		"pos2.dml",
	}

	cfg, err := parseArgs(args, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.LineLength != 80 {
		t.Errorf("expected LineLength 80, got %d", cfg.LineLength)
	}
	if cfg.TargetVersion != "1.4" {
		t.Errorf("expected TargetVersion 1.4, got %q", cfg.TargetVersion)
	}
	if !cfg.SkipSourceFirstLine {
		t.Errorf("expected SkipSourceFirstLine=true")
	}
	if !cfg.Unstable || !cfg.Preview {
		t.Errorf("expected Unstable=true and Preview=true (implied)")
	}
	if cfg.UnstableFeature != "experimental_wrapping" {
		t.Errorf(
			"expected UnstableFeature 'experimental_wrapping', got %q",
			cfg.UnstableFeature,
		)
	}
	if !cfg.Check || !cfg.Diff || !cfg.Color {
		t.Errorf("expected Check, Diff, and Color to be true")
	}
	if cfg.LineRanges != "10-20" {
		t.Errorf("expected LineRanges '10-20', got %q", cfg.LineRanges)
	}
	if !cfg.Fast || cfg.Safe {
		t.Errorf("expected Fast=true and Safe=false")
	}
	if cfg.Workers != 4 {
		t.Errorf("expected Workers 4, got %d", cfg.Workers)
	}
	if !cfg.Verbose || cfg.Quiet {
		t.Errorf("expected Verbose=true and Quiet=false")
	}
	if !cfg.NoCache {
		t.Errorf("expected NoCache=true")
	}
	if !cfg.NerdStats {
		t.Errorf("expected NerdStats=true")
	}
	if len(cfg.Sources) != 2 {
		t.Fatalf("expected 2 positional sources, got %v", cfg.Sources)
	}
}

func TestParseArgsQuietPrecedence(t *testing.T) {
	var stderr bytes.Buffer
	cfg, err := parseArgs([]string{"-q", "-v", "dev.dml"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Quiet {
		t.Errorf("expected Quiet=true")
	}
	if cfg.Verbose {
		t.Errorf("expected Verbose=false when Quiet is set")
	}

	stderr.Reset()
	cfg, err = parseArgs([]string{"-v", "--quiet", "dev.dml"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Quiet {
		t.Errorf("expected Quiet=true")
	}
	if cfg.Verbose {
		t.Errorf("expected Verbose=false when Quiet is set")
	}
}

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--help"}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf("expected exit code %d, got %d", ExitCodeOK, code)
	}
	helpOut := stdout.String()
	expectedFlags := []string{
		"-c, --code TEXT",
		"-l, --line-length INTEGER",
		// "-t, --target-version [1.2|1.4]",
		"--check",
		// "--diff",
		// "--line-ranges START-END",
		"--fast / --safe",
		"--include TEXT",
		"-W, --workers INTEGER",
		"-v, --verbose",
		"--nerd-stats",
	}
	for _, flagStr := range expectedFlags {
		if !strings.Contains(helpOut, flagStr) {
			t.Errorf("help output missing %q", flagStr)
		}
	}
	if !strings.Contains(helpOut, "which can be") ||
		!strings.Contains(helpOut, "silenced by 2>/dev/null)") {
		t.Errorf("help output missing corrected quiet description")
	}
}

func TestUnsupportedAndUnimplementedFlagsExit123(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"unsupported flag --unknown", []string{"--unknown"}},
		{"unsupported flag -S", []string{"-S"}},
		{"unsupported flag -C", []string{"-C"}},
		{
			"unsupported flag --required-version",
			[]string{"--required-version", "0.1.0"},
		},
		{"unsupported flag --pyi", []string{"--pyi"}},
		{"unimplemented -t", []string{"-t", "1.4"}},
		{"unimplemented --target-version", []string{"--target-version", "1.4"}},
		{"unimplemented -x", []string{"-x"}},
		{
			"unimplemented --skip-source-first-line",
			[]string{"--skip-source-first-line"},
		},
		{"unimplemented --preview", []string{"--preview"}},
		{"unimplemented --unstable", []string{"--unstable"}},
		{
			"unimplemented --enable-unstable-feature",
			[]string{"--enable-unstable-feature", "foo"},
		},
		{"unimplemented --diff", []string{"--diff"}},
		{"unimplemented --line-ranges", []string{"--line-ranges", "1-10"}},
		{"unsupported flag --config", []string{"--config", "cfg.toml"}},
		{"unimplemented --no-cache", []string{"--no-cache"}},
		{"unimplemented --cpuprofile", []string{"--cpuprofile", "cpu.prof"}},
		{"unimplemented --memprofile", []string{"--memprofile", "mem.prof"}},
		{"unsupported flag -file", []string{"-file", "f1.dml"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(tc.args, &stdout, &stderr)
			if code != ExitCodeInternalError {
				t.Errorf(
					"expected exit code %d, got %d",
					ExitCodeInternalError,
					code,
				)
			}
			if stderr.Len() == 0 {
				t.Errorf(
					"expected error message on stderr for %v, got empty stderr",
					tc.args,
				)
			}
		})
	}
}

func TestRunFile(t *testing.T) {
	var stdout, stderr bytes.Buffer

	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "test.dml")
	// Write unformatted code
	if err := os.WriteFile(
		dmlPath,
		[]byte("dml   1.4  ;\n  device   test_dev  ;\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write temp dml file: %v", err)
	}

	code := run([]string{dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK for positional argument, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedStdout := "reformatted " + dmlPath + "\n1 file reformatted.\n"
	if stdout.String() != expectedStdout {
		t.Errorf("expected stdout %q, got %q", expectedStdout, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// Verify formatted content
	content, err := os.ReadFile(filepath.Clean(dmlPath))
	if err != nil {
		t.Fatalf("failed to read reformatted file: %v", err)
	}
	expectedContent := "dml 1.4;\ndevice test_dev;\n"
	if string(content) != expectedContent {
		t.Errorf(
			"expected formatted content %q, got %q",
			expectedContent,
			string(content),
		)
	}

	// Verify permissions preserved (0600)
	info, err := os.Stat(dmlPath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %o", info.Mode().Perm())
	}

	// Run again: unchanged file outputs 1 file left unchanged and returns ExitCodeOK
	stdout.Reset()
	stderr.Reset()
	code = run([]string{dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK for unchanged file, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedUnchangedStdout := "1 file left unchanged.\n"
	if stdout.String() != expectedUnchangedStdout {
		t.Errorf(
			"expected stdout %q for unchanged file, got %q",
			expectedUnchangedStdout,
			stdout.String(),
		)
	}

	// Syntax error in file exits 123 and leaves file untouched
	badPath := filepath.Join(tmpDir, "bad.dml")
	badContent := []byte("dml 1.4; @@@ invalid")
	if err := os.WriteFile(badPath, badContent, 0600); err != nil {
		t.Fatalf("failed to write temp bad dml file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{badPath}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Errorf(
			"expected ExitCodeInternalError for syntax error, got %d",
			code,
		)
	}
	afterBadContent, err := os.ReadFile(filepath.Clean(badPath))
	if err != nil {
		t.Fatalf("failed to read bad file: %v", err)
	}
	if !bytes.Equal(afterBadContent, badContent) {
		t.Errorf(
			"syntax error file was modified: got %q, want %q",
			afterBadContent,
			badContent,
		)
	}

	// Directory argument traverses and reformats files in directory
	subDir := t.TempDir()
	subDml := filepath.Join(subDir, "dir_test.dml")
	if err := os.WriteFile(
		subDml,
		[]byte("dml   1.4  ;\n  device   dir_dev  ;\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write sub dml file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{subDir}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK for directory argument, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if !strings.Contains(stdout.String(), "reformatted "+subDml) {
		t.Errorf(
			"expected stdout to mention reformatted %s, got: %q",
			subDml,
			stdout.String(),
		)
	}

	// Non-existent file exits 123
	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{filepath.Join(tmpDir, "nonexistent.dml")},
		&stdout,
		&stderr,
	)
	if code != ExitCodeInternalError {
		t.Errorf(
			"expected ExitCodeInternalError for non-existent file, got %d",
			code,
		)
	}

	// No arguments returns error requiring SRC or code
	stdout.Reset()
	stderr.Reset()
	code = run([]string{}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Errorf(
			"expected ExitCodeInternalError when no inputs given, got %d",
			code,
		)
	}
	expectedNoSrc := "Usage: dmlfmt [OPTIONS] SRC ...\n\nOne of 'SRC' or 'code' is required.\n"
	if stderr.String() != expectedNoSrc {
		t.Errorf("expected stderr %q, got %q", expectedNoSrc, stderr.String())
	}

	// Multiple files: reformats all files
	multiPath1 := filepath.Join(tmpDir, "m1.dml")
	multiPath2 := filepath.Join(tmpDir, "m2.dml")
	if err := os.WriteFile(
		multiPath1,
		[]byte("dml   1.4  ;\n  device   m1  ;\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write m1: %v", err)
	}
	if err := os.WriteFile(
		multiPath2,
		[]byte("dml   1.4  ;\n  device   m2  ;\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write m2: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{multiPath1, multiPath2}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK for multiple files, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	outStr := stdout.String()
	if !strings.Contains(outStr, "reformatted "+multiPath1+"\n") ||
		!strings.Contains(outStr, "reformatted "+multiPath2+"\n") {
		t.Errorf(
			"expected stdout to contain both reformatted files, got %q",
			outStr,
		)
	}

	// Multiple files: stops on first error
	stdout.Reset()
	stderr.Reset()
	code = run([]string{badPath, multiPath1}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Errorf(
			"expected ExitCodeInternalError when first file has error, got %d",
			code,
		)
	}
}

func TestRunFile_LineLengthFlags(t *testing.T) {
	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "linelen.dml")
	if err := os.WriteFile(
		dmlPath,
		[]byte("dml 1.4;\ndevice test_dev;\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-l", "80", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Errorf(
			"expected ExitCodeOK with -l 80, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{"--line-length", "100", dmlPath},
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Errorf(
			"expected ExitCodeOK with --line-length 100, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
}

func TestRunFile_NerdStats(t *testing.T) {
	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "nerd.dml")
	content := "dml   1.4;\n\n// sample comment\ndevice   test_dev;\n"
	if err := os.WriteFile(dmlPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"--nerd-stats", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK with --nerd-stats, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	errOut := stderr.String()
	expectedSections := []string{
		"--- Nerd Stats ---",
		"Total read bytes:",
		"Total processed tokens:",
		"Total processed lines:",
		"Total time spent:",
		"Total IO time spent:",
		"walk",
		"read",
		"write",
		"Total format time spent:",
		"tokenize",
		"printing",
		"validate",
		"Throughput:",
		"lines/ms",
		"tokens/ms",
		"GC stats:",
		"allocs",
		"frees",
		"GC cycles",
		"Memory consumption:",
		"final",
		"peak",
		"------------------",
	}
	for _, sec := range expectedSections {
		if !strings.Contains(errOut, sec) {
			t.Errorf("expected stderr to contain %q, but got:\n%s", sec, errOut)
		}
	}

	expectedStdout := "reformatted " + dmlPath + "\n1 file reformatted.\n"
	if stdout.String() != expectedStdout {
		t.Errorf("expected stdout %q, got %q", expectedStdout, stdout.String())
	}
}

func TestRunFile_NerdStats_OnError(t *testing.T) {
	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "syntax_err.dml")
	if err := os.WriteFile(
		dmlPath,
		[]byte("dml invalid syntax ;;;;\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"--nerd-stats", dmlPath}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError on syntax error, got %d", code)
	}

	errOut := stderr.String()
	if !strings.Contains(errOut, "--- Nerd Stats ---") ||
		!strings.Contains(errOut, "------------------") {
		t.Errorf(
			"expected nerd stats printed on stderr on error, got stderr:\n%s",
			errOut,
		)
	}
}

func TestNerdStatsCollector_Calculations(t *testing.T) {
	collector := newNerdStatsCollector()
	collector.RecordReadBytes(1024)
	collector.RecordReadBytes(2048)

	// Test formatReadBytes
	if s := formatReadBytes(500); s != "500 bytes" {
		t.Errorf("expected '500 bytes', got %q", s)
	}
	if s := formatReadBytes(1024); s != "1024 bytes (1.00 KiB)" {
		t.Errorf("expected '1024 bytes (1.00 KiB)', got %q", s)
	}

	// Test formatBytes
	if s := formatBytes(512); s != "512 B" {
		t.Errorf("expected '512 B', got %q", s)
	}
	if s := formatBytes(1024 * 1024 * 2); s != "2.00 MiB" {
		t.Errorf("expected '2.00 MiB', got %q", s)
	}

	// Test empty tokens
	collector.RecordTokensAndLines(nil)
	if collector.totalTokens != 0 || collector.totalLines != 0 {
		t.Errorf("expected 0 tokens and lines for nil tokens")
	}

	var buf bytes.Buffer
	collector.PrintSummary(&buf)
	res := buf.String()
	if !strings.Contains(
		res,
		"Total read bytes:        3072 bytes (3.00 KiB)",
	) {
		t.Errorf("expected 3072 bytes in output, got:\n%s", res)
	}

	// Test IO timings formatting
	ioCollector := newNerdStatsCollector()
	ioCollector.RecordWalkIo(100 * time.Millisecond)
	ioCollector.RecordReadIo(200 * time.Millisecond)
	ioCollector.RecordWriteIo(700 * time.Millisecond)

	var ioBuf bytes.Buffer
	ioCollector.PrintSummary(&ioBuf)
	ioRes := ioBuf.String()
	if !strings.Contains(
		ioRes,
		"Total IO time spent:     1.00 s (walk 10.0%, read 20.0%, write 70.0%)",
	) {
		t.Errorf(
			"expected 'Total IO time spent:     1.00 s (walk 10.0%%, read 20.0%%, write 70.0%%)', got:\n%s",
			ioRes,
		)
	}

	// Test zero IO timings
	zeroCollector := newNerdStatsCollector()
	var zeroBuf bytes.Buffer
	zeroCollector.PrintSummary(&zeroBuf)
	zeroRes := zeroBuf.String()
	if !strings.Contains(
		zeroRes,
		"Total IO time spent:     0.00 s (walk 0.0%, read 0.0%, write 0.0%)",
	) {
		t.Errorf(
			"expected 'Total IO time spent:     0.00 s (walk 0.0%%, read 0.0%%, write 0.0%%)', got:\n%s",
			zeroRes,
		)
	}

	// Test format timings formatting
	fmtCollector := newNerdStatsCollector()
	fmtCollector.RecordTokenize(100 * time.Millisecond)
	fmtCollector.RecordPrinting(800 * time.Millisecond)
	fmtCollector.RecordValidate(100 * time.Millisecond)

	var fmtBuf bytes.Buffer
	fmtCollector.PrintSummary(&fmtBuf)
	fmtRes := fmtBuf.String()
	if !strings.Contains(
		fmtRes,
		"Total format time spent: 1.00 s (tokenize 10.0%, printing 80.0%, validate 10.0%)",
	) {
		t.Errorf(
			"expected 'Total format time spent: 1.00 s (tokenize 10.0%%, printing 80.0%%, validate 10.0%%)', got:\n%s",
			fmtRes,
		)
	}

	// Test zero format timings
	zeroFmtCollector := newNerdStatsCollector()
	var zeroFmtBuf bytes.Buffer
	zeroFmtCollector.PrintSummary(&zeroFmtBuf)
	zeroFmtRes := zeroFmtBuf.String()
	if !strings.Contains(
		zeroFmtRes,
		"Total format time spent: 0.00 s (tokenize 0.0%, printing 0.0%, validate 0.0%)",
	) {
		t.Errorf(
			"expected 'Total format time spent: 0.00 s (tokenize 0.0%%, printing 0.0%%, validate 0.0%%)', got:\n%s",
			zeroFmtRes,
		)
	}
}

func TestRunVersion(t *testing.T) {
	for _, flag := range []string{"--version", "-version"} {
		t.Run(flag, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run([]string{flag}, &stdout, &stderr)
			if code != ExitCodeOK {
				t.Fatalf(
					"expected exit code %d, got %d. stderr: %s",
					ExitCodeOK,
					code,
					stderr.String(),
				)
			}
			out := stdout.String()
			if !strings.HasPrefix(out, "dmlfmt version ") {
				t.Errorf(
					"expected version output to start with 'dmlfmt version ', got: %q",
					out,
				)
			}
			if !strings.Contains(out, "License:") ||
				!strings.Contains(out, "GPL-2.0") {
				t.Errorf(
					"expected license information in version output, got: %q",
					out,
				)
			}
		})
	}
}

func TestParseArgs_CodeAndStdinFlags(t *testing.T) {
	var stderr bytes.Buffer

	cfg, err := parseArgs([]string{"-c", "dml 1.4;"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.CodeSet || cfg.Code != "dml 1.4;" {
		t.Errorf(
			"expected CodeSet=true and Code 'dml 1.4;', got %v, %q",
			cfg.CodeSet,
			cfg.Code,
		)
	}

	cfg, err = parseArgs([]string{"--code", "dml 1.4;"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.CodeSet || cfg.Code != "dml 1.4;" {
		t.Errorf(
			"expected CodeSet=true and Code 'dml 1.4;', got %v, %q",
			cfg.CodeSet,
			cfg.Code,
		)
	}

	cfg, err = parseArgs([]string{"-c", ""}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.CodeSet || cfg.Code != "" {
		t.Errorf(
			"expected CodeSet=true and Code '', got %v, %q",
			cfg.CodeSet,
			cfg.Code,
		)
	}

	cfg, err = parseArgs([]string{"--stdin-filename", "foo.dml"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.StdinFilename != "foo.dml" {
		t.Errorf("expected StdinFilename 'foo.dml', got %q", cfg.StdinFilename)
	}
}

func TestRun_ArgumentValidation(t *testing.T) {
	// Both -c and a file argument
	var stdout, stderr bytes.Buffer
	code := run([]string{"-c", "dml 1.4;", "test.dml"}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Errorf("expected ExitCodeInternalError, got %d", code)
	}
	if !strings.Contains(
		stderr.String(),
		"'SRC' and 'code' cannot be passed simultaneously",
	) {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}

	// Both -c and -
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-c", "dml 1.4;", "-"}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Errorf("expected ExitCodeInternalError, got %d", code)
	}
	if !strings.Contains(
		stderr.String(),
		"'SRC' and 'code' cannot be passed simultaneously",
	) {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}

	// Neither SRC nor code produces error
	stdout.Reset()
	stderr.Reset()
	code = run([]string{}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Errorf(
			"expected ExitCodeInternalError for empty arguments, got %d",
			code,
		)
	}
	expectedNoSrcErr := "Usage: dmlfmt [OPTIONS] SRC ...\n\nOne of 'SRC' or 'code' is required.\n"
	if stderr.String() != expectedNoSrcErr {
		t.Errorf(
			"expected stderr %q, got %q",
			expectedNoSrcErr,
			stderr.String(),
		)
	}

	// Mixing - with another file
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-", "file.dml"}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Errorf("expected ExitCodeInternalError, got %d", code)
	}
	if !strings.Contains(
		stderr.String(),
		"Cannot mix stdin ('-') with other file paths",
	) {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}
}

func TestRun_Code(t *testing.T) {
	// Format unformatted code
	var stdout, stderr bytes.Buffer
	code := run(
		[]string{"-c", "dml   1.4  ;\n  device   dev  ;\n"},
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expected := "dml 1.4;\ndevice dev;\n"
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr for -c, got %q", stderr.String())
	}

	// Format already formatted code
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--code", "dml 1.4;\ndevice dev;\n"}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr for --code, got %q", stderr.String())
	}

	// Empty string code
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-c", ""}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK for empty code, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// Syntax error on code: echoes input on stdout, writes error on stderr, exits 123
	stdout.Reset()
	stderr.Reset()
	badInput := "dml 1.4;\ninvalid syntax :::\n"
	code = run([]string{"-c", badInput}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError on syntax error, got %d", code)
	}
	if stdout.String() != badInput {
		t.Errorf(
			"expected original input echoed on stdout, got %q",
			stdout.String(),
		)
	}
	if !strings.Contains(stderr.String(), "<string>") ||
		!strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf(
			"expected syntax error message referencing <string>, got %q",
			stderr.String(),
		)
	}

	// Line length flag respected
	stdout.Reset()
	stderr.Reset()
	longCode := "dml 1.4;\nconstant msg = \"a very long string that should wrap when formatted with short line length limit\";\n"
	code = run([]string{"-c", longCode, "-l", "40"}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() == 0 {
		t.Errorf("expected non-empty formatted output")
	}
}

func TestRun_Stdin(t *testing.T) {
	// Reformat from stdin
	var stdout, stderr bytes.Buffer
	unformatted := "dml   1.4  ;\n  device   dev  ;\n"
	code := runWithIO(
		[]string{"-"},
		strings.NewReader(unformatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expected := "dml 1.4;\ndevice dev;\n"
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.String() != "reformatted -\n" {
		t.Errorf("expected stderr 'reformatted -\\n', got %q", stderr.String())
	}

	// Unchanged stdin is silent on stderr
	stdout.Reset()
	stderr.Reset()
	code = runWithIO(
		[]string{"-"},
		strings.NewReader(expected),
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr for unchanged stdin, got %q",
			stderr.String(),
		)
	}

	// Empty stdin
	stdout.Reset()
	stderr.Reset()
	code = runWithIO([]string{"-"}, strings.NewReader(""), &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK for empty stdin, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf(
			"expected empty stdout for empty stdin, got %q",
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr for empty stdin, got %q",
			stderr.String(),
		)
	}

	// Stdin with --stdin-filename (reformatted)
	stdout.Reset()
	stderr.Reset()
	code = runWithIO(
		[]string{"--stdin-filename", "custom.dml", "-"},
		strings.NewReader(unformatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.String() != "reformatted custom.dml\n" {
		t.Errorf(
			"expected stderr 'reformatted custom.dml\\n', got %q",
			stderr.String(),
		)
	}

	// Stdin with --stdin-filename (syntax error)
	stdout.Reset()
	stderr.Reset()
	badInput := "dml 1.4;\ninvalid syntax :::\n"
	code = runWithIO(
		[]string{"--stdin-filename", "custom.dml", "-"},
		strings.NewReader(badInput),
		&stdout,
		&stderr,
	)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError on syntax error, got %d", code)
	}
	if stdout.String() != badInput {
		t.Errorf(
			"expected original input echoed on stdout, got %q",
			stdout.String(),
		)
	}
	if !strings.Contains(stderr.String(), "custom.dml") ||
		!strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf(
			"expected syntax error referencing custom.dml, got %q",
			stderr.String(),
		)
	}

	// Stdin without custom name (syntax error)
	stdout.Reset()
	stderr.Reset()
	code = runWithIO(
		[]string{"-"},
		strings.NewReader(badInput),
		&stdout,
		&stderr,
	)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError on syntax error, got %d", code)
	}
	if stdout.String() != badInput {
		t.Errorf(
			"expected original input echoed on stdout, got %q",
			stdout.String(),
		)
	}
	if !strings.Contains(stderr.String(), "-:") ||
		!strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf(
			"expected syntax error referencing '-', got %q",
			stderr.String(),
		)
	}

	// --stdin-filename ignored for file on disk
	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "disk.dml")
	if err := os.WriteFile(dmlPath, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to write disk.dml: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{"--stdin-filename", "ignored.dml", dmlPath},
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != "reformatted "+dmlPath+"\n1 file reformatted.\n" {
		t.Errorf(
			"expected stdout 'reformatted %s', got %q",
			dmlPath,
			stdout.String(),
		)
	}

	// --stdin-filename ignored for -c
	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{"--stdin-filename", "ignored.dml", "-c", unformatted},
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr for -c with --stdin-filename, got %q",
			stderr.String(),
		)
	}
}

func TestRun_Check_Files(t *testing.T) {
	tmpDir := t.TempDir()
	unformatted := "dml   1.4  ;\n  device   dev  ;\n"
	formatted := "dml 1.4;\ndevice dev;\n"

	unformattedFile := filepath.Join(tmpDir, "unformatted.dml")
	if err := os.WriteFile(
		unformattedFile,
		[]byte(unformatted),
		0600,
	); err != nil {
		t.Fatalf("failed to write unformatted.dml: %v", err)
	}

	formattedFile := filepath.Join(tmpDir, "formatted.dml")
	if err := os.WriteFile(formattedFile, []byte(formatted), 0600); err != nil {
		t.Fatalf("failed to write formatted.dml: %v", err)
	}

	badFile := filepath.Join(tmpDir, "bad.dml")
	if err := os.WriteFile(
		badFile,
		[]byte("dml 1.4;\ninvalid syntax :::\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write bad.dml: %v", err)
	}

	// 1. Single unformatted file: returns 1, prints "would reformat <file>\n" to stdout, stderr empty, file untouched
	var stdout, stderr bytes.Buffer
	code := run([]string{"--check", unformattedFile}, &stdout, &stderr)
	if code != ExitCodeWouldReformat {
		t.Fatalf(
			"expected ExitCodeWouldReformat (1), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedStdout := "would reformat " + unformattedFile + "\n1 file would be reformatted.\n"
	if stdout.String() != expectedStdout {
		t.Errorf("expected stdout %q, got %q", expectedStdout, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}
	content, err := os.ReadFile(filepath.Clean(unformattedFile))
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != unformatted {
		t.Errorf("expected file on disk to remain unchanged")
	}

	// 2. Single formatted file: returns 0, outputs 1 file would be left unchanged, stderr empty, file untouched
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--check", formattedFile}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedFormattedStdout := "1 file would be left unchanged.\n"
	if stdout.String() != expectedFormattedStdout {
		t.Errorf(
			"expected stdout %q, got %q",
			expectedFormattedStdout,
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// 3. Syntax error file: returns 123, outputs 1 file would fail to reformat on stdout, error on stderr
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--check", badFile}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	expectedBadStdout := "1 file would fail to reformat.\n"
	if stdout.String() != expectedBadStdout {
		t.Errorf(
			"expected stdout %q, got %q",
			expectedBadStdout,
			stdout.String(),
		)
	}
	if !strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf("expected syntax error on stderr, got %q", stderr.String())
	}

	// 4. Multiple files (one unformatted, one formatted): returns 1
	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{"--check", unformattedFile, formattedFile},
		&stdout,
		&stderr,
	)
	if code != ExitCodeWouldReformat {
		t.Fatalf(
			"expected ExitCodeWouldReformat (1), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedMultiStdout := "would reformat " + unformattedFile + "\n1 file would be reformatted, 1 file would be left unchanged.\n"
	if stdout.String() != expectedMultiStdout {
		t.Errorf(
			"expected stdout %q, got %q",
			expectedMultiStdout,
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// 5. Multiple files (one unformatted, one syntax error): returns 123 (error takes precedence)
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--check", unformattedFile, badFile}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf(
			"expected ExitCodeInternalError (123), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if !strings.Contains(stdout.String(), "would reformat "+unformattedFile) {
		t.Errorf("expected would reformat on stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf("expected syntax error on stderr, got %q", stderr.String())
	}
}

func TestRun_Check_Code(t *testing.T) {
	// Unformatted code: returns 1, stdout & stderr empty
	var stdout, stderr bytes.Buffer
	code := run(
		[]string{"--check", "-c", "dml   1.4  ;\n  device   dev  ;\n"},
		&stdout,
		&stderr,
	)
	if code != ExitCodeWouldReformat {
		t.Fatalf(
			"expected ExitCodeWouldReformat (1), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr for -c, got %q", stderr.String())
	}

	// Formatted code: returns 0, stdout & stderr empty
	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{"--check", "-c", "dml 1.4;\ndevice dev;\n"},
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// Empty code: returns 0, stdout & stderr empty
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--check", "-c", ""}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// Syntax error on code: returns 123, stdout empty (does not echo), error on stderr
	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{"--check", "-c", "dml 1.4;\ninvalid :::\n"},
		&stdout,
		&stderr,
	)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout with --check, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "<string>") ||
		!strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf(
			"expected syntax error message on stderr, got %q",
			stderr.String(),
		)
	}
}

func TestRun_Check_Stdin(t *testing.T) {
	// Unformatted stdin: returns 1, stdout has "would reformat -\n", stderr empty
	var stdout, stderr bytes.Buffer
	unformatted := "dml   1.4  ;\n  device   dev  ;\n"
	code := runWithIO(
		[]string{"--check", "-"},
		strings.NewReader(unformatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeWouldReformat {
		t.Fatalf(
			"expected ExitCodeWouldReformat (1), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != "would reformat -\n" {
		t.Errorf(
			"expected stdout 'would reformat -\\n', got %q",
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// Formatted stdin: returns 0, stdout & stderr empty
	stdout.Reset()
	stderr.Reset()
	formatted := "dml 1.4;\ndevice dev;\n"
	code = runWithIO(
		[]string{"--check", "-"},
		strings.NewReader(formatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// Stdin with --stdin-filename: returns 1, stdout has "would reformat custom.dml\n", stderr empty
	stdout.Reset()
	stderr.Reset()
	code = runWithIO(
		[]string{"--check", "--stdin-filename", "custom.dml", "-"},
		strings.NewReader(unformatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeWouldReformat {
		t.Fatalf(
			"expected ExitCodeWouldReformat (1), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != "would reformat custom.dml\n" {
		t.Errorf(
			"expected stdout 'would reformat custom.dml\\n', got %q",
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// Syntax error on stdin: returns 123, stdout empty, error on stderr
	stdout.Reset()
	stderr.Reset()
	code = runWithIO(
		[]string{"--check", "-"},
		strings.NewReader("dml 1.4;\ninvalid :::\n"),
		&stdout,
		&stderr,
	)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout with --check, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "-:") ||
		!strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf("expected syntax error on stderr, got %q", stderr.String())
	}
}

func TestColorAndNoColorFlagSupport(t *testing.T) {
	// 1. Verify flag parsing and precedence
	var stderr bytes.Buffer
	cfg, err := parseArgs([]string{"--color", "dev.dml"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error parsing --color: %v", err)
	}
	if !cfg.Color || cfg.NoColor {
		t.Errorf(
			"expected Color=true, NoColor=false; got Color=%v, NoColor=%v",
			cfg.Color,
			cfg.NoColor,
		)
	}

	stderr.Reset()
	cfg, err = parseArgs([]string{"--no-color", "dev.dml"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error parsing --no-color: %v", err)
	}
	if cfg.Color || !cfg.NoColor {
		t.Errorf(
			"expected Color=false, NoColor=true; got Color=%v, NoColor=%v",
			cfg.Color,
			cfg.NoColor,
		)
	}

	// Last flag wins: --color then --no-color => Color=false
	stderr.Reset()
	cfg, err = parseArgs([]string{"--color", "--no-color", "dev.dml"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error parsing --color --no-color: %v", err)
	}
	if cfg.Color {
		t.Errorf("expected Color=false when --no-color follows --color")
	}

	// Last flag wins: --no-color then --color => Color=true
	stderr.Reset()
	cfg, err = parseArgs([]string{"--no-color", "--color", "dev.dml"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error parsing --no-color --color: %v", err)
	}
	if !cfg.Color {
		t.Errorf("expected Color=true when --color follows --no-color")
	}

	// No flags provided => ColorMode == ColorAuto
	stderr.Reset()
	cfg, err = parseArgs([]string{"dev.dml"}, &stderr)
	if err != nil {
		t.Fatalf("unexpected error parsing dev.dml: %v", err)
	}
	if cfg.ColorMode != ColorAuto {
		t.Errorf(
			"expected ColorMode=ColorAuto by default, got %v",
			cfg.ColorMode,
		)
	}

	// 2. Test run with --color on an unformatted file (formats and outputs colored message on stdout)
	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "color_test.dml")
	if err := os.WriteFile(
		dmlPath,
		[]byte("dml   1.4  ;\n  device   dev  ;\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	var stdout bytes.Buffer
	stderr.Reset()
	code := run([]string{"--color", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	outStr := stdout.String()
	if !strings.Contains(outStr, "\033[1;32mreformatted\033[0m") ||
		!strings.Contains(outStr, "\033[1m"+dmlPath+"\033[0m") {
		t.Errorf(
			"expected colored reformatted output on stdout, got %q",
			outStr,
		)
	}

	// 3. Test run with --color --check on an unformatted file (outputs colored would reformat on stdout)
	if err := os.WriteFile(
		dmlPath,
		[]byte("dml   1.4  ;\n  device   dev  ;\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--color", "--check", dmlPath}, &stdout, &stderr)
	if code != ExitCodeWouldReformat {
		t.Fatalf("expected ExitCodeWouldReformat (1), got %d", code)
	}
	outStr = stdout.String()
	if !strings.Contains(outStr, "\033[1;34mwould reformat\033[0m") ||
		!strings.Contains(outStr, "\033[1m"+dmlPath+"\033[0m") {
		t.Errorf(
			"expected colored would reformat output on stdout, got %q",
			outStr,
		)
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// 4. Test run with --color on syntax error file (outputs colored syntax error on stderr)
	badPath := filepath.Join(tmpDir, "bad.dml")
	if err := os.WriteFile(
		badPath,
		[]byte("dml 1.4; @@@ invalid"),
		0600,
	); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--color", badPath}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	errStr := stderr.String()
	if !strings.Contains(errStr, "\033[1;31mSyntax error\033[0m") {
		t.Errorf("expected colored syntax error on stderr, got %q", errStr)
	}

	// 5. Test run with --no-color on unformatted file (outputs plain text)
	if err := os.WriteFile(
		dmlPath,
		[]byte("dml   1.4  ;\n  device   dev  ;\n"),
		0600,
	); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--no-color", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf("expected ExitCodeOK (0), got %d", code)
	}
	outStr = stdout.String()
	if strings.Contains(outStr, "\033") {
		t.Errorf("expected uncolored output with --no-color, got %q", outStr)
	}
	expectedPlain := "reformatted " + dmlPath + "\n1 file reformatted.\n"
	if outStr != expectedPlain {
		t.Errorf(
			"expected plain reformatted string %q, got %q",
			expectedPlain,
			outStr,
		)
	}
}

func TestDebugFlag_HelpOutput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--help"}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf("expected ExitCodeOK, got %d", code)
	}
	helpOut := stdout.String()
	if isDebugEnabled {
		for _, expected := range []string{
			"Debug Options:",
			"--cpuprofile FILE               Write CPU profile to FILE.",
			"--memprofile",
		} {
			if !strings.Contains(helpOut, expected) {
				t.Errorf(
					"expected help output to contain %q when debug is enabled",
					expected,
				)
			}
		}
	} else {
		if strings.Contains(helpOut, "Debug Options:") {
			t.Errorf(
				"expected help output not to contain 'Debug Options:' when ndebug is enabled",
			)
		}
	}
}

func TestDebugFlag_ProfileFlagsUnlockedInDebug(t *testing.T) {
	if isDebugEnabled {
		if !isFlagImplemented("cpuprofile") {
			t.Errorf("expected cpuprofile to be unlocked when debug is enabled")
		}
		if !isFlagImplemented("memprofile") {
			t.Errorf("expected memprofile to be unlocked when debug is enabled")
		}
	} else {
		if isFlagImplemented("cpuprofile") {
			t.Errorf("expected cpuprofile to be locked when ndebug is enabled")
		}
		if isFlagImplemented("memprofile") {
			t.Errorf("expected memprofile to be locked when ndebug is enabled")
		}
	}
}

func TestRun_ProfileFlagsWhenDebugDisabled(t *testing.T) {
	if isDebugEnabled {
		t.Skip("skipping in debug build")
	}
	var stdout, stderr bytes.Buffer
	code := run(
		[]string{"--cpuprofile", "foo.prof", "-c", "dml 1.4;\n"},
		&stdout,
		&stderr,
	)
	if code == ExitCodeOK {
		t.Errorf(
			"expected error when debug flag is passed in ndebug build, got code %d",
			code,
		)
	}
}

func TestRun_CPUProfile(t *testing.T) {
	if !isDebugEnabled {
		t.Skip("skipping in ndebug build")
	}

	tmpDir := t.TempDir()
	cpuProf := filepath.Join(tmpDir, "cpu.prof")

	var stdout, stderr bytes.Buffer
	code := run(
		[]string{"--cpuprofile", cpuProf, "-c", "dml 1.4;\ndevice test_dev;\n"},
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	info, err := os.Stat(cpuProf)
	if err != nil {
		t.Fatalf("expected cpu profile to exist: %v", err)
	}
	if info.Size() == 0 {
		t.Errorf("expected non-empty cpu profile, got size 0")
	}
}

func TestRun_MemProfile(t *testing.T) {
	if !isDebugEnabled {
		t.Skip("skipping in ndebug build")
	}

	tmpDir := t.TempDir()
	memProf := filepath.Join(tmpDir, "mem.prof")

	var stdout, stderr bytes.Buffer
	code := run(
		[]string{"--memprofile", memProf, "-c", "dml 1.4;\ndevice test_dev;\n"},
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	info, err := os.Stat(memProf)
	if err != nil {
		t.Fatalf("expected mem profile to exist: %v", err)
	}
	if info.Size() == 0 {
		t.Errorf("expected non-empty mem profile, got size 0")
	}
}

func TestRun_ProfileEmptyPathErrors(t *testing.T) {
	if !isDebugEnabled {
		t.Skip("skipping in ndebug build")
	}

	flags := []string{"--cpuprofile", "--memprofile"}
	for _, flag := range flags {
		var stdout, stderr bytes.Buffer
		code := run([]string{flag, "", "-c", "dml 1.4;\n"}, &stdout, &stderr)
		if code != ExitCodeInternalError {
			t.Errorf(
				"expected ExitCodeInternalError for %s='', got %d",
				flag,
				code,
			)
		}
		if !strings.Contains(
			stderr.String(),
			"requires a non-empty file path",
		) {
			t.Errorf(
				"expected non-empty error message, got: %s",
				stderr.String(),
			)
		}
	}
}

func TestRun_ProfileCreationError(t *testing.T) {
	if !isDebugEnabled {
		t.Skip("skipping in ndebug build")
	}

	invalidPath := filepath.Join(t.TempDir(), "nonexistent", "dir", "cpu.prof")
	var stdout, stderr bytes.Buffer
	code := run(
		[]string{"--cpuprofile", invalidPath, "-c", "dml 1.4;\n"},
		&stdout,
		&stderr,
	)
	if code != ExitCodeInternalError {
		t.Errorf(
			"expected ExitCodeInternalError for invalid path, got %d",
			code,
		)
	}
	if !strings.Contains(stderr.String(), "Error creating CPU profile") {
		t.Errorf(
			"expected 'Error creating CPU profile' in stderr, got: %s",
			stderr.String(),
		)
	}
}

func TestRun_DirectoryArgument(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "dev1.dml")
	file2 := filepath.Join(tmpDir, "sub", "dev2.dml")
	venvFile := filepath.Join(tmpDir, "venv", "dev3.dml")
	txtFile := filepath.Join(tmpDir, "readme.txt")

	_ = os.MkdirAll(filepath.Dir(file2), 0750)
	_ = os.MkdirAll(filepath.Dir(venvFile), 0750)

	unformatted := []byte("dml   1.4  ;\n  device   dev  ;\n")
	_ = os.WriteFile(file1, unformatted, 0600)
	_ = os.WriteFile(file2, unformatted, 0600)
	_ = os.WriteFile(venvFile, unformatted, 0600)
	_ = os.WriteFile(txtFile, []byte("text file"), 0600)

	var stdout, stderr bytes.Buffer
	code := run([]string{tmpDir}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	outStr := stdout.String()
	if !strings.Contains(outStr, "reformatted "+file1) {
		t.Errorf("expected reformatted %s, got: %s", file1, outStr)
	}
	if !strings.Contains(outStr, "reformatted "+file2) {
		t.Errorf("expected reformatted %s, got: %s", file2, outStr)
	}
	if strings.Contains(outStr, "dev3.dml") {
		t.Errorf("expected venv file not to be formatted, got: %s", outStr)
	}
	if strings.Contains(outStr, "readme.txt") {
		t.Errorf("expected non-dml file not to be formatted, got: %s", outStr)
	}

	// Verify formatted contents
	c1, _ := os.ReadFile(filepath.Clean(file1))
	if string(c1) != "dml 1.4;\ndevice dev;\n" {
		t.Errorf("file1 not formatted properly: %q", string(c1))
	}
	cVenv, _ := os.ReadFile(filepath.Clean(venvFile))
	if string(cVenv) != string(unformatted) {
		t.Errorf("venv file was modified: %q", string(cVenv))
	}
}

func TestRun_EmptySRCFailsWithUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf(
			"expected ExitCodeInternalError for empty SRC, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	expectedErr := "Usage: dmlfmt [OPTIONS] SRC ...\n\nOne of 'SRC' or 'code' is required.\n"
	if stderr.String() != expectedErr {
		t.Errorf(
			"expected stderr %q, got: %q",
			expectedErr,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout, got %q", stdout.String())
	}
}

func TestRun_ExcludeAndExtendExclude(t *testing.T) {
	tmpDir := t.TempDir()
	customDir := filepath.Join(tmpDir, "custom")
	extraDir := filepath.Join(tmpDir, "extra")
	keepDir := filepath.Join(tmpDir, "keep")

	_ = os.MkdirAll(customDir, 0750)
	_ = os.MkdirAll(extraDir, 0750)
	_ = os.MkdirAll(keepDir, 0750)

	unformatted := []byte("dml   1.4  ;\n  device   dev  ;\n")
	_ = os.WriteFile(filepath.Join(customDir, "f.dml"), unformatted, 0600)
	_ = os.WriteFile(filepath.Join(extraDir, "f.dml"), unformatted, 0600)
	_ = os.WriteFile(filepath.Join(keepDir, "f.dml"), unformatted, 0600)

	var stdout, stderr bytes.Buffer
	code := run(
		[]string{
			"--exclude",
			"/custom/",
			"--extend-exclude",
			"/extra/",
			tmpDir,
		},
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	outStr := stdout.String()
	if strings.Contains(outStr, "custom") || strings.Contains(outStr, "extra") {
		t.Errorf("excluded files were reported as formatted: %s", outStr)
	}
	if !strings.Contains(outStr, filepath.Join(keepDir, "f.dml")) {
		t.Errorf("expected keep/f.dml to be formatted, got: %s", outStr)
	}
}

func TestRun_ForceExcludeExplicitFile(t *testing.T) {
	tmpDir := t.TempDir()
	f1 := filepath.Join(tmpDir, "f1.dml")
	f2 := filepath.Join(tmpDir, "force_skip.dml")

	unformatted := []byte("dml   1.4  ;\n  device   dev  ;\n")
	_ = os.WriteFile(f1, unformatted, 0600)
	_ = os.WriteFile(f2, unformatted, 0600)

	var stdout, stderr bytes.Buffer
	code := run(
		[]string{"--force-exclude", "force_skip", f1, f2},
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	outStr := stdout.String()
	if !strings.Contains(outStr, f1) {
		t.Errorf("expected f1 to be formatted, got: %s", outStr)
	}
	if strings.Contains(outStr, f2) {
		t.Errorf("expected f2 to be force-excluded, got: %s", outStr)
	}

	c2, _ := os.ReadFile(filepath.Clean(f2))
	if string(c2) != string(unformatted) {
		t.Errorf("force-excluded file was modified")
	}
}

func TestRun_StdinForceExclude(t *testing.T) {
	input := "dml   1.4  ;\n  device   dev  ;\n"

	// Case 1: Stdin without check; content passed through untouched
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader(input)
	code := runWithIO(
		[]string{
			"--stdin-filename",
			"path/to/gen_skip.dml",
			"--force-exclude",
			"gen_skip",
			"-",
		},
		stdin,
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf("expected ExitCodeOK for force-excluded stdin, got %d", code)
	}
	if stdout.String() != input {
		t.Errorf("expected stdout %q, got %q", input, stdout.String())
	}

	// Case 2: Stdin with --check; exits 0 without writing
	stdout.Reset()
	stderr.Reset()
	stdin = strings.NewReader(input)
	code = runWithIO(
		[]string{
			"--stdin-filename", "path/to/gen_skip.dml",
			"--force-exclude", "gen_skip",
			"--check",
			"-",
		},
		stdin,
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK for force-excluded stdin with check, got %d",
			code,
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout with check, got: %q", stdout.String())
	}
}

func TestRun_InvalidRegexExit123(t *testing.T) {
	flags := []string{
		"--include",
		"--exclude",
		"--extend-exclude",
		"--force-exclude",
	}
	for _, f := range flags {
		var stdout, stderr bytes.Buffer
		code := run([]string{f, "[", "dummy.dml"}, &stdout, &stderr)
		if code != ExitCodeInternalError {
			t.Errorf(
				"expected ExitCodeInternalError for %s '[', got %d",
				f,
				code,
			)
		}
		if !strings.Contains(stderr.String(), "invalid regular expression") {
			t.Errorf(
				"expected 'invalid regular expression' in stderr for %s, got: %s",
				f,
				stderr.String(),
			)
		}
	}
}

func TestRun_PipelineWorkers(t *testing.T) {
	workerCounts := []string{"1", "4", "8"}
	for _, w := range workerCounts {
		t.Run("workers_"+w, func(t *testing.T) {
			tmpDir := t.TempDir()
			unformatted := []byte("dml   1.4  ;\n  device   dev  ;\n")
			expected := "dml 1.4;\ndevice dev;\n"

			files := make([]string, 0, 5)
			for i := range 5 {
				path := filepath.Join(tmpDir, fmt.Sprintf("f%d.dml", i))
				_ = os.WriteFile(path, unformatted, 0600)
				files = append(files, path)
			}

			var stdout, stderr bytes.Buffer
			args := append([]string{"-W", w}, files...)
			code := run(args, &stdout, &stderr)
			if code != ExitCodeOK {
				t.Fatalf(
					"expected ExitCodeOK with -W %s, got %d (stderr: %s)",
					w,
					code,
					stderr.String(),
				)
			}

			for _, f := range files {
				content, err := os.ReadFile(filepath.Clean(f))
				if err != nil {
					t.Fatalf("failed to read file %s: %v", f, err)
				}
				if string(content) != expected {
					t.Errorf(
						"file %s not properly formatted: %q",
						f,
						string(content),
					)
				}
			}
		})
	}
}

func TestRun_LargeBatchBackpressure(t *testing.T) {
	tmpDir := t.TempDir()
	unformatted := []byte("dml   1.4  ;\n  device   dev  ;\n")
	expected := "dml 1.4;\ndevice dev;\n"

	// Create 100 files to verify queue backpressure and interleaving across workers
	files := make([]string, 0, 100)
	for i := range 100 {
		path := filepath.Join(tmpDir, fmt.Sprintf("file_%03d.dml", i))
		_ = os.WriteFile(path, unformatted, 0600)
		files = append(files, path)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-W", "4", tmpDir}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK for 100 files, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	for _, f := range files {
		content, err := os.ReadFile(filepath.Clean(f))
		if err != nil {
			t.Fatalf("failed to read file %s: %v", f, err)
		}
		if string(content) != expected {
			t.Errorf("file %s not properly formatted", f)
		}
	}
}

func TestRun_PipelineErrorAccumulation(t *testing.T) {
	tmpDir := t.TempDir()
	goodFile1 := filepath.Join(tmpDir, "good1.dml")
	badFile := filepath.Join(tmpDir, "bad.dml")
	goodFile2 := filepath.Join(tmpDir, "good2.dml")

	unformatted := []byte("dml   1.4  ;\n  device   dev  ;\n")
	badContent := []byte("dml 1.4; @@@ invalid syntax")

	_ = os.WriteFile(goodFile1, unformatted, 0600)
	_ = os.WriteFile(badFile, badContent, 0600)
	_ = os.WriteFile(goodFile2, unformatted, 0600)

	var unreadableFile string
	if runtime.GOOS != "windows" && os.Geteuid() != 0 {
		unreadableFile = filepath.Join(tmpDir, "unreadable.dml")
		_ = os.WriteFile(unreadableFile, unformatted, 0000)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-W", "2", tmpDir}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf(
			"expected ExitCodeInternalError when file has error, got %d",
			code,
		)
	}

	// Verify both good files were still reformatted
	c1, _ := os.ReadFile(filepath.Clean(goodFile1))
	if string(c1) != "dml 1.4;\ndevice dev;\n" {
		t.Errorf(
			"good1.dml was not reformatted despite pipeline error in bad.dml",
		)
	}
	c2, _ := os.ReadFile(filepath.Clean(goodFile2))
	if string(c2) != "dml 1.4;\ndevice dev;\n" {
		t.Errorf(
			"good2.dml was not reformatted despite pipeline error in bad.dml",
		)
	}

	// Verify bad file was not corrupted
	cBad, _ := os.ReadFile(filepath.Clean(badFile))
	if !bytes.Equal(cBad, badContent) {
		t.Errorf("bad.dml was modified")
	}

	if unreadableFile != "" {
		if !strings.Contains(stderr.String(), "opening file") &&
			!strings.Contains(stderr.String(), "reading file") {
			t.Errorf(
				"expected unreadable file error on stderr, got: %s",
				stderr.String(),
			)
		}
	}
}

func TestWorkers_CLIValidation(t *testing.T) {
	invalidFlags := []struct {
		name string
		args []string
		val  string
	}{
		{"-W 0", []string{"-W", "0"}, "0"},
		{"-W -1", []string{"-W", "-1"}, "-1"},
		{"--workers 0", []string{"--workers", "0"}, "0"},
		{"--workers -5", []string{"--workers", "-5"}, "-5"},
	}

	for _, tc := range invalidFlags {
		t.Run("file_"+tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "dev.dml")
			unformatted := []byte("dml   1.4  ;\n  device   dev  ;\n")
			_ = os.WriteFile(filePath, unformatted, 0600)

			var stdout, stderr bytes.Buffer
			args := append(tc.args, filePath)
			code := run(args, &stdout, &stderr)
			if code != ExitCodeInternalError {
				t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
			}
			expectedMsg := fmt.Sprintf(
				"Error: Invalid value for -W / --workers: %s is not in the range x>=1.",
				tc.val,
			)
			if !strings.Contains(stderr.String(), expectedMsg) {
				t.Errorf(
					"expected stderr to contain %q, got: %s",
					expectedMsg,
					stderr.String(),
				)
			}

			// Verify file was untouched
			after, _ := os.ReadFile(filepath.Clean(filePath))
			if !bytes.Equal(after, unformatted) {
				t.Errorf("file should not have been modified")
			}
		})

		t.Run("code_"+tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			args := append(tc.args, "-c", "dml 1.4;\n")
			code := run(args, &stdout, &stderr)
			if code != ExitCodeInternalError {
				t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
			}
			expectedMsg := fmt.Sprintf(
				"Error: Invalid value for -W / --workers: %s is not in the range x>=1.",
				tc.val,
			)
			if !strings.Contains(stderr.String(), expectedMsg) {
				t.Errorf(
					"expected stderr to contain %q, got: %s",
					expectedMsg,
					stderr.String(),
				)
			}
		})

		t.Run("stdin_"+tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			args := append(tc.args, "-")
			code := runWithIO(
				args,
				strings.NewReader("dml 1.4;\n"),
				&stdout,
				&stderr,
			)
			if code != ExitCodeInternalError {
				t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
			}
			expectedMsg := fmt.Sprintf(
				"Error: Invalid value for -W / --workers: %s is not in the range x>=1.",
				tc.val,
			)
			if !strings.Contains(stderr.String(), expectedMsg) {
				t.Errorf(
					"expected stderr to contain %q, got: %s",
					expectedMsg,
					stderr.String(),
				)
			}
		})
	}
}

func TestWorkers_EnvVarConfiguration(t *testing.T) {
	t.Run("valid integer", func(t *testing.T) {
		t.Setenv("DMLFMT_NUM_WORKERS", "4")
		var stderr bytes.Buffer
		cfg, err := parseArgs([]string{"dev.dml"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Workers != 4 {
			t.Errorf("expected Workers=4, got %d", cfg.Workers)
		}
		if cfg.WorkersErr != "" {
			t.Errorf("unexpected WorkersErr: %s", cfg.WorkersErr)
		}
	})

	t.Run("whitespace trimmed", func(t *testing.T) {
		t.Setenv("DMLFMT_NUM_WORKERS", "  8  ")
		var stderr bytes.Buffer
		cfg, err := parseArgs([]string{"dev.dml"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Workers != 8 {
			t.Errorf("expected Workers=8, got %d", cfg.Workers)
		}
	})

	t.Run("empty treated as unset", func(t *testing.T) {
		t.Setenv("DMLFMT_NUM_WORKERS", "")
		var stderr bytes.Buffer
		cfg, err := parseArgs([]string{"dev.dml"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expectedDefaultWorkers := min(2, runtime.NumCPU())
		if cfg.Workers != expectedDefaultWorkers {
			t.Errorf(
				"expected default Workers=%d, got %d",
				expectedDefaultWorkers,
				cfg.Workers,
			)
		}
		if cfg.WorkersErr != "" {
			t.Errorf("unexpected WorkersErr: %s", cfg.WorkersErr)
		}
	})

	t.Run("formats files with env workers", func(t *testing.T) {
		t.Setenv("DMLFMT_NUM_WORKERS", "2")
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "f.dml")
		unformatted := []byte("dml   1.4  ;\n  device   dev  ;\n")
		_ = os.WriteFile(filePath, unformatted, 0600)

		var stdout, stderr bytes.Buffer
		code := run([]string{filePath}, &stdout, &stderr)
		if code != ExitCodeOK {
			t.Fatalf(
				"expected ExitCodeOK, got %d (stderr: %s)",
				code,
				stderr.String(),
			)
		}
		content, _ := os.ReadFile(filepath.Clean(filePath))
		if string(content) != "dml 1.4;\ndevice dev;\n" {
			t.Errorf("file not reformatted correctly: %q", string(content))
		}
	})
}

func TestWorkers_EnvVarValidation(t *testing.T) {
	cases := []struct {
		name        string
		envVal      string
		expectedMsg string
	}{
		{
			"zero",
			"0",
			"Error: Invalid value for DMLFMT_NUM_WORKERS: 0 is not in the range x>=1.",
		},
		{
			"negative",
			"-2",
			"Error: Invalid value for DMLFMT_NUM_WORKERS: -2 is not in the range x>=1.",
		},
		{
			"non_integer",
			"abc",
			"Error: Invalid value for DMLFMT_NUM_WORKERS: \"abc\" is not a valid integer.",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("DMLFMT_NUM_WORKERS", tc.envVal)
			var stdout, stderr bytes.Buffer
			code := run([]string{"dev.dml"}, &stdout, &stderr)
			if code != ExitCodeInternalError {
				t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
			}
			if !strings.Contains(stderr.String(), tc.expectedMsg) {
				t.Errorf(
					"expected stderr to contain %q, got: %s",
					tc.expectedMsg,
					stderr.String(),
				)
			}
		})
	}
}

func TestWorkers_FlagOverridesEnvVar(t *testing.T) {
	t.Run("explicit flag overrides valid env var", func(t *testing.T) {
		t.Setenv("DMLFMT_NUM_WORKERS", "8")
		var stderr bytes.Buffer
		cfg, err := parseArgs([]string{"-W", "2", "dev.dml"}, &stderr)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Workers != 2 {
			t.Errorf("expected Workers=2 from flag, got %d", cfg.Workers)
		}
	})

	t.Run("explicit flag overrides invalid env var", func(t *testing.T) {
		t.Setenv("DMLFMT_NUM_WORKERS", "invalid")
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "f.dml")
		_ = os.WriteFile(filePath, []byte("dml 1.4;\ndevice dev;\n"), 0600)

		var stdout, stderr bytes.Buffer
		code := run([]string{"-W", "2", filePath}, &stdout, &stderr)
		if code != ExitCodeOK {
			t.Fatalf(
				"expected ExitCodeOK with -W override, got %d (stderr: %s)",
				code,
				stderr.String(),
			)
		}
	})
}

func TestWorkers_HelpAndVersionBypass(t *testing.T) {
	cases := []struct {
		name   string
		envVal string
		args   []string
	}{
		{"-W 0 --help", "", []string{"-W", "0", "--help"}},
		{"--workers -1 -h", "", []string{"--workers", "-1", "-h"}},
		{"-W 0 --version", "", []string{"-W", "0", "--version"}},
		{"env 0 --help", "0", []string{"--help"}},
		{"env invalid --version", "invalid", []string{"--version"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envVal != "" {
				t.Setenv("DMLFMT_NUM_WORKERS", tc.envVal)
			}
			var stdout, stderr bytes.Buffer
			code := run(tc.args, &stdout, &stderr)
			if code != ExitCodeOK {
				t.Fatalf(
					"expected ExitCodeOK (0), got %d (stderr: %s)",
					code,
					stderr.String(),
				)
			}
		})
	}
}

func TestRun_Dml12_Ignored_FileMode(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "v12.dml")
	content := []byte("dml 1.2;\nparameter foo = 1;\ncall $bar();\n")
	if err := os.WriteFile(filePath, content, 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{filePath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedMsg := fmt.Sprintf(
		"ignored %s (dml 1.2)\n1 file ignored due to dml 1.2.\n",
		filePath,
	)
	if stdout.String() != expectedMsg {
		t.Errorf("expected stdout %q, got %q", expectedMsg, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	afterContent, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		t.Fatalf("failed to read file after run: %v", err)
	}
	if !bytes.Equal(afterContent, content) {
		t.Errorf("file was modified: got %q, want %q", afterContent, content)
	}
}

func TestRun_Dml12_Ignored_CheckMode(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "v12.dml")
	content := []byte("dml 1.2;\nparameter foo = 1;\ncall $bar();\n")
	if err := os.WriteFile(filePath, content, 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"--check", filePath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedMsg := fmt.Sprintf(
		"ignored %s (dml 1.2)\n1 file ignored due to dml 1.2.\n",
		filePath,
	)
	if stdout.String() != expectedMsg {
		t.Errorf("expected stdout %q, got %q", expectedMsg, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRun_Dml12_MixedDirectory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	v12Path := filepath.Join(tmpDir, "v12.dml")
	v12Content := []byte("dml 1.2;\nparameter foo = 1;\ncall $bar();\n")
	if err := os.WriteFile(v12Path, v12Content, 0600); err != nil {
		t.Fatalf("failed to write v12 file: %v", err)
	}

	v14CleanPath := filepath.Join(tmpDir, "v14_clean.dml")
	v14CleanContent := []byte("dml 1.4;\n\ndevice dev;\n")
	if err := os.WriteFile(v14CleanPath, v14CleanContent, 0600); err != nil {
		t.Fatalf("failed to write clean v14 file: %v", err)
	}

	v14DirtyPath := filepath.Join(tmpDir, "v14_dirty.dml")
	v14DirtyContent := []byte("dml 1.4; device dirty;\n")
	if err := os.WriteFile(v14DirtyPath, v14DirtyContent, 0600); err != nil {
		t.Fatalf("failed to write dirty v14 file: %v", err)
	}

	// 1. Run check mode: v12 is ignored, v14_clean is ok, v14_dirty needs reformat -> exits 1
	var stdoutCheck, stderrCheck bytes.Buffer
	codeCheck := run([]string{"--check", tmpDir}, &stdoutCheck, &stderrCheck)
	if codeCheck != ExitCodeWouldReformat {
		t.Fatalf(
			"expected ExitCodeWouldReformat (1), got %d (stderr: %s)",
			codeCheck,
			stderrCheck.String(),
		)
	}
	if !strings.Contains(stdoutCheck.String(), "ignored") ||
		!strings.Contains(stdoutCheck.String(), "v12.dml (dml 1.2)") {
		t.Errorf(
			"stdout missing ignored message for v12.dml: %s",
			stdoutCheck.String(),
		)
	}
	if !strings.Contains(stdoutCheck.String(), "v14_dirty.dml") {
		t.Errorf(
			"stdout missing would reformat for v14_dirty.dml: %s",
			stdoutCheck.String(),
		)
	}
	if stderrCheck.Len() != 0 {
		t.Errorf("expected empty stderrCheck, got %s", stderrCheck.String())
	}

	// 2. Run format mode: v12 untouched, v14_dirty reformatted -> exits 0
	var stdoutFmt, stderrFmt bytes.Buffer
	codeFmt := run([]string{tmpDir}, &stdoutFmt, &stderrFmt)
	if codeFmt != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			codeFmt,
			stderrFmt.String(),
		)
	}
	if !strings.Contains(stdoutFmt.String(), "v12.dml (dml 1.2)") ||
		!strings.Contains(stdoutFmt.String(), "ignored") {
		t.Errorf(
			"stdout missing ignored message for v12.dml: %s",
			stdoutFmt.String(),
		)
	}
	if !strings.Contains(stdoutFmt.String(), "v14_dirty.dml") {
		t.Errorf(
			"stdout missing reformatted message for v14_dirty.dml: %s",
			stdoutFmt.String(),
		)
	}
	if stderrFmt.Len() != 0 {
		t.Errorf("expected empty stderrFmt, got %s", stderrFmt.String())
	}

	// Verify v12 file remained untouched
	afterV12, err := os.ReadFile(filepath.Clean(v12Path))
	if err != nil {
		t.Fatalf("failed to read v12 file: %v", err)
	}
	if !bytes.Equal(afterV12, v12Content) {
		t.Errorf("v12 file was modified: got %q, want %q", afterV12, v12Content)
	}
}

func TestRun_Dml12_Stdin_Errors(t *testing.T) {
	t.Parallel()

	input := "dml 1.2;\nparameter foo = 1;\ncall $bar();\n"
	stdin := strings.NewReader(input)
	var stdout, stderr bytes.Buffer

	code := runWithIO([]string{"-"}, stdin, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	if !strings.Contains(stderr.String(), "Syntax error") ||
		!strings.Contains(
			stderr.String(),
			"Error Unimplemented: dml 1.2 syntax is not supported yet",
		) {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}
	if stdout.String() != input {
		t.Errorf("expected stdout to echo input, got %q", stdout.String())
	}
}

func TestRun_Dml12_Code_Errors(t *testing.T) {
	t.Parallel()

	input := "dml 1.2;\ndevice dev;\n"
	var stdout, stderr bytes.Buffer

	code := run([]string{"-c", input}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	if !strings.Contains(stderr.String(), "Syntax error") ||
		!strings.Contains(
			stderr.String(),
			"Error Unimplemented: dml 1.2 syntax is not supported yet",
		) {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}
	if stdout.String() != input {
		t.Errorf("expected stdout to echo input, got %q", stdout.String())
	}
}

func TestRun_Dml12_CheckStdin_Errors(t *testing.T) {
	t.Parallel()

	input := "dml 1.2;\ndevice dev;\n"
	stdin := strings.NewReader(input)
	var stdout, stderr bytes.Buffer

	code := runWithIO([]string{"--check", "-"}, stdin, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	if !strings.Contains(stderr.String(), "Syntax error") ||
		!strings.Contains(
			stderr.String(),
			"Error Unimplemented: dml 1.2 syntax is not supported yet",
		) {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout in check mode, got %q", stdout.String())
	}
}

func TestRun_OversizedFile_Ignored(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "huge.dml")
	f, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		t.Fatalf("failed to create huge file: %v", err)
	}
	if err := f.Truncate(51 << 20); err != nil {
		_ = f.Close()
		t.Fatalf("failed to truncate huge file: %v", err)
	}
	_ = f.Close()

	var stdout, stderr bytes.Buffer
	code := run([]string{filePath}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf(
			"expected ExitCodeInternalError (123), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedMsg := fmt.Sprintf(
		"ignored %s (file size exceeds 50 MiB limit)\n1 file failed to reformat.\n",
		filePath,
	)
	if stdout.String() != expectedMsg {
		t.Errorf("expected stdout %q, got %q", expectedMsg, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRun_OversizedFile_CheckMode(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "huge.dml")
	f, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		t.Fatalf("failed to create huge file: %v", err)
	}
	if err := f.Truncate(51 << 20); err != nil {
		_ = f.Close()
		t.Fatalf("failed to truncate huge file: %v", err)
	}
	_ = f.Close()

	var stdout, stderr bytes.Buffer
	code := run([]string{"--check", filePath}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf(
			"expected ExitCodeInternalError (123), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedMsg := fmt.Sprintf(
		"ignored %s (file size exceeds 50 MiB limit)\n1 file would fail to reformat.\n",
		filePath,
	)
	if stdout.String() != expectedMsg {
		t.Errorf("expected stdout %q, got %q", expectedMsg, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}
}

func TestRun_OversizedFile_MixedWithValid(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	hugePath := filepath.Join(tmpDir, "huge.dml")
	f, err := os.Create(filepath.Clean(hugePath))
	if err != nil {
		t.Fatalf("failed to create huge file: %v", err)
	}
	if err := f.Truncate(51 << 20); err != nil {
		_ = f.Close()
		t.Fatalf("failed to truncate huge file: %v", err)
	}
	_ = f.Close()

	validPath := filepath.Join(tmpDir, "valid.dml")
	validContent := []byte("dml  1.4 ;\n\ndevice  dev ;\n")
	if err := os.WriteFile(validPath, validContent, 0600); err != nil {
		t.Fatalf("failed to write valid file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{tmpDir}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf(
			"expected ExitCodeInternalError (123), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	outStr := stdout.String()
	expectedIgnored := fmt.Sprintf(
		"ignored %s (file size exceeds 50 MiB limit)\n",
		hugePath,
	)
	if !strings.Contains(outStr, expectedIgnored) {
		t.Errorf(
			"expected stdout to contain %q, got %q",
			expectedIgnored,
			outStr,
		)
	}
	expectedReformatted := fmt.Sprintf(
		"reformatted %s\n",
		validPath,
	)
	if !strings.Contains(outStr, expectedReformatted) {
		t.Errorf(
			"expected stdout to contain %q, got %q",
			expectedReformatted,
			outStr,
		)
	}
	if !strings.Contains(
		outStr,
		"1 file reformatted, 1 file failed to reformat.\n",
	) {
		t.Errorf("expected summary line in stdout, got %q", outStr)
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	afterValid, err := os.ReadFile(filepath.Clean(validPath))
	if err != nil {
		t.Fatalf("failed to read valid file: %v", err)
	}
	expectedValid := "dml 1.4;\n\ndevice dev;\n"
	if string(afterValid) != expectedValid {
		t.Errorf(
			"expected valid file reformatted to %q, got %q",
			expectedValid,
			string(afterValid),
		)
	}
}

func TestRun_Quiet_Files(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "test.dml")
	unformatted := "dml   1.4  ;\n  device   test_dev  ;\n"
	expected := "dml 1.4;\ndevice test_dev;\n"

	// 1. -q on unformatted file formats silently
	if err := os.WriteFile(dmlPath, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-q", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout with -q, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr with -q, got %q", stderr.String())
	}
	content, err := os.ReadFile(filepath.Clean(dmlPath))
	if err != nil {
		t.Fatalf("failed to read formatted file: %v", err)
	}
	if string(content) != expected {
		t.Errorf("expected content %q, got %q", expected, string(content))
	}

	// 2. --quiet on unformatted file formats silently
	if err := os.WriteFile(dmlPath, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to rewrite test file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--quiet", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout with --quiet, got %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr with --quiet, got %q", stderr.String())
	}

	// 3. -q on already formatted file remains silent
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-q", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf(
			"expected empty stdout on formatted file, got %q",
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr on formatted file, got %q",
			stderr.String(),
		)
	}

	// 4. -q on multiple files formats silently
	m1 := filepath.Join(tmpDir, "m1.dml")
	m2 := filepath.Join(tmpDir, "m2.dml")
	if err := os.WriteFile(m1, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to write m1: %v", err)
	}
	if err := os.WriteFile(m2, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to write m2: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-q", m1, m2}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf("expected ExitCodeOK (0) for multiple files, got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf(
			"expected empty stdout with -q on multiple files, got %q",
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr with -q on multiple files, got %q",
			stderr.String(),
		)
	}
}

func TestRun_Quiet_Files_Errors(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	badPath := filepath.Join(tmpDir, "bad.dml")
	if err := os.WriteFile(
		badPath,
		[]byte("dml 1.4; @@@ invalid"),
		0600,
	); err != nil {
		t.Fatalf("failed to write bad file: %v", err)
	}

	// Syntax error with -q should still output error on stderr
	var stdout, stderr bytes.Buffer
	code := run([]string{"-q", badPath}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf(
			"expected empty stdout on syntax error, got %q",
			stdout.String(),
		)
	}
	if !strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf(
			"expected syntax error on stderr with -q, got: %s",
			stderr.String(),
		)
	}

	// Non-existent file with -q should still output error on stderr
	stdout.Reset()
	stderr.Reset()
	code = run(
		[]string{"-q", filepath.Join(tmpDir, "nonexistent.dml")},
		&stdout,
		&stderr,
	)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf(
			"expected empty stdout on non-existent file, got %q",
			stdout.String(),
		)
	}
	if !strings.Contains(stderr.String(), "Error:") {
		t.Errorf("expected Error: on stderr with -q, got: %s", stderr.String())
	}
}

func TestRun_Quiet_Check(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "check_test.dml")
	unformatted := "dml   1.4  ;\n  device   check_dev  ;\n"
	formatted := "dml 1.4;\ndevice check_dev;\n"

	// Unformatted file with -q --check returns 1 with silent output
	if err := os.WriteFile(dmlPath, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-q", "--check", dmlPath}, &stdout, &stderr)
	if code != ExitCodeWouldReformat {
		t.Fatalf("expected ExitCodeWouldReformat (1), got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf(
			"expected empty stdout in check mode with -q, got %q",
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr in check mode with -q, got %q",
			stderr.String(),
		)
	}

	// Formatted file with -q --check returns 0 with silent output
	if err := os.WriteFile(dmlPath, []byte(formatted), 0600); err != nil {
		t.Fatalf("failed to write formatted file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-q", "--check", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf("expected ExitCodeOK (0), got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf(
			"expected empty stdout in check mode with -q, got %q",
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr in check mode with -q, got %q",
			stderr.String(),
		)
	}
}

func TestRun_Quiet_Stdin(t *testing.T) {
	t.Parallel()

	unformatted := "dml   1.4  ;\n  device   dev  ;\n"
	expected := "dml 1.4;\ndevice dev;\n"

	// 1. -q with unformatted stdin outputs formatted code to stdout, stderr empty
	var stdout, stderr bytes.Buffer
	code := runWithIO(
		[]string{"-q", "-"},
		strings.NewReader(unformatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr with -q on stdin, got %q",
			stderr.String(),
		)
	}

	// 2. --quiet with unformatted stdin
	stdout.Reset()
	stderr.Reset()
	code = runWithIO(
		[]string{"--quiet", "-"},
		strings.NewReader(unformatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr with --quiet on stdin, got %q",
			stderr.String(),
		)
	}

	// 3. -q with formatted stdin
	stdout.Reset()
	stderr.Reset()
	code = runWithIO(
		[]string{"-q", "-"},
		strings.NewReader(expected),
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr with -q on formatted stdin, got %q",
			stderr.String(),
		)
	}

	// 4. -q with --stdin-filename
	stdout.Reset()
	stderr.Reset()
	code = runWithIO(
		[]string{"-q", "--stdin-filename", "custom.dml", "-"},
		strings.NewReader(unformatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr with -q on custom.dml stdin, got %q",
			stderr.String(),
		)
	}
}

func TestRun_Quiet_Stdin_Check(t *testing.T) {
	t.Parallel()

	unformatted := "dml   1.4  ;\n  device   dev  ;\n"
	formatted := "dml 1.4;\ndevice dev;\n"

	// Unformatted stdin with -q --check returns 1 with silent output
	var stdout, stderr bytes.Buffer
	code := runWithIO(
		[]string{"-q", "--check", "-"},
		strings.NewReader(unformatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeWouldReformat {
		t.Fatalf("expected ExitCodeWouldReformat (1), got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf(
			"expected empty stdout with -q --check -, got %q",
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr with -q --check -, got %q",
			stderr.String(),
		)
	}

	// Formatted stdin with -q --check returns 0 with silent output
	stdout.Reset()
	stderr.Reset()
	code = runWithIO(
		[]string{"-q", "--check", "-"},
		strings.NewReader(formatted),
		&stdout,
		&stderr,
	)
	if code != ExitCodeOK {
		t.Fatalf("expected ExitCodeOK (0), got %d", code)
	}
	if stdout.Len() != 0 {
		t.Errorf(
			"expected empty stdout with -q --check -, got %q",
			stdout.String(),
		)
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr with -q --check -, got %q",
			stderr.String(),
		)
	}
}

func TestRun_Quiet_Stdin_Errors(t *testing.T) {
	t.Parallel()

	badInput := "dml 1.4; @@@ invalid\n"
	var stdout, stderr bytes.Buffer
	code := runWithIO(
		[]string{"-q", "-"},
		strings.NewReader(badInput),
		&stdout,
		&stderr,
	)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	if !strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf(
			"expected syntax error on stderr with -q -, got %s",
			stderr.String(),
		)
	}
	// Syntax error in stdin non-check mode echoes input to stdout
	if stdout.String() != badInput {
		t.Errorf("expected stdout to echo bad input, got %q", stdout.String())
	}
}

func TestRun_Quiet_Code(t *testing.T) {
	t.Parallel()

	unformatted := "dml   1.4  ;\n  device   dev  ;\n"
	expected := "dml 1.4;\ndevice dev;\n"

	// 1. -q -c outputs formatted code to stdout, stderr empty
	var stdout, stderr bytes.Buffer
	code := run([]string{"-q", "-c", unformatted}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr with -q -c, got %q", stderr.String())
	}

	// 2. --quiet --code
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--quiet", "--code", unformatted}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expected {
		t.Errorf("expected stdout %q, got %q", expected, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf(
			"expected empty stderr with --quiet --code, got %q",
			stderr.String(),
		)
	}

	// 3. Syntax error with -q -c still writes syntax error to stderr
	stdout.Reset()
	stderr.Reset()
	badCode := "dml 1.4; @@@ invalid\n"
	code = run([]string{"-q", "-c", badCode}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}
	if !strings.Contains(stderr.String(), "Syntax error") {
		t.Errorf(
			"expected syntax error on stderr with -q -c, got %s",
			stderr.String(),
		)
	}
}

func TestParseArgs_FastSafeToggles(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		args     []string
		wantSafe bool
		wantFast bool
	}{
		{"default without flags", []string{"file.dml"}, true, false},
		{"explicit safe", []string{"--safe", "file.dml"}, true, false},
		{"explicit fast", []string{"--fast", "file.dml"}, false, true},
		{
			"safe then fast",
			[]string{"--safe", "--fast", "file.dml"},
			false,
			true,
		},
		{
			"fast then safe",
			[]string{"--fast", "--safe", "file.dml"},
			true,
			false,
		},
		{"single dash safe", []string{"-safe", "file.dml"}, true, false},
		{"single dash fast", []string{"-fast", "file.dml"}, false, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var stderr bytes.Buffer
			cfg, err := parseArgs(tc.args, &stderr)
			if err != nil {
				t.Fatalf("unexpected error parsing %v: %v", tc.args, err)
			}
			if cfg.Safe != tc.wantSafe {
				t.Errorf("expected Safe=%v, got %v", tc.wantSafe, cfg.Safe)
			}
			if cfg.Fast != tc.wantFast {
				t.Errorf("expected Fast=%v, got %v", tc.wantFast, cfg.Fast)
			}
		})
	}
}

func TestRun_Safe_Files(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "safe_test.dml")
	unformatted := "dml   1.4  ;\n  device   safe_dev  ;\n"
	expected := "dml 1.4;\ndevice safe_dev;\n"

	// 1. Unformatted file with --safe succeeds and reformats
	if err := os.WriteFile(dmlPath, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"--safe", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	content, err := os.ReadFile(filepath.Clean(dmlPath))
	if err != nil {
		t.Fatalf("failed to read reformatted file: %v", err)
	}
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}

	// 2. Already formatted file with --safe succeeds
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--safe", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	// 3. Check mode with --safe on unformatted file returns 1
	if err := os.WriteFile(dmlPath, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to rewrite file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--safe", "--check", dmlPath}, &stdout, &stderr)
	if code != ExitCodeWouldReformat {
		t.Fatalf("expected ExitCodeWouldReformat (1), got %d", code)
	}

	// 4. Default run without flags succeeds and formats with safe mode
	if err := os.WriteFile(dmlPath, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to rewrite file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	content, err = os.ReadFile(filepath.Clean(dmlPath))
	if err != nil {
		t.Fatalf("failed to read reformatted file: %v", err)
	}
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}
}

func TestRun_Fast_Files(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dmlPath := filepath.Join(tmpDir, "fast_test.dml")
	unformatted := "dml   1.4  ;\n  device   fast_dev  ;\n"
	expected := "dml 1.4;\ndevice fast_dev;\n"

	// 1. Unformatted file with --fast succeeds
	if err := os.WriteFile(dmlPath, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"--fast", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK (0), got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	content, err := os.ReadFile(filepath.Clean(dmlPath))
	if err != nil {
		t.Fatalf("failed to read reformatted file: %v", err)
	}
	if string(content) != expected {
		t.Errorf("expected %q, got %q", expected, string(content))
	}

	// 2. Overriding --safe with --fast
	if err := os.WriteFile(dmlPath, []byte(unformatted), 0600); err != nil {
		t.Fatalf("failed to rewrite file: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--safe", "--fast", dmlPath}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf("expected ExitCodeOK (0), got %d", code)
	}
}

func TestProcessFormattingJob_SafeValidation(t *testing.T) {
	t.Parallel()

	ui := NewUI(&bytes.Buffer{}, &bytes.Buffer{}, ColorNever)
	stats := &noopStatsCollector{}

	unformatted := "dml   1.4  ;\n  device   dev  ;\n"
	buf := stream.NewHashedBufferFromString("test.dml", unformatted)

	job := InputJob{
		Path:      "test.dml",
		OrigMode:  0600,
		Size:      len(unformatted),
		Buffer:    buf,
		InputHash: buf.Hash(),
	}

	// 1. Safe mode formatting valid unformatted job succeeds
	outJob := processFormattingJob(job, 99, true, ui, stats)
	if outJob.Err != nil {
		t.Fatalf("unexpected error in safe mode: %v", outJob.Err)
	}
	if !outJob.WouldReformat {
		t.Errorf("expected WouldReformat=true")
	}

	// 2. Fast mode formatting valid unformatted job succeeds
	bufFast := stream.NewHashedBufferFromString("test.dml", unformatted)
	jobFast := InputJob{
		Path:      "test.dml",
		OrigMode:  0600,
		Size:      len(unformatted),
		Buffer:    bufFast,
		InputHash: bufFast.Hash(),
	}
	outJobFast := processFormattingJob(jobFast, 99, false, ui, stats)
	if outJobFast.Err != nil {
		t.Fatalf("unexpected error in fast mode: %v", outJobFast.Err)
	}
	if !outJobFast.WouldReformat {
		t.Errorf("expected WouldReformat=true")
	}
}

func TestRun_SummaryReport_EmptyDirectory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	// Normal run on empty directory
	code := run([]string{tmpDir}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	expectedMsg := "No DML files are present to be formatted. Nothing to do.\n"
	if stdout.String() != expectedMsg {
		t.Errorf("expected stdout %q, got %q", expectedMsg, stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("expected empty stderr, got %q", stderr.String())
	}

	// Check mode run on empty directory
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"--check", tmpDir}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.String() != expectedMsg {
		t.Errorf("expected stdout %q, got %q", expectedMsg, stdout.String())
	}

	// Quiet mode on empty directory suppresses notice
	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-q", tmpDir}, &stdout, &stderr)
	if code != ExitCodeOK {
		t.Fatalf(
			"expected ExitCodeOK, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected empty stdout with -q, got %q", stdout.String())
	}
}

func TestRun_SummaryReport_AllCategories(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	reformatPath := filepath.Join(tmpDir, "reformat.dml")
	unchangedPath := filepath.Join(tmpDir, "unchanged.dml")
	dml12Path := filepath.Join(tmpDir, "v12.dml")
	badPath := filepath.Join(tmpDir, "bad.dml")

	_ = os.WriteFile(reformatPath, []byte("dml 1.4; device dev;\n"), 0600)
	_ = os.WriteFile(unchangedPath, []byte("dml 1.4;\ndevice dev;\n"), 0600)
	_ = os.WriteFile(dml12Path, []byte("dml 1.2;\ncall $foo();\n"), 0600)
	_ = os.WriteFile(badPath, []byte("dml 1.4; @@@ invalid"), 0600)

	var stdout, stderr bytes.Buffer
	code := run([]string{"--no-color", tmpDir}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}

	outStr := stdout.String()
	expectedSummary := "1 file reformatted, 1 file left unchanged, " +
		"1 file ignored due to dml 1.2, 1 file failed to reformat.\n"
	if !strings.HasSuffix(outStr, expectedSummary) {
		t.Errorf(
			"expected stdout to end with %q, but got:\n%s",
			expectedSummary,
			outStr,
		)
	}
}

func TestRun_SummaryReport_CheckMode(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	reformatPath := filepath.Join(tmpDir, "reformat.dml")
	unchangedPath := filepath.Join(tmpDir, "unchanged.dml")
	dml12Path := filepath.Join(tmpDir, "v12.dml")
	badPath := filepath.Join(tmpDir, "bad.dml")

	_ = os.WriteFile(reformatPath, []byte("dml 1.4; device dev;\n"), 0600)
	_ = os.WriteFile(unchangedPath, []byte("dml 1.4;\ndevice dev;\n"), 0600)
	_ = os.WriteFile(dml12Path, []byte("dml 1.2;\ncall $foo();\n"), 0600)
	_ = os.WriteFile(badPath, []byte("dml 1.4; @@@ invalid"), 0600)

	var stdout, stderr bytes.Buffer
	code := run([]string{"--check", "--no-color", tmpDir}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}

	outStr := stdout.String()
	expectedSummary := "1 file would be reformatted, 1 file would be left unchanged, " +
		"1 file ignored due to dml 1.2, 1 file would fail to reformat.\n"
	if !strings.HasSuffix(outStr, expectedSummary) {
		t.Errorf(
			"expected stdout to end with %q, but got:\n%s",
			expectedSummary,
			outStr,
		)
	}
}

func TestRun_SummaryReport_Colored(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	reformatPath := filepath.Join(tmpDir, "reformat.dml")
	unchangedPath := filepath.Join(tmpDir, "unchanged.dml")
	dml12Path := filepath.Join(tmpDir, "v12.dml")
	badPath := filepath.Join(tmpDir, "bad.dml")

	_ = os.WriteFile(reformatPath, []byte("dml 1.4; device dev;\n"), 0600)
	_ = os.WriteFile(unchangedPath, []byte("dml 1.4;\ndevice dev;\n"), 0600)
	_ = os.WriteFile(dml12Path, []byte("dml 1.2;\ncall $foo();\n"), 0600)
	_ = os.WriteFile(badPath, []byte("dml 1.4; @@@ invalid"), 0600)

	var stdout, stderr bytes.Buffer
	code := run([]string{"--color", tmpDir}, &stdout, &stderr)
	if code != ExitCodeInternalError {
		t.Fatalf("expected ExitCodeInternalError (123), got %d", code)
	}

	outStr := stdout.String()
	expectedSummary := "\033[1;34m1 file \033[0m\033[1mreformatted\033[0m, " +
		"\033[34m1 file \033[0mleft unchanged, " +
		"\033[33m1 file \033[0mignored due to dml 1.2, " +
		"\033[31m1 file failed to reformat\033[0m.\n"
	if !strings.HasSuffix(outStr, expectedSummary) {
		t.Errorf("expected colored summary at end of stdout, got:\n%q", outStr)
	}
}

func TestVerboseMode(t *testing.T) {
	t.Parallel()

	t.Run("ProcessingAndExclusions_NoColor", func(t *testing.T) {
		tmpDir := t.TempDir()
		validPath := filepath.Join(tmpDir, "valid.dml")
		excludedDml := filepath.Join(tmpDir, "skip.dml")
		txtPath := filepath.Join(tmpDir, "readme.txt")
		cPath := filepath.Join(tmpDir, "helper.c")

		subDir := filepath.Join(tmpDir, "sub_excluded")
		_ = os.Mkdir(subDir, 0750)
		subDml := filepath.Join(subDir, "nested.dml")

		_ = os.WriteFile(validPath, []byte("dml 1.4;\ndevice dev;\n"), 0600)
		_ = os.WriteFile(excludedDml, []byte("dml 1.4;\ndevice skip;\n"), 0600)
		_ = os.WriteFile(txtPath, []byte("hello world"), 0600)
		_ = os.WriteFile(cPath, []byte("int x = 1;"), 0600)
		_ = os.WriteFile(subDml, []byte("dml 1.4;\ndevice nested;\n"), 0600)

		var stdout, stderr bytes.Buffer
		args := []string{
			"-v",
			"--no-color",
			"--extend-exclude", "skip|sub_excluded",
			tmpDir,
		}
		code := run(args, &stdout, &stderr)
		if code != ExitCodeOK {
			t.Fatalf(
				"expected ExitCodeOK (0), got %d (stderr: %s)",
				code,
				stderr.String(),
			)
		}

		outStr := stdout.String()

		// 1. Processing message for valid file
		expectedProc := fmt.Sprintf("Processing %s...\n", validPath)
		if !strings.Contains(outStr, expectedProc) {
			t.Errorf(
				"expected stdout to contain %q, got:\n%s",
				expectedProc,
				outStr,
			)
		}

		// 2. Directory ignored message
		expectedDirIgnored := fmt.Sprintf(
			"%s directory ignored: matches the --extend-exclude regular expression\n",
			subDir,
		)
		if !strings.Contains(outStr, expectedDirIgnored) {
			t.Errorf(
				"expected stdout to contain %q, got:\n%s",
				expectedDirIgnored,
				outStr,
			)
		}

		// 3. Excluded DML file ignored message
		expectedFileIgnored := fmt.Sprintf(
			"%s ignored: matches the --extend-exclude regular expression\n",
			excludedDml,
		)
		if !strings.Contains(outStr, expectedFileIgnored) {
			t.Errorf(
				"expected stdout to contain %q, got:\n%s",
				expectedFileIgnored,
				outStr,
			)
		}

		// 4. Unchanged message for valid file
		expectedUnchanged := fmt.Sprintf(
			"%s already well formatted, good job.\n",
			validPath,
		)
		if !strings.Contains(outStr, expectedUnchanged) {
			t.Errorf(
				"expected stdout to contain %q, got:\n%s",
				expectedUnchanged,
				outStr,
			)
		}

		// 5. Non-DML files must NOT be mentioned in stdout
		if strings.Contains(outStr, "readme.txt") {
			t.Errorf(
				"non-DML file readme.txt should be ignored silently, got:\n%s",
				outStr,
			)
		}
		if strings.Contains(outStr, "helper.c") {
			t.Errorf(
				"non-DML file helper.c should be ignored silently, got:\n%s",
				outStr,
			)
		}

		// 6. Must not contain ANSI escape codes
		if strings.Contains(outStr, "\033[") {
			t.Errorf(
				"stdout should not contain ANSI codes with --no-color, got:\n%s",
				outStr,
			)
		}
	})

	t.Run("ProcessingAndExclusions_Color", func(t *testing.T) {
		tmpDir := t.TempDir()
		validPath := filepath.Join(tmpDir, "valid.dml")
		excludedDml := filepath.Join(tmpDir, "skip.dml")

		_ = os.WriteFile(validPath, []byte("dml 1.4;\ndevice dev;\n"), 0600)
		_ = os.WriteFile(excludedDml, []byte("dml 1.4;\ndevice skip;\n"), 0600)

		var stdout, stderr bytes.Buffer
		args := []string{
			"-v",
			"--color",
			"--extend-exclude", "skip",
			tmpDir,
		}
		code := run(args, &stdout, &stderr)
		if code != ExitCodeOK {
			t.Fatalf(
				"expected ExitCodeOK (0), got %d (stderr: %s)",
				code,
				stderr.String(),
			)
		}

		outStr := stdout.String()

		expectedProc := fmt.Sprintf(
			"\033[90mProcessing %s...\033[0m\n",
			validPath,
		)
		if !strings.Contains(outStr, expectedProc) {
			t.Errorf(
				"expected colored processing message %q, got:\n%s",
				expectedProc,
				outStr,
			)
		}

		expectedIgnored := fmt.Sprintf(
			"\033[90m%s ignored: matches the --extend-exclude regular expression\033[0m\n",
			excludedDml,
		)
		if !strings.Contains(outStr, expectedIgnored) {
			t.Errorf(
				"expected colored ignored message %q, got:\n%s",
				expectedIgnored,
				outStr,
			)
		}

		expectedUnchanged := fmt.Sprintf(
			"\033[90m%s already well formatted, good job.\033[0m\n",
			validPath,
		)
		if !strings.Contains(outStr, expectedUnchanged) {
			t.Errorf(
				"expected colored unchanged message %q, got:\n%s",
				expectedUnchanged,
				outStr,
			)
		}
	})

	t.Run("QuietSuppressesVerbose", func(t *testing.T) {
		tmpDir := t.TempDir()
		validPath := filepath.Join(tmpDir, "valid.dml")
		_ = os.WriteFile(validPath, []byte("dml 1.4;\ndevice dev;\n"), 0600)

		var stdout, stderr bytes.Buffer
		code := run([]string{"-v", "-q", tmpDir}, &stdout, &stderr)
		if code != ExitCodeOK {
			t.Fatalf("expected ExitCodeOK (0), got %d", code)
		}
		if stdout.Len() != 0 {
			t.Errorf("expected empty stdout with -q, got:\n%s", stdout.String())
		}
	})

	t.Run("IncludeFilterIgnoresDML", func(t *testing.T) {
		tmpDir := t.TempDir()
		p1 := filepath.Join(tmpDir, "keep.dml")
		p2 := filepath.Join(tmpDir, "skip.dml")
		_ = os.WriteFile(p1, []byte("dml 1.4;\ndevice keep;\n"), 0600)
		_ = os.WriteFile(p2, []byte("dml 1.4;\ndevice skip;\n"), 0600)

		var stdout, stderr bytes.Buffer
		args := []string{
			"--verbose",
			"--no-color",
			"--include", "keep",
			tmpDir,
		}
		code := run(args, &stdout, &stderr)
		if code != ExitCodeOK {
			t.Fatalf("expected ExitCodeOK (0), got %d", code)
		}

		outStr := stdout.String()
		expectedIncludeIgnored := fmt.Sprintf(
			"%s ignored: matches the --include regular expression\n",
			p2,
		)
		if !strings.Contains(outStr, expectedIncludeIgnored) {
			t.Errorf(
				"expected stdout to contain %q, got:\n%s",
				expectedIncludeIgnored,
				outStr,
			)
		}
	})

	t.Run("StdinAndCodeDoNotEmitProcessing", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := runWithIO(
			[]string{"-v", "-"},
			strings.NewReader("dml 1.4;\ndevice dev;\n"),
			&stdout,
			&stderr,
		)
		if code != ExitCodeOK {
			t.Fatalf("expected ExitCodeOK (0), got %d", code)
		}
		if strings.Contains(stdout.String(), "Processing") {
			t.Errorf(
				"stdin should not emit Processing message, got:\n%s",
				stdout.String(),
			)
		}

		stdout.Reset()
		stderr.Reset()
		code = run(
			[]string{"-v", "-c", "dml 1.4;\ndevice dev;\n"},
			&stdout,
			&stderr,
		)
		if code != ExitCodeOK {
			t.Fatalf("expected ExitCodeOK (0), got %d", code)
		}
		if strings.Contains(stdout.String(), "Processing") {
			t.Errorf(
				"--code should not emit Processing message, got:\n%s",
				stdout.String(),
			)
		}
	})

	t.Run("Unchanged_CheckMode", func(t *testing.T) {
		tmpDir := t.TempDir()
		validPath := filepath.Join(tmpDir, "valid.dml")
		_ = os.WriteFile(validPath, []byte("dml 1.4;\ndevice dev;\n"), 0600)

		var stdout, stderr bytes.Buffer
		code := run(
			[]string{"-v", "--check", "--no-color", validPath},
			&stdout,
			&stderr,
		)
		if code != ExitCodeOK {
			t.Fatalf(
				"expected ExitCodeOK (0), got %d (stderr: %s)",
				code,
				stderr.String(),
			)
		}

		outStr := stdout.String()
		expectedProc := fmt.Sprintf("Processing %s...\n", validPath)
		if !strings.Contains(outStr, expectedProc) {
			t.Errorf(
				"expected stdout to contain %q, got:\n%s",
				expectedProc,
				outStr,
			)
		}
		expectedUnchanged := fmt.Sprintf(
			"%s already well formatted, good job.\n",
			validPath,
		)
		if !strings.Contains(outStr, expectedUnchanged) {
			t.Errorf(
				"expected stdout to contain %q, got:\n%s",
				expectedUnchanged,
				outStr,
			)
		}
		expectedReport := "1 file would be left unchanged.\n"
		if !strings.Contains(outStr, expectedReport) {
			t.Errorf(
				"expected stdout to contain report %q, got:\n%s",
				expectedReport,
				outStr,
			)
		}
	})

	t.Run("MixedFiles_Verbose", func(t *testing.T) {
		tmpDir := t.TempDir()
		cleanPath := filepath.Join(tmpDir, "clean.dml")
		dirtyPath := filepath.Join(tmpDir, "dirty.dml")

		_ = os.WriteFile(cleanPath, []byte("dml 1.4;\ndevice clean;\n"), 0600)
		_ = os.WriteFile(dirtyPath, []byte("dml 1.4;   device dirty;   "), 0600)

		var stdout, stderr bytes.Buffer
		code := run([]string{"-v", "--no-color", tmpDir}, &stdout, &stderr)
		if code != ExitCodeOK {
			t.Fatalf(
				"expected ExitCodeOK (0), got %d (stderr: %s)",
				code,
				stderr.String(),
			)
		}

		outStr := stdout.String()
		expectedCleanUnchanged := fmt.Sprintf(
			"%s already well formatted, good job.\n",
			cleanPath,
		)
		if !strings.Contains(outStr, expectedCleanUnchanged) {
			t.Errorf(
				"expected stdout to contain %q, got:\n%s",
				expectedCleanUnchanged,
				outStr,
			)
		}
		expectedDirtyReformatted := fmt.Sprintf("reformatted %s\n", dirtyPath)
		if !strings.Contains(outStr, expectedDirtyReformatted) {
			t.Errorf(
				"expected stdout to contain %q, got:\n%s",
				expectedDirtyReformatted,
				outStr,
			)
		}
		expectedReport := "1 file reformatted, 1 file left unchanged.\n"
		if !strings.Contains(outStr, expectedReport) {
			t.Errorf(
				"expected stdout to contain report %q, got:\n%s",
				expectedReport,
				outStr,
			)
		}
	})

	t.Run("GitIgnoreAndIgnoreFile_Verbose", func(t *testing.T) {
		tmpDir := t.TempDir()
		_ = os.Mkdir(filepath.Join(tmpDir, ".git"), 0750)
		_ = os.WriteFile(
			filepath.Join(tmpDir, ".gitignore"),
			[]byte("git_skip.dml\ngit_ignored_dir/\nnotes.txt\n"),
			0600,
		)
		_ = os.WriteFile(
			filepath.Join(tmpDir, ".ignore"),
			[]byte("dotignore_skip.dml\ndotignore_dir/\n"),
			0600,
		)

		keepPath := filepath.Join(tmpDir, "keep.dml")
		gitSkipPath := filepath.Join(tmpDir, "git_skip.dml")
		gitDir := filepath.Join(tmpDir, "git_ignored_dir")
		dotIgnoreSkipPath := filepath.Join(tmpDir, "dotignore_skip.dml")
		dotIgnoreDir := filepath.Join(tmpDir, "dotignore_dir")
		notesPath := filepath.Join(tmpDir, "notes.txt")

		_ = os.Mkdir(gitDir, 0750)
		_ = os.Mkdir(dotIgnoreDir, 0750)
		_ = os.WriteFile(keepPath, []byte("dml 1.4;\ndevice keep;\n"), 0600)
		_ = os.WriteFile(
			gitSkipPath,
			[]byte("dml 1.4;\ndevice gskip;\n"),
			0600,
		)
		_ = os.WriteFile(
			filepath.Join(gitDir, "inner.dml"),
			[]byte("dml 1.4;\ndevice gdir;\n"),
			0600,
		)
		_ = os.WriteFile(
			dotIgnoreSkipPath,
			[]byte("dml 1.4;\ndevice dskip;\n"),
			0600,
		)
		_ = os.WriteFile(
			filepath.Join(dotIgnoreDir, "inner.dml"),
			[]byte("dml 1.4;\ndevice ddir;\n"),
			0600,
		)
		_ = os.WriteFile(notesPath, []byte("ignored text"), 0600)

		var stdout, stderr bytes.Buffer
		args := []string{"-v", "--no-color", tmpDir}
		code := run(args, &stdout, &stderr)
		if code != ExitCodeOK {
			t.Fatalf(
				"expected ExitCodeOK (0), got %d (stderr: %s)",
				code,
				stderr.String(),
			)
		}

		outStr := stdout.String()
		expectedMsgs := []string{
			fmt.Sprintf("Processing %s...\n", keepPath),
			fmt.Sprintf(
				"%s ignored: matches the .gitignore file content\n",
				gitSkipPath,
			),
			fmt.Sprintf(
				"%s directory ignored: matches the .gitignore file content\n",
				gitDir,
			),
			fmt.Sprintf(
				"%s ignored: matches the .ignore file content\n",
				dotIgnoreSkipPath,
			),
			fmt.Sprintf(
				"%s directory ignored: matches the .ignore file content\n",
				dotIgnoreDir,
			),
			fmt.Sprintf("%s already well formatted, good job.\n", keepPath),
		}
		for _, want := range expectedMsgs {
			if !strings.Contains(outStr, want) {
				t.Errorf(
					"expected stdout to contain %q, got:\n%s",
					want,
					outStr,
				)
			}
		}

		if strings.Contains(outStr, "notes.txt") {
			t.Errorf(
				"non-DML file notes.txt should not appear in verbose output, got:\n%s",
				outStr,
			)
		}
	})
}
