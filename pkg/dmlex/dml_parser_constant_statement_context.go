// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IConstantStatementContext is an interface to support dynamic dispatch.
type IConstantStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CONSTANT() antlr.TerminalNode
	Ident() IIdentContext
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext
	SEMI() antlr.TerminalNode

	// IsConstantStatementContext differentiates from other interfaces.
	IsConstantStatementContext()
}

type ConstantStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ConstantStatementContext) IsConstantStatementContext() {}

func (s *ConstantStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ConstantStatementContext) CONSTANT() antlr.TerminalNode {
	return s.GetToken(dmlParserCONSTANT, 0)
}

func (s *ConstantStatementContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *ConstantStatementContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserASSIGN, 0)
}

func (s *ConstantStatementContext) Expression() IExpressionContext {
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

func (s *ConstantStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *ConstantStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstantStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConstantStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterConstantStatement(s)
	}
}

func (s *ConstantStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitConstantStatement(s)
	}
}

func NewEmptyConstantStatementContext() *ConstantStatementContext {
	var p = new(ConstantStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_constantStatement
	return p
}

func InitEmptyConstantStatementContext(p *ConstantStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_constantStatement
}

func NewConstantStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstantStatementContext {
	var p = new(ConstantStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_constantStatement

	return p
}
