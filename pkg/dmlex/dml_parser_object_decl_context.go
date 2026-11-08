// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IObject_declContext is an interface to support dynamic dispatch.
type IObject_declContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	REGISTER() antlr.TerminalNode
	Objident() IObjidentContext
	ObjectBody() IObjectBodyContext
	Array_list() IArray_listContext
	SIZE() antlr.TerminalNode
	SizeSpec() ISizeSpecContext
	OffsetSpec() IOffsetSpecContext
	FIELD() antlr.TerminalNode
	BitrangeDecl() IBitrangeDeclContext
	Object_type() IObject_typeContext
	IN() antlr.TerminalNode

	// IsObject_declContext differentiates from other interfaces.
	IsObject_declContext()
}

type Object_declContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Object_declContext) IsObject_declContext() {}

func (s *Object_declContext) GetParser() antlr.Parser { return s.parser }

func (s *Object_declContext) REGISTER() antlr.TerminalNode {
	return s.GetToken(dmlParserREGISTER, 0)
}

func (s *Object_declContext) Objident() IObjidentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjidentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjidentContext)
}

func (s *Object_declContext) ObjectBody() IObjectBodyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectBodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectBodyContext)
}

func (s *Object_declContext) Array_list() IArray_listContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArray_listContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArray_listContext)
}

func (s *Object_declContext) SIZE() antlr.TerminalNode {
	return s.GetToken(dmlParserSIZE, 0)
}

func (s *Object_declContext) SizeSpec() ISizeSpecContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISizeSpecContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISizeSpecContext)
}

func (s *Object_declContext) OffsetSpec() IOffsetSpecContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOffsetSpecContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOffsetSpecContext)
}

func (s *Object_declContext) FIELD() antlr.TerminalNode {
	return s.GetToken(dmlParserFIELD, 0)
}

func (s *Object_declContext) BitrangeDecl() IBitrangeDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBitrangeDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBitrangeDeclContext)
}

func (s *Object_declContext) Object_type() IObject_typeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObject_typeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObject_typeContext)
}

func (s *Object_declContext) IN() antlr.TerminalNode {
	return s.GetToken(dmlParserIN, 0)
}

func (s *Object_declContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Object_declContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Object_declContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterObject_decl(s)
	}
}

func (s *Object_declContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitObject_decl(s)
	}
}

func NewEmptyObject_declContext() *Object_declContext {
	var p = new(Object_declContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_object_decl
	return p
}

func InitEmptyObject_declContext(p *Object_declContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_object_decl
}

func NewObject_declContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Object_declContext {
	var p = new(Object_declContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_object_decl

	return p
}
