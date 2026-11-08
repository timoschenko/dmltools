// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IOffsetSpecContext is an interface to support dynamic dispatch.
type IOffsetSpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AT() antlr.TerminalNode
	Expression() IExpressionContext

	// IsOffsetSpecContext differentiates from other interfaces.
	IsOffsetSpecContext()
}

type OffsetSpecContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*OffsetSpecContext) IsOffsetSpecContext() {}

func (s *OffsetSpecContext) GetParser() antlr.Parser { return s.parser }

func (s *OffsetSpecContext) AT() antlr.TerminalNode {
	return s.GetToken(dmlParserAT, 0)
}

func (s *OffsetSpecContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *OffsetSpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OffsetSpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OffsetSpecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterOffsetSpec(s)
	}
}

func (s *OffsetSpecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitOffsetSpec(s)
	}
}

func NewEmptyOffsetSpecContext() *OffsetSpecContext {
	var p = new(OffsetSpecContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_offsetSpec
	return p
}

func InitEmptyOffsetSpecContext(p *OffsetSpecContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_offsetSpec
}

func NewOffsetSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OffsetSpecContext {
	var p = new(OffsetSpecContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_offsetSpec

	return p
}
