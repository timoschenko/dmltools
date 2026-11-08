// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IDmlContext is an interface to support dynamic dispatch.
type IDmlContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DmlCode() IDmlCodeContext
	EOF() antlr.TerminalNode
	DmlTag() IDmlTagContext

	// IsDmlContext differentiates from other interfaces.
	IsDmlContext()
}

type DmlContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*DmlContext) IsDmlContext() {}

func (s *DmlContext) GetParser() antlr.Parser { return s.parser }

func (s *DmlContext) DmlCode() IDmlCodeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDmlCodeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDmlCodeContext)
}

func (s *DmlContext) EOF() antlr.TerminalNode {
	return s.GetToken(dmlParserEOF, 0)
}

func (s *DmlContext) DmlTag() IDmlTagContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDmlTagContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDmlTagContext)
}

func (s *DmlContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DmlContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DmlContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDml(s)
	}
}

func (s *DmlContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDml(s)
	}
}

func NewEmptyDmlContext() *DmlContext {
	var p = new(DmlContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_dml
	return p
}

func InitEmptyDmlContext(p *DmlContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_dml
}

func NewDmlContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DmlContext {
	var p = new(DmlContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_dml

	return p
}
