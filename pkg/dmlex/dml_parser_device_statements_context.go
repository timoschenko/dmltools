// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IDeviceStatementsContext is an interface to support dynamic dispatch.
type IDeviceStatementsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllDeviceStatement() []IDeviceStatementContext
	DeviceStatement(i int) IDeviceStatementContext

	// IsDeviceStatementsContext differentiates from other interfaces.
	IsDeviceStatementsContext()
}

type DeviceStatementsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*DeviceStatementsContext) IsDeviceStatementsContext() {}

func (s *DeviceStatementsContext) GetParser() antlr.Parser { return s.parser }

func (s *DeviceStatementsContext) AllDeviceStatement() []IDeviceStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IDeviceStatementContext); ok {
			len++
		}
	}

	tst := make([]IDeviceStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IDeviceStatementContext); ok {
			tst[i] = t.(IDeviceStatementContext)
			i++
		}
	}

	return tst
}

func (s *DeviceStatementsContext) DeviceStatement(i int) IDeviceStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeviceStatementContext); ok {
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

	return t.(IDeviceStatementContext)
}

func (s *DeviceStatementsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeviceStatementsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeviceStatementsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDeviceStatements(s)
	}
}

func (s *DeviceStatementsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDeviceStatements(s)
	}
}

func NewEmptyDeviceStatementsContext() *DeviceStatementsContext {
	var p = new(DeviceStatementsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_deviceStatements
	return p
}

func InitEmptyDeviceStatementsContext(p *DeviceStatementsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_deviceStatements
}

func NewDeviceStatementsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeviceStatementsContext {
	var p = new(DeviceStatementsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_deviceStatements

	return p
}
