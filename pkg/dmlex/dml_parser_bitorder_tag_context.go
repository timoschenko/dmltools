// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IBitorderTagContext is an interface to support dynamic dispatch.
type IBitorderTagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BITORDER() antlr.TerminalNode
	ENDIAN() antlr.TerminalNode
	SEMI() antlr.TerminalNode

	// IsBitorderTagContext differentiates from other interfaces.
	IsBitorderTagContext()
}

type BitorderTagContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*BitorderTagContext) IsBitorderTagContext() {}

func (s *BitorderTagContext) GetParser() antlr.Parser { return s.parser }

func (s *BitorderTagContext) BITORDER() antlr.TerminalNode {
	return s.GetToken(dmlParserBITORDER, 0)
}

func (s *BitorderTagContext) ENDIAN() antlr.TerminalNode {
	return s.GetToken(dmlParserENDIAN, 0)
}

func (s *BitorderTagContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *BitorderTagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BitorderTagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BitorderTagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterBitorderTag(s)
	}
}

func (s *BitorderTagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitBitorderTag(s)
	}
}

func NewEmptyBitorderTagContext() *BitorderTagContext {
	var p = new(BitorderTagContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitorderTag
	return p
}

func InitEmptyBitorderTagContext(p *BitorderTagContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitorderTag
}

func NewBitorderTagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BitorderTagContext {
	var p = new(BitorderTagContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_bitorderTag

	return p
}
