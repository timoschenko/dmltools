// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IObjectStatementContext is an interface to support dynamic dispatch.
type IObjectStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Object() IObjectContext
	Param() IParamContext
	MethodUnified() IMethodUnifiedContext
	IstemplateStatement() IIstemplateStatementContext
	ObjectIf() IObjectIfContext
	ErrorStatement() IErrorStatementContext
	IneachStatement() IIneachStatementContext

	// IsObjectStatementContext differentiates from other interfaces.
	IsObjectStatementContext()
}

type ObjectStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ObjectStatementContext) IsObjectStatementContext() {}

func (s *ObjectStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectStatementContext) Object() IObjectContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectContext)
}

func (s *ObjectStatementContext) Param() IParamContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *ObjectStatementContext) MethodUnified() IMethodUnifiedContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMethodUnifiedContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMethodUnifiedContext)
}

func (s *ObjectStatementContext) IstemplateStatement() IIstemplateStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIstemplateStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIstemplateStatementContext)
}

func (s *ObjectStatementContext) ObjectIf() IObjectIfContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectIfContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectIfContext)
}

func (s *ObjectStatementContext) ErrorStatement() IErrorStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IErrorStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IErrorStatementContext)
}

func (s *ObjectStatementContext) IneachStatement() IIneachStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIneachStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIneachStatementContext)
}

func (s *ObjectStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterObjectStatement(s)
	}
}

func (s *ObjectStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitObjectStatement(s)
	}
}

func NewEmptyObjectStatementContext() *ObjectStatementContext {
	var p = new(ObjectStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectStatement
	return p
}

func InitEmptyObjectStatementContext(p *ObjectStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectStatement
}

func NewObjectStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectStatementContext {
	var p = new(ObjectStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_objectStatement

	return p
}
