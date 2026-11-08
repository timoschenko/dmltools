// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICaseOrDefaultStatementContext is an interface to support dynamic dispatch.
type ICaseOrDefaultStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	Expression() IExpressionContext
	DEFAULT() antlr.TerminalNode

	// IsCaseOrDefaultStatementContext differentiates from other interfaces.
	IsCaseOrDefaultStatementContext()
}

type CaseOrDefaultStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*CaseOrDefaultStatementContext) IsCaseOrDefaultStatementContext() {}

func (s *CaseOrDefaultStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *CaseOrDefaultStatementContext) CASE() antlr.TerminalNode {
	return s.GetToken(dmlParserCASE, 0)
}

func (s *CaseOrDefaultStatementContext) Expression() IExpressionContext {
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

func (s *CaseOrDefaultStatementContext) DEFAULT() antlr.TerminalNode {
	return s.GetToken(dmlParserDEFAULT, 0)
}

func (s *CaseOrDefaultStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CaseOrDefaultStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CaseOrDefaultStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCaseOrDefaultStatement(s)
	}
}

func (s *CaseOrDefaultStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCaseOrDefaultStatement(s)
	}
}

func NewEmptyCaseOrDefaultStatementContext() *CaseOrDefaultStatementContext {
	var p = new(CaseOrDefaultStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_caseOrDefaultStatement
	return p
}

func InitEmptyCaseOrDefaultStatementContext(p *CaseOrDefaultStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_caseOrDefaultStatement
}

func NewCaseOrDefaultStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaseOrDefaultStatementContext {
	var p = new(CaseOrDefaultStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_caseOrDefaultStatement

	return p
}
