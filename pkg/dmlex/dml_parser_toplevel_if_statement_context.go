// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IToplevelIfStatementContext is an interface to support dynamic dispatch.
type IToplevelIfStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	DeviceStatements() IDeviceStatementsContext
	RBRACE() antlr.TerminalNode

	// IsToplevelIfStatementContext differentiates from other interfaces.
	IsToplevelIfStatementContext()
}

type ToplevelIfStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ToplevelIfStatementContext) IsToplevelIfStatementContext() {}

func (s *ToplevelIfStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ToplevelIfStatementContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *ToplevelIfStatementContext) DeviceStatements() IDeviceStatementsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeviceStatementsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeviceStatementsContext)
}

func (s *ToplevelIfStatementContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *ToplevelIfStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToplevelIfStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ToplevelIfStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterToplevelIfStatement(s)
	}
}

func (s *ToplevelIfStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitToplevelIfStatement(s)
	}
}

func NewEmptyToplevelIfStatementContext() *ToplevelIfStatementContext {
	var p = new(ToplevelIfStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_toplevelIfStatement
	return p
}

func InitEmptyToplevelIfStatementContext(p *ToplevelIfStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_toplevelIfStatement
}

func NewToplevelIfStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ToplevelIfStatementContext {
	var p = new(ToplevelIfStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_toplevelIfStatement

	return p
}
