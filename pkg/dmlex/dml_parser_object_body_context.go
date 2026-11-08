// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IObjectBodyContext is an interface to support dynamic dispatch.
type IObjectBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SEMI() antlr.TerminalNode
	IstemplateExpr() IIstemplateExprContext
	ObjectDescSpec() IObjectDescSpecContext
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllObjectStatement() []IObjectStatementContext
	ObjectStatement(i int) IObjectStatementContext

	// IsObjectBodyContext differentiates from other interfaces.
	IsObjectBodyContext()
}

type ObjectBodyContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ObjectBodyContext) IsObjectBodyContext() {}

func (s *ObjectBodyContext) GetParser() antlr.Parser { return s.parser }

func (s *ObjectBodyContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *ObjectBodyContext) IstemplateExpr() IIstemplateExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIstemplateExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIstemplateExprContext)
}

func (s *ObjectBodyContext) ObjectDescSpec() IObjectDescSpecContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjectDescSpecContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjectDescSpecContext)
}

func (s *ObjectBodyContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACE, 0)
}

func (s *ObjectBodyContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACE, 0)
}

func (s *ObjectBodyContext) AllObjectStatement() []IObjectStatementContext {
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

func (s *ObjectBodyContext) ObjectStatement(i int) IObjectStatementContext {
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

func (s *ObjectBodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ObjectBodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ObjectBodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterObjectBody(s)
	}
}

func (s *ObjectBodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitObjectBody(s)
	}
}

func NewEmptyObjectBodyContext() *ObjectBodyContext {
	var p = new(ObjectBodyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectBody
	return p
}

func InitEmptyObjectBodyContext(p *ObjectBodyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_objectBody
}

func NewObjectBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ObjectBodyContext {
	var p = new(ObjectBodyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_objectBody

	return p
}
