// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type SizeofExprContext struct {
	ExpressionContext
}

func (s *SizeofExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SizeofExprContext) SIZEOF() antlr.TerminalNode {
	return s.GetToken(dmlParserSIZEOF, 0)
}

func (s *SizeofExprContext) Expression() IExpressionContext {
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

func (s *SizeofExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSizeofExpr(s)
	}
}

func (s *SizeofExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSizeofExpr(s)
	}
}

func NewSizeofExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SizeofExprContext {
	var p = new(SizeofExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
