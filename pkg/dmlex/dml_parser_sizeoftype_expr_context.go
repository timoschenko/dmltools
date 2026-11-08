// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type SizeoftypeExprContext struct {
	ExpressionContext
}

func (s *SizeoftypeExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SizeoftypeExprContext) SIZEOFTYPE() antlr.TerminalNode {
	return s.GetToken(dmlParserSIZEOFTYPE, 0)
}

func (s *SizeoftypeExprContext) Typeoparg() ITypeopargContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeopargContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeopargContext)
}

func (s *SizeoftypeExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSizeoftypeExpr(s)
	}
}

func (s *SizeoftypeExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSizeoftypeExpr(s)
	}
}

func NewSizeoftypeExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SizeoftypeExprContext {
	var p = new(SizeoftypeExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
