// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IObjectElseStatementContext is an interface to support dynamic dispatch.
type IObjectElseStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllObjectStatement() []IObjectStatementContext
	ObjectStatement(i int) IObjectStatementContext

	// IsObjectElseStatementContext differentiates from other interfaces.
	IsObjectElseStatementContext()
}

type ObjectElseStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ObjectElseStatementContext) IsObjectElseStatementContext() {}

func (s *ObjectElseStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectElseStatementContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *ObjectElseStatementContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *ObjectElseStatementContext) AllObjectStatement() []IObjectStatementContext {
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

func (s *ObjectElseStatementContext) ObjectStatement(i int) IObjectStatementContext {
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

func (s *ObjectElseStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectElseStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectElseStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterObjectElseStatement(s)
	}
}

func (s *ObjectElseStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitObjectElseStatement(s)
	}
}

func NewEmptyObjectElseStatementContext() *ObjectElseStatementContext {
	var p = new(ObjectElseStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectElseStatement
	return p
}

func InitEmptyObjectElseStatementContext(p *ObjectElseStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectElseStatement
}

func NewObjectElseStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectElseStatementContext {
	var p = new(ObjectElseStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_objectElseStatement

	return p
}
