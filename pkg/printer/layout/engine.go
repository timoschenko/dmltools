// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package layout

import (
	"fmt"
	"github.com/timoschenko/dmltools/pkg/dmlex"
	"github.com/timoschenko/dmltools/pkg/printer/format"
	"io"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"
)

// Settings encapsulates layout-specific configuration parameters.
type Settings interface {
	LineCap() int
	TabWidth() int
}

//go:inline
func calculateEffectiveIndent(
	hardIndent, softIndent int, hasHardIndentMarker bool, tabSize int,
) int {
	if hasHardIndentMarker {
		return hardIndent
	}
	// Non-absorb mode for soft ident: return hardIndent + soft
	return max(hardIndent, min(softIndent, hardIndent+tabSize))
}

func simulateBuffer(
	leBuffer []BufferedToken, // readonly
	state State, // pass copy
	stopOnExceed bool,
	tabWidth int,
) (int, int) {
	bestIdx := -1
	bestPriority := -1
	var bestGroup *Group
	for i := range leBuffer {
		bt := &leBuffer[i]
		if bt.IsSynthetic() {
			switch bt.synthetic {
			case SlcSoftIdentBegin:
				state.pushSoftIndent(state.HardIndent, tabWidth)
			case SlcSoftIdentEnd:
				state.popSoftIndent()
			}
			continue
		}

		if bt.Flag.HasIndentDecrease() {
			state.HardIndent = max(0, state.HardIndent-tabWidth)
		}

		switch bt.ResolvedSepBefore {
		case format.SepMandatoryWrap, format.SepWrapButNoSpace:
			if stopOnExceed && state.CursorColumn > state.LineCap &&
				bestIdx != -1 {
				return 0, bestIdx
			}
			state.LastChar = '\n'
			// BUG FIX: Set currentLineLen to 0 instead of simulatedIndent to avoid double-counting
			// the indent, because the block below already adds simulatedIndent when LastChar is '\n'.
			state.CursorColumn = 0
			bestIdx = -1
			bestPriority = -1
			bestGroup = nil
		case format.SepSpace, format.SepSpaceMaybeWrap:
			if state.LastChar != '\n' && state.LastChar != 0 {
				state.CursorColumn++
				state.LastChar = ' '
			}
		}

		if state.LastChar == '\n' {
			isAbsorbMarker := bt.Flag.HasIndentDecrease() ||
				(bt.Flag.HasIndentIncrease() &&
					(bt.ResolvedSepBefore != format.SepMandatoryWrap ||
						(bt.Token != nil && dmlex.GetTokenType(bt.Token) == dmlex.TokenLbrace)))
			if !isAbsorbMarker {
				state.materializeSoftIndent()
			}
			effectiveIndent := calculateEffectiveIndent(
				state.HardIndent,
				state.SoftIndentLevel,
				isAbsorbMarker,
				tabWidth,
			)

			if effectiveIndent > 0 {
				state.CursorColumn += effectiveIndent
				state.LastChar = ' '
			}
		}

		if bt.ResolvedSepBefore == format.SepSpaceMaybeWrap ||
			bt.ResolvedSepBefore == format.SepNoSpaceMaybeWrap {
			priority := int(bt.ResolvedBreakPriority)
			group := bt.AssingedGroup
			if bestIdx == -1 {
				bestIdx = i
				bestPriority = priority
				bestGroup = group
			} else if group == bestGroup {
				if priority >= bestPriority {
					bestIdx = i
					bestPriority = priority
					bestGroup = group
				}
			} else if isAncestor(group, bestGroup) {
				bestIdx = i
				bestPriority = priority
				bestGroup = group
			} else if isAncestor(bestGroup, group) {
				// Outer group takes precedence over inner group
			} else {
				if priority >= bestPriority {
					bestIdx = i
					bestPriority = priority
					bestGroup = group
				}
			}
		}

		if bt.Flag.HasIndentIncrease() {
			state.HardIndent += tabWidth
		}

		if bt.Token != nil {
			if bt.cache == nil {
				panic("cache must be initialized at Push(token) call")
			}

			text := bt.Token.GetText()
			if len(text) > 0 {
				r, _ := utf8.DecodeLastRuneInString(text)
				state.LastChar = r
			}
			if !bt.cache.HasNewlines {
				state.CursorColumn += int(bt.cache.PrefixWidth)
			} else {
				// 1. Prefix ends the current line
				state.CursorColumn += int(bt.cache.PrefixWidth)
				if stopOnExceed && state.CursorColumn > state.LineCap &&
					bestIdx != -1 {
					return 0, bestIdx
				}

				// 2. Middle lines are independent lines (internal to token)
				for _, w := range bt.cache.MiddleWidths {
					state.CursorColumn = int(w)
					if stopOnExceed && state.CursorColumn > state.LineCap &&
						bestIdx != -1 {
						return 0, bestIdx
					}
				}

				// 3. Suffix starts a new line
				state.CursorColumn = int(bt.cache.SuffixWidth)
				bestIdx = -1
				bestPriority = -1
				bestGroup = nil
			}
		}

		// Check if we exceeded linecap after processing the token
		if stopOnExceed && state.CursorColumn > state.LineCap &&
			bestIdx != -1 {
			return 0, bestIdx
		}
	}

	return state.CursorColumn, -1
}

// Engine is responsible for handling the physical layout, spacing, and wrapping
// of the token stream. It acts as a token buffer, separator resolver, and line wrapper.
//
// Architectural Responsibilities:
//   - Buffers tokens received from the Printer.
//   - Resolves conflicts between adjacent tokens using TokenPriority and FormatFlag properties.
//   - Automatically wraps lines when the linecap is exceeded dynamically based on the highest
//     priority SepSpaceMaybeWrap or SepNoSpaceMaybeWrap.
//
// Indentation Constraints & Risks:
//   - Applies conditional indentation (IndentIncrease, IndentDecrease) only when a wrap actually occurs.
//   - WARNING (Indentation Leakage): The Engine maintains internal currentIndent state across
//     token flushes. IndentIncrease must ONLY be used when a corresponding token with IndentDecrease
//     exists to wrap symmetrically. Otherwise, the indentation will leak to following tokens.
type Engine struct {
	settings Settings
	state    *State

	// nil if no parent
	parent *Engine
	tracer TokenTracer

	buffer          []BufferedToken
	writer          io.Writer
	activeGroups    []*Group
	currentGroup    *Group
	triggeredGroups map[*Group]bool
}

// NewEngine creates a new Engine with the given writer and shared context.
func NewEngine(
	settings Settings, writer io.Writer, state *State,
) *Engine {
	if state == nil {
		state = EmptyState()
		state.LineCap = settings.LineCap() // sync lineCap with maxLineCap
	}
	return &Engine{
		settings:        settings,
		state:           state,
		parent:          nil,
		writer:          writer,
		triggeredGroups: make(map[*Group]bool),
	}
}

// HasActiveGroups returns true if any wrap groups are currently active on the group stack.
func (e *Engine) HasActiveGroups() bool {
	return len(e.activeGroups) > 0
}

// TriggerCurrentGroup marks the currently active group and its collective ancestors as triggered.
func (e *Engine) TriggerCurrentGroup() {
	e.triggerGroup(e.currentActiveGroup())
}

// SetTracer attaches an optional token resolution tracer to the Engine.
func (e *Engine) SetTracer(tracer TokenTracer) {
	e.tracer = tracer
}

// Push adds a token to the internal buffer and resolves the separator before it.
// If the token's flag indicates a mandatory wrap after, it triggers a line flush.
// Returns an error if the underlying writer fails during a flush.
func (e *Engine) Push(token antlr.Token, flag format.FormatFlag) error {
	bt := BufferedToken{
		Token: token,
		Flag:  flag,
	}
	if token != nil {
		bt.cache = NewTokenStringCache(token.GetText())
	}

	prevFlag, prevToken, prevGroup, hasPrev := e.getPreviousToken()

	prevTokenType := dmlex.TokenEOF
	if prevToken != nil {
		prevTokenType = dmlex.GetTokenType(prevToken)
	}

	if hasPrev {
		if !e.overrideResolveDesign_Workaround(
			token, &bt, prevToken, prevTokenType,
		) {
			// NOTE: Real formatting design
			bt.ResolveWrappingCache = ResolveSeparator(
				prevFlag, flag, prevGroup, e.currentActiveGroup(),
			)
		}
	} else {
		bt.ResolvedSepBefore = flag.GetSeparatorBefore()
		bt.ResolvedBreakPriority = flag.GetBreakBeforePriority()
		bt.AssingedGroup = e.currentActiveGroup()
		bt.CurrentGroup = bt.AssingedGroup
	}

	isMandatoryWrap := bt.ResolvedSepBefore == format.SepMandatoryWrap ||
		bt.ResolvedSepBefore == format.SepWrapButNoSpace

	if isMandatoryWrap && bt.AssingedGroup != nil &&
		bt.AssingedGroup.WrapMode != WrapModeGreedy {
		targetGroup := bt.AssingedGroup
		// TODO: AI WA: Avoid breaking binary expressions when preceding token is a comment
		if prevTokenType.IsComment() {
			switch targetGroup.Name {
			case "BinaryExpr":
				targetGroup = targetGroup.Parent
			}
		}
		if targetGroup != nil {
			e.triggerGroup(targetGroup)
		}
	}

	if e.isGroupTriggered(bt.AssingedGroup) {
		switch bt.ResolvedSepBefore {
		case format.SepNoSpaceMaybeWrap:
			bt.ResolvedSepBefore = format.SepWrapButNoSpace
			isMandatoryWrap = true
		case format.SepSpaceMaybeWrap:
			bt.ResolvedSepBefore = format.SepMandatoryWrap
			isMandatoryWrap = true
		}
	}

	if hasDebug {
		if e.tracer != nil {
			e.tracePush(
				token, flag, &bt, prevToken, prevFlag, prevGroup, hasPrev,
			)
		}
	}

	e.buffer = append(e.buffer, bt)

	if isMandatoryWrap && len(e.activeGroups) == 0 {
		return e.flushLine()
	}

	return nil
}

// flushLine performs a greedy line-wrapping strategy. It keeps finding the highest
// priority optional break and converting it to a mandatory break until the total
// length of all segments in the buffer fits within the linecap.
func (e *Engine) flushLine() error {
	hasBreak := 0
	for i := 0; ; /* nothing */ i++ {
		if i > 100_000 {
			panic("potential infinite loop in flushLine")
		}

		estLen, bestIdx := simulateBuffer(
			e.buffer,
			*e.state,
			true,
			e.settings.TabWidth(),
		)
		if bestIdx == -1 {
			_ = estLen
			break
		}
		e.applyBreak(bestIdx)
		hasBreak++
	}

	return e.emitBuffer()
}

// BeginSoftIndent manually injects a marker to begin a soft indent group.
func (e *Engine) BeginSoftIndent() func() {
	e.buffer = append(e.buffer, BufferedToken{
		synthetic: SlcSoftIdentBegin,
	})

	return e.EndSoftIndent
}

// EndSoftIndent manually injects a marker to end a soft indent group.
func (e *Engine) EndSoftIndent() {
	e.buffer = append(e.buffer, BufferedToken{
		synthetic: SlcSoftIdentEnd,
	})
}

func (e *Engine) beginGroup(
	name string,
	mode WrapMode,
) (*Group, func()) {
	lg := &Group{
		Name:     name,
		Parent:   e.currentActiveGroup(),
		WrapMode: mode,
	}

	e.currentGroup = lg

	if mode == WrapModeCollective {
		e.buffer = append(e.buffer, BufferedToken{
			synthetic:         SlcWrappingGroupBegin,
			ResolvedSepBefore: format.SepNone,
			AssingedGroup:     nil,
			CurrentGroup:      nil,
		})
	}
	e.activeGroups = append(e.activeGroups, lg)

	return lg, func() {
		e.endGroup(lg)
	}
}

func (e *Engine) endGroup(expected *Group) {
	if len(e.activeGroups) == 0 {
		if expected != nil {
			panic(fmt.Sprintf(
				"endGroup called with empty activeGroups stack (expected %q)",
				expected.Name,
			))
		}
		panic("endGroup called with empty activeGroups stack")
	}

	top := e.activeGroups[len(e.activeGroups)-1]
	if expected != nil && top != expected {
		panic(fmt.Sprintf(
			"group mismatch: popped %q (WrapMode=%s), expected %q (WrapMode=%s)",
			top.Name,
			top.WrapMode,
			expected.Name,
			expected.WrapMode,
		))
	}

	e.activeGroups = e.activeGroups[:len(e.activeGroups)-1]

	if top.WrapMode == WrapModeCollective {
		e.buffer = append(e.buffer, BufferedToken{
			synthetic:         SlcWrappingGroupEnd,
			ResolvedSepBefore: format.SepNone,
			AssingedGroup:     nil,
			CurrentGroup:      nil,
		})
	}

	e.currentGroup = nil
	if sz := len(e.activeGroups); sz > 0 {
		e.currentGroup = e.activeGroups[sz-1]
	}
}

// BeginGroup pushes a new layout group with the specified WrapMode onto the active group stack.
func (e *Engine) BeginGroup(name string, mode WrapMode) func() {
	_, end := e.beginGroup(name, mode)
	return end
}

// BeginWrapGroup pushes a new collective wrap group onto the active group stack.
// Breaks within this group are applied collectively across all peer items when triggered.
func (e *Engine) BeginWrapGroup(name string) func() {
	return e.BeginGroup(name, WrapModeCollective)
}

// EndWrapGroup closes the most recently opened collective wrap group.
// It panics if the active groups stack is empty or if the top group is not a collective wrap group.
func (e *Engine) EndWrapGroup() {
	if len(e.activeGroups) == 0 {
		panic("EndWrapGroup called without matching BeginWrapGroup")
	}
	top := e.activeGroups[len(e.activeGroups)-1]
	if top.WrapMode != WrapModeCollective {
		panic(
			fmt.Sprintf(
				"EndWrapGroup called on non-collective group %q (WrapMode=%s)",
				top.Name,
				top.WrapMode,
			),
		)
	}
	e.endGroup(top)
}

// BeginGreedyGroup pushes a new greedy wrapping group onto the active group stack.
// Unlike collective wrap groups, line breaks within this group are applied independently
// and greedily on the current level without triggering peer items to wrap collectively.
// It also acts as an isolation boundary: upward trigger propagation from child
// wrap groups terminates at this boundary and does not propagate to grandparent groups.
func (e *Engine) BeginGreedyGroup(name string) func() {
	return e.BeginGroup(name, WrapModeGreedy)
}

// EndGreedyGroup closes the most recently opened greedy group.
// It panics if the active groups stack is empty or if the top group is not a greedy group.
func (e *Engine) EndGreedyGroup() {
	if len(e.activeGroups) == 0 {
		panic("EndGreedyGroup called without matching BeginGreedyGroup")
	}
	top := e.activeGroups[len(e.activeGroups)-1]
	if top.WrapMode != WrapModeGreedy {
		panic(
			fmt.Sprintf(
				"EndGreedyGroup called on non-greedy group %q (WrapMode=%s)",
				top.Name,
				top.WrapMode,
			),
		)
	}
	e.endGroup(top)
}

// EndGroup closes the most recently opened group regardless of WrapMode.
func (e *Engine) EndGroup() {
	e.endGroup(nil)
}

// getPreviousToken finds the previous actual token, its format flag, and group,
// checking the current buffer first (skipping synthetic markers) and falling back to the engine state.
func (e *Engine) getPreviousToken() (
	prevFlag format.FormatFlag,
	prevToken antlr.Token,
	prevGroup *Group,
	hasPrev bool,
) {
	if len(e.buffer) > 0 {
		// Find the last actual token in the buffer to resolve the separator against,
		// skipping any group markers.
		for i := range slices.Backward(e.buffer) {
			if e.buffer[i].IsSynthetic() {
				continue
			}

			prevFlag = e.buffer[i].Flag
			prevToken = e.buffer[i].Token
			prevGroup = e.buffer[i].CurrentGroup
			hasPrev = true
			return
		}
	}

	if e.state.HasLastFormatFlag {
		prevFlag = e.state.LastFormatFlag
		prevToken = e.state.LastToken
		prevGroup = e.currentActiveGroup()
		hasPrev = true
		return
	}

	return
}

func (e *Engine) tracePush(
	token antlr.Token,
	flag format.FormatFlag,
	bt *BufferedToken,
	prevToken antlr.Token,
	prevFlag format.FormatFlag,
	prevGroup *Group,
	hasPrev bool,
) {
	var tr TraceRecord
	tr.Left.FormatFlag = prevFlag
	if hasPrev && prevToken != nil {
		tr.Left.Text = prevToken.GetText()
		tr.Left.TokenType = prevToken.GetTokenType()
		tr.Left.Line = prevToken.GetLine()
		tr.Left.Column = prevToken.GetColumn()
		tr.Left.TokenIndex = prevToken.GetTokenIndex()
	}
	if hasPrev && prevGroup != nil {
		tr.Left.GroupPath = prevGroup.Path()
	}

	tr.Right.FormatFlag = flag
	if token != nil {
		tr.Right.Text = token.GetText()
		tr.Right.TokenType = token.GetTokenType()
		tr.Right.Line = token.GetLine()
		tr.Right.Column = token.GetColumn()
		tr.Right.TokenIndex = token.GetTokenIndex()
	}
	if activeGroup := e.currentActiveGroup(); activeGroup != nil {
		tr.Right.GroupPath = activeGroup.Path()
	}

	tr.ResolvedSep = bt.ResolvedSepBefore
	tr.BreakPriority = bt.ResolvedBreakPriority
	tr.Winner = bt.Winner

	e.tracer.Record(&tr)
}

func (e *Engine) currentActiveGroup() *Group {
	return e.currentGroup
}

func (e *Engine) isGroupTriggered(g *Group) bool {
	if g == nil || g.WrapMode == WrapModeGreedy ||
		len(e.triggeredGroups) == 0 {
		return false
	}
	return e.triggeredGroups[g]
}

// triggerGroup marks the given group and its collective wrapping ancestors as triggered.
// Upward trigger propagation stops upon encountering a greedy group boundary (gid.WrapMode == WrapModeGreedy),
// isolating inner wrap explosions from outer ancestor wrap groups.
func (e *Engine) triggerGroup(g *Group) {
	if g == nil {
		return
	}
	if e.triggeredGroups == nil {
		e.triggeredGroups = make(map[*Group]bool)
	}
	for gid := g; gid != nil; gid = gid.Parent {
		if gid.WrapMode == WrapModeGreedy {
			break
		}
		e.triggeredGroups[gid] = true
	}

	for i := range e.buffer {
		bt := &e.buffer[i]
		if bt.IsSynthetic() {
			continue
		}
		if e.isGroupTriggered(bt.AssingedGroup) {
			switch bt.ResolvedSepBefore {
			case format.SepNoSpaceMaybeWrap:
				bt.ResolvedSepBefore = format.SepWrapButNoSpace
			case format.SepSpaceMaybeWrap:
				bt.ResolvedSepBefore = format.SepMandatoryWrap
			}
		}
	}
}

func (e *Engine) applyBreak(idx int) {
	bt := &e.buffer[idx]
	targetGroup := bt.AssingedGroup
	if targetGroup == nil || targetGroup.WrapMode == WrapModeGreedy {
		// Not in a group (or in a greedy group), apply single break
		switch bt.ResolvedSepBefore {
		case format.SepNoSpaceMaybeWrap:
			bt.ResolvedSepBefore = format.SepWrapButNoSpace
		case format.SepSpaceMaybeWrap:
			bt.ResolvedSepBefore = format.SepMandatoryWrap
		}
		return
	}

	e.triggerGroup(targetGroup)

	// Ensure candidate itself is converted
	switch bt.ResolvedSepBefore {
	case format.SepNoSpaceMaybeWrap:
		bt.ResolvedSepBefore = format.SepWrapButNoSpace
	case format.SepSpaceMaybeWrap:
		bt.ResolvedSepBefore = format.SepMandatoryWrap
	}
}

func (e *Engine) emitBuffer() error {
	if len(e.buffer) == 0 {
		return nil
	}

	tabWidth := e.settings.TabWidth()
	for _, bt := range e.buffer {
		if bt.IsSynthetic() {
			switch bt.synthetic {
			case SlcSoftIdentBegin:
				e.state.pushSoftIndent(e.state.HardIndent, tabWidth)
			case SlcSoftIdentEnd:
				e.state.popSoftIndent()
			}
			continue
		}

		// 1. Update Sticky Indentation (applies to this and subsequent lines)
		if bt.Flag.HasIndentDecrease() {
			e.state.HardIndent = max(0, e.state.HardIndent-tabWidth)
		}

		if bt.ResolvedSepBefore == format.SepNone &&
			e.state.LastChar != '\n' &&
			e.state.LastChar != 0 &&
			e.state.LastToken != nil &&
			bt.Token != nil {
			if WillTokensMalform(e.state.LastToken, bt.Token) {
				if hasDebug {
					panic(fmt.Sprintf(
						"tokens would be malformed if printed adjacent with"+
							" SepNone: %q (%s) and %q (%s)",
						e.state.LastToken.GetText(),
						dmlex.GetTokenType(e.state.LastToken),
						bt.Token.GetText(),
						dmlex.GetTokenType(bt.Token),
					))
				} else {
					bt.ResolvedSepBefore = format.SepSpace
				}
			}
		}

		// 2. Handle Separator Before
		switch bt.ResolvedSepBefore {
		case format.SepMandatoryWrap, format.SepWrapButNoSpace:
			if e.state.LastChar != '\n' || bt.Token == nil {
				if _, err := io.WriteString(e.writer, "\n"); err != nil {
					return err
				}
				e.state.LastChar = '\n'
				e.state.CursorColumn = 0
			}
		case format.SepSpace, format.SepSpaceMaybeWrap:
			// SepSpaceMaybeWrap becomes Space if it reached here without being
			// converted to Mandatory
			if e.state.LastChar != '\n' && e.state.LastChar != 0 &&
				e.state.LastChar != ' ' {
				if _, err := io.WriteString(e.writer, " "); err != nil {
					return err
				}
				e.state.LastChar = ' '
				e.state.CursorColumn++
			}
		case format.SepNone, format.SepNoSpaceMaybeWrap:
			// No character emitted
		}

		if e.state.LastChar == '\n' {
			isAbsorbMarker := bt.Flag.HasIndentDecrease() ||
				(bt.Flag.HasIndentIncrease() &&
					(bt.ResolvedSepBefore != format.SepMandatoryWrap ||
						(bt.Token != nil && dmlex.GetTokenType(bt.Token) == dmlex.TokenLbrace)))
			if !isAbsorbMarker {
				e.state.materializeSoftIndent()
			}
			effectiveIndent := calculateEffectiveIndent(
				e.state.HardIndent,
				e.state.SoftIndentLevel,
				isAbsorbMarker,
				e.settings.TabWidth(),
			)

			if effectiveIndent > 0 {
				ws := whitespace(effectiveIndent)
				if _, err := io.WriteString(e.writer, ws); err != nil {
					return err
				}
				e.state.LastChar = ' ' // Not a newline
				e.state.CursorColumn += effectiveIndent
			}
		}

		if bt.Flag.HasIndentIncrease() {
			e.state.HardIndent += tabWidth
		}

		// 3. Emit Token Text
		if bt.Token == nil {
			continue
		}

		text := bt.Token.GetText()
		if _, err := io.WriteString(e.writer, text); err != nil {
			return err
		}

		if ofs := strings.LastIndexByte(text, '\n'); ofs != -1 {
			e.state.CursorColumn = len(text) - ofs - 1
		} else {
			e.state.CursorColumn += len(text)
		}

		if len(text) > 0 {
			r, _ := utf8.DecodeLastRuneInString(text)
			e.state.LastChar = r
		}

		e.state.LastFormatFlag = bt.Flag
		e.state.HasLastFormatFlag = true
		e.state.LastToken = bt.Token
	}

	e.buffer = e.buffer[:0]
	clear(e.triggeredGroups)
	return nil
}

// Flush forces the engine to process and emit all currently buffered tokens.
func (e *Engine) Flush() error {
	if len(e.activeGroups) != 0 {
		panic(
			fmt.Sprintf(
				"Engine.Flush called with unclosed wrap groups: %v",
				e.activeGroups,
			),
		)
	}
	return e.flushLine()
}

// WriteRaw flushes the current buffer and then writes the provided text
// directly to the output writer. Returns an error if either operation fails.
func (e *Engine) WriteRaw(text string) error {
	if err := e.Flush(); err != nil {
		return err
	}
	if len(text) > 0 {
		r, _ := utf8.DecodeLastRuneInString(text)
		e.state.LastChar = r

		if ofs := strings.LastIndexByte(text, '\n'); ofs != -1 {
			e.state.CursorColumn = len(text) - ofs - 1
		} else {
			e.state.CursorColumn += len(text)
		}
		e.state.HasLastFormatFlag = false
		e.state.LastToken = nil
	}
	_, err := io.WriteString(e.writer, text)
	return err
}
