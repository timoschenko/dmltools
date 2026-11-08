// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type CdeclIdentContext struct {
	Cdecl_ident_array_funcContext
}

func (s *CdeclIdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CdeclIdentContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *CdeclIdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdeclIdent(s)
	}
}

func (s *CdeclIdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdeclIdent(s)
	}
}

func NewCdeclIdentContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CdeclIdentContext {
	var p = new(CdeclIdentContext)

	InitEmptyCdecl_ident_array_funcContext(&p.Cdecl_ident_array_funcContext)
	p.parser = parser
	p.CopyAll(ctx.(*Cdecl_ident_array_funcContext))

	return p
}
