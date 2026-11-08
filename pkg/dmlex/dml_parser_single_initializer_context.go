// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISingleInitializerContext is an interface to support dynamic dispatch.
type ISingleInitializerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	CdeclInitializer() ICdeclInitializerContext

	// IsSingleInitializerContext differentiates from other interfaces.
	IsSingleInitializerContext()
}

type SingleInitializerContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*SingleInitializerContext) IsSingleInitializerContext() {}

func (s *SingleInitializerContext) GetParser() antlr.Parser { return s.parser }

func (s *SingleInitializerContext) Expression() IExpressionContext {
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

func (s *SingleInitializerContext) CdeclInitializer() ICdeclInitializerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdeclInitializerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdeclInitializerContext)
}

func (s *SingleInitializerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SingleInitializerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SingleInitializerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSingleInitializer(s)
	}
}

func (s *SingleInitializerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSingleInitializer(s)
	}
}

func NewEmptySingleInitializerContext() *SingleInitializerContext {
	var p = new(SingleInitializerContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_singleInitializer
	return p
}

func InitEmptySingleInitializerContext(p *SingleInitializerContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_singleInitializer
}

func NewSingleInitializerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SingleInitializerContext {
	var p = new(SingleInitializerContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_singleInitializer

	return p
}
