// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type DefinedExprContext struct {
	ExpressionContext
}

func (s *DefinedExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DefinedExprContext) DEFINED() antlr.TerminalNode {
	return s.GetToken(dmlParserDEFINED, 0)
}

func (s *DefinedExprContext) Expression() IExpressionContext {
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

func (s *DefinedExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDefinedExpr(s)
	}
}

func (s *DefinedExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDefinedExpr(s)
	}
}

func NewDefinedExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DefinedExprContext {
	var p = new(DefinedExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
