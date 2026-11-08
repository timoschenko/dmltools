// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type PostIncrementExprContext struct {
	ExpressionContext
	op antlr.Token
}

func (s *PostIncrementExprContext) GetOp() antlr.Token { return s.op }

func (s *PostIncrementExprContext) SetOp(v antlr.Token) { s.op = v }

func (s *PostIncrementExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PostIncrementExprContext) Expression() IExpressionContext {
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

func (s *PostIncrementExprContext) INC() antlr.TerminalNode {
	return s.GetToken(dmlParserINC, 0)
}

func (s *PostIncrementExprContext) DEC() antlr.TerminalNode {
	return s.GetToken(dmlParserDEC, 0)
}

func (s *PostIncrementExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterPostIncrementExpr(s)
	}
}

func (s *PostIncrementExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitPostIncrementExpr(s)
	}
}

func NewPostIncrementExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PostIncrementExprContext {
	var p = new(PostIncrementExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
