// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type EachExprContext struct {
	ExpressionContext
}

func (s *EachExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EachExprContext) EACH() antlr.TerminalNode {
	return s.GetToken(dmlParserEACH, 0)
}

func (s *EachExprContext) Objident() IObjidentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjidentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjidentContext)
}

func (s *EachExprContext) IN() antlr.TerminalNode {
	return s.GetToken(dmlParserIN, 0)
}

func (s *EachExprContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *EachExprContext) Expression() IExpressionContext {
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

func (s *EachExprContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *EachExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterEachExpr(s)
	}
}

func (s *EachExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitEachExpr(s)
	}
}

func NewEachExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EachExprContext {
	var p = new(EachExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
