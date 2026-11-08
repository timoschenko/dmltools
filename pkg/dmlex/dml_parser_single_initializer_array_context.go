// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISingleInitializerArrayContext is an interface to support dynamic dispatch.
type ISingleInitializerArrayContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSingleInitializer() []ISingleInitializerContext
	SingleInitializer(i int) ISingleInitializerContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsSingleInitializerArrayContext differentiates from other interfaces.
	IsSingleInitializerArrayContext()
}

type SingleInitializerArrayContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*SingleInitializerArrayContext) IsSingleInitializerArrayContext() {}

func (s *SingleInitializerArrayContext) GetParser() antlr.Parser { return s.parser }

func (s *SingleInitializerArrayContext) AllSingleInitializer() []ISingleInitializerContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISingleInitializerContext); ok {
			len++
		}
	}

	tst := make([]ISingleInitializerContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISingleInitializerContext); ok {
			tst[i] = t.(ISingleInitializerContext)
			i++
		}
	}

	return tst
}

func (s *SingleInitializerArrayContext) SingleInitializer(i int) ISingleInitializerContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleInitializerContext); ok {
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

	return t.(ISingleInitializerContext)
}

func (s *SingleInitializerArrayContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCOMMA)
}

func (s *SingleInitializerArrayContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, i)
}

func (s *SingleInitializerArrayContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SingleInitializerArrayContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SingleInitializerArrayContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSingleInitializerArray(s)
	}
}

func (s *SingleInitializerArrayContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSingleInitializerArray(s)
	}
}

func NewEmptySingleInitializerArrayContext() *SingleInitializerArrayContext {
	var p = new(SingleInitializerArrayContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_singleInitializerArray
	return p
}

func InitEmptySingleInitializerArrayContext(p *SingleInitializerArrayContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_singleInitializerArray
}

func NewSingleInitializerArrayContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SingleInitializerArrayContext {
	var p = new(SingleInitializerArrayContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_singleInitializerArray

	return p
}
