// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

// Package dmldefaults provides baseline token formatting tables, expression rules,
// and C-declaration formatting configurations for the DML formatter.
//
// # Overview
//
// The DML formatter relies on baseline token configuration tables implementing [format.TokenMap].
// These maps assign initial [format.FormatFlag] attributes (such as spacing separators,
// break priorities, indentation changes, and pair conflict priorities) to each [dmlex.TokenType].
//
// # Key Tables
//
//   - [Tokens]: Global baseline formatting map covering all DML token types.
//   - [Expression]: Baseline configuration for arithmetic, logical, and relational expressions.
//   - [Cdecl]: Baseline configuration for C-type declarations and pointer syntax.
package dmldefaults
