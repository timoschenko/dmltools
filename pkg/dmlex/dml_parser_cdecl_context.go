// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICdeclContext is an interface to support dynamic dispatch.
type ICdeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Basetype() IBasetypeContext
	Cdecl_qualifier_ident_array_func() ICdecl_qualifier_ident_array_funcContext
	CONST() antlr.TerminalNode

	// IsCdeclContext differentiates from other interfaces.
	IsCdeclContext()
}

type CdeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*CdeclContext) IsCdeclContext() {}

func (s *CdeclContext) GetParser() antlr.Parser { return s.parser }

func (s *CdeclContext) Basetype() IBasetypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBasetypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBasetypeContext)
}

func (s *CdeclContext) Cdecl_qualifier_ident_array_func() ICdecl_qualifier_ident_array_funcContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdecl_qualifier_ident_array_funcContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdecl_qualifier_ident_array_funcContext)
}

func (s *CdeclContext) CONST() antlr.TerminalNode {
	return s.GetToken(dmlParserCONST, 0)
}

func (s *CdeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CdeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CdeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdecl(s)
	}
}

func (s *CdeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdecl(s)
	}
}

func NewEmptyCdeclContext() *CdeclContext {
	var p = new(CdeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl
	return p
}

func InitEmptyCdeclContext(p *CdeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl
}

func NewCdeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CdeclContext {
	var p = new(CdeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_cdecl

	return p
}
