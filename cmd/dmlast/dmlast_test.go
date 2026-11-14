// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrintVersion(t *testing.T) {
	var buf bytes.Buffer
	printVersion(&buf)
	out := buf.String()

	if !strings.HasPrefix(out, "dmlast version ") {
		t.Errorf(
			"expected version output to start with 'dmlast version ', got: %q",
			out,
		)
	}
	if !strings.Contains(out, "License:") ||
		!strings.Contains(out, "GPL-2.0-only") {
		t.Errorf("expected license information in version output, got: %q", out)
	}
	if !strings.Contains(out, previewWarning) {
		t.Errorf(
			"expected preview warning in version output, got: %q",
			out,
		)
	}
}

func TestHelpFlag(t *testing.T) {
	for _, flag := range []string{"--help", "-h", "-help"} {
		var stdout, stderr bytes.Buffer
		code := run([]string{flag}, &stdout, &stderr)
		if code != 0 {
			t.Errorf("expected exit code 0 for %s, got %d", flag, code)
		}
		out := stdout.String()
		if !strings.Contains(out, "Usage: dmlast") {
			t.Errorf("expected usage output on stdout, got: %q", out)
		}
		if !strings.Contains(out, previewWarning) {
			t.Errorf(
				"expected preview warning in usage output, got: %q",
				out,
			)
		}
	}
}

func TestRunPreviewWarning(t *testing.T) {
	tmpDir := t.TempDir()
	dmlFile := filepath.Join(tmpDir, "test.dml")
	if err := os.WriteFile(
		dmlFile,
		[]byte("dml 1.4;\ndevice dev;\n"),
		0o600,
	); err != nil {
		t.Fatalf("failed to write temp DML file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-file", dmlFile}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf(
			"expected exit code 0, got %d (stderr: %s)",
			code,
			stderr.String(),
		)
	}

	if !strings.Contains(stderr.String(), previewWarning) {
		t.Errorf(
			"expected preview warning on stderr, got: %q",
			stderr.String(),
		)
	}

	var parsed JsonNode
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		t.Errorf("expected valid JSON on stdout, got error: %v", err)
	}
}
