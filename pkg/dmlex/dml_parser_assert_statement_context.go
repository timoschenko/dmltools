// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IAssertStatementContext is an interface to support dynamic dispatch.
type IAssertStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASSERT() antlr.TerminalNode
	ConditionExpression() IConditionExpressionContext
	SEMI() antlr.TerminalNode

	// IsAssertStatementContext differentiates from other interfaces.
	IsAssertStatementContext()
}

type AssertStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*AssertStatementContext) IsAssertStatementContext() {}

func (s *AssertStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *AssertStatementContext) ASSERT() antlr.TerminalNode {
	return s.GetToken(dmlParserASSERT, 0)
}

func (s *AssertStatementContext) ConditionExpression() IConditionExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConditionExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConditionExpressionContext)
}

func (s *AssertStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *AssertStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssertStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssertStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterAssertStatement(s)
	}
}

func (s *AssertStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitAssertStatement(s)
	}
}

func NewEmptyAssertStatementContext() *AssertStatementContext {
	var p = new(AssertStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_assertStatement
	return p
}

func InitEmptyAssertStatementContext(p *AssertStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_assertStatement
}

func NewAssertStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssertStatementContext {
	var p = new(AssertStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_assertStatement

	return p
}
