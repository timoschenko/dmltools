// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

//go:build !dmltools.ndebug

package layout_test

import (
	"bytes"
	"encoding/json"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer"
	"github.com/timoschenko/dmltools/pkg/printer/layout"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

func TestJSONTokenTracer(t *testing.T) {
	// Simple DML code snippet
	code := `dml 1.4;`

	stream := antlr.NewInputStream(code)
	file := dmlex.ParseStream(stream)
	if file.Parser.HasError() {
		t.Fatalf("Failed to parse DML")
	}

	tracer := layout.NewJSONTokenTracer()
	var buf bytes.Buffer
	p := printer.NewPrinterDebug(80, file, &buf, tracer)

	// In dmlex, the root is typically Dml() or similar.
	// We just need a valid node.
	tree := file.Parser.Dml()
	if err := p.Print(tree); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	if err := p.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	// Verify tracer captured entries
	if len(tracer.Entries) == 0 {
		t.Errorf("Expected tracer to capture entries, got 0")
	}

	for _, entry := range tracer.Entries {
		if entry.PID == 0 {
			t.Errorf("Expected PID > 0, got %d", entry.PID)
		}
		if entry.Name != "TokenFormatDecision" {
			t.Errorf("Expected name 'TokenFormatDecision', got %q", entry.Name)
		}
		if entry.Args == nil {
			t.Errorf("Expected Args to be populated")
		}
	}

	// Serialize to JSON
	jsonBytes, err := tracer.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize to JSON: %v", err)
	}

	jsonStr := string(jsonBytes)
	if !strings.Contains(jsonStr, "TokenFormatDecision") {
		t.Errorf("Expected JSON to contain 'TokenFormatDecision'")
	}

	var unmarshaled []map[string]any
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	if len(unmarshaled) != len(tracer.Entries) {
		t.Errorf(
			"Expected %d entries in JSON, got %d",
			len(tracer.Entries),
			len(unmarshaled),
		)
	}
}
