// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IParamContext is an interface to support dynamic dispatch.
type IParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PARAM() antlr.TerminalNode
	Objident() IObjidentContext
	SEMI() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Ctypedecl() ICtypedeclContext
	Expression() IExpressionContext
	NEW_ASSIGN() antlr.TerminalNode
	DEFAULT_NEW() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	DEFAULT() antlr.TerminalNode
	AUTO() antlr.TerminalNode

	// IsParamContext differentiates from other interfaces.
	IsParamContext()
}

type ParamContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*ParamContext) IsParamContext() {}

func (s *ParamContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamContext) PARAM() antlr.TerminalNode {
	return s.GetToken(dmlParserPARAM, 0)
}

func (s *ParamContext) Objident() IObjidentContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjidentContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IObjidentContext)
}

func (s *ParamContext) SEMI() antlr.TerminalNode {
	return s.GetToken(dmlParserSEMI, 0)
}

func (s *ParamContext) COLON() antlr.TerminalNode {
	return s.GetToken(dmlParserCOLON, 0)
}

func (s *ParamContext) Ctypedecl() ICtypedeclContext {
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

func (s *ParamContext) Expression() IExpressionContext {
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

func (s *ParamContext) NEW_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserNEW_ASSIGN, 0)
}

func (s *ParamContext) DEFAULT_NEW() antlr.TerminalNode {
	return s.GetToken(dmlParserDEFAULT_NEW, 0)
}

func (s *ParamContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserASSIGN, 0)
}

func (s *ParamContext) DEFAULT() antlr.TerminalNode {
	return s.GetToken(dmlParserDEFAULT, 0)
}

func (s *ParamContext) AUTO() antlr.TerminalNode {
	return s.GetToken(dmlParserAUTO, 0)
}

func (s *ParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterParam(s)
	}
}

func (s *ParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitParam(s)
	}
}

func NewEmptyParamContext() *ParamContext {
	var p = new(ParamContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_param
	return p
}

func InitEmptyParamContext(p *ParamContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_param
}

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext {
	var p = new(ParamContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_param

	return p
}
