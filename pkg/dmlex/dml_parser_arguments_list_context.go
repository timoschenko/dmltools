// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IArgumentsListContext is an interface to support dynamic dispatch.
type IArgumentsListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	SingleInitializerArray() ISingleInitializerArrayContext
	RPAREN() antlr.TerminalNode

	// IsArgumentsListContext differentiates from other interfaces.
	IsArgumentsListContext()
}

type ArgumentsListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ArgumentsListContext) IsArgumentsListContext() {}

func (s *ArgumentsListContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgumentsListContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *ArgumentsListContext) SingleInitializerArray() ISingleInitializerArrayContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleInitializerArrayContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISingleInitializerArrayContext)
}

func (s *ArgumentsListContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *ArgumentsListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgumentsListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgumentsListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterArgumentsList(s)
	}
}

func (s *ArgumentsListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitArgumentsList(s)
	}
}

func NewEmptyArgumentsListContext() *ArgumentsListContext {
	var p = new(ArgumentsListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_argumentsList
	return p
}

func InitEmptyArgumentsListContext(p *ArgumentsListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_argumentsList
}

func NewArgumentsListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentsListContext {
	var p = new(ArgumentsListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_argumentsList

	return p
}
