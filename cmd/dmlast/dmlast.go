// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/antlr4-go/antlr/v4"

	"github.com/timoschenko/dmltools/pkg/dmlex"
)

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

type JsonNode struct {
	Type     string     `json:"type"`
	Kind     string     `json:"kind"`
	Text     string     `json:"text"`
	Children []JsonNode `json:"children,omitempty"`
}

// buildJsonTree is a recursive function that converts an ANTLR ParseTree
// into our serializable JsonNode struct.
func buildJsonTree(
	tree antlr.ParseTree,
	ruleNames []string,
	symbolicNames []string,
) JsonNode {
	const kindRule = "rule"
	const kindToken = "token"

	node := JsonNode{}
	node.Text = tree.GetText()

	// Check the type of the node (Rule vs. Terminal)
	switch v := tree.(type) {
	case antlr.RuleContext:
		// It's a Rule node (e.g., 'method', 'expression')
		// Get the rule name from the parser's list
		node.Kind = kindRule
		node.Type = ruleNames[v.GetRuleIndex()]
	case antlr.TerminalNode:
		// It's a Terminal node (a Token e.g., 'METHOD', 'IDENTIFIER', ';')
		// Get the symbolic name from the parser's list
		node.Kind = kindToken
		tokenType := v.GetSymbol().GetTokenType()
		if tokenType > 0 { // Avoid EOF token which can be -1
			node.Type = symbolicNames[tokenType]
		} else {
			node.Type = "EOF"
		}
	}

	// Recursively build children
	childCount := tree.GetChildCount()
	if childCount > 0 {
		node.Children = make([]JsonNode, 0, childCount)
		for i := range childCount {
			child := tree.GetChild(i)
			parseTreeChild, ok := child.(antlr.ParseTree)
			if !ok {
				continue
			}

			parsedChild := buildJsonTree(
				parseTreeChild,
				ruleNames,
				symbolicNames,
			)
			node.Children = append(node.Children, parsedChild)
		}
	}

	return node
}

type AstSerializerContext struct {
	RuleNames     []string
	SymbolicNames []string
}

func NewAstSerializerContext(p *antlr.BaseLexer) *AstSerializerContext {
	return &AstSerializerContext{
		RuleNames:     p.GetRuleNames(),
		SymbolicNames: p.GetSymbolicNames(),
	}
}

func printVersion(w io.Writer) {
	ver := Version
	commit := GitCommit
	date := BuildDate

	if (ver == "dev" || ver == "") && commit == "none" {
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				ver = info.Main.Version
			}
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					if len(setting.Value) > 7 {
						commit = setting.Value[:7]
					} else {
						commit = setting.Value
					}
				case "vcs.time":
					date = setting.Value
				}
			}
		}
	}

	if commit != "none" || date != "unknown" {
		_, _ = fmt.Fprintf(
			w,
			"dmlast version %s (commit: %s, built: %s)\n",
			ver,
			commit,
			date,
		)
	} else {
		_, _ = fmt.Fprintf(w, "dmlast version %s\n\n", ver)
	}
	_, _ = fmt.Fprintln(
		w,
		"License: GNU General Public License v2.0 only (GPL-2.0-only)",
	)
	_, _ = fmt.Fprintln(
		w,
		"This is free software: you are free to change and redistribute it.",
	)
	_, _ = fmt.Fprintln(
		w,
		"There is NO WARRANTY, to the extent permitted by law.",
	)
}

func main() {
	filePathPtr := flag.String("file", "", "The path to the DML file to parse.")
	versionPtr := flag.Bool("version", false, "Show the version and exit.")
	vPtr := flag.Bool("v", false, "Show the version and exit.")
	VPtr := flag.Bool("V", false, "Show the version and exit.")
	flag.Parse() // Parse the command-line arguments

	if *versionPtr || *vPtr || *VPtr {
		printVersion(os.Stdout)
		return
	}

	// Check if the file path was provided
	if *filePathPtr == "" {
		fmt.Fprintf(os.Stderr, "Error: No file path provided.\n")
		fmt.Fprintf(
			os.Stderr,
			"Usage: %s -file <path-to-dml-file>\n",
			os.Args[0],
		)
		os.Exit(1)
	}

	filename := *filePathPtr
	errorListener := dmlex.NewErrorListener(filename, 0)
	errorListener.Handler = func(e dmlex.SyntaxErrorInfo) {
		fmt.Fprintf(
			os.Stderr,
			"ERROR %s:%d column %d: Syntax error - %s\n",
			filename, e.Line, e.Column, e.Msg,
		)
	}
	file, err := dmlex.ParseFile(filename, errorListener)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		os.Exit(1)
	}

	p := file.Parser // alias
	tree := p.Dml()

	if errorCount := len(errorListener.Errors); errorCount > 0 {
		fmt.Fprintf(
			os.Stderr,
			"Parsing failed with %d syntax error(s).\n",
			errorCount,
		)
		os.Exit(1)
	}

	jsonTree := buildJsonTree(tree, p.GetRuleNames(), p.GetSymbolicNames())

	// Marshal to JSON with pretty-printing
	jsonData, err := json.MarshalIndent(jsonTree, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error serializing tree to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write to standard output
	if _, err := fmt.Println(string(jsonData)); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON to stdout: %v\n", err)
		os.Exit(1)
	}
}
