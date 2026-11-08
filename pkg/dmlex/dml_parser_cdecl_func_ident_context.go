// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type CdeclFuncIdentContext struct {
	Cdecl_ident_array_funcContext
}

func (s *CdeclFuncIdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CdeclFuncIdentContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *CdeclFuncIdentContext) Cdecl_qualifier_ident_array_func() ICdecl_qualifier_ident_array_funcContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdecl_qualifier_ident_array_funcContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdecl_qualifier_ident_array_funcContext)
}

func (s *CdeclFuncIdentContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *CdeclFuncIdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdeclFuncIdent(s)
	}
}

func (s *CdeclFuncIdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdeclFuncIdent(s)
	}
}

func NewCdeclFuncIdentContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CdeclFuncIdentContext {
	var p = new(CdeclFuncIdentContext)

	InitEmptyCdecl_ident_array_funcContext(&p.Cdecl_ident_array_funcContext)
	p.parser = parser
	p.CopyAll(ctx.(*Cdecl_ident_array_funcContext))

	return p
}
