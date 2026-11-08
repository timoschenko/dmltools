// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IExpressionAccessIdentContext is an interface to support dynamic dispatch.
type IExpressionAccessIdentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEFAULT() antlr.TerminalNode
	THIS() antlr.TerminalNode
	TEMPLATE() antlr.TerminalNode
	Objident() IObjidentContext

	// IsExpressionAccessIdentContext differentiates from other interfaces.
	IsExpressionAccessIdentContext()
}

type ExpressionAccessIdentContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ExpressionAccessIdentContext) IsExpressionAccessIdentContext() {}

func (s *ExpressionAccessIdentContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionAccessIdentContext) DEFAULT() antlr.TerminalNode {
	return s.GetToken(dmlParserDEFAULT, 0)
}

func (s *ExpressionAccessIdentContext) THIS() antlr.TerminalNode {
	return s.GetToken(dmlParserTHIS, 0)
}

func (s *ExpressionAccessIdentContext) TEMPLATE() antlr.TerminalNode {
	return s.GetToken(dmlParserTEMPLATE, 0)
}

func (s *ExpressionAccessIdentContext) Objident() IObjidentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjidentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjidentContext)
}

func (s *ExpressionAccessIdentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionAccessIdentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionAccessIdentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterExpressionAccessIdent(s)
	}
}

func (s *ExpressionAccessIdentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitExpressionAccessIdent(s)
	}
}

func NewEmptyExpressionAccessIdentContext() *ExpressionAccessIdentContext {
	var p = new(ExpressionAccessIdentContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_expressionAccessIdent
	return p
}

func InitEmptyExpressionAccessIdentContext(p *ExpressionAccessIdentContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_expressionAccessIdent
}

func NewExpressionAccessIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionAccessIdentContext {
	var p = new(ExpressionAccessIdentContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_expressionAccessIdent

	return p
}
