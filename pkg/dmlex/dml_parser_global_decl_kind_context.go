// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IGlobal_decl_kindContext is an interface to support dynamic dispatch.
type IGlobal_decl_kindContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SESSION() antlr.TerminalNode
	SAVED() antlr.TerminalNode

	// IsGlobal_decl_kindContext differentiates from other interfaces.
	IsGlobal_decl_kindContext()
}

type Global_decl_kindContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Global_decl_kindContext) IsGlobal_decl_kindContext() {}

func (s *Global_decl_kindContext) GetParser() antlr.Parser { return s.parser }

func (s *Global_decl_kindContext) SESSION() antlr.TerminalNode {
	return s.GetToken(dmlParserSESSION, 0)
}

func (s *Global_decl_kindContext) SAVED() antlr.TerminalNode {
	return s.GetToken(dmlParserSAVED, 0)
}

func (s *Global_decl_kindContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Global_decl_kindContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Global_decl_kindContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterGlobal_decl_kind(s)
	}
}

func (s *Global_decl_kindContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitGlobal_decl_kind(s)
	}
}

func NewEmptyGlobal_decl_kindContext() *Global_decl_kindContext {
	var p = new(Global_decl_kindContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_global_decl_kind
	return p
}

func InitEmptyGlobal_decl_kindContext(p *Global_decl_kindContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_global_decl_kind
}

func NewGlobal_decl_kindContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Global_decl_kindContext {
	var p = new(Global_decl_kindContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_global_decl_kind

	return p
}
