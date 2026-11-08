// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IInnerForLoopStatementContext is an interface to support dynamic dispatch.
type IInnerForLoopStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AssignStatement() IAssignStatementContext
	Assignop() IAssignopContext
	Expression() IExpressionContext

	// IsInnerForLoopStatementContext differentiates from other interfaces.
	IsInnerForLoopStatementContext()
}

type InnerForLoopStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*InnerForLoopStatementContext) IsInnerForLoopStatementContext() {}

func (s *InnerForLoopStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *InnerForLoopStatementContext) AssignStatement() IAssignStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAssignStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAssignStatementContext)
}

func (s *InnerForLoopStatementContext) Assignop() IAssignopContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAssignopContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAssignopContext)
}

func (s *InnerForLoopStatementContext) Expression() IExpressionContext {
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

func (s *InnerForLoopStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InnerForLoopStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InnerForLoopStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterInnerForLoopStatement(s)
	}
}

func (s *InnerForLoopStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitInnerForLoopStatement(s)
	}
}

func NewEmptyInnerForLoopStatementContext() *InnerForLoopStatementContext {
	var p = new(InnerForLoopStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_innerForLoopStatement
	return p
}

func InitEmptyInnerForLoopStatementContext(p *InnerForLoopStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_innerForLoopStatement
}

func NewInnerForLoopStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InnerForLoopStatementContext {
	var p = new(InnerForLoopStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_innerForLoopStatement

	return p
}
