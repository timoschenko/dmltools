// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

//go:inline
func (p *Printer) extractToken(tokenIdx TokenIndex) antlr.Token {
	return p.getTokenStream().Get(tokenIdx.Int())
}

//go:inline
func (p *Printer) gatherTokenIndexFromRightiestToken(
	token antlr.Token,
) TokenIndex {
	if token == nil {
		panic("token is nil")
	}

	tokenIdx := token.GetTokenIndex()
	// lookup for comments on the same line
	lineno := token.GetLine()
	ts := p.getTokenStream()
	commetns := ts.GetHiddenTokensToRight(tokenIdx, antlr.TokenHiddenChannel)
	for commentIdx := range commetns {
		if commetns[commentIdx].GetLine() != lineno {
			break
		}

		tokenIdx = commetns[commentIdx].GetTokenIndex()
	}

	return TokenIndex(tokenIdx)
}

//go:inline
func (p *Printer) gatherTokenIndexFromRightiestTerminalNode(
	node antlr.TerminalNode,
) TokenIndex {
	if node == nil {
		panic("node is nil")
	}

	token := node.GetSymbol()
	return p.gatherTokenIndexFromRightiestToken(token)
}

// gatherTokenIndexFromRightiest reads trailing hidden tokens (like comments)
// on the same line as the given token, returning the index of the right-most
// hidden token found.
//
//go:inline
func (p *Printer) gatherTokenIndexFromRightiest(token any) TokenIndex {
	switch t := token.(type) {
	case antlr.TerminalNode:
		return p.gatherTokenIndexFromRightiestTerminalNode(t)
	case antlr.Token:
		return p.gatherTokenIndexFromRightiestToken(t)
	default:
		panic("type is not supported")
	}
}

// getTokenStream returns the underlying ANTLR CommonTokenStream used by the Printer.
//
//go:inline
func (p *Printer) getTokenStream() *antlr.CommonTokenStream {
	return p.cts
}

type trimmedToken struct {
	antlr.Token
	text string
}

func (t *trimmedToken) GetText() string {
	return t.text
}

func trimBlockCommentClosingLine(line string) string {
	line = strings.TrimRight(line, " \t\r")
	if !strings.HasSuffix(line, "*/") {
		return line
	}

	starCount := 0
	for i := len(line) - 2; i >= 0 && line[i] == '*'; i-- {
		starCount++
	}
	delimLen := starCount + 1
	if delimLen >= len(line) {
		return line
	}

	content := line[:len(line)-delimLen]
	delim := line[len(line)-delimLen:]

	trimmedContent := strings.TrimRight(content, " \t\r")
	if len(trimmedContent) == len(content) {
		return line
	}

	// If content is purely leading indentation before */, preserve it.
	if trimmedContent == "" {
		return line
	}

	// Otherwise, collapse trailing whitespace before */ to a single space.
	return trimmedContent + " " + delim
}

func trimCommentToken(
	token antlr.Token,
	flag format.FormatFlag,
) antlr.Token {
	tt := dmlex.GetTokenType(token)
	if !tt.IsComment() {
		return token
	}

	text := token.GetText()
	hasNewline := strings.HasSuffix(text, "\n")
	body := text
	if hasNewline {
		body = text[:len(text)-1]
	}
	body = strings.TrimSuffix(body, "\r")

	if tt == dmlex.TokenComment {
		trimmed := strings.TrimRight(body, " \t\r")
		var newText string
		if hasNewline {
			sepAfter := flag.GetSeparatorAfter()
			if sepAfter != format.SepMandatoryWrap &&
				sepAfter != format.SepWrapButNoSpace {
				newText = trimmed
			} else {
				newText = trimmed + "\n"
			}
		} else {
			newText = trimmed
		}
		if newText != text {
			return &trimmedToken{Token: token, text: newText}
		}
		return token
	}

	// Block comments: TokenCommentBlock, TokenCommentBlockEol,
	// TokenCommentInline, TokenCommentInlineEol
	body = strings.TrimRight(body, " \t")

	if strings.Contains(body, "\n") {
		lines := strings.Split(body, "\n")
		for i := range lines {
			if i == len(lines)-1 {
				lines[i] = trimBlockCommentClosingLine(lines[i])
			} else {
				lines[i] = strings.TrimRight(lines[i], " \t\r")
			}
		}
		body = strings.Join(lines, "\n")
	} else {
		body = trimBlockCommentClosingLine(body)
	}

	var newText string
	if hasNewline {
		sepAfter := flag.GetSeparatorAfter()
		if sepAfter != format.SepMandatoryWrap &&
			sepAfter != format.SepWrapButNoSpace {
			newText = body
		} else {
			newText = body + "\n"
		}
	} else {
		newText = body
	}

	if newText != text {
		return &trimmedToken{Token: token, text: newText}
	}
	return token
}

func (p *Printer) processToken(token antlr.Token, ctx *PrintContext) error {
	if ctx == nil {
		ctx = p.context
	}

	tt := dmlex.GetTokenType(token)
	if tt == dmlex.TokenWs || tt == antlr.TokenEOF {
		return nil
	}

	flag := ctx.GetFormatFlag(tt)
	token = trimCommentToken(token, flag)

	if err := p.LayoutEngine.Push(token, flag); err != nil {
		return err
	}

	return nil
}

func (p *Printer) pushTokenRange(endIdx TokenIndex) error {
	for i := p.lastProcessedTokenIndex; i <= endIdx; i++ {
		token := p.extractToken(i)
		tt := dmlex.GetTokenType(token)

		if tt == dmlex.TokenWs {
			p.lastProcessedTokenIndex = i + 1
			continue
		}
		if tt == antlr.TokenEOF {
			p.lastProcessedTokenIndex = i + 1
			continue
		}

		flag := p.context.GetFormatFlag(tt)
		token = trimCommentToken(token, flag)

		if err := p.LayoutEngine.Push(token, flag); err != nil {
			return err
		}

		p.lastProcessedTokenIndex = i + 1
	}
	return nil
}

//go:inline
func (p *Printer) pushSynthecic(
	token *dmlex.SyntheticToken,
) error {
	tt := dmlex.GetTokenType(token)
	flag := p.context.GetFormatFlag(tt)
	if err := p.LayoutEngine.Push(token, flag); err != nil {
		return err
	}

	return nil
}

//go:inline
func (p *Printer) pushTokenWithCustomFormat(
	node antlr.TerminalNode,
	flagFormatter func(format.FormatFlag) format.FormatFlag,
) error {

	token := node.GetSymbol()
	endIdx := TokenIndex(token.GetTokenIndex())
	if p.lastProcessedTokenIndex < endIdx {
		if err := p.pushTokenRange(endIdx - 1); err != nil {
			return err
		}
	}

	tt := dmlex.GetTokenType(token)
	assertLazy(func() bool {
		switch tt {
		case dmlex.TokenWs, antlr.TokenEOF:
			// has to be skipped
			return false
		case dmlex.TokenCommentInlineEol, dmlex.TokenCommentBlockEol:
			// trimmed
			return false
		}
		return true
	}, "has to be dml-code token")
	flagOrign := p.context.GetFormatFlag(tt)
	flag := flagFormatter(flagOrign)
	if err := p.LayoutEngine.Push(token, flag); err != nil {
		return err
	}

	p.lastProcessedTokenIndex = endIdx + 1
	return nil
}

//go:inline
func (p *Printer) pushTokenRangeFromRightiest(token antlr.Tree) error {
	stopToken := castStopToken(token)
	endToken := p.gatherTokenIndexFromRightiest(stopToken)
	return p.pushTokenRange(endToken)
}

// printSimpleRuleContextEndsSemi
// Pushes tokens via pushTokenRange with a SEMI override context flushes the engine.
func (p *Printer) printSimpleRuleContextEndsSemi(
	node dmlex.ISimpleRuleContextEndsSemi,
) error {
	assertLazy(func() bool {
		const expectedChildren = 3
		tree := node.GetChildren()
		return len(tree) == expectedChildren
	}, "unexpected number of children")

	tokenIndexEnd := p.gatherTokenIndexFromRightiestTerminalNode(node.SEMI())
	if err := p.pushTokenRange(tokenIndexEnd); err != nil {
		return err
	}
	return p.Flush()
}

// writeNewLines writes out the requested number of new lines using the LayoutEngine.
// It preserves writeNewLine's deduplication: if lastWrittenChar is already '\n',
// one fewer newline is emitted.
//
// Warning: This method resets the Printer's linecap tracking and should typically
// only be called by container rules between statements.
func (p *Printer) writeNewLines(amount int) error {
	if p.LayoutEngine.GetLastChar_Workaround() == '\n' {
		amount--
	}

	if amount <= 0 {
		return nil
	}

	var text string
	switch amount {
	case 1:
		text = "\n"
	case 2:
		text = "\n\n"
	case 3:
		text = "\n\n\n"
	default:
		text = strings.Repeat("\n", amount)
	}

	if err := p.LayoutEngine.WriteRaw(text); err != nil {
		return err
	}
	return nil
}

// printUnimplemented extracts and writes the raw text of an unimplemented parser rule context
// directly from the input stream, preserving its original formatting.
//
//nolint:unused // Retained for scaffolding new grammar rules per SKILL.md
func (p *Printer) printUnimplemented(node antlr.ParserRuleContext) error {
	startToken := node.GetStart()
	if startToken == nil {
		panic("unreachable: " + fmt.Sprintf("%T", node))
	}

	// Look for comments from left?
	tokenIndexEnd := p.gatherTokenIndexFromRightiestToken(startToken) - 1
	if tokenIndexEnd > p.lastProcessedTokenIndex {
		if err := p.pushTokenRange(tokenIndexEnd); err != nil {
			return err
		}
	}

	stopToken := node.GetStop()
	if stopToken == nil {
		panic("unreachable: " + fmt.Sprintf("%T", node))
	}
	// expect stopToken to have the same
	input := startToken.GetInputStream()
	if input == nil {
		panic("unreachable: " + fmt.Sprintf("%T", node))
	}
	fileSize := input.Size()
	startOffset := startToken.GetStart()
	stopOffset := stopToken.GetStop()
	if startOffset > stopOffset || startOffset > fileSize ||
		stopOffset > fileSize {
		panic("unreachable: " + fmt.Sprintf("%T", node))
	}
	text := input.GetTextFromInterval(
		antlr.NewInterval(startOffset, stopOffset),
	)
	if err := p.LayoutEngine.WriteRaw(text); err != nil {
		return err
	}

	stopTokenIndex := getTokenIndexFromToken(stopToken)
	rightiestToken := p.gatherTokenIndexFromRightiestToken(stopToken)
	// Mark stopToken as processed
	p.lastProcessedTokenIndex = stopTokenIndex.Next()
	if rightiestToken != stopTokenIndex {
		if err := p.pushTokenRange(rightiestToken); err != nil {
			return err
		}
	}

	return nil
}

// isSignificantToken returns true if the token is not whitespace or EOF.
func isSignificantToken(token antlr.Token) bool {
	if token == nil {
		return false
	}
	tt := dmlex.GetTokenType(token)
	return tt != dmlex.TokenWs && tt != antlr.TokenEOF
}
