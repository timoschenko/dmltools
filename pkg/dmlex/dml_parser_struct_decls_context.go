// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IStructDeclsContext is an interface to support dynamic dispatch.
type IStructDeclsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Cdecl() ICdeclContext
	SEMI() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Expression() IExpressionContext

	// IsStructDeclsContext differentiates from other interfaces.
	IsStructDeclsContext()
}

type StructDeclsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*StructDeclsContext) IsStructDeclsContext() {}

func (s *StructDeclsContext) GetParser() antlr.Parser { return s.parser }

func (s *StructDeclsContext) Cdecl() ICdeclContext {
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

func (s *StructDeclsContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *StructDeclsContext) COLON() antlr.TerminalNode {
	return s.GetToken(dmlParserCOLON, 0)
}

func (s *StructDeclsContext) Expression() IExpressionContext {
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

func (s *StructDeclsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StructDeclsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StructDeclsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterStructDecls(s)
	}
}

func (s *StructDeclsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitStructDecls(s)
	}
}

func NewEmptyStructDeclsContext() *StructDeclsContext {
	var p = new(StructDeclsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_structDecls
	return p
}

func InitEmptyStructDeclsContext(p *StructDeclsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_structDecls
}

func NewStructDeclsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StructDeclsContext {
	var p = new(StructDeclsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_structDecls

	return p
}
