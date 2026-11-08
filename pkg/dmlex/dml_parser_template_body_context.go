// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ITemplateBodyContext is an interface to support dynamic dispatch.
type ITemplateBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ObjectStatement() IObjectStatementContext
	SHARED() antlr.TerminalNode
	Hook_decl() IHook_declContext

	// IsTemplateBodyContext differentiates from other interfaces.
	IsTemplateBodyContext()
}

type TemplateBodyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*TemplateBodyContext) IsTemplateBodyContext() {}

func (s *TemplateBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *TemplateBodyContext) ObjectStatement() IObjectStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectStatementContext)
}

func (s *TemplateBodyContext) SHARED() antlr.TerminalNode {
	return s.GetToken(dmlParserSHARED, 0)
}

func (s *TemplateBodyContext) Hook_decl() IHook_declContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHook_declContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHook_declContext)
}

func (s *TemplateBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TemplateBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TemplateBodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterTemplateBody(s)
	}
}

func (s *TemplateBodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitTemplateBody(s)
	}
}

func NewEmptyTemplateBodyContext() *TemplateBodyContext {
	var p = new(TemplateBodyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_templateBody
	return p
}

func InitEmptyTemplateBodyContext(p *TemplateBodyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_templateBody
}

func NewTemplateBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TemplateBodyContext {
	var p = new(TemplateBodyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_templateBody

	return p
}
