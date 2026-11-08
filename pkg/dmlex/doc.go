// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

// Package dmlex provides lexing, parsing, and AST representations for
// the Device Modeling Language (DML 1.4 & 1.2).
//
// # Overview
//
// The package embeds an ANTLR4-generated lexer and parser derived from the
// formal grammar definition in grammar/dml/dml.g4. To optimize compilation
// performance and code maintainability, monolithic ANTLR parser files are
// decomposed during the build process via the internal baker tool into modular,
// per-rule context structures.
//
// # Key Components
//
//   - [File]: Represents a parsed DML source file, encapsulating an [antlr.CharStream]
//     and the underlying parser instance.
//   - [ParseFile] and [ParseStream]: High-level entry points for loading and
//     parsing DML source code into an AST. They accept optional error and
//     parse tree listeners.
//   - [TokenType]: Strongly typed enumeration classifying DML tokens
//     (keywords, operators, delimiters, and literals) with string conversion and
//     JSON serialization support.
//   - [ErrorListener]: An ANTLR error listener implementation that collects
//     syntax diagnostics and maps error positions to source lines and columns.
//   - Context Rules: Typed parse tree nodes (e.g. [DmlContext], [DeviceStatementContext],
//     [MethodUnifiedContext]) providing strongly typed access to DML AST elements.
//
// # Usage Example
//
//	file, err := dmlex.ParseFile("model.dml")
//	if err != nil {
//	    log.Fatalf("failed to parse DML file: %v", err)
//	}
//
//	// Obtain the root AST context
//	tree := file.Parser.Dml()
//	fmt.Println("Parsed DML AST root:", tree.GetText())
package dmlex
