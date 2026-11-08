// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout_test

import (
	"github.com/timoschenko/dmltools/pkg/printer"
	"github.com/timoschenko/dmltools/pkg/printer/format"
	"github.com/timoschenko/dmltools/pkg/printer/layout"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

// PoC test: Reproduces the bug where the closing bracket `}` doesn't get a
// newline before it when the inner wrap group triggers wrapping.
//
// The structure mirrors printCdeclInitializer + printSingleInitializerArray:
//   - Outer group "CdeclInitializer" contains: { [inner group] }
//   - Inner group "SingleInitializer" contains: elem1, elem2, elem3
//
// When the inner group wraps (commas become mandatory breaks), the `}`
// which is in the OUTER group should also get its optional wrap converted
// to mandatory — but currently it doesn't because:
//
//  1. applyBreak only processes tokens within the immediate group boundaries
//     of the break candidate (the inner group), so `}` (in the outer group)
//     is never considered.
//
//  2. After all inner breaks are applied, the line `    c}` fits within
//     linecap, so findBreakCandidate returns -1 and the loop exits without
//     ever converting `}`'s SepNoSpaceMaybeWrap to mandatory.
func TestLayoutEngine_WrapGroups_ClosingBracketBug(t *testing.T) {
	t.Run("inner_wrap_should_propagate_to_outer_bracket", func(t *testing.T) {
		var sb strings.Builder
		// linecap=20: "prefix {a_long_element_name, b, c}" exceeds 20
		le := layout.NewEngine(
			printer.DefaultSettings(20),
			&sb,
			nil,
		)

		push := func(token antlr.Token, flag format.FormatFlag) {
			if err := le.Push(token, flag); err != nil {
				t.Fatal(err)
			}
		}

		// "prefix " — simulates "saved tp x = "
		push(&mockToken{text: "prefix"}, format.FormatFlag(0).
			WithSeparatorAfter(format.SepSpace))

		// Outer group: CdeclInitializer
		endOuter := le.BeginWrapGroup("CdeclInitializer")

		// `{` token — same flags as printCdeclInitializer
		push(&mockToken{text: "{"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityLowest).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithIndentIncrease().
			WithBreakAfterPriority(format.PriorityVeryHigh).
			WithSeparatorAfterPriority(format.PriorityMax))

		// Inner group: SingleInitializer
		endInner := le.BeginWrapGroup("SingleInitializer")

		// elem1 (long element that forces wrapping)
		// SepBefore resolved between `{` (SepAfter=OptionalWrapNoSpace, prio Max)
		// and elem1 default. Since `{` has PriorityMax, `{`'s SepAfter wins.
		push(&mockToken{text: "a_long_element_name"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityVeryHigh))

		// comma after elem1
		push(&mockToken{text: ","}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithBreakAfterPriority(format.PriorityHigh).
			WithSeparatorPriority(format.PriorityMax))

		// elem2
		push(&mockToken{text: "b"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityHigh))

		// comma after elem2
		push(&mockToken{text: ","}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithBreakAfterPriority(format.PriorityHigh).
			WithSeparatorPriority(format.PriorityMax))

		// elem3
		push(&mockToken{text: "c"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityHigh))

		endInner()

		// `}` token — same flags as printCdeclInitializer
		// This is OUTSIDE the inner group but INSIDE the outer group.
		// SepBefore is resolved between elem3 and `}` flags.
		// `}` has SepBeforePriority=PriorityMax so its SepBefore wins.
		push(&mockToken{text: "}"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityMax).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityLowest).
			WithIndentDecrease().
			WithBreakBeforePriority(format.PriorityMax))

		endOuter()

		// `;` — glued to `}` with no separator
		push(&mockToken{text: ";"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityMax))

		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		// BUG: Currently produces "prefix {\n    a_long_element_name,\n    b,\n    c};"
		// The `}` should be on its own line:
		expected := "prefix {\n    a_long_element_name,\n    b,\n    c\n};"
		got := sb.String()

		t.Logf("GOT:      %q", got)
		t.Logf("EXPECTED: %q", expected)

		if got != expected {
			t.Errorf(
				"closing bracket not wrapped correctly.\n  expected: %q\n       got: %q",
				expected,
				got,
			)
		}
	})

	t.Run("no_wrap_when_fits", func(t *testing.T) {
		var sb strings.Builder
		// linecap=80: everything fits on one line
		le := layout.NewEngine(
			printer.DefaultSettings(80),
			&sb,
			nil,
		)

		push := func(token antlr.Token, flag format.FormatFlag) {
			if err := le.Push(token, flag); err != nil {
				t.Fatal(err)
			}
		}

		push(&mockToken{text: "prefix"}, format.FormatFlag(0).
			WithSeparatorAfter(format.SepSpace))

		endOuter := le.BeginWrapGroup("CdeclInitializer")

		push(&mockToken{text: "{"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpace).
			WithSeparatorBeforePriority(format.PriorityLowest).
			WithSeparatorAfter(format.SepNoSpaceMaybeWrap).
			WithIndentIncrease().
			WithBreakAfterPriority(format.PriorityVeryHigh).
			WithSeparatorAfterPriority(format.PriorityMax))

		endInner := le.BeginWrapGroup("SingleInitializer")

		push(&mockToken{text: "a"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityVeryHigh))

		push(&mockToken{text: ","}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorAfter(format.SepSpaceMaybeWrap).
			WithBreakAfterPriority(format.PriorityHigh).
			WithSeparatorPriority(format.PriorityMax))

		push(&mockToken{text: "b"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepSpaceMaybeWrap).
			WithBreakBeforePriority(format.PriorityHigh))

		endInner()

		push(&mockToken{text: "}"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNoSpaceMaybeWrap).
			WithSeparatorBeforePriority(format.PriorityMax).
			WithSeparatorAfter(format.SepNone).
			WithSeparatorAfterPriority(format.PriorityLowest).
			WithIndentDecrease().
			WithBreakBeforePriority(format.PriorityMax))

		endOuter()

		push(&mockToken{text: ";"}, format.FormatFlag(0).
			WithSeparatorBefore(format.SepNone).
			WithSeparatorBeforePriority(format.PriorityMax))

		if err := le.Flush(); err != nil {
			t.Fatal(err)
		}

		// Should NOT wrap when everything fits
		expected := "prefix {a, b};"
		got := sb.String()

		if got != expected {
			t.Errorf(
				"should not wrap when fits.\n  expected: %q\n       got: %q",
				expected,
				got,
			)
		}
	})
}
