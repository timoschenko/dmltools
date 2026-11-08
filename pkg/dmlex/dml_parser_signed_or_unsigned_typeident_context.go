// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISigned_or_unsigned_typeidentContext is an interface to support dynamic dispatch.
type ISigned_or_unsigned_typeidentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Typeident() ITypeidentContext
	UNSIGNED() antlr.TerminalNode
	SIGNED() antlr.TerminalNode

	// IsSigned_or_unsigned_typeidentContext differentiates from other interfaces.
	IsSigned_or_unsigned_typeidentContext()
}

type Signed_or_unsigned_typeidentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Signed_or_unsigned_typeidentContext) IsSigned_or_unsigned_typeidentContext() {}

func (s *Signed_or_unsigned_typeidentContext) GetParser() antlr.Parser { return s.parser }

func (s *Signed_or_unsigned_typeidentContext) Typeident() ITypeidentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeidentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeidentContext)
}

func (s *Signed_or_unsigned_typeidentContext) UNSIGNED() antlr.TerminalNode {
	return s.GetToken(dmlParserUNSIGNED, 0)
}

func (s *Signed_or_unsigned_typeidentContext) SIGNED() antlr.TerminalNode {
	return s.GetToken(dmlParserSIGNED, 0)
}

func (s *Signed_or_unsigned_typeidentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Signed_or_unsigned_typeidentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Signed_or_unsigned_typeidentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSigned_or_unsigned_typeident(s)
	}
}

func (s *Signed_or_unsigned_typeidentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSigned_or_unsigned_typeident(s)
	}
}

func NewEmptySigned_or_unsigned_typeidentContext() *Signed_or_unsigned_typeidentContext {
	var p = new(Signed_or_unsigned_typeidentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_signed_or_unsigned_typeident
	return p
}

func InitEmptySigned_or_unsigned_typeidentContext(p *Signed_or_unsigned_typeidentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_signed_or_unsigned_typeident
}

func NewSigned_or_unsigned_typeidentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Signed_or_unsigned_typeidentContext {
	var p = new(Signed_or_unsigned_typeidentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_signed_or_unsigned_typeident

	return p
}
