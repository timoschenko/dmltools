// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpFlag(t *testing.T) {
	for _, flag := range []string{"--help", "-h", "-help"} {
		var stdout, stderr bytes.Buffer
		stdin := strings.NewReader("")

		code := run([]string{flag}, stdin, &stdout, &stderr)
		if code != 0 {
			t.Errorf("expected exit code 0 for %s, got %d", flag, code)
		}

		out := stdout.String()
		if !strings.Contains(out, "Usage: dmlsp") {
			t.Errorf("expected usage output on stdout, got: %s", out)
		}
		if !strings.Contains(out, previewWarning) {
			t.Errorf(
				"expected preview warning in usage output, got: %s",
				out,
			)
		}
	}
}

func TestVersionFlag(t *testing.T) {
	for _, flag := range []string{"--version", "-version", "-V"} {
		var stdout, stderr bytes.Buffer
		stdin := strings.NewReader("")

		code := run([]string{flag}, stdin, &stdout, &stderr)
		if code != 0 {
			t.Errorf("expected exit code 0 for %s, got %d", flag, code)
		}

		out := stdout.String()
		if !strings.Contains(out, "dmlsp version") {
			t.Errorf("expected version output on stdout, got: %s", out)
		}
		if !strings.Contains(out, "GPL-2.0-only") {
			t.Errorf("expected license info on stdout, got: %s", out)
		}
		if !strings.Contains(out, previewWarning) {
			t.Errorf(
				"expected preview warning in version output, got: %s",
				out,
			)
		}
	}
}

func TestStartupPreviewWarning(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	_ = run(nil, stdin, &stdout, &stderr)
	if !strings.Contains(stderr.String(), previewWarning) {
		t.Errorf(
			"expected preview warning on stderr at startup, got: %q",
			stderr.String(),
		)
	}
}

func TestInvalidFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")

	code := run(
		[]string{"--invalid-flag-that-does-not-exist"},
		stdin,
		&stdout,
		&stderr,
	)
	if code != 2 {
		t.Errorf("expected exit code 2 for invalid flag, got %d", code)
	}
}
