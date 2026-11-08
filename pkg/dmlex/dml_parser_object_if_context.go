// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IObjectIfContext is an interface to support dynamic dispatch.
type IObjectIfContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HASH_IF() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ConditionExpression() IConditionExpressionContext
	RPAREN() antlr.TerminalNode
	ObjectIfStatement() IObjectIfStatementContext
	HASH_ELSE() antlr.TerminalNode
	ObjectElseStatement() IObjectElseStatementContext

	// IsObjectIfContext differentiates from other interfaces.
	IsObjectIfContext()
}

type ObjectIfContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ObjectIfContext) IsObjectIfContext() {}

func (s *ObjectIfContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectIfContext) HASH_IF() antlr.TerminalNode {
	return s.GetToken(dmlParserHASH_IF, 0)
}

func (s *ObjectIfContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *ObjectIfContext) ConditionExpression() IConditionExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConditionExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConditionExpressionContext)
}

func (s *ObjectIfContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *ObjectIfContext) ObjectIfStatement() IObjectIfStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectIfStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectIfStatementContext)
}

func (s *ObjectIfContext) HASH_ELSE() antlr.TerminalNode {
	return s.GetToken(dmlParserHASH_ELSE, 0)
}

func (s *ObjectIfContext) ObjectElseStatement() IObjectElseStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectElseStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectElseStatementContext)
}

func (s *ObjectIfContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectIfContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectIfContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterObjectIf(s)
	}
}

func (s *ObjectIfContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitObjectIf(s)
	}
}

func NewEmptyObjectIfContext() *ObjectIfContext {
	var p = new(ObjectIfContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectIf
	return p
}

func InitEmptyObjectIfContext(p *ObjectIfContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectIf
}

func NewObjectIfContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectIfContext {
	var p = new(ObjectIfContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_objectIf

	return p
}
