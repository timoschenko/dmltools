// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IIstemplateExprContext is an interface to support dynamic dispatch.
type IIstemplateExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IS() antlr.TerminalNode
	IstemplateList() IIstemplateListContext

	// IsIstemplateExprContext differentiates from other interfaces.
	IsIstemplateExprContext()
}

type IstemplateExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*IstemplateExprContext) IsIstemplateExprContext() {}

func (s *IstemplateExprContext) GetParser() antlr.Parser { return s.parser }

func (s *IstemplateExprContext) IS() antlr.TerminalNode {
	return s.GetToken(dmlParserIS, 0)
}

func (s *IstemplateExprContext) IstemplateList() IIstemplateListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIstemplateListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIstemplateListContext)
}

func (s *IstemplateExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IstemplateExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IstemplateExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterIstemplateExpr(s)
	}
}

func (s *IstemplateExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitIstemplateExpr(s)
	}
}

func NewEmptyIstemplateExprContext() *IstemplateExprContext {
	var p = new(IstemplateExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_istemplateExpr
	return p
}

func InitEmptyIstemplateExprContext(p *IstemplateExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_istemplateExpr
}

func NewIstemplateExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IstemplateExprContext {
	var p = new(IstemplateExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_istemplateExpr

	return p
}
