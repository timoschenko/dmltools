// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IDoWhileLoopStatementContext is an interface to support dynamic dispatch.
type IDoWhileLoopStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DO() antlr.TerminalNode
	LoopStatement() ILoopStatementContext
	WHILE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ConditionExpression() IConditionExpressionContext
	RPAREN() antlr.TerminalNode
	SEMI() antlr.TerminalNode

	// IsDoWhileLoopStatementContext differentiates from other interfaces.
	IsDoWhileLoopStatementContext()
}

type DoWhileLoopStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*DoWhileLoopStatementContext) IsDoWhileLoopStatementContext() {}

func (s *DoWhileLoopStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *DoWhileLoopStatementContext) DO() antlr.TerminalNode {
	return s.GetToken(dmlParserDO, 0)
}

func (s *DoWhileLoopStatementContext) LoopStatement() ILoopStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILoopStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILoopStatementContext)
}

func (s *DoWhileLoopStatementContext) WHILE() antlr.TerminalNode {
	return s.GetToken(dmlParserWHILE, 0)
}

func (s *DoWhileLoopStatementContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *DoWhileLoopStatementContext) ConditionExpression() IConditionExpressionContext {
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

func (s *DoWhileLoopStatementContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *DoWhileLoopStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *DoWhileLoopStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DoWhileLoopStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DoWhileLoopStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDoWhileLoopStatement(s)
	}
}

func (s *DoWhileLoopStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDoWhileLoopStatement(s)
	}
}

func NewEmptyDoWhileLoopStatementContext() *DoWhileLoopStatementContext {
	var p = new(DoWhileLoopStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_doWhileLoopStatement
	return p
}

func InitEmptyDoWhileLoopStatementContext(p *DoWhileLoopStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_doWhileLoopStatement
}

func NewDoWhileLoopStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DoWhileLoopStatementContext {
	var p = new(DoWhileLoopStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_doWhileLoopStatement

	return p
}
