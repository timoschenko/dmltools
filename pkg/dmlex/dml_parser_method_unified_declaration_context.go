// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type MethodUnifiedDeclarationContext struct {
	MethodUnifiedContext
}

func (s *MethodUnifiedDeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MethodUnifiedDeclarationContext) METHOD() antlr.TerminalNode {
	return s.GetToken(dmlParserMETHOD, 0)
}

func (s *MethodUnifiedDeclarationContext) Objident() IObjidentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjidentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjidentContext)
}

func (s *MethodUnifiedDeclarationContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *MethodUnifiedDeclarationContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *MethodUnifiedDeclarationContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *MethodUnifiedDeclarationContext) MethodQualifiers() IMethodQualifiersContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMethodQualifiersContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMethodQualifiersContext)
}

func (s *MethodUnifiedDeclarationContext) MethodParamsMaybeUntyped() IMethodParamsMaybeUntypedContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMethodParamsMaybeUntypedContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMethodParamsMaybeUntypedContext)
}

func (s *MethodUnifiedDeclarationContext) MethodOutParams() IMethodOutParamsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMethodOutParamsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMethodOutParamsContext)
}

func (s *MethodUnifiedDeclarationContext) THROWS() antlr.TerminalNode {
	return s.GetToken(dmlParserTHROWS, 0)
}

func (s *MethodUnifiedDeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterMethodUnifiedDeclaration(s)
	}
}

func (s *MethodUnifiedDeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitMethodUnifiedDeclaration(s)
	}
}

func NewMethodUnifiedDeclarationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MethodUnifiedDeclarationContext {
	var p = new(MethodUnifiedDeclarationContext)

	InitEmptyMethodUnifiedContext(&p.MethodUnifiedContext)
	p.parser = parser
	p.CopyAll(ctx.(*MethodUnifiedContext))

	return p
}
