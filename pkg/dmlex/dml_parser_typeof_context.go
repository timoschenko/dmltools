// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ITypeofContext is an interface to support dynamic dispatch.
type ITypeofContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPEOF() antlr.TerminalNode
	Expression() IExpressionContext

	// IsTypeofContext differentiates from other interfaces.
	IsTypeofContext()
}

type TypeofContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*TypeofContext) IsTypeofContext() {}

func (s *TypeofContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeofContext) TYPEOF() antlr.TerminalNode {
	return s.GetToken(dmlParserTYPEOF, 0)
}

func (s *TypeofContext) Expression() IExpressionContext {
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

func (s *TypeofContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeofContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeofContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterTypeof(s)
	}
}

func (s *TypeofContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitTypeof(s)
	}
}

func NewEmptyTypeofContext() *TypeofContext {
	var p = new(TypeofContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typeof
	return p
}

func InitEmptyTypeofContext(p *TypeofContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typeof
}

func NewTypeofContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeofContext {
	var p = new(TypeofContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_typeof

	return p
}
