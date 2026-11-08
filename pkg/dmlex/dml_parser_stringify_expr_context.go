// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type StringifyExprContext struct {
	ExpressionContext
}

func (s *StringifyExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringifyExprContext) STRINGIFY() antlr.TerminalNode {
	return s.GetToken(dmlParserSTRINGIFY, 0)
}

func (s *StringifyExprContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *StringifyExprContext) Expression() IExpressionContext {
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

func (s *StringifyExprContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *StringifyExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterStringifyExpr(s)
	}
}

func (s *StringifyExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitStringifyExpr(s)
	}
}

func NewStringifyExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringifyExprContext {
	var p = new(StringifyExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
