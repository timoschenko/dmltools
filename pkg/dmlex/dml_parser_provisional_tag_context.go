// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IProvisionalTagContext is an interface to support dynamic dispatch.
type IProvisionalTagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PROVISIONAL() antlr.TerminalNode
	IdentList() IIdentListContext
	SEMI() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode

	// IsProvisionalTagContext differentiates from other interfaces.
	IsProvisionalTagContext()
}

type ProvisionalTagContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ProvisionalTagContext) IsProvisionalTagContext() {}

func (s *ProvisionalTagContext) GetParser() antlr.Parser { return s.parser }

func (s *ProvisionalTagContext) PROVISIONAL() antlr.TerminalNode {
	return s.GetToken(dmlParserPROVISIONAL, 0)
}

func (s *ProvisionalTagContext) IdentList() IIdentListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentListContext)
}

func (s *ProvisionalTagContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *ProvisionalTagContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *ProvisionalTagContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *ProvisionalTagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProvisionalTagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProvisionalTagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterProvisionalTag(s)
	}
}

func (s *ProvisionalTagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitProvisionalTag(s)
	}
}

func NewEmptyProvisionalTagContext() *ProvisionalTagContext {
	var p = new(ProvisionalTagContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_provisionalTag
	return p
}

func InitEmptyProvisionalTagContext(p *ProvisionalTagContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_provisionalTag
}

func NewProvisionalTagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProvisionalTagContext {
	var p = new(ProvisionalTagContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_provisionalTag

	return p
}
