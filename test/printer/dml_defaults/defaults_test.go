// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package dmldefaults_test

import (
	"encoding/json"
	"testing"

	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer"
	dmldefaults "github.com/timoschenko/dmltools/pkg/printer/dml_defaults"
	"github.com/timoschenko/dmltools/pkg/printer/format"
	"github.com/timoschenko/dmltools/pkg/printer/layout"
)

func TestTokenFormatDefaultsParity(t *testing.T) {
	expr := printer.NewPrintContext(
		dmldefaults.Tokens(),
		nil,
	)

	tests := []struct {
		view  string
		left  dmlex.TokenType
		right dmlex.TokenType
		sep   format.Separator
	}{
		{"a b", dmlex.TokenIdentifier, dmlex.TokenIdentifier, format.SepSpace},
		{"a,", dmlex.TokenIdentifier, dmlex.TokenComma, format.SepNone},
		{
			", a",
			dmlex.TokenComma,
			dmlex.TokenIdentifier,
			format.SepSpaceMaybeWrap,
		},
		{"a;", dmlex.TokenIdentifier, dmlex.TokenSemi, format.SepNone},
		{
			"; a",
			dmlex.TokenSemi,
			dmlex.TokenIdentifier,
			format.SepMandatoryWrap,
		},
		// Comments - Comment //
		{"a //c", dmlex.TokenIdentifier, dmlex.TokenComment, format.SepSpace},
		{"; //c", dmlex.TokenSemi, dmlex.TokenComment, format.SepSpace},
		{", //c", dmlex.TokenComma, dmlex.TokenComment, format.SepSpace},
		{
			"//c\\na",
			dmlex.TokenComment,
			dmlex.TokenIdentifier,
			format.SepMandatoryWrap,
		},
		// Comments - Inline /**/
		{
			"a /**/",
			dmlex.TokenIdentifier,
			dmlex.TokenCommentInline,
			format.SepSpace,
		},
		{"; /**/", dmlex.TokenSemi, dmlex.TokenCommentInline, format.SepSpace},
		{
			", /**/",
			dmlex.TokenComma,
			dmlex.TokenCommentInline,
			format.SepSpaceMaybeWrap,
		},
		{
			"/**/ a",
			dmlex.TokenCommentInline,
			dmlex.TokenIdentifier,
			format.SepSpace,
		},
		{"/**/,", dmlex.TokenCommentInline, dmlex.TokenComma, format.SepNone},
		{"/**/;", dmlex.TokenCommentInline, dmlex.TokenSemi, format.SepNone},
	}

	for _, tt := range tests {
		t.Run(tt.left.String()+"_"+tt.right.String(), func(t *testing.T) {
			leftFlag := expr.GetFormatFlag(tt.left)
			rightFlag := expr.GetFormatFlag(tt.right)
			cache := layout.ResolveSeparator(leftFlag, rightFlag, nil, nil)
			sep := cache.ResolvedSepBefore
			if sep != tt.sep {
				t.Errorf(
					"Pair [%s]:\nExpected sep=%v, got %v\n"+
						"Left  % 18v, SA: % 20v, TAP: %v;\n"+
						"Right % 18v, SB: % 20v, TBP: %v;",
					tt.view,
					tt.sep,
					sep,
					tt.left,
					leftFlag.GetSeparatorAfter(),
					leftFlag.GetTokenAfterPriority(),
					tt.right,
					rightFlag.GetSeparatorBefore(),
					rightFlag.GetTokenBeforePriority(),
				)
			}
		})
	}
}

func TestExpressionDefaultsParity(t *testing.T) {
	expr := printer.NewPrintContext(
		dmldefaults.Expression(),
		printer.NewPrintContext(dmldefaults.Tokens(), nil),
	)

	tests := []struct {
		left     dmlex.TokenType
		right    dmlex.TokenType
		hasSpace bool
	}{
		{dmlex.TokenIdentifier, dmlex.TokenIdentifier, true},
		{dmlex.TokenIdentifier, dmlex.TokenDot, false},
		{dmlex.TokenDot, dmlex.TokenIdentifier, false},
		// function call style baseline
		{dmlex.TokenIdentifier, dmlex.TokenLparen, false},
		{dmlex.TokenRparen, dmlex.TokenIdentifier, true},
		{dmlex.TokenIdentifier, dmlex.TokenComma, false},
		{dmlex.TokenComma, dmlex.TokenIdentifier, true},
		{dmlex.TokenIdentifier, dmlex.TokenSemi, false},
		// mandatory wrap usually resolved as space if on one line
		{dmlex.TokenSemi, dmlex.TokenIdentifier, true},
	}

	for _, tt := range tests {
		t.Run(tt.left.String()+"_"+tt.right.String(), func(t *testing.T) {
			leftFlag := expr.GetFormatFlag(tt.left)
			rightFlag := expr.GetFormatFlag(tt.right)
			cache := layout.ResolveSeparator(leftFlag, rightFlag, nil, nil)
			sep := cache.ResolvedSepBefore
			// For this test, we consider anything not SepNone as "having space or break"
			actualHasSpace := sep != format.SepNone &&
				sep != format.SepWrapButNoSpace &&
				sep != format.SepNoSpaceMaybeWrap
			if actualHasSpace != tt.hasSpace {
				t.Errorf(
					"Pair (%v, %v): expected hasSpace=%v, got %v (sep=%s)",
					tt.left, tt.right, tt.hasSpace, actualHasSpace, sep,
				)
			}
		})
	}
}

func TestCdeclDefaultsParity(t *testing.T) {
	cdecl := printer.NewPrintContext(
		dmldefaults.Cdecl(),
		printer.NewPrintContext(dmldefaults.Tokens(), nil),
	)

	tests := []struct {
		left     dmlex.TokenType
		right    dmlex.TokenType
		hasSpace bool
	}{
		{dmlex.TokenIdentifier, dmlex.TokenMul, true},     // int *
		{dmlex.TokenMul, dmlex.TokenIdentifier, false},    // *p
		{dmlex.TokenMul, dmlex.TokenMul, false},           // **
		{dmlex.TokenIdentifier, dmlex.TokenLparen, false}, // func(
		{dmlex.TokenLparen, dmlex.TokenMul, false},        // (*
		{dmlex.TokenLparen, dmlex.TokenIdentifier, false}, // (p
		{dmlex.TokenIdentifier, dmlex.TokenRparen, false}, // p)
		{dmlex.TokenIdentifier, dmlex.TokenLbrack, false}, // p[
		{dmlex.TokenLbrack, dmlex.TokenIdentifier, false}, // [5
		{dmlex.TokenIdentifier, dmlex.TokenRbrack, false}, // 5]
		{dmlex.TokenBitand, dmlex.TokenIdentifier, false}, // &p
		{dmlex.TokenIdentifier, dmlex.TokenDot, false},    // p.
		{dmlex.TokenIdentifier, dmlex.TokenSemi, false},   // p;
		{dmlex.TokenIdentifier, dmlex.TokenComma, false},  // p,
		{dmlex.TokenIdentifier, dmlex.TokenColon, false},  // p:
		{dmlex.TokenColon, dmlex.TokenIdentifier, true},   // : 5
	}

	for _, tt := range tests {
		t.Run(
			"Cdecl_"+tt.left.String()+"_"+tt.right.String(),
			func(t *testing.T) {
				leftFlag := cdecl.GetFormatFlag(tt.left)
				rightFlag := cdecl.GetFormatFlag(tt.right)
				cache := layout.ResolveSeparator(leftFlag, rightFlag, nil, nil)
				sep := cache.ResolvedSepBefore
				actualHasSpace := sep != format.SepNone &&
					sep != format.SepWrapButNoSpace &&
					sep != format.SepNoSpaceMaybeWrap
				if actualHasSpace != tt.hasSpace {
					t.Errorf(
						"Pair (%v, %v): expected hasSpace=%v, got %v (sep=%s)",
						tt.left, tt.right, tt.hasSpace, actualHasSpace, sep,
					)
				}
			},
		)
	}
}

func TestHashTokenFlagParity(t *testing.T) {
	defaults := dmldefaults.Tokens()

	pairs := []struct {
		analog    dmlex.TokenType
		hashToken dmlex.TokenType
	}{
		{dmlex.TokenIf, dmlex.TokenHashIf},
		{dmlex.TokenElse, dmlex.TokenHashElse},
		{dmlex.TokenForeach, dmlex.TokenHashForeach},
		{dmlex.TokenSelect, dmlex.TokenHashSelect},
		{dmlex.TokenQuestion, dmlex.TokenAltQuestion},
		{dmlex.TokenColon, dmlex.TokenAltColon},
	}

	for _, pair := range pairs {
		t.Run(
			pair.analog.String()+"_vs_"+pair.hashToken.String(),
			func(t *testing.T) {
				analogFlag := defaults.GetFormatFlag(pair.analog)
				hashFlag := defaults.GetFormatFlag(pair.hashToken)
				if analogFlag != hashFlag {
					t.Errorf(
						"Token flag mismatch between %v and %v:\nanalog:    %v\nhashToken: %v",
						pair.analog,
						pair.hashToken,
						analogFlag,
						hashFlag,
					)
				}
			},
		)
	}
}

func TestDefaultTablesMarshal(t *testing.T) {
	tables := []struct {
		name string
		tfm  format.TokenMap
	}{
		{"Tokens", dmldefaults.Tokens()},
		{"Expression", dmldefaults.Expression()},
		{"Cdecl", dmldefaults.Cdecl()},
	}

	for _, tc := range tables {
		data, err := json.Marshal(tc.tfm)
		if err != nil {
			t.Fatalf("Marshal(%s) error: %v", tc.name, err)
		}

		decoded, err := format.UnmarshalTokenMap(data)
		if err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", tc.name, err)
		}

		// Validate some known tokens
		for _, tt := range []dmlex.TokenType{
			dmlex.TokenSemi,
			dmlex.TokenComma,
			dmlex.TokenIdentifier,
		} {
			if got, want := decoded.GetFormatFlag(
				tt,
			), tc.tfm.GetFormatFlag(
				tt,
			); got != want {
				t.Errorf("%s[%v] = %v, want %v", tc.name, tt, got, want)
			}
		}
	}
}
