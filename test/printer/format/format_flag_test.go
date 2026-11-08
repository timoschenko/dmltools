// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format_test

import (
	"testing"

	"github.com/timoschenko/dmltools/pkg/printer/format"
)

func TestFormatFlag_RoundTrip(t *testing.T) {
	t.Run("BreakBeforePriority", func(t *testing.T) {
		for i := range format.Priority(16) {
			f := format.FormatFlag(0).WithBreakBeforePriority(i)
			if got := f.GetBreakBeforePriority(); got != i {
				t.Errorf("GetBreakBeforePriority() = %v, want %v", got, i)
			}
			// Verify isolation
			if f&^format.BreakBeforePriorityMask != 0 {
				t.Errorf("Field corruption: bits outside mask set: %032b", f)
			}
		}
	})

	t.Run("SeparatorBefore", func(t *testing.T) {
		seps := []format.Separator{
			format.SepSpace,
			format.SepNone,
			format.SepSpaceMaybeWrap,
			format.SepMandatoryWrap,
		}
		for _, s := range seps {
			f := format.FormatFlag(0).WithSeparatorBefore(s)
			if got := f.GetSeparatorBefore(); got != s {
				t.Errorf("GetSeparatorBefore() = %v, want %v", got, s)
			}
			// Verify isolation
			if f&^format.SepBeforeMask != 0 {
				t.Errorf("Field corruption: bits outside mask set: %032b", f)
			}
		}
	})

	t.Run("BreakAfterPriority", func(t *testing.T) {
		for i := range format.Priority(16) {
			f := format.FormatFlag(0).WithBreakAfterPriority(i)
			if got := f.GetBreakAfterPriority(); got != i {
				t.Errorf("GetBreakAfterPriority() = %v, want %v", got, i)
			}
			// Verify isolation
			if f&^format.BreakAfterPriorityMask != 0 {
				t.Errorf("Field corruption: bits outside mask set: %032b", f)
			}
		}
	})

	t.Run("SeparatorAfter", func(t *testing.T) {
		seps := []format.Separator{
			format.SepSpace,
			format.SepNone,
			format.SepSpaceMaybeWrap,
			format.SepMandatoryWrap,
		}
		for _, s := range seps {
			f := format.FormatFlag(0).WithSeparatorAfter(s)
			if got := f.GetSeparatorAfter(); got != s {
				t.Errorf("GetSeparatorAfter() = %v, want %v", got, s)
			}
			// Verify isolation
			if f&^format.SepAfterMask != 0 {
				t.Errorf("Field corruption: bits outside mask set: %032b", f)
			}
		}
	})

	t.Run("IndentIncrease", func(t *testing.T) {
		f := format.FormatFlag(0).WithIndentIncrease()
		if !f.HasIndentIncrease() {
			t.Error("HasIndentIncrease() = false, want true")
		}
		if f&^format.IndentIncMask != 0 {
			t.Errorf("Field corruption: bits outside mask set: %032b", f)
		}
		f = f.WithoutIndentIncrease()
		if f.HasIndentIncrease() {
			t.Error("WithoutIndentIncrease() still has IndentIncrease")
		}
	})

	t.Run("IndentDecrease", func(t *testing.T) {
		f := format.FormatFlag(0).WithIndentDecrease()
		if !f.HasIndentDecrease() {
			t.Error("HasIndentDecrease() = false, want true")
		}
		if f&^format.IndentDecMask != 0 {
			t.Errorf("Field corruption: bits outside mask set: %032b", f)
		}
		f = f.WithoutIndentDecrease()
		if f.HasIndentDecrease() {
			t.Error("WithoutIndentDecrease() still has IndentDecrease")
		}
	})

	t.Run("ForceWrap", func(t *testing.T) {
		f := format.FormatFlag(0).WithForceWrap()
		if !f.IsForceWrap() {
			t.Error("IsForceWrap() = false, want true")
		}
		if f&^format.ForceWrapMask != 0 {
			t.Errorf("Field corruption: bits outside mask set: %032b", f)
		}
	})

	t.Run("TokenBeforePriority", func(t *testing.T) {
		for i := range format.Priority(16) {
			f := format.FormatFlag(0).WithSeparatorBeforePriority(i)
			if got := f.GetTokenBeforePriority(); got != i {
				t.Errorf("GetTokenBeforePriority() = %v, want %v", got, i)
			}
			// Verify isolation
			if f&^format.TokenBeforePriorityMask != 0 {
				t.Errorf("Field corruption: bits outside mask set: %032b", f)
			}
		}
	})

	t.Run("TokenAfterPriority", func(t *testing.T) {
		for i := range format.Priority(16) {
			f := format.FormatFlag(0).WithSeparatorAfterPriority(i)
			if got := f.GetTokenAfterPriority(); got != i {
				t.Errorf("GetTokenAfterPriority() = %v, want %v", got, i)
			}
			// Verify isolation
			if f&^format.TokenAfterPriorityMask != 0 {
				t.Errorf("Field corruption: bits outside mask set: %032b", f)
			}
		}
	})

	t.Run("TokenPriorityShim", func(t *testing.T) {
		for i := range format.Priority(16) {
			f := format.FormatFlag(0).WithSeparatorPriority(i)
			if got := f.GetTokenBeforePriority(); got != i {
				t.Errorf("GetTokenBeforePriority() = %v, want %v", got, i)
			}
			if got := f.GetTokenAfterPriority(); got != i {
				t.Errorf("GetTokenAfterPriority() = %v, want %v", got, i)
			}
			if got := f.GetTokenPriority(); got != i {
				t.Errorf("GetTokenPriority() = %v, want %v", got, i)
			}
		}
	})

	t.Run("PriorityOverflow", func(t *testing.T) {
		// Value 16 (0x10) should be masked to 0
		f := format.FormatFlag(0).WithSeparatorBeforePriority(16)
		if got := f.GetTokenBeforePriority(); got != 0 {
			t.Errorf("WithTokenBeforePriority(16) got %v, want 0", got)
		}

		f = format.FormatFlag(0).WithSeparatorAfterPriority(16)
		if got := f.GetTokenAfterPriority(); got != 0 {
			t.Errorf("WithTokenAfterPriority(16) got %v, want 0", got)
		}

		// Value 17 (0x11) should be masked to 1
		f = format.FormatFlag(0).WithSeparatorBeforePriority(17)
		if got := f.GetTokenBeforePriority(); got != 1 {
			t.Errorf("WithTokenBeforePriority(17) got %v, want 1", got)
		}

		f = format.FormatFlag(0).WithSeparatorAfterPriority(17)
		if got := f.GetTokenAfterPriority(); got != 1 {
			t.Errorf("WithTokenAfterPriority(17) got %v, want 1", got)
		}
	})
}

func TestFormatFlag_AllFieldsCombined(t *testing.T) {
	f := format.FormatFlag(0).
		WithBreakBeforePriority(5).
		WithSeparatorBefore(format.SepNone).
		WithBreakAfterPriority(10).
		WithSeparatorAfter(format.SepMandatoryWrap).
		WithIndentIncrease().
		WithIndentDecrease().
		WithForceWrap().
		WithSeparatorBeforePriority(14).
		WithSeparatorAfterPriority(15)

	if f.GetBreakBeforePriority() != 5 {
		t.Errorf(
			"GetBreakBeforePriority() = %v, want 5",
			f.GetBreakBeforePriority(),
		)
	}
	if f.GetSeparatorBefore() != format.SepNone {
		t.Errorf(
			"GetSeparatorBefore() = %v, want %v",
			f.GetSeparatorBefore(),
			format.SepNone,
		)
	}
	if f.GetBreakAfterPriority() != 10 {
		t.Errorf(
			"GetBreakAfterPriority() = %v, want 10",
			f.GetBreakAfterPriority(),
		)
	}
	if f.GetSeparatorAfter() != format.SepMandatoryWrap {
		t.Errorf(
			"GetSeparatorAfter() = %v, want %v",
			f.GetSeparatorAfter(),
			format.SepMandatoryWrap,
		)
	}
	if !f.HasIndentIncrease() {
		t.Error("HasIndentIncrease() = false, want true")
	}
	if !f.HasIndentDecrease() {
		t.Error("HasIndentDecrease() = false, want true")
	}
	if !f.IsForceWrap() {
		t.Error("IsForceWrap() = false, want true")
	}
	if f.GetTokenBeforePriority() != 14 {
		t.Errorf(
			"GetTokenBeforePriority() = %v, want 14",
			f.GetTokenBeforePriority(),
		)
	}
	if f.GetTokenAfterPriority() != 15 {
		t.Errorf(
			"GetTokenAfterPriority() = %v, want 15",
			f.GetTokenAfterPriority(),
		)
	}
}

func TestFormatFlag_String(t *testing.T) {
	f := format.FormatFlag(0).
		WithBreakBeforePriority(5).
		WithSeparatorBefore(format.SepSpace).
		WithBreakAfterPriority(3).
		WithSeparatorAfter(format.SepSpaceMaybeWrap).
		WithIndentIncrease().
		WithSeparatorPriority(8)

	want := "FormatFlag{BBP:5 SB:SepSpace BAP:3 SA:SepSpaceMaybeWrap IndentInc:true IndentDec:false ForceWrap:false TBP:8 TAP:8}"
	if got := f.String(); got != want {
		t.Errorf("FormatFlag.String() = %v, want %v", got, want)
	}
}

func TestFormatFlag_AllBitsSet(t *testing.T) {
	f := format.FormatFlag(0xFFFFFFFF)
	// Getters should mask and return max values without panic
	if f.GetBreakBeforePriority() != 15 {
		t.Errorf(
			"GetBreakBeforePriority() = %v, want 15",
			f.GetBreakBeforePriority(),
		)
	}
	if f.GetSeparatorBefore() != format.Separator(7) {
		t.Errorf("GetSeparatorBefore() = %v, want 7", f.GetSeparatorBefore())
	}
	if f.GetBreakAfterPriority() != 15 {
		t.Errorf(
			"GetBreakAfterPriority() = %v, want 15",
			f.GetBreakAfterPriority(),
		)
	}
	if f.GetSeparatorAfter() != format.Separator(7) {
		t.Errorf("GetSeparatorAfter() = %v, want 7", f.GetSeparatorAfter())
	}
	if !f.HasIndentIncrease() {
		t.Error("HasIndentIncrease() = false, want true")
	}
	if f.GetTokenBeforePriority() != 15 {
		t.Errorf(
			"GetTokenBeforePriority() = %v, want 15",
			f.GetTokenBeforePriority(),
		)
	}
	if f.GetTokenAfterPriority() != 15 {
		t.Errorf(
			"GetTokenAfterPriority() = %v, want 15",
			f.GetTokenAfterPriority(),
		)
	}
	if f.GetTokenPriority() != 15 {
		t.Errorf("GetTokenPriority() = %v, want 15", f.GetTokenPriority())
	}
}
