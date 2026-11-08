// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ILayoutContext is an interface to support dynamic dispatch.
type ILayoutContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LAYOUT() antlr.TerminalNode
	STRING_LITERAL() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllLayoutDecls() []ILayoutDeclsContext
	LayoutDecls(i int) ILayoutDeclsContext

	// IsLayoutContext differentiates from other interfaces.
	IsLayoutContext()
}

type LayoutContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*LayoutContext) IsLayoutContext() {}

func (s *LayoutContext) GetParser() antlr.Parser { return s.parser }

func (s *LayoutContext) LAYOUT() antlr.TerminalNode {
	return s.GetToken(dmlParserLAYOUT, 0)
}

func (s *LayoutContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(dmlParserSTRING_LITERAL, 0)
}

func (s *LayoutContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *LayoutContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *LayoutContext) AllLayoutDecls() []ILayoutDeclsContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILayoutDeclsContext); ok {
			len++
		}
	}

	tst := make([]ILayoutDeclsContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILayoutDeclsContext); ok {
			tst[i] = t.(ILayoutDeclsContext)
			i++
		}
	}

	return tst
}

func (s *LayoutContext) LayoutDecls(i int) ILayoutDeclsContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILayoutDeclsContext); ok {
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

	return t.(ILayoutDeclsContext)
}

func (s *LayoutContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LayoutContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LayoutContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLayout(s)
	}
}

func (s *LayoutContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLayout(s)
	}
}

func NewEmptyLayoutContext() *LayoutContext {
	var p = new(LayoutContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_layout
	return p
}

func InitEmptyLayoutContext(p *LayoutContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_layout
}

func NewLayoutContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LayoutContext {
	var p = new(LayoutContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_layout

	return p
}
