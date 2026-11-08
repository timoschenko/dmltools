// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISessionOrSavedDeclContext is an interface to support dynamic dispatch.
type ISessionOrSavedDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Global_decl_kind() IGlobal_decl_kindContext
	CvarDecl() ICvarDeclContext

	// IsSessionOrSavedDeclContext differentiates from other interfaces.
	IsSessionOrSavedDeclContext()
}

type SessionOrSavedDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*SessionOrSavedDeclContext) IsSessionOrSavedDeclContext() {}

func (s *SessionOrSavedDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *SessionOrSavedDeclContext) Global_decl_kind() IGlobal_decl_kindContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGlobal_decl_kindContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGlobal_decl_kindContext)
}

func (s *SessionOrSavedDeclContext) CvarDecl() ICvarDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICvarDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICvarDeclContext)
}

func (s *SessionOrSavedDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SessionOrSavedDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SessionOrSavedDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSessionOrSavedDecl(s)
	}
}

func (s *SessionOrSavedDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSessionOrSavedDecl(s)
	}
}

func NewEmptySessionOrSavedDeclContext() *SessionOrSavedDeclContext {
	var p = new(SessionOrSavedDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_sessionOrSavedDecl
	return p
}

func InitEmptySessionOrSavedDeclContext(p *SessionOrSavedDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_sessionOrSavedDecl
}

func NewSessionOrSavedDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SessionOrSavedDeclContext {
	var p = new(SessionOrSavedDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_sessionOrSavedDecl

	return p
}
