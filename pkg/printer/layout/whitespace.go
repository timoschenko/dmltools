// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

import "strings"

const Whitespace = "                                                                                                                                 " // 129

// whitespace generates a string containing spaces corresponding to the requested amount.
// It uses a pre-allocated cache up to the length of the internal Whitespace constant for speed.
func whitespace(amount int) string {
	if amount >= len(Whitespace) {
		return strings.Repeat(" ", amount)
	}

	return Whitespace[:amount]
}
