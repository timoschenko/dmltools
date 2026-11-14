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

const previewWarning = "WARNING: dmlast is a preview tool; " +
	"features and behavior may change without notice."

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
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, previewWarning)
}

func printUsage(out io.Writer, fs *flag.FlagSet) {
	_, _ = fmt.Fprintf(out, "Usage: dmlast -file <path-to-dml-file>\n\n")
	_, _ = fmt.Fprintf(out, "%s\n\n", previewWarning)
	_, _ = fmt.Fprintf(out, "Options:\n")
	fs.SetOutput(out)
	fs.PrintDefaults()
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("dmlast", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var filePath string
	var showVersion bool
	var showHelp bool

	fs.StringVar(&filePath, "file", "", "The path to the DML file to parse.")
	fs.BoolVar(&showVersion, "version", false, "Show the version and exit.")
	fs.BoolVar(&showVersion, "v", false, "Show the version and exit.")
	fs.BoolVar(&showVersion, "V", false, "Show the version and exit.")
	fs.BoolVar(&showHelp, "help", false, "Show help usage and exit.")
	fs.BoolVar(&showHelp, "h", false, "Show help usage and exit.")
	fs.Usage = func() {
		printUsage(stderr, fs)
	}

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if showHelp {
		printUsage(stdout, fs)
		return 0
	}

	if showVersion {
		printVersion(stdout)
		return 0
	}

	if filePath == "" {
		_, _ = fmt.Fprintf(stderr, "Error: No file path provided.\n")
		printUsage(stderr, fs)
		return 1
	}

	_, _ = fmt.Fprintln(stderr, previewWarning)

	errorListener := dmlex.NewErrorListener(filePath, 0)
	errorListener.Handler = func(e dmlex.SyntaxErrorInfo) {
		_, _ = fmt.Fprintf(
			stderr,
			"ERROR %s:%d column %d: Syntax error - %s\n",
			filePath, e.Line, e.Column, e.Msg,
		)
	}
	file, err := dmlex.ParseFile(filePath, errorListener)
	if err != nil {
		_, _ = fmt.Fprintf(
			stderr,
			"Error reading file %s: %v\n",
			filePath,
			err,
		)
		return 1
	}

	p := file.Parser // alias
	tree := p.Dml()

	if errorCount := len(errorListener.Errors); errorCount > 0 {
		_, _ = fmt.Fprintf(
			stderr,
			"Parsing failed with %d syntax error(s).\n",
			errorCount,
		)
		return 1
	}

	jsonTree := buildJsonTree(tree, p.GetRuleNames(), p.GetSymbolicNames())

	// Marshal to JSON with pretty-printing
	jsonData, err := json.MarshalIndent(jsonTree, "", "  ")
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "Error serializing tree to JSON: %v\n", err)
		return 1
	}

	// Write to standard output
	if _, err := fmt.Fprintln(stdout, string(jsonData)); err != nil {
		_, _ = fmt.Fprintf(stderr, "Error writing JSON to stdout: %v\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
