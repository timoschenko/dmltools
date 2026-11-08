// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICvarDeclContext is an interface to support dynamic dispatch.
type ICvarDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Cdecl() ICdeclContext
	ASSIGN() antlr.TerminalNode
	Initializer() IInitializerContext
	LPAREN() antlr.TerminalNode
	Cdecl_list() ICdecl_listContext
	RPAREN() antlr.TerminalNode

	// IsCvarDeclContext differentiates from other interfaces.
	IsCvarDeclContext()
}

type CvarDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*CvarDeclContext) IsCvarDeclContext() {}

func (s *CvarDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *CvarDeclContext) Cdecl() ICdeclContext {
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

func (s *CvarDeclContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserASSIGN, 0)
}

func (s *CvarDeclContext) Initializer() IInitializerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInitializerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInitializerContext)
}

func (s *CvarDeclContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *CvarDeclContext) Cdecl_list() ICdecl_listContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdecl_listContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdecl_listContext)
}

func (s *CvarDeclContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *CvarDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CvarDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CvarDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCvarDecl(s)
	}
}

func (s *CvarDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCvarDecl(s)
	}
}

func NewEmptyCvarDeclContext() *CvarDeclContext {
	var p = new(CvarDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cvarDecl
	return p
}

func InitEmptyCvarDeclContext(p *CvarDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cvarDecl
}

func NewCvarDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CvarDeclContext {
	var p = new(CvarDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_cvarDecl

	return p
}
