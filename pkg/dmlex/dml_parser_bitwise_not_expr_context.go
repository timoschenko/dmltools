// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type BitwiseNotExprContext struct {
	ExpressionContext
}

func (s *BitwiseNotExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BitwiseNotExprContext) TILDE() antlr.TerminalNode {
	return s.GetToken(dmlParserTILDE, 0)
}

func (s *BitwiseNotExprContext) Expression() IExpressionContext {
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

func (s *BitwiseNotExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterBitwiseNotExpr(s)
	}
}

func (s *BitwiseNotExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitBitwiseNotExpr(s)
	}
}

func NewBitwiseNotExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BitwiseNotExprContext {
	var p = new(BitwiseNotExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
