// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IArgumentsListEmptyContext is an interface to support dynamic dispatch.
type IArgumentsListEmptyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode

	// IsArgumentsListEmptyContext differentiates from other interfaces.
	IsArgumentsListEmptyContext()
}

type ArgumentsListEmptyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ArgumentsListEmptyContext) IsArgumentsListEmptyContext() {}

func (s *ArgumentsListEmptyContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgumentsListEmptyContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *ArgumentsListEmptyContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *ArgumentsListEmptyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgumentsListEmptyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgumentsListEmptyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterArgumentsListEmpty(s)
	}
}

func (s *ArgumentsListEmptyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitArgumentsListEmpty(s)
	}
}

func NewEmptyArgumentsListEmptyContext() *ArgumentsListEmptyContext {
	var p = new(ArgumentsListEmptyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_argumentsListEmpty
	return p
}

func InitEmptyArgumentsListEmptyContext(p *ArgumentsListEmptyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_argumentsListEmpty
}

func NewArgumentsListEmptyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentsListEmptyContext {
	var p = new(ArgumentsListEmptyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_argumentsListEmpty

	return p
}
