// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

import (
	"github.com/timoschenko/dmltools/pkg/printer/format"

	"github.com/antlr4-go/antlr/v4"
)

type softIndentEntry struct {
	baseHardIndent int
	targetIndent   int
	appliedDelta   int
}

// State represents the physical layout state during token rendering.
type State struct {
	// Deprecated: Use Settings.LineCap
	LineCap int

	CursorColumn      int
	HardIndent        int
	SoftIndentLevel   int
	LastChar          rune
	LastFormatFlag    format.FormatFlag
	HasLastFormatFlag bool
	LastToken         antlr.Token

	softIndentStack [8]softIndentEntry
	softIndentDepth int
}

func (s *State) pushSoftIndent(hardIndent, tabWidth int) {
	target := hardIndent + tabWidth
	if s.softIndentDepth < len(s.softIndentStack) {
		s.softIndentStack[s.softIndentDepth] = softIndentEntry{
			baseHardIndent: hardIndent,
			targetIndent:   target,
			appliedDelta:   0,
		}
		s.softIndentDepth++
	}
	s.SoftIndentLevel = target
}

func (s *State) popSoftIndent() {
	if s.softIndentDepth > 0 {
		s.softIndentDepth--
		entry := s.softIndentStack[s.softIndentDepth]
		if entry.appliedDelta > 0 {
			s.HardIndent = max(0, s.HardIndent-entry.appliedDelta)
		}
	}
	if s.softIndentDepth > 0 {
		s.SoftIndentLevel = s.softIndentStack[s.softIndentDepth-1].targetIndent
	} else {
		s.SoftIndentLevel = 0
	}
}

func (s *State) materializeSoftIndent() {
	if s.softIndentDepth == 0 {
		return
	}
	top := &s.softIndentStack[s.softIndentDepth-1]
	if top.appliedDelta == 0 && s.HardIndent == top.baseHardIndent {
		delta := top.targetIndent - s.HardIndent
		if delta > 0 {
			s.HardIndent += delta
			top.appliedDelta = delta
		}
	}
}

// EmptyState returns an initialized empty layout state.
func EmptyState() *State {
	return &State{
		LineCap:           0,
		CursorColumn:      0,
		HardIndent:        0,
		SoftIndentLevel:   0,
		LastChar:          0,
		HasLastFormatFlag: false,
		LastToken:         nil,
	}
}
