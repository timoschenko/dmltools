// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IEndianflagContext is an interface to support dynamic dispatch.
type IEndianflagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COMMA() antlr.TerminalNode
	Ident() IIdentContext

	// IsEndianflagContext differentiates from other interfaces.
	IsEndianflagContext()
}

type EndianflagContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*EndianflagContext) IsEndianflagContext() {}

func (s *EndianflagContext) GetParser() antlr.Parser { return s.parser }

func (s *EndianflagContext) COMMA() antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, 0)
}

func (s *EndianflagContext) Ident() IIdentContext {
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

func (s *EndianflagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EndianflagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EndianflagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterEndianflag(s)
	}
}

func (s *EndianflagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitEndianflag(s)
	}
}

func NewEmptyEndianflagContext() *EndianflagContext {
	var p = new(EndianflagContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_endianflag
	return p
}

func InitEmptyEndianflagContext(p *EndianflagContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_endianflag
}

func NewEndianflagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EndianflagContext {
	var p = new(EndianflagContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_endianflag

	return p
}
