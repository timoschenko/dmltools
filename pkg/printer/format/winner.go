// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package format

// Winner indicates which side won in token pair conflict resolution.
type Winner uint8

const (
	WinnerNone Winner = iota
	WinnerLeft
	WinnerRight
	WinnerEqual
)
