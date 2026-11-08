// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICdecl_list_variadicContext is an interface to support dynamic dispatch.
type ICdecl_list_variadicContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ELLIPSIS() antlr.TerminalNode
	Cdecl_list() ICdecl_listContext
	COMMA() antlr.TerminalNode

	// IsCdecl_list_variadicContext differentiates from other interfaces.
	IsCdecl_list_variadicContext()
}

type Cdecl_list_variadicContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Cdecl_list_variadicContext) IsCdecl_list_variadicContext() {}

func (s *Cdecl_list_variadicContext) GetParser() antlr.Parser { return s.parser }

func (s *Cdecl_list_variadicContext) ELLIPSIS() antlr.TerminalNode {
	return s.GetToken(dmlParserELLIPSIS, 0)
}

func (s *Cdecl_list_variadicContext) Cdecl_list() ICdecl_listContext {
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

func (s *Cdecl_list_variadicContext) COMMA() antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, 0)
}

func (s *Cdecl_list_variadicContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Cdecl_list_variadicContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Cdecl_list_variadicContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdecl_list_variadic(s)
	}
}

func (s *Cdecl_list_variadicContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdecl_list_variadic(s)
	}
}

func NewEmptyCdecl_list_variadicContext() *Cdecl_list_variadicContext {
	var p = new(Cdecl_list_variadicContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl_list_variadic
	return p
}

func InitEmptyCdecl_list_variadicContext(p *Cdecl_list_variadicContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl_list_variadic
}

func NewCdecl_list_variadicContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Cdecl_list_variadicContext {
	var p = new(Cdecl_list_variadicContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_cdecl_list_variadic

	return p
}
