// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICdecl_qualifier_ident_array_funcContext is an interface to support dynamic dispatch.
type ICdecl_qualifier_ident_array_funcContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Cdecl_ident_array_func() ICdecl_ident_array_funcContext
	AllTypeQualifier() []ITypeQualifierContext
	TypeQualifier(i int) ITypeQualifierContext

	// IsCdecl_qualifier_ident_array_funcContext differentiates from other interfaces.
	IsCdecl_qualifier_ident_array_funcContext()
}

type Cdecl_qualifier_ident_array_funcContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Cdecl_qualifier_ident_array_funcContext) IsCdecl_qualifier_ident_array_funcContext() {}

func (s *Cdecl_qualifier_ident_array_funcContext) GetParser() antlr.Parser { return s.parser }

func (s *Cdecl_qualifier_ident_array_funcContext) Cdecl_ident_array_func() ICdecl_ident_array_funcContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdecl_ident_array_funcContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdecl_ident_array_funcContext)
}

func (s *Cdecl_qualifier_ident_array_funcContext) AllTypeQualifier() []ITypeQualifierContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeQualifierContext); ok {
			len++
		}
	}

	tst := make([]ITypeQualifierContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeQualifierContext); ok {
			tst[i] = t.(ITypeQualifierContext)
			i++
		}
	}

	return tst
}

func (s *Cdecl_qualifier_ident_array_funcContext) TypeQualifier(i int) ITypeQualifierContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeQualifierContext); ok {
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

	return t.(ITypeQualifierContext)
}

func (s *Cdecl_qualifier_ident_array_funcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Cdecl_qualifier_ident_array_funcContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Cdecl_qualifier_ident_array_funcContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdecl_qualifier_ident_array_func(s)
	}
}

func (s *Cdecl_qualifier_ident_array_funcContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdecl_qualifier_ident_array_func(s)
	}
}

func NewEmptyCdecl_qualifier_ident_array_funcContext() *Cdecl_qualifier_ident_array_funcContext {
	var p = new(Cdecl_qualifier_ident_array_funcContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl_qualifier_ident_array_func
	return p
}

func InitEmptyCdecl_qualifier_ident_array_funcContext(p *Cdecl_qualifier_ident_array_funcContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl_qualifier_ident_array_func
}

func NewCdecl_qualifier_ident_array_funcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Cdecl_qualifier_ident_array_funcContext {
	var p = new(Cdecl_qualifier_ident_array_funcContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_cdecl_qualifier_ident_array_func

	return p
}
