// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISwitchBodyHashIfStatementContext is an interface to support dynamic dispatch.
type ISwitchBodyHashIfStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HASH_IF() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ConditionExpression() IConditionExpressionContext
	RPAREN() antlr.TerminalNode
	AllLBRACE() []antlr.TerminalNode
	LBRACE(i int) antlr.TerminalNode
	SwitchBodyStatement() ISwitchBodyStatementContext
	AllRBRACE() []antlr.TerminalNode
	RBRACE(i int) antlr.TerminalNode
	HASH_ELSE() antlr.TerminalNode
	SwitchBodyHashElseStatement() ISwitchBodyHashElseStatementContext

	// IsSwitchBodyHashIfStatementContext differentiates from other interfaces.
	IsSwitchBodyHashIfStatementContext()
}

type SwitchBodyHashIfStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*SwitchBodyHashIfStatementContext) IsSwitchBodyHashIfStatementContext() {}

func (s *SwitchBodyHashIfStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *SwitchBodyHashIfStatementContext) HASH_IF() antlr.TerminalNode {
	return s.GetToken(dmlParserHASH_IF, 0)
}

func (s *SwitchBodyHashIfStatementContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *SwitchBodyHashIfStatementContext) ConditionExpression() IConditionExpressionContext {
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

func (s *SwitchBodyHashIfStatementContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *SwitchBodyHashIfStatementContext) AllLBRACE() []antlr.TerminalNode {
	return s.GetTokens(dmlParserLBRACE)
}

func (s *SwitchBodyHashIfStatementContext) LBRACE(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, i)
}

func (s *SwitchBodyHashIfStatementContext) SwitchBodyStatement() ISwitchBodyStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISwitchBodyStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISwitchBodyStatementContext)
}

func (s *SwitchBodyHashIfStatementContext) AllRBRACE() []antlr.TerminalNode {
	return s.GetTokens(dmlParserRBRACE)
}

func (s *SwitchBodyHashIfStatementContext) RBRACE(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, i)
}

func (s *SwitchBodyHashIfStatementContext) HASH_ELSE() antlr.TerminalNode {
	return s.GetToken(dmlParserHASH_ELSE, 0)
}

func (s *SwitchBodyHashIfStatementContext) SwitchBodyHashElseStatement() ISwitchBodyHashElseStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISwitchBodyHashElseStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISwitchBodyHashElseStatementContext)
}

func (s *SwitchBodyHashIfStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SwitchBodyHashIfStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SwitchBodyHashIfStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSwitchBodyHashIfStatement(s)
	}
}

func (s *SwitchBodyHashIfStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSwitchBodyHashIfStatement(s)
	}
}

func NewEmptySwitchBodyHashIfStatementContext() *SwitchBodyHashIfStatementContext {
	var p = new(SwitchBodyHashIfStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_switchBodyHashIfStatement
	return p
}

func InitEmptySwitchBodyHashIfStatementContext(p *SwitchBodyHashIfStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_switchBodyHashIfStatement
}

func NewSwitchBodyHashIfStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchBodyHashIfStatementContext {
	var p = new(SwitchBodyHashIfStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_switchBodyHashIfStatement

	return p
}
