// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICtypedecl_ptrContext is an interface to support dynamic dispatch.
type ICtypedecl_ptrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Ctypedecl_simple() ICtypedecl_simpleContext
	AllMUL() []antlr.TerminalNode
	MUL(i int) antlr.TerminalNode
	AllCONST() []antlr.TerminalNode
	CONST(i int) antlr.TerminalNode

	// IsCtypedecl_ptrContext differentiates from other interfaces.
	IsCtypedecl_ptrContext()
}

type Ctypedecl_ptrContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Ctypedecl_ptrContext) IsCtypedecl_ptrContext() {}

func (s *Ctypedecl_ptrContext) GetParser() antlr.Parser { return s.parser }

func (s *Ctypedecl_ptrContext) Ctypedecl_simple() ICtypedecl_simpleContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICtypedecl_simpleContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICtypedecl_simpleContext)
}

func (s *Ctypedecl_ptrContext) AllMUL() []antlr.TerminalNode {
	return s.GetTokens(dmlParserMUL)
}

func (s *Ctypedecl_ptrContext) MUL(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserMUL, i)
}

func (s *Ctypedecl_ptrContext) AllCONST() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCONST)
}

func (s *Ctypedecl_ptrContext) CONST(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCONST, i)
}

func (s *Ctypedecl_ptrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Ctypedecl_ptrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Ctypedecl_ptrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCtypedecl_ptr(s)
	}
}

func (s *Ctypedecl_ptrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCtypedecl_ptr(s)
	}
}

func NewEmptyCtypedecl_ptrContext() *Ctypedecl_ptrContext {
	var p = new(Ctypedecl_ptrContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ctypedecl_ptr
	return p
}

func InitEmptyCtypedecl_ptrContext(p *Ctypedecl_ptrContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ctypedecl_ptr
}

func NewCtypedecl_ptrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Ctypedecl_ptrContext {
	var p = new(Ctypedecl_ptrContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_ctypedecl_ptr

	return p
}
