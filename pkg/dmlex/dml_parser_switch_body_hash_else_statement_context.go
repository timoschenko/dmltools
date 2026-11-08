// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISwitchBodyHashElseStatementContext is an interface to support dynamic dispatch.
type ISwitchBodyHashElseStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SwitchBodyStatement() ISwitchBodyStatementContext

	// IsSwitchBodyHashElseStatementContext differentiates from other interfaces.
	IsSwitchBodyHashElseStatementContext()
}

type SwitchBodyHashElseStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*SwitchBodyHashElseStatementContext) IsSwitchBodyHashElseStatementContext() {}

func (s *SwitchBodyHashElseStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *SwitchBodyHashElseStatementContext) SwitchBodyStatement() ISwitchBodyStatementContext {
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

func (s *SwitchBodyHashElseStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SwitchBodyHashElseStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SwitchBodyHashElseStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSwitchBodyHashElseStatement(s)
	}
}

func (s *SwitchBodyHashElseStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSwitchBodyHashElseStatement(s)
	}
}

func NewEmptySwitchBodyHashElseStatementContext() *SwitchBodyHashElseStatementContext {
	var p = new(SwitchBodyHashElseStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_switchBodyHashElseStatement
	return p
}

func InitEmptySwitchBodyHashElseStatementContext(p *SwitchBodyHashElseStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_switchBodyHashElseStatement
}

func NewSwitchBodyHashElseStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchBodyHashElseStatementContext {
	var p = new(SwitchBodyHashElseStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_switchBodyHashElseStatement

	return p
}
