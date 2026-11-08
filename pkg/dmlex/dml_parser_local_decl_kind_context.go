// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ILocal_decl_kindContext is an interface to support dynamic dispatch.
type ILocal_decl_kindContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LOCAL() antlr.TerminalNode
	SESSION() antlr.TerminalNode
	STATIC() antlr.TerminalNode
	SAVED() antlr.TerminalNode

	// IsLocal_decl_kindContext differentiates from other interfaces.
	IsLocal_decl_kindContext()
}

type Local_decl_kindContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Local_decl_kindContext) IsLocal_decl_kindContext() {}

func (s *Local_decl_kindContext) GetParser() antlr.Parser { return s.parser }

func (s *Local_decl_kindContext) LOCAL() antlr.TerminalNode {
	return s.GetToken(dmlParserLOCAL, 0)
}

func (s *Local_decl_kindContext) SESSION() antlr.TerminalNode {
	return s.GetToken(dmlParserSESSION, 0)
}

func (s *Local_decl_kindContext) STATIC() antlr.TerminalNode {
	return s.GetToken(dmlParserSTATIC, 0)
}

func (s *Local_decl_kindContext) SAVED() antlr.TerminalNode {
	return s.GetToken(dmlParserSAVED, 0)
}

func (s *Local_decl_kindContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Local_decl_kindContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Local_decl_kindContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterLocal_decl_kind(s)
	}
}

func (s *Local_decl_kindContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitLocal_decl_kind(s)
	}
}

func NewEmptyLocal_decl_kindContext() *Local_decl_kindContext {
	var p = new(Local_decl_kindContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_local_decl_kind
	return p
}

func InitEmptyLocal_decl_kindContext(p *Local_decl_kindContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_local_decl_kind
}

func NewLocal_decl_kindContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Local_decl_kindContext {
	var p = new(Local_decl_kindContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_local_decl_kind

	return p
}
