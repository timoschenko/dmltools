// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IComposed_string_literalContext is an interface to support dynamic dispatch.
type IComposed_string_literalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllUtf8_sconst() []IUtf8_sconstContext
	Utf8_sconst(i int) IUtf8_sconstContext
	AllADD() []antlr.TerminalNode
	ADD(i int) antlr.TerminalNode

	// IsComposed_string_literalContext differentiates from other interfaces.
	IsComposed_string_literalContext()
}

type Composed_string_literalContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Composed_string_literalContext) IsComposed_string_literalContext() {}

func (s *Composed_string_literalContext) GetParser() antlr.Parser { return s.parser }

func (s *Composed_string_literalContext) AllUtf8_sconst() []IUtf8_sconstContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IUtf8_sconstContext); ok {
			len++
		}
	}

	tst := make([]IUtf8_sconstContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IUtf8_sconstContext); ok {
			tst[i] = t.(IUtf8_sconstContext)
			i++
		}
	}

	return tst
}

func (s *Composed_string_literalContext) Utf8_sconst(i int) IUtf8_sconstContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUtf8_sconstContext); ok {
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

	return t.(IUtf8_sconstContext)
}

func (s *Composed_string_literalContext) AllADD() []antlr.TerminalNode {
	return s.GetTokens(dmlParserADD)
}

func (s *Composed_string_literalContext) ADD(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserADD, i)
}

func (s *Composed_string_literalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Composed_string_literalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Composed_string_literalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterComposed_string_literal(s)
	}
}

func (s *Composed_string_literalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitComposed_string_literal(s)
	}
}

func NewEmptyComposed_string_literalContext() *Composed_string_literalContext {
	var p = new(Composed_string_literalContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_composed_string_literal
	return p
}

func InitEmptyComposed_string_literalContext(p *Composed_string_literalContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_composed_string_literal
}

func NewComposed_string_literalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Composed_string_literalContext {
	var p = new(Composed_string_literalContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_composed_string_literal

	return p
}
