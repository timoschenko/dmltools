// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IDesignatedStructInitializerContext is an interface to support dynamic dispatch.
type IDesignatedStructInitializerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DOT() antlr.TerminalNode
	Ident() IIdentContext
	ASSIGN() antlr.TerminalNode
	SingleInitializer() ISingleInitializerContext

	// IsDesignatedStructInitializerContext differentiates from other interfaces.
	IsDesignatedStructInitializerContext()
}

type DesignatedStructInitializerContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*DesignatedStructInitializerContext) IsDesignatedStructInitializerContext() {}

func (s *DesignatedStructInitializerContext) GetParser() antlr.Parser { return s.parser }

func (s *DesignatedStructInitializerContext) DOT() antlr.TerminalNode {
	return s.GetToken(dmlParserDOT, 0)
}

func (s *DesignatedStructInitializerContext) Ident() IIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentContext)
}

func (s *DesignatedStructInitializerContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserASSIGN, 0)
}

func (s *DesignatedStructInitializerContext) SingleInitializer() ISingleInitializerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleInitializerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISingleInitializerContext)
}

func (s *DesignatedStructInitializerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DesignatedStructInitializerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DesignatedStructInitializerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterDesignatedStructInitializer(s)
	}
}

func (s *DesignatedStructInitializerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitDesignatedStructInitializer(s)
	}
}

func NewEmptyDesignatedStructInitializerContext() *DesignatedStructInitializerContext {
	var p = new(DesignatedStructInitializerContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_designatedStructInitializer
	return p
}

func InitEmptyDesignatedStructInitializerContext(p *DesignatedStructInitializerContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_designatedStructInitializer
}

func NewDesignatedStructInitializerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DesignatedStructInitializerContext {
	var p = new(DesignatedStructInitializerContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_designatedStructInitializer

	return p
}
