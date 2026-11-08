// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IErrorStatementContext is an interface to support dynamic dispatch.
type IErrorStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ERROR() antlr.TerminalNode
	SEMI() antlr.TerminalNode
	Bracketed_string_literal() IBracketed_string_literalContext

	// IsErrorStatementContext differentiates from other interfaces.
	IsErrorStatementContext()
}

type ErrorStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ErrorStatementContext) IsErrorStatementContext() {}

func (s *ErrorStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ErrorStatementContext) ERROR() antlr.TerminalNode {
	return s.GetToken(dmlParserERROR, 0)
}

func (s *ErrorStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *ErrorStatementContext) Bracketed_string_literal() IBracketed_string_literalContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBracketed_string_literalContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBracketed_string_literalContext)
}

func (s *ErrorStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ErrorStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ErrorStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterErrorStatement(s)
	}
}

func (s *ErrorStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitErrorStatement(s)
	}
}

func NewEmptyErrorStatementContext() *ErrorStatementContext {
	var p = new(ErrorStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_errorStatement
	return p
}

func InitEmptyErrorStatementContext(p *ErrorStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_errorStatement
}

func NewErrorStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ErrorStatementContext {
	var p = new(ErrorStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_errorStatement

	return p
}
