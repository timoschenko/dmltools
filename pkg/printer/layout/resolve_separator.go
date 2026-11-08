// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

import (
	"github.com/timoschenko/dmltools/pkg/printer/format"
)

// ResolveSeparator resolves the effective separator between two adjacent tokens
// using their respective FormatFlags and layout groups.
//
// Conflict resolution proceeds in three deterministic steps:
//  1. Higher Token Priority wins: Compare left token's TokenAfterPriority with right
//     token's TokenBeforePriority. The side with greater priority dictates the separator.
//  2. Equal token priority: Higher separator value wins. The ordering is:
//     SepSpace < SepNone < SepSpaceMaybeWrap < SepNoSpaceMaybeWrap < SepMandatoryWrap < SepWrapButNoSpace.
//  3. Equal token priority and separator: Higher break priority wins. Ties prefer the
//     outer (right) layout group.
func ResolveSeparator(
	left, right format.FormatFlag,
	lg, rg *Group,
) ResolveWrappingCache {
	lp := left.GetTokenAfterPriority()
	rp := right.GetTokenBeforePriority()
	ls := left.GetSeparatorAfter()
	rs := right.GetSeparatorBefore()

	if lp != rp {
		if lp > rp {
			return ResolveWrappingCache{
				ls,
				left.GetBreakAfterPriority(),
				format.WinnerLeft,
				lg,
				rg,
			}
		}
		return ResolveWrappingCache{
			rs,
			right.GetBreakBeforePriority(),
			format.WinnerRight,
			rg,
			rg,
		}
	}

	if ls != rs {
		if ls > rs {
			return ResolveWrappingCache{
				ls,
				left.GetBreakAfterPriority(),
				format.WinnerLeft,
				lg,
				rg,
			}
		}
		return ResolveWrappingCache{
			rs,
			right.GetBreakBeforePriority(),
			format.WinnerRight,
			rg,
			rg,
		}
	}

	// Equal token priority and equal separator: higher break priority wins
	bg := rg // Prefer the outer (right) group on ties
	bp := left.GetBreakAfterPriority()
	winner := format.WinnerEqual
	rpb := right.GetBreakBeforePriority()
	if rpb > bp {
		bg = rg
		bp = rpb
		winner = format.WinnerRight
	} else if bp > rpb {
		bg = lg
		winner = format.WinnerLeft
	}
	return ResolveWrappingCache{ls, bp, winner, bg, rg}
}
