// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICdeclInitializerContext is an interface to support dynamic dispatch.
type ICdeclInitializerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	SingleInitializerArray() ISingleInitializerArrayContext
	RBRACE() antlr.TerminalNode
	Designated_struct_initializer_list() IDesignated_struct_initializer_listContext

	// IsCdeclInitializerContext differentiates from other interfaces.
	IsCdeclInitializerContext()
}

type CdeclInitializerContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*CdeclInitializerContext) IsCdeclInitializerContext() {}

func (s *CdeclInitializerContext) GetParser() antlr.Parser { return s.parser }

func (s *CdeclInitializerContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *CdeclInitializerContext) SingleInitializerArray() ISingleInitializerArrayContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISingleInitializerArrayContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISingleInitializerArrayContext)
}

func (s *CdeclInitializerContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *CdeclInitializerContext) Designated_struct_initializer_list() IDesignated_struct_initializer_listContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDesignated_struct_initializer_listContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDesignated_struct_initializer_listContext)
}

func (s *CdeclInitializerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CdeclInitializerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CdeclInitializerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdeclInitializer(s)
	}
}

func (s *CdeclInitializerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdeclInitializer(s)
	}
}

func NewEmptyCdeclInitializerContext() *CdeclInitializerContext {
	var p = new(CdeclInitializerContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdeclInitializer
	return p
}

func InitEmptyCdeclInitializerContext(p *CdeclInitializerContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdeclInitializer
}

func NewCdeclInitializerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CdeclInitializerContext {
	var p = new(CdeclInitializerContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_cdeclInitializer

	return p
}
