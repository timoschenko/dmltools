// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type MultiplicativeExprContext struct {
	ExpressionContext
	op antlr.Token
}

func (s *MultiplicativeExprContext) GetOp() antlr.Token { return s.op }

func (s *MultiplicativeExprContext) SetOp(v antlr.Token) { s.op = v }

func (s *MultiplicativeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicativeExprContext) AllExpression() []IExpressionContext {
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

func (s *MultiplicativeExprContext) Expression(i int) IExpressionContext {
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

func (s *MultiplicativeExprContext) MUL() antlr.TerminalNode {
	return s.GetToken(dmlParserMUL, 0)
}

func (s *MultiplicativeExprContext) DIV() antlr.TerminalNode {
	return s.GetToken(dmlParserDIV, 0)
}

func (s *MultiplicativeExprContext) MOD() antlr.TerminalNode {
	return s.GetToken(dmlParserMOD, 0)
}

func (s *MultiplicativeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterMultiplicativeExpr(s)
	}
}

func (s *MultiplicativeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitMultiplicativeExpr(s)
	}
}

func NewMultiplicativeExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MultiplicativeExprContext {
	var p = new(MultiplicativeExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
