// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IDmlTagContext is an interface to support dynamic dispatch.
type IDmlTagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DML() antlr.TerminalNode
	FLOAT_LITERAL() antlr.TerminalNode
	SEMI() antlr.TerminalNode

	// IsDmlTagContext differentiates from other interfaces.
	IsDmlTagContext()
}

type DmlTagContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*DmlTagContext) IsDmlTagContext() {}

func (s *DmlTagContext) GetParser() antlr.Parser { return s.parser }

func (s *DmlTagContext) DML() antlr.TerminalNode {
	return s.GetToken(dmlParserDML, 0)
}

func (s *DmlTagContext) FLOAT_LITERAL() antlr.TerminalNode {
	return s.GetToken(dmlParserFLOAT_LITERAL, 0)
}

func (s *DmlTagContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *DmlTagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DmlTagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DmlTagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDmlTag(s)
	}
}

func (s *DmlTagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDmlTag(s)
	}
}

func NewEmptyDmlTagContext() *DmlTagContext {
	var p = new(DmlTagContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_dmlTag
	return p
}

func InitEmptyDmlTagContext(p *DmlTagContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_dmlTag
}

func NewDmlTagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DmlTagContext {
	var p = new(DmlTagContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_dmlTag

	return p
}
