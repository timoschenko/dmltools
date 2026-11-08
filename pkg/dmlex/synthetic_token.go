// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmlex

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

const IvalidTokenIndex = antlr.TokenEpsilon

func readTextFromInputStream(t antlr.Token) string {
	input := t.GetInputStream()
	if input == nil {
		return ""
	}
	start := t.GetStart()
	stop := t.GetStop()
	n := input.Size()
	if start >= 0 && stop >= 0 && start < n && stop < n {
		return input.GetTextFromInterval(antlr.NewInterval(start, stop))
	}
	return "<EOF>"
}

// TokenSource is a workarounds needed to unroll
// private antlr.TokenSourceCharStreamPair structure.
type TokenSource struct {
	source *antlr.TokenSourceCharStreamPair

	tokenSource antlr.TokenSource
	charStream  antlr.CharStream
}

func NewTokenSouce(donor antlr.Token) *TokenSource {
	return &TokenSource{
		source:      donor.GetSource(),
		tokenSource: donor.GetTokenSource(),
		charStream:  donor.GetInputStream(),
	}
}

type SyntheticToken struct {
	source    *TokenSource
	tokenType TokenType // token type of the token
	line      int       // line=1..n of the 1st character
	column    int       // beginning of the line at which it occurs, 0..n-1
	text      string    // text of the token.
}

var _ antlr.Token = (*SyntheticToken)(nil)

func NewSyntheticToken(
	source *TokenSource, tp TokenType, text string, line, column int,
) *SyntheticToken {
	if source == nil {
		panic("NewToken - source can not be nil")
	}
	return &SyntheticToken{
		source:    source,
		tokenType: tp,
		line:      1,
		column:    0,
		text:      text,
	}
}

// Create a "fake" token on the same line but just after.
//
// Token will overlap if next "real" token on the same line.
func NewSyntheticTokenAfter(
	donor antlr.Token, tp TokenType, text string,
) *SyntheticToken {
	src := NewTokenSouce(donor)

	line := donor.GetLine()         // keep same line
	column := donor.GetColumn() + 1 // insert just after

	return NewSyntheticToken(src, tp, text, line, column)
}

func NewFakeLBrace(source *TokenSource) *SyntheticToken {
	return NewSyntheticToken(source, TokenLbrace, "{", 1, 1)
}

func NewFakeRBrace(source *TokenSource) *SyntheticToken {
	return NewSyntheticToken(source, TokenRbrace, "}", 1, 1)
}

func (t *SyntheticToken) GetChannel() int {
	return antlr.TokenDefaultChannel
}

func (t *SyntheticToken) SetChannel(channel int) {
	panic("dmlex.Token.SetChannel is readonly")
}

func (t *SyntheticToken) GetStart() int {
	return -1
}

func (t *SyntheticToken) SetStart(start int) {
	panic("dmlex.Token.SetStart is readonly")
}

func (t *SyntheticToken) GetStop() int {
	return -1
}

func (t *SyntheticToken) SetStop(stop int) {
	panic("dmlex.Token.SetStop is readonly")
}

func (t *SyntheticToken) GetLine() int {
	return t.line
}

func (t *SyntheticToken) SetLine(line int) {
	panic("dmlex.Token.SetLine is readonly")
}

func (t *SyntheticToken) GetColumn() int {
	return t.column
}

func (t *SyntheticToken) SetColumn(column int) {
	panic("dmlex.Token.SetColumn is readonly")
}

func (t *SyntheticToken) GetTokenType() int {
	return int(t.tokenType)
}

func (t *SyntheticToken) SetTokenType(tt TokenType) {
	panic("dmlex.Token.SetTokenType is readonly")
}

func (t *SyntheticToken) TokenType() TokenType {
	return t.tokenType
}

func (t *SyntheticToken) GetSource() *antlr.TokenSourceCharStreamPair {
	return t.source.source
}

func (t *SyntheticToken) SetSource(source *antlr.TokenSourceCharStreamPair) {
	panic("dmlex.Token.SetSource is readonly")
}

func (t *SyntheticToken) GetText() string {
	if t.text != "" {
		return t.text
	}

	return readTextFromInputStream(t)
}

func (t *SyntheticToken) SetText(s string) {
	panic("dmlex.Token.SetText is readonly")
}

func (t *SyntheticToken) GetTokenIndex() int {
	return IvalidTokenIndex
}

func (t *SyntheticToken) SetTokenIndex(v int) {
	panic("dmlex.Token.SetTokenIndex is readonly")
}

func (t *SyntheticToken) GetTokenSource() antlr.TokenSource {
	return t.source.tokenSource
}

func (t *SyntheticToken) GetInputStream() antlr.CharStream {
	return t.source.charStream
}

func (t *SyntheticToken) String() string {
	txt := t.GetText()
	if txt != "" {
		txt = strings.ReplaceAll(txt, "\n", "\\n")
		txt = strings.ReplaceAll(txt, "\r", "\\r")
		txt = strings.ReplaceAll(txt, "\t", "\\t")
	} else {
		txt = "<no text>"
	}

	return "[@dmlex.SyntheticToken='" + txt + "',<" +
		strconv.Itoa(int(t.tokenType)) + ">" +
		"," + strconv.Itoa(t.line) + ":" + strconv.Itoa(t.column) + "]"
}
