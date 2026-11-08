// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

// Package printer provides the code formatting engine for Device Modeling Language (DML).
//
// The formatter operates on a Token Stream Processing model using Unified
// Formatting Flags (UFF). Formatting instructions are decoupled from AST traversal:
// AST printer methods visit ANTLR parse tree nodes and push scoped [PrintContext]
// overrides containing [format.FormatFlag] definitions, while a unified [layout.Engine]
// evaluates token buffering, line capacity, wrapping priorities, and indentation.
//
// # Architecture Overview
//
//	     ANTLR4 AST Walk (pkg/printer/dml_*.go)
//	                       │
//	                       ▼
//	      Assign PrintContext & FormatFlags
//	                       │
//	                       ▼
//	           pushTokenRange / processToken
//	                       │
//	                       ▼
//	         LayoutEngine (pkg/printer/layout)
//	┌──────────────────────────────────────────────┐
//	│ - Token Buffer ([]layout.BufferedToken)      │
//	│ - ResolveSeparator (TokenPriority / Layout)  │
//	│ - simulateBuffer (Linecap & Wrap Points)     │
//	│ - Indentation Tracking (Hard & Soft)         │
//	│ - Wrapping Groups (*layout.Group)            │
//	└──────────────────────────────────────────────┘
//	                       │
//	                       ▼
//	              Formatted Output Stream
//
// # 1. Unified Formatting Flags (UFF)
//
// To minimize memory overhead and streamline formatting decisions, spacing and
// line-break rules are packed into a 32-bit bitfield ([format.FormatFlag]).
//
// # Bitfield Layout
//
//	Bits  0-3:   Line break before priority (0-15)
//	Bits  4-6:   Separator Before ([format.Separator])
//	Bit   7:     Reserved
//	Bits  8-11:  Line break after priority (0-15)
//	Bits 12-14:  Separator After ([format.Separator])
//	Bit  15:     Reserved
//	Bit  16:     Hard Indent increase (applies only if a line break occurs here)
//	Bit  17:     Hard Indent decrease (applies only if a line break occurs here)
//	Bit  18:     ForceWrap (indicates indentation represents wrap context, not block scope)
//	Bits 19-23:  Reserved
//	Bits 24-27:  Token Before Priority (0-15) — used for pair conflict resolution
//	Bits 28-31:  Token After Priority (0-15) — used for pair conflict resolution
//
// # Separators
//
// Spacing, optional line breaks, and mandatory line breaks are represented by the
// [format.Separator] enumeration as mutually exclusive states:
//
//   - [format.SepSpace]: Single space between tokens.
//   - [format.SepNone]: No space (clinging adjacent tokens, e.g. parentheses or dots).
//   - [format.SepSpaceMaybeWrap]: Space if single-line; wraps if line capacity is exceeded.
//   - [format.SepNoSpaceMaybeWrap]: No space if single-line; wraps if line capacity is exceeded.
//   - [format.SepMandatoryWrap]: Forced newline between tokens.
//   - [format.SepWrapButNoSpace]: Forced newline without leading/trailing space.
//
// # Priorities (4-bit Scale: 0–15)
//
// Priorities ([format.Priority]) govern both break selection and token conflict resolution:
//
//   - Priority 14–15 ([format.PriorityCritical] / [format.PriorityMax]): Structural delimiters and mandatory boundaries.
//   - Priority 10–12 ([format.PriorityHigh] / [format.PriorityHighest]): Comma separators, outer list elements, statement split points.
//   - Priority 8–9 ([format.PriorityAboveNormal] / [format.PriorityMediumHigh]): Major binary operators (assignment, logical `&&`, `||`).
//   - Priority 6–7 ([format.PriorityNormal] / [format.PriorityDefault]): Comparison and relational operators (`==`, `!=`, `<`, `>`).
//   - Priority 4–5 ([format.PriorityBelowNormal] / [format.PriorityMediumLow]): Additive arithmetic operators (`+`, `-`).
//   - Priority 2–3 ([format.PriorityLow] / [format.PriorityVeryLow]): Multiplicative operators (`*`, `/`, `%`).
//   - Priority 1 ([format.PriorityLowest]): Clinging operators and member access (`.`, `->`).
//
// # Pair Conflict Resolution
//
// When two adjacent tokens emit conflicting separator requirements (such as the left
// token's SeparatorAfter versus the right token's SeparatorBefore), [layout.ResolveSeparator]
// deterministically decides the outcome using three tie-breaking rules:
//
//  1. Higher Token Priority wins: Compare left token's TokenAfterPriority with right
//     token's TokenBeforePriority. The side with greater priority dictates the separator.
//  2. Equal token priority: Higher separator value wins. Ordering:
//     [format.SepSpace] < [format.SepNone] < [format.SepSpaceMaybeWrap] < [format.SepNoSpaceMaybeWrap] < [format.SepMandatoryWrap] < [format.SepWrapButNoSpace].
//  3. Equal token priority and separator: Higher break priority wins. Ties prefer the
//     outer (right) layout group.
//
// # Comments as High-Priority Tokens
//
// Comments are treated as first-class tokens carrying high token priority
// ([format.PriorityMax] or [format.PriorityAboveHighest]), ensuring comments win conflict resolution
// against adjacent code tokens and preserve correct spacing and newlines.
//
// # Indentation & Wrap Context
//
//   - Hard Indentation (Bits 16 & 17): [format.FormatFlag.WithIndentIncrease] and
//     [format.FormatFlag.WithIndentDecrease] adjust indentation in [TabSize] (4 spaces) increments.
//     Crucially, indentation shifts trigger ONLY if a line break actually occurs at that token.
//   - ForceWrap (Bit 18): [format.FormatFlag.WithForceWrap] indicates that indentation describes
//     a multi-line wrapping context (such as an expression continued on the next line)
//     rather than a structural block scope.
//   - Soft Indentation: The layout engine additionally supports soft indentation markers
//     for statements with hanging indentation requirements.
//
// # 2. Scoped Contexts via Chained Maps
//
// Formatting configurations are scoped using [PrintContext]. Contexts form a parent
// hierarchy: looking up a token flag searches the local overrides map and recurses
// up to the parent chain, falling back to [DefaultFlagFor] if unmapped.
//
// Predefined baseline tables implement [format.TokenMap] in subpackage [dmldefaults]:
//   - [dmldefaults.Tokens]: Global baseline table covering all token types.
//   - [dmldefaults.Expression]: Baseline configuration for arithmetic, logical, and relational expressions.
//   - [dmldefaults.Cdecl]: Baseline configuration for C-type declarations and pointer syntax.
//
// Contexts are managed on [Printer] via scoped lifecycles:
//
//	p.newContext(format.NewTokenMap(overrides))
//	defer p.freeContext()
//
// # 3. Architectural Constraints & Boundaries
//
//   - Layout Engine Delegation: AST printing methods in `pkg/printer/` MUST NOT perform
//     manual string length calculations, manual whitespace injection, or manual line
//     wrapping within statements. All formatting decisions must be delegated to the
//     layout engine.
//   - Statement Boundaries: The sole exception to layout engine delegation is statement
//     separation (e.g. blank lines between declarations), which is handled at statement
//     boundaries via writeNewLines.
//   - Rule Annotation: AST print methods must be annotated with their corresponding
//     grammar rule using the tag comment: `// g4-tag: <rule_name>`.
//
// # 4. Token Pushing APIs
//
// AST visitor functions stream tokens to the layout engine using dedicated methods on [Printer]:
//
//   - `processToken(token antlr.Token, ctx *PrintContext) error`: Pushes a single token
//     with its context-resolved formatting flag.
//   - `pushTokenRange(endIdx TokenIndex) error`: Pushes all tokens between the last processed
//     index and `endIdx` (inclusive), skipping whitespace while handling hidden comments.
//   - `pushTokenRangeFromRightiest(token antlr.Tree) error`: Finds the rightmost token in an
//     AST subtree (including trailing same-line comments) and pushes the range.
//   - `pushTokenWithCustomFormat(node antlr.TerminalNode, flagFormatter func(format.FormatFlag) format.FormatFlag) error`:
//     Pushes a specific terminal node with an on-the-fly flag transformation.
//   - Passthrough Mode (`printUnimplemented`): Unimplemented grammar rules bypass formatting
//     by reading the exact character interval from the ANTLR input stream and emitting raw text.
//
// # Usage Example
//
//	file, err := dmlex.ParseFile("device.dml")
//	if err != nil {
//	    log.Fatalf("parse error: %v", err)
//	}
//
//	var out bytes.Buffer
//	lineCap := 100
//	if err := printer.PrinterWriteTo(file.Parser.Dml(), file, &out, lineCap); err != nil {
//	    log.Fatalf("format error: %v", err)
//	}
//	fmt.Print(out.String())
package printer
