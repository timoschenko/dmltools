// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IIstemplateStatementContext is an interface to support dynamic dispatch.
type IIstemplateStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IstemplateExpr() IIstemplateExprContext
	SEMI() antlr.TerminalNode

	// IsIstemplateStatementContext differentiates from other interfaces.
	IsIstemplateStatementContext()
}

type IstemplateStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*IstemplateStatementContext) IsIstemplateStatementContext() {}

func (s *IstemplateStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *IstemplateStatementContext) IstemplateExpr() IIstemplateExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIstemplateExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIstemplateExprContext)
}

func (s *IstemplateStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *IstemplateStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IstemplateStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IstemplateStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterIstemplateStatement(s)
	}
}

func (s *IstemplateStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitIstemplateStatement(s)
	}
}

func NewEmptyIstemplateStatementContext() *IstemplateStatementContext {
	var p = new(IstemplateStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_istemplateStatement
	return p
}

func InitEmptyIstemplateStatementContext(p *IstemplateStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_istemplateStatement
}

func NewIstemplateStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IstemplateStatementContext {
	var p = new(IstemplateStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_istemplateStatement

	return p
}
