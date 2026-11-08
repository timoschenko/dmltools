// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IBitfieldsContext is an interface to support dynamic dispatch.
type IBitfieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BITFIELDS() antlr.TerminalNode
	INTEGER_LITERAL() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllBitfields_decls() []IBitfields_declsContext
	Bitfields_decls(i int) IBitfields_declsContext

	// IsBitfieldsContext differentiates from other interfaces.
	IsBitfieldsContext()
}

type BitfieldsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*BitfieldsContext) IsBitfieldsContext() {}

func (s *BitfieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *BitfieldsContext) BITFIELDS() antlr.TerminalNode {
	return s.GetToken(dmlParserBITFIELDS, 0)
}

func (s *BitfieldsContext) INTEGER_LITERAL() antlr.TerminalNode {
	return s.GetToken(dmlParserINTEGER_LITERAL, 0)
}

func (s *BitfieldsContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *BitfieldsContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *BitfieldsContext) AllBitfields_decls() []IBitfields_declsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IBitfields_declsContext); ok {
			len++
		}
	}

	tst := make([]IBitfields_declsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IBitfields_declsContext); ok {
			tst[i] = t.(IBitfields_declsContext)
			i++
		}
	}

	return tst
}

func (s *BitfieldsContext) Bitfields_decls(i int) IBitfields_declsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBitfields_declsContext); ok {
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

	return t.(IBitfields_declsContext)
}

func (s *BitfieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BitfieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BitfieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterBitfields(s)
	}
}

func (s *BitfieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitBitfields(s)
	}
}

func NewEmptyBitfieldsContext() *BitfieldsContext {
	var p = new(BitfieldsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitfields
	return p
}

func InitEmptyBitfieldsContext(p *BitfieldsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitfields
}

func NewBitfieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BitfieldsContext {
	var p = new(BitfieldsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_bitfields

	return p
}
