// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IAfterDelayConditionContext is an interface to support dynamic dispatch.
type IAfterDelayConditionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	IDENTIFIER() antlr.TerminalNode
	ARROW() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	IdentList() IIdentListContext
	Ident() IIdentContext

	// IsAfterDelayConditionContext differentiates from other interfaces.
	IsAfterDelayConditionContext()
}

type AfterDelayConditionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*AfterDelayConditionContext) IsAfterDelayConditionContext() {}

func (s *AfterDelayConditionContext) GetParser() antlr.Parser { return s.parser }

func (s *AfterDelayConditionContext) Expression() IExpressionContext {
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

func (s *AfterDelayConditionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(dmlParserIDENTIFIER, 0)
}

func (s *AfterDelayConditionContext) ARROW() antlr.TerminalNode {
	return s.GetToken(dmlParserARROW, 0)
}

func (s *AfterDelayConditionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *AfterDelayConditionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *AfterDelayConditionContext) IdentList() IIdentListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentListContext)
}

func (s *AfterDelayConditionContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *AfterDelayConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AfterDelayConditionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AfterDelayConditionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterAfterDelayCondition(s)
	}
}

func (s *AfterDelayConditionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitAfterDelayCondition(s)
	}
}

func NewEmptyAfterDelayConditionContext() *AfterDelayConditionContext {
	var p = new(AfterDelayConditionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_afterDelayCondition
	return p
}

func InitEmptyAfterDelayConditionContext(p *AfterDelayConditionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_afterDelayCondition
}

func NewAfterDelayConditionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AfterDelayConditionContext {
	var p = new(AfterDelayConditionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_afterDelayCondition

	return p
}
