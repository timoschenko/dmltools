// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ISelectStatementContext is an interface to support dynamic dispatch.
type ISelectStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HASH_SELECT() antlr.TerminalNode
	Ident() IIdentContext
	IN() antlr.TerminalNode
	AllLPAREN() []antlr.TerminalNode
	LPAREN(i int) antlr.TerminalNode
	Expression() IExpressionContext
	AllRPAREN() []antlr.TerminalNode
	RPAREN(i int) antlr.TerminalNode
	WHERE() antlr.TerminalNode
	ConditionExpression() IConditionExpressionContext
	LoopStatement() ILoopStatementContext
	HASH_ELSE() antlr.TerminalNode
	ElseStatement() IElseStatementContext
	SELECT() antlr.TerminalNode
	ELSE() antlr.TerminalNode

	// IsSelectStatementContext differentiates from other interfaces.
	IsSelectStatementContext()
}

type SelectStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*SelectStatementContext) IsSelectStatementContext() {}

func (s *SelectStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectStatementContext) HASH_SELECT() antlr.TerminalNode {
	return s.GetToken(dmlParserHASH_SELECT, 0)
}

func (s *SelectStatementContext) Ident() IIdentContext {
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

func (s *SelectStatementContext) IN() antlr.TerminalNode {
	return s.GetToken(dmlParserIN, 0)
}

func (s *SelectStatementContext) AllLPAREN() []antlr.TerminalNode {
	return s.GetTokens(dmlParserLPAREN)
}

func (s *SelectStatementContext) LPAREN(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, i)
}

func (s *SelectStatementContext) Expression() IExpressionContext {
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

func (s *SelectStatementContext) AllRPAREN() []antlr.TerminalNode {
	return s.GetTokens(dmlParserRPAREN)
}

func (s *SelectStatementContext) RPAREN(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, i)
}

func (s *SelectStatementContext) WHERE() antlr.TerminalNode {
	return s.GetToken(dmlParserWHERE, 0)
}

func (s *SelectStatementContext) ConditionExpression() IConditionExpressionContext {
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

func (s *SelectStatementContext) LoopStatement() ILoopStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILoopStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILoopStatementContext)
}

func (s *SelectStatementContext) HASH_ELSE() antlr.TerminalNode {
	return s.GetToken(dmlParserHASH_ELSE, 0)
}

func (s *SelectStatementContext) ElseStatement() IElseStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IElseStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IElseStatementContext)
}

func (s *SelectStatementContext) SELECT() antlr.TerminalNode {
	return s.GetToken(dmlParserSELECT, 0)
}

func (s *SelectStatementContext) ELSE() antlr.TerminalNode {
	return s.GetToken(dmlParserELSE, 0)
}

func (s *SelectStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterSelectStatement(s)
	}
}

func (s *SelectStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitSelectStatement(s)
	}
}

func NewEmptySelectStatementContext() *SelectStatementContext {
	var p = new(SelectStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_selectStatement
	return p
}

func InitEmptySelectStatementContext(p *SelectStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_selectStatement
}

func NewSelectStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectStatementContext {
	var p = new(SelectStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_selectStatement

	return p
}
