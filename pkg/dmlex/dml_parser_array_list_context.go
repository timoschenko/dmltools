// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IArray_listContext is an interface to support dynamic dispatch.
type IArray_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLBRACK() []antlr.TerminalNode
	LBRACK(i int) antlr.TerminalNode
	AllArraydef() []IArraydefContext
	Arraydef(i int) IArraydefContext
	AllRBRACK() []antlr.TerminalNode
	RBRACK(i int) antlr.TerminalNode

	// IsArray_listContext differentiates from other interfaces.
	IsArray_listContext()
}

type Array_listContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*Array_listContext) IsArray_listContext() {}

func (s *Array_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Array_listContext) AllLBRACK() []antlr.TerminalNode {
	return s.GetTokens(dmlParserLBRACK)
}

func (s *Array_listContext) LBRACK(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACK, i)
}

func (s *Array_listContext) AllArraydef() []IArraydefContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IArraydefContext); ok {
			len++
		}
	}

	tst := make([]IArraydefContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IArraydefContext); ok {
			tst[i] = t.(IArraydefContext)
			i++
		}
	}

	return tst
}

func (s *Array_listContext) Arraydef(i int) IArraydefContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArraydefContext); ok {
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

	return t.(IArraydefContext)
}

func (s *Array_listContext) AllRBRACK() []antlr.TerminalNode {
	return s.GetTokens(dmlParserRBRACK)
}

func (s *Array_listContext) RBRACK(i int) antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACK, i)
}

func (s *Array_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Array_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Array_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterArray_list(s)
	}
}

func (s *Array_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitArray_list(s)
	}
}

func NewEmptyArray_listContext() *Array_listContext {
	var p = new(Array_listContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_array_list
	return p
}

func InitEmptyArray_listContext(p *Array_listContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_array_list
}

func NewArray_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Array_listContext {
	var p = new(Array_listContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_array_list

	return p
}
