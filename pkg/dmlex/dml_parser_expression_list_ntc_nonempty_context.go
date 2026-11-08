// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IExpression_list_ntc_nonemptyContext is an interface to support dynamic dispatch.
type IExpression_list_ntc_nonemptyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsExpression_list_ntc_nonemptyContext differentiates from other interfaces.
	IsExpression_list_ntc_nonemptyContext()
}

type Expression_list_ntc_nonemptyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Expression_list_ntc_nonemptyContext) IsExpression_list_ntc_nonemptyContext() {}

func (s *Expression_list_ntc_nonemptyContext) GetParser() antlr.Parser { return s.parser }

func (s *Expression_list_ntc_nonemptyContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Expression_list_ntc_nonemptyContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *Expression_list_ntc_nonemptyContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCOMMA)
}

func (s *Expression_list_ntc_nonemptyContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, i)
}

func (s *Expression_list_ntc_nonemptyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Expression_list_ntc_nonemptyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Expression_list_ntc_nonemptyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterExpression_list_ntc_nonempty(s)
	}
}

func (s *Expression_list_ntc_nonemptyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitExpression_list_ntc_nonempty(s)
	}
}

func NewEmptyExpression_list_ntc_nonemptyContext() *Expression_list_ntc_nonemptyContext {
	var p = new(Expression_list_ntc_nonemptyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_expression_list_ntc_nonempty
	return p
}

func InitEmptyExpression_list_ntc_nonemptyContext(p *Expression_list_ntc_nonemptyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_expression_list_ntc_nonempty
}

func NewExpression_list_ntc_nonemptyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Expression_list_ntc_nonemptyContext {
	var p = new(Expression_list_ntc_nonemptyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_expression_list_ntc_nonempty

	return p
}
