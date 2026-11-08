// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISizeSpecContext is an interface to support dynamic dispatch.
type ISizeSpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext

	// IsSizeSpecContext differentiates from other interfaces.
	IsSizeSpecContext()
}

type SizeSpecContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*SizeSpecContext) IsSizeSpecContext() {}

func (s *SizeSpecContext) GetParser() antlr.Parser { return s.parser }

func (s *SizeSpecContext) Expression() IExpressionContext {
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

func (s *SizeSpecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SizeSpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SizeSpecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSizeSpec(s)
	}
}

func (s *SizeSpecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSizeSpec(s)
	}
}

func NewEmptySizeSpecContext() *SizeSpecContext {
	var p = new(SizeSpecContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_sizeSpec
	return p
}

func InitEmptySizeSpecContext(p *SizeSpecContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_sizeSpec
}

func NewSizeSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SizeSpecContext {
	var p = new(SizeSpecContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_sizeSpec

	return p
}
