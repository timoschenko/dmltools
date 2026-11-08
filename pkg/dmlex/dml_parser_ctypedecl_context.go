// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICtypedeclContext is an interface to support dynamic dispatch.
type ICtypedeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Basetype() IBasetypeContext
	Ctypedecl_ptr() ICtypedecl_ptrContext
	CONST() antlr.TerminalNode

	// IsCtypedeclContext differentiates from other interfaces.
	IsCtypedeclContext()
}

type CtypedeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*CtypedeclContext) IsCtypedeclContext() {}

func (s *CtypedeclContext) GetParser() antlr.Parser { return s.parser }

func (s *CtypedeclContext) Basetype() IBasetypeContext {
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

func (s *CtypedeclContext) Ctypedecl_ptr() ICtypedecl_ptrContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICtypedecl_ptrContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICtypedecl_ptrContext)
}

func (s *CtypedeclContext) CONST() antlr.TerminalNode {
	return s.GetToken(dmlParserCONST, 0)
}

func (s *CtypedeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CtypedeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CtypedeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCtypedecl(s)
	}
}

func (s *CtypedeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCtypedecl(s)
	}
}

func NewEmptyCtypedeclContext() *CtypedeclContext {
	var p = new(CtypedeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ctypedecl
	return p
}

func InitEmptyCtypedeclContext(p *CtypedeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ctypedecl
}

func NewCtypedeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CtypedeclContext {
	var p = new(CtypedeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_ctypedecl

	return p
}
