// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ILayoutDeclsContext is an interface to support dynamic dispatch.
type ILayoutDeclsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Cdecl() ICdeclContext
	SEMI() antlr.TerminalNode

	// IsLayoutDeclsContext differentiates from other interfaces.
	IsLayoutDeclsContext()
}

type LayoutDeclsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*LayoutDeclsContext) IsLayoutDeclsContext() {}

func (s *LayoutDeclsContext) GetParser() antlr.Parser { return s.parser }

func (s *LayoutDeclsContext) Cdecl() ICdeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdeclContext)
}

func (s *LayoutDeclsContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *LayoutDeclsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LayoutDeclsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LayoutDeclsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLayoutDecls(s)
	}
}

func (s *LayoutDeclsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLayoutDecls(s)
	}
}

func NewEmptyLayoutDeclsContext() *LayoutDeclsContext {
	var p = new(LayoutDeclsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_layoutDecls
	return p
}

func InitEmptyLayoutDeclsContext(p *LayoutDeclsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_layoutDecls
}

func NewLayoutDeclsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LayoutDeclsContext {
	var p = new(LayoutDeclsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_layoutDecls

	return p
}
