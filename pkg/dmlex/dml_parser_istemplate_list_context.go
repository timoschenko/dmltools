// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IIstemplateListContext is an interface to support dynamic dispatch.
type IIstemplateListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllObjident() []IObjidentContext
	Objident(i int) IObjidentContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsIstemplateListContext differentiates from other interfaces.
	IsIstemplateListContext()
}

type IstemplateListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*IstemplateListContext) IsIstemplateListContext() {}

func (s *IstemplateListContext) GetParser() antlr.Parser { return s.parser }

func (s *IstemplateListContext) AllObjident() []IObjidentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IObjidentContext); ok {
			len++
		}
	}

	tst := make([]IObjidentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IObjidentContext); ok {
			tst[i] = t.(IObjidentContext)
			i++
		}
	}

	return tst
}

func (s *IstemplateListContext) Objident(i int) IObjidentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IObjidentContext); ok {
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

	return t.(IObjidentContext)
}

func (s *IstemplateListContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserLPAREN, 0)
}

func (s *IstemplateListContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(dmlParserRPAREN, 0)
}

func (s *IstemplateListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(dmlParserCOMMA)
}

func (s *IstemplateListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserCOMMA, i)
}

func (s *IstemplateListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IstemplateListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IstemplateListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterIstemplateList(s)
	}
}

func (s *IstemplateListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitIstemplateList(s)
	}
}

func NewEmptyIstemplateListContext() *IstemplateListContext {
	var p = new(IstemplateListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_istemplateList
	return p
}

func InitEmptyIstemplateListContext(p *IstemplateListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_istemplateList
}

func NewIstemplateListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IstemplateListContext {
	var p = new(IstemplateListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_istemplateList

	return p
}
