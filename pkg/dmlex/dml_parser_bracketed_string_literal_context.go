// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IBracketed_string_literalContext is an interface to support dynamic dispatch.
type IBracketed_string_literalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Composed_string_literal() IComposed_string_literalContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode

	// IsBracketed_string_literalContext differentiates from other interfaces.
	IsBracketed_string_literalContext()
}

type Bracketed_string_literalContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Bracketed_string_literalContext) IsBracketed_string_literalContext() {}

func (s *Bracketed_string_literalContext) GetParser() antlr.Parser { return s.parser }

func (s *Bracketed_string_literalContext) Composed_string_literal() IComposed_string_literalContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComposed_string_literalContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComposed_string_literalContext)
}

func (s *Bracketed_string_literalContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *Bracketed_string_literalContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *Bracketed_string_literalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Bracketed_string_literalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Bracketed_string_literalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterBracketed_string_literal(s)
	}
}

func (s *Bracketed_string_literalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitBracketed_string_literal(s)
	}
}

func NewEmptyBracketed_string_literalContext() *Bracketed_string_literalContext {
	var p = new(Bracketed_string_literalContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bracketed_string_literal
	return p
}

func InitEmptyBracketed_string_literalContext(p *Bracketed_string_literalContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bracketed_string_literal
}

func NewBracketed_string_literalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Bracketed_string_literalContext {
	var p = new(Bracketed_string_literalContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_bracketed_string_literal

	return p
}
