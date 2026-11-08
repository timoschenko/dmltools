// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IThenLevelContext is an interface to support dynamic dispatch.
type IThenLevelContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext

	// IsThenLevelContext differentiates from other interfaces.
	IsThenLevelContext()
}

type ThenLevelContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ThenLevelContext) IsThenLevelContext() {}

func (s *ThenLevelContext) GetParser() antlr.Parser { return s.parser }

func (s *ThenLevelContext) Expression() IExpressionContext {
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

func (s *ThenLevelContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ThenLevelContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ThenLevelContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterThenLevel(s)
	}
}

func (s *ThenLevelContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitThenLevel(s)
	}
}

func NewEmptyThenLevelContext() *ThenLevelContext {
	var p = new(ThenLevelContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_thenLevel
	return p
}

func InitEmptyThenLevelContext(p *ThenLevelContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_thenLevel
}

func NewThenLevelContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ThenLevelContext {
	var p = new(ThenLevelContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_thenLevel

	return p
}
