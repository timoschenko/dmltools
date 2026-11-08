// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISwitchBodyStatementContext is an interface to support dynamic dispatch.
type ISwitchBodyStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCaseStatement() []ICaseStatementContext
	CaseStatement(i int) ICaseStatementContext

	// IsSwitchBodyStatementContext differentiates from other interfaces.
	IsSwitchBodyStatementContext()
}

type SwitchBodyStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*SwitchBodyStatementContext) IsSwitchBodyStatementContext() {}

func (s *SwitchBodyStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *SwitchBodyStatementContext) AllCaseStatement() []ICaseStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICaseStatementContext); ok {
			len++
		}
	}

	tst := make([]ICaseStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICaseStatementContext); ok {
			tst[i] = t.(ICaseStatementContext)
			i++
		}
	}

	return tst
}

func (s *SwitchBodyStatementContext) CaseStatement(i int) ICaseStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICaseStatementContext); ok {
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

	return t.(ICaseStatementContext)
}

func (s *SwitchBodyStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SwitchBodyStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SwitchBodyStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSwitchBodyStatement(s)
	}
}

func (s *SwitchBodyStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSwitchBodyStatement(s)
	}
}

func NewEmptySwitchBodyStatementContext() *SwitchBodyStatementContext {
	var p = new(SwitchBodyStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_switchBodyStatement
	return p
}

func InitEmptySwitchBodyStatementContext(p *SwitchBodyStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_switchBodyStatement
}

func NewSwitchBodyStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchBodyStatementContext {
	var p = new(SwitchBodyStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_switchBodyStatement

	return p
}
