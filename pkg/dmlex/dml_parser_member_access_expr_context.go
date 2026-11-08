// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

type MemberAccessExprContext struct {
	ExpressionContext
}

func (s *MemberAccessExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberAccessExprContext) Expression() IExpressionContext {
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

func (s *MemberAccessExprContext) AllExpressionAccessIdent() []IExpressionAccessIdentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionAccessIdentContext); ok {
			len++
		}
	}

	tst := make([]IExpressionAccessIdentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionAccessIdentContext); ok {
			tst[i] = t.(IExpressionAccessIdentContext)
			i++
		}
	}

	return tst
}

func (s *MemberAccessExprContext) ExpressionAccessIdent(i int) IExpressionAccessIdentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionAccessIdentContext); ok {
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

	return t.(IExpressionAccessIdentContext)
}

func (s *MemberAccessExprContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(dmlParserDOT)
}

func (s *MemberAccessExprContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserDOT, i)
}

func (s *MemberAccessExprContext) AllARROW() []antlr.TerminalNode {
	return s.GetTokens(dmlParserARROW)
}

func (s *MemberAccessExprContext) ARROW(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserARROW, i)
}

func (s *MemberAccessExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterMemberAccessExpr(s)
	}
}

func (s *MemberAccessExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitMemberAccessExpr(s)
	}
}

func NewMemberAccessExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberAccessExprContext {
	var p = new(MemberAccessExprContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}
