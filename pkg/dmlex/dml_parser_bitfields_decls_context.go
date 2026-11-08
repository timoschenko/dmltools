// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IBitfields_declsContext is an interface to support dynamic dispatch.
type IBitfields_declsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Cdecl() ICdeclContext
	BitrangeDecl() IBitrangeDeclContext
	SEMI() antlr.TerminalNode

	// IsBitfields_declsContext differentiates from other interfaces.
	IsBitfields_declsContext()
}

type Bitfields_declsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Bitfields_declsContext) IsBitfields_declsContext() {}

func (s *Bitfields_declsContext) GetParser() antlr.Parser { return s.parser }

func (s *Bitfields_declsContext) Cdecl() ICdeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdeclContext)
}

func (s *Bitfields_declsContext) BitrangeDecl() IBitrangeDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBitrangeDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBitrangeDeclContext)
}

func (s *Bitfields_declsContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *Bitfields_declsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Bitfields_declsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Bitfields_declsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterBitfields_decls(s)
	}
}

func (s *Bitfields_declsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitBitfields_decls(s)
	}
}

func NewEmptyBitfields_declsContext() *Bitfields_declsContext {
	var p = new(Bitfields_declsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitfields_decls
	return p
}

func InitEmptyBitfields_declsContext(p *Bitfields_declsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitfields_decls
}

func NewBitfields_declsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Bitfields_declsContext {
	var p = new(Bitfields_declsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_bitfields_decls

	return p
}
