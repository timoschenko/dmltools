// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

// Group represents a logical formatting group with a name, parent, and wrap mode.
type Group struct {
	Name     string
	Parent   *Group
	WrapMode WrapMode
}

func (g *Group) Path() string {
	if g == nil {
		return ""
	}
	if g.Parent == nil {
		return g.Name
	}
	return g.Parent.Path() + "/" + g.Name
}

//go:inline
func isAncestor(ancestor, desc *Group) bool {
	if ancestor == desc || desc == nil {
		return false
	}
	if ancestor == nil {
		return true
	}
	for cur := desc.Parent; cur != nil; cur = cur.Parent {
		if cur == ancestor {
			return true
		}
	}
	return false
}
