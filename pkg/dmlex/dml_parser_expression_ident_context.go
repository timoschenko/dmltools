// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IExpressionIdentContext is an interface to support dynamic dispatch.
type IExpressionIdentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NumericLiteral() INumericLiteralContext
	StringLiteral() IStringLiteralContext
	UNDEFINED() antlr.TerminalNode
	ExpressionAccessIdent() IExpressionAccessIdentContext

	// IsExpressionIdentContext differentiates from other interfaces.
	IsExpressionIdentContext()
}

type ExpressionIdentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ExpressionIdentContext) IsExpressionIdentContext() {}

func (s *ExpressionIdentContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionIdentContext) NumericLiteral() INumericLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INumericLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(INumericLiteralContext)
}

func (s *ExpressionIdentContext) StringLiteral() IStringLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStringLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStringLiteralContext)
}

func (s *ExpressionIdentContext) UNDEFINED() antlr.TerminalNode {
	return s.GetToken(dmlParserUNDEFINED, 0)
}

func (s *ExpressionIdentContext) ExpressionAccessIdent() IExpressionAccessIdentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionAccessIdentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionAccessIdentContext)
}

func (s *ExpressionIdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionIdentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionIdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterExpressionIdent(s)
	}
}

func (s *ExpressionIdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitExpressionIdent(s)
	}
}

func NewEmptyExpressionIdentContext() *ExpressionIdentContext {
	var p = new(ExpressionIdentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_expressionIdent
	return p
}

func InitEmptyExpressionIdentContext(p *ExpressionIdentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_expressionIdent
}

func NewExpressionIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionIdentContext {
	var p = new(ExpressionIdentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_expressionIdent

	return p
}
