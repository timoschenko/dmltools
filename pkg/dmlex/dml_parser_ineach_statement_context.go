// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IIneachStatementContext is an interface to support dynamic dispatch.
type IIneachStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IN() antlr.TerminalNode
	EACH() antlr.TerminalNode
	IstemplateList() IIstemplateListContext
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllObjectStatement() []IObjectStatementContext
	ObjectStatement(i int) IObjectStatementContext

	// IsIneachStatementContext differentiates from other interfaces.
	IsIneachStatementContext()
}

type IneachStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*IneachStatementContext) IsIneachStatementContext() {}

func (s *IneachStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *IneachStatementContext) IN() antlr.TerminalNode {
	return s.GetToken(dmlParserIN, 0)
}

func (s *IneachStatementContext) EACH() antlr.TerminalNode {
	return s.GetToken(dmlParserEACH, 0)
}

func (s *IneachStatementContext) IstemplateList() IIstemplateListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIstemplateListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIstemplateListContext)
}

func (s *IneachStatementContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *IneachStatementContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *IneachStatementContext) AllObjectStatement() []IObjectStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IObjectStatementContext); ok {
			len++
		}
	}

	tst := make([]IObjectStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IObjectStatementContext); ok {
			tst[i] = t.(IObjectStatementContext)
			i++
		}
	}

	return tst
}

func (s *IneachStatementContext) ObjectStatement(i int) IObjectStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectStatementContext); ok {
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

	return t.(IObjectStatementContext)
}

func (s *IneachStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IneachStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IneachStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterIneachStatement(s)
	}
}

func (s *IneachStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitIneachStatement(s)
	}
}

func NewEmptyIneachStatementContext() *IneachStatementContext {
	var p = new(IneachStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ineachStatement
	return p
}

func InitEmptyIneachStatementContext(p *IneachStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_ineachStatement
}

func NewIneachStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IneachStatementContext {
	var p = new(IneachStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_ineachStatement

	return p
}
