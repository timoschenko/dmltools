// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICtypedecl_simpleContext is an interface to support dynamic dispatch.
type ICtypedecl_simpleContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	Ctypedecl_ptr() ICtypedecl_ptrContext
	RPAREN() antlr.TerminalNode

	// IsCtypedecl_simpleContext differentiates from other interfaces.
	IsCtypedecl_simpleContext()
}

type Ctypedecl_simpleContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Ctypedecl_simpleContext) IsCtypedecl_simpleContext() {}

func (s *Ctypedecl_simpleContext) GetParser() antlr.Parser { return s.parser }

func (s *Ctypedecl_simpleContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *Ctypedecl_simpleContext) Ctypedecl_ptr() ICtypedecl_ptrContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICtypedecl_ptrContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICtypedecl_ptrContext)
}

func (s *Ctypedecl_simpleContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *Ctypedecl_simpleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Ctypedecl_simpleContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Ctypedecl_simpleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCtypedecl_simple(s)
	}
}

func (s *Ctypedecl_simpleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCtypedecl_simple(s)
	}
}

func NewEmptyCtypedecl_simpleContext() *Ctypedecl_simpleContext {
	var p = new(Ctypedecl_simpleContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ctypedecl_simple
	return p
}

func InitEmptyCtypedecl_simpleContext(p *Ctypedecl_simpleContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ctypedecl_simple
}

func NewCtypedecl_simpleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Ctypedecl_simpleContext {
	var p = new(Ctypedecl_simpleContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_ctypedecl_simple

	return p
}
