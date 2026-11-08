// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IObjidentContext is an interface to support dynamic dispatch.
type IObjidentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	REGISTER() antlr.TerminalNode
	Ident() IIdentContext

	// IsObjidentContext differentiates from other interfaces.
	IsObjidentContext()
}

type ObjidentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ObjidentContext) IsObjidentContext() {}

func (s *ObjidentContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjidentContext) REGISTER() antlr.TerminalNode {
	return s.GetToken(dmlParserREGISTER, 0)
}

func (s *ObjidentContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ObjidentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjidentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjidentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterObjident(s)
	}
}

func (s *ObjidentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitObjident(s)
	}
}

func NewEmptyObjidentContext() *ObjidentContext {
	var p = new(ObjidentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objident
	return p
}

func InitEmptyObjidentContext(p *ObjidentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objident
}

func NewObjidentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjidentContext {
	var p = new(ObjidentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_objident

	return p
}
