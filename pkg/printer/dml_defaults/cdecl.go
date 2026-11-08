// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmldefaults

import (
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"
)

var cdeclDefaults staticTokenMap

// Cdecl returns the baseline FormatFlag for every TokenType in the C-declaration context.
func Cdecl() format.TokenMap {
	cdeclDefaults.once.Do(initCdeclDefaults)
	return &cdeclDefaults.data
}

func initCdeclDefaults() {
	tfd := Tokens()
	cdeclDefaults.data = make(format.TokenMapImpl)

	// Helper to get the base flag: use Expression override if it exists,
	// otherwise fall back to Tokens.
	exprOverrides := Expression().(*format.TokenMapImpl)

	//go:inline
	getFlag := func(tt dmlex.TokenType) format.FormatFlag {
		if f, ok := cdeclDefaults.data[tt]; ok {
			return f
		}
		if f, ok := (*exprOverrides)[tt]; ok {
			return f
		}
		return tfd.GetFormatFlag(tt)
	}

	// 1. Inherit relevant structural rules from Expression
	for tt := range *exprOverrides {
		cdeclDefaults.data[tt] = getFlag(tt)
	}

	// 2. Add Cdecl specific overrides
	cdeclDefaults.data[dmlex.TokenArrow] = tfd.GetFormatFlag(dmlex.TokenArrow).
		WithSeparatorBefore(format.SepSpace).
		WithBreakBeforePriority(format.PriorityAboveNormal).
		// just a bit more than ()
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorAfterPriority(format.PriorityAboveNormal)
	// just a bit more than ()

	cdeclDefaults.data[dmlex.TokenMul] = getFlag(dmlex.TokenMul).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepNone).
		WithSeparatorPriority(format.PriorityAboveNormal)

	cdeclDefaults.data[dmlex.TokenConst] = getFlag(dmlex.TokenConst).
		WithSeparatorBefore(format.SepSpace)

	cdeclDefaults.data[dmlex.TokenBitand] = getFlag(dmlex.TokenBitand).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepNone)

	// Lparen in cdecl never has a space before it (unlike in expressions)
	cdeclDefaults.data[dmlex.TokenLparen] = tfd.GetFormatFlag(dmlex.TokenLparen)
	cdeclDefaults.data[dmlex.TokenRparen] = tfd.GetFormatFlag(dmlex.TokenRparen)

	cdeclDefaults.data[dmlex.TokenColon] = getFlag(dmlex.TokenColon).
		WithSeparatorBefore(format.SepNone).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorPriority(format.PriorityHigh)

	cdeclDefaults.data[dmlex.TokenAt] = getFlag(dmlex.TokenAt).
		WithSeparatorAfter(format.SepSpace).
		WithSeparatorPriority(format.PriorityHigh)

	clingers := []dmlex.TokenType{
		dmlex.TokenIdentifier,
		dmlex.TokenRegister,
		dmlex.TokenAttribute,
		dmlex.TokenT__0,
		dmlex.TokenAsync,
		dmlex.TokenAuto,
		dmlex.TokenAwait,
		dmlex.TokenBank,
		dmlex.TokenBitorder,
		dmlex.TokenCall,
		dmlex.TokenClass,
		dmlex.TokenConnect,
		dmlex.TokenConstant,
		dmlex.TokenCritical,
		dmlex.TokenData,
		dmlex.TokenDevice,
		dmlex.TokenEndian,
		dmlex.TokenEnum,
		dmlex.TokenEvent,
		dmlex.TokenField,
		dmlex.TokenFooter,
		dmlex.TokenGroup,
		dmlex.TokenHeader,
		dmlex.TokenHeader2,
		dmlex.TokenImplement,
		dmlex.TokenImport,
		dmlex.TokenIndependent,
		dmlex.TokenInfo,
		dmlex.TokenInterface,
		dmlex.TokenLoggroup,
		dmlex.TokenMemoized,
		dmlex.TokenMethod,
		dmlex.TokenNamespace,
		dmlex.TokenNothrow,
		dmlex.TokenParam,
		dmlex.TokenPort,
		dmlex.TokenPrivate,
		dmlex.TokenProtected,
		dmlex.TokenProvisional,
		dmlex.TokenPublic,
		dmlex.TokenRestrict,
		dmlex.TokenSaved,
		dmlex.TokenSelect,
		dmlex.TokenSize,
		dmlex.TokenSpecViol,
		dmlex.TokenStartup,
		dmlex.TokenStatic,
		dmlex.TokenSubdevice,
		dmlex.TokenThen,
		dmlex.TokenThrows,
		dmlex.TokenUnimpl,
		dmlex.TokenUnion,
		dmlex.TokenUsing,
		dmlex.TokenVirtual,
		dmlex.TokenVolatile,
		dmlex.TokenWith,
		dmlex.TokenHook,
		dmlex.TokenSequence,
	}
	for _, tt := range clingers {
		cdeclDefaults.data[tt] = getFlag(tt).
			WithSeparatorAfterPriority(format.PriorityLow)
	}

	cdeclDefaults.data[dmlex.TokenTypeof] = getFlag(dmlex.TokenTypeof).
		WithSeparatorAfter(format.SepNone).
		WithSeparatorAfterPriority(format.PriorityAboveHighest)
	// has to beat `(` which has high prio

	cdeclDefaults.data[dmlex.TokenLbrace] = getFlag(dmlex.TokenLbrace).
		WithSeparatorBefore(format.SepSpace).
		WithSeparatorAfter(format.SepMandatoryWrap).
		WithSeparatorPriority(format.PriorityMax)
	cdeclDefaults.data[dmlex.TokenRbrace] = getFlag(dmlex.TokenRbrace).
		WithSeparatorBefore(format.SepMandatoryWrap).
		WithSeparatorBeforePriority(format.PriorityMax)
}
