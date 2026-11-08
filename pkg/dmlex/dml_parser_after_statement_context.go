// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IAfterStatementContext is an interface to support dynamic dispatch.
type IAfterStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AFTER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Expression() IExpressionContext
	SEMI() antlr.TerminalNode
	AfterDelayCondition() IAfterDelayConditionContext

	// IsAfterStatementContext differentiates from other interfaces.
	IsAfterStatementContext()
}

type AfterStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*AfterStatementContext) IsAfterStatementContext() {}

func (s *AfterStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *AfterStatementContext) AFTER() antlr.TerminalNode {
	return s.GetToken(dmlParserAFTER, 0)
}

func (s *AfterStatementContext) COLON() antlr.TerminalNode {
	return s.GetToken(dmlParserCOLON, 0)
}

func (s *AfterStatementContext) Expression() IExpressionContext {
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

func (s *AfterStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *AfterStatementContext) AfterDelayCondition() IAfterDelayConditionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAfterDelayConditionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAfterDelayConditionContext)
}

func (s *AfterStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AfterStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AfterStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterAfterStatement(s)
	}
}

func (s *AfterStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitAfterStatement(s)
	}
}

func NewEmptyAfterStatementContext() *AfterStatementContext {
	var p = new(AfterStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_afterStatement
	return p
}

func InitEmptyAfterStatementContext(p *AfterStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_afterStatement
}

func NewAfterStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AfterStatementContext {
	var p = new(AfterStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_afterStatement

	return p
}
