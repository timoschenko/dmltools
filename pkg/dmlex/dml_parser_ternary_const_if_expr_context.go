// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type TernaryConstIfExprContext struct {
	ExpressionContext
}

func (s *TernaryConstIfExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TernaryConstIfExprContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *TernaryConstIfExprContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *TernaryConstIfExprContext) ALT_QUESTION() antlr.TerminalNode {
	return s.GetToken(dmlParserALT_QUESTION, 0)
}

func (s *TernaryConstIfExprContext) ALT_COLON() antlr.TerminalNode {
	return s.GetToken(dmlParserALT_COLON, 0)
}

func (s *TernaryConstIfExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterTernaryConstIfExpr(s)
	}
}

func (s *TernaryConstIfExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitTernaryConstIfExpr(s)
	}
}

func NewTernaryConstIfExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TernaryConstIfExprContext {
	var p = new(TernaryConstIfExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
