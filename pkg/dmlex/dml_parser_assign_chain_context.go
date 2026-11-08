// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IAssign_chainContext is an interface to support dynamic dispatch.
type IAssign_chainContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	ASSIGN() antlr.TerminalNode
	Assign_chain() IAssign_chainContext
	Initializer() IInitializerContext

	// IsAssign_chainContext differentiates from other interfaces.
	IsAssign_chainContext()
}

type Assign_chainContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Assign_chainContext) IsAssign_chainContext() {}

func (s *Assign_chainContext) GetParser() antlr.Parser { return s.parser }

func (s *Assign_chainContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Assign_chainContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserASSIGN, 0)
}

func (s *Assign_chainContext) Assign_chain() IAssign_chainContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAssign_chainContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAssign_chainContext)
}

func (s *Assign_chainContext) Initializer() IInitializerContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInitializerContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInitializerContext)
}

func (s *Assign_chainContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Assign_chainContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Assign_chainContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterAssign_chain(s)
	}
}

func (s *Assign_chainContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitAssign_chain(s)
	}
}

func NewEmptyAssign_chainContext() *Assign_chainContext {
	var p = new(Assign_chainContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_assign_chain
	return p
}

func InitEmptyAssign_chainContext(p *Assign_chainContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_assign_chain
}

func NewAssign_chainContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Assign_chainContext {
	var p = new(Assign_chainContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_assign_chain

	return p
}
