// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IBitSliceContext is an interface to support dynamic dispatch.
type IBitSliceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	COLON() antlr.TerminalNode
	Expression() IExpressionContext
	Endianflag() IEndianflagContext

	// IsBitSliceContext differentiates from other interfaces.
	IsBitSliceContext()
}

type BitSliceContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*BitSliceContext) IsBitSliceContext() {}

func (s *BitSliceContext) GetParser() antlr.Parser { return s.parser }

func (s *BitSliceContext) COLON() antlr.TerminalNode {
	return s.GetToken(dmlParserCOLON, 0)
}

func (s *BitSliceContext) Expression() IExpressionContext {
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

func (s *BitSliceContext) Endianflag() IEndianflagContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEndianflagContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEndianflagContext)
}

func (s *BitSliceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BitSliceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BitSliceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterBitSlice(s)
	}
}

func (s *BitSliceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitBitSlice(s)
	}
}

func NewEmptyBitSliceContext() *BitSliceContext {
	var p = new(BitSliceContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitSlice
	return p
}

func InitEmptyBitSliceContext(p *BitSliceContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitSlice
}

func NewBitSliceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BitSliceContext {
	var p = new(BitSliceContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_bitSlice

	return p
}
