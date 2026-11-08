// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IObjectIfStatementContext is an interface to support dynamic dispatch.
type IObjectIfStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllObjectStatement() []IObjectStatementContext
	ObjectStatement(i int) IObjectStatementContext

	// IsObjectIfStatementContext differentiates from other interfaces.
	IsObjectIfStatementContext()
}

type ObjectIfStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ObjectIfStatementContext) IsObjectIfStatementContext() {}

func (s *ObjectIfStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectIfStatementContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *ObjectIfStatementContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *ObjectIfStatementContext) AllObjectStatement() []IObjectStatementContext {
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

func (s *ObjectIfStatementContext) ObjectStatement(i int) IObjectStatementContext {
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

func (s *ObjectIfStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectIfStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectIfStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterObjectIfStatement(s)
	}
}

func (s *ObjectIfStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitObjectIfStatement(s)
	}
}

func NewEmptyObjectIfStatementContext() *ObjectIfStatementContext {
	var p = new(ObjectIfStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectIfStatement
	return p
}

func InitEmptyObjectIfStatementContext(p *ObjectIfStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectIfStatement
}

func NewObjectIfStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectIfStatementContext {
	var p = new(ObjectIfStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_objectIfStatement

	return p
}
