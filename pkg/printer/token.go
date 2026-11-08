// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import "github.com/antlr4-go/antlr/v4"

type TokenIndex int

// Next returns the index of the next token after the given tokenIdx.
//
//go:inline
func (t TokenIndex) Next() TokenIndex {
	return TokenIndex(int(t) + 1)
}

// Int returns the integer value of the TokenIndex.
//
//go:inline
func (t TokenIndex) Int() int {
	return int(t)
}

//go:inline
func TokenIndexFromToken(token antlr.Token) TokenIndex {
	if token == nil {
		panic("token is nil")
	}
	return TokenIndex(token.GetTokenIndex())
}

//go:inline
func TokenIndexFromTerminalNode(node antlr.TerminalNode) TokenIndex {
	if node == nil {
		panic("node is nil")
	}
	return TokenIndex(node.GetSymbol().GetTokenIndex())
}

//go:inline
func TokenIndexFromStart(ctx antlr.ParserRuleContext) TokenIndex {
	if ctx == nil {
		panic("ctx is nil")
	}
	return TokenIndex(ctx.GetStart().GetTokenIndex())
}

//go:inline
func TokenIndexFromStop(ctx antlr.ParserRuleContext) TokenIndex {
	return TokenIndex(ctx.GetStop().GetTokenIndex())
}

// getTokenIndexFromToken safely extracts the TokenIndex from an antlr.Token,
// panicking if it's nil.
//
//go:inline
func getTokenIndexFromToken(token antlr.Token) TokenIndex {
	if token == nil {
		panic("token is nil")
	}
	return TokenIndex(token.GetTokenIndex())
}

// castStartToken extracts the starting antlr.Token from a given antlr.Tree node
// (context or terminal).
//
//go:inline
func castStartToken(token antlr.Tree) antlr.Token {
	switch v := token.(type) {
	case antlr.ParserRuleContext:
		return v.GetStart()
	case antlr.TerminalNode:
		return v.GetSymbol()
	default:
		panic("antlr.Tree has no Token")
	}
}

// castStopToken extracts the stopping antlr.Token from a given antlr.Tree node.
//
//go:inline
func castStopToken(token antlr.Tree) antlr.Token {
	switch v := token.(type) {
	case antlr.ParserRuleContext:
		return v.GetStop()
	case antlr.TerminalNode:
		return v.GetSymbol()
	default:
		panic("antlr.Tree has no Token")
	}
}
