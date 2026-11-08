// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ITypeopargContext is an interface to support dynamic dispatch.
type ITypeopargContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Ctypedecl() ICtypedeclContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode

	// IsTypeopargContext differentiates from other interfaces.
	IsTypeopargContext()
}

type TypeopargContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*TypeopargContext) IsTypeopargContext() {}

func (s *TypeopargContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeopargContext) Ctypedecl() ICtypedeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICtypedeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICtypedeclContext)
}

func (s *TypeopargContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *TypeopargContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *TypeopargContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeopargContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeopargContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterTypeoparg(s)
	}
}

func (s *TypeopargContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitTypeoparg(s)
	}
}

func NewEmptyTypeopargContext() *TypeopargContext {
	var p = new(TypeopargContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typeoparg
	return p
}

func InitEmptyTypeopargContext(p *TypeopargContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_typeoparg
}

func NewTypeopargContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeopargContext {
	var p = new(TypeopargContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_typeoparg

	return p
}
