// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
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
}
