// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type CdeclFuncArgsContext struct {
	Cdecl_ident_array_funcContext
}

func (s *CdeclFuncArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CdeclFuncArgsContext) Cdecl_ident_array_func() ICdecl_ident_array_funcContext {
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

func (s *CdeclFuncArgsContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *CdeclFuncArgsContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *CdeclFuncArgsContext) Cdecl_list_variadic() ICdecl_list_variadicContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdecl_list_variadicContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdecl_list_variadicContext)
}

func (s *CdeclFuncArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdeclFuncArgs(s)
	}
}

func (s *CdeclFuncArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdeclFuncArgs(s)
	}
}

func NewCdeclFuncArgsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CdeclFuncArgsContext {
	var p = new(CdeclFuncArgsContext)

	InitEmptyCdecl_ident_array_funcContext(&p.Cdecl_ident_array_funcContext)
	p.parser = parser
	p.CopyAll(ctx.(*Cdecl_ident_array_funcContext))

	return p
}
