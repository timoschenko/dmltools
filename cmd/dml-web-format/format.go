// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"bytes"
	"fmt"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer"
	"github.com/timoschenko/dmltools/pkg/printer/layout"

	"github.com/antlr4-go/antlr/v4"
)

type TokenInfo struct {
	Type      int    `json:"type"`
	Name      string `json:"name"`
	Text      string `json:"text"`
	Start     int    `json:"start"`
	Stop      int    `json:"stop"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	GroupPath string `json:"groupPath,omitempty"`
}

type FormatResult struct {
	Result      string      `json:"result,omitempty"`
	Error       string      `json:"error,omitempty"`
	Line        int         `json:"line,omitempty"`
	Column      int         `json:"column,omitempty"`
	Length      int         `json:"length,omitempty"`
	SyntaxError bool        `json:"syntax_error,omitempty"`
	Trace       string      `json:"trace,omitempty"`
	TraceCount  int         `json:"trace_count,omitempty"`
	Tokens      []TokenInfo `json:"tokens,omitempty"`
}

func FormatDML(
	inputString string,
	lineCap int,
	enableTrace ...bool,
) FormatResult {
	if lineCap <= 0 {
		lineCap = 109
	}

	shouldTrace := true
	if len(enableTrace) > 0 {
		shouldTrace = enableTrace[0]
	}

	errorListener := dmlex.NewErrorListener("", 0)
	inputStream := antlr.NewInputStream(inputString)
	file := dmlex.ParseStream(inputStream, errorListener)

	if len(errorListener.Errors) > 0 {
		firstErr := errorListener.Errors[0]
		return FormatResult{
			Error: fmt.Sprintf(
				"line %d:%d: %s", firstErr.Line, firstErr.Column, firstErr.Msg,
			),
			Line:        firstErr.Line,
			Column:      firstErr.Column,
			Length:      firstErr.Length,
			SyntaxError: true,
		}
	}

	tree := file.Parser.Dml()
	if len(errorListener.Errors) > 0 {
		firstErr := errorListener.Errors[0]
		errStr := fmt.Sprintf(
			"line %d:%d: %s", firstErr.Line, firstErr.Column, firstErr.Msg,
		)
		return FormatResult{
			Error:       errStr,
			Line:        firstErr.Line,
			Column:      firstErr.Column,
			Length:      firstErr.Length,
			SyntaxError: true,
		}
	}
	if file.Parser.HasError() {
		return FormatResult{
			Error:       "syntax error in DML",
			SyntaxError: true,
		}
	}

	var tracer *layout.JSONTokenTracer
	if shouldTrace {
		tracer = layout.NewJSONTokenTracer()
	}

	var buf bytes.Buffer
	var printErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				printErr = fmt.Errorf("%v", r)
			}
		}()
		var t layout.TokenTracer
		if tracer != nil {
			t = tracer
		}
		printErr = printer.PrinterWriteToDebug(tree, file, &buf, lineCap, t)
	}()

	if printErr != nil {
		return FormatResult{
			Error: printErr.Error(),
		}
	}

	res := FormatResult{
		Result: buf.String(),
	}

	formattedStream := antlr.NewInputStream(res.Result)
	lexer := dmlex.NewdmlLexer(formattedStream)
	for {
		t := lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		res.Tokens = append(res.Tokens, TokenInfo{
			Type:   t.GetTokenType(),
			Name:   dmlex.TokenType(t.GetTokenType()).String(),
			Text:   t.GetText(),
			Start:  t.GetStart(),
			Stop:   t.GetStop(),
			Line:   t.GetLine(),
			Column: t.GetColumn(),
		})
	}

	if tracer != nil {
		traceBytes, err := tracer.ToJSON()
		if err == nil {
			res.Trace = string(traceBytes)
			res.TraceCount = len(tracer.Entries)
		}
		if len(tracer.Entries) == len(res.Tokens) {
			for i := range res.Tokens {
				if tracer.Entries[i].Args != nil {
					res.Tokens[i].GroupPath = tracer.Entries[i].Args.Right.GroupPath
				}
			}
		} else {
			entryIdx := 0
			for i := range res.Tokens {
				for entryIdx < len(tracer.Entries) {
					if tracer.Entries[entryIdx].Args == nil {
						entryIdx++
						continue
					}
					right := tracer.Entries[entryIdx].Args.Right
					if right.TokenType == res.Tokens[i].Type &&
						right.Text == res.Tokens[i].Text {
						res.Tokens[i].GroupPath = right.GroupPath
						entryIdx++
						break
					}
					entryIdx++
				}
			}
		}
	}
	return res
}
