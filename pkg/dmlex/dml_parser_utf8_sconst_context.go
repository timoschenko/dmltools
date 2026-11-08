// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IUtf8_sconstContext is an interface to support dynamic dispatch.
type IUtf8_sconstContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING_LITERAL() antlr.TerminalNode

	// IsUtf8_sconstContext differentiates from other interfaces.
	IsUtf8_sconstContext()
}

type Utf8_sconstContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Utf8_sconstContext) IsUtf8_sconstContext() {}

func (s *Utf8_sconstContext) GetParser() antlr.Parser { return s.parser }

func (s *Utf8_sconstContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(dmlParserSTRING_LITERAL, 0)
}

func (s *Utf8_sconstContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Utf8_sconstContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Utf8_sconstContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterUtf8_sconst(s)
	}
}

func (s *Utf8_sconstContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitUtf8_sconst(s)
	}
}

func NewEmptyUtf8_sconstContext() *Utf8_sconstContext {
	var p = new(Utf8_sconstContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_utf8_sconst
	return p
}

func InitEmptyUtf8_sconstContext(p *Utf8_sconstContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_utf8_sconst
}

func NewUtf8_sconstContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Utf8_sconstContext {
	var p = new(Utf8_sconstContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_utf8_sconst

	return p
}
