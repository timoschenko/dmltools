// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IExportStatementContext is an interface to support dynamic dispatch.
type IExportStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXPORT() antlr.TerminalNode
	Expression() IExpressionContext
	AS() antlr.TerminalNode
	ExportAsExpression() IExportAsExpressionContext
	SEMI() antlr.TerminalNode

	// IsExportStatementContext differentiates from other interfaces.
	IsExportStatementContext()
}

type ExportStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ExportStatementContext) IsExportStatementContext() {}

func (s *ExportStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ExportStatementContext) EXPORT() antlr.TerminalNode {
	return s.GetToken(dmlParserEXPORT, 0)
}

func (s *ExportStatementContext) Expression() IExpressionContext {
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

func (s *ExportStatementContext) AS() antlr.TerminalNode {
	return s.GetToken(dmlParserAS, 0)
}

func (s *ExportStatementContext) ExportAsExpression() IExportAsExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExportAsExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExportAsExpressionContext)
}

func (s *ExportStatementContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *ExportStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExportStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExportStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterExportStatement(s)
	}
}

func (s *ExportStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitExportStatement(s)
	}
}

func NewEmptyExportStatementContext() *ExportStatementContext {
	var p = new(ExportStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_exportStatement
	return p
}

func InitEmptyExportStatementContext(p *ExportStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_exportStatement
}

func NewExportStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExportStatementContext {
	var p = new(ExportStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_exportStatement

	return p
}
