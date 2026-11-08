// Code generated from ./dml.g4 by ANTLR 4.13.2. DO NOT EDIT.

package dmlex

import (
	"github.com/antlr4-go/antlr/v4"
)

// IAssignopContext is an interface to support dynamic dispatch.
type IAssignopContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	ADD_ASSIGN() antlr.TerminalNode
	SUB_ASSIGN() antlr.TerminalNode
	MUL_ASSIGN() antlr.TerminalNode
	DIV_ASSIGN() antlr.TerminalNode
	MOD_ASSIGN() antlr.TerminalNode
	OR_ASSIGN() antlr.TerminalNode
	AND_ASSIGN() antlr.TerminalNode
	XOR_ASSIGN() antlr.TerminalNode
	LSHIFT_ASSIGN() antlr.TerminalNode
	RSHIFT_ASSIGN() antlr.TerminalNode

	// IsAssignopContext differentiates from other interfaces.
	IsAssignopContext()
}

type AssignopContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func (*AssignopContext) IsAssignopContext() {}

func (s *AssignopContext) GetParser() antlr.Parser { return s.parser }

func (s *AssignopContext) AllExpression() []IExpressionContext {
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

func (s *AssignopContext) Expression(i int) IExpressionContext {
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

func (s *AssignopContext) ADD_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserADD_ASSIGN, 0)
}

func (s *AssignopContext) SUB_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserSUB_ASSIGN, 0)
}

func (s *AssignopContext) MUL_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserMUL_ASSIGN, 0)
}

func (s *AssignopContext) DIV_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserDIV_ASSIGN, 0)
}

func (s *AssignopContext) MOD_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserMOD_ASSIGN, 0)
}

func (s *AssignopContext) OR_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserOR_ASSIGN, 0)
}

func (s *AssignopContext) AND_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserAND_ASSIGN, 0)
}

func (s *AssignopContext) XOR_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserXOR_ASSIGN, 0)
}

func (s *AssignopContext) LSHIFT_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserLSHIFT_ASSIGN, 0)
}

func (s *AssignopContext) RSHIFT_ASSIGN() antlr.TerminalNode {
	return s.GetToken(dmlParserRSHIFT_ASSIGN, 0)
}

func (s *AssignopContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssignopContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssignopContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.EnterAssignop(s)
	}
}

func (s *AssignopContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(dmlListener); ok {
		listenerT.ExitAssignop(s)
	}
}

func NewEmptyAssignopContext() *AssignopContext {
	var p = new(AssignopContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_assignop
	return p
}

func InitEmptyAssignopContext(p *AssignopContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = dmlParserRULE_assignop
}

func NewAssignopContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignopContext {
	var p = new(AssignopContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = dmlParserRULE_assignop

	return p
}
