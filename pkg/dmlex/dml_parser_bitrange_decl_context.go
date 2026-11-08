// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IBitrangeDeclContext is an interface to support dynamic dispatch.
type IBitrangeDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AT() antlr.TerminalNode
	LBRACK() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	RBRACK() antlr.TerminalNode
	COLON() antlr.TerminalNode

	// IsBitrangeDeclContext differentiates from other interfaces.
	IsBitrangeDeclContext()
}

type BitrangeDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*BitrangeDeclContext) IsBitrangeDeclContext() {}

func (s *BitrangeDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *BitrangeDeclContext) AT() antlr.TerminalNode {
	return s.GetToken(dmlParserAT, 0)
}

func (s *BitrangeDeclContext) LBRACK() antlr.TerminalNode {
	return s.GetToken(dmlParserLBRACK, 0)
}

func (s *BitrangeDeclContext) AllExpression() []IExpressionContext {
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

func (s *BitrangeDeclContext) Expression(i int) IExpressionContext {
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

func (s *BitrangeDeclContext) RBRACK() antlr.TerminalNode {
	return s.GetToken(dmlParserRBRACK, 0)
}

func (s *BitrangeDeclContext) COLON() antlr.TerminalNode {
	return s.GetToken(dmlParserCOLON, 0)
}

func (s *BitrangeDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BitrangeDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BitrangeDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterBitrangeDecl(s)
	}
}

func (s *BitrangeDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitBitrangeDecl(s)
	}
}

func NewEmptyBitrangeDeclContext() *BitrangeDeclContext {
	var p = new(BitrangeDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitrangeDecl
	return p
}

func InitEmptyBitrangeDeclContext(p *BitrangeDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_bitrangeDecl
}

func NewBitrangeDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BitrangeDeclContext {
	var p = new(BitrangeDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_bitrangeDecl

	return p
}
