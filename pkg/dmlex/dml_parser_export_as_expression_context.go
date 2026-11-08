// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IExportAsExpressionContext is an interface to support dynamic dispatch.
type IExportAsExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext

	// IsExportAsExpressionContext differentiates from other interfaces.
	IsExportAsExpressionContext()
}

type ExportAsExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ExportAsExpressionContext) IsExportAsExpressionContext() {}

func (s *ExportAsExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExportAsExpressionContext) Expression() IExpressionContext {
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

func (s *ExportAsExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExportAsExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExportAsExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterExportAsExpression(s)
	}
}

func (s *ExportAsExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitExportAsExpression(s)
	}
}

func NewEmptyExportAsExpressionContext() *ExportAsExpressionContext {
	var p = new(ExportAsExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_exportAsExpression
	return p
}

func InitEmptyExportAsExpressionContext(p *ExportAsExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_exportAsExpression
}

func NewExportAsExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExportAsExpressionContext {
	var p = new(ExportAsExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_exportAsExpression

	return p
}
