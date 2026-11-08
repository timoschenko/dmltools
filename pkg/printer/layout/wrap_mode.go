// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

import "fmt"

// WrapMode defines the wrapping behavior for a layout group.
type WrapMode uint8

const (
	// WrapModeCollective wraps all peer break candidates collectively when the group is triggered.
	WrapModeCollective WrapMode = iota
	// WrapModeGreedy wraps break candidates independently and greedily as line length overflows,
	// without triggering peer candidates, and acts as an upward trigger propagation isolation boundary.
	WrapModeGreedy
)

func (m WrapMode) String() string {
	switch m {
	case WrapModeCollective:
		return "Collective"
	case WrapModeGreedy:
		return "Greedy"
	default:
		return fmt.Sprintf("WrapMode(%d)", m)
	}
}
