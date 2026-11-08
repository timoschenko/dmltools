// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type PreIncrementExprContext struct {
	ExpressionContext
	op antlr.Token
}

func (s *PreIncrementExprContext) GetOp() antlr.Token { return s.op }

func (s *PreIncrementExprContext) SetOp(v antlr.Token) { s.op = v }

func (s *PreIncrementExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PreIncrementExprContext) Expression() IExpressionContext {
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

func (s *PreIncrementExprContext) INC() antlr.TerminalNode {
	return s.GetToken(dmlParserINC, 0)
}

func (s *PreIncrementExprContext) DEC() antlr.TerminalNode {
	return s.GetToken(dmlParserDEC, 0)
}

func (s *PreIncrementExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterPreIncrementExpr(s)
	}
}

func (s *PreIncrementExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitPreIncrementExpr(s)
	}
}

func NewPreIncrementExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PreIncrementExprContext {
	var p = new(PreIncrementExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
