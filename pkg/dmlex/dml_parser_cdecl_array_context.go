// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type CdeclArrayContext struct {
	Cdecl_ident_array_funcContext
}

func (s *CdeclArrayContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CdeclArrayContext) Cdecl_ident_array_func() ICdecl_ident_array_funcContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdecl_ident_array_funcContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdecl_ident_array_funcContext)
}

func (s *CdeclArrayContext) LBRACK() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACK, 0)
}

func (s *CdeclArrayContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *CdeclArrayContext) RBRACK() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACK, 0)
}

func (s *CdeclArrayContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdeclArray(s)
	}
}

func (s *CdeclArrayContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdeclArray(s)
	}
}

func NewCdeclArrayContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CdeclArrayContext {
	var p = new(CdeclArrayContext)

	InitEmptyCdecl_ident_array_funcContext(&p.Cdecl_ident_array_funcContext)
	p.parser = parser
	p.CopyAll(ctx.(*Cdecl_ident_array_funcContext))

	return p
}
