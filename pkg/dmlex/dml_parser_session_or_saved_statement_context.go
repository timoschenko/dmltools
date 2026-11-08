// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISessionOrSavedStatementContext is an interface to support dynamic dispatch.
type ISessionOrSavedStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SessionOrSavedDecl() ISessionOrSavedDeclContext
	SEMI() antlr.TerminalNode

	// IsSessionOrSavedStatementContext differentiates from other interfaces.
	IsSessionOrSavedStatementContext()
}

type SessionOrSavedStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*SessionOrSavedStatementContext) IsSessionOrSavedStatementContext() {}

func (s *SessionOrSavedStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *SessionOrSavedStatementContext) SessionOrSavedDecl() ISessionOrSavedDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISessionOrSavedDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISessionOrSavedDeclContext)
}

func (s *SessionOrSavedStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *SessionOrSavedStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SessionOrSavedStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SessionOrSavedStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSessionOrSavedStatement(s)
	}
}

func (s *SessionOrSavedStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSessionOrSavedStatement(s)
	}
}

func NewEmptySessionOrSavedStatementContext() *SessionOrSavedStatementContext {
	var p = new(SessionOrSavedStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_sessionOrSavedStatement
	return p
}

func InitEmptySessionOrSavedStatementContext(p *SessionOrSavedStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_sessionOrSavedStatement
}

func NewSessionOrSavedStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SessionOrSavedStatementContext {
	var p = new(SessionOrSavedStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_sessionOrSavedStatement

	return p
}
