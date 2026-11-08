// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"github.com/timoschenko/dmltools/pkg/dmlex"
	dmldefaults "github.com/timoschenko/dmltools/pkg/printer/dml_defaults"
	"github.com/timoschenko/dmltools/pkg/printer/format"
)

// PrintContext provides scoped FormatFlag overrides via a parent-chain lookup.
// Rules push a new PrintContext with specific overrides; Lookup walks the chain
// and falls back to defaultFlagFor(tt) if no override is found.
type PrintContext struct {
	parent    *PrintContext
	overrides format.TokenMap
}

// NewPrintContext creates a new PrintContext with the given parent.
func NewPrintContext(
	overrides format.TokenMap,
	parent *PrintContext,
) *PrintContext {
	if overrides == nil {
		panic("NewPrintContext: overrides cannot be nil")
	}
	return &PrintContext{
		parent:    parent,
		overrides: overrides,
	}
}

// GetFormatFlag searches for the FormatFlag for the given TokenType in this context
// and its parent chain, falling back to defaultFlagFor if no override is found.
func (pc *PrintContext) GetFormatFlag(tt dmlex.TokenType) format.FormatFlag {
	if pc == nil {
		return 0 // Or handle appropriately if settings is missing, but typically pc != nil
	}

	if m, ok := pc.overrides.(format.TokenMapImpl); ok {
		if flag, ok := m[tt]; ok {
			return flag
		}
	} else if m, ok := pc.overrides.(*format.TokenMapImpl); ok {
		if flag, ok := (*m)[tt]; ok {
			return flag
		}
	} else {
		// Fallback if it's some other implementation of TokenMap
		flag := pc.overrides.GetFormatFlag(tt)
		if flag != 0 {
			return flag
		}
	}

	if pc.parent != nil {
		return pc.parent.GetFormatFlag(tt)
	}
	// TODO: Pass Printer context thru whole stack
	return DefaultFlagFor(tt) // leak from "settings"
}

// Parent returns the parent context of this context.
//
//go:inline
func (pc *PrintContext) Parent() *PrintContext {
	return pc.parent
}

// DefaultFlagFor returns the default FormatFlag for a given TokenType.
//
//go:inline
func DefaultFlagFor(tt dmlex.TokenType) format.FormatFlag {
	return dmldefaults.Tokens().GetFormatFlag(tt)
}
