// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer_test

import (
	"bytes"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer"
	"github.com/timoschenko/dmltools/pkg/printer/format"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

type mockToken struct {
	antlr.Token
	text string
}

func (t *mockToken) GetText() string {
	return t.text
}

func (t *mockToken) GetLine() int {
	return 1
}

func (t *mockToken) GetTokenType() int {
	return 0
}

func TestLayoutEngineIntegration_PassthroughFlush(t *testing.T) {
	var buf bytes.Buffer
	input := antlr.NewInputStream("")
	file := dmlex.ParseStream(input)
	p := printer.NewPrinter(80, file, &buf)

	// Manually push a token to the layout engine
	tok := &mockToken{text: "BUFFERED_TOKEN"}
	flag := format.FormatFlag(0) // No special flags
	if err := p.LayoutEngine.Push(tok, flag); err != nil {
		t.Fatalf("Failed to push token: %v", err)
	}

	// Now call WriteRaw via a public-facing way or directly if possible.
	// We want to verify that printUnimplemented (which uses WriteRaw) flushes the engine.

	rawText := "UNIMPLEMENTED_TEXT"
	if err := p.LayoutEngine.WriteRaw(rawText); err != nil {
		t.Fatalf("WriteRaw failed: %v", err)
	}

	expected := "BUFFERED_TOKENUNIMPLEMENTED_TEXT"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}
