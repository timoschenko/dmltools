// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package printer

import (
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer"
	"github.com/timoschenko/dmltools/pkg/printer/format"
	"testing"
)

func TestPrintContext_LookupDefault(t *testing.T) {
	ctx := printer.NewPrintContext(format.EmptyTokenMap, nil)
	tt := dmlex.TokenIdentifier
	got := ctx.GetFormatFlag(tt)
	want := printer.DefaultFlagFor(tt)
	if got != want {
		t.Errorf("Lookup() = %v, want default %v", got, want)
	}
}

func TestPrintContext_LookupOverride(t *testing.T) {
	tt := dmlex.TokenIdentifier
	flag := format.FormatFlag(100)
	overrides := format.NewTokenMap(
		map[dmlex.TokenType]format.FormatFlag{
			tt: flag,
		},
	)
	ctx := printer.NewPrintContext(overrides, nil)

	got := ctx.GetFormatFlag(tt)
	if got != flag {
		t.Errorf("Lookup() = %v, want override %v", got, flag)
	}
}

func TestPrintContext_ParentChain(t *testing.T) {
	tt1 := dmlex.TokenIdentifier
	tt2 := dmlex.TokenLparen
	flag1 := format.FormatFlag(100)
	flag2 := format.FormatFlag(200)
	flag3 := format.FormatFlag(300)

	parentOverrides := format.NewTokenMap(
		map[dmlex.TokenType]format.FormatFlag{
			tt1: flag1,
			tt2: flag2,
		},
	)
	parent := printer.NewPrintContext(parentOverrides, nil)

	childOverrides := format.NewTokenMap(
		map[dmlex.TokenType]format.FormatFlag{
			tt1: flag3,
		},
	)
	child := printer.NewPrintContext(childOverrides, parent)

	// Child override
	if got := child.GetFormatFlag(tt1); got != flag3 {
		t.Errorf("Lookup(tt1) = %v, want child override %v", got, flag3)
	}

	// Parent fallback
	if got := child.GetFormatFlag(tt2); got != flag2 {
		t.Errorf("Lookup(tt2) = %v, want parent fallback %v", got, flag2)
	}

	// Global fallback
	got := child.GetFormatFlag(dmlex.TokenRparen)
	if got != printer.DefaultFlagFor(dmlex.TokenRparen) {
		t.Errorf("Lookup(TokenRparen) = %v, want global fallback", got)
	}
}

func TestPrintContext_ThreeLevelChain(t *testing.T) {
	tt := dmlex.TokenIdentifier
	flagGP := format.FormatFlag(1)
	flagP := format.FormatFlag(2)
	flagC := format.FormatFlag(3)

	gp := printer.NewPrintContext(format.NewTokenMap(
		map[dmlex.TokenType]format.FormatFlag{tt: flagGP}), nil,
	)
	p := printer.NewPrintContext(format.EmptyTokenMap, gp)
	c := printer.NewPrintContext(format.EmptyTokenMap, p)

	if got := c.GetFormatFlag(tt); got != flagGP {
		t.Errorf("Grandparent lookup failed: got %v, want %v", got, flagGP)
	}

	p = printer.NewPrintContext(format.NewTokenMap(
		map[dmlex.TokenType]format.FormatFlag{tt: flagP}), gp,
	)
	c = printer.NewPrintContext(format.EmptyTokenMap, p)
	if got := c.GetFormatFlag(tt); got != flagP {
		t.Errorf("Parent lookup failed: got %v, want %v", got, flagP)
	}

	c = printer.NewPrintContext(format.NewTokenMap(
		map[dmlex.TokenType]format.FormatFlag{tt: flagC}), p,
	)
	if got := c.GetFormatFlag(tt); got != flagC {
		t.Errorf("Child lookup failed: got %v, want %v", got, flagC)
	}
}

func TestPrintContext_Parent(t *testing.T) {
	parent := printer.NewPrintContext(
		format.EmptyTokenMap,
		nil,
	)
	child := printer.NewPrintContext(
		format.EmptyTokenMap,
		parent,
	)

	if child.Parent() != parent {
		t.Errorf("Parent() = %p, want %p", child.Parent(), parent)
	}
	if parent.Parent() != nil {
		t.Errorf("Parent() = %p, want nil", parent.Parent())
	}
}

func TestNewPrintContext_NilParent(t *testing.T) {
	ctx := printer.NewPrintContext(format.EmptyTokenMap, nil)
	if ctx == nil {
		t.Fatal("NewPrintContext(nil) returned nil")
	}
	if ctx.Parent() != nil {
		t.Errorf("NewPrintContext(nil).parent = %p, want nil", ctx.Parent())
	}
}
