// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

// Package layout provides token buffering, conflict resolution, line wrapping,
// indentation management, and stream emission for the DML formatter.
//
// # Overview
//
// The layout package receives tokens alongside their [format.FormatFlag] metadata,
// resolves separators and break priorities between adjacent tokens, tracks hierarchical
// layout groups, and dynamically wraps lines when the configured line capacity is exceeded.
//
// # Key Components
//
//   - [Engine]: Core layout engine buffering tokens, resolving separators, applying breaks,
//     and emitting formatted text to the destination writer.
//   - [State]: Indentation and column tracking state across token flushes.
//   - [Group]: Hierarchical grouping structure governing collective or greedy line breaks.
//   - [BufferedToken]: Buffered token paired with resolved layout cache and synthetic markers.
//   - [ResolveSeparator]: Three-stage conflict resolution between adjacent token format flags.
//   - [TokenTracer]: Diagnostic interface for recording formatting decisions between tokens.
package layout
