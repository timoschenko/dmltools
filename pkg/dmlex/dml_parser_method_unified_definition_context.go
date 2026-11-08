// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type MethodUnifiedDefinitionContext struct {
	MethodUnifiedContext
}

func (s *MethodUnifiedDefinitionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MethodUnifiedDefinitionContext) METHOD() antlr.TerminalNode {
	return s.GetToken(dmlParserMETHOD, 0)
}

func (s *MethodUnifiedDefinitionContext) Objident() IObjidentContext {
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

func (s *MethodUnifiedDefinitionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *MethodUnifiedDefinitionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *MethodUnifiedDefinitionContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *MethodUnifiedDefinitionContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *MethodUnifiedDefinitionContext) LBRACE_NEW() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE_NEW, 0)
}

func (s *MethodUnifiedDefinitionContext) MethodQualifiers() IMethodQualifiersContext {
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

func (s *MethodUnifiedDefinitionContext) MethodParamsMaybeUntyped() IMethodParamsMaybeUntypedContext {
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

func (s *MethodUnifiedDefinitionContext) MethodOutParams() IMethodOutParamsContext {
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

func (s *MethodUnifiedDefinitionContext) THROWS() antlr.TerminalNode {
	return s.GetToken(dmlParserTHROWS, 0)
}

func (s *MethodUnifiedDefinitionContext) AllStatement() []IStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IStatementContext); ok {
			len++
		}
	}

	tst := make([]IStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IStatementContext); ok {
			tst[i] = t.(IStatementContext)
			i++
		}
	}

	return tst
}

func (s *MethodUnifiedDefinitionContext) Statement(i int) IStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *MethodUnifiedDefinitionContext) DEFAULT() antlr.TerminalNode {
	return s.GetToken(dmlParserDEFAULT, 0)
}

func (s *MethodUnifiedDefinitionContext) DEFAULT_NEW() antlr.TerminalNode {
	return s.GetToken(dmlParserDEFAULT_NEW, 0)
}

func (s *MethodUnifiedDefinitionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterMethodUnifiedDefinition(s)
	}
}

func (s *MethodUnifiedDefinitionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitMethodUnifiedDefinition(s)
	}
}

func NewMethodUnifiedDefinitionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MethodUnifiedDefinitionContext {
	var p = new(MethodUnifiedDefinitionContext)

	InitEmptyMethodUnifiedContext(&p.MethodUnifiedContext)
	p.parser = parser
	p.CopyAll(ctx.(*MethodUnifiedContext))

	return p
}
