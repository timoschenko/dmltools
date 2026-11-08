// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type BitSliceExprContext struct {
	ExpressionContext
}

func (s *BitSliceExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BitSliceExprContext) AllExpression() []IExpressionContext {
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

func (s *BitSliceExprContext) Expression(i int) IExpressionContext {
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

func (s *BitSliceExprContext) LBRACK() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACK, 0)
}

func (s *BitSliceExprContext) BitSlice() IBitSliceContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBitSliceContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBitSliceContext)
}

func (s *BitSliceExprContext) RBRACK() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACK, 0)
}

func (s *BitSliceExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterBitSliceExpr(s)
	}
}

func (s *BitSliceExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitBitSliceExpr(s)
	}
}

func NewBitSliceExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BitSliceExprContext {
	var p = new(BitSliceExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
