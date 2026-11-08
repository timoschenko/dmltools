// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type AddressofExprContext struct {
	ExpressionContext
}

func (s *AddressofExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AddressofExprContext) BITAND() antlr.TerminalNode {
	return s.GetToken(dmlParserBITAND, 0)
}

func (s *AddressofExprContext) Expression() IExpressionContext {
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

func (s *AddressofExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterAddressofExpr(s)
	}
}

func (s *AddressofExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitAddressofExpr(s)
	}
}

func NewAddressofExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AddressofExprContext {
	var p = new(AddressofExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
