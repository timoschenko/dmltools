// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type CdeclNoneContext struct {
	Cdecl_ident_array_funcContext
}

func (s *CdeclNoneContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CdeclNoneContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdeclNone(s)
	}
}

func (s *CdeclNoneContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdeclNone(s)
	}
}

func NewCdeclNoneContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CdeclNoneContext {
	var p = new(CdeclNoneContext)

	InitEmptyCdecl_ident_array_funcContext(&p.Cdecl_ident_array_funcContext)
	p.parser = parser
	p.CopyAll(ctx.(*Cdecl_ident_array_funcContext))

	return p
}
