// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format

import (
	"fmt"
)

// FormatFlag is a 32-bit bitfield encoding spacing rules, break priorities,
// indentation shifts, and pair conflict priorities for a token.
//
// Bitfield Layout:
//   - Bits 0-3:   Line break before priority (0-15)
//   - Bits 4-6:   Separator Before (Separator)
//   - Bit 7:      Reserved
//   - Bits 8-11:  Line break after priority (0-15)
//   - Bits 12-14: Separator After (Separator)
//   - Bit 15:     Reserved
//   - Bit 16:     Hard Indent increase (triggers only if a line break occurs here)
//   - Bit 17:     Hard Indent decrease (triggers only if a line break occurs here)
//   - Bit 18:     ForceWrap (indent bits describe wrap context, not real indentation)
//   - Bits 19-23: Reserved
//   - Bits 24-27: Token Before Priority (0-15) — used for pair conflict resolution
//   - Bits 28-31: Token After Priority (0-15) — used for pair conflict resolution
type FormatFlag uint32

const (
	PriorityMask uint32 = 0xF
	SepMask      uint32 = 0x7

	BreakBeforePriorityOffset uint32 = 0
	SepBeforeOffset           uint32 = 4
	BreakAfterPriorityOffset  uint32 = 8
	SepAfterOffset            uint32 = 12
	IndentIncOffset           uint32 = 16
	IndentDecOffset           uint32 = 17
	ForceWrapOffset           uint32 = 18
	TokenBeforePriorityOffset uint32 = 24
	TokenAfterPriorityOffset  uint32 = 28

	BreakBeforePriorityMask FormatFlag = FormatFlag(
		PriorityMask << BreakBeforePriorityOffset,
	)
	SepBeforeMask          FormatFlag = FormatFlag(SepMask << SepBeforeOffset)
	BreakAfterPriorityMask FormatFlag = FormatFlag(
		PriorityMask << BreakAfterPriorityOffset,
	)
	SepAfterMask            FormatFlag = FormatFlag(SepMask << SepAfterOffset)
	IndentIncMask           FormatFlag = 1 << IndentIncOffset
	IndentDecMask           FormatFlag = 1 << IndentDecOffset
	ForceWrapMask           FormatFlag = 1 << ForceWrapOffset
	TokenBeforePriorityMask FormatFlag = FormatFlag(
		PriorityMask << TokenBeforePriorityOffset,
	)
	TokenAfterPriorityMask FormatFlag = FormatFlag(
		PriorityMask << TokenAfterPriorityOffset,
	)

	priorityMask = PriorityMask
	sepMask      = SepMask

	breakBeforePriorityOffset = BreakBeforePriorityOffset
	sepBeforeOffset           = SepBeforeOffset
	breakAfterPriorityOffset  = BreakAfterPriorityOffset
	sepAfterOffset            = SepAfterOffset
	indentIncOffset           = IndentIncOffset
	indentDecOffset           = IndentDecOffset
	forceWrapOffset           = ForceWrapOffset
	tokenBeforePriorityOffset = TokenBeforePriorityOffset
	tokenAfterPriorityOffset  = TokenAfterPriorityOffset

	breakBeforePriorityMask = BreakBeforePriorityMask
	sepBeforeMask           = SepBeforeMask
	breakAfterPriorityMask  = BreakAfterPriorityMask
	sepAfterMask            = SepAfterMask
	indentIncMask           = IndentIncMask
	indentDecMask           = IndentDecMask
	forceWrapMask           = ForceWrapMask
	tokenBeforePriorityMask = TokenBeforePriorityMask
	tokenAfterPriorityMask  = TokenAfterPriorityMask
)

func (f FormatFlag) GetBreakBeforePriority() Priority {
	return Priority(f & breakBeforePriorityMask)
}

func (f FormatFlag) GetSeparatorBefore() Separator {
	return Separator((f & sepBeforeMask) >> sepBeforeOffset)
}

func (f FormatFlag) GetBreakAfterPriority() Priority {
	return Priority((f & breakAfterPriorityMask) >> breakAfterPriorityOffset)
}

func (f FormatFlag) GetSeparatorAfter() Separator {
	return Separator((f & sepAfterMask) >> sepAfterOffset)
}

func (f FormatFlag) HasIndentIncrease() bool {
	return (f & indentIncMask) != 0
}

func (f FormatFlag) HasIndentDecrease() bool {
	return (f & indentDecMask) != 0
}

func (f FormatFlag) IsForceWrap() bool {
	return (f & forceWrapMask) != 0
}

func (f FormatFlag) GetTokenBeforePriority() Priority {
	return Priority((f & tokenBeforePriorityMask) >> tokenBeforePriorityOffset)
}

func (f FormatFlag) GetTokenAfterPriority() Priority {
	return Priority((f & tokenAfterPriorityMask) >> tokenAfterPriorityOffset)
}

// GetTokenPriority is a backward-compatible shim that returns the maximum of the two split priorities.
func (f FormatFlag) GetTokenPriority() Priority {
	tbp := f.GetTokenBeforePriority()
	tap := f.GetTokenAfterPriority()
	if tbp > tap {
		return tbp
	}
	return tap
}

func (f FormatFlag) WithBreakBeforePriority(p Priority) FormatFlag {
	return (f & ^breakBeforePriorityMask) | (FormatFlag(uint32(p)&priorityMask) << breakBeforePriorityOffset)
}

func (f FormatFlag) WithSeparatorBefore(s Separator) FormatFlag {
	return (f & ^sepBeforeMask) | (FormatFlag(uint32(s)&sepMask) << sepBeforeOffset)
}

func (f FormatFlag) WithBreakAfterPriority(p Priority) FormatFlag {
	return (f & ^breakAfterPriorityMask) | (FormatFlag(uint32(p)&priorityMask) << breakAfterPriorityOffset)
}

func (f FormatFlag) WithSeparatorAfter(s Separator) FormatFlag {
	return (f & ^sepAfterMask) | (FormatFlag(uint32(s)&sepMask) << sepAfterOffset)
}

func (f FormatFlag) WithIndentIncrease() FormatFlag {
	return f | indentIncMask
}

func (f FormatFlag) WithoutIndentIncrease() FormatFlag {
	return f & ^indentIncMask
}

func (f FormatFlag) WithIndentDecrease() FormatFlag {
	return f | indentDecMask
}

func (f FormatFlag) WithoutIndentDecrease() FormatFlag {
	return f & ^indentDecMask
}

func (f FormatFlag) WithForceWrap() FormatFlag {
	return f | forceWrapMask
}

func (f FormatFlag) WithSeparatorBeforePriority(p Priority) FormatFlag {
	return (f & ^tokenBeforePriorityMask) | (FormatFlag(uint32(p)&priorityMask) << tokenBeforePriorityOffset)
}

func (f FormatFlag) WithSeparatorAfterPriority(p Priority) FormatFlag {
	return (f & ^tokenAfterPriorityMask) | (FormatFlag(uint32(p)&priorityMask) << tokenAfterPriorityOffset)
}

// WithSeparatorPriority is a backward-compatible shim that sets both TokenBeforePriority and TokenAfterPriority.
func (f FormatFlag) WithSeparatorPriority(p Priority) FormatFlag {
	return f.WithSeparatorBeforePriority(p).WithSeparatorAfterPriority(p)
}

func (f FormatFlag) String() string {
	return fmt.Sprintf(
		"FormatFlag{BBP:%d SB:%s BAP:%d SA:%s IndentInc:%v IndentDec:%v ForceWrap:%v TBP:%d TAP:%d}",
		f.GetBreakBeforePriority(),
		f.GetSeparatorBefore(),
		f.GetBreakAfterPriority(),
		f.GetSeparatorAfter(),
		f.HasIndentIncrease(),
		f.HasIndentDecrease(),
		f.IsForceWrap(),
		f.GetTokenBeforePriority(),
		f.GetTokenAfterPriority(),
	)
}
