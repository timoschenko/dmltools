// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

//go:build js && wasm

package main

import (
	"syscall/js"
)

func formatDML(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return map[string]any{
			"error": "expected at least 1 argument (DML string)",
		}
	}
	inputString := args[0].String()

	lineCap := 109 // default MaxLineCap
	if len(args) >= 2 && args[1].Type() == js.TypeNumber {
		lineCap = args[1].Int()
	}

	enableTrace := true
	if len(args) >= 3 && args[2].Type() == js.TypeBoolean {
		enableTrace = args[2].Bool()
	}

	res := FormatDML(inputString, lineCap, enableTrace)
	out := map[string]any{}
	if res.Result != "" {
		out["result"] = res.Result
	}
	if res.Error != "" {
		out["error"] = res.Error
	}
	if res.SyntaxError {
		out["syntax_error"] = true
		out["line"] = res.Line
		out["column"] = res.Column
		out["length"] = res.Length
	}
	if res.Trace != "" {
		out["trace"] = res.Trace
		out["trace_count"] = res.TraceCount
	}
	if res.Tokens != nil {
		jsTokens := make([]any, len(res.Tokens))
		for i, tok := range res.Tokens {
			tokMap := map[string]any{
				"type":   tok.Type,
				"name":   tok.Name,
				"text":   tok.Text,
				"start":  tok.Start,
				"stop":   tok.Stop,
				"line":   tok.Line,
				"column": tok.Column,
			}
			if tok.GroupPath != "" {
				tokMap["groupPath"] = tok.GroupPath
			}
			jsTokens[i] = tokMap
		}
		out["tokens"] = jsTokens
	}
	return out
}

func main() {
	js.Global().Set("formatDML", js.FuncOf(formatDML))
	// Block forever to keep the WebAssembly module running
	select {}
}
