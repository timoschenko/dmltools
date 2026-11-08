// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

import (
	"unicode/utf8"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"

	"github.com/antlr4-go/antlr/v4"
)

// Deprecated: Only layoutEngine is aware of LastChar and it is not intended to be set externally.
// This method is exposed temporarily to bridge legacy output and should not be used for new code.
//
// Deprecated: This is workaround
func (e *Engine) GetLastChar_Workaround() rune {
	for i := len(e.buffer) - 1; i >= 0; i-- {
		bt := &e.buffer[i]
		if bt.IsSynthetic() {
			continue
		}
		if bt.Token == nil {
			if bt.ResolvedSepBefore == format.SepMandatoryWrap ||
				bt.ResolvedSepBefore == format.SepWrapButNoSpace {
				return '\n'
			}
			continue
		}
		text := bt.Token.GetText()
		if len(text) > 0 {
			r, _ := utf8.DecodeLastRuneInString(text)
			return r
		}
	}

	return e.state.LastChar
}

func (e *Engine) overrideResolveDesign_Workaround(
	token antlr.Token,
	bt *BufferedToken,
	prevToken antlr.Token,
	prevTokenType dmlex.TokenType,
) bool {
	tokenType := dmlex.TokenEOF
	if token != nil {
		tokenType = dmlex.GetTokenType(token)
	}

	sameLine := false
	if prevToken != nil && token != nil {
		sameLine = token.GetLine() == prevToken.GetLine()
	}

	switch {
	// () - always has no space, nowrap
	case prevTokenType == dmlex.TokenLparen && tokenType == dmlex.TokenRparen:
		bt.ResolveWrappingCache = ResolveWrappingCache{
			format.SepNone,
			format.PriorityMax,
			format.WinnerNone,
			// no groups assigned
			e.currentActiveGroup(),
			e.currentActiveGroup(),
		}
	// `, // comment` avoid wrapping
	case prevTokenType == dmlex.TokenComma && tokenType == dmlex.TokenComment:
		bt.ResolveWrappingCache = ResolveWrappingCache{
			format.SepSpace,
			format.PriorityMax,
			format.WinnerNone,
			// no groups assigned
			e.currentActiveGroup(),
			e.currentActiveGroup(),
		}
	// TODO: AI WA: Resolve line comment separation for delimiters and statement boundaries (Issue #19, #24)
	// {{{
	case tokenType == dmlex.TokenComment && !sameLine &&
		(prevTokenType == dmlex.TokenLbrace ||
			prevTokenType == dmlex.TokenLparen ||
			prevTokenType == dmlex.TokenLbrack ||
			prevTokenType == dmlex.TokenColon ||
			prevTokenType == dmlex.TokenSemi):
		bt.ResolveWrappingCache = ResolveWrappingCache{
			format.SepMandatoryWrap,
			format.PriorityMax,
			format.WinnerRight,
			nil,
			e.currentActiveGroup(),
		}
	case prevTokenType == dmlex.TokenLbrace && tokenType == dmlex.TokenComment && sameLine:
		bt.ResolveWrappingCache = ResolveWrappingCache{
			format.SepSpace,
			format.PriorityMax,
			format.WinnerNone,
			nil,
			e.currentActiveGroup(),
		}
	// }}}
	default:
		return false
	}

	return true
}

// Expose for testing purposeses
func TesOnly_SimulateBuffer(
	leBuffer []BufferedToken,
	state State,
	stopOnExceed bool,
	tabWidth int,
) (int, int) {
	return simulateBuffer(
		leBuffer, state, stopOnExceed, tabWidth,
	)
}

// Expose for testing purposeses
func (e *Engine) TesOnly_GetState() State {
	return *e.state
}

// Expose for testing purposeses
func (e *Engine) TesOnly_GetLeBuffer() []BufferedToken {
	return e.buffer
}

// Expose for testing purposeses
func (e *Engine) TesOnly_ApplyBreak(idx int) {
	e.applyBreak(idx)
}
