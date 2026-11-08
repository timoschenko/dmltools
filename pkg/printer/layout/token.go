// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

import (
	"unicode/utf8"

	"github.com/timoschenko/dmltools/pkg/printer/format"

	"github.com/antlr4-go/antlr/v4"
)

type noCopy struct{}

// FlaggedToken pairs a parsed ANTLR token with its computed formatting flag.
type FlaggedToken struct {
	Token antlr.Token
	Flag  format.FormatFlag
}

// ResolveWrappingCache holds the resolved line-wrapping configuration (separator and
// break priority) computed between a token and the token immediately preceding it.
type ResolveWrappingCache struct {
	ResolvedSepBefore     format.Separator
	ResolvedBreakPriority format.Priority
	Winner                format.Winner
	AssingedGroup         *Group
	CurrentGroup          *Group
}

// Simple alias. Feel free to change in case of overflow
type LineLength uint32

// Type annotation helper
const LengthOfEmptyString LineLength = 0

type TokenStringCache struct {
	_ noCopy

	// Pre-calculated metrics to avoid redundant string scanning
	PrefixWidth  LineLength
	SuffixWidth  LineLength
	HasNewlines  bool
	MiddleWidths []LineLength
}

func NewTokenStringCache(text string) *TokenStringCache {
	runeCount := LengthOfEmptyString
	lastNLRuneCount := LengthOfEmptyString

	cache := TokenStringCache{}
	for i := 0; i < len(text); {
		r, size := utf8.DecodeRuneInString(text[i:])
		runeCount++
		if r == '\n' {
			if !cache.HasNewlines {
				cache.PrefixWidth = runeCount - 1
				cache.HasNewlines = true
			} else {
				cache.MiddleWidths = append(
					cache.MiddleWidths, (runeCount-1)-lastNLRuneCount,
				)
			}
			lastNLRuneCount = runeCount
		}
		i += size
	}
	if cache.HasNewlines {
		cache.SuffixWidth = runeCount - lastNLRuneCount
	} else {
		cache.PrefixWidth = runeCount
	}
	return &cache
}

// SyntheticLayoutConfig defines special internal markers injected into the layout buffer
// to represent structural or layout-specific boundaries without consuming ANTLR tokens.
type SyntheticLayoutConfig uint32

const (
	// SlcNone represents a standard token without any synthetic layout significance.
	SlcNone SyntheticLayoutConfig = iota

	// SlcWrappingGroupBegin marks the start of a logical wrapping group.
	// When a line wrap occurs within this group, the layout engine will attempt to
	// apply breaks collectively to all peers at this nesting level.
	SlcWrappingGroupBegin

	// SlcWrappingGroupEnd marks the end of a logical wrapping group.
	SlcWrappingGroupEnd

	// SlcSoftIdentBegin marks the start of a logical soft ident group.
	// The layout engine will attempt to increase ident level for new line wraps
	// within this group, however if the line wrap increasing ident by itself
	// (via by FormatFlags+ResolveWrappingCache) then soft ident is "absorbed"
	// and ignored.
	SlcSoftIdentBegin

	// SlcSoftIdentEnd marks the end of a logical soft ident group.
	SlcSoftIdentEnd
)

// BufferedToken represents a token along with its formatting metadata
// and the separator that was resolved between it and the previous token.
type BufferedToken struct {
	_ noCopy

	FlaggedToken
	ResolveWrappingCache

	synthetic SyntheticLayoutConfig
	cache     *TokenStringCache
}

func (bt *BufferedToken) IsSynthetic() bool {
	return bt.synthetic != SlcNone
}

func (bt *BufferedToken) Sanitize() {
	if bt.synthetic != SlcNone {
		if bt.Token != nil {
			panic("BufferedToken and synthetic cannot be set at the same time")
		}
	} else if bt.Token == nil {
		panic("BufferedToken and synthetic cannot be nil at the same time")
	}
}
