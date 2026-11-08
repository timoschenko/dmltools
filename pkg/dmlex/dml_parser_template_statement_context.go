// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ITemplateStatementContext is an interface to support dynamic dispatch.
type ITemplateStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TEMPLATE() antlr.TerminalNode
	Objident() IObjidentContext
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	IstemplateExpr() IIstemplateExprContext
	AllTemplateBody() []ITemplateBodyContext
	TemplateBody(i int) ITemplateBodyContext

	// IsTemplateStatementContext differentiates from other interfaces.
	IsTemplateStatementContext()
}

type TemplateStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*TemplateStatementContext) IsTemplateStatementContext() {}

func (s *TemplateStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *TemplateStatementContext) TEMPLATE() antlr.TerminalNode {
	return s.GetToken(dmlParserTEMPLATE, 0)
}

func (s *TemplateStatementContext) Objident() IObjidentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjidentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjidentContext)
}

func (s *TemplateStatementContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *TemplateStatementContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *TemplateStatementContext) IstemplateExpr() IIstemplateExprContext {
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

func (s *TemplateStatementContext) AllTemplateBody() []ITemplateBodyContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITemplateBodyContext); ok {
			len++
		}
	}

	tst := make([]ITemplateBodyContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITemplateBodyContext); ok {
			tst[i] = t.(ITemplateBodyContext)
			i++
		}
	}

	return tst
}

func (s *TemplateStatementContext) TemplateBody(i int) ITemplateBodyContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITemplateBodyContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITemplateBodyContext)
}

func (s *TemplateStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TemplateStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TemplateStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterTemplateStatement(s)
	}
}

func (s *TemplateStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitTemplateStatement(s)
	}
}

func NewEmptyTemplateStatementContext() *TemplateStatementContext {
	var p = new(TemplateStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_templateStatement
	return p
}

func InitEmptyTemplateStatementContext(p *TemplateStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_templateStatement
}

func NewTemplateStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TemplateStatementContext {
	var p = new(TemplateStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_templateStatement

	return p
}
