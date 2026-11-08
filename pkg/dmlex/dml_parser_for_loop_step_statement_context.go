// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IForLoopStepStatementContext is an interface to support dynamic dispatch.
type IForLoopStepStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllInnerForLoopStatement() []IInnerForLoopStatementContext
	InnerForLoopStatement(i int) IInnerForLoopStatementContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsForLoopStepStatementContext differentiates from other interfaces.
	IsForLoopStepStatementContext()
}

type ForLoopStepStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ForLoopStepStatementContext) IsForLoopStepStatementContext() {}

func (s *ForLoopStepStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ForLoopStepStatementContext) AllInnerForLoopStatement() []IInnerForLoopStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IInnerForLoopStatementContext); ok {
			len++
		}
	}

	tst := make([]IInnerForLoopStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IInnerForLoopStatementContext); ok {
			tst[i] = t.(IInnerForLoopStatementContext)
			i++
		}
	}

	return tst
}

func (s *ForLoopStepStatementContext) InnerForLoopStatement(i int) IInnerForLoopStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInnerForLoopStatementContext); ok {
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

	return t.(IInnerForLoopStatementContext)
}

func (s *ForLoopStepStatementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCOMMA)
}

func (s *ForLoopStepStatementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, i)
}

func (s *ForLoopStepStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForLoopStepStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForLoopStepStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterForLoopStepStatement(s)
	}
}

func (s *ForLoopStepStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitForLoopStepStatement(s)
	}
}

func NewEmptyForLoopStepStatementContext() *ForLoopStepStatementContext {
	var p = new(ForLoopStepStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_forLoopStepStatement
	return p
}

func InitEmptyForLoopStepStatementContext(p *ForLoopStepStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_forLoopStepStatement
}

func NewForLoopStepStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForLoopStepStatementContext {
	var p = new(ForLoopStepStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_forLoopStepStatement

	return p
}
