// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IHook_declContext is an interface to support dynamic dispatch.
type IHook_declContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	HOOK() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	Ident() IIdentContext
	SEMI() antlr.TerminalNode
	Cdecl_list() ICdecl_listContext
	AllLBRACK() []antlr.TerminalNode
	LBRACK(i int) antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllRBRACK() []antlr.TerminalNode
	RBRACK(i int) antlr.TerminalNode

	// IsHook_declContext differentiates from other interfaces.
	IsHook_declContext()
}

type Hook_declContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Hook_declContext) IsHook_declContext() {}

func (s *Hook_declContext) GetParser() antlr.Parser { return s.parser }

func (s *Hook_declContext) HOOK() antlr.TerminalNode {
	return s.GetToken(dmlParserHOOK, 0)
}

func (s *Hook_declContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *Hook_declContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *Hook_declContext) Ident() IIdentContext {
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

func (s *Hook_declContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *Hook_declContext) Cdecl_list() ICdecl_listContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdecl_listContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICdecl_listContext)
}

func (s *Hook_declContext) AllLBRACK() []antlr.TerminalNode {
	return s.GetTokens(dmlParserLBRACK)
}

func (s *Hook_declContext) LBRACK(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACK, i)
}

func (s *Hook_declContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *Hook_declContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
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

	return t.(IExpressionContext)
}

func (s *Hook_declContext) AllRBRACK() []antlr.TerminalNode {
	return s.GetTokens(dmlParserRBRACK)
}

func (s *Hook_declContext) RBRACK(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACK, i)
}

func (s *Hook_declContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Hook_declContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Hook_declContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterHook_decl(s)
	}
}

func (s *Hook_declContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitHook_decl(s)
	}
}

func NewEmptyHook_declContext() *Hook_declContext {
	var p = new(Hook_declContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_hook_decl
	return p
}

func InitEmptyHook_declContext(p *Hook_declContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_hook_decl
}

func NewHook_declContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Hook_declContext {
	var p = new(Hook_declContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_hook_decl

	return p
}
