// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICdecl_ident_array_funcContext is an interface to support dynamic dispatch.
type ICdecl_ident_array_funcContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsCdecl_ident_array_funcContext differentiates from other interfaces.
	IsCdecl_ident_array_funcContext()
}

type Cdecl_ident_array_funcContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Cdecl_ident_array_funcContext) IsCdecl_ident_array_funcContext() {}

func (s *Cdecl_ident_array_funcContext) GetParser() antlr.Parser { return s.parser }

func (s *Cdecl_ident_array_funcContext) CopyAll(ctx *Cdecl_ident_array_funcContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *Cdecl_ident_array_funcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Cdecl_ident_array_funcContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func NewEmptyCdecl_ident_array_funcContext() *Cdecl_ident_array_funcContext {
	var p = new(Cdecl_ident_array_funcContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl_ident_array_func
	return p
}

func InitEmptyCdecl_ident_array_funcContext(p *Cdecl_ident_array_funcContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl_ident_array_func
}

func NewCdecl_ident_array_funcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Cdecl_ident_array_funcContext {
	var p = new(Cdecl_ident_array_funcContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_cdecl_ident_array_func

	return p
}
