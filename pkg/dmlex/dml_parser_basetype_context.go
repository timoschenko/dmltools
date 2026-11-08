// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IBasetypeContext is an interface to support dynamic dispatch.
type IBasetypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Signed_or_unsigned_typeident() ISigned_or_unsigned_typeidentContext
	Struct_() IStructContext
	Union() IUnionContext
	Layout() ILayoutContext
	Bitfields() IBitfieldsContext
	Typeof() ITypeofContext
	SEQUENCE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	HOOK() antlr.TerminalNode
	Cdecl_list() ICdecl_listContext

	// IsBasetypeContext differentiates from other interfaces.
	IsBasetypeContext()
}

type BasetypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*BasetypeContext) IsBasetypeContext() {}

func (s *BasetypeContext) GetParser() antlr.Parser { return s.parser }

func (s *BasetypeContext) Signed_or_unsigned_typeident() ISigned_or_unsigned_typeidentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISigned_or_unsigned_typeidentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISigned_or_unsigned_typeidentContext)
}

func (s *BasetypeContext) Struct_() IStructContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStructContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStructContext)
}

func (s *BasetypeContext) Union() IUnionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionContext)
}

func (s *BasetypeContext) Layout() ILayoutContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILayoutContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILayoutContext)
}

func (s *BasetypeContext) Bitfields() IBitfieldsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBitfieldsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBitfieldsContext)
}

func (s *BasetypeContext) Typeof() ITypeofContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeofContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeofContext)
}

func (s *BasetypeContext) SEQUENCE() antlr.TerminalNode {
	return s.GetToken(dmlParserSEQUENCE, 0)
}

func (s *BasetypeContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *BasetypeContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *BasetypeContext) HOOK() antlr.TerminalNode {
	return s.GetToken(dmlParserHOOK, 0)
}

func (s *BasetypeContext) Cdecl_list() ICdecl_listContext {
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

func (s *BasetypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BasetypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BasetypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterBasetype(s)
	}
}

func (s *BasetypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitBasetype(s)
	}
}

func NewEmptyBasetypeContext() *BasetypeContext {
	var p = new(BasetypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_basetype
	return p
}

func InitEmptyBasetypeContext(p *BasetypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_basetype
}

func NewBasetypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BasetypeContext {
	var p = new(BasetypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_basetype

	return p
}
