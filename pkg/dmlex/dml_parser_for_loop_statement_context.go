// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IForLoopStatementContext is an interface to support dynamic dispatch.
type IForLoopStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FOR() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AllSEMI() []antlr.TerminalNode
	SEMI(i int) antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	LoopStatement() ILoopStatementContext
	ForLoopInitStatement() IForLoopInitStatementContext
	ConditionExpression() IConditionExpressionContext
	ForLoopStepStatement() IForLoopStepStatementContext

	// IsForLoopStatementContext differentiates from other interfaces.
	IsForLoopStatementContext()
}

type ForLoopStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ForLoopStatementContext) IsForLoopStatementContext() {}

func (s *ForLoopStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ForLoopStatementContext) FOR() antlr.TerminalNode {
	return s.GetToken(dmlParserFOR, 0)
}

func (s *ForLoopStatementContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *ForLoopStatementContext) AllSEMI() []antlr.TerminalNode {
	return s.GetTokens(dmlParserSEMI)
}

func (s *ForLoopStatementContext) SEMI(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, i)
}

func (s *ForLoopStatementContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *ForLoopStatementContext) LoopStatement() ILoopStatementContext {
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

func (s *ForLoopStatementContext) ForLoopInitStatement() IForLoopInitStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForLoopInitStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForLoopInitStatementContext)
}

func (s *ForLoopStatementContext) ConditionExpression() IConditionExpressionContext {
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

func (s *ForLoopStatementContext) ForLoopStepStatement() IForLoopStepStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IForLoopStepStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IForLoopStepStatementContext)
}

func (s *ForLoopStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForLoopStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForLoopStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterForLoopStatement(s)
	}
}

func (s *ForLoopStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitForLoopStatement(s)
	}
}

func NewEmptyForLoopStatementContext() *ForLoopStatementContext {
	var p = new(ForLoopStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_forLoopStatement
	return p
}

func InitEmptyForLoopStatementContext(p *ForLoopStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_forLoopStatement
}

func NewForLoopStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForLoopStatementContext {
	var p = new(ForLoopStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_forLoopStatement

	return p
}
