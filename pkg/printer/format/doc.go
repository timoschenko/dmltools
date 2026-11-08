// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

// Package format defines token formatting primitives, Unified Formatting Flags (UFF),
// spacing separators, break priorities, and JSON serialization support for the DML formatter.
//
// # Overview
//
// The DML formatter uses Unified Formatting Flags (UFF) to represent token-level
// formatting decisions as compact 32-bit bitfields ([FormatFlag]). Each flag encodes
// break priorities before and after, spacing separators ([Separator]), indentation shifts,
// force-wrap directives, and token pair conflict priorities ([Priority]).
//
// # Key Components
//
//   - [FormatFlag]: 32-bit bitfield encoding spacing rules, break priorities,
//     indentation shifts, and pair conflict priorities for a token.
//   - [Separator]: Enumeration of mutually exclusive whitespace and line break states
//     between adjacent tokens (e.g. [SepSpace], [SepNone], [SepSpaceMaybeWrap]).
//   - [Priority]: 4-bit priority level (0-15) governing line break candidate selection
//     and token pair conflict resolution.
//   - [Winner]: Enumeration representing the winner of token pair conflict resolution
//     ([WinnerLeft], [WinnerRight], [WinnerEqual]).
//   - [SideFlagJSON]: JSON representation of grouped boundary properties for a token.
//   - [TokenMap]: Interface providing format flag lookups for token types.
//   - [TokenMapImpl]: Concrete map of token types to format flags with JSON serialization support.
package format
