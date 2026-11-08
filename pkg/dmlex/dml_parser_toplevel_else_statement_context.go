// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IToplevelElseStatementContext is an interface to support dynamic dispatch.
type IToplevelElseStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	DeviceStatements() IDeviceStatementsContext
	RBRACE() antlr.TerminalNode
	ToplevelIf() IToplevelIfContext

	// IsToplevelElseStatementContext differentiates from other interfaces.
	IsToplevelElseStatementContext()
}

type ToplevelElseStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ToplevelElseStatementContext) IsToplevelElseStatementContext() {}

func (s *ToplevelElseStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ToplevelElseStatementContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *ToplevelElseStatementContext) DeviceStatements() IDeviceStatementsContext {
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

func (s *ToplevelElseStatementContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *ToplevelElseStatementContext) ToplevelIf() IToplevelIfContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IToplevelIfContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IToplevelIfContext)
}

func (s *ToplevelElseStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToplevelElseStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ToplevelElseStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterToplevelElseStatement(s)
	}
}

func (s *ToplevelElseStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitToplevelElseStatement(s)
	}
}

func NewEmptyToplevelElseStatementContext() *ToplevelElseStatementContext {
	var p = new(ToplevelElseStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_toplevelElseStatement
	return p
}

func InitEmptyToplevelElseStatementContext(p *ToplevelElseStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_toplevelElseStatement
}

func NewToplevelElseStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ToplevelElseStatementContext {
	var p = new(ToplevelElseStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_toplevelElseStatement

	return p
}
