// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IToplevelIfContext is an interface to support dynamic dispatch.
type IToplevelIfContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HASH_IF() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ConditionExpression() IConditionExpressionContext
	RPAREN() antlr.TerminalNode
	ToplevelIfStatement() IToplevelIfStatementContext
	HASH_ELSE() antlr.TerminalNode
	ToplevelElseStatement() IToplevelElseStatementContext

	// IsToplevelIfContext differentiates from other interfaces.
	IsToplevelIfContext()
}

type ToplevelIfContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ToplevelIfContext) IsToplevelIfContext() {}

func (s *ToplevelIfContext) GetParser() antlr.Parser { return s.parser }

func (s *ToplevelIfContext) HASH_IF() antlr.TerminalNode {
	return s.GetToken(dmlParserHASH_IF, 0)
}

func (s *ToplevelIfContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *ToplevelIfContext) ConditionExpression() IConditionExpressionContext {
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

func (s *ToplevelIfContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *ToplevelIfContext) ToplevelIfStatement() IToplevelIfStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IToplevelIfStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IToplevelIfStatementContext)
}

func (s *ToplevelIfContext) HASH_ELSE() antlr.TerminalNode {
	return s.GetToken(dmlParserHASH_ELSE, 0)
}

func (s *ToplevelIfContext) ToplevelElseStatement() IToplevelElseStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IToplevelElseStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IToplevelElseStatementContext)
}

func (s *ToplevelIfContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ToplevelIfContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ToplevelIfContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterToplevelIf(s)
	}
}

func (s *ToplevelIfContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitToplevelIf(s)
	}
}

func NewEmptyToplevelIfContext() *ToplevelIfContext {
	var p = new(ToplevelIfContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_toplevelIf
	return p
}

func InitEmptyToplevelIfContext(p *ToplevelIfContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_toplevelIf
}

func NewToplevelIfContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ToplevelIfContext {
	var p = new(ToplevelIfContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_toplevelIf

	return p
}
