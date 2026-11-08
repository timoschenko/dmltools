// SPDX-License-Identifier: GPL-2.0-only
// Copyright (C) 2026 Sergei Timoshchenko

package test_dmlex

import (
	"testing"

	lex "github.com/timoschenko/dmltools/pkg/dmlex"

	"github.com/antlr4-go/antlr/v4"
)

type MethodCollector struct {
	*lex.BasedmlListener
	t       *testing.T
	Methods map[string]lex.IMethodUnifiedContext
}

func NewMethodCollector(t *testing.T) *MethodCollector {
	return &MethodCollector{
		BasedmlListener: new(lex.BasedmlListener),
		t:               t,
		Methods:         make(map[string]lex.IMethodUnifiedContext),
	}
}

func (v *MethodCollector) EnterMethodUnifiedDeclaration(
	ctx *lex.MethodUnifiedDeclarationContext,
) {
	name := lex.MethodName(ctx) // TODO: use qname instead
	if _, ok := v.Methods[name]; ok {
		v.t.Errorf("Method '%s' already defined", name)
	} else {
		v.Methods[name] = ctx
	}
}

func (v *MethodCollector) EnterMethodUnifiedDefinition(
	ctx *lex.MethodUnifiedDefinitionContext,
) {
	name := lex.MethodName(ctx) // TODO: use qname instead
	if _, ok := v.Methods[name]; ok {
		v.t.Errorf("Method '%s' already defined", name)
	} else {
		v.Methods[name] = ctx
	}
}

func ParseDmlFileFromText(
	t *testing.T,
	data string,
	filename string,
	linenoOffset int,
) *lex.File {
	input := antlr.NewInputStream(data)
	errorCounter := NewDmlLexErrorCounterListener(t, filename, linenoOffset)
	file := lex.ParseStream(input, errorCounter)

	return file
}

// Parsing text code
//
// Deprecated: ParseDmlText is deprecated, you have to use ParseDmlFileFromText instead
func ParseDmlText(
	t *testing.T,
	data string,
	filename string,
	linenoOffset int,
) lex.IDmlContext {
	file := ParseDmlFileFromText(t, data, filename, linenoOffset)
	tree := file.Parser.Dml()

	return tree
}
