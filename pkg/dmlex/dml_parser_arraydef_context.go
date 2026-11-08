// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IArraydefContext is an interface to support dynamic dispatch.
type IArraydefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Ident() IIdentContext
	LT() antlr.TerminalNode
	Expression() IExpressionContext
	ELLIPSIS() antlr.TerminalNode

	// IsArraydefContext differentiates from other interfaces.
	IsArraydefContext()
}

type ArraydefContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ArraydefContext) IsArraydefContext() {}

func (s *ArraydefContext) GetParser() antlr.Parser { return s.parser }

func (s *ArraydefContext) Ident() IIdentContext {
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

func (s *ArraydefContext) LT() antlr.TerminalNode {
	return s.GetToken(dmlParserLT, 0)
}

func (s *ArraydefContext) Expression() IExpressionContext {
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

func (s *ArraydefContext) ELLIPSIS() antlr.TerminalNode {
	return s.GetToken(dmlParserELLIPSIS, 0)
}

func (s *ArraydefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArraydefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArraydefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterArraydef(s)
	}
}

func (s *ArraydefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitArraydef(s)
	}
}

func NewEmptyArraydefContext() *ArraydefContext {
	var p = new(ArraydefContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_arraydef
	return p
}

func InitEmptyArraydefContext(p *ArraydefContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_arraydef
}

func NewArraydefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArraydefContext {
	var p = new(ArraydefContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_arraydef

	return p
}
