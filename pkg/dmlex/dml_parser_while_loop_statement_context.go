// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IWhileLoopStatementContext is an interface to support dynamic dispatch.
type IWhileLoopStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WHILE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ConditionExpression() IConditionExpressionContext
	RPAREN() antlr.TerminalNode
	LoopStatement() ILoopStatementContext

	// IsWhileLoopStatementContext differentiates from other interfaces.
	IsWhileLoopStatementContext()
}

type WhileLoopStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*WhileLoopStatementContext) IsWhileLoopStatementContext() {}

func (s *WhileLoopStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *WhileLoopStatementContext) WHILE() antlr.TerminalNode {
	return s.GetToken(dmlParserWHILE, 0)
}

func (s *WhileLoopStatementContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *WhileLoopStatementContext) ConditionExpression() IConditionExpressionContext {
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

func (s *WhileLoopStatementContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *WhileLoopStatementContext) LoopStatement() ILoopStatementContext {
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

func (s *WhileLoopStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhileLoopStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WhileLoopStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterWhileLoopStatement(s)
	}
}

func (s *WhileLoopStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitWhileLoopStatement(s)
	}
}

func NewEmptyWhileLoopStatementContext() *WhileLoopStatementContext {
	var p = new(WhileLoopStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_whileLoopStatement
	return p
}

func InitEmptyWhileLoopStatementContext(p *WhileLoopStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_whileLoopStatement
}

func NewWhileLoopStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhileLoopStatementContext {
	var p = new(WhileLoopStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_whileLoopStatement

	return p
}
