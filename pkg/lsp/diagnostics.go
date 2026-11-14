// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package lsp

import (
	"errors"

	"github.com/antlr4-go/antlr/v4"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/fs/stream"
	"go.lsp.dev/protocol"
)

const copyTextIntoTokenOnParse = true

func init() {
	// NOTE: Enable "writing" all token content inside AST during parsing.
	// It is expected that ANTLR runtime will call HashedBuffer.GetText()
	// and save a pointer to string from the buffer instead of copying it.
	antlr.CommonTokenFactoryDEFAULT = antlr.NewCommonTokenFactory(
		copyTextIntoTokenOnParse,
	)
}

// Diagnose parses DML content and returns syntax errors as LSP diagnostics,
// along with the parsed AST (which is non-nil unless completely unparseable).
func Diagnose(content string) ([]protocol.Diagnostic, antlr.ParseTree) {
	diags, _, tree := diagnoseWithTokens(content)
	return diags, tree
}

func diagnoseWithTokens(
	content string,
) ([]protocol.Diagnostic, *antlr.CommonTokenStream, antlr.ParseTree) {
	errListener := dmlex.NewErrorListener("memory", 0)

	tokens, tree := parseStreamAndTokens("memory", content, errListener)

	if len(errListener.Errors) == 0 {
		return []protocol.Diagnostic{}, tokens, tree
	}

	lineOffsets := computeLineOffsets(content)
	diags := make([]protocol.Diagnostic, 0, len(errListener.Errors))
	for _, e := range errListener.Errors {
		lineIdx := 0
		if e.Line > 0 {
			lineIdx = e.Line - 1
		}
		line := uint32(lineIdx) //nolint:gosec // lineIdx is non-negative

		lineText := ""
		if lineIdx < len(lineOffsets) {
			lineStart := lineOffsets[lineIdx]
			lineEnd := len(content)
			if lineIdx+1 < len(lineOffsets) {
				lineEnd = lineOffsets[lineIdx+1]
			}
			lineText = content[lineStart:lineEnd]
		}
		col := runeColToUTF16Col(lineText, e.Column)

		length := uint32(1)
		if e.Symbol != "" && e.Symbol != "<EOF>" {
			if symLen := utf16Len(e.Symbol); symLen > 0 {
				length = symLen
			}
		} else if e.Length > 0 {
			length = uint32(e.Length) //nolint:gosec // e.Length is positive
		}

		severity := protocol.DiagnosticSeverityError
		var code protocol.ProgressToken
		if errors.Is(e.Err, dmlex.EUNIMPL) {
			severity = protocol.DiagnosticSeverityWarning
			code = protocol.String("EUNIMPL")
		}

		diags = append(diags, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: line, Character: col},
				End:   protocol.Position{Line: line, Character: col + length},
			},
			Severity: severity,
			Code:     code,
			Source:   protocol.NewOptional("dml"),
			Message:  protocol.String(e.Msg),
		})
	}

	return diags, tokens, tree
}

func parseStreamAndTokens(
	name string,
	content string,
	errListener *dmlex.ErrorListener,
) (*antlr.CommonTokenStream, antlr.ParseTree) {
	input := stream.NewHashedBufferFromString(name, content)
	file := dmlex.ParseStream(input, errListener)
	tokens, _ := file.Parser.GetTokenStream().(*antlr.CommonTokenStream)
	tree := file.Parser.Dml()
	if tokens != nil {
		tokens.Fill()
	}
	return tokens, tree
}
