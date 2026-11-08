// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IMethodOutParamsContext is an interface to support dynamic dispatch.
type IMethodOutParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ARROW() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	Cdecl_list() ICdecl_listContext

	// IsMethodOutParamsContext differentiates from other interfaces.
	IsMethodOutParamsContext()
}

type MethodOutParamsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*MethodOutParamsContext) IsMethodOutParamsContext() {}

func (s *MethodOutParamsContext) GetParser() antlr.Parser { return s.parser }

func (s *MethodOutParamsContext) ARROW() antlr.TerminalNode {
	return s.GetToken(dmlParserARROW, 0)
}

func (s *MethodOutParamsContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *MethodOutParamsContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *MethodOutParamsContext) Cdecl_list() ICdecl_listContext {
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

func (s *MethodOutParamsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MethodOutParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MethodOutParamsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterMethodOutParams(s)
	}
}

func (s *MethodOutParamsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitMethodOutParams(s)
	}
}

func NewEmptyMethodOutParamsContext() *MethodOutParamsContext {
	var p = new(MethodOutParamsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_methodOutParams
	return p
}

func InitEmptyMethodOutParamsContext(p *MethodOutParamsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_methodOutParams
}

func NewMethodOutParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MethodOutParamsContext {
	var p = new(MethodOutParamsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_methodOutParams

	return p
}
