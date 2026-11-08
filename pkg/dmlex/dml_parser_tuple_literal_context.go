// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ITuple_literalContext is an interface to support dynamic dispatch.
type ITuple_literalContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	Expression() IExpressionContext
	COMMA() antlr.TerminalNode
	Expression_list_ntc_nonempty() IExpression_list_ntc_nonemptyContext
	RPAREN() antlr.TerminalNode

	// IsTuple_literalContext differentiates from other interfaces.
	IsTuple_literalContext()
}

type Tuple_literalContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Tuple_literalContext) IsTuple_literalContext() {}

func (s *Tuple_literalContext) GetParser() antlr.Parser { return s.parser }

func (s *Tuple_literalContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *Tuple_literalContext) Expression() IExpressionContext {
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

func (s *Tuple_literalContext) COMMA() antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, 0)
}

func (s *Tuple_literalContext) Expression_list_ntc_nonempty() IExpression_list_ntc_nonemptyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpression_list_ntc_nonemptyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpression_list_ntc_nonemptyContext)
}

func (s *Tuple_literalContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *Tuple_literalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tuple_literalContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tuple_literalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterTuple_literal(s)
	}
}

func (s *Tuple_literalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitTuple_literal(s)
	}
}

func NewEmptyTuple_literalContext() *Tuple_literalContext {
	var p = new(Tuple_literalContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_tuple_literal
	return p
}

func InitEmptyTuple_literalContext(p *Tuple_literalContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_tuple_literal
}

func NewTuple_literalContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tuple_literalContext {
	var p = new(Tuple_literalContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_tuple_literal

	return p
}
