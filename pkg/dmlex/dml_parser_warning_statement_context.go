// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IWarningStatementContext is an interface to support dynamic dispatch.
type IWarningStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	WARNING() antlr.TerminalNode
	Bracketed_string_literal() IBracketed_string_literalContext
	SEMI() antlr.TerminalNode

	// IsWarningStatementContext differentiates from other interfaces.
	IsWarningStatementContext()
}

type WarningStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*WarningStatementContext) IsWarningStatementContext() {}

func (s *WarningStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *WarningStatementContext) WARNING() antlr.TerminalNode {
	return s.GetToken(dmlParserWARNING, 0)
}

func (s *WarningStatementContext) Bracketed_string_literal() IBracketed_string_literalContext {
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

func (s *WarningStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *WarningStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WarningStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WarningStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterWarningStatement(s)
	}
}

func (s *WarningStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitWarningStatement(s)
	}
}

func NewEmptyWarningStatementContext() *WarningStatementContext {
	var p = new(WarningStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_warningStatement
	return p
}

func InitEmptyWarningStatementContext(p *WarningStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_warningStatement
}

func NewWarningStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WarningStatementContext {
	var p = new(WarningStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_warningStatement

	return p
}
