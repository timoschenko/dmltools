// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type LiteralAtomContext struct {
	ExpressionContext
}

func (s *LiteralAtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralAtomContext) ExpressionIdent() IExpressionIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionIdentContext)
}

func (s *LiteralAtomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLiteralAtom(s)
	}
}

func (s *LiteralAtomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLiteralAtom(s)
	}
}

func NewLiteralAtomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LiteralAtomContext {
	var p = new(LiteralAtomContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
