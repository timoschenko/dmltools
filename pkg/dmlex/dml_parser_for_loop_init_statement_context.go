// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IForLoopInitStatementContext is an interface to support dynamic dispatch.
type IForLoopInitStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VariableDecl() IVariableDeclContext
	AllInnerForLoopStatement() []IInnerForLoopStatementContext
	InnerForLoopStatement(i int) IInnerForLoopStatementContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsForLoopInitStatementContext differentiates from other interfaces.
	IsForLoopInitStatementContext()
}

type ForLoopInitStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ForLoopInitStatementContext) IsForLoopInitStatementContext() {}

func (s *ForLoopInitStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ForLoopInitStatementContext) VariableDecl() IVariableDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVariableDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVariableDeclContext)
}

func (s *ForLoopInitStatementContext) AllInnerForLoopStatement() []IInnerForLoopStatementContext {
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

func (s *ForLoopInitStatementContext) InnerForLoopStatement(i int) IInnerForLoopStatementContext {
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

func (s *ForLoopInitStatementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCOMMA)
}

func (s *ForLoopInitStatementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, i)
}

func (s *ForLoopInitStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ForLoopInitStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ForLoopInitStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterForLoopInitStatement(s)
	}
}

func (s *ForLoopInitStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitForLoopInitStatement(s)
	}
}

func NewEmptyForLoopInitStatementContext() *ForLoopInitStatementContext {
	var p = new(ForLoopInitStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_forLoopInitStatement
	return p
}

func InitEmptyForLoopInitStatementContext(p *ForLoopInitStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_forLoopInitStatement
}

func NewForLoopInitStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForLoopInitStatementContext {
	var p = new(ForLoopInitStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_forLoopInitStatement

	return p
}
