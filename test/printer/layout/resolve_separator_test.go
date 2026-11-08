// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout_test

import (
	"testing"

	"github.com/timoschenko/dmltools/pkg/printer/format"
	"github.com/timoschenko/dmltools/pkg/printer/layout"
)

func TestResolveSeparator(t *testing.T) {
	tests := []struct {
		name         string
		left         format.FormatFlag
		right        format.FormatFlag
		wantSep      format.Separator
		wantPriority format.Priority
	}{
		{
			name: "HigherLeftPriority",
			left: format.FormatFlag(
				0,
			).WithSeparatorAfterPriority(10).
				WithSeparatorAfter(format.SepSpace).
				WithBreakAfterPriority(5),
			right: format.FormatFlag(
				0,
			).WithSeparatorBeforePriority(5).
				WithSeparatorBefore(format.SepNone).
				WithBreakBeforePriority(10),
			wantSep:      format.SepSpace,
			wantPriority: 5,
		},
		{
			name: "HigherRightPriority",
			left: format.FormatFlag(
				0,
			).WithSeparatorAfterPriority(5).
				WithSeparatorAfter(format.SepSpace).
				WithBreakAfterPriority(5),
			right: format.FormatFlag(
				0,
			).WithSeparatorBeforePriority(10).
				WithSeparatorBefore(format.SepNone).
				WithBreakBeforePriority(10),
			wantSep:      format.SepNone,
			wantPriority: 10,
		},
		{
			name: "EqualPriority_MaxWins_NoneOverSpace",
			left: format.FormatFlag(
				0,
			).WithSeparatorAfterPriority(8).
				WithSeparatorAfter(format.SepNone).
				WithBreakAfterPriority(5),
			right: format.FormatFlag(
				0,
			).WithSeparatorBeforePriority(8).
				WithSeparatorBefore(format.SepSpace).
				WithBreakBeforePriority(10),
			wantSep:      format.SepNone, // SepSpace(0) < SepNone(1)
			wantPriority: 5,
		},
		{
			name: "EqualPriority_MaxWins_MandatoryOverOptional",
			left: format.FormatFlag(
				0,
			).WithSeparatorAfterPriority(8).
				WithSeparatorAfter(format.SepSpaceMaybeWrap).
				WithBreakAfterPriority(5),
			right: format.FormatFlag(
				0,
			).WithSeparatorBeforePriority(8).
				WithSeparatorBefore(format.SepMandatoryWrap).
				WithBreakBeforePriority(10),
			wantSep:      format.SepMandatoryWrap,
			wantPriority: 10,
		},
		{
			name: "Priority0vs15",
			left: format.FormatFlag(
				0,
			).WithSeparatorAfterPriority(0).
				WithSeparatorAfter(format.SepMandatoryWrap).
				WithBreakAfterPriority(15),
			right: format.FormatFlag(
				0,
			).WithSeparatorBeforePriority(15).
				WithSeparatorBefore(format.SepSpace).
				WithBreakBeforePriority(0),
			wantSep:      format.SepSpace,
			wantPriority: 0,
		},
		{
			name: "Equal_Sep_Higher_Break_Priority_Wins",
			left: format.FormatFlag(
				0,
			).WithSeparatorAfterPriority(8).
				WithSeparatorAfter(format.SepSpaceMaybeWrap).
				WithBreakAfterPriority(5),
			right: format.FormatFlag(
				0,
			).WithSeparatorBeforePriority(8).
				WithSeparatorBefore(format.SepSpaceMaybeWrap).
				WithBreakBeforePriority(10),
			wantSep:      format.SepSpaceMaybeWrap,
			wantPriority: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCache := layout.ResolveSeparator(tt.left, tt.right, nil, nil)
			if gotCache.ResolvedSepBefore != tt.wantSep {
				t.Errorf(
					"ResolveSeparator() gotSep = %v, want %v",
					gotCache.ResolvedSepBefore,
					tt.wantSep,
				)
			}
			if gotCache.ResolvedBreakPriority != tt.wantPriority {
				t.Errorf(
					"ResolveSeparator() gotPriority = %v, want %v",
					gotCache.ResolvedBreakPriority,
					tt.wantPriority,
				)
			}
		})
	}
}
