// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// ICdecl_listContext is an interface to support dynamic dispatch.
type ICdecl_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllCdecl() []ICdeclContext
	Cdecl(i int) ICdeclContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsCdecl_listContext differentiates from other interfaces.
	IsCdecl_listContext()
}

type Cdecl_listContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Cdecl_listContext) IsCdecl_listContext() {}

func (s *Cdecl_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Cdecl_listContext) AllCdecl() []ICdeclContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICdeclContext); ok {
			len++
		}
	}

	tst := make([]ICdeclContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICdeclContext); ok {
			tst[i] = t.(ICdeclContext)
			i++
		}
	}

	return tst
}

func (s *Cdecl_listContext) Cdecl(i int) ICdeclContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICdeclContext); ok {
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

	return t.(ICdeclContext)
}

func (s *Cdecl_listContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCOMMA)
}

func (s *Cdecl_listContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, i)
}

func (s *Cdecl_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Cdecl_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Cdecl_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterCdecl_list(s)
	}
}

func (s *Cdecl_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitCdecl_list(s)
	}
}

func NewEmptyCdecl_listContext() *Cdecl_listContext {
	var p = new(Cdecl_listContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl_list
	return p
}

func InitEmptyCdecl_listContext(p *Cdecl_listContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_cdecl_list
}

func NewCdecl_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Cdecl_listContext {
	var p = new(Cdecl_listContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_cdecl_list

	return p
}
